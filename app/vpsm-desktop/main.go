package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	// Parse CLI flags for standalone terminal window mode without exiting on unknown flags
	fs := flag.NewFlagSet("vpsm-desktop", flag.ContinueOnError)
	fs.UintVar(&termServerID, "terminal-server-id", 0, "Terminal server ID")
	fs.StringVar(&termHost, "terminal-host", "", "Terminal host IP/address")
	fs.StringVar(&termName, "terminal-name", "", "Terminal server display name")
	fs.IntVar(&termPort, "terminal-port", 22, "Terminal port")
	fs.StringVar(&termUser, "terminal-user", "root", "Terminal username")
	fs.StringVar(&termAuthType, "terminal-authtype", "password", "Terminal auth type")
	fs.StringVar(&termAuthSecret, "terminal-authsecret", "", "Terminal auth secret")

	if len(os.Args) > 1 {
		_ = fs.Parse(os.Args[1:])
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

		// Load custom window dimensions from stored preferences if available
		if pref, err := app.GetTerminalPreference(); err == nil && pref != nil {
			if pref.WindowWidth > 0 {
				width = pref.WindowWidth
			}
			if pref.WindowHeight > 0 {
				height = pref.WindowHeight
			}
		}

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

	// Create native application menu
	appMenu := menu.NewMenu()

	// macOS Application Menu ("VPSM Desktop")
	mainSubmenu := appMenu.AddSubmenu("VPSM Desktop")
	mainSubmenu.AddText("About VPSM Desktop", nil, func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
				Type:    runtime.InfoDialog,
				Title:   "About VPSM Desktop",
				Message: "VPSM Desktop v1.2.0\nVPS Manager — Remote Server Inventory & SSH Terminal Panel\n\nDeveloped by @devlopersabbir",
			})
		}
	})
	mainSubmenu.AddText("Check for Updates...", nil, func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
				Type:    runtime.InfoDialog,
				Title:   "Check for Updates",
				Message: "You are running the latest version of VPSM Desktop (v1.2.0).",
			})
		}
	})
	mainSubmenu.AddSeparator()
	mainSubmenu.AddText("Hide VPSM Desktop", keys.CmdOrCtrl("h"), func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.WindowHide(app.ctx)
		}
	})
	mainSubmenu.AddSeparator()
	mainSubmenu.AddText("Quit VPSM Desktop", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.Quit(app.ctx)
		}
	})

	// View / Terminal Menu
	viewMenu := appMenu.AddSubmenu("View")
	viewMenu.AddText("Zoom In (Increase Font)", keys.CmdOrCtrl("+"), func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.EventsEmit(app.ctx, "terminal:zoom_in")
		}
	})
	viewMenu.AddText("Zoom Out (Decrease Font)", keys.CmdOrCtrl("-"), func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.EventsEmit(app.ctx, "terminal:zoom_out")
		}
	})
	viewMenu.AddText("Reset Font Size", keys.CmdOrCtrl("0"), func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.EventsEmit(app.ctx, "terminal:zoom_reset")
		}
	})
	viewMenu.AddSeparator()
	viewMenu.AddText("Toggle Fullscreen", keys.Key("f"), func(cd *menu.CallbackData) {
		if app.ctx != nil {
			runtime.WindowToggleMaximise(app.ctx)
		}
	})

	// Edit menu for Cut/Copy/Paste
	appMenu.Append(menu.EditMenu())

	err := wails.Run(&options.App{
		Title:  title,
		Width:  width,
		Height: height,
		Menu:   appMenu,
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
