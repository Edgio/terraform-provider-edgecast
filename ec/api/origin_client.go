// Copyright 2021 Edgecast Inc. Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package api

import (
	"fmt"
	"log"
)

//OriginAPIClient -
type OriginAPIClient struct {
	Config        *ClientConfig
	BaseAPIClient *BaseClient
	AccountNumber string
}

//AddOriginRequest - Don't change variable names, otherwise API call to api server will get an error message as response.
type AddOriginRequest struct {
	DirectoryName      string
	HostHeader         string
	HTTPLoadBalancing  string                     `json:"HttpLoadBalancing,omitempty"`
	HTTPSLoadBalancing string                     `json:"HttpsLoadBalancing,omitempty"`
	HTTPHostnames      []AddOriginRequestHostname `json:"HttpHostnames,omitempty"`
	HTTPSHostnames     []AddOriginRequestHostname `json:"HttpsHostnames,omitempty"`
}

//AddOriginRequestHostname -
type AddOriginRequestHostname struct {
	Name string
}

//AddOriginResponse -
type AddOriginResponse struct {
	CustomerOriginID int
}

//UpdateOriginRequest - Don't change variable names, otherwise API call to api server will get an error message as response.
type UpdateOriginRequest struct {
	DirectoryName      string
	HostHeader         string
	HTTPHostnames      []UpdateOriginRequestHostname `json:"HttpHostnames,omitempty"`
	HTTPSHostnames     []UpdateOriginRequestHostname `json:"HttpsHostnames,omitempty"`
	HTTPLoadBalancing  string                        `json:"HttpLoadBalancing,omitempty"`
	HTTPSLoadBalancing string                        `json:"HttpsLoadBalancing,omitempty"`
}

//UpdateOriginRequestHostname -
type UpdateOriginRequestHostname struct {
	Name string
}

//UpdateOriginResponse -
type UpdateOriginResponse struct {
	CustomerOriginID int
}

//Origin - Don't change variable names, otherwise API call to api server will get an error message as response.
type Origin struct {
	ID                int
	DirectoryName     string
	HostHeader        string
	HttpHostnames     []OriginHostname
	HttpLoadBalancing string
}

//OriginHostname -
type OriginHostname struct {
	Name string
}

//NewOriginAPIClient -
func NewOriginAPIClient(config *ClientConfig) *OriginAPIClient {
	apiClient := &OriginAPIClient{
		Config:        config,
		BaseAPIClient: config.BaseClientLegacy,
		AccountNumber: config.AccountNumber,
	}

	return apiClient
}

//AddOrigin -
func (c *OriginAPIClient) AddOrigin(origin *AddOriginRequest, mediaType string) (*AddOriginResponse, error) {
	log.Printf("AddOrigin>>origin:%v mediaType:%s", origin, mediaType)
	request, err := c.BaseAPIClient.BuildRequest("POST", fmt.Sprintf("v2/mcc/customers/%s/origins/%s", c.AccountNumber, mediaType), origin, false)
	if err != nil {
		return nil, fmt.Errorf("AddOrigin: %v", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("UpdateOrigin: %v", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("GetOrigin: %v", err)
	}

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
	if err != nil {
		return fmt.Errorf("DeleteOrigin: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteOrigin: %v", err)
	}

	return nil
}
