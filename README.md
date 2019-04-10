# GORouter [![Build Status](https://travis-ci.org/Illuminasy/gorouter.svg?branch=master)](https://travis-ci.org/Illuminasy/gorouter) [![GoDoc](https://godoc.org/github.com/Illuminasy/gorouter?status.svg)](https://godoc.org/github.com/Illuminasy/gorouter) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Illuminasy/gorouter/blob/master/LICENSE.md)

GoRouter a package which wraps around https://github.com/julienschmidt/httprouter
allows to start router with middlewares

Currently supported middlewares:
1) NewRelic - APM monitoring tool
2) Bugsnag - Log Management
 
# Usage

Get the library:

    $ go get -v github.com/Illuminasy/gorouter

Just the router (no middleware)
```go
func startServer() {
	fmt.Println("Listening for http on port 80")
	router := gorouter.GetRouter(routes(), []string{})
	log.Fatal(http.ListenAndServe(":80", router))
}

func routes() gorouter.Routes {
	return gorouter.Routes{
		List: []gorouter.Route{
			gorouter.Route{
				Method: "GET",
				Path:   "/robots.txt",
				Handler: gorouter.PlainTextHandler([]string{
					"User-agent: *",
					"Disallow: /",
				}),
			},
			gorouter.Route{
				Method:  "GET",
				Path:    "/healthz",
				Handler: gorouter.JsonHandler(healthzHandler),
			},
		},
		PanicHandler: func(w http.ResponseWriter, r *http.Request, err interface{}) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, "")
			log.Println(err)
		},
	}
}

func healthzHandler() string {
	return `{"status":"up"}`
}

```

Router with middlewares
```go
package somepackge

import "github.com/Illuminasy/gorouter"

func startServer() {
	mc := middleware.Config{
		Bugsnag: middleware.BugsnagConfig{
			APIKey:       "apikey",
			ReleaseStage: "testing",
			ProjectPackages: []string{
				"main",
			},
			NotifyReleaseStages: []string{
				"testing",
			},
			AppVersion: "0.1.0",
		},
		Newrelic: middleware.NewrelicConfig{
			AppName: "Testing",
			License: "newrelickey",
		},
	}

	fmt.Println("Listening for http on port 80")
	router := gorouter.GetRouterWithMiddleware(mc, routes(), []string{})
	log.Fatal(http.ListenAndServe(":80", router))
}

func routes() gorouter.Routes {
	return gorouter.Routes{
		List: []gorouter.Route{
			gorouter.Route{
				Method: "GET",
				Path:   "/robots.txt",
				Handler: gorouter.PlainTextHandler([]string{
					"User-agent: *",
					"Disallow: /",
				}),
			},
			gorouter.Route{
				Method:  "GET",
				Path:    "/healthz",
				Handler: gorouter.JsonHandler(healthzHandler),
			},
		},
		PanicHandler: func(w http.ResponseWriter, r *http.Request, err interface{}) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, "")
			log.Println(err)
		},
	}
}

func healthzHandler() string {
	return `{"status":"up"}`
}
```

## Thanks to
https://github.com/julienschmidt/httprouter

https://github.com/newrelic/go-agent

https://github.com/bugsnag/bugsnag-go
