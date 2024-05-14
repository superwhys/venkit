package vgin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

type jsonHandler struct {
	JsonDataStr     string  `vjson:"name" form:"name"`
	JsonDataInt     int     `vjson:"age" form:"age"`
	JsonDataFloat64 float64 `vjson:"money" form:"money"`
}

func (h *jsonHandler) Name() string {
	return "jsonHandler"
}

func (h *jsonHandler) InitHandler() IsolatedHandler {
	return &jsonHandler{}
}

func (h *jsonHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return &Ret{
		Code: 200,
		Data: h,
	}
}

func TestParseJson(t *testing.T) {
	r := NewWithEngine(gin.New(), gin.Recovery())
	r.POST("/json", ParamsIn(&jsonHandler{}))

	reqGetter := func(body string) *http.Request {
		req := httptest.NewRequest("POST", "/json", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		return req
	}

	respCheck := func(req *http.Request, expected *jsonHandler) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		ret := new(Ret)
		ret.Data = new(jsonHandler)
		err := json.Unmarshal(w.Body.Bytes(), ret)
		assert.Nil(t, err)
		assert.Equal(t, expected, ret.Data)
	}

	func() {
		t.Log("test parse json")
		reqBody := `{"name": "John Doe", "age": 18, "money": 123.456}`
		req := reqGetter(reqBody)
		expected := &jsonHandler{
			JsonDataStr:     "John Doe",
			JsonDataInt:     18,
			JsonDataFloat64: 123.456,
		}
		respCheck(req, expected)
	}()

	func() {
		t.Log("test parse json with type error -> (want int get string)")
		reqBody := `{"name": "John Doe", "age": "18", "money": 123.456}`
		req := reqGetter(reqBody)
		expected := &jsonHandler{
			JsonDataStr:     "John Doe",
			JsonDataInt:     0,
			JsonDataFloat64: 123.456,
		}
		respCheck(req, expected)
	}()

	func() {
		t.Log("test parse json with type error -> (want float64 get string)")
		reqBody := `{"name": "John Doe", "age": 18, "money": "123.456"}`
		req := reqGetter(reqBody)
		expected := &jsonHandler{
			JsonDataStr:     "John Doe",
			JsonDataInt:     18,
			JsonDataFloat64: 0,
		}
		respCheck(req, expected)
	}()

	func() {
		t.Log("test parse json with type error -> (want string get int)")
		reqBody := `{"name": 12345, "age": 18, "money": 123.456}`
		req := reqGetter(reqBody)
		expected := &jsonHandler{
			JsonDataStr:     "",
			JsonDataInt:     18,
			JsonDataFloat64: 123.456,
		}
		respCheck(req, expected)
	}()
}

type queryHandler struct {
	QueryInt    int    `vquery:"query1"`
	QueryString string `vquery:"query2"`
}

func (h *queryHandler) InitHandler() IsolatedHandler {
	return &queryHandler{}
}

func (h *queryHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return &Ret{
		Code: 200,
		Data: h,
	}
}

func TestParseQuery(t *testing.T) {
	r := NewWithEngine(gin.New(), gin.Recovery())
	r.GET("/query", ParamsIn(&queryHandler{}))

	reqGetter := func(v url.Values) *http.Request {
		req := httptest.NewRequest("GET", "/query", nil)
		req.URL.RawQuery = v.Encode()
		return req
	}

	respCheck := func(req *http.Request, expected *queryHandler, wantErr bool) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		ret := new(Ret)
		ret.Data = new(queryHandler)
		err := json.Unmarshal(w.Body.Bytes(), ret)
		assert.Nil(t, err)
		if wantErr {
			assert.NotEqual(t, 200, ret.Code)
		} else {
			assert.Equal(t, expected, ret.Data)
		}
	}

	func() {
		t.Log("test parse query")
		v := url.Values{}
		v.Add("query1", "123")
		v.Add("query2", "ssssss")
		req := reqGetter(v)
		expected := &queryHandler{QueryInt: 123, QueryString: "ssssss"}
		respCheck(req, expected, false)
	}()

	func() {
		t.Log("test parse query with type error -> (want int get string) ")
		v := url.Values{}
		v.Add("query1", "abc")
		v.Add("query2", "ssssss")
		req := reqGetter(v)
		respCheck(req, nil, true)
	}()
}

type headerHandler struct {
	HeaderInt    int    `vheader:"Header-Id"`
	HeaderString string `vheader:"Token-Id"`
}

func (h *headerHandler) InitHandler() IsolatedHandler {
	return &headerHandler{}
}

func (h *headerHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return &Ret{
		Code: 200,
		Data: h,
	}
}

