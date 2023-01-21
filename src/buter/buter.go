package buter

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/edpryk/buter/src/docs"
)

var (
	rePayloadPosition = regexp.MustCompile("(![^!]+!)")
)

/*
- Config includes payloads and url ? Attack Type ?
- Or it would be Cluster/Sniper instead of Run ?
and Attack type will operate on top level
*/
type Config struct {
	PayloadSet [][]string
	AttackType string
	Url        string
	Variants   int
	Ctx        context.Context
}

type UrlProvider chan string

func Run(config Config) (UrlProvider, error) {
	provider := make(UrlProvider)
	// text := "?param1=!x!&param2=!y!&param3=!z!"
	// payload1 := []string{"1", "2"}
	// payload2 := []string{"a", "b", "c"}
	// payload3 := []string{"L", "M", "N", "O"}

	// payloadsSet := [][]string{payload1, payload2, payload3}
	entryNode, err := transformPayload(config.Url, config.PayloadSet)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(config.Ctx)
	defer cancel()

	// Validate match length == payloads length
	go func() {
		/*
			- Handling text substitute like Cluster Bobm in Burp Suite
			- Prepare paylaods set
		*/

		// go func() {
		// 	select {
		// 	case <-ctx.Done():
		// 		cancel()
		// 		fmt.Println("Timeout")
		// 	}
		// }()

		switch config.AttackType {
		case docs.ClusterAttack:
			Cluster(ctx, config.Url, provider, &entryNode)
		default:
			fmt.Println(errAttackNotSupported)
			os.Exit(1)
		}
	}()

	return provider, nil
}
