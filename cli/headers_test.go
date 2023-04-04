package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	t.Run("Expected: headers is set successfully", func(t *testing.T) {
		// arrange
		rawInput := "Content-Type:application/json  Connection:keep-alive"
		headers := Headers{}

		// act
		headers.Set(rawInput)

		// assert
		require.Equal(t, headers["Content-Type"], "application/json")
		require.Equal(t, headers["Connection"], "keep-alive")
	})
}
