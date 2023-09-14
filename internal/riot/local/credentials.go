package local

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var lockfilePath = filepath.Join(os.Getenv("LocalAppData"), "Riot Games/Riot Client/Config/lockfile")
var credentialsCache *LocalCredentials

type LocalCredentials struct {
	HttpEndpoint string
	WssEndpoint  string
	AuthHeader   string
}

func init() {
	onChange := func(event fsnotify.Event) {
		EventBus.Publish(TopicCredentialsChanged, event.Op)
	}
	if _, err := watchFileChanges(lockfilePath, onChange); err != nil {
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

	return &credentials, nil
}

func watchFileChanges(path string, callback func(fsnotify.Event)) (func() error, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func(events chan fsnotify.Event) {
		for event := range events {
			switch {
			default:
				log.Printf("lockfile changed: %s", event.Op.String())
				callback(event)
			}
		}
	}(watcher.Events)

	if err := watcher.Add(path); err != nil {
		watcher.Close()
		return nil, err
	}

	return watcher.Close, nil
}

func updateCredentialsCache() error {
	credentials, err := readLockfileCredentials(lockfilePath)
	if err != nil {
		credentialsCache = nil
		return err
	}
	credentialsCache = credentials
	return nil
}

func updateCredentialsCacheFromEvent(data interface{}) {
	op, ok := data.(fsnotify.Op)
	if !ok {
		log.Printf("Invalid data type received for credentials update: %T", data)
		return
	}

	if op == fsnotify.Remove {
		credentialsCache = nil
		EventBus.Publish(TopicCredentialsDeleted, nil)
		log.Print("Credentials cache deleted")
		return
	}

	if err := updateCredentialsCache(); err != nil {
		log.Printf("Failed to update credentials cache\n%s", err)
	}
	log.Print("Credentials cache updated")
}
