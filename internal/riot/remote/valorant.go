package remote

import (
	"fmt"

	"github.com/Cufee/valorant-account-tracker-go/internal/riot"
)

/* Lifted straight from docs - https://valapidocs.techchrism.me/endpoint/player-info */
type PlayerInfo struct {
	Country       string      `json:"country"`
	Sub           string      `json:"sub"`
	EmailVerified bool        `json:"email_verified"`
	PlayerLocale  interface{} `json:"player_plocale"`
	CountryAt     int64       `json:"country_at"`
	Pw            struct {
		CngAt     int64 `json:"cng_at"`
		Reset     bool  `json:"reset"`
		MustReset bool  `json:"must_reset"`
	} `json:"pw"`
	PhoneNumberVerified        bool        `json:"phone_number_verified"`
	AccountVerified            bool        `json:"account_verified"`
	Ppid                       interface{} `json:"ppid"`
	FederatedIdentityProviders []string    `json:"federated_identity_providers"`
	PlayerLocale2              string      `json:"player_locale"`
	Username                   string      `json:"username"`
	Acct                       struct {
		Type      int    `json:"type"`
		State     string `json:"state"`
		Adm       bool   `json:"adm"`
		GameName  string `json:"game_name"`
		TagLine   string `json:"tag_line"`
		CreatedAt int64  `json:"created_at"`
	} `json:"acct"`
	Age      int               `json:"age"`
	Jti      string            `json:"jti"`
	Affinity map[string]string `json:"affinity"`
}

func GetPlayerInfo(credentials riot.AccessTokens) (PlayerInfo, error) {
	var trimmedCredentials riot.AccessTokens
	trimmedCredentials.AccessToken = credentials.AccessToken

	var data PlayerInfo
	err := authedRequest("GET", "https://auth.riotgames.com/userinfo", &data, nil, trimmedCredentials)
	if err != nil {
		return PlayerInfo{}, err
	}

	return data, nil
}

type SeasonMMR struct {
	CompetitiveTier int    `json:"CompetitiveTier"`
	SeasonID        string `json:"SeasonID"`
	Rank            int    `json:"Rank"`
}

type PlayerMMR struct {
	QueueSkills map[string]struct {
		SeasonalInfoBySeasonID map[string]SeasonMMR `json:"SeasonalInfoBySeasonID"`
	} `json:"QueueSkills"`
	LatestCompetitiveUpdate struct {
		SeasonID                string `json:"SeasonID"`
		TierAfterUpdate         int    `json:"TierAfterUpdate"`
		RankedRatingAfterUpdate int    `json:"RankedRatingAfterUpdate"`
		RankedRatingEarned      int    `json:"RankedRatingEarned"`
	} `json:"LatestCompetitiveUpdate"`
}

func GetPlayerMMR(puuid string, credentials riot.AccessTokens, opts RequestOptions) (PlayerMMR, error) {
	var data PlayerMMR
	err := authedRequest("GET", fmt.Sprintf("https://pd.%s.a.pvp.net/mmr/v1/players/%s", opts.Shard, puuid), &data, nil, credentials, opts)
	if err != nil {
		return PlayerMMR{}, err
	}

	return data, nil
}

type Season struct {
	Content `json:"embedded"`
	Type    string `json:"Type"`
}

type Content struct {
	ID        string `json:"ID"`
	Name      string `json:"Name"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
	IsActive  bool   `json:"IsActive"`
}

type GameContent struct {
	Seasons []Season  `json:"Seasons"`
	Events  []Content `json:"Events"`
}

func GetGameContent(opts RequestOptions) (GameContent, error) {
	var data GameContent
	err := authedRequest("GET", fmt.Sprintf("https://shared.%s.a.pvp.net/content-service/v3/content", opts.Shard), &data, nil, riot.AccessTokens{}, opts)
	if err != nil {
		return GameContent{}, err
	}

	return data, nil
}
