package middleware

import (
	"net/http"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

// BugsnagConfig bugsnag config
type BugsnagConfig struct {
	APIKey              string
	ReleaseStage        string
	AppType             string
	AppVersion          string
	ProjectPackages     []string
	NotifyReleaseStages []string
	ParamsFilters       []string
	PanicHandler        func()
}

func configureBugsnag(config BugsnagConfig) {
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

// BugsnagMiddleware configures and wraps bugsnag handler around the router
func BugsnagMiddleware(router http.Handler, config BugsnagConfig) http.Handler {
	configureBugsnag(config)
	return bugsnag.Handler(router)
}
