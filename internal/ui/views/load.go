package views

import (
	_ "embed"

	"github.com/Cufee/valorant-account-tracker-go/internal/types"
)

type AccountViewProps struct {
	Accounts []types.Account
}

//go:embed accounts.gohtml
var accounts string
var Accounts View

//go:embed home.gohtml
var home string
var Home View

func init() {
	var err error
	Accounts, err = newView("accounts", accounts)
	if err != nil {
		panic(err)
	}

	Home, err = newView("home", home, Accounts)
	if err != nil {
		panic(err)
	}
}
