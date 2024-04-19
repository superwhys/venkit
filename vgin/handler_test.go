package vgin

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	JsonDataStr   string `vjson:"name"`
	JsonDataInt   int    `vjson:"age"`
	QueryDataStr  string `vquery:"query1"`
	QueryDataInt  int    `vquery:"query2"`
	PathData      int    `vpath:"test_id"`
	HeaderDataStr string `vheader:"Token"`
	HeaderDataInt int    `vheader:"header_id"`
}

func (h *PingHandler) HandleFunc(ctx context.Context, c *gin.Context) HandleResponse {
	return &Ret{
		Code: 200,
		Data: h,
	}
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
