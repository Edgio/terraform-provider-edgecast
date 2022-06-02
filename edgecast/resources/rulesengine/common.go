// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package rulesengine

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast/eclog"
	"terraform-provider-edgecast/edgecast/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/rulesengine"
)

// buildRulesEngineService builds the SDK Rules Engine service to manage Rule
// resources
func buildRulesEngineService(
	config api.ClientConfig,
) (*rulesengine.RulesEngineService, error) {

	idsCredentials := edgecast.IDSCredentials{
		ClientID:     config.IdsClientID,
		ClientSecret: config.IdsClientSecret,
		Scope:        config.IdsScope,
	}

	sdkConfig := edgecast.NewSDKConfig()
	sdkConfig.IDSCredentials = idsCredentials
	sdkConfig.BaseAPIURL = *config.APIURL
	sdkConfig.BaseIDSURL = *config.IdsURL
	sdkConfig.Logger = eclog.NewStandardLogger()

	return rulesengine.New(sdkConfig)
}
