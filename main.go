package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cufee/valorant-account-tracker-go/internal/logic"
	"github.com/Cufee/valorant-account-tracker-go/internal/riot/local"
	"github.com/Cufee/valorant-account-tracker-go/internal/ui"
)

func main() {
	// Init UI
	_, err := ui.StartWebserver()
	if err != nil {
		log.Panicf("failed to start a web server: %s", err)
	}

	log.Print("Webserver started")

	ui.RegisterSystrayIcon()

	err = ui.OpenAppInBrowser()
	if err != nil {
		log.Panicf("failed to start a web server: %s", err)
	}

	log.Print("UI started")

	// Init background logic
	go logic.ListenForAccountUpdates()

	log.Print("Background logic started")

	local.EventBus.Publish(local.TopicSocketMessage, map[string]interface{}{local.EventNameAuthorization: nil})

	// Wait for Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for range c {
		os.Exit(0)
	}
}
