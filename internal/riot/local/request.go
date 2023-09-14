package local

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ErrorResponseBody struct {
	HttpStatus int    `json:"httpStatus"`
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
}

var httpClient *http.Client
var ErrNilResponse = errors.New("response is blank")
var ErrResourceNotFound = errors.New("resource not found")

func init() {

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
	request.Header.Set("Authorization", credentials.AuthHeader)

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	if response == nil {
		return ErrNilResponse
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var errorData ErrorResponseBody
	err = json.Unmarshal(b, &errorData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body into error: %s", err)
	}
	if errorData.ErrorCode == "RESOURCE_NOT_FOUND" {
		return ErrResourceNotFound
	}
	if errorData.Message != "" {
		return fmt.Errorf("game api request '%s' returned an error: %s", path, errorData.Message)
	}

	return json.Unmarshal(b, target)

}
