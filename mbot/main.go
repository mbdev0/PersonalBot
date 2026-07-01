package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"personal_bot/backend/infrastructure/tui"
	"personal_bot/backend/pkg/logger"
	"personal_bot/backend/server"
	"strconv"

	"log"

	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hpcloud/tail"
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
	if len(os.Args) > 1 && os.Args[1] == "tui" {
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			logger.Error(err)
			return
		}
		runTUI(id)
		return
	}

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

func runTUI(id int64) {
	exe, _ := os.Executable()
	filename := filepath.Join(filepath.Dir(exe), fmt.Sprintf("logs/task-%d-logs", id))
	tail, err := tail.TailFile(filename, tail.Config{
		Follow: true,
		ReOpen: true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: io.SeekStart,
		},
	})

	if err != nil {
		logger.Error(err)
		return
	}

	p := tea.NewProgram(tui.NewModel(tail))
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
