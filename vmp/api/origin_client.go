// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
)

//OriginAPIClient -
type OriginAPIClient struct {
	Config        *ClientConfig
	BaseAPIClient *BaseClient
	AccountNumber string
}

//AddOriginRequest -
type AddOriginRequest struct {
	DirectoryName     string
	HostHeader        string
	HTTPHostnames     []AddOriginRequestHostname
	HTTPSHostnames    []AddOriginRequestHostname
	HTTPLoadBalancing string
}

//AddOriginRequestHostname -
type AddOriginRequestHostname struct {
	Name string
}

//AddOriginResponse -
type AddOriginResponse struct {
	CustomerOriginID int
}

//UpdateOriginRequest -
type UpdateOriginRequest struct {
	DirectoryName     string
	HostHeader        string
	HTTPHostnames     []UpdateOriginRequestHostname
	HTTPSHostnames    []UpdateOriginRequestHostname
	HTTPLoadBalancing string
}

//UpdateOriginRequestHostname -
type UpdateOriginRequestHostname struct {
	Name string
}

//UpdateOriginResponse -
type UpdateOriginResponse struct {
	CustomerOriginID int
}

//Origin -
type Origin struct {
	ID                int
	DirectoryName     string
	HostHeader        string
	HTTPHostnames     []OriginHostname
	HTTPLoadBalancing string
}

//OriginHostname -
type OriginHostname struct {
	Name string
}

//NewOriginAPIClient -
func NewOriginAPIClient(config *ClientConfig) *OriginAPIClient {
	apiClient := &OriginAPIClient{
		Config:        config,
		BaseAPIClient: config.BaseClient,
		AccountNumber: config.AccountNumber,
	}

	return apiClient
}

//AddOrigin -
func (c *OriginAPIClient) AddOrigin(origin *AddOriginRequest, mediaType string) (*AddOriginResponse, error) {
	request, err := c.BaseAPIClient.BuildRequest("POST", fmt.Sprintf("v2/mcc/customers/%s/origins/%s", c.AccountNumber, mediaType), origin, false)
	parsedResponse := &AddOriginResponse{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("AddOrigin: %v", err)
	}

	return parsedResponse, nil
}

//UpdateOrigin -
func (c *OriginAPIClient) UpdateOrigin(origin *UpdateOriginRequest, originID int, mediaType string) (*UpdateOriginResponse, error) {
	request, err := c.BaseAPIClient.BuildRequest("PUT", fmt.Sprintf("v2/mcc/customers/%s/origins/%s/%d", c.AccountNumber, mediaType, originID), origin, false)
	parsedResponse := &UpdateOriginResponse{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("UpdateOrigin: %v", err)
	}

	return parsedResponse, nil
}

//GetOrigin -
func (c *OriginAPIClient) GetOrigin(id int, mediaType string) (*Origin, error) {
	request, err := c.BaseAPIClient.BuildRequest("GET", fmt.Sprintf("v2/mcc/customers/%s/origins/%s/%d", c.AccountNumber, mediaType, id), nil, false)
	parsedResponse := &Origin{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetOrigin: %v", err)
	}

	return parsedResponse, nil
}

//DeleteOrigin -
func (c *OriginAPIClient) DeleteOrigin(id int) error {
	request, err := c.BaseAPIClient.BuildRequest("DELETE", fmt.Sprintf("v2/mcc/customers/%s/origins/%d", c.AccountNumber, id), nil, false)
	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteOrigin: %v", err)
	}

	return nil
}
