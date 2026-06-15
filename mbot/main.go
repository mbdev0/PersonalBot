package main

import (
	"context"
	"embed"
	"personal_bot/backend/pkg/logger"
	"personal_bot/backend/server"

	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	server, err := server.New(ctx)
	if err != nil {
		logger.Error(err)
	}
	go server.Start()

	app := application.New(application.Options{
		Name:        "mbot",
		Description: "A demo of using raw HTML & CSS",
		Services:    []application.Service{},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		}, OnShutdown: func() { server.Shutdown(ctx); cancel() },
	})

	macTitleBar := application.MacTitleBar{
		AppearsTransparent:   false,
		Hide:                 false,
		HideTitle:            false,
		FullSizeContent:      true,
		UseToolbar:           false,
		HideToolbarSeparator: false,
	}

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Personal Bot",
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
			TitleBar: macTitleBar,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "http://localhost:9245",
	})

	// Run the application. This blocks until the application has been exited.
	err = app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
