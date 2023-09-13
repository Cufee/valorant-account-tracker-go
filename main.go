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
	ui.OpenAppWindow()

	// account, err := logic.GetCurrentPlayerAccount()
	// if err != nil {
	// 	log.Panicf("Failed to get game session info: %s", err)
	// }

	// db, err := database.GetClient()
	// if err != nil {
	// 	log.Panicf("Failed to get database client: %s", err)
	// }

	// encoded, err := json.Marshal(account)
	// if err != nil {
	// 	log.Panicf("Failed to encode account: %s", err)
	// }
	// err = db.Set("accounts", account.ID, encoded)
	// if err != nil {
	// 	log.Panicf("Failed to set account: %s", err)
	// }

	// var decoded types.Account
	// err = db.GetEncoded("accounts", account.ID, &decoded, json.Unmarshal)
	// if err != nil {
	// 	log.Panicf("Failed to get account: %s", err)
	// }

	// log.Printf("Decoded: %+v", decoded)

	// Wait for Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for range c {
		os.Exit(0)
	}
}
