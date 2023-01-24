package docs

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

	urlUsage        = "(Url) -u <http://localhost?param1=!abc!&param_N=!efg!> (payload wrapped into '!' char)"
	payloadUsage    = "(Payload) -p <payload-file_1> -p <payload-file_N>"
	attackTypeUsage = fmt.Sprintf("(AttackType) %10s <%s>", "-a", strings.Join([]string{ClusterAttack}, "/"))
	threadsUsage    = "(Max Concurrent Threads) -t 5"
	headersUsage    = "(Headers) -h '{ \"User-Agent\": \"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/99.1\" }'"
	delayUsage      = fmt.Sprintf("(Delay In Milliseconds)%10s", "-d 800 ")
	methodUsage     = fmt.Sprintf("(Method)%10s", "-m get")

	defaultAttackType = ""
	defaultUrl        = ""
	defaultThreads    = 3
	defaultHeaders    = ""
	defaultDealy      = 800
	defaultMethod     = "get"
)
