package gracefulhttp

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// GracefulServer is a composite type of http.Server and the GracefulServer function. Must be created with New.
type GracefulServer struct {
	*http.Server
	logger  *zap.Logger
	context context.Context
}

func New(ctx context.Context, server *http.Server, options ...Option) *GracefulServer {
	gs := &GracefulServer{
		Server:  server,
		context: ctx,
		logger:  zap.NewNop(),
	}

	// Call the options to modify our graceful server
	for _, fn := range options {
		fn(gs)
	}

	return gs
}

func (gs *GracefulServer) GracefulServe(listener net.Listener) (err error) {
	// Create context based on the parent context so that the caller can cancel our context.
	// This context can also be cancelled when the HTTP server stops serving to get us out of a deadlock.
	httpCtx, cancel := context.WithCancel(gs.context)
	defer cancel()

	gs.logger = gs.logger.With(
		zap.String("addr", listener.Addr().String()),
		zap.String("network", listener.Addr().Network()),
	)

	// Create a logger with the log level Error that the HTTP server can use
	// to log its errors.
	gs.ErrorLog, err = zap.NewStdLogAt(gs.logger.Named("http-server"), zap.ErrorLevel)
	if err != nil {
		gs.logger.Error("failed to create standard logger for HTTP server",
			zap.Error(err),
		)
		return err
	}

	// Create a wait group to allow us to wait until the go routines have finished before we return to the caller.
	// A deferred call to Wait() will make sure that we always wait for the go routines to finish before we return to the caller regardless of what's happening-
	wg := sync.WaitGroup{}
	defer wg.Wait()

	// Start serving in a separate go routine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		gs.logger.Info("starting to serve")

		err = gs.Serve(listener)
		if err != http.ErrServerClosed && err != nil {
			gs.logger.Error("http server error",
				zap.Error(err),
			)
		}

		gs.logger.Info("stopped serving",
			zap.NamedError("cause", err),
		)
	}()

	// Let's get stuck until HTTP server context gets canceled. The caller and our go routine can cancel this
	// so this should never result in a deadlock.
	<-httpCtx.Done()

	// Make sure that we gracefully shut down the HTTP server. If the context was cancelled by the caller then we
	// must do this. If the context was cancelled by our go routine then this will essentially do nothing.
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeoutCancel()

	// Add some informative logging.
	deadline, hasDeadline := timeoutCtx.Deadline()
	if !hasDeadline {
		deadline = time.Time{}
	}
	gs.logger.Info("shutting down",
		zap.NamedError("cause", httpCtx.Err()),
		zap.Time("deadline", deadline),
	)

	if err = gs.Shutdown(timeoutCtx); err != nil {
		gs.logger.Error("http server return an error when shutting down",
			zap.Error(err),
		)
	}

	return err
}
