package dispatcher

import (
	"context"
	"errors"

	"github.com/edpryk/buter/cli"
)

type AttackConfig struct {
	AttackCompletedSig chan int

	cli.UserConfig
}

type Runner func(ctx context.Context, config AttackConfig)

func DispatchAttack(attack string) (Runner, error) {
	switch attack {
	case cli.SniperAttack:
		fallthrough
	case cli.ClusterAttack:
		return attackWithPayload, nil
	case cli.DOSAttack:
		return dosAttack, nil
	}

	return nil, errors.New("Attack not supported")
}
