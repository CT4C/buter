package cli

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Filter struct {
	Value    string
	Operator string
}

type Filters map[string]any

func (f *Filters) Set(value string) error {
	filerSeparator := ";"
	filterValueSeparator := ","
	filterRegexp := regexp.MustCompile("([^:]+):(.+)")
	rawFilters := strings.Split(value, filerSeparator)

	for _, filter := range rawFilters {
		matched := filterRegexp.FindStringSubmatch(filter)
		if len(matched) <= 0 {
			continue
		}

		filterName := matched[1]

		if strings.Contains(strings.Join(knownFilters, "/"), filterName) == false {
			return fmt.Errorf("unknown filter %s", filterName)
		}

		filterValue := strings.Split(matched[2], filterValueSeparator)

		switch filterName {
		case knownFilters[0]:
			fallthrough
		case knownFilters[1]:
			for _, stringedInteger := range filterValue {
				_, ok := (*f)[filterName]
				if ok == false {
					(*f)[filterName] = make([]int, 0)
				}

				converted, err := strconv.Atoi(stringedInteger)
				if err != nil {
					return err
				}

				(*f)[filterName] = append((*f)[filterName].([]int), converted)
			}
		}

	}

	return nil
}

func (f *Filters) String() string {
	b, err := json.Marshal(*f)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (f Filters) Status() []int {
	return f[knownFilters[0]].([]int)
}

func (f Filters) Length() []int {
	return f[knownFilters[1]].([]int)
}
