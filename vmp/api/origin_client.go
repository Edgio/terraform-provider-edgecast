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

func (c *OriginApiClient) AddHttpLargeOrigin(origin *AddOriginRequest) (*AddOriginResponse, error) {
	request, err := c.BaseApiClient.BuildRequest("POST", fmt.Sprintf("mcc/customers/%s/origins/httplarge", c.AccountNumber), origin)

	parsedResponse := &AddOriginResponse{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	return parsedResponse, err
}

func (c *OriginApiClient) GetHttpLargeOrigin(id int) (*Origin, error) {
	request, err := c.BaseApiClient.BuildRequest("GET", fmt.Sprintf("mcc/customers/%s/origins/httplarge/%d", c.AccountNumber, id), nil)

	parsedResponse := &Origin{}

	_, err = c.BaseApiClient.SendRequest(request, &parsedResponse)

	return parsedResponse, err
}

func (c *OriginApiClient) DeleteOrigin(id int) error {
	request, err := c.BaseApiClient.BuildRequest("DELETE", fmt.Sprintf("mcc/customers/%s/origins/%d", c.AccountNumber, id), nil)

	_, err = c.BaseApiClient.SendRequest(request, nil)

	return err
}
