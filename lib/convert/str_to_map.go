package convert

import (
	"regexp"
	"strings"
)

type MapJoin interface {
	Join(key string, value any)
}

func StringToKeyValue[T comparable](text string, lineSeparator string, valueSeparator string, keyValuePattern string, mapJoin MapJoin) {
	re := regexp.MustCompile(keyValuePattern)
	lines := strings.Split(text, lineSeparator)

	for _, line := range lines {
		matched := re.FindStringSubmatch(line)
		if len(matched) <= 0 {
			continue
		}

		key := matched[1]

		for _, value := range strings.Split(matched[2], valueSeparator) {
			mapJoin.Join(key, value)
		}
	}
}
