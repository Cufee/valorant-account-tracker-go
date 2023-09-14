package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Cufee/valorant-account-tracker-go/internal/logic"
	"github.com/Cufee/valorant-account-tracker-go/internal/ui"
)

func main() {
	// Init UI
	go ui.RegisterSystrayIcon()
	go func() {
		ui.StartWebserver()
		ui.OpenAppInBrowser()
	}()

	// Init background logic
	go logic.ListenForAccountUpdates()

	// Wait for Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for range c {
		os.Exit(0)
	}
}
