package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testData struct {
	Name string `json:"name"`
}

func TestToString(t *testing.T) {
	t.Run("Expected: Struct stringed", func(t *testing.T) {
		// arrange
		target := testData{
			Name: "John",
		}
		expectedString := "{\"name\":\"John\"}"

		// act
		result, err := ToString(target)
		require.Equal(t, nil, err)
		require.Equal(t, expectedString, result)
	})
}
