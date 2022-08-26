// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"terraform-provider-edgecast/edgecast/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps"
)

// buildCPSService builds the SDK CPS service to manage CPS
// resources
func buildCPSService(
	config api.ClientConfig,
) (*cps.CpsService, error) {

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
	sdkConfig.UserAgent = config.UserAgent

	return cps.New(sdkConfig)
}
