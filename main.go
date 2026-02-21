package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Worker mode: headless subprocess for parallel PDF rendering
	if len(os.Args) > 1 && os.Args[1] == "--worker" {
		runWorker(os.Args[2:])
		return
	}

	// CLI mode: command-line conversion without GUI
	if len(os.Args) > 1 && os.Args[1] == "--cli" {
		runCLI(os.Args[2:])
		return
	}

	// GUI mode
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "PDF2Image",
		Width:            1280,
		Height:           800,
		MinWidth:         960,
		MinHeight:        600,
		Frameless:        true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 24, G: 24, B: 27, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
