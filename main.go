package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	backendpkg "ynab-monthly-expenses-manager/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	app := backendpkg.SetupApp()

	backend := backendpkg.SetupBackend()

	wails.Run(&options.App{
		Title:         "YNAB Monthly Expenses Manager",
		Width:         1024,
		Height:        768,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(context context.Context) {
			app.Startup(context)
			backend.Startup(context)
		},
		OnDomReady: func(context context.Context) {
			backend.DomReady(context)
		},
		Bind: []interface{}{
			app,
			backend,
		},
	})
}
