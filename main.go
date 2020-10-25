package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	os.Exit(start())
}

func start() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := log.New(os.Stderr, "", 0)

	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		signalHandler(ctx, logger)
	}()

	err := Run(ctx, logger, nil)
	if err != nil {
		logger.Printf("http server returned with the error: %v", err)
		return -1
	}

	return 0
}
