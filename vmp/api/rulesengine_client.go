// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
	"time"
)

// RulesEngineAPIClient interacts with the Verizon Media Rules Engine API
type RulesEngineAPIClient struct {
	BaseAPIClient *ApiClient
}

// NewRulesEngineAPIClient -
func NewRulesEngineAPIClient(baseAPIClient *ApiClient) *RulesEngineAPIClient {
	return &RulesEngineAPIClient{
		BaseAPIClient: baseAPIClient,
	}
}

// GetPolicyResponse -
type GetPolicyResponse struct {
	ID          int
	Name        string
	Description string
	PolicyType  string `json:"policy_type"`
	State       string
	Platform    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rule -
type Rule struct {
	ID          int
	Name        string
	Description string
	Ordinal     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Matches     []Match
}

// Match -
// Features are treated as a map because their properties are indeterminate
type Match struct {
	Type     string
	Ordinal  int
	Value    string
	Matches  []Match
	Features map[string]string
}

// GetPolicy -
func (apiClient *RulesEngineAPIClient) GetPolicy(policyID int) (*GetPolicyResponse, error) {
	relURL := fmt.Sprintf("rules-engine/v.1/policies/%d", policyID)
	request, err := apiClient.BaseAPIClient.BuildRequest("GET", relURL, nil)

	if err != nil {
		return nil, err
	}

	parsedResponse := &GetPolicyResponse{}

	_, err = apiClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, err
	}

	return parsedResponse, nil
}
