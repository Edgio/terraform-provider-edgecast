// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package waf

import (
	"terraform-provider-ec/ec/api"

	sdk "github.com/EdgeCast/ec-sdk-go/edgecast"
	sdkclient "github.com/EdgeCast/ec-sdk-go/edgecast/client"
	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

// buildWAFService builds the SDK WAF service to managed WAF resources
func buildWAFService(config api.ClientConfig) (*sdkwaf.WAFService, error) {
	/*
		TODO SDK Re-work: unify WAFConfig with SDK Config

		This is not dev-friendly... should not have to construct the WAF Service in this way

		Something like this would be ideal:
			sdkConfig := sdk.Config{...}
			wafService := waf.New(sdkConfig)

		The below workaround will work in the meantime.
	*/
	authProvider, err := sdkclient.NewLegacyAuthorizationHeaderProvider(config.APIToken)

	if err != nil {
		return nil, err
	}

	logger := sdk.NewStandardLogger()

	return &sdkwaf.WAFService{
		Logger: logger,
		Client: sdkclient.NewClient(sdkclient.ClientConfig{
			AuthHeaderProvider: authProvider,
			BaseAPIURL:         config.APIURLLegacy,
			UserAgent:          config.UserAgent,
			Logger:             logger,
		}),
	}, nil
}
