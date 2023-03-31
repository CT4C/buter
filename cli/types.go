package cli

import (
	"encoding/json"
	"regexp"
	"strings"
)

type PayloadFiles []string

func (ps *PayloadFiles) Set(value string) error {
	*ps = append(*ps, value)
	return nil
}

func (ps *PayloadFiles) String() string {
	s := ""
	for _, v := range *ps {
		s += v
	}
	return s
}

type Headers map[string]string

func (h *Headers) Set(value string) error {
	headerPattern := regexp.MustCompile("(?P<key>[^:]+):[ ]{0,1}(?P<value>.+)")

	*h = make(map[string]string)

	if len(value) == 0 {
		return nil
	}

	for _, subString := range strings.Split(value, "  ") {
		matched := headerPattern.FindAllStringSubmatch(subString, 1)
		if matched != nil {
			(*h)[matched[0][1]] = matched[0][2]
		}
	}

	return nil
}

func (h Headers) String() string {
	dataB, err := json.Marshal(h)
	if err != nil {
		return "-"
	}

	return string(dataB)
}
