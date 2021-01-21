// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package api

import (
	"strconv"
	"strings"
	"time"
)

type RulesEngineApiClient struct {
	BaseApiClient *ApiClient
}

type AddPolicyResponse struct {
	ID          string    `json:"id,omitempty"`
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

// Rule -
type Rule struct {
	ID          string                   `json:"id,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Ordinal     int                      `json:"ordinal,omitempty"`
	CreatedAt   time.Time                `json:"created_at,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at,omitempty"`
	Matches     []map[string]interface{} `json:"matches,omitempty"`
}

// Match -
// Features are treated as a map because their properties are indeterminate
type Match struct {
	Type       string `json:@type`
	Id         int
	Ordinal    int    `json:"ordinal,omitempty"`
	Value      string `json:"value,omitempty"`
	codes      string `json:"codes,omitempty"`
	compare    string `json:"compare,omitempty"`
	encoded    bool   `json:"encoded,omitempty"`
	hostnames  string `json:"hostnames,omitempty"`
	ignoreCase bool   `json:"ignore-case,omitempty"`
	name       string `json:"name,omitempty"`
	relativeTo string `json:"relative-to,omitempty"`
	result     string `json:"result,omitempty"`
	Matches    []map[string]interface{}
	Features   []map[string]interface{}
}

type Feature struct {
	action          string   `json:action,omitempty`
	code            string   `json:code,omitempty`
	destination     string   `json:destination,omitempty`
	eanbled         bool     `json:enabled,omitempty`
	expires         int      `json:expires,omitempty`
	extensions      string   `json:extensions,omitempty`
	format          string   `json:format,omitempty`
	headerName      string   `json:header-name,omitempty`
	headerValue     string   `json:header-value,omitempty`
	instance        string   `json:instance,omitempty`
	kbytesPerSecond int      `json:kbytes-per-second,omitempty`
	mediaTypes      []string `json:mediaTypes,omitempty`
	methods         string   `json:methods,omitempty`
	milliseconds    int      `json:milliseconds,omitempty`
	mode            string   `json:mode,omitempty`
	name            string   `json:name,omitempty`
	names           []string `json:names,omitempty`
	parameters      string   `json:parameters,omitempty`
	prebufSeconds   int      `json:prebuf-seconds,omitempty`
	requests        int      `json:requests,omitempty`
	seconds         int      `json:seconds,omitempty`
	seekEnd         string   `json:seekEnd,omitempty`
	seekStart       string   `json:seekStart,omitempty`
	site            string   `json:site,omitempty`
	source          string   `json:source,omitempty`
	status          string   `json:status,omitempty`
	tags            string   `json:tags,omitempty`
	treatment       string   `json:treatment,omitempty`
	units           string   `json:units,omitempty`
	value           string   `json:value,omitempty`
}

func NewRulesEngineApiClient(baseApiClient *ApiClient) *RulesEngineApiClient {
	apiClient := &RulesEngineApiClient{
		BaseApiClient: baseApiClient,
	}

	return apiClient
}

func (c *RulesEngineApiClient) AddPolicy(policy string, customerid string, portaltypeid string, customeruserid string) (*AddPolicyResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("POST", "rules-engine/v1.1/policies", policy, true)
	// convert hex string to int
	output, err := strconv.ParseInt(hexaNumberToInteger(customerid), 16, 64)
	InfoLogger.Printf("customerId: %d\n", output)
	request.Header.Set("Portals_CustomerId", strconv.FormatInt(output, 10))
	request.Header.Set("Portals_UserId", customeruserid)
	request.Header.Set("Portals_PortalTypeId", portaltypeid)
	InfoLogger.Printf("policy from terraform.tfvars: %s\n", policy)
	parsedResponse := &AddPolicyResponse{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	InfoLogger.Printf("RE policy response: %+v\n", parsedResponse)
	return parsedResponse, err
}

/*
func (c *RulesEngineApiClient) GetHttpLargeOrigin(id int) (*Origin, error) {
	request, err := c.BaseApiClient.BuildRequest("GET", fmt.Sprintf("rules-engine/v1.1/policies/%d", id), nil, true)

	parsedResponse := &Origin{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	return parsedResponse, err
}


func (c *RulesEngineApiClient) DeleteOrigin(id int) error {
	request, err := c.BaseApiClient.BuildRequest("DELETE", fmt.Sprintf("mcc/customers/%s/origins/%d", c.AccountNumber, id), nil)

	_, err = c.BaseApiClient.SendRequest(request, nil)

	return err
}*/

func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}
