package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cufee/valorant-account-tracker-go/internal/ui"
)

func main() {
	_, err := ui.StartWebserver()
	if err != nil {
		log.Panicf("failed to start a web server: %s", err)
	}

	ui.RegisterSystrayIcon()

	err = ui.OpenAppInBrowser()
	if err != nil {
		log.Panicf("failed to start a web server: %s", err)
	}

	// Wait for Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for range c {
		os.Exit(0)
	}
}
