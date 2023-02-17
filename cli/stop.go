package cli

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Stopper map[string]any

func (f *Stopper) Set(value string) error {
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

		if strings.Contains(strings.Join(knownStoppers, "/"), filterName) == false {
			return fmt.Errorf("unknown stopper %s", filterName)
		}

		filterValue := strings.Split(matched[2], filterValueSeparator)

		switch filterName {
		case knownStoppers[0]:
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

func (f *Stopper) String() string {
	b, err := json.Marshal(*f)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (f Stopper) Status() []int {
	return f[knownStoppers[0]].([]int)
}
