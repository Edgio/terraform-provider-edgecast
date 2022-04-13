// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package waf

import (
	"terraform-provider-edgecast/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

// buildWAFService builds the SDK WAF service to managed WAF resources
func buildWAFService(config api.ClientConfig) (*sdkwaf.WAFService, error) {

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

	return sdkwaf.New(sdkConfig)
}
