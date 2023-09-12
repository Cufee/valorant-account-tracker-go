package local

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var httpClient *http.Client
var ErrNilResponse = errors.New("response is blank")
var ErrNotJSON = errors.New("response type is not json")

func init() {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(riotCertificate)
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
}

func request(path string, target interface{}) error {
	credentials, err := GetLocalCredentials()
	if err != nil {
		return err
	}

	requestUrl, err := url.JoinPath(credentials.HttpEndpoint, path)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", credentials.HttpAuthHeader)

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	if response == nil {
		return ErrNilResponse
	}
	if response.Header.Get("Content-Type") != "application/json" {
		return ErrNotJSON
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, target)

}
