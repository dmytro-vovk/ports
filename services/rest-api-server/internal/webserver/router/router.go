package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func New(routes ...Route) *Router {
	r := Router{
		router: mux.NewRouter(),
	}

	for i := range routes {
		r.router.
			HandleFunc(routes[i].Path, routes[i].Handler).
			Methods(routes[i].Method)
	}

	return &r
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Printf("[%s] %s %s", request.RemoteAddr, request.Method, request.URL.Path)

	r.router.ServeHTTP(writer, request)
}
