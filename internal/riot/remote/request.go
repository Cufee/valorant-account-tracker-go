package remote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Cufee/valorant-account-tracker-go/internal/riot"
)

type ErrorResponseBody struct {
	HttpStatus int    `json:"httpStatus"`
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
}

type RequestOptions struct {
	ContentType string

	ClientPlatform string
	ClientVersion  string
	Region         string
	Shard          string
}

var ErrNilResponse = errors.New("response is blank")
var ErrResourceNotFound = errors.New("resource not found")

func authedRequest(method, url string, target interface{}, body []byte, credentials riot.AccessTokens, o ...RequestOptions) error {
	opts := RequestOptions{}
	if len(o) == 1 {
		opts = o[0]
	}

	var bodyReader bytes.Reader
	if body != nil {
		bodyReader = *bytes.NewReader(body)
	}

	request, err := http.NewRequest(method, url, &bodyReader)
	if err != nil {
		return err
	}
	if body != nil {
		request.Header.Set("Content-Type", stringOr(opts.ContentType, "application/json"))
	}
	if credentials.AccessToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
	}
	if credentials.EntitlementToken != "" {
		request.Header.Set("X-Riot-Entitlements-JWT", credentials.EntitlementToken)
	}
	request.Header.Set("X-Riot-ClientPlatform", opts.ClientPlatform)
	request.Header.Set("X-Riot-ClientVersion", opts.ClientVersion)

	response, err := http.DefaultClient.Do(request)
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
		return fmt.Errorf("riot api request '%s' returned an error: %s", url, errorData.Message)
	}

	return json.Unmarshal(b, target)
}

func plainRequest(method, url string, target interface{}) error {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
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

	return json.Unmarshal(b, target)
}

func stringOr(values ...string) string {
	for _, v := range values {
		if len(v) > 0 {
			return v
		}
	}
	return ""
}
