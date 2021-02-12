// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
)

type CnameApiClient struct {
	BaseApiClient *ApiClient
	AccountNumber string
}

type Cname struct {
	Id                  int
	Name                string
	DirPath             string
	EnableCustomReports int
	OriginId            int
	OriginString        string
}

type AddCnameRequest struct {
	Name                string
	DirPath             string
	EnableCustomReports int
	MediaTypeId         int
	OriginId            int
	OriginType          int
}

type AddCnameResponse struct {
	CnameId int
}

type UpdateCnameRequest struct {
	Name                string
	DirPath             string
	EnableCustomReports int
	MediaTypeId         int
	OriginId            int
	OriginType          int
}

type UpdateCnameResponse struct {
	CnameId int
}

func NewCnameApiClient(baseApiClient *ApiClient, accountNumber string) *CnameApiClient {
	apiClient := &CnameApiClient{
		BaseApiClient: baseApiClient,
		AccountNumber: accountNumber,
	}

	return apiClient
}

func (c *CnameApiClient) AddCname(cname *AddCnameRequest) (*AddCnameResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("POST", fmt.Sprintf("v2/mcc/customers/%s/cnames", c.AccountNumber), cname, false)
	parsedResponse := &AddCnameResponse{}
	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("AddCname: %v", err)
	}

	return parsedResponse, nil
}

func (c *CnameApiClient) UpdateCname(cname *UpdateCnameRequest, cnameId int) (*UpdateCnameResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("PUT", fmt.Sprintf("mcc/customers/%s/cnames/%d", c.AccountNumber, cnameId), cname, false)
	parsedResponse := &UpdateCnameResponse{}
	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("UpdateCname: %v", err)
	}

	return parsedResponse, nil
}

func (c *CnameApiClient) GetCname(id int) (*Cname, error) {
	request, err := c.BaseApiClient.BuildRequest("GET", fmt.Sprintf("v2/mcc/customers/%s/cnames/%d", c.AccountNumber, id), nil, false)
	parsedResponse := &Cname{}
	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetCname: %v", err)
	}

	return parsedResponse, nil
}

func (c *CnameApiClient) DeleteCname(id int) error {
	request, err := c.BaseApiClient.BuildRequest("DELETE", fmt.Sprintf("v2/mcc/customers/%s/cnames/%d", c.AccountNumber, id), nil, false)
	_, err = c.BaseApiClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteOrigin: %v", err)
	}

	return nil

}
