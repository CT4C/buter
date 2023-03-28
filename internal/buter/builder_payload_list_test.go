package buter

import (
	"strings"
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

var targetString = "His name is !Name! and he lives in the Mexico"
var points = [2]int{12, 18}
var payloadList1 = []string{"Bill", "Matthew"}

func TestBuildPayloadList(t *testing.T) {
	// arrange
	targetStringTest1 := targetString

	entryPayloadNode := PayloadNode{
		Number:            0,
		Points:            points,
		PayloadList:       payloadList1,
		WorkingPayload:    targetStringTest1,
		CurrentPayloadIdx: 0,
		PreviousNode:      nil,
	}

	payloadListener := PayloadListener{}
	payloadListener.
		On("onUpdate", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Return()

	t.Run("Expected: Updated targetString with payload", func(t *testing.T) {
		// act
		produced := buildPayloadList(targetStringTest1, &entryPayloadNode, payloadListener.onUpdate)

		require.Equal(t, 2, produced)
	})
}

func TestInsertPayload(t *testing.T) {
	t.Run("Expected: payload inserted successfully", func(t *testing.T) {
		// arrange
		targetStringTest2 := targetString
		payload := "Bill"
		payloadPattern := "!Name!"
		expectedString := strings.Replace(targetStringTest2, payloadPattern, payload, 1)

		// act
		updatedTargetString := insertPayload(targetStringTest2, payload, points)

		// assert
		require.Equal(t, expectedString, updatedTargetString)
	})
}
