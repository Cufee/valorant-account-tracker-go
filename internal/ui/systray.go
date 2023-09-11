package ui

import (
	"os"

	"github.com/Cufee/valorant-account-tracker-go/internal/ui/icons"
	"github.com/getlantern/systray"
)

func RegisterSystrayIcon() {
	go systray.Run(onReady, nil)
}

func onReady() {
	systray.SetTitle("Account Tracker")
	systray.SetTooltip("Valorant Account Tracker")
	systray.SetIcon(icons.V)

	show := systray.AddMenuItem("Show UI", "Shows the app UI")
	go func() {
		for range show.ClickedCh {
			OpenHomeView()
		}
	}()

	quit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for range quit.ClickedCh {
			os.Exit(0)
		}
	}()
}
