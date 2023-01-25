package cli

import (
	"fmt"
	"strings"
)

var (
	attackTypeFlag = "a"
	urlFlag        = "u"
	payloadFlag    = "p"
	threadsFlag    = "t"
	headersFlag    = "h"
	delayFlag      = "d"
	methodFlag     = "m"

	urlUsage        = "-u <http://localhost?param1=!abc!&param_N=!efg!> (payload wrapped into '!' char) (Url) "
	payloadUsage    = "-p <payload-file_1> -p <payload-file_N> (Payload)"
	attackTypeUsage = fmt.Sprintf("%s <%s> (AttackType)", "-a", strings.Join([]string{ClusterAttack}, "/"))
	threadsUsage    = "-t 5 (Max Concurrent Threads) "
	headersUsage    = "-h '{ \"User-Agent\": \"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/99.1\" }' (Headers) "
	delayUsage      = fmt.Sprintf("%s (Delay In Milliseconds)", "-d 800 ")
	methodUsage     = fmt.Sprintf("%s (Method)", "-m get")

	defaultAttackType = ""
	defaultUrl        = ""
	defaultThreads    = 3
	defaultHeaders    = ""
	defaultDealy      = 800
	defaultMethod     = "get"
)
