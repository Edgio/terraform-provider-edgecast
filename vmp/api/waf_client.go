// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
	"strconv"
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
	AllowedHTTPMethods         []string
	AllowedRequestContentTypes []string
	ASNAccessControls          AccessControls `json:"asn"`
	CookieAccessControls       AccessControls `json:"cookie"`
	CountryAccessControls      AccessControls `json:"country"`
	CustomerID                 string
	DisallowedExtensions       []string
	DisallowedHeaders          []string
	IPAccessControls           AccessControls `json:"ip"`
	Name                       string
	RefererAccessControls      AccessControls `json:"referer"`
	ResponseHeaderName         string
	URLAccessControls          AccessControls `json:"url"`
}

type AccessControls struct {
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
	Success bool
	Errors  []WAFError
}

func (APIClient *WAFAPIClient) AddAccessRule(accessRule AccessRule) (int, error) {
	url := fmt.Sprintf("/v2/mcc/customers/%s/waf/v1.0/acl", accessRule.CustomerID)

	request, err := APIClient.BaseAPIClient.BuildRequest("POST", url, accessRule, false)

	if err != nil {
		return 0, fmt.Errorf("AddAccessRule: %v", err)
	}

	parsedResponse := &AddAccessRuleResponse{}

	_, err = APIClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return 0, fmt.Errorf("AddAccessRule: %v", err)
	}

	if !parsedResponse.Success || len(parsedResponse.Errors) > 0 {
		return 0, fmt.Errorf("AddAccessRule: Errors: %v", flattenWAFErrors(parsedResponse.Errors))
	}

	return parsedResponse.Id, nil
}

func flattenWAFErrors(errors []WAFError) string {
	error := ""

	for i, v := range errors {
		if i > 0 {
			error += ","
		}

		error += strconv.Itoa(v.Code) + ":" + v.Message
	}

	return error
}
