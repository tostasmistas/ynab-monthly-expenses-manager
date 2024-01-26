package backend

import (
	"context"
)

// App represents the Wails application
type App struct {
	Context context.Context
}

// SetupApp creates a new App instance
func SetupApp() *App {
	return &App{}
}

// Startup sets the application context
func (app *App) Startup(context context.Context) {
	app.Context = context
}
