package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"sync"
	"time"
)

func Run(ctx context.Context, logger *log.Logger, tlsConfig *tls.Config) (err error) {
	httpCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", httpRouter)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		TLSConfig:         tlsConfig,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		MaxHeaderBytes:    4096,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          logger,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		var err error
		if server.TLSConfig == nil {
			err = server.ListenAndServe()
		} else {
			err = server.ListenAndServeTLS("", "")
		}

		if err != http.ErrServerClosed && err != nil {
			logger.Printf("http server error: %v", err)
		}
	}()

	<-httpCtx.Done()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutCancel()

	if err = server.Shutdown(timeoutCtx); err != http.ErrServerClosed && err != nil {
		logger.Printf("http server shutdown error: %v", err)
	}

	wg.Wait()

	return err
}
