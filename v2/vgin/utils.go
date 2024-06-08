package vgin

import "github.com/superwhys/venkit/lg"

func guessHandlerName(handler Handler) string {

	return lg.FuncName(handler)
}
