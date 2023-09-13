package ui

import (
	"fmt"

	"github.com/pkg/browser"
)

func OpenAppInBrowser() error {
	return browser.OpenURL(fmt.Sprintf("http://127.0.0.1:%v", getWebserverPort()))
}
