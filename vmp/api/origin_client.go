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

func NewOriginApiClient(baseApiClient *ApiClient, accountNumber string) *OriginApiClient {
	apiClient := &OriginApiClient{
		BaseApiClient: baseApiClient,
		AccountNumber: accountNumber,
	}

	return apiClient
}

func (c *OriginApiClient) AddOrigin(origin *AddOriginRequest, mediaType string) (*AddOriginResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("POST", fmt.Sprintf("mcc/customers/%s/origins/%s", c.AccountNumber, mediaType), origin, false)
	InfoLogger.Printf("Add origin in %s [POST] Url: %s\n", mediaType, request.URL)

	parsedResponse := &AddOriginResponse{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	return parsedResponse, err
}

func (c *OriginApiClient) UpdateOrigin(origin *UpdateOriginRequest, originID int, mediaType string) (*UpdateOriginResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("PUT", fmt.Sprintf("mcc/customers/%s/origins/%s/%d", c.AccountNumber, mediaType, originID), origin, false)
	InfoLogger.Printf("Update origin in %s [PUT] Url: %s\n", mediaType, request.URL)

	parsedResponse := &UpdateOriginResponse{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	return parsedResponse, err
}

func (c *OriginApiClient) GetOrigin(id int, mediaType string) (*Origin, error) {
	request, err := c.BaseApiClient.BuildRequest("GET", fmt.Sprintf("mcc/customers/%s/origins/%s/%d", c.AccountNumber, mediaType, id), nil, false)
	InfoLogger.Printf("Get origin in %s [GET] Url: %s\n", mediaType, request.URL)

	parsedResponse := &Origin{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	return parsedResponse, err
}

func (c *OriginApiClient) DeleteOrigin(id int) error {
	request, err := c.BaseApiClient.BuildRequest("DELETE", fmt.Sprintf("mcc/customers/%s/origins/%d", c.AccountNumber, id), nil, false)
	InfoLogger.Printf("DeleteOrigin [DELETE] Url: %s\n", request.URL)

	_, err = c.BaseApiClient.SendRequest(request, nil)

	return err
}
