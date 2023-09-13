package riot

/* Returned from a local game api and used to authenticate all requests to remote api */
type AccessTokens struct {
	AccessToken      string
	EntitlementToken string
}
