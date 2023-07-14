package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/runner"
)

var (
	configs           []cli.UserConfig
	rootContext       context.Context
	cancelRootContext context.CancelFunc

	sigEnd             = make(chan os.Signal)
	attackCompletedSig = make(chan int)
)

func main() {
	attackStartTime := time.Now()
	configs = cli.ParseFlags()

	cli.PrintInfo()

	signal.Notify(sigEnd, syscall.SIGINT)
	log.SetFlags(2)

	for _, config := range configs {
		if config.Timeout > 0 {
			rootContext, cancelRootContext = context.WithTimeout(context.Background(), time.Duration(10*time.Second))
		} else {
			rootContext, cancelRootContext = context.WithCancel(context.Background())
		}

		if config.Delay <= 0 {
			config.Delay = 1
		}

		defer cancelRootContext()
		go runner.RunAttack(rootContext, runner.AttackConfig{
			AttackCompletedSig: attackCompletedSig,
			UserConfig:         config,
		})

		select {
		case <-sigEnd:
			log.Println("Closed by Interruption")
			cancelRootContext()
		case <-attackCompletedSig:
		}

		log.Printf("Attack completed in %s\n", time.Now().Sub(attackStartTime))
		runtime.GC()
	}
}
