// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const rulesEngineRelUrlFormat = "rules-engine/v1.1/%s"

type RulesEngineApiClient struct {
	BaseApiClient *ApiClient
}

type GetPolicyResponse struct {
	Id          string    `json:"id,omitempty"`
	Type        string    `json:"@type,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	PolicyType  string    `json:"policy_type,omitempty"`
	State       string    `json:"state,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Rules       []Rule    `json:"rules,omitempty"`
}

type AddPolicyResponse struct {
	Id          string    `json:"id,omitempty"`
	Type        string    `json:"@type,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	PolicyType  string    `json:"policy_type,omitempty"`
	State       string    `json:"state,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Rules       []Rule    `json:"rules,omitempty"`
}

type Rule struct {
	Id          string                   `json:"id,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Ordinal     int                      `json:"ordinal,omitempty"`
	CreatedAt   time.Time                `json:"created_at,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at,omitempty"`
	Matches     []map[string]interface{} `json:"matches,omitempty"`
}

type Match struct {
	Id         int    `json:"id`
	Type       string `json:"@type"`
	Ordinal    int    `json:"ordinal,omitempty"`
	Value      string `json:"value,omitempty"`
	Codes      string `json:"codes,omitempty"`
	Compare    string `json:"compare,omitempty"`
	Encoded    bool   `json:"encoded,omitempty"`
	Hostnames  string `json:"hostnames,omitempty"`
	IgnoreCase bool   `json:"ignore-case,omitempty"`
	Name       string `json:"name,omitempty"`
	RelativeTo string `json:"relative-to,omitempty"`
	Result     string `json:"result,omitempty"`
	Matches    []map[string]interface{}
	Features   []map[string]interface{}
}

type Feature struct {
	Action          string   `json:"action,omitempty"`
	Code            string   `json:"code,omitempty"`
	Destination     string   `json:"destination,omitempty"`
	Enabled         bool     `json:"enabled,omitempty"`
	Expires         int      `json:"expires,omitempty"`
	Extensions      string   `json:"extensions,omitempty"`
	Format          string   `json:"format,omitempty"`
	HeaderName      string   `json:"header-name,omitempty"`
	HeaderValue     string   `json:"header-value,omitempty"`
	Instance        string   `json:"instance,omitempty"`
	KbytesPerSecond int      `json:"kbytes-per-second,omitempty"`
	MediaTypes      []string `json:"mediaTypes,omitempty"`
	Methods         string   `json:"methods,omitempty"`
	Milliseconds    int      `json:"milliseconds,omitempty"`
	Mode            string   `json:"mode,omitempty"`
	Name            string   `json:"name,omitempty"`
	Names           []string `json:"names,omitempty"`
	Parameters      string   `json:"parameters,omitempty"`
	PrebufSeconds   int      `json:"prebuf-seconds,omitempty"`
	Requests        int      `json:"requests,omitempty"`
	Seconds         int      `json:"seconds,omitempty"`
	SeekEnd         string   `json:"seekEnd,omitempty"`
	SeekStart       string   `json:"seekStart,omitempty"`
	Site            string   `json:"site,omitempty"`
	Source          string   `json:"source,omitempty"`
	Status          string   `json:"status,omitempty"`
	Tags            string   `json:"tags,omitempty"`
	Treatment       string   `json:"treatment,omitempty"`
	Units           string   `json:"units,omitempty"`
	Value           string   `json:"value,omitempty"`
}

func NewRulesEngineApiClient(baseApiClient *ApiClient) *RulesEngineApiClient {
	apiClient := &RulesEngineApiClient{
		BaseApiClient: baseApiClient,
	}

	return apiClient
}

// GetPolicy -
func (apiClient *RulesEngineApiClient) GetPolicy(customerId int, customerUserId string, portalTypeId string, policyId int) (map[string]interface{}, error) {
	relURL := formatRulesEngineRelURL("policies/%d", policyId)
	request, err := apiClient.BaseApiClient.BuildRequest("GET", relURL, nil, true)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Portals_CustomerId", strconv.Itoa(customerId))
	request.Header.Set("Portals_UserId", customerUserId)
	request.Header.Set("Portals_PortalTypeId", portalTypeId)

	InfoLogger.Printf("GetPolicy [GET] Url: %s\n", request.URL)
	//parsedResponse := &GetPolicyResponse{}
	parsedResponse := make(map[string]interface{})

	_, err = apiClient.BaseApiClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, err
	}

	InfoLogger.Printf("RE policy response: %+v\n", parsedResponse)
	return parsedResponse, nil
}

func (c *RulesEngineApiClient) AddPolicy(policy string, accountNumber string, portalTypeId string, customerUserId string) (*AddPolicyResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("POST", "rules-engine/v1.1/policies", policy, true)

	// account number hex string -> customer ID
	customerId, err := strconv.ParseInt(accountNumber, 16, 64)

	if err != nil {
		return nil, err
	}

	InfoLogger.Printf("customerId: %d\n", customerId)

	request.Header.Set("Portals_CustomerId", strconv.FormatInt(customerId, 10))
	request.Header.Set("Portals_UserId", customerUserId)
	request.Header.Set("Portals_PortalTypeId", portalTypeId)
	InfoLogger.Printf("policy from terraform.tfvars: %s\n", policy)
	parsedResponse := &AddPolicyResponse{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	InfoLogger.Printf("RE policy response: %+v\n", parsedResponse)
	return parsedResponse, err
}

func removeHexPrefix(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

func formatRulesEngineRelURL(subFormat string, params ...interface{}) string {
	subPath := fmt.Sprintf(subFormat, params...)
	return fmt.Sprintf(rulesEngineRelUrlFormat, subPath)
}
