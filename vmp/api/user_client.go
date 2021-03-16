// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
)

// UserAPIClient interacts with the Verizon Media API
type UserAPIClient struct {
	Config        *ClientConfig
	BaseAPIClient *BaseClient
	PartnerID     int
}

// NewUserAPIClient -
func NewUserAPIClient(config *ClientConfig) *UserAPIClient {
	return &UserAPIClient{
		Config:        config,
		BaseAPIClient: config.BaseClientLegacy,
		PartnerID:     config.PartnerID,
	}
}

// CustomerUser -
type CustomerUser struct {
	Address1      string
	Address2      string
	City          string
	Country       string
	CustomID      string `json:"CustomId"` // Read-only
	Email         string
	Fax           string
	FirstName     string
	IsAdmin       int8
	LastName      string
	Mobile        string
	Password      *string // nullable
	Phone         string
	State         string
	TimeZoneID    *int `json:"TimeZoneId"` // nullable
	Title         string
	ZIP           string `json:"Zip"`
	LastLoginDate string // Read-only
}

// GetCustomerUser -
func (APIClient *UserAPIClient) GetCustomerUser(accountNumber string, customerUserID int) (*CustomerUser, error) {
	// TODO: support custom id types, not just Hex ID ANs
	baseURL := fmt.Sprintf("v2/pcc/customers/users/%d?idtype=an&id=%s", customerUserID, accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, APIClient.PartnerID)

	request, err := APIClient.BaseAPIClient.BuildRequest("GET", relURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetCustomerUser: %v", err)
	}

	parsedResponse := &CustomerUser{}

	_, err = APIClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetCustomerUser: %v", err)
	}

	return parsedResponse, nil
}

// AddCustomerUser -
func (APIClient *UserAPIClient) AddCustomerUser(accountNumber string, body *CustomerUser) (int, error) {
	// TODO: support custom id types, not just Hex ID ANs
	baseURL := fmt.Sprintf("v2/pcc/customers/users?idtype=an&id=%s", accountNumber)
	if APIClient.PartnerID == 0 {
		return 0, fmt.Errorf("partner_id was not provided.")
	}
	relURL := FormatURLAddPartnerID(baseURL, APIClient.PartnerID)

	request, err := APIClient.BaseAPIClient.BuildRequest("POST", relURL, body, false)

	if err != nil {
		return 0, fmt.Errorf("AddCustomerUser: %v", err)
	}

	parsedResponse := &struct {
		CustomerUserID int `json:"CustomerUserId"`
	}{}

	_, err = APIClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return 0, fmt.Errorf("AddCustomerUser: %v", err)
	}

	return parsedResponse.CustomerUserID, nil
}

// UpdateCustomerUser -
func (APIClient *UserAPIClient) UpdateCustomerUser(accountNumber string, customerUserID int, body *CustomerUser) error {
	// TODO: support custom ids for accounts
	baseURL := fmt.Sprintf("v2/pcc/customers/users/%d?idtype=an&id=%s", customerUserID, accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, APIClient.PartnerID)

	request, err := APIClient.BaseAPIClient.BuildRequest("PUT", relURL, body, false)

	if err != nil {
		return fmt.Errorf("UpdateCustomerUser: %v", err)
	}

	_, err = APIClient.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("UpdateCustomerUser: %v", err)
	}

	return nil
}

// DeleteCustomerUser -
func (APIClient *UserAPIClient) DeleteCustomerUser(accountNumber string, customerUserID int) error {
	// TODO: support custom ids for accounts
	baseURL := fmt.Sprintf("v2/pcc/customers/users/%d?idtype=an&id=%s", customerUserID, accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, APIClient.PartnerID)

	request, err := APIClient.BaseAPIClient.BuildRequest("DELETE", relURL, nil, false)

	if err != nil {
		return fmt.Errorf("DeleteCustomerUser: %v", err)
	}

	_, err = APIClient.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteCustomerUser: %v", err)
	}

	return nil
}
