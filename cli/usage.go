package cli

import (
	"fmt"
	"os"
	"strings"
)

var (
	urlUsage           = fmt.Sprintf("-%-3s %s %s", urlFlag, "<http://localhost?param1=!abc!&param_N=!efg!>", "(Url)")
	payloadUsage       = fmt.Sprintf("-%-3s %s %s", payloadFlag, "<payload-file_1> -p <payload-file_N>", "(Payload)")
	attackTypeUsage    = fmt.Sprintf("-%-3s %s %s", attackTypeFlag, strings.Join([]string{ClusterAttack, SniperAttack, DOSAttack}, "/"), "(AttackType)")
	threadsUsage       = fmt.Sprintf("-%-3s %s %s", threadsFlag, "5", "(Max Concurrent Threads)")
	headersUsage       = fmt.Sprintf("-%-3s %s %s", headersFlag, "'{ \"User-Agent\": \"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/99.1\" }'", "(Headers)")
	delayUsage         = fmt.Sprintf("-%-3s %s %s", delayFlag, "800", "(Delay in milliseconds)")
	methodUsage        = fmt.Sprintf("-%-3s %s %s", methodFlag, "POST", "(HTTP method)")
	retriesAmountUsage = fmt.Sprintf("-%-3s %s %s", retriesAmountFlag, "3", "(Retries on request error)")
	retryDelayUsage    = fmt.Sprintf("-%-3s %s %s", retriesDelayFlag, "1000", "(Retry delay in milliseconds)")
	bodyUsage          = fmt.Sprintf("-%-3s %s %s", bodyFlag, "{\"email\":\"user_nameg@mail.com\",\"password\":\"12345\"}", "(request body)")
	timeoutUsage       = fmt.Sprintf("-%-3s %s %s", timeoutFlag, "10", "(Request timeout in Seconds)")
	dosReqUsage        = fmt.Sprintf("-%-3s %s %s", dosRequestsFlag, "10", "(request amount in DOS mode)")
	filterOutUsage     = fmt.Sprintf("-%-3s %s %s", filterOutFlag, "status:200,201;length:1553", "(Output filters)")
	stopUsage          = fmt.Sprintf("-%-3s %s %s", stopFlag, "status:200", "(Stop attack on event)")
)

func printUsage() {
	fmt.Println()
	fmt.Printf("Usage: %s\n", os.Args[0])
	fmt.Printf("\t%s\n\n", "Any payload/fuzzing position must be highlighted with ! char")
	fmt.Printf("\t%s\n", urlUsage)
	fmt.Printf("\t%s\n", payloadUsage)
	fmt.Printf("\t%s\n", attackTypeUsage)
	fmt.Printf("\t%s\n", threadsUsage)
	fmt.Printf("\t%s\n", headersUsage)
	fmt.Printf("\t%s\n", delayUsage)
	fmt.Printf("\t%s\n", methodUsage)
	fmt.Printf("\t%s\n", retriesAmountUsage)
	fmt.Printf("\t%s\n", retryDelayUsage)
	fmt.Printf("\t%s\n", bodyUsage)
	fmt.Printf("\t%s\n", timeoutUsage)
	fmt.Printf("\t%s\n", dosReqUsage)
	fmt.Printf("\t%s\n", filterOutUsage)
	fmt.Printf("\t%s\n", stopUsage)
	fmt.Println()
}
