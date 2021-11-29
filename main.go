package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/somebadcode/gracefulhttp/demoapp"
	"github.com/somebadcode/gracefulhttp/internal/gracefulhttp"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Create a base context so that we don't want to allow the context derived from the signal handler to cancel
	// ongoing connections currently served by http.Server.
	baseCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new context based on the base context. Canceling the base context will stop the signal handling.
	ctx, stop := signal.NotifyContext(baseCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
	defer stop()

	// Create a logger. Defer sync of logger.
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	ws := zapcore.Lock(zapcore.AddSync(os.Stderr))
	core := zapcore.NewCore(encoder, ws, zap.DebugLevel)
	logger := zap.New(core)
	// We would normally defer a call to logger.Sync() here, but we don't have to do this if we're using stdout or
	// stderr since it's not supported and will therefore cause an error.

	loggerMain := logger.Named("main")
	defer func() {
		loggerMain.Info("exiting",
			zap.NamedError("cause", ctx.Err()),
		)
	}()

	// Create the app and register it with the multiplexer.
	app := demoapp.New("Hello, World!", demoapp.WithZapLogger(logger.Named("app")))
	mux := http.NewServeMux()
	mux.Handle("/", app)

	// Configure standard HTTP server. Use the base context for connections.
	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			return baseCtx
		},
	}

	// Wrap standard HTTP server in a piece of code that handles shutdowns and errors more gracefully.
	graceful := gracefulhttp.New(ctx, server, gracefulhttp.WithZapLogger(logger.Named("gracefulhttp")))

	// Create the network listener. The listener will be closed by http.Server. If this wasn't the case, we would have
	// to defer the closing of the listener but this is done by http.Server already.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		loggerMain.Error("listener failed",
			zap.Error(err),
		)
		return -1
	}
	loggerMain.Info("binding listener",
		zap.String("addr", listener.Addr().String()),
		zap.String("network", listener.Addr().Network()),
	)

	// Start serving using the graceful wrapper.
	if err = graceful.GracefulServe(listener); err != nil {
		loggerMain.Error("http server returned an error",
			zap.Error(err),
		)
		return -2
	}

	return 0
}
