package local

import (
	"github.com/fsnotify/fsnotify"
)

/* Single place to init multiple events/function in order to avoid race conditions */
func init() {
	// Register a credentials update listener
	onCredentialsUpdate := EventBus.Subscribe(TopicCredentialsChanged)
	go func() {
		for e := range onCredentialsUpdate {
			updateCredentialsCacheFromEvent(e.Data)
			updateSocketFromEvent(e.Data)
		}
	}()

	// Initial load
	EventBus.Publish(TopicCredentialsChanged, fsnotify.Create)

}
