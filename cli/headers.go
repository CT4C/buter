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
	valueSeparator := "  "

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
