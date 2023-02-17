package convert

import "encoding/json"

func ToString(o any) (string, error) {
	b, err := json.Marshal(o)
	return string(b), err
}
