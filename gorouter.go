package gorouter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/illuminasy/gorouter/middleware"

	"github.com/julienschmidt/httprouter"
)

type Routes struct {
	List         []Route
	PanicHandler func(w http.ResponseWriter, r *http.Request, err interface{})
}

type Route struct {
	Method  string
	Path    string
	Handler httprouter.Handle
}

var headersAllowedByCORS = []string{
	"Host",
	"Content-Type",
	"Connection",
	"User-Agent",
	"Cache-Control",
	"Accept-Encoding",
}

// PlainTextHandler handles plain text responses with appropriate headers
func PlainTextHandler(lines []string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		contents := []byte(strings.Join(lines, "\n"))
		w.Header().Set("Content-Type", "text/plain")
		w.Write(contents)
	}
}

// JsonHandler handles json responses with appropriate headers
func JsonHandler(handler func(http.ResponseWriter, *http.Request) string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		body := handler(w, r)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}
}

// GetRouter returns a router, optionally additional headers can be passed to set
func GetRouter(routes Routes, additionalHeaders []string) *httprouter.Router {
	return createRouter(routes, additionalHeaders, middleware.MiddlewareConfig{})
}

// GetRouterWithMiddleware returns a router with middlewares wrapped around it, optionally additional headers can be passed to set
func GetRouterWithMiddleware(mc middleware.Config, routes Routes, additionalHeaders []string) http.Handler {
	router := createRouter(routes, additionalHeaders, mc)
	return middleware.ErrorReportingMiddleware(router, mc.ErrorReportingConfig)
}

// createRouter creates a router
func createRouter(routes Routes, additionalHeaders []string, mc middleware.Config) *httprouter.Router {
	apiMethods := map[string][]string{}

	router := httprouter.New()
	router.PanicHandler = routes.PanicHandler
	for _, route := range routes.List {
		apiMethods[route.Path] = append(apiMethods[route.Path], route.Method)
		router.Handle(route.Method, route.Path, wrapMiddlewares(route.Handler, route.Path, mc))
	}

	for k, v := range apiMethods {
		router.OPTIONS(k, constructOptions(v, additionalHeaders))
	}

	return router
}

func wrapMiddlewares(handler httprouter.Handle, path string, mc middleware.Config) httprouter.Handle {
	return middleware.Wrapper(handler, path, mc)
}

func constructOptions(methods []string, additionalHeaders []string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	methodCsv := strings.Join(append(methods, "OPTIONS"), ",")
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		decorateWithCORS(w.Header(), methodCsv, additionalHeaders)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{}")
	}
}

func decorateWithCORS(headers http.Header, methods string, additionalHeaders []string) {
	allowedHeaders := append(headersAllowedByCORS, additionalHeaders...)
	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Access-Control-Allow-Methods", methods)
	headers.Set("Access-Control-Allow-Headers",
		strings.Join(allowedHeaders, ","),
	)
}
