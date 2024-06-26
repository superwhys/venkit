package vgin

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	
	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/v2/lg"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type OneQueryData struct {
	User string `vquery:"user" form:"user"`
}

type JsonHandlerData struct {
	JsonDataStr     string  `vjson:"name" form:"name" json:"name"`
	JsonDataInt     int     `vjson:"age" form:"age" json:"age"`
	JsonDataFloat64 float64 `vjson:"money" form:"money" json:"money"`
	Address         string  `vjson:"address" form:"address" json:"address"`
	City            string  `vjson:"city" form:"city" json:"city"`
}

func UserHandler(ctx context.Context, c *Context, data *OneQueryData) HandleResponse {
	return SuccessRet(data.User)
}

func JsonHandler(ctx context.Context, c *Context, data *JsonHandlerData) HandleResponse {
	return SuccessRet(data)
}

func BenchmarkVginOneQuery(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery())
	
	r.POST("/user", UserHandler)
	
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
		data := new(OneQueryData)
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

func BenchmarkVgin(b *testing.B) {
	r := NewWithEngine(gin.New(), gin.Recovery())
	r.POST("/ping", JsonHandler)
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18, "money": 100, "address": "asdadad", "city": "city"}`
		req := httptest.NewRequest("POST", "/ping", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkOriginGin(b *testing.B) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/ping", func(ctx *gin.Context) {
		data := &JsonHandlerData{}
		if err := ctx.ShouldBind(data); err != nil {
			lg.Error(err)
			ctx.JSON(400, "parse data error")
			return
		}
	})
	
	for i := 0; i < b.N; i++ {
		reqBody := `{"name": "John Doe", "age": 18, "money": 100, "address": "asdadad", "city": "city"}`
		req := httptest.NewRequest("POST", "/ping", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
