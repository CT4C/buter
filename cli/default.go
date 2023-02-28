package cli

import "net/http"

var (
	defaultAttackType     = ""
	defaultUrl            = ""
	defaultThreads        = 10
	defaultHeaders        = ""
	defaultDelay          = 150
	defaultMethod         = http.MethodGet
	defaultRetriesAmount  = 3
	defaultRetryDelay     = 2000
	defaultBody           = ""
	defaultTimeout        = 0
	defaultDosRequests    = 10
	defaultFilters        = "NotSet"
	defaultStop           = ""
	defaultConfig         = ""
	defaultConfigTemplate = 0
)
