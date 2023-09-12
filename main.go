package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Cufee/valorant-account-tracker-go/internal/ui"
)

func main() {
	ui.StartWebserver()

	ui.RegisterSystrayIcon()
	// ui.OpenAppWindow()

	// session, err := local.GetGameSessionInfo()
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Printf("%#v", session)

	// Wait for Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for range c {
		os.Exit(0)
	}
}
