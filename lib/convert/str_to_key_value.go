package convert

import (
	"strings"
)

type MapJoin interface {
	Join(key string, value any)
}

func StringToKeyValue[T comparable](text string, lineSeparator string, valueSeparator string, keyValueSeparator string, mapJoin MapJoin) {
	lines := strings.Split(text, lineSeparator)

	for _, line := range lines {
		keyValue := strings.Split(line, keyValueSeparator)
		if len(keyValue) <= 0 {
			continue
		}

		key := keyValue[0]

		for _, value := range strings.Split(keyValue[1], valueSeparator) {
			mapJoin.Join(key, value)
		}
	}
}
