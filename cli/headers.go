package cli

import (
	"encoding/json"

	"github.com/edpryk/buter/lib/convert"
)

type Headers map[string]string

func (h *Headers) Join(key string, value any) {
	(*h)[key] = value.(string)
}

func (h *Headers) Set(value string) error {
	lineSeparator := "  "
	keyValueSeparator := ":"
	valueSeparator := " "
	// headerPattern := regexp.MustCompile("(?P<key>[^:]+):[ ]{0,1}(?P<value>.+)")

	// *h = make(map[string]string)

	// if len(value) == 0 {
	// 	return nil
	// }

	// for _, subString := range strings.Split(value, "  ") {
	// 	matched := headerPattern.FindAllStringSubmatch(subString, 1)
	// 	if matched != nil {
	// 		(*h)[matched[0][1]] = matched[0][2]
	// 	}
	// }
	convert.StringToKeyValue[string](value, lineSeparator, valueSeparator, keyValueSeparator, h)

	return nil
}

func (h Headers) String() string {
	dataB, err := json.Marshal(h)
	if err != nil {
		return "-"
	}

	return string(dataB)
}
