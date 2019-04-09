package middleware

import "github.com/julienschmidt/httprouter"

// Config middleware config
type Config struct {
	ErrorReportingConfig  ErrorReportingConfig
	MetricCollectorConfig MetricCollectorConfig
}

func Wrapper(handler httprouter.Handle, path string, mc Config) httprouter.Handle {
	if mc.MetricCollectorConfig.Enabled {
		handler = metricCollectorMiddleware(handler, path, mc.MetricCollectorConfig)
	}

	return handler
}
