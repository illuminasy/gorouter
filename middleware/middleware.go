package middleware

// MiddlewareConfig middleware config
type MiddlewareConfig struct {
	Bugsnag  BugsnagConfig
	Newrelic NewrelicConfig
}
