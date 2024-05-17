package vgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Data gin.H

type HandleResponse interface {
	GetCode() int
	GetError() error
	GetData() any
	GetMessage() string
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

func (r *Ret) SuccessRet(data any) *Ret {
	r.Code = http.StatusOK
	r.Data = data
	return r
}

func (r *Ret) FailedRet(code int, err error, message string) *Ret {
	r.Code = code
	r.Err = err
	r.Message = message
	return r
}

func (r *Ret) PackContent(code int, data any, err error, message string) *Ret {
	r.Err = err
	r.Message = message
	r.Code = code
	r.Data = data
	return r
}

func FailedRet(code int, message string) *Ret {
	return ErrorRet(code, errors.New(message), message)
}

func ErrorRet(code int, err error, message string) *Ret {
	ret := &Ret{}
	return ret.FailedRet(code, err, message)
}

func SuccessRet(data any) *Ret {
	ret := &Ret{}
	return ret.SuccessRet(data)
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
