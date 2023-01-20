package docs

import "errors"

var (
	errNoAttackType = errors.New("No AttackType provided")
	errNoUrl        = errors.New("No URL provided")
	errNoPayloads   = errors.New("No Payloads provided")
)

func validateInput(in Input) error {
	if in.AttackType == "" {
		return errNoAttackType
	}
	if in.Url == "" {
		return errNoUrl
	}
	if len(in.PayloadFiles) == 0 {
		return errNoPayloads
	}

	return nil
}
