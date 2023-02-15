package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/modules/dispatcher"
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

	signal.Notify(sigEnd, syscall.SIGINT)
	log.SetFlags(2)

	/*
		Need to test target connection before start
	*/

	config = cli.ParseFlags()

	if config.Timeout > 0 {
		rootContext, cancelRootContext = context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	} else {
		rootContext, cancelRootContext = context.WithCancel(context.Background())
	}
	defer cancelRootContext()

	attackRunner, err := dispatcher.DispatchAttack(config.AttackType)
	if err != nil {
		fmt.Println(err)
		cancelRootContext()
		os.Exit(1)
	}

	go attackRunner(rootContext, dispatcher.AttackConfig{
		UserConfig:         config,
		AttackCompletedSig: attackCompletedSig,
	})

	select {
	case <-sigEnd:
		log.Println("Closed by Interruption")
	case <-attackCompletedSig:
	}

	log.Printf("%3s Attack completed in %s\n", "", time.Now().Sub(attackStartTime))
}
