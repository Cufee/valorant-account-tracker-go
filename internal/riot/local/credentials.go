package local

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/rjeczalik/notify"
)

var lockfilePath = filepath.Join(os.Getenv("LocalAppData"), "Riot Games/Riot Client/Config/lockfile")
var credentialsCache *LocalCredentials

type LocalCredentials struct {
	HttpEndpoint string
	WssEndpoint  string
	AuthHeader   string
	raw          string
}

func init() {
	onChange := func(event notify.Event) {
		EventBus.Publish(TopicCredentialsChanged, event)
	}
	if _, err := watchFileChanges(lockfilePath, onChange, time.Second*1); err != nil {
		log.Panicf("Failed to register an event watcher for lockfile\n%s", err)
	}
}

func GetLocalCredentials() (*LocalCredentials, error) {
	if credentialsCache == nil {
		return nil, errors.New("credentials not available")
	}
	return credentialsCache, nil
}

func readLockfileCredentials(path string) (*LocalCredentials, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, os.ErrNotExist
	}

	// name:pid:port:password:protocol
	credsSlice := strings.SplitN(string(content), ":", 5)
	if len(credsSlice) != 5 {
		return nil, errors.New("invalid credentials inside lockfile")
	}

	var credentials LocalCredentials
	credentials.WssEndpoint = fmt.Sprintf("wss://127.0.0.1:%s", credsSlice[2])
	credentials.HttpEndpoint = fmt.Sprintf("%s://127.0.0.1:%s", credsSlice[4], credsSlice[2])
	encodedAuthHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("riot:%s", credsSlice[3])))
	credentials.AuthHeader = fmt.Sprintf("Basic %s", encodedAuthHeader)
	credentials.raw = string(content)

	return &credentials, nil
}

func watchFileChanges(path string, callback func(notify.Event), debounceTimeout time.Duration) (func(), error) {
	c := make(chan notify.EventInfo, 1)
	var debounce = debounce.New(debounceTimeout)

	if err := notify.Watch(filepath.Dir(path), c, notify.Create, notify.Write, notify.Remove); err != nil {
		return nil, err
	}

	go func() {
		for {
			info := <-c
			if info.Path() == path {
				debounce(func() { callback(info.Event()) })
			}
		}
	}()

	return func() { notify.Stop(c) }, nil
}

func updateCredentialsCache() (bool, error) {
	credentials, err := readLockfileCredentials(lockfilePath)
	if err != nil {
		credentialsCache = nil
		return true, err
	}
	if credentialsCache != nil && credentialsCache.raw == credentials.raw {
		return false, nil
	}
	credentialsCache = credentials
	return true, nil
}

func updateCredentialsCacheFromEvent(event notify.Event) {
	if event == notify.Remove {
		credentialsCache = nil
		EventBus.Publish(TopicCredentialsDeleted, nil)
		log.Print("Credentials cache deleted")
		return
	}

	if _, err := updateCredentialsCache(); err != nil {
		log.Printf("Failed to update credentials cache\n%s", err)
	}
	log.Print("Credentials cache updated")
}
