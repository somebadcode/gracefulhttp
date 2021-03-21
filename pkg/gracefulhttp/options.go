package gracefulhttp

import (
	"go.uber.org/zap"
)

// Option is a function that can modify GracefulServer. This allows for a more flexible optional arguments.
type Option func(*GracefulServer)

// WithZapLogger sets app logger.
func WithZapLogger(l *zap.Logger) Option {
	return func(gs *GracefulServer) {
		gs.logger = l
	}
}