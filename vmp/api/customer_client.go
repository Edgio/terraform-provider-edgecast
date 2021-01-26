// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
)

// CustomerAPIClient interacts with the Verizon Media API
type CustomerAPIClient struct {
	BaseAPIClient *ApiClient
	PartnerUserID *int
	PartnerID     *int
}

// CustomerCreateUpdate represents a request to add a new customer
type CustomerCreateUpdate struct {
	AccountID                 string // Ignored when Updating
	Address1                  string
	Address2                  string
	BandwidthUsageLimit       int64
	BillingAccountTag         string
	BillingAddress1           string
	BillingAddress2           string
	BillingCity               string
	BillingContactEmail       string
	BillingContactFax         string
	BillingContactFirstName   string
	BillingContactLastName    string
	BillingContactMobile      string
	BillingContactPhone       string
	BillingContactTitle       string
	BillingCountry            string
	BillingRateInfo           string
	BillingState              string
	BillingZIP                string
	City                      string
	CompanyName               string
	ContactEmail              string
	ContactFax                string
	ContactFirstName          string
	ContactLastName           string
	ContactMobile             string
	ContactPhone              string
	ContactTitle              string
	Country                   string
	DataTransferredUsageLimit int64
	Notes                     string
	ServiceLevelCode          string
	State                     string
	Status                    int // Ignored when Updating
	Website                   string
	ZIP                       string
}

// NewCustomerAPIClient -
func NewCustomerAPIClient(baseAPIClient *ApiClient, partnerUserID *int, partnerID *int) *CustomerAPIClient {
	return &CustomerAPIClient{
		BaseAPIClient: baseAPIClient,
		PartnerUserID: partnerUserID,
		PartnerID:     partnerID,
	}
}

// AddCustomer -
func (apiClient *CustomerAPIClient) AddCustomer(body *CustomerCreateUpdate) (string, error) {
	relURL := "v2/pcc/customers"
	if apiClient.PartnerUserID != nil {
		relURL = relURL + fmt.Sprintf("?partneruserid=%d", *apiClient.PartnerUserID)
	}

	request, err := apiClient.BaseAPIClient.BuildRequest("POST", relURL, body, false)
	InfoLogger.Printf("AddCustomer [POST] Url: %s\n", request.URL)

	parsedResponse := &struct {
		AccountNumber string
	}{}

	_, err = apiClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return "", err
	}

	return parsedResponse.AccountNumber, err
}

// GetCustomerResponse -
type GetCustomerResponse struct {
	ID                        int32
	CustomID                  string
	HexID                     string
	Address1                  string
	Address2                  string
	BandwidthUsageLimit       int64
	UsageLimitUpdateDate      string
	BillingAccountTag         string
	BillingAddress1           string
	BillingAddress2           string
	BillingCity               string
	BillingContactEmail       string
	BillingContactFax         string
	BillingContactFirstName   string
	BillingContactLastName    string
	BillingContactMobile      string
	BillingContactPhone       string
	BillingContactTitle       string
	BillingCountry            string
	BillingRateInfo           string
	BillingState              string
	BillingZIP                string
	City                      string
	CompanyName               string
	ContactEmail              string
	ContactFax                string
	ContactFirstName          string
	ContactLastName           string
	ContactMobile             string
	ContactPhone              string
	ContactTitle              string
	Country                   string
	DataTransferredUsageLimit int64
	Notes                     string
	ServiceLevelCode          string
	State                     string
	Status                    int8
	Website                   string
	ZIP                       string
	PartnerID                 int
	PartnerName               string
	WholesaleID               int
	WholesaleName             string
}

// GetCustomer retrieves a Customer's info using the Hex Account Number
func (apiClient *CustomerAPIClient) GetCustomer(accountNumber string) (*GetCustomerResponse, error) {
	relURL := fmt.Sprintf("v2/pcc/customers/%s", accountNumber)
	request, err := apiClient.BaseAPIClient.BuildRequest("GET", relURL, nil, false)
	InfoLogger.Printf("AddHttpLargeOrigin [POST] Url: %s\n", request.URL)

	if err != nil {
		return nil, err
	}

	parsedResponse := &GetCustomerResponse{}

	_, err = apiClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, err
	}

	return parsedResponse, nil
}

