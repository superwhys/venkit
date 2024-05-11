package vgin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type PingHandler struct {
	JsonDataStr     string  `vjson:"name"`
	JsonDataInt     int     `vjson:"age"`
	JsonDataFloat64 float64 `vjson:"money"`
	QueryDataStr    string  `vquery:"query1"`
	QueryDataInt    int     `vquery:"query2"`
	PathData        int     `vpath:"test_id"`
	HeaderDataStr   string  `vheader:"Token"`
	HeaderDataInt   int     `vheader:"header_id"`
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
	r := NewWithEngine(gin.New(), gin.Recovery())
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
