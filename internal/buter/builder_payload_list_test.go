package buter

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type PayloadListener struct {
	mock.Mock
}

func (listener *PayloadListener) onUpdate(updatedTargetString string, payloadInserted string, payloadNumber int) {
	listener.Called(updatedTargetString, payloadInserted, payloadNumber)
}

func TestBuildPayloadList(t *testing.T) {
	// arrange
	targetString := "His name is !Name! and he lives in the Mexico"
	payloadList1 := []string{"Bill", "Matthew"}

	entryPayloadNode := PayloadNode{
		Number:            0,
		Points:            [2]int{0, 6},
		PayloadList:       payloadList1,
		WorkingPayload:    targetString,
		CurrentPayloadIdx: 0,
		PreviousNode:      nil,
	}

	payloadListener := PayloadListener{}
	payloadListener.
		On("onUpdate", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Return()

	t.Run("Expected: Updated targetString with payload", func(t *testing.T) {
		// act
		produced := buildPayloadList(targetString, &entryPayloadNode, payloadListener.onUpdate)

		require.Equal(t, 2, produced)
	})
}