func TestParseHeader(t *testing.T) {
	r := NewWithEngine(gin.New(), gin.Recovery())
	r.GET("/header", ParamsIn(&headerHandler{}))

	reqGetter := func(intV string, strV string) *http.Request {
		req := httptest.NewRequest("GET", "/header", nil)
		req.Header.Set("Header-id", intV)
		req.Header.Set("Token-id", strV)

		return req
	}

	respCheck := func(req *http.Request, expected *headerHandler, wantErr bool) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		ret := new(Ret)
		ret.Data = new(headerHandler)
		err := json.Unmarshal(w.Body.Bytes(), ret)
		assert.Nil(t, err)
		if wantErr {
			assert.NotEqual(t, 200, ret.Code)
		} else {
			assert.Equal(t, expected, ret.Data)
		}
	}

	func() {
		t.Log("test parse header")
		req := reqGetter("123", "abc")
		expected := &headerHandler{HeaderInt: 123, HeaderString: "abc"}
		respCheck(req, expected, false)
	}()

	func() {
		t.Log("test parse header with type Error")
		req := reqGetter("1a1", "abc")
		respCheck(req, nil, true)
	}()
}

type pathHandler struct {
	PathInt    int    `vpath:"int_val"`
	PathString string `vpath:"str_val"`
}

func (h *pathHandler) InitHandler() IsolatedHandler {
	return &pathHandler{}
}

func (h *pathHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return &Ret{
		Code: 200,
		Data: h,
	}
}

func TestParsePath(t *testing.T) {
	r := NewWithEngine(gin.New(), gin.Recovery())
	r.GET("/path/:int_val/:str_val", ParamsIn(&pathHandler{}))

	reqGetter := func(intV string, strV string) *http.Request {
		req := httptest.NewRequest("GET", fmt.Sprintf("/path/%v/%v", intV, strV), nil)
		return req
	}

	respCheck := func(req *http.Request, expected *pathHandler, wantErr bool) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		ret := new(Ret)
		ret.Data = new(pathHandler)
		err := json.Unmarshal(w.Body.Bytes(), ret)
		assert.Nil(t, err)
		if wantErr {
			assert.NotEqual(t, 200, ret.Code)
		} else {
			assert.Equal(t, expected, ret.Data)
		}
	}

	func() {
		t.Log("test parse path")
		req := reqGetter("123", "abc")
		expected := &pathHandler{PathInt: 123, PathString: "abc"}
		respCheck(req, expected, false)
	}()

	func() {
		t.Log("test parse path with type Error")
		req := reqGetter("1a1", "abc")
		respCheck(req, nil, true)
	}()
}

func BenchmarkVgin(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery(), BodyBufferMiddleware())
	r.POST("/ping", ParamsIn(&jsonHandler{}))

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18}`
		req := httptest.NewRequest("POST", "/ping", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkVginWithMiddlewareAndPoolOptimize(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery(), BodyBufferMiddleware())
	r.POST("/json", ParamsIn(&jsonHandler{}), ParamsIn(&jsonHandler{}))

	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18, "money": 123.45}`
		req := httptest.NewRequest("POST", "/json", strings.NewReader(reqBody))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkOriginGin(b *testing.B) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/ping", func(ctx *gin.Context) {
		data := &jsonHandler{}
		if err := ctx.ShouldBind(data); err != nil {
			fmt.Println(err)
			ctx.JSON(400, "parse data error")
			return
		}
	})

	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18}`
		req := httptest.NewRequest("POST", "/ping", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

type SimpleHandler struct{}

func (h *SimpleHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return SuccessRet("success")
}

func BenchmarkVginSimple(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery())

	r.POST("/ping", &SimpleHandler{})

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/ping", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGinSimple(b *testing.B) {
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/ping", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

type oneJSonHandler struct {
	User string `vjson:"user" form:"user"`
}

func (h *oneJSonHandler) InitHandler() IsolatedHandler {
	return &oneJSonHandler{}
}

func (h *oneJSonHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return SuccessRet(h)
}

func BenchmarkVginOneJsonParams(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery())

	r.POST("/user", ParamsIn(&oneJSonHandler{}))

	for i := 0; i < b.N; i++ {
		reqBody := `{"user": "hoven"}`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGinOneJsonParams(b *testing.B) {
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/user", func(ctx *gin.Context) {
		data := new(oneJSonHandler)
		if err := ctx.ShouldBind(data); err != nil {
			ctx.JSON(400, "failed")
			return
		}
		ctx.JSON(200, data)
	})

	for i := 0; i < b.N; i++ {
		reqBody := `{"user": "hoven"}`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

type oneQueryHandler struct {
	User string `vquery:"user" form:"user"`
}

func (h *oneQueryHandler) InitHandler() IsolatedHandler {
	return &oneQueryHandler{}
}

func (h *oneQueryHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	// data := new(oneJSonHandler)
	// if err := c.ShouldBindJSON(data); err != nil {
	// 	return ErrorRet(400, err, err.Error())
	// }
	return SuccessRet(h)
}

func BenchmarkVginOneQuery(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery())

	r.POST("/user", ParamsIn(&oneQueryHandler{}))

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/user?user=hoven", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGinOneQuery(b *testing.B) {
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/user", func(ctx *gin.Context) {
		data := new(oneQueryHandler)
		if err := ctx.ShouldBind(data); err != nil {
			ctx.JSON(400, "failed")
			return
		}
		ctx.JSON(200, data)
	})

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/user?user=hoven", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
