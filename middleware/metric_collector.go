package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	newrelic "github.com/newrelic/go-agent"
)

var metricCollector MetricCollector

// MetricCollectorConfig Metric Collector config
type MetricCollectorConfig struct {
	Enabled         bool
	Newrelic        bool
	Debug           bool
	AppName         string
	License         string
	Labels          map[string]string
	HostDisplayName string
}

// MetricCollector Middleware with list of metric collectors to use
type MetricCollector struct {
	NewrelicApp newrelic.Application
}

// MetricCollectorTxn Metric Collector transaction
type MetricCollectorTxn struct {
	Txn newrelic.Transaction
}

// MetricCollectorSegment Metric Collector segment
type MetricCollectorSegment struct {
	Segment *newrelic.Segment
}

// MetricCollectorDatastoreSegment Metric Collector datastore segment
type MetricCollectorDatastoreSegment struct {
	DatastoreSegment newrelic.DatastoreSegment
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

// metricCollectorMiddleware Creates and starts Metric Collector middlware
func metricCollectorMiddleware(handler httprouter.Handle, path string, config MetricCollectorConfig) httprouter.Handle {
	if config.Newrelic && len(config.License) > 0 {
		return newrelicMiddleware(handler, path, config)
	}

	return handler
}

// GetMetricCollectorTransaction returns new or existing metric collector transanction
func GetMetricCollectorTransaction(txnID string, txnName string, w http.ResponseWriter, r *http.Request) MetricCollectorTxn {
	var mcTxn MetricCollectorTxn
	if newRelicApp != nil {
		mcTxn.Txn = getNewrelicTransaction(txnID, txnName, w, r)
	}

	return mcTxn
}

// StartMetricCollectorSegment Starts and retuns a metric collector segment for a transaction
func StartMetricCollectorSegment(txnID string, txnName string, segmentName string, w http.ResponseWriter, r *http.Request) MetricCollectorSegment {
	var mxnSgmt MetricCollectorSegment

	if metricCollector.NewrelicApp != nil {
		mxnSgmt.Segment = startNewrelicSegment(txnID, txnName, segmentName, w, r)
	}

	return mxnSgmt
}

// StartMetricCollectorDataStoreSegment Starts and retuns a metric collector datastore segment for a transaction
func StartMetricCollectorDataStoreSegment(txnID string, txnName string, datastore DataStore, w http.ResponseWriter, r *http.Request) MetricCollectorDatastoreSegment {
	var mxnDsSgmt MetricCollectorDatastoreSegment

	if metricCollector.NewrelicApp != nil {
		mxnDsSgmt.DatastoreSegment = startNewrelicDataStoreSegment(txnID, txnName, datastore, w, r)
	}

	return mxnDsSgmt
}

// MetricCollectorNoticeError Send error to newrelic
func MetricCollectorNoticeError(txnID string, txnName string, err error, w http.ResponseWriter, r *http.Request) {
	if metricCollector.NewrelicApp != nil {
		newrelicNoticeError(txnID, txnName, err, w, r)
	}
}
