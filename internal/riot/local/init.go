package local

import (
	"log"
	"reflect"

	"github.com/rjeczalik/notify"
)

/* Single place to init multiple events/function in order to avoid race conditions */
func init() {
	// Register a credentials update listener
	onCredentialsUpdate := EventBus.Subscribe(TopicCredentialsChanged)
	go func() {
		for e := range onCredentialsUpdate {
			event, ok := e.Data.(notify.Event)
			if !ok {
				log.Printf("Failed to cast event of type: %v", reflect.TypeOf(e.Data))
				continue
			}
			updateCredentialsCacheFromEvent(event)
			updateSocketFromEvent(event)
		}
	}()

	// Initial load
	EventBus.Publish(TopicCredentialsChanged, notify.Create)

}
