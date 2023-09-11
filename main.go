package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Cufee/valorant-account-tracker-go/internal/ui"
)

func main() {
	// Create a tray icon and open the UI
	ui.RegisterSystrayIcon()
	ui.OpenHomeView()

	// Wait for Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for range c {
		os.Exit(0)
	}
}
