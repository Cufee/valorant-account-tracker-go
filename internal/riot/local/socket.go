package local

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rjeczalik/notify"
)

var openSocket *websocket.Conn

var (
	defaultEvents = []string{EventNameAuthorization}

	EventNameAllApiEvents  = "OnJsonApiEvent"
	EventNameAuthorization = "OnJsonApiEvent_rso-auth_v2_authorizations"

	/* A list of all events for reference */
	_ = []string{
		"AgentResourceEvent",
		"OnClientFlash",
		"OnClientFocus",
		"OnClientMinimize",
		"OnJsonApiEvent",
		"OnJsonApiEvent_agent_v1_requests",
		"OnJsonApiEvent_agent_v1_session",
		"OnJsonApiEvent_anti-addiction_v1_products",
		"OnJsonApiEvent_chat_v1_session",
		"OnJsonApiEvent_chat_v1_settings",
		"OnJsonApiEvent_chat_v4_friends",
		"OnJsonApiEvent_chat_v4_presences",
		"OnJsonApiEvent_chat_v5_messages",
		"OnJsonApiEvent_chat_v5_participants",
		"OnJsonApiEvent_chat_v6_conversations",
		"OnJsonApiEvent_chat_v6_friendrequests",
		"OnJsonApiEvent_client-config_v2_namespace",
		"OnJsonApiEvent_entitlements_v1_token",
		"OnJsonApiEvent_eula_v1_agreement",
		"OnJsonApiEvent_external-message-handler_v1_patch-request",
		"OnJsonApiEvent_ga-warning_v1_warnings",
		"OnJsonApiEvent_patch-proxy_v1_active-updates",
		"OnJsonApiEvent_patch-proxy_v1_patch-states",
		"OnJsonApiEvent_pay-mobile_v1_productListResult",
		"OnJsonApiEvent_player-account_aliases_v1",
		"OnJsonApiEvent_player-affinity_config_v1",
		"OnJsonApiEvent_player-reporting_v1_reporter-feedback",
		"OnJsonApiEvent_plugin-manager_v1_status",
		"OnJsonApiEvent_process-control_v1_process",
		"OnJsonApiEvent_product-integration_v1_app-update",
		"OnJsonApiEvent_product-integration_v1_locale",
		"OnJsonApiEvent_product-launcher_v1_launching_from_xbgp",
		"OnJsonApiEvent_product-metadata_v2_products",
		"OnJsonApiEvent_publishing-content_v1_news-feed",
		"OnJsonApiEvent_publishing-content_v1_promo",
		"OnJsonApiEvent_rc-auth_v1_xbgp",
		"OnJsonApiEvent_restriction_v1_launchRestrictedProducts",
		"OnJsonApiEvent_restriction_v1_launchRestrictions",
		"OnJsonApiEvent_restriction_v1_ready",
		"OnJsonApiEvent_riot-client-app-command_v1_uri-handler",
		"OnJsonApiEvent_riot-client-lifecycle-state_v1_state",
		"OnJsonApiEvent_riot-client-lifecycle_v1_league-region-election",
		"OnJsonApiEvent_riot-client-lifecycle_v1_product-context",
		"OnJsonApiEvent_riot-client-lifecycle_v1_ux-command",
		"OnJsonApiEvent_riot-messaging-service_v1_message",
		"OnJsonApiEvent_riot-messaging-service_v1_messages",
		"OnJsonApiEvent_riot-messaging-service_v1_out-of-sync",
		"OnJsonApiEvent_riot-messaging-service_v1_session",
		"OnJsonApiEvent_riot-messaging-service_v1_state",
		"OnJsonApiEvent_riot-messaging-service_v1_user",
		"OnJsonApiEvent_riotclient_affinity",
		"OnJsonApiEvent_riotclient_zoom-scale",
		"OnJsonApiEvent_riotclientapp_v1_isXbgpRunning",
		"OnJsonApiEvent_riotclientapp_v1_new-args",
		"OnJsonApiEvent_rnet-lifecycle_v1_league-region-election",
		"OnJsonApiEvent_rnet-lifecycle_v1_product-context",
		"OnJsonApiEvent_rnet-lifecycle_v1_product-context-phase",
		"OnJsonApiEvent_rnet-lifecycle_v2_ux-command",
		"OnJsonApiEvent_rnet-pft_v1_surveys",
		"OnJsonApiEvent_rnet-product-registry_v1_background-patching",
		"OnJsonApiEvent_rnet-product-registry_v1_install-states",
		"OnJsonApiEvent_rnet-product-registry_v1_move-install-states",
		"OnJsonApiEvent_rnet-product-registry_v4_active-updates",
		"OnJsonApiEvent_rnet-product-registry_v4_available-product-locales",
		"OnJsonApiEvent_rnet-product-registry_v4_install-settings",
		"OnJsonApiEvent_rnet-product-registry_v4_patch-states",
		"OnJsonApiEvent_rnet-product-registry_v4_player-products-state",
		"OnJsonApiEvent_rnet-product-registry_v4_priority-patch-requests",
		"OnJsonApiEvent_rnet-product-registry_v4_products",
		"OnJsonApiEvent_rnet-product-registry_v4_public-products-state",
		"OnJsonApiEvent_rnet-self-update_v1_status",
		"OnJsonApiEvent_rso-auth_configuration_v3",
		"OnJsonApiEvent_rso-auth_v1_userinfo",
		"OnJsonApiEvent_rso-auth_v2_authorizations",
		"OnJsonApiEvent_startup-config_v1_registry-config",
		"OnJsonApiEvent_vanguard_v1_status",
		"OnJsonApiEvent_voice-chat_v1_audio-properties",
		"OnJsonApiEvent_voice-chat_v2_devices",
		"OnJsonApiEvent_voice-chat_v4_sessions",
		"OnLcdsEvent",
		"OnRegionLocaleChanged",
		"OnServiceProxyAsyncEvent",
		"OnServiceProxyMethodEvent",
		"OnServiceProxyUuidEvent",
	}
)

