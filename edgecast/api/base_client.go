// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"terraform-provider-edgecast/edgecast/helper"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Deprecated - remove once migrated over to SDK
	userAgent string = "edgecast/terraform:1.0.0"
)

var (
	idsToken *IDSToken
)

// IDSToken holds the OAuth 2.0 token for calling Edgecast APIs
type IDSToken struct {
	AccessToken    string
	ExpirationTime time.Time
}

/* TODO: Rename ClientConfig to ProviderConfig and move to provider.go */

// ClientConfig holds configuration values for the provider.
type ClientConfig struct {
	APIToken         string `json:"-"` // sensitive.
	AccountNumber    string
	IdsClientID      string `json:"-"` // sensitive.
	IdsClientSecret  string `json:"-"` // sensitive.
	IdsScope         string `json:"-"` // sensitive.
	IDSAddress       string
	APIAddress       string
	APIAddressLegacy string
	APIURL           *url.URL
	IdsURL           *url.URL
	APIURLLegacy     *url.URL
	PartnerID        int
	PartnerUserID    int
	UserAgent        string
}

/*  TODO: Rename to ExpandProviderConfig and move to provider.go */

// ExpandClientConfig reads ClientConfig using the TF Resource Data.
func ExpandClientConfig(d *schema.ResourceData) (*ClientConfig, error) {
	config := &ClientConfig{
		APIToken:         d.Get("api_token").(string),
		AccountNumber:    d.Get("account_number").(string),
		IdsClientID:      d.Get("ids_client_id").(string),
		IdsClientSecret:  d.Get("ids_client_secret").(string),
		IdsScope:         d.Get("ids_scope").(string),
		IDSAddress:       d.Get("ids_address").(string),
		APIAddress:       d.Get("api_address").(string),
		APIAddressLegacy: d.Get("api_address_legacy").(string),
	}

	var err error

	config.IdsURL, err = url.Parse(config.IDSAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IDS URL: %w", err)
	}

	config.APIURL, err = url.Parse(config.APIAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	config.APIURLLegacy, err = url.Parse(config.APIAddressLegacy)
	if err != nil {
		return nil, fmt.Errorf("failed to parse legacy API URL: %w", err)
	}

	if partnerUserIDValue, ok := d.GetOk("partner_user_id"); ok {
		config.PartnerUserID = partnerUserIDValue.(int)
	}

	if partnerIDValue, ok := d.GetOk("partner_id"); ok {
		config.PartnerID = partnerIDValue.(int)
	}

	return config, nil
}

// BaseClient -
type BaseClient struct {
	Config          *ClientConfig
	BaseURL         *url.URL
	IdsURL          *url.URL
	UserAgent       string
	Token           string
	IdsClientID     string
	IdsClientSecret string
	IdsScope        string
	HTTPClient      *retryablehttp.Client
}

// NewBaseClient -
func NewBaseClient(config *ClientConfig) *BaseClient {
	return newClient(config, false)
}

// NewLegacyBaseClient -
func NewLegacyBaseClient(config *ClientConfig) *BaseClient {
	return newClient(config, true)
}

func newClient(config *ClientConfig, isLegacy bool) *BaseClient {
	var baseURL *url.URL

	if isLegacy {
		baseURL = config.APIURLLegacy
	} else {
		baseURL = config.APIURL
	}

	// Use PassthroughErrorHandler so that retryablehttp.Client does not obscure API errors
	httpClient := retryablehttp.NewClient()
	httpClient.ErrorHandler = retryablehttp.PassthroughErrorHandler

	return &BaseClient{
		Config:          config,
		BaseURL:         baseURL,
		IdsURL:          config.IdsURL,
		Token:           config.APIToken,
		IdsClientID:     config.IdsClientID,
		IdsClientSecret: config.IdsClientSecret,
		IdsScope:        config.IdsScope,
		HTTPClient:      httpClient,
	}
}

// BuildRequest creates a new Request for a Edgecast API, adding appropriate headers
func (BaseClient *BaseClient) BuildRequest(method, path string, body interface{}, isUsingIdsToken bool) (*retryablehttp.Request, error) {
	relativeURL, err := url.Parse(path)

	if err != nil {
		return nil, fmt.Errorf("BuildRequest: url.Parse: %v", err)
	}

	absoluteURL := BaseClient.BaseURL.ResolveReference(relativeURL)

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
			//helper.LogRequestBody(method, absoluteURL.String(), body)
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
		idsToken, err := BaseClient.GetIdsToken()

		if err != nil {
			return nil, fmt.Errorf("BuildRequest: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+idsToken)
	} else {
		if len(BaseClient.Token) == 0 {
			return nil, errors.New("BuildRequest: API Token is required")
		}

		req.Header.Set("Authorization", "TOK:"+BaseClient.Token)
	}

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// SendRequest sends an HTTP request and, if applicable, sets the response to parsedResponse
func (BaseClient *BaseClient) SendRequest(req *retryablehttp.Request, parsedResponse interface{}) (*http.Response, error) {
	resp, err := BaseClient.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("SendRequest: Do: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	bodyAsString := string(body)
	helper.LogPrettyJson("Response", bodyAsString)
	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		if err != nil {
			return nil, fmt.Errorf("SendRequest: io.ReadAll: %v", err)
		}

		return nil, fmt.Errorf("SendRequest failed: %s", bodyAsString)
	}
	if parsedResponse == nil {
		return nil, nil
	}

	var f interface{}
	if jsonUnmarshalErr := json.Unmarshal(body, &f); err != nil {
		return nil, fmt.Errorf("Malformed Json response:%v", jsonUnmarshalErr)
	}

	if helper.IsInterfaceArray(f) {
		log.Print("isJsonArray")
		if jsonArryErr := json.Unmarshal([]byte(body), parsedResponse); jsonArryErr != nil {
			return nil, fmt.Errorf("Malformed Json Array response:%v", jsonArryErr)
		}
	} else {
		if helper.IsJSONString(bodyAsString) {
			log.Print("Is Not JsonArray")

			err = json.Unmarshal([]byte(bodyAsString), parsedResponse)

			if err != nil {
				return nil, fmt.Errorf("SendRequest: Decode error: %v", err)
			}
		} else {

			// if response is not json string
			switch v := parsedResponse.(type) {
			case LiteralResponse:
				rs, ok := parsedResponse.(LiteralResponse)
				if ok {
					rs.Value = bodyAsString
					parsedResponse = rs
				}
			case float64:
				fmt.Println("float64:", v)
			default:
				fmt.Println("unknown")
			}

		}
	}
	return resp, nil
}

// SendRequest sends an HTTP request and, if applicable, sets the response to parsedResponse
func (BaseClient *BaseClient) SendRequestWithStringResponse(req *retryablehttp.Request) (*string, error) {
	resp, err := BaseClient.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("SendRequest: Do: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	bodyAsString := string(body)

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		if err != nil {
			return nil, fmt.Errorf("SendRequest: io.ReadAll: %v", err)
		}

		return nil, fmt.Errorf("SendRequest failed: %s", bodyAsString)
	}
	// Do not delete this.
	log.Printf("[DEBUG] Raw Response Body:base_client>>SendRequest:%s", body)

	return &bodyAsString, nil
}

// GetIdsToken returns the cached token, refreshing it if it has expired
func (BaseClient *BaseClient) GetIdsToken() (string, error) {

	if len(BaseClient.IdsClientID) == 0 || len(BaseClient.IdsClientSecret) == 0 || len(BaseClient.IdsScope) == 0 {
		return "", errors.New("GetIdsToken: IDS Client ID, Secret, and Scope required")
	}

	if idsToken == nil || idsToken.ExpirationTime.Before(time.Now()) {
		data := url.Values{}
		data.Set("grant_type", "client_credentials")
		data.Add("scope", BaseClient.IdsScope)
		data.Add("client_id", BaseClient.IdsClientID)
		data.Add("client_secret", BaseClient.IdsClientSecret)

		idsTokenEndpoint := fmt.Sprintf("%s/connect/token", BaseClient.IdsURL.String())
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
func FormatURLAddPartnerID(originalURL string, partnerID int) string {
	if partnerID != 0 {
		return originalURL + fmt.Sprintf("&partnerid=%d", partnerID)
	}

	return originalURL
}

// LiteralResponse -
// TODO: Delete this in final refactor to remove base_client.
type LiteralResponse struct {
	Value interface{}
}
