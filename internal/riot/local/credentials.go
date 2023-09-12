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
	WssEndpoint    string
	HttpEndpoint   string
	HttpAuthHeader string
}

/* Load credentials and register a watcher to update cache on lockfile changes */
func init() {
	log.Print("Loading local game credentials")
	credentials, err := readLockfileCredentials(lockfilePath)
	if err != nil {
		log.Panicf("Failed to load local game credentials\n%s", err)
	}
	credentialsCache = credentials

	_, err = watchLockfileFileChanges(lockfilePath)
	if err != nil {
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
	credentials.WssEndpoint = fmt.Sprintf("wss://riot:%s@127.0.0.1:%s", credsSlice[3], credsSlice[2])
	credentials.HttpEndpoint = fmt.Sprintf("%s://127.0.0.1:%s", credsSlice[4], credsSlice[2])
	encodedAuthHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("riot:%s", credsSlice[3])))
	credentials.HttpAuthHeader = fmt.Sprintf("Basic %s", encodedAuthHeader)

	return &credentials, nil
}

func watchLockfileFileChanges(path string) (func() error, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func(events chan fsnotify.Event) {
		for event := range events {
			switch {
			case event.Op.Has(fsnotify.Remove):
				credentialsCache = nil
			default:
				log.Print("lockfile changed")
				credentials, err := readLockfileCredentials(path)
				if err != nil {
					return
				}
				credentialsCache = credentials
			}
		}
	}(watcher.Events)

	if err := watcher.Add(path); err != nil {
		watcher.Close()
		return nil, err
	}

	return watcher.Close, nil
}
