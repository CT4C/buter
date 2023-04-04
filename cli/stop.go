package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edpryk/buter/lib/convert"
)

type Stopper map[string][]string

func (f *Stopper) Join(key string, value any) {
	if !strings.Contains(strings.Join(knownStoppers, "/"), key) {
		log.Println(fmt.Errorf("unknown filter %s", key))
		os.Exit(1)
	}

	_, ok := (*f)[key]
	if !ok {
		(*f)[key] = make([]string, 0)
	}

	(*f)[key] = append((*f)[key], fmt.Sprint(value))
}

func (f *Stopper) Set(value string) error {
	filersSeparator := ";"
	filterparseValueSeparator := ","
	parseKeyparseValueSeparator := ":"
	convert.StringToKeyValue[string](value, filersSeparator, filterparseValueSeparator, parseKeyparseValueSeparator, f)
	return nil
}

func (f *Stopper) String() string {
	b, err := json.Marshal(*f)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (f Stopper) Status() []string {
	return f[knownStoppers[0]]
}