func newSocketConnection(events ...string) error {
	if len(events) == 0 {
		events = defaultEvents
	}

	if openSocket != nil {
		openSocket.Close()
		openSocket = nil
		EventBus.Publish(TopicSocketClosed, nil)
	}

	credentials, err := GetLocalCredentials()
	if err != nil {
		return err
	}

	dialer := websocket.Dialer{}
	dialer.TLSClientConfig = &tls.Config{
		RootCAs: caCertPool,
	}

	header := http.Header{}
	header.Set("Authorization", credentials.AuthHeader)

	ws, _, err := dialer.Dial(credentials.WssEndpoint, header)
	if err != nil {
		log.Printf("Failed to connect to websocket: %s", err)
		return err
	}

	// Register events to listen to
	for _, event := range events {
		var payload []interface{} = []interface{}{5, event}
		err = ws.WriteJSON(payload)
		if err != nil {
			log.Printf("Failed to write to websocket: %s", err)
			openSocket = nil
			ws.Close()
			return err
		}
	}

	openSocket = ws
	EventBus.Publish(TopicSocketConnected, nil)

	go func() {
		for {
			messageType, reader, err := ws.NextReader()
			if err != nil {
				log.Printf("Failed to get websocket reader: %s", err)
				return
			}
			if messageType == websocket.CloseMessage {
				openSocket = nil
				EventBus.Publish(TopicSocketClosed, nil)
				return
			}

			message, err := io.ReadAll(reader)
			if err != nil {
				log.Printf("Failed to read websocket message: %s", err)
				continue
			}

			if len(message) == 0 {
				continue
			}

			// [code, name, payload]
			var data []interface{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Printf("Failed to unmarshal websocket message: %s", err)
				continue
			}
			if len(data) != 3 {
				log.Printf("Invalid websocket message: %s", message)
				continue
			}

			EventBus.Publish(TopicSocketMessage, map[string]interface{}{fmt.Sprint(data[1]): data[2]})
			log.Printf("Received websocket message: %s", data[1])
		}
	}()

	return nil
}

func updateSocketFromEvent(event notify.Event) {
	if event == notify.Remove {
		if openSocket != nil {
			openSocket.Close()
			openSocket = nil
			EventBus.Publish(TopicSocketClosed, nil)
		}
		return
	}

	err := newSocketConnection()
	if err != nil {
		log.Printf("Failed to connect to websocket: %s", err)
		return
	}
	log.Print("Connected to websocket")
}
