package definer

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var (
	AttackValueSeparator = "^"
	rePartsOfAttackValue = regexp.MustCompile(fmt.Sprintf("([^%s])+", AttackValueSeparator))
)

type AttackValue struct {
	Url     string
	Headers map[string]string
}

func ParseAttackValues(value string) (AttackValue, error) {
	attackValue := AttackValue{
		Headers: make(map[string]string),
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
		if err := json.Unmarshal([]byte(rawHeaders), &attackValue.Headers); err != nil {
			return attackValue, err
		}
	}

	return attackValue, nil
}
