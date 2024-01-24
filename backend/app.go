package backend

import (
	"context"
)

type App struct {
	Context context.Context
}

func SetupApp() *App {
	return &App{}
}

func (app *App) Startup(context context.Context) {
	app.Context = context
}
