package vgin

import (
	"context"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"sync"
	
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/slices/v2"
	v1Vgin "github.com/superwhys/venkit/vgin"
)

const (
	handlerFormatInvalid = "handler supports only gin.HandlerFunc, v1Vgin.Handler, and func(ctx context.Context, c *vgin.Context, data *DataStruct) vgin.HandleResponse"
)

type Context struct {
	*gin.Context
}

// Handler must be a format of func(ctx context.Context, c *vgin.Context, data *DataStruct) vgin.HandleResponse
type Handler any

func WrapHandler(ctx context.Context, handlers ...Handler) []gin.HandlerFunc {
	handlerFuncs := make([]gin.HandlerFunc, 0, len(handlers))
	for _, handler := range handlers {
		switch h := handler.(type) {
		case gin.HandlerFunc:
			handlerFuncs = append(handlerFuncs, h)
		case func(*gin.Context):
			handlerFuncs = append(handlerFuncs, h)
		case v1Vgin.Handler:
			handlerFuncs = append(handlerFuncs, v1Vgin.WrapHandler(ctx, h)...)
		default:
			handlerFuncs = append(handlerFuncs, wrapHandler(ctx, h))
		}
	}
	
	return handlerFuncs
}

func wrapHandler(ctx context.Context, handler Handler) gin.HandlerFunc {
	fv := reflect.ValueOf(handler)
	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		lg.Fatal(handlerFormatInvalid)
	}
	
	ctx = lg.With(ctx, "[%s]", runtime.FuncForPC(fv.Pointer()).Name())
	
	if ft.NumIn() < 2 {
		lg.Fatal(handlerFormatInvalid)
	}
	
	if ft.NumOut() != 1 {
		lg.Fatal(handlerFormatInvalid)
	}
	
	if !reflect.TypeOf((*HandleResponse)(nil)).Elem().Implements(ft.Out(0)) {
		lg.Fatal(handlerFormatInvalid)
	}
	
	return wrapHandlerFunc(ctx, fv, ft)
}

func wrapHandlerFunc(ctx context.Context, fv reflect.Value, ft reflect.Type) gin.HandlerFunc {
	funcParamsNum := ft.NumIn()
	argsPool := sync.Pool{
		New: func() any {
			args := make([]reflect.Value, funcParamsNum)
			args[0] = reflect.ValueOf(ctx)
			return &args
		},
	}
	
	return func(c *gin.Context) {
		args := *(argsPool.Get().(*[]reflect.Value))
		defer func() {
			for i := 1; i < funcParamsNum; i++ {
				args[i] = reflect.Value{}
			}
			argsPool.Put(&args)
		}()
		
		vc := &Context{Context: c}
		
		args[1] = reflect.ValueOf(vc)
		prepareParams(ctx, vc, ft, funcParamsNum, args)
		
		responses := fv.Call(args)
		var ret HandleResponse
		if r := responses[0].Interface(); r != nil {
			ret = r.(HandleResponse)
		}
		if checkRet(ctx, vc, ret) {
			return
		}
		
		if c.IsAborted() {
			return
		}
		
		if ret != nil {
			ReturnWithStatus(c, ret.GetCode(), ret.GetData())
		}
	}
}

func prepareParams(ctx context.Context, vc *Context, ft reflect.Type, funcParamsNum int, args []reflect.Value) {
	if funcParamsNum == 2 {
		return
	}
	
	for i := 2; i < funcParamsNum; i++ {
		params := ft.In(i)
		if params.Kind() == reflect.Ptr {
			params = params.Elem()
		}
		
		if params.Kind() != reflect.Struct {
			continue
		}
		
		paramsValue := reflect.New(params)
		tags := findStructTag(params)
		if tags.Length() != 0 {
			err := parseParams(ctx, vc, tags, paramsValue.Interface())
			if err != nil {
				lg.Errorc(ctx, "params params: %v error: %v", params, err)
				continue
			}
		}
		
		args[i] = paramsValue
	}
}

func checkRet(ctx context.Context, c *Context, ret HandleResponse) (hasErr bool) {
	if ret == nil {
		return
	}
	
	if ret.GetCode() != 200 && ret.GetError() != nil {
		lg.Errorc(ctx, "handle err: %v", ret.GetError())
		AbortWithError(c.Context, ret.GetCode(), ret.GetMessage())
		hasErr = true
	} else if ret.GetCode() != 200 {
		AbortWithError(c.Context, ret.GetCode(), ret.GetMessage())
		hasErr = true
	}
	return
}

var (
	pattern     = regexp.MustCompile(`(\w+):"[^"]+"`)
	tagMap      = make(map[string]slices.StringSet)
	tagMapMutex sync.RWMutex
)

func findStructTag(t reflect.Type) slices.StringSet {
	numField := t.NumField()
	if numField == 0 {
		return nil
	}
	structName := t.String()
	
	tagMapMutex.RLock()
	if r, exists := tagMap[structName]; exists {
		tagMapMutex.RUnlock()
		return r
	}
	tagMapMutex.RUnlock()
	
	tagMapMutex.Lock()
	defer tagMapMutex.Unlock()
	if r, exists := tagMap[structName]; exists {
		return r
	}
	
	tags := slices.NewStringSet()
	for idx := 0; idx < numField; idx++ {
		field := t.Field(idx)
		ret := pattern.FindAllStringSubmatch(string(field.Tag), -1)
		for _, r := range ret {
			if len(r) != 2 {
				continue
			}
			tags.Add(r[1])
		}
	}
	
	tagMap[structName] = tags
	return tags
}

const (
	ParamsJsonTag      = "vjson"
	ParamsMultiFormTag = "vform"
	ParamsQueryTag     = "vquery"
	ParamsPathTag      = "vpath"
	ParamsHeaderTag    = "vheader"
	
	defaultMemory = 32 << 20
)

func parseParams(ctx context.Context, c *Context, tags slices.StringSet, params any) (err error) {
	switch c.ContentType() {
	case gin.MIMEJSON:
		if !tags.Contains(ParamsJsonTag) {
			break
		}
		var raw []byte
		raw, err = BodyRawData(c)
		if err != nil {
			break
		}
		err = parseJson(ctx, raw, params)
	case gin.MIMEMultipartPOSTForm:
		if !tags.Contains(ParamsMultiFormTag) {
			break
		}
		var form *multipart.Form
		form, err = c.MultipartForm()
		if err != nil {
			break
		}
		err = parseMultiForm(form.Value, params)
	case gin.MIMEPOSTForm:
		if !tags.Contains(ParamsMultiFormTag) {
			break
		}
		
		if err = c.Request.ParseForm(); err != nil {
			break
		}
		if err = c.Request.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
			break
		}
		err = parseMultiForm(c.Request.Form, params)
	}
	if err != nil {
		return errors.Wrap(err, "parse contentType data")
	}
	
	alwaysParse := []func(*Context, any) error{
		parseQuery(tags.Contains(ParamsQueryTag)),
		parsePath(tags.Contains(ParamsPathTag)),
		parseHeader(tags.Contains(ParamsHeaderTag)),
	}
	
	for _, parser := range alwaysParse {
		if err := parser(c, params); err != nil {
			return err
		}
	}
	return nil
}
