package service

import "net/http/pprof"

func WithPprof() ServiceOption {
	return func(vk *VkService) {
		vk.httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		vk.httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		vk.httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		vk.httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		vk.httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		vk.httpMux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		vk.httpMux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		vk.httpMux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		vk.httpMux.Handle("/debug/pprof/block", pprof.Handler("block"))
	}
}
