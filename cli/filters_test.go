package cli

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilters(t *testing.T) {

	t.Run("Expected: Successfully parsed filters", func(t *testing.T) {
		// arrange
		f := Filters{}
		expected := Filters{}
		expected["status"] = []int{http.StatusOK, http.StatusCreated}
		expected["length"] = []int{1337}
		rawFilters := "status:200,201;length:1337"

		// act
		err := f.Set(rawFilters)

		// assert
		require.Equal(t, nil, err)
		require.Equal(t, expected.String(), f.String())
	})
}