// AccessModule represents a module that a customer has access to
type AccessModule struct {
	ID       int
	Name     string
	ParentID *int
}

// GetCustomerAccessModules retrieves a Customer's Access Module info using the Hex Account Number
func (apiClient *CustomerAPIClient) GetCustomerAccessModules(accountNumber string) (*[]AccessModule, error) {
	relURL := fmt.Sprintf("v2/pcc/customers/%s/accessmodules", accountNumber)
	request, err := apiClient.BaseAPIClient.BuildRequest("GET", relURL, nil, false)
	InfoLogger.Printf("GetCustomerAccessModules [GET] Url: %s\n", request.URL)

	if err != nil {
		return nil, err
	}

	var accessModules []AccessModule

	_, err = apiClient.BaseAPIClient.SendRequest(request, &accessModules)

	if err != nil {
		return nil, err
	}

	return &accessModules, nil
}

// Service -
type Service struct {
	ID       int
	Name     string
	ParentID int
	Status   int8
}

// GetAvailableCustomerServices gets all service information available for a partner to administor to thier customers
func (apiClient *CustomerAPIClient) GetAvailableCustomerServices() (*[]Service, error) {
	request, err := apiClient.BaseAPIClient.BuildRequest("GET", "v2/pcc/customers/services", nil, false)
	InfoLogger.Printf("GetAvailableCustomerServices [GET] Url: %s\n", request.URL)

	if err != nil {
		return nil, err
	}

	var services []Service

	_, err = apiClient.BaseAPIClient.SendRequest(request, &services)

	if err != nil {
		return nil, err
	}

	return &services, nil
}

// GetCustomerServices gets the list of services available to a customer and whether each is active for the customer
func (apiClient *CustomerAPIClient) GetCustomerServices(accountNumber string) ([]Service, error) {
	relURL := fmt.Sprintf("v2/pcc/customers/%s/services", accountNumber)
	request, err := apiClient.BaseAPIClient.BuildRequest("GET", relURL, nil, false)
	InfoLogger.Printf("GetCustomerServices [GET] Url: %s\n", request.URL)

	if err != nil {
		return nil, err
	}

	var services []Service

	_, err = apiClient.BaseAPIClient.SendRequest(request, &services)

	if err != nil {
		return nil, err
	}

	return services, nil
}

// UpdateCustomerServices -
func (apiClient *CustomerAPIClient) UpdateCustomerServices(accountNumber string, serviceIDs []int, status int8) error {
	relURL := fmt.Sprintf("v2/pcc/customers/%s/services", accountNumber)

	body := &struct {
		Status      int8
		ServiceCode []int
	}{
		Status:      status,
		ServiceCode: serviceIDs,
	}

	request, err := apiClient.BaseAPIClient.BuildRequest("PUT", relURL, body, false)
	InfoLogger.Printf("UpdateCustomerServicers [PUT] Url: %s\n", request.URL)

	if err != nil {
		return err
	}

	resp, err := apiClient.BaseAPIClient.SendRequest(request, nil)

	if err == nil && resp.StatusCode != 200 {
		return fmt.Errorf("failed to set customer services, please contact an administrator")
	}

	if err != nil {
		return err
	}

	return nil
}

// GetCustomerDeliveryRegion gets the current active delivery region set for the customer
func (apiClient *CustomerAPIClient) GetCustomerDeliveryRegion(accountNumber string) (int, error) {
	relURL := fmt.Sprintf("v2/pcc/customers/%s/deliveryregions", accountNumber)

	request, err := apiClient.BaseAPIClient.BuildRequest("GET", relURL, nil, false)
	InfoLogger.Printf("GetCustomerDeliveryRegion [GET] Url: %s\n", request.URL)

	if err != nil {
		return 0, err
	}

	parsedResponse := &struct {
		AccountNumber    string
		CustomID         string
		DeliveryRegionID int
	}{}

	_, err = apiClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return 0, err
	}

	return parsedResponse.DeliveryRegionID, nil
}

