// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import "fmt"

// WAFAPIClient interacts with the Verizon Media API
type WAFAPIClient struct {
	Config        *ClientConfig
	BaseAPIClient *BaseClient
}

// WAFAPIClient -
func NewWAFAPIClient(config *ClientConfig) *WAFAPIClient {
	return &WAFAPIClient{
		Config:        config,
		BaseAPIClient: config.BaseClientLegacy,
	}
}

type AccessRule struct {
	AllowedHTTPMethods         []string
	AllowedRequestContentTypes []string
	ASNAccessControls          *AcessControls `json:"asn"`
	CookieAccessControls       *AcessControls `json:"cookie"`
	CountryAccessControls      *AcessControls `json:"country"`
	CustomerID                 string
	DisallowedExtensions       []string
	DisallowedHeaders          []string
	IPAccessControls           *AcessControls `json:"ip"`
	Name                       string
	RefererAccessControls      *AcessControls `json:"referer"`
	ResponseHeaderName         string
	URLAccessControls          *AcessControls `json:"url"`
}

type AcessControls struct {
	AccessList []interface{}
	Blacklist  []interface{}
	Whitelist  []interface{}
}

type WAFError struct {
	Code    int
	Message string
}

type AddAccessRuleResponse struct {
	Id      int
	Status  string
	Success string
	Errors  []WAFError
}

func (APIClient *WAFAPIClient) AddAccessRule(accessRule *AccessRule) (*AddAccessRuleResponse, error) {
	url := fmt.Sprintf("/v2/mcc/customers/%s/waf/v1.0/acl", accessRule.CustomerID)

	request, err := APIClient.BaseAPIClient.BuildRequest("POST", url, accessRule, false)

	if err != nil {
		return nil, fmt.Errorf("AddAccessRule: %v", err)
	}

	parsedResponse := &AddAccessRuleResponse{}

	_, err = APIClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("AddAccessRule: %v", err)
	}

	return parsedResponse, nil
}
