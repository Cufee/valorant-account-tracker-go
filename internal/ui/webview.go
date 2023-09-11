package ui

import (
	"log"

	"github.com/Cufee/valorant-account-tracker-go/internal/types"
	"github.com/Cufee/valorant-account-tracker-go/internal/ui/views"
	"github.com/jchv/go-webview2"
)

func startWebview(title, html string, width, height uint) {
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  title,
			Width:  width,
			Height: height,
			IconId: 2, // icon not implemented?
			Center: true,
		},
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()
	w.SetSize(800, 600, webview2.HintFixed)
	w.SetHtml(html)
	w.Run()
}

func OpenHomeView() error {
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

	go startWebview("Valorant Account Tracker", html, 800, 600)
	return nil
}
