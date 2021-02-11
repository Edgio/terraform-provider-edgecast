// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	userAgent string = "verizonmedia/terraform:1.0.0"
)

var (
	idsToken *IDSToken
)

// IDSToken holds the OAuth 2.0 token for calling Verizon Media APIs
type IDSToken struct {
	AccessToken    string
	ExpirationTime time.Time
}

// ApiClient -
type ApiClient struct {
	BaseUrl         *url.URL
	IdsUrl          *url.URL
	UserAgent       string
	Token           string
	IdsClientId     string
	IdsClientSecret string
	IdsScope        string

	HttpClient *retryablehttp.Client
}

// NewApiClient
func NewApiClient(apiBaseUri string, idsUri string, apiToken string, idsClientId string, idsClientSecret string, idsScope string) (*ApiClient, error) {
	baseUrl, err := url.Parse(apiBaseUri)

	if err != nil {
		return nil, fmt.Errorf("NewApiClient: Parse API Base URL: %v", err)
	}

	idsUrl, err := url.Parse(idsUri)

	if err != nil {
		return nil, fmt.Errorf("NewApiClient: Parse IDS URL: %v", err)
	}

	return &ApiClient{
		BaseUrl:         baseUrl,
		IdsUrl:          idsUrl,
		Token:           apiToken,
		IdsClientId:     idsClientId,
		IdsClientSecret: idsClientSecret,
		IdsScope:        idsScope,
		HttpClient:      retryablehttp.NewClient(),
	}, nil

}

// BuildRequest creates a new Request for a Verizon Media API, adding appropriate headers
func (apiClient *ApiClient) BuildRequest(method, path string, body interface{}, isUsingIdsToken bool) (*retryablehttp.Request, error) {
	relativeURL, err := url.Parse(path)

	if err != nil {
		return nil, fmt.Errorf("BuildRequest: url.Parse: %v", err)
	}

	absoluteURL := apiClient.BaseUrl.ResolveReference(relativeURL)

	var payload interface{}

	if body != nil {
		switch body.(type) {
		case string:
			payload = []byte(body.(string))
		default:
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
			payload = buf
		}
	}

	req, err := retryablehttp.NewRequest(method, absoluteURL.String(), payload)

	if err != nil {
		return nil, fmt.Errorf("BuildRequest: NewRequest: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	if isUsingIdsToken {
		idsToken, err := apiClient.GetIdsToken()

		if err != nil {
			return nil, fmt.Errorf("BuildRequest: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+idsToken)
	} else {
		if len(apiClient.Token) == 0 {
			return nil, errors.New("BuildRequest: API Token is required")
		}

		req.Header.Set("Authorization", "TOK:"+apiClient.Token)
	}

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// SendRequest sends an HTTP request and, if applicable, sets the response to parsedResponse
func (apiClient *ApiClient) SendRequest(req *retryablehttp.Request, parsedResponse interface{}) (*http.Response, error) {
	log.Printf("[INFO] SendRequest >> [%s] %s", req.Method, req.URL.String())
	resp, err := apiClient.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("SendRequest: Do: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf("SendRequest: ioutil.ReadAll: %v", err)
		}

		bodyAsString := string(body)
		return nil, fmt.Errorf("SendRequest failed: %s", bodyAsString)
	}

	if parsedResponse != nil {
		err = json.NewDecoder(resp.Body).Decode(parsedResponse)

		if err != nil {
			return nil, fmt.Errorf("SendRequest: Decode error: %v", err)
		}
	}

	return resp, nil
}

// GetIdsToken returns the cached token, refreshing it if it has expired
func (apiClient *ApiClient) GetIdsToken() (string, error) {

	if len(apiClient.IdsClientId) == 0 || len(apiClient.IdsClientSecret) == 0 || len(apiClient.IdsScope) == 0 {
		return "", errors.New("GetIdsToken: IDS Client ID, Secret, and Scope required")
	}

	if idsToken == nil || idsToken.ExpirationTime.Before(time.Now()) {
		data := url.Values{}
		data.Set("grant_type", "client_credentials")
		data.Add("scope", apiClient.IdsScope)
		data.Add("client_id", apiClient.IdsClientId)
		data.Add("client_secret", apiClient.IdsClientSecret)

		idsTokenEndpoint := fmt.Sprintf("%s/connect/token", apiClient.IdsUrl.String())
		log.Printf("IDS TOKEN ENDPOINDT: %s", idsTokenEndpoint)
		newTokenRequest, err := http.NewRequest("POST", idsTokenEndpoint, bytes.NewBufferString(data.Encode()))

		if err != nil {
			return "", fmt.Errorf("GetIdsToken: NewRequest: %v", err)
		}

		newTokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		newTokenRequest.Header.Add("Cache-Control", "no-cache")
		newTokenRequest.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		log.Printf("[INFO] GetIdsToken: [POST] %s", idsTokenEndpoint)
		httpClient := &http.Client{}
		newTokenResponse, err := httpClient.Do(newTokenRequest)
		if err != nil {
			return "", fmt.Errorf("GetIdsToken: Do: %v", err)
		}

		var tokenMap map[string]interface{}
		err = json.NewDecoder(newTokenResponse.Body).Decode(&tokenMap)
		if err != nil {
			return "", fmt.Errorf("GetIdsToken: Decode: %v", err)
		}

		expiresIn := time.Second * time.Duration((tokenMap["expires_in"].(float64)))

		idsToken = &IDSToken{
			AccessToken:    tokenMap["access_token"].(string),
			ExpirationTime: time.Now().Add(expiresIn),
		}
	}

	return idsToken.AccessToken, nil
}

// FormatURLAddPartnerID is a utility function for adding the optional partner ID query string param
func FormatURLAddPartnerID(originalURL string, partnerID *int) string {
	if partnerID != nil {
		return originalURL + fmt.Sprintf("&partnerid=%d", *partnerID)
	}

	return originalURL
}
