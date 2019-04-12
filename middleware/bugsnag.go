package middleware

import (
	"fmt"
	"net/http"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

var bugsnagConfigured = false

func configureBugsnag(config ErrorReportingConfig) {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:              config.APIKey,
		ReleaseStage:        config.ReleaseStage,
		AppType:             config.AppType,
		AppVersion:          config.AppVersion,
		ProjectPackages:     config.ProjectPackages,
		NotifyReleaseStages: config.NotifyReleaseStages,
		ParamsFilters:       []string{"password", "secret"},
		PanicHandler:        func() {},
		Hostname:            config.Hostname,
	})

	bugsnagConfigured = true
}

// bugsnagMiddleware configures and wraps bugsnag handler around the router
func bugsnagMiddleware(router http.Handler, config ErrorReportingConfig) http.Handler {
	configureBugsnag(config)
	return bugsnag.Handler(router)
}

func sendError(errorClass string, err error, a ...interface{}) error {
	if bugsnagConfigured {
		// append error class so bugsnag can group errors using this
		a = append([]interface{}{bugsnag.ErrorClass{Name: errorClass}}, a...)
		return bugsnag.Notify(err, a...)
	}

	return fmt.Errorf("Bugsnag has not been configured")
}
