package docs

var (
	attackTypeFlag = "a"
	urlFlag        = "u"
	payloadFlag    = "p"

	urlUsage        = "-u <http://localhost?param1=!abc!&param_N=!efg!> (Payload wrapped into '!' char)"
	payloadUsage    = "-p <payload-file1> -p <payload-fileN>"
	attackTypeUsage = "-a <Sniper/Cluster/PichFork>"

	defaultAttackType = "Sniper"
	defaultUrl        = ""
)
