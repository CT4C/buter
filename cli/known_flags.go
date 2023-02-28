package cli

var (
	attackTypeFlag     = "a"
	urlFlag            = "u"
	payloadFlag        = "p"
	threadsFlag        = "t"
	headersFlag        = "h"
	delayFlag          = "d"
	methodFlag         = "m"
	retriesAmountFlag  = "r"
	retriesDelayFlag   = "rd"
	bodyFlag           = "b"
	timeoutFlag        = "T"
	dosRequestsFlag    = "R"
	filterOutFlag      = "f"
	stopFlag           = "S"
	configFlag         = "c"
	configTemplateFlag = "cT"
)

var knownFilters = []string{"status", "length", "duration"}
var knownStoppers = []string{"status"}
