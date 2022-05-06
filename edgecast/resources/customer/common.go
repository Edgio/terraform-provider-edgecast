// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package customer

import (
	"terraform-provider-edgecast/edgecast/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/customer"
)

// buildCustomerService builds the SDK Customer service to manage Customer
// resources
func buildCustomerService(
	config api.ClientConfig,
) (*customer.CustomerService, error) {

	idsCredentials := edgecast.IDSCredentials{
		ClientID:     config.IdsClientID,
		ClientSecret: config.IdsClientSecret,
		Scope:        config.IdsScope,
	}

	sdkConfig := edgecast.NewSDKConfig()
	sdkConfig.APIToken = config.APIToken
	sdkConfig.IDSCredentials = idsCredentials
	sdkConfig.BaseAPIURL = *config.APIURL
	sdkConfig.BaseAPIURLLegacy = *config.APIURLLegacy
	sdkConfig.BaseIDSURL = *config.IdsURL

	return customer.New(sdkConfig)
}
