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

type newrelicConfig struct {
	Enabled         bool
	AppName         string
	License         string
	Labels          map[string]string
	HostDisplayName string
}

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

func getnewrelicApp(config NewrelicConfig) (newrelic.Application, error) {
	var err error
	if newRelicApp == nil {
		err = configureNewRelic(config)
	}

	return newRelicApp, err
}

func newrelicMiddleware(handler httprouter.Handle, path string, config NewrelicConfig) httprouter.Handle {
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

func getNewrelicTransaction(id string, name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	if txn, ok := w.(newrelic.Transaction); ok {
		transactionList[id][name] = txn
		return txn
	}

	if list, ok := transactionList[id]; ok {
		if txn, ok := list[name]; ok {
			return txn
		}
	}

	var txn newrelic.Transaction

	if newRelicApp != nil {
		txn = newRelicApp.StartTransaction(name, w, r)
		transactionList[id][name] = txn
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

func newrelicNoticeError(err error) {
	txn := getNewrelicTransaction(txnID, txnName, nil, nil)
	txn.NoticeError(err)
}

func configureNewRelic(config newrelicConfig) error {
	cfg := newrelic.NewConfig(config.AppName, config.License)
	cfg.Enabled = config.Enabled
	cfg.Labels = config.Labels
	cfg.HostDisplayName = config.HostDisplayName
	cfg.Logger = newrelic.NewLogger(os.Stdout)
	cfg.DistributedTracer.Enabled = true

	err := cfg.Validate()
	if err != nil {
		return err
	}

	newRelicApp, err = newrelic.NewApplication(cfg)
	return err
}
