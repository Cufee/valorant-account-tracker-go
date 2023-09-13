package ui

import (
	"encoding/json"
	"net"

	"github.com/Cufee/valorant-account-tracker-go/internal/database"
	"github.com/Cufee/valorant-account-tracker-go/internal/types"
	"github.com/Cufee/valorant-account-tracker-go/internal/ui/views"
	"github.com/gofiber/fiber/v2"
)

var webServerPort int

func StartWebserver() (int, error) {
	dbClient, err := database.GetClient()
	if err != nil {
		return 0, err
	}

	// Get a random available port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, nil
	}
	webServerPort = l.Addr().(*net.TCPAddr).Port

	// Start the Fiber server
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		accountsRaw, err := dbClient.List("accounts")
		if err != nil {
			return err
		}

		var accounts []types.Account
		for _, accountRaw := range accountsRaw {
			var account types.Account
			err = json.Unmarshal(accountRaw, &account)
			if err != nil {
				return err
			}
			accounts = append(accounts, account)
		}

		props := make(map[string]any)
		props["AccountsProps"] = views.AccountViewProps{
			Accounts: accounts,
		}

		html, err := views.Home.Render(props)
		if err != nil {
			return err
		}

		c.Response().Header.Set("Content-Type", "text/html")
		c.SendString(html)
		return nil
	})

	go app.Listener(l)
	return webServerPort, nil
}

func getWebserverPort() int {
	return webServerPort
}
