package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"simonpartner/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	svc := app.New()

	err := wails.Run(&options.App{
		Title:  "Simon Desktop",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 13, G: 13, B: 13, A: 1},
		OnStartup:        svc.Startup,
		OnShutdown:       svc.Shutdown,
		Bind: []interface{}{
			svc,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
