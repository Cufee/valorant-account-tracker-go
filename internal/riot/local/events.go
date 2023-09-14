package local

import "github.com/Cufee/valorant-account-tracker-go/internal/events"

var (
	EventBus = events.NewBus()

	TopicCredentialsChanged = "riot:local:credentials_changed"
	TopicCredentialsDeleted = "riot:local:credentials_deleted"

	TopicSocketConnected = "riot:local:socket_connected"
	TopicSocketMessage   = "riot:local:socket_message"
	TopicSocketClosed    = "riot:local:socket_closed"
	TopicSocketError     = "riot:local:socket_error"
)

type (
	SocketMessage interface{}
	SocketError   error
)
