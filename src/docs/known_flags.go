package docs

import (
	"fmt"
	"strings"
)

var (
	attackTypeFlag = "a"
	urlFlag        = "u"
	payloadFlag    = "p"

	urlUsage        = "-u <http://localhost?param1=!abc!&param_N=!efg!> (Payload wrapped into '!' char)"
	payloadUsage    = "-p <payload-file_1> -p <payload-file_N>"
	attackTypeUsage = fmt.Sprintf("-a <%s>", strings.Join([]string{ClusterAttack}, "/"))

	defaultAttackType = ""
	defaultUrl        = ""
)
