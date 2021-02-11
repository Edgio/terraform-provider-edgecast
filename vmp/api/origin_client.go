// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"fmt"
)

type OriginApiClient struct {
	BaseApiClient *ApiClient
	AccountNumber string
}

type AddOriginRequest struct {
	DirectoryName     string
	HostHeader        string
	HttpHostnames     []AddOriginRequestHostname
	HttpsHostnames    []AddOriginRequestHostname
	HttpLoadBalancing string
}

type AddOriginRequestHostname struct {
	Name string
}

type AddOriginResponse struct {
	CustomerOriginId int
}

type UpdateOriginRequest struct {
	DirectoryName     string
	HostHeader        string
	HttpHostnames     []UpdateOriginRequestHostname
	HttpsHostnames    []UpdateOriginRequestHostname
	HttpLoadBalancing string
}

type UpdateOriginRequestHostname struct {
	Name string
}

type UpdateOriginResponse struct {
	CustomerOriginId int
}

type Origin struct {
	Id                int
	DirectoryName     string
	HostHeader        string
	HttpHostnames     []OriginHostname
	HttpLoadBalancing string
}

type OriginHostname struct {
	Name string
}

// NewOriginApiClient -
func NewOriginApiClient(baseApiClient *ApiClient, accountNumber string) *OriginApiClient {
	apiClient := &OriginApiClient{
		BaseApiClient: baseApiClient,
		AccountNumber: accountNumber,
	}

	return apiClient
}

func (c *OriginApiClient) AddOrigin(origin *AddOriginRequest, mediaType string) (*AddOriginResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("POST", fmt.Sprintf("v2/mcc/customers/%s/origins/%s", c.AccountNumber, mediaType), origin, false)
	parsedResponse := &AddOriginResponse{}
	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("AddOrigin: %v", err)
	}

	return parsedResponse, nil
}

func (c *OriginApiClient) UpdateOrigin(origin *UpdateOriginRequest, originID int, mediaType string) (*UpdateOriginResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("PUT", fmt.Sprintf("v2/mcc/customers/%s/origins/%s/%d", c.AccountNumber, mediaType, originID), origin, false)
	parsedResponse := &UpdateOriginResponse{}
	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("UpdateOrigin: %v", err)
	}

	return parsedResponse, nil
}

func (c *OriginApiClient) GetOrigin(id int, mediaType string) (*Origin, error) {
	request, err := c.BaseApiClient.BuildRequest("GET", fmt.Sprintf("v2/mcc/customers/%s/origins/%s/%d", c.AccountNumber, mediaType, id), nil, false)
	parsedResponse := &Origin{}
	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetOrigin: %v", err)
	}

	return parsedResponse, nil
}

func (c *OriginApiClient) DeleteOrigin(id int) error {
	request, err := c.BaseApiClient.BuildRequest("DELETE", fmt.Sprintf("v2/mcc/customers/%s/origins/%d", c.AccountNumber, id), nil, false)
	_, err = c.BaseApiClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteOrigin: %v", err)
	}

	return nil
}
