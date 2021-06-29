// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
)

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
	AllowedHTTPMethods         []string        `json:"allowed_http_methods"`
	AllowedRequestContentTypes []string        `json:"allowed_request_content_types"`
	ASNAccessControls          *AccessControls `json:"asn"`
	CookieAccessControls       *AccessControls `json:"cookie"`
	CountryAccessControls      *AccessControls `json:"country"`
	CustomerID                 string          `json:"customer_id"`
	DisallowedExtensions       []string        `json:"disallowed_extensions"`
	DisallowedHeaders          []string        `json:"disallowed_headers"`
	IPAccessControls           *AccessControls `json:"ip"`
	Name                       string          `json:"name"`
	RefererAccessControls      *AccessControls `json:"referer"`
	ResponseHeaderName         string          `json:"response_header_name"`
	URLAccessControls          *AccessControls `json:"url"`
	UserAgentAccessControls    *AccessControls `json:"user_agent"`
}

// AccessControls contains entries that identify traffic. Note: ASN Access Controls must be integers, all other types are strings.
type AccessControls struct {
	Accesslist []interface{} `json:"accesslist"`
	Blacklist  []interface{} `json:"blacklist"`
	Whitelist  []interface{} `json:"whitelist"`
}

type AddAccessRuleResponse struct {
	Id string
}

func (APIClient *WAFAPIClient) AddAccessRule(accessRule AccessRule) (string, error) {
	url := fmt.Sprintf("/v2/mcc/customers/%s/waf/v1.0/acl", accessRule.CustomerID)

	request, err := APIClient.BaseAPIClient.BuildRequest("POST", url, accessRule, false)

	if err != nil {
		return "", fmt.Errorf("AddAccessRule: %v", err)
	}

	parsedResponse := &AddAccessRuleResponse{}

	_, err = APIClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return "", fmt.Errorf("AddAccessRule: %v", err)
	}

	return parsedResponse.Id, nil
}
