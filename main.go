package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/runner"
)

var (
	config            cli.UserConfig
	payloadSet        [][]string
	totalPayloads     int
	rootContext       context.Context
	cancelRootContext context.CancelFunc

	err error

	sigEnd             = make(chan os.Signal)
	attackCompletedSig = make(chan int)
)

func main() {
	attackStartTime := time.Now()

	cli.PrintInfo()
	config = cli.ParseFlags()

	signal.Notify(sigEnd, syscall.SIGINT)
	log.SetFlags(2)

	if config.Timeout > 0 {
		rootContext, cancelRootContext = context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	} else {
		rootContext, cancelRootContext = context.WithCancel(context.Background())
	}
	defer cancelRootContext()

	go runner.RunAttack(rootContext, runner.AttackConfig{
		AttackCompletedSig: attackCompletedSig,
		UserConfig:         config,
	})

	select {
	case <-sigEnd:
		log.Printf("%3s Closed by Interruption\n", "")
	case <-attackCompletedSig:
	}

	log.Printf("%7s Attack completed in %s\n", "Summary:", time.Now().Sub(attackStartTime))
}
