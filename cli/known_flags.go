package cli

var (
	attackTypeFlag    = "a"
	urlFlag           = "u"
	payloadFlag       = "p"
	threadsFlag       = "t"
	headersFlag       = "h"
	delayFlag         = "d"
	methodFlag        = "m"
	retriesAmountFlag = "r"
	retriesDelayFlag  = "rd"
	bodyFlag          = "b"
	timeoutFlag       = "T"
	dosRequestsFlag   = "R"
	filterOutFlag     = "f"
	stopFlag          = "S"
)

var knownFilters = []string{"status", "length"}
var knownStoppers = []string{"status"}
