package logic

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Cufee/valorant-account-tracker-go/internal/database"
	"github.com/Cufee/valorant-account-tracker-go/internal/riot/local"
)

func RegisterAccountUpdateListener() {
	onSocketMessage := local.EventBus.Subscribe(local.TopicSocketMessage)
	go func() {
		for e := range onSocketMessage {
			data, ok := e.Data.(map[string]interface{})
			if !ok || data[local.EventNameAuthorization] == nil {
				continue
			}

			account, err := GetCurrentPlayerAccount()
			if err != nil {
				if errors.Is(err, local.ErrNoGameSession) {
					// no game session, no need to update anything
					return
				}
				log.Printf("Failed to get game session info: %s", err)
				return
			}

			encoded, err := json.Marshal(account)
			if err != nil {
				log.Printf("Failed to encode account: %s", err)
				return
			}

			db, err := database.GetClient()
			if err != nil {
				log.Printf("Failed to get database client: %s", err)
				return
			}
			err = db.Set("accounts", account.ID, encoded)
			if err != nil {
				log.Printf("Failed to set account: %s", err)
				return
			}

			log.Printf("Account updated: %s", account.Username)
		}
	}()
}