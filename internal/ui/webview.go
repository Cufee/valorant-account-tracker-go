package ui

import (
	"fmt"
	"log"

	"github.com/jchv/go-webview2"
)

var openWindow webview2.WebView

func startWebview(title, url string, width, height uint) {
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
	openWindow = w
	defer func() {
		openWindow = nil
		w.Destroy()
	}()
	w.SetSize(800, 600, webview2.HintFixed)
	w.Navigate(url)
	w.Run()
}

func OpenAppWindow() error {
	if openWindow != nil {
		return nil
	}
	go startWebview("Valorant Account Tracker", fmt.Sprintf("http://127.0.0.1:%v", getWebserverPort()), 800, 600)
	return nil
}
