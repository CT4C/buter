package prepare

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var (
	AttackValueSeparator = "^"
	rePartsOfAttackValue = regexp.MustCompile(fmt.Sprintf("([^%s]+)", AttackValueSeparator))
)

type AttackValue struct {
	Url     string
	Headers map[string]string
	Body    map[string]string
}

func ParseAttackValue(value string) (AttackValue, error) {
	attackValue := AttackValue{
		Headers: make(map[string]string),
		Body:    make(map[string]string),
	}

	matchedValues := rePartsOfAttackValue.FindAllString(value, -1)
	if matchedValues == nil {
		return attackValue, errors.New("Invalid attack value " + value)
	}

	if len(matchedValues) > 0 {
		attackValue.Url = matchedValues[0]
	}

	if len(matchedValues) > 1 {
		rawHeaders := matchedValues[1]

		if err := json.Unmarshal([]byte(rawHeaders), &attackValue.Headers); err != nil && len(rawHeaders) > 2 {
			return attackValue, err
		}
	}

	if len(matchedValues) > 2 {
		rawBody := matchedValues[2]
		if err := json.Unmarshal([]byte(rawBody), &attackValue.Body); err != nil {
			return attackValue, err
		}
	}

	return attackValue, nil
}