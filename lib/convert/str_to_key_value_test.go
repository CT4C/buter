package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type box map[string][]string

func (b *box) Join(key string, value any) {
	if _, ok := (*b)[key]; !ok {
		(*b)[key] = make([]string, 0)
	}

	(*b)[key] = append((*b)[key], value.(string))
}

func TestStrToKeyValue(t *testing.T) {
	inputText := "status:201,301;length:1337"
	lineSeparator := ";"
	keyValueSeparator := ":"
	valueSeparator := ","

	t.Run("Expected: Successfully convert string to key value", func(t *testing.T) {
		// arrange
		b := box{}
		expectedStatuses := []string{"201", "301"}
		expectedLength := []string{"1337"}

		// act
		StringToKeyValue[string](inputText, lineSeparator, valueSeparator, keyValueSeparator, &b)

		// assert
		require.Equal(t, expectedStatuses, b["status"])
		require.Equal(t, expectedLength, b["length"])
	})
}
