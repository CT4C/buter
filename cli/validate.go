package cli

import (
	"errors"
	"net/http"
)

var (
	errNoAttackType          = errors.New("No AttackType provided")
	errNoUrl                 = errors.New("No URL provided")
	errNoPayloads            = errors.New("No Payloads provided")
	errFewPayloadsForCluster = errors.New("To few payloads files for Cluster attack")
	errMethodGetAndBodyUsage = errors.New("Cannot user method GET when Body provided")
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
	if len(in.Body.String()) > 2 && in.Method == http.MethodGet {
		return errMethodGetAndBodyUsage
	}

	return nil
}
