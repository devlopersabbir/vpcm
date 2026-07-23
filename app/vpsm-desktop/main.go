package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	termServerID   uint
	termHost       string
	termName       string
	termPort       int
	termUser       string
	termAuthType   string
	termAuthSecret string
)

func main() {
	// Parse CLI flags for standalone terminal window mode
	flag.UintVar(&termServerID, "terminal-server-id", 0, "Terminal server ID")
	flag.StringVar(&termHost, "terminal-host", "", "Terminal host IP/address")
	flag.StringVar(&termName, "terminal-name", "", "Terminal server display name")
	flag.IntVar(&termPort, "terminal-port", 22, "Terminal port")
	flag.StringVar(&termUser, "terminal-user", "root", "Terminal username")
	flag.StringVar(&termAuthType, "terminal-authtype", "password", "Terminal auth type")
	flag.StringVar(&termAuthSecret, "terminal-authsecret", "", "Terminal auth secret")

	// Avoid breaking when standard args are passed
	if len(os.Args) > 1 {
		flag.CommandLine.Parse(os.Args[1:])
	}

	app := NewApp()

	title := "VPSM - VPS Manager"
	width := 1024
	height := 768

	if termHost != "" || termServerID > 0 {
		if termName != "" {
			title = fmt.Sprintf("VPSM Terminal — %s (%s@%s)", termName, termUser, termHost)
		} else if termHost != "" {
			title = fmt.Sprintf("VPSM Terminal — %s@%s", termUser, termHost)
		} else {
			title = fmt.Sprintf("VPSM Terminal — Server #%d", termServerID)
		}
		width = 900
		height = 600
		app.isTerminalOnly = true
		app.terminalServerID = termServerID
		app.terminalParams = SSHConnectionParams{
			Host:       termHost,
			Port:       termPort,
			Username:   termUser,
			AuthType:   termAuthType,
			AuthSecret: termAuthSecret,
		}
	}

	err := wails.Run(&options.App{
		Title:  title,
		Width:  width,
		Height: height,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 13, G: 17, B: 23, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