// UpdateCustomerDomainURL -
func (apiClient *CustomerAPIClient) UpdateCustomerDomainURL(accountNumber string, domainType int, url string) error {
	// TODO: support custom ids for accounts
	baseURL := fmt.Sprintf("v2/pcc/customers/domains/%d/url?idtype=an&id=%s", domainType, accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, apiClient.PartnerID)

	body := &struct {
		URL string `json:"Url"`
	}{
		URL: url,
	}

	request, err := apiClient.BaseAPIClient.BuildRequest("PUT", relURL, body, false)
	InfoLogger.Printf("UpdateCustomerDomainURL [PUT] Url: %s\n", request.URL)

	if err != nil {
		return err
	}

	_, err = apiClient.BaseAPIClient.SendRequest(request, nil)

	return err
}

// UpdateCustomerDeliveryRegion -
func (apiClient *CustomerAPIClient) UpdateCustomerDeliveryRegion(accountNumber string, deliveryRegionID int) error {
	// TODO: support custom ids for accounts
	baseURL := fmt.Sprintf("v2/pcc/customers/deliveryregions?idtype=an&id=%s", accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, apiClient.PartnerID)

	body := &struct {
		ID int `json:"Id"`
	}{
		ID: deliveryRegionID,
	}

	request, err := apiClient.BaseAPIClient.BuildRequest("PUT", relURL, body, false)
	InfoLogger.Printf("UpdateCustomerDeliveryRegion [PUT] Url: %s\n", request.URL)

	if err != nil {
		return err
	}

	_, err = apiClient.BaseAPIClient.SendRequest(request, nil)

	return err
}

// DeleteCustomer -
func (apiClient *CustomerAPIClient) DeleteCustomer(accountNumber string) error {
	// TODO: support custom ids for accounts
	baseURL := fmt.Sprintf("v2/pcc/customers?idtype=an&id=%s", accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, apiClient.PartnerID)

	request, err := apiClient.BaseAPIClient.BuildRequest("DELETE", relURL, nil, false)
	InfoLogger.Printf("DeleteCustomer [DELETE] Url: %s\n", request.URL)

	if err != nil {
		return err
	}

	_, err = apiClient.BaseAPIClient.SendRequest(request, nil)

	return err
}

// UpdateCustomer -
func (apiClient *CustomerAPIClient) UpdateCustomer(accountNumber string, body *CustomerCreateUpdate) error {
	// TODO: support custom ids for accounts
	baseURL := fmt.Sprintf("v2/pcc/customers?idtype=an&id=%s", accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, apiClient.PartnerID)

	request, err := apiClient.BaseAPIClient.BuildRequest("PUT", relURL, body, false)
	InfoLogger.Printf("UpdateCustomer [PUT] Url: %s\n", request.URL)

	if err != nil {
		return err
	}

	_, err = apiClient.BaseAPIClient.SendRequest(request, nil)

	return err
}

// UpdateCustomerAccessModule -
func (apiClient *CustomerAPIClient) UpdateCustomerAccessModule(accountNumber string, accessModuleID int) error {
	// TODO: support custom ids for accounts
	baseURL := fmt.Sprintf("v2/pcc/customers/accessmodules/%d/status?idtype=an&id=%s", accessModuleID, accountNumber)
	relURL := FormatURLAddPartnerID(baseURL, apiClient.PartnerID)
	body := &struct{ Status int8 }{Status: 1}

	request, err := apiClient.BaseAPIClient.BuildRequest("PUT", relURL, body, false)
	InfoLogger.Printf("UpdateCustomerAccessModule [PUT] Url: %s\n", request.URL)

	if err != nil {
		return err
	}

	_, err = apiClient.BaseAPIClient.SendRequest(request, nil)

	return err
}
