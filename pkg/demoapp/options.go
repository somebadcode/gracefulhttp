package demoapp

import (
	"go.uber.org/zap"
)

type Option func(*App)

// WithZapLogger sets app logger.
func WithZapLogger(l *zap.Logger) Option {
	return func(app *App) {
		app.logger = l
	}
}