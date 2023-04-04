package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/edpryk/buter/lib/convert"
)

type Filters map[string][]string

func (f *Filters) Join(key string, value any) {
	if !strings.Contains(strings.Join(knownFilters, "/"), key) {
		log.Println(fmt.Errorf("unknown filter %s", key))
		os.Exit(1)
	}

	_, ok := (*f)[key]
	if !ok {
		(*f)[key] = make([]string, 0)
	}

	(*f)[key] = append((*f)[key], fmt.Sprint(value))
}

func (f *Filters) Set(value string) error {
	lineSeparator := ";"
	valueSeparator := ","
	keyValuePattern := "([^:]+):(.+)"
	convert.StringToKeyValue[string](value, lineSeparator, valueSeparator, keyValuePattern, f)
	return nil
}

func (f *Filters) String() string {
	b, err := json.Marshal(*f)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (f Filters) Status() []string {
	return f[knownFilters[0]]
}

func (f Filters) Length() []string {
	return f[knownFilters[1]]
}

/*
	TODO: Move converter to single pkg
*/
func (f Filters) Duration() []int {
	d := f[knownFilters[2]]
	i := make([]int, 0)

	for _, filter := range d {
		convertedInt, err := strconv.Atoi(filter)
		if err != nil {
			panic(err)
		}
		i = append(i, convertedInt)
	}

	return i
}
