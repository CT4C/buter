package cli

import "errors"

var (
	errNoAttackType          = errors.New("No AttackType provided")
	errNoUrl                 = errors.New("No URL provided")
	errNoPayloads            = errors.New("No Payloads provided")
	errFewPayloadsForCluster = errors.New("To few payloads files for Cluster attack")
)

func validateInput(in Input) error {
	if in.AttackType == "" {
		return errNoAttackType
	}
	if in.AttackType == ClusterAttack && len(in.PayloadFiles) < 2 {
		return errFewPayloadsForCluster
	}
	if in.Url == "" {
		return errNoUrl
	}
	if len(in.PayloadFiles) == 0 {
		return errNoPayloads
	}

	return nil
}
