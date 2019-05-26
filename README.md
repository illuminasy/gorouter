# GORouter [![Build Status](https://travis-ci.org/illuminasy/gorouter.svg?branch=master)](https://travis-ci.org/illuminasy/gorouter) [![Coverage Status](https://coveralls.io/repos/github/illuminasy/gorouter/badge.svg?branch=master)](https://coveralls.io/github/illuminasy/gorouter?branch=master) [![GoDoc](https://godoc.org/github.com/illuminasy/gorouter?status.svg)](https://godoc.org/github.com/illuminasy/gorouter) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/illuminasy/gorouter/blob/master/LICENSE.md)

GoRouter a package which wraps around https://github.com/julienschmidt/httprouter
allows to start router with middlewares

Currently supported middlewares:

1) NewRelic - APM monitoring tool
2) Bugsnag - Log Management
 
## Usage

Get the library:

    $ go get -v github.com/illuminasy/gorouter

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
                Handler: gorouter.JSONHandler(healthzHandler),
            },
            gorouter.Route{
                Method:  "GET",
                Path:    "/test",
                Handler: gorouter.HTMLHandler(healthzHandler),
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

func healthzHandler(w http.ResponseWriter, r *http.Request) (string, int) {
    return `{"status":"up"}`, http.StatusOK
}

```

Router with middlewares

```go
package somepackge

import "github.com/illuminasy/gorouter"
import "github.com/illuminasy/gorouter/middleware"

func startServer() {
    mc := middleware.Config{
        ErrorReportingConfig: middleware.ErrorReportingConfig{
            Enabled:      true,
            Bugsnag:      true,
            APIKey:       "testing",
            AppType:      "router",
            ReleaseStage: "Dev",
            AppVersion:   "0.1.0",
            ProjectPackages: []string{
                "main",
            },
            NotifyReleaseStages: []string{
                "Dev",
            },
            PanicHandler: func() {},
            Hostname: "localhost",
        },
        MetricCollectorConfig: middleware.MetricCollectorConfig{
            Enabled:  true,
            Newrelic: true,
            Debug:    true,
            AppName:  "TestApp",
            License:  "testing",
            Labels: map[string]string{
                "Environment": "Dev",
                "Version":     "0.1.0",
            },
            HostDisplayName: "localhost",
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
                Handler: gorouter.JSONHandler(healthzHandler),
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

func healthzHandler(w http.ResponseWriter, r *http.Request) (string, int) {
    return `{"status":"up"}`, http.StatusOK
}
```

Gorouter starts a web transaction by wrapping middleware around the handlers
but if you still need to create custom transactions (maybe non web stuff)

```go
    // w, r are w http.ResponseWriter, r *http.Request
    // these are optionals, just pass nil instead for non web transactions
    txn := middleware.GetMetricCollectorTransaction("someTxnID", "someTxnName", w, r)
    // Do something
    txn.End()
```

To create segments within a transaction

```go
    // w, r are w http.ResponseWriter, r *http.Request
    // these are optionals, just pass nil instead
    // if transaction doesnot exist then it creates new one.
    s := middleware.StartMetricCollectorSegment("someTxnID", "someTxnName", "someSegmentName", w, r)
    // Do something
    s.End()
```

To create datastore segments within a transaction say for mysql

```go
    dataStore := middleware.DataStore {
        Product: "mysql",
        Collection: "users",
        Operation: "INSERT",
        ParameterizedQuery: `INSERT INTO users (name, age) VALUES ($1, $2)"`,
        QueryParameters: map[string]interface{}{
            "name": "Dracula",
            "age": 439,
        },
        Host: "mysql-server-1",
        PortPathOrID: "3306",
        DatabaseName: "my_database",
    }
    // w, r are w http.ResponseWriter, r *http.Request
    // these are optionals, just pass nil instead
    // if transaction doesnot exist then it creates new one.
    s := middleware.StartMetricCollectorDataStoreSegment("someTxnID", "someTxnName", dataStore, w, r)
    // Do something
    s.End()
```

To report errors to metric collectors

```go
    // w, r are w http.ResponseWriter, r *http.Request
    // these are optionals, just pass nil instead
    // if transaction doesnot exist then it creates new one.
    middleware.MetricCollectorNoticeError("someTxnID", "someTxnName", errors.New("Invalid API config"), w, r)
```

Send Errors to bugsnag

```go
    middleware.ReportErrorToBugsnag("someErrorClass", errors.New("Invalid API config"))
```

## Thanks to

https://github.com/julienschmidt/httprouter

https://github.com/newrelic/go-agent

https://github.com/bugsnag/bugsnag-go
