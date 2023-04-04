package cli

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStopper(t *testing.T) {
	t.Run("Expected: parse stopper successfully", func(t *testing.T) {
		// arrange
		rawInput := "status:200"
		expectedStopper := Stopper{}
		expectedStatuses := []string{fmt.Sprint(http.StatusOK)}

		// act
		expectedStopper.Set(rawInput)

		// assert
		require.Equal(t, expectedStatuses, expectedStopper.Status())
	})
}
