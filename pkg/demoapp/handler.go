package demoapp

import (
	"embed"
	"go.uber.org/zap"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed templates
var content embed.FS

type App struct {
	logger   *zap.Logger
	template *template.Template
	data     map[string]string
}

func New(msg string, options ...Option) *App {
	app := &App{
		data: map[string]string{
			"title":   "Demonstration",
			"message": msg,
		},
		logger: zap.NewNop(),
	}

	// Apply optional arguments.
	for _, fn := range options {
		fn(app)
	}

	// Get templates/ directory listing.
	tmplDir, err := fs.Sub(content, "templates")
	if err != nil {
		app.logger.Panic("failed to find embedded directory",
			zap.Error(err),
		)
	}

	// Parse all templates and panic if it fails.
	app.template, err = template.ParseFS(tmplDir, "*")
	if err != nil {
		app.logger.Panic("failed to parse embedded templates",
			zap.Error(err),
		)
	}

	return app
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// For demonstrative purpose, log each request.
	app.logger.Debug("incoming request",
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("remote", r.RemoteAddr),
	)

	// For demonstrative purpose, only allow the GET method.
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		app.logger.Debug("client used a disallowed method",
			zap.String("remote", r.RemoteAddr),
			zap.String("method", r.Method),
		)
		return
	}

	// Respond with 200.
	w.WriteHeader(http.StatusOK)

	// Write the app's message.
	if err := app.template.ExecuteTemplate(w, "index.html", app.data); err != nil {
		app.logger.Error("failed to write response",
			zap.String("remote", r.RemoteAddr),
			zap.Error(err),
		)
	}
}
