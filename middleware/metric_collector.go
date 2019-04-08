package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	newrelic "github.com/newrelic/go-agent"
)

var metricCollector MetricCollector

// MetricCollector Middleware with list of metric collectors to use
type MetricCollector struct {
	NewrelicApp newrelic.Application
}

// MetricCollectorTxn Metric Collector transaction
type MetricCollectorTxn struct {
	NewrelicTxn newrelic.Transaction
}

// MetricCollectorSegment Metric Collector segment
type MetricCollectorSegment struct {
	NewrelicSegment *newrelic.Segment
}

// MetricCollectorDatastoreSegment Metric Collector datastore segment
type MetricCollectorDatastoreSegment struct {
	NewrelicDatastoreSegment *newrelic.DatastoreSegment
}

// Config Metric Collector config
type Config struct {
	Enabled         bool
	Newrelic        bool
	AppName         string
	License         string
	Labels          map[string]string
	HostDisplayName string
}

// DataStore Metric Collector datastore
type DataStore struct {
	Product            string
	Collection         string
	Operation          string
	ParameterizedQuery string
	QueryParameters    map[string]interface{}
	Host               string
	PortPathOrID       string
	DatabaseName       string
}

// GetMetricCollector creates and returns new metric collector
func GetMetricCollector(config Config) (MetricCollector, error) {
	var err error

	if !Config.Enabled {
		return nil, err
	}

	if metricCollector == (MetricCollector{}) || (config.Newrelic && metricCollector.NewrelicApp == nil) {
		err = configureNewRelic(config)
		metricCollector.NewrelicApp = newRelicApp
	}

	return MetricCollector, err
}

// MetricCollectorMiddleware Creates and starts Metric Collector middlware
func MetricCollectorMiddleware(handler httprouter.Handle, path string, config Config) httprouter.Handle {
	if config.Newrelic {
		return newrelicMiddleware()
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handler(w, r, ps)
	}
}

// GetMetricCollectorTransaction returns new or existing metric collector transanction
func GetMetricCollectorTransaction(txnID string, txnName string, w http.ResponseWriter, r *http.Request) MetricCollectorTxn {
	var mcTxn MetricCollectorTxn
	if metricCollector.NewrelicApp != nil {
		mcTxn.NewrelicTxn = getNewrelicTransaction(txnID, txnName, w, r)
	}

	return mcTxn
}

// StartMetricCollectorSegment Starts and retuns a metric collector segment for a transaction
func StartMetricCollectorSegment(txnID string, txnName string, segmentName string, w http.ResponseWriter, r *http.Request) MetricCollectorSegment {
	var mxnSgmt MetricCollectorSegment

	if metricCollector.NewrelicApp != nil {
		mxnSgmt.NewrelicSegment = startNewrelicSegment(txn, segmentName, w, r)
	}

	return mxnSgmt
}

// StartMetricCollectorDataStoreSegment Starts and retuns a metric collector datastore segment for a transaction
func StartMetricCollectorDataStoreSegment(txnID string, txnName string, datastore DataStore, w http.ResponseWriter, r *http.Request) MetricCollectorDatastoreSegment {
	var mxnDsSgmt MetricCollectorDatastoreSegment

	if metricCollector.NewrelicApp != nil {
		mxnSgmt.NewrelicSegment = startNewrelicDataStoreSegment(txn, segmentName, datastore, w, r)
	}

	return mxnSgmt
}

// MetricCollectorNoticeError Send error to newrelic
func MetricCollectorNoticeError(err error) {
	if metricCollector.NewrelicApp != nil {
		newrelicNoticeError(err)
	}
}
