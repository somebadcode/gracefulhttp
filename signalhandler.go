package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func signalHandler(ctx context.Context, logger *log.Logger) {
	ch := make(chan os.Signal, 5)
	defer close(ch)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(ch)

	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGINT:
				logger.Print("terminating")
				return
			default:
				logger.Print("unhandled signal")
				return
			}
		case <-ctx.Done():
			logger.Print("terminating due context cancellation")
			return
		}
	}
}
