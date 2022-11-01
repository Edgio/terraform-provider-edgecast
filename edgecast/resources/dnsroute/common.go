// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package dnsroute

import (
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
)

// buildRouteDNSService builds the SDK Route DNS service to manage DNS
// resources
func buildRouteDNSService(
	config internal.ProviderConfig,
) (*routedns.RouteDNSService, error) {

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

	return routedns.New(sdkConfig)
}
