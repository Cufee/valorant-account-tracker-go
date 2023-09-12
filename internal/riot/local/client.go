package local

import (
	"errors"
	"strings"
)

var ErrNoGameSession = errors.New("game session not found")
var ErrRiotClientNotRunning = errors.New("riot client is not running")

// https://valapidocs.techchrism.me/endpoint/player-mmr#client-platform
var clientPlatformValue = "ew0KCSJwbGF0Zm9ybVR5cGUiOiAiUEMiLA0KCSJwbGF0Zm9ybU9TIjogIldpbmRvd3MiLA0KCSJwbGF0Zm9ybU9TVmVyc2lvbiI6ICIxMC4wLjE5MDQyLjEuMjU2LjY0Yml0IiwNCgkicGxhdGZvcm1DaGlwc2V0IjogIlVua25vd24iDQp9"

type ClientRegionResponse struct {
	Region string `json:"region"`
}

func GetClientRegion() (string, error) {
	var data ClientRegionResponse
	err := request("/riotclient/region-locale", &data)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			// This seems to indicate that riot client was not started
			return "", ErrRiotClientNotRunning
		}
		return "", err
	}
	if data.Region == "" {
		return "", errors.New("api returned blank region")
	}
	return data.Region, nil
}

type Session struct {
	Configuration LaunchConfiguration `json:"launchConfiguration"`
	ProductId     string              `json:"productId"`
	ExitCode      int                 `json:"exitCode"`
	Version       string              `json:"version"`
	Phase         string              `json:"phase"`
}

type LaunchConfiguration struct {
	Arguments  []string `json:"arguments"`
	Executable string   `json:"executable"`
}

func GetSessions() ([]Session, error) {
	var data map[string]Session
	err := request("/product-session/v1/external-sessions", &data)
	if err != nil {
		return nil, err
	}

	var sessions []Session
	for _, s := range data {
		sessions = append(sessions, s)
	}
	return sessions, nil
}

type GameSession struct {
	Platform string
	Version  string
	Region   string
	Shard    string
}

func GetGameSessionInfo() (GameSession, error) {
	region, err := GetClientRegion()
	if err != nil {
		return GameSession{}, err
	}

	sessions, err := GetSessions()
	if err != nil {
		return GameSession{}, err
	}

	// Find a session of valorant that has a valid configuration endpoint
	for _, session := range sessions {
		if session.ProductId == "valorant" {
			for _, arg := range session.Configuration.Arguments {
				if strings.HasPrefix(arg, "-config-endpoint=") {
					// https://shared.na.a.pvp.net
					endpointSlice := strings.SplitN(strings.ReplaceAll(arg, "-config-endpoint=", ""), ".", 4)
					if len(endpointSlice) != 4 {
						return GameSession{}, ErrNoGameSession
					}
					return GameSession{
						Platform: clientPlatformValue,
						Version:  session.Version,
						Region:   region,
						Shard:    endpointSlice[1],
					}, nil
				}
			}
		}
	}

	return GameSession{}, ErrNoGameSession

}
