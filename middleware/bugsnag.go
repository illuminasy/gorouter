package middleware

import (
	"net/http"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func configureBugsnag(config ErrorReportingConfig) {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:              config.APIKey,
		ReleaseStage:        config.ReleaseStage,
		AppType:             "mysqldatamanager",
		AppVersion:          config.AppVersion,
		ProjectPackages:     config.ProjectPackages,
		NotifyReleaseStages: config.NotifyReleaseStages,
		ParamsFilters:       []string{"password", "secret"},
		PanicHandler:        func() {},
	})
}

// bugsnagMiddleware configures and wraps bugsnag handler around the router
func bugsnagMiddleware(router http.Handler, config ErrorReportingConfig) http.Handler {
	configureBugsnag(config)
	return bugsnag.Handler(router)
}

func sendError(errorClass string, err error, a ...interface{}) error {
	// append error class so bugsnag can group errors using this
	a = append([]interface{}{bugsnag.ErrorClass{Name: errorClass}}, a...)
	return bugsnag.Notify(err, a...)
}
