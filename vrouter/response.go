package vrouter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/superwhys/venkit/lg/v2"
)

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

func SuccessResponse(data any) *Ret {
	return &Ret{
		Code: http.StatusOK,
		Data: data,
	}
}

func ErrorResponse(code int, message string) *Ret {
	return &Ret{
		Code:    code,
		Message: message,
		Err:     errors.New(message),
	}
}
