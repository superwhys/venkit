package vgin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type PingHandler struct {
	JsonDataStr     string  `vjson:"name" form:"name"`
	JsonDataInt     int     `vjson:"age" form:"age"`
	JsonDataFloat64 float64 `vjson:"money" form:"money"`
	QueryDataStr    string  `vquery:"query1"`
	QueryDataInt    int     `vquery:"query2"`
	PathData        int     `vpath:"test_id"`
	HeaderDataStr   string  `vheader:"Token"`
	HeaderDataInt   int     `vheader:"header_id" `
}

func (h *PingHandler) Name() string {
	return "PingHandler"
}

func (h *PingHandler) InitHandler() IsolatedHandler {
	return &PingHandler{}
}

func (h *PingHandler) HandleFunc(ctx context.Context, c *gin.Context) HandleResponse {
	return &Ret{
		Code: 200,
		Data: h,
	}
}

func TestParseJson(t *testing.T) {
	r := NewWithEngine(gin.New(), gin.Recovery())
	r.POST("/ping/:test_id", ParamsIn(&PingHandler{}))

	reqGetter := func(body string) *http.Request {
		req := httptest.NewRequest("POST", fmt.Sprintf("/ping/%v?query1=asd&query2=123", "12"), strings.NewReader(body))
		req.Header.Set("Token", "asdfasdfasfdasdfasdfadsfasdfasdfasddf")
		req.Header.Set("header_id", "123123123")
		req.Header.Set("Content-Type", "application/json")

		return req
	}

	respCheck := func(req *http.Request, expected string) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, expected, w.Body.String())
	}

	func() {
		t.Log("test parse json")
		reqBody := `{"name": "John Doe", "age": 18, "money": 123.456}`
		req := reqGetter(reqBody)
		expected := `{"code":200,"data":{"JsonDataStr":"John Doe","JsonDataInt":18,"JsonDataFloat64":123.456,"QueryDataStr":"asd","QueryDataInt":123,"PathData":12,"HeaderDataStr":"asdfasdfasfdasdfasdfadsfasdfasdfasddf","HeaderDataInt":123123123}}`
		respCheck(req, expected)
	}()

	func() {
		t.Log("test parse json with type error -> (want int get string)")
		reqBody := `{"name": "John Doe", "age": "18", "money": 123.456}`
		req := reqGetter(reqBody)
		expected := `{"code":200,"data":{"JsonDataStr":"John Doe","JsonDataInt":0,"JsonDataFloat64":123.456,"QueryDataStr":"asd","QueryDataInt":123,"PathData":12,"HeaderDataStr":"asdfasdfasfdasdfasdfadsfasdfasdfasddf","HeaderDataInt":123123123}}`
		respCheck(req, expected)
	}()

	func() {
		t.Log("test parse json with type error -> (want float64 get string)")
		reqBody := `{"name": "John Doe", "age": 18, "money": "123.456"}`
		req := reqGetter(reqBody)
		expected := `{"code":200,"data":{"JsonDataStr":"John Doe","JsonDataInt":18,"JsonDataFloat64":0,"QueryDataStr":"asd","QueryDataInt":123,"PathData":12,"HeaderDataStr":"asdfasdfasfdasdfasdfadsfasdfasdfasddf","HeaderDataInt":123123123}}`
		respCheck(req, expected)
	}()

	func() {
		t.Log("test parse json with type error -> (want string get int)")
		reqBody := `{"name": 12345, "age": 18, "money": 123.456}`
		req := reqGetter(reqBody)
		expected := `{"code":200,"data":{"JsonDataStr":"","JsonDataInt":18,"JsonDataFloat64":123.456,"QueryDataStr":"asd","QueryDataInt":123,"PathData":12,"HeaderDataStr":"asdfasdfasfdasdfasdfadsfasdfasdfasddf","HeaderDataInt":123123123}}`
		respCheck(req, expected)
	}()

}

func BenchmarkVgin(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery(), BodyBufferMiddleware())
	r.POST("/ping/:test_id", ParamsIn(&PingHandler{}))

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18}`
		req := httptest.NewRequest("POST", fmt.Sprintf("/ping/%v?query1=asd&query2=123", "12"), strings.NewReader(reqBody))
		req.Header.Set("Token", "asdfasdfasfdasdfasdfadsfasdfasdfasddf")
		req.Header.Set("header_id", "123123123")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkVginWithMiddlewareAndPoolOptimize(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery(), BodyBufferMiddleware())
	r.POST("/ping/:test_id", ParamsIn(&PingHandler{}), ParamsIn(&PingHandler{}))

	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18}`
		req := httptest.NewRequest("POST", fmt.Sprintf("/ping/%v?query1=asd&query2=123", "12"), strings.NewReader(reqBody))
		req.Header.Set("Token", "asdfasdfasfdasdfasdfadsfasdfasdfasddf")
		req.Header.Set("header_id", "123123123")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkOriginGin(b *testing.B) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/ping/:test_id", func(ctx *gin.Context) {
		data := &PingHandler{}
		if err := ctx.ShouldBind(data); err != nil {
			fmt.Println(err)
			ctx.JSON(400, "parse data error")
			return
		}

		testId, _ := strconv.ParseInt(ctx.Param("test_id"), 10, 64)
		data.PathData = int(testId)
		data.QueryDataStr = ctx.Query("query1")
		query2Int, _ := strconv.ParseInt(ctx.Query("query2"), 10, 64)
		data.QueryDataInt = int(query2Int)
		data.HeaderDataStr = ctx.GetHeader("Token")
		headerId, _ := strconv.ParseInt(ctx.GetHeader("header_id"), 10, 64)
		data.HeaderDataInt = int(headerId)
	})

	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18}`
		req := httptest.NewRequest("POST", fmt.Sprintf("/ping/%v?query1=asd&query2=123", "12"), strings.NewReader(reqBody))
		req.Header.Set("Token", "asdfasdfasfdasdfasdfadsfasdfasdfasddf")
		req.Header.Set("header_id", "123123123")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

type SimpleHandler struct{}

func (h *SimpleHandler) HandleFunc(ctx context.Context, c *gin.Context) HandleResponse {
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
