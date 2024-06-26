package vgin

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/lg/v2"
)

type Data gin.H

type HandleResponse interface {
	GetCode() int
	GetError() error
	GetData() any
	GetMessage() string
	GetJson() string
}

type Ret struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"-"`
}

func (r *Ret) GetCode() int {
	return r.Code
}

func (r *Ret) GetError() error {
	return r.Err
}

func (r *Ret) GetData() any {
	return r
}

func (r *Ret) GetMessage() string {
	return r.Message
}

func (r *Ret) GetJson() string {
	b, err := json.Marshal(r)
	lg.PanicError(err)
	return string(b)
}

func FailedRet(code int, message string) *Ret {
	return ErrorRet(code, errors.New(message), message)
}

func ErrorRet(code int, err error, message string) *Ret {
	return &Ret{Code: code, Err: err, Message: message}
}

func SuccessRet(data any) *Ret {
	return &Ret{Code: http.StatusOK, Data: data}
}

func AbortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Ret{
		Code:    code,
		Message: message,
	})
	c.Error(errors.New(message))
}

func StatusOk(c *gin.Context, data any) {
	ReturnWithStatus(c, 200, data)
}

func ReturnWithStatus(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}
