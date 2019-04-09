package middleware

import (
	"net/http"
)

// ErrorReportingConfig Error reporting config
type ErrorReportingConfig struct {
	Enabled             bool
	Bugsnag             bool
	APIKey              string
	ReleaseStage        string
	AppType             string
	AppVersion          string
	ProjectPackages     []string
	NotifyReleaseStages []string
	ParamsFilters       []string
	PanicHandler        func()
}

// ErrorReportingMiddleware configures and wraps error reporting handler around the router
func ErrorReportingMiddleware(router http.Handler, config ErrorReportingConfig) http.Handler {
	if config.Enabled && config.Bugsnag {
		return bugsnagMiddleware(router, config)
	}

	return router
}
