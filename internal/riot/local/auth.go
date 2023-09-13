package local

import "github.com/Cufee/valorant-account-tracker-go/internal/riot"

type EntitlementsTokenResponse struct {
	AccessToken  string        `json:"accessToken"`
	Entitlements []interface{} `json:"entitlements"`
	Issuer       string        `json:"issuer"`
	Subject      string        `json:"subject"`
	Token        string        `json:"token"`
}

func GetAccessTokens() (riot.AccessTokens, error) {
	var data EntitlementsTokenResponse
	var tokens riot.AccessTokens

	err := request("/entitlements/v1/token", &data)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = data.AccessToken
	tokens.EntitlementToken = data.Token
	return tokens, nil
}
