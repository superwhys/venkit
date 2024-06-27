package main

import (
	"context"
	"fmt"
	"net/http"
	
	"github.com/gorilla/mux"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/vrouter/v2"
)

type routers []vrouter.Route

func (r routers) Routes() []vrouter.Route {
	return r
}

var (
	myRouters = routers{vrouter.NewRoute(http.MethodGet, "/test", helloHandler)}
)

func helloHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) vrouter.Response {
	return vrouter.SuccessResponse("hello world")
}

func testMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) vrouter.Response {
	fmt.Println(111)
	return nil
}

func main() {
	vrouter.SetMode(vrouter.DebugMode)
	router := vrouter.NewVrouter(vrouter.WithHost("0.0.0.0"))
	
	router.HandleRoute(http.MethodGet, "/hello/{name}", helloHandler, func(route *mux.Route) *mux.Route {
		return route.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")
	})
	router.HandlerRouter(myRouters)
	router.Static("/static", "./content", nil)
	
	group := router.Group("/group1")
	group.UseMiddleware(vrouter.HandleFunc(testMiddleware))
	group.HandleRoute(http.MethodGet, "/hello2/{name}", helloHandler, nil)
	
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	lg.PanicError(srv.ListenAndServe())
}
