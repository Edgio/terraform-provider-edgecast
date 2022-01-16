// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package origin

import (
	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/auth"
	"github.com/EdgeCast/ec-sdk-go/edgecast/origin"
)

// buildOriginService builds the SDK Origin service to manage Origin
// resources
func buildOriginService(
	config api.ClientConfig,
) (*origin.OriginService, error) {

	idsCredentials := auth.OAuth2Credentials{
		ClientID:     config.IdsClientID,
		ClientSecret: config.IdsClientSecret,
		Scope:        config.IdsScope,
	}

	sdkConfig := edgecast.NewSDKConfig(config.APIToken, idsCredentials)

	sdkConfig.BaseAPIURL = *config.APIURL
	sdkConfig.BaseAPIURLLegacy = *config.APIURLLegacy
	sdkConfig.BaseIDSURL = *config.IdsURL

	return origin.New(sdkConfig)
}