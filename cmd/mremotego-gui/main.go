package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/yourusername/mremotego/internal/config"
	"github.com/yourusername/mremotego/internal/gui"
)

func main() {
	// Set locale to avoid warnings on non-standard systems
	if os.Getenv("LANG") == "" || os.Getenv("LANG") == "C" {
		os.Setenv("LANG", "en_US.UTF-8")
	}

	// Create Fyne application
	myApp := app.NewWithID("com.mremotego.app")
	myApp.Settings().SetTheme(&customTheme{})

	// Set application icon (ignore errors - icon is optional)
	if icon := gui.GetAppIcon(); icon != nil {
		myApp.SetIcon(icon)
	}

	// Get config path
	cfgPath, err := config.GetDefaultConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config path: %v\n", err)
		os.Exit(1)
	}

	// Create config manager
	manager := config.NewManager(cfgPath)

	// Load config (create if doesn't exist)
	if err := manager.Load(); err != nil {
		// If config doesn't exist, show welcome dialog
		w := myApp.NewWindow("MremoteGO - Welcome")
		dialog.ShowConfirm("Welcome to MremoteGO",
			"No configuration file found. Would you like to create one?",
			func(create bool) {
				if create {
					if err := manager.Save(); err != nil {
						dialog.ShowError(fmt.Errorf("Failed to create config: %w", err), w)
						return
					}
					w.Close()
					showMainWindow(myApp, manager)
				} else {
					myApp.Quit()
				}
			}, w)
		w.ShowAndRun()
		return
	}

	showMainWindow(myApp, manager)
}

func showMainWindow(myApp fyne.App, manager *config.Manager) {
	mainWindow := gui.NewMainWindow(myApp, manager)
	mainWindow.Show()
}
