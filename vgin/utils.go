package vgin

import "github.com/superwhys/venkit/lg"

func guessHandlerName(handler Handler) string {
	var handlerName string
	nh, ok := handler.(NameHandler)
	if ok {
		handlerName = nh.Name()
	} else {
		handlerName = lg.StructName(handler)
	}
	return handlerName
}
