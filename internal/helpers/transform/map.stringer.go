package transform

import "encoding/json"

type MapStringer[MapValue any] struct {
	Value map[string]MapValue
}

func (ms MapStringer[any]) String() string {
	b, err := json.Marshal(ms.Value)
	if err != nil {
		return ""
	}

	return string(b)
}

func (ms MapStringer[any]) Map() map[string]any {
	return ms.Value
}

func NewMapStringer[MapValue any](v map[string]MapValue) MapStringer[MapValue] {
	return MapStringer[MapValue]{
		Value: v,
	}
}
