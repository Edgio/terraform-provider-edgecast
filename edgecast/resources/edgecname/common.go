// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package edgecname

import (
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/edgecname"
)

// buildEdgeCnameService builds the SDK Edge CNAME service to manage Edge CNAME
// resources
func buildEdgeCnameService(
	config internal.ProviderConfig,
) (*edgecname.EdgeCnameService, error) {

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

	return edgecname.New(sdkConfig)
}
