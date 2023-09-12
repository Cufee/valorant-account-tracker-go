package local

type EntitlementsTokenResponse struct {
	AccessToken  string        `json:"accessToken"`
	Entitlements []interface{} `json:"entitlements"`
	Issuer       string        `json:"issuer"`
	Subject      string        `json:"subject"`
	Token        string        `json:"token"`
}

type AccessTokens struct {
	AccessToken      string
	EntitlementToken string
}

func GetAccessTokens() (AccessTokens, error) {
	var data EntitlementsTokenResponse
	var tokens AccessTokens

	err := request("/entitlements/v1/token", &data)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = data.AccessToken
	tokens.EntitlementToken = data.Token
	return tokens, nil
}
