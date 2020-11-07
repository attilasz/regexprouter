package regexprouter

import (
	//	"log"
	"net/http"
	"regexp"
	"strings"
	//	"reflect"
)

type (
	PathParams     map[string]string
	RequestHandler func(http.ResponseWriter, *http.Request, interface{})
	RouterSub      map[string]RequestHandler
	RouterTop      map[string]RouterSub
	Router         struct {
		routes          RouterTop
		notFoundHandler RequestHandler
	}
)

func (router *Router) AddHandler(method string, path string, f RequestHandler) {
	if router.routes == nil {
		router.routes = make(RouterTop)
	}
	if router.routes[method] == nil {
		router.routes[method] = make(RouterSub)
	}
	router.routes[method][path] = f
}

func (router *Router) GET(path string, f RequestHandler) {
	router.AddHandler(http.MethodGet, path, f)
}

func (router *Router) POST(path string, f RequestHandler) {
	router.AddHandler(http.MethodPost, path, f)
}

func (router *Router) PUT(path string, f RequestHandler) {
	router.AddHandler(http.MethodPut, path, f)
}

func (router *Router) OPTIONS(path string, f RequestHandler) {
	router.AddHandler(http.MethodDelete, path, f)
}

func (router *Router) DELETE(path string, f RequestHandler) {
	router.AddHandler(http.MethodDelete, path, f)
}

func setupCORS(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	path := request.URL.Path
	handler := router.routes[method][path]
	//	log.Printf("got: %s\n", path)
	setupCORS(writer)
	if method == "OPTIONS" {
		return
	}
	if handler == nil {
		for route := range router.routes[method] {
			if strings.HasPrefix(route, ":") {
				if regex, err := regexp.Compile(route[1:]); err == nil {
					if match := regex.FindStringSubmatch(path); len(match) > 0 {
						//						log.Printf("path matched (%s): %s -> %s", regex, path, route)
						ret := make(PathParams)
						for i, name := range regex.SubexpNames() {
							if i > 0 && name != "" && match[i] != "" {
								ret[name] = match[i]
							}
						}
						router.routes[method][route](writer, request, ret)
						return
					}
				}
			}
		}

		if router.notFoundHandler == nil {
			http.NotFound(writer, request)
		} else {
			router.notFoundHandler(writer, request, nil)
		}
	} else {
		handler(writer, request, nil)
	}
}
