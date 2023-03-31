package reporter

import (
	"fmt"
)

type reportLine [][]any

func (line *reportLine) add(key string, value any) {
	(*line) = append(*line, []any{key, value})
}

func (line reportLine) string() string {
	report := ""

	for _, properties := range line {
		report += fmt.Sprintf("%s: %9s ", properties[0], fmt.Sprint(properties[1]))
	}

	return report
}
