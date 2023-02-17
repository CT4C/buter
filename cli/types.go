package cli

import (
	"encoding/json"
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
	*h = make(map[string]string)

	if len(value) == 0 {
		return nil
	}

	if err := json.Unmarshal([]byte(value), h); err != nil {
		return err
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

type Body map[string]string

func (b *Body) Set(value string) error {
	if b == nil {
		*b = make(map[string]string)
	}

	if len(value) == 0 {
		b = nil
		return nil
	}

	if err := json.Unmarshal([]byte(value), b); err != nil {
		return err
	}

	return nil
}

func (b Body) String() string {
	dataB, err := json.Marshal(b)
	if err != nil {
		return "-"
	}

	return string(dataB)
}
