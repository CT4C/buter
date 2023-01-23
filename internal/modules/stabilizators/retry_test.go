package stabilizator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	t.Run("Succesfully after fail", func(t *testing.T) {
		// arrange
		expected := "test"
		expectedErr := errors.New("")
		expectedErr = nil
		failedAttempts := 2
		counter := 0

		caller := func() (any, error) {
			if counter != failedAttempts {
				counter++
				return "", errors.New("Timeout")
			}

			return expected, nil
		}

		// act
		result, err := Retry(caller, 3, 500)

		require.Equal(t, expectedErr, err)
		require.Equal(t, expected, result)
	})
	t.Run("UnSuccesfully after a lot fails", func(t *testing.T) {
		// arrange
		expected := "test"
		expectedErr := errors.New("Timeout")
		failedAttempts := 4
		counter := 0

		caller := func() (any, error) {
			if counter != failedAttempts {
				counter++
				return "", expectedErr
			}

			return expected, nil
		}

		// act
		result, err := Retry(caller, 3, 500)

		require.Equal(t, expectedErr, err)
		require.Equal(t, nil, result)
	})
}
