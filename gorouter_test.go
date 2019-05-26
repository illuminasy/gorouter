package gorouter

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/illuminasy/gorouter/middleware"
	"github.com/stretchr/testify/assert"
)

var routes = Routes{
	List: []Route{
		Route{
			Method: "GET",
			Path:   "/robots.txt",
			Handler: PlainTextHandler([]string{
				"User-agent: *",
				"Disallow: /",
			}),
		},
		Route{
			Method: "GET",
			Path:   "/healthz",
			Handler: JSONHandler(func(w http.ResponseWriter, r *http.Request) (string, int) {
				return `{"status":"up"}`, http.StatusOK
			}),
		},
		Route{
			Method: "GET",
			Path:   "/panic",
			Handler: JSONHandler(func(w http.ResponseWriter, r *http.Request) (string, int) {
				panic("testing...")
			}),
		},
	},
	PanicHandler: func(w http.ResponseWriter, r *http.Request, err interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "")
		log.Println(err)
	},
}

var mc = middleware.Config{
	ErrorReportingConfig: middleware.ErrorReportingConfig{
		Enabled:      true,
		Bugsnag:      true,
		APIKey:       "testing",
		ReleaseStage: "Testing",
		ProjectPackages: []string{
			"main",
		},
		NotifyReleaseStages: []string{
			"Testing",
		},
		AppVersion: "0.1.0",
	},
	MetricCollectorConfig: middleware.MetricCollectorConfig{
		Enabled:  true,
		Newrelic: true,
		Debug:    true,
		AppName:  "Testing",
		License:  "testing",
		Labels: map[string]string{
			"Environment": "Dev",
			"Version":     "0.1.0",
		},
		HostDisplayName: "testing.localhost",
	},
}

func TestGetRobots(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/robots.txt", nil)

	router := GetRouter(routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusOK, respRec.Code)
	assert.Equal(t, "User-agent: *\nDisallow: /", string(respRec.Body.Bytes()))
}

func TestGetHealthz(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)

	router := GetRouter(routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusOK, respRec.Code)
	assert.Contains(t, string(respRec.Body.Bytes()), `{"status":"up"}`)
}

func TestNonExistent(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/non-existent-endpoint", nil)

	router := GetRouter(routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusNotFound, respRec.Code)
	assert.Equal(t, string(respRec.Body.Bytes()), "404 page not found\n")
}

func TestPanic(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)

	router := GetRouter(routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusInternalServerError, respRec.Code)
}

func TestGetRobotsWithMiddleware(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/robots.txt", nil)

	router := GetRouterWithMiddleware(mc, routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusOK, respRec.Code)
	assert.Equal(t, "User-agent: *\nDisallow: /", string(respRec.Body.Bytes()))
}

func TestGetHealthzWithMiddleware(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)

	router := GetRouterWithMiddleware(mc, routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusOK, respRec.Code)
	assert.Contains(t, string(respRec.Body.Bytes()), `{"status":"up"}`)
}

func TestNonExistentWithMiddleware(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/non-existent-endpoint", nil)

	router := GetRouterWithMiddleware(mc, routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusNotFound, respRec.Code)
	assert.Equal(t, string(respRec.Body.Bytes()), "404 page not found\n")
}

func TestPanicWithMiddleware(t *testing.T) {
	respRec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)

	router := GetRouterWithMiddleware(mc, routes, []string{})

	router.ServeHTTP(respRec, req)
	assert.Equal(t, http.StatusInternalServerError, respRec.Code)
}
