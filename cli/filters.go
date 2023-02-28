package cli

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Filters map[string][]string

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
					(*f)[filterName] = make([]string, 0)
				}

				(*f)[filterName] = append((*f)[filterName], stringedInteger)
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
