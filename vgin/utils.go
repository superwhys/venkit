package vgin

import "github.com/superwhys/venkit/v2/lg"

func guessHandlerName(handler Handler) string {
	
	return lg.FuncName(handler)
}
