package ui

import (
	"log"
	"net"

	"github.com/Cufee/valorant-account-tracker-go/internal/types"
	"github.com/Cufee/valorant-account-tracker-go/internal/ui/views"
	"github.com/gofiber/fiber/v2"
)

var webServerPort int

func StartWebserver() error {
	// Get a random available port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil
	}
	webServerPort = l.Addr().(*net.TCPAddr).Port
	log.Print(webServerPort)

	// Start the Fiber server
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		var accounts views.AccountViewProps = views.AccountViewProps{
			Accounts: []types.Account{{
				Tag:      "0000",
				Name:     "NameHere",
				Username: "Username",
				LastRank: types.Rank{
					Icon: "",
				},
			}, {
				Tag:      "0001",
				Name:     "Name 2",
				Username: "username 2",
				LastRank: types.Rank{
					Icon: "",
				},
			}},
		}

		props := make(map[string]any)
		props["AccountsProps"] = accounts

		html, err := views.Home.Render(props)
		if err != nil {
			return err
		}

		c.Response().Header.Set("Content-Type", "text/html")
		c.SendString(html)
		return nil
	})

	go app.Listener(l)
	return nil
}

func getWebserverPort() int {
	return webServerPort
}
