// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package edgecname

import (
	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/auth"
	"github.com/EdgeCast/ec-sdk-go/edgecast/edgecname"
)

// buildCustomerService builds the SDK Customer service to manage Customer
// resources
func buildEdgeCnameService(
	config api.ClientConfig,
) (*edgecname.EdgeCnameService, error) {

	idsCredentials := auth.OAuth2Credentials{
		ClientID:     config.IdsClientID,
		ClientSecret: config.IdsClientSecret,
		Scope:        config.IdsScope,
	}

	sdkConfig := edgecast.NewSDKConfig(config.APIToken, idsCredentials)

	sdkConfig.BaseAPIURL = *config.APIURL
	sdkConfig.BaseAPIURLLegacy = *config.APIURLLegacy
	sdkConfig.BaseIDSURL = *config.IdsURL

	return edgecname.New(sdkConfig)
}
