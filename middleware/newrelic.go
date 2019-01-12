package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	newrelic "github.com/newrelic/go-agent"
)

var newRelicApp newrelic.Application

// NewrelicConfig newrelic config
type NewrelicConfig struct {
	AppName string
	License string
}

// GetNewRelicApp creates and returns new relic app
func GetNewRelicApp(config NewrelicConfig) (newrelic.Application, error) {
	var err error
	if newRelicApp == nil {
		err = configureNewRelic(config)
	}

	return newRelicApp, err
}

func configureNewRelic(config NewrelicConfig) error {
	cfg := newrelic.NewConfig(config.AppName, config.License)
	cfg.Logger = newrelic.NewLogger(os.Stdout)

	err := cfg.Validate()
	if err != nil {
		return err
	}

	newRelicApp, err = newrelic.NewApplication(cfg)
	return err
}

// NewRelicMiddleware creates newrelic app and starts transaction
func NewRelicMiddleware(handler httprouter.Handle, path string, config NewrelicConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		newRelicApp, err := GetNewRelicApp(config)
		if err != nil {
			log.Println(err)
		}
		if newRelicApp != nil && err == nil {
			txn := newRelicApp.StartTransaction(path, w, r)
			defer func() {
				err := txn.End()
				log.Println(err)
			}()
			handler(txn, r, ps)
			return
		}

		handler(w, r, ps)
	}
}
