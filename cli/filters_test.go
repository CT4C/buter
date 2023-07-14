package cli

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilters(t *testing.T) {

	t.Run("Expected: Successfully parsed filters", func(t *testing.T) {
		// arrange
		f := Filters{}
		expected := Filters{}
		expected["status"] = []string{fmt.Sprint(http.StatusOK), fmt.Sprint(http.StatusCreated)}
		expected["length"] = []string{"1337"}
		expected["duration"] = []string{"1337"}
		expectedDuration := []int{1337}
		rawFilters := "status:200,201;length:1337;duration:1337"

		// act
		err := f.Set(rawFilters)

		// assert
		require.Equal(t, nil, err)
		require.Equal(t, expected.String(), f.String())
		require.Equal(t, expected["status"], f.Status())
		require.Equal(t, expected["length"], f.Length())
		require.Equal(t, expectedDuration, f.Duration())
	})
}
