package logging

type LoggingHandle string

const (
	JsonHandler    = LoggingHandle("json")
	TextHandler    = LoggingHandle("text")
	SentryHandler  = LoggingHandle("sentry")
	RollbarHandler = LoggingHandle("rollbar")
)

type BackendLoggingHandle string

const (
	BackendJsonHandler = BackendLoggingHandle("json")
	BackendTextHandler = BackendLoggingHandle("text")
)
