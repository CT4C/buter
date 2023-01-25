package cli

import (
	"encoding/json"
	"fmt"
	"os"
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

type Headers string

func (h *Headers) Set(value string) error {
	if len(value) < 1 {
		return nil
	}

	d := make(map[string]string)

	if err := json.Unmarshal([]byte(value), &d); err != nil {
		fmt.Println("Can't parse headers", err)
		os.Exit(1)
	}
	b, _ := json.Marshal(d)

	parsed := Headers(b)

	(*h) = parsed

	return nil
}

func (h *Headers) String() string {
	return string(*h)
}
