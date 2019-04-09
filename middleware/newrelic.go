package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	newrelic "github.com/newrelic/go-agent"
)

var newRelicApp newrelic.Application

var newrelicTransactionList map[string]map[string]newrelic.Transaction

type newrelicDataStore struct {
	Product            string
	Collection         string
	Operation          string
	ParameterizedQuery string
	QueryParameters    map[string]interface{}
	Host               string
	PortPathOrID       string
	DatabaseName       string
}

func newrelicMiddleware(handler httprouter.Handle, path string, config MetricCollectorConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		newRelicApp, err := getnewrelicApp(config)
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

func getnewrelicApp(config MetricCollectorConfig) (newrelic.Application, error) {
	var err error
	if newRelicApp == nil {
		err = configureNewRelic(config)
	}

	return newRelicApp, err
}

func configureNewRelic(config MetricCollectorConfig) error {
	cfg := newrelic.NewConfig(config.AppName, config.License)
	cfg.Enabled = config.Enabled
	cfg.Labels = config.Labels
	cfg.HostDisplayName = config.HostDisplayName
	cfg.DistributedTracer.Enabled = true

	if config.Debug {
		cfg.Logger = newrelic.NewDebugLogger(os.Stdout)
	} else {
		cfg.Logger = newrelic.NewLogger(os.Stdout)
	}

	err := cfg.Validate()
	if err != nil {
		return err
	}

	newRelicApp, err = newrelic.NewApplication(cfg)
	return err
}

func getNewrelicTransaction(id string, name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	if txn, ok := w.(newrelic.Transaction); ok {
		newrelicTransactionList[id][name] = txn
		return txn
	}

	if list, ok := newrelicTransactionList[id]; ok {
		if txn, ok := list[name]; ok {
			return txn
		}
	}

	var txn newrelic.Transaction

	if newRelicApp != nil {
		txn = newRelicApp.StartTransaction(name, w, r)
		newrelicTransactionList[id][name] = txn
	}

	return txn
}

func startNewrelicSegment(txnID string, txnName string, segmentName string, w http.ResponseWriter, r *http.Request) *newrelic.Segment {
	txn := getNewrelicTransaction(txnID, txnName, w, r)
	segment := newrelic.StartSegment(txn, segmentName)

	return segment
}

func startNewrelicDataStoreSegment(txnID string, txnName string, datastore DataStore, w http.ResponseWriter, r *http.Request) newrelic.DatastoreSegment {
	txn := getNewrelicTransaction(txnID, txnName, w, r)
	s := newrelic.DatastoreSegment{
		Product:            newrelic.DatastoreMySQL,
		Collection:         "users",
		Operation:          "INSERT",
		ParameterizedQuery: "INSERT INTO users (name, age) VALUES ($1, $2)",
		QueryParameters: map[string]interface{}{
			"name": "Dracula",
			"age":  439,
		},
		Host:         "mysql-server-1",
		PortPathOrID: "3306",
		DatabaseName: "my_database",
	}
	s.StartTime = newrelic.StartSegmentNow(txn)

	return s
}

func newrelicNoticeError(txnID string, txnName string, err error, w http.ResponseWriter, r *http.Request) {
	txn := getNewrelicTransaction(txnID, txnName, w, r)
	txn.NoticeError(err)
}
