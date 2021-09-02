// Copyright Edgecast, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
)

//CnameAPIClient -
type CnameAPIClient struct {
	Config        *ClientConfig
	BaseAPIClient *BaseClient
	AccountNumber string
}

//Cname -
type Cname struct {
	ID                  int
	Name                string
	DirPath             string
	EnableCustomReports int
	OriginID            int `json:"OriginId,omitempty"`
	OriginString        string
}

//AddCnameRequest -
type AddCnameRequest struct {
	Name                string
	DirPath             string
	EnableCustomReports int
	MediaTypeID         int `json:"MediaTypeId,omitempty"`
	OriginID            int `json:"OriginId,omitempty"`
	OriginType          int
}

//AddCnameResponse -
type AddCnameResponse struct {
	CnameID int
}

// UpdateCnameRequest -
type UpdateCnameRequest struct {
	Name                string
	DirPath             string
	EnableCustomReports int
	MediaTypeID         int `json:"MediaTypeId,omitempty"`
	OriginID            int `json:"OriginId,omitempty"`
	OriginType          int
}

//UpdateCnameResponse -
type UpdateCnameResponse struct {
	CnameID int `json:"CnameId,omitempty"`
}

//NewCnameAPIClient -
func NewCnameAPIClient(config *ClientConfig) *CnameAPIClient {
	apiClient := &CnameAPIClient{
		Config:        config,
		BaseAPIClient: config.BaseClientLegacy,
		AccountNumber: config.AccountNumber,
	}

	return apiClient
}

//AddCname -
func (c *CnameAPIClient) AddCname(cname *AddCnameRequest) (*AddCnameResponse, error) {
	request, err := c.BaseAPIClient.BuildRequest("POST", fmt.Sprintf("v2/mcc/customers/%s/cnames", c.AccountNumber), cname, false)
	if err != nil {
		return nil, fmt.Errorf("AddCname: %v", err)
	}

	parsedResponse := &AddCnameResponse{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("AddCname: %v", err)
	}

	return parsedResponse, nil
}

//UpdateCname =
func (c *CnameAPIClient) UpdateCname(cname *UpdateCnameRequest, cnameID int) (*UpdateCnameResponse, error) {
	request, err := c.BaseAPIClient.BuildRequest("PUT", fmt.Sprintf("v2/mcc/customers/%s/cnames/%d", c.AccountNumber, cnameID), cname, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateCname: %v", err)
	}

	parsedResponse := &UpdateCnameResponse{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("UpdateCname: %v", err)
	}

	return parsedResponse, nil
}

//GetCname -
func (c *CnameAPIClient) GetCname(id int) (*Cname, error) {
	request, err := c.BaseAPIClient.BuildRequest("GET", fmt.Sprintf("v2/mcc/customers/%s/cnames/%d", c.AccountNumber, id), nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetCname: %v", err)
	}

	parsedResponse := &Cname{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetCname: %v", err)
	}

	return parsedResponse, nil
}

// DeleteCname -
func (c *CnameAPIClient) DeleteCname(id int) error {
	request, err := c.BaseAPIClient.BuildRequest("DELETE", fmt.Sprintf("v2/mcc/customers/%s/cnames/%d", c.AccountNumber, id), nil, false)
	if err != nil {
		return fmt.Errorf("DeleteCname: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteCname: %v", err)
	}

	return nil

}
