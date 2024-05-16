package vgin

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/slices"
)

const (
	ParamsJsonTag      = "vjson"
	ParamsMultiFormTag = "vform"
	ParamsQueryTag     = "vquery"
	ParamsPathTag      = "vpath"
	ParamsHeaderTag    = "vheader"
)

type WrapInHandler interface {
	IsolatedHandler
	OriginHandler() Handler
}

type paramsInHandler struct {
	Into IsolatedHandler
}

func (ph *paramsInHandler) Name() string {
	return guessHandlerName(ph.Into)
}

func (ph *paramsInHandler) InitHandler() IsolatedHandler {
	return &paramsInHandler{ph.Into.InitHandler()}
}

func (ph *paramsInHandler) OriginHandler() Handler {
	return ph.Into.InitHandler()
}

func (ph *paramsInHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	if err := ParseMapParams(ctx, c, ph.Into); err != nil {
		return &Ret{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Err:     err,
			Message: "Parse params error",
		}
	}

	return ph.Into.HandleFunc(ctx, c)
}

func ParamsIn(handler IsolatedHandler) Handler {
	if checkIsWebsocket(handler) {
		lg.Fatal("Websocket can not use `ParamsIn` inject handler")
	}
	return &paramsInHandler{handler}
}

var (
	pattern     = regexp.MustCompile(`(\w+):"[^"]+"`)
	tagMap      = make(map[string]slices.StringSet)
	tagMapMutex sync.RWMutex
)

func findHandlerTag(h Handler) slices.StringSet {
	handlerName := lg.StructName(h)
	tagMapMutex.RLock()
	if r, exists := tagMap[handlerName]; exists {
		tagMapMutex.RUnlock()
		return r
	}
	tagMapMutex.RUnlock()

	tagMapMutex.Lock()
	defer tagMapMutex.Unlock()
	if r, exists := tagMap[lg.StructName(h)]; exists {
		return r
	}

	uv := reflect.TypeOf(h)

	if uv.Kind() == reflect.Pointer {
		uv = uv.Elem()
	}

	numField := uv.NumField()
	if numField == 0 {
		return nil
	}

	tagSet := slices.NewStringSet()
	for idx := 0; idx < numField; idx++ {
		field := uv.Field(idx)
		ret := pattern.FindAllStringSubmatch(string(field.Tag), -1)
		for _, r := range ret {
			if len(r) != 2 {
				continue
			}
			tagSet.Add(r[1])
		}
	}

	tagMap[handlerName] = tagSet
	return tagSet
}

func ParseMapParams(ctx context.Context, c *Context, into Handler) (err error) {
	tagSet := findHandlerTag(into)
	switch c.ContentType() {
	case gin.MIMEJSON:
		if !tagSet.Contains(ParamsJsonTag) {
			break
		}
		var raw []byte
		raw, err = BodyRawData(c)
		if err != nil {
			return err
		}
		err = parseJson(ctx, raw, into)
	case gin.MIMEMultipartPOSTForm:
		if !tagSet.Contains(ParamsMultiFormTag) {
			break
		}
		var form *multipart.Form
		form, err = c.MultipartForm()
		if err != nil {
			return err
		}
		err = parseMultiForm(form.Value, into)
	}
	if err != nil {
		return errors.Wrap(err, "parse contentType data")
	}

	allwaysParse := []func(*Context, Handler) error{
		parseQuery(tagSet.Contains(ParamsQueryTag)),
		parsePath(tagSet.Contains(ParamsPathTag)),
		parseHeader(tagSet.Contains(ParamsHeaderTag)),
	}

	for _, parser := range allwaysParse {
		if err := parser(c, into); err != nil {
			return err
		}
	}

	return nil
}

func parseHeader(needDo bool) func(c *Context, intoHandler Handler) error {
	return func(c *Context, intoHandler Handler) error {
		if !needDo || len(c.Request.Header) == 0 {
			return nil
		}

		return mapFormByTag(intoHandler, c.Request.Header, ParamsHeaderTag)
	}
}

func parsePath(needDo bool) func(c *Context, intoHandler Handler) error {
	return func(c *Context, intoHandler Handler) error {
		if !needDo || len(c.Params) == 0 {
			return nil
		}

		tmp := make(map[string][]string)
		for _, p := range c.Params {
			tmp[p.Key] = []string{p.Value}
		}
		return mapFormByTag(intoHandler, tmp, ParamsPathTag)
	}
}

func parseQuery(needDo bool) func(c *Context, intoHandler Handler) error {
	return func(c *Context, intoHandler Handler) error {
		if !needDo {
			return nil
		}

		queryMap := c.Request.URL.Query()
		return mapFormByTag(intoHandler, queryMap, ParamsQueryTag)
	}
}

func mapstructureDecoder(tag string, result Handler) (*mapstructure.Decoder, error) {
	config := &mapstructure.DecoderConfig{
		Result:  result,
		TagName: tag,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	return decoder, nil
}

func parseJson(ctx context.Context, data []byte, intoHandler Handler) error {
	data = bytes.TrimSpace(data)
	tmp := make(map[string]any)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "decode json")
	}
	decoder, err := mapstructureDecoder(ParamsJsonTag, intoHandler)
	if err != nil {
		return errors.Wrap(err, "map decode")
	}

	err = decoder.Decode(tmp)
	if err != nil {
		if strings.Contains(err.Error(), "unconvertible type") {
			lg.Warnc(ctx, "Type error: %v", err)
		} else {
			return errors.Wrap(err, "map decode json")
		}
	}

	return nil
}

func parseMultiForm(data map[string][]string, intoHandler Handler) error {
	return mapFormByTag(intoHandler, data, ParamsMultiFormTag)
}
