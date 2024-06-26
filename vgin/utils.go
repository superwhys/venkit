package vgin

import "github.com/superwhys/venkit/lg/v2"

func guessHandlerName(handler Handler) string {

	return lg.FuncName(handler)
}
