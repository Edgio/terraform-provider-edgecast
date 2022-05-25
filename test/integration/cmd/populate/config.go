package main

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/eclog"
	"net/url"
	"os"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

func createConfig() edgecast.SDKConfig {
	config := edgecast.NewSDKConfig()
	config.APIToken = os.Getenv("API_TOKEN")
	config.BaseAPIURL = envURL("API_ADDRESS")
	config.BaseAPIURLLegacy = envURL("API_ADDRESS_LEGACY")
	config.BaseIDSURL = envURL("IDS_ADDRESS")
	config.IDSCredentials = edgecast.IDSCredentials{
		ClientID:     os.Getenv("IDS_CLIENT_ID"),
		ClientSecret: os.Getenv("IDS_CLIENT_SECRET"),
		Scope:        os.Getenv("IDS_SCOPE"),
	}
	config.Logger = eclog.NewStandardLogger()
	return config
}

func envURL(s string) url.URL {
	return *internal.Check(url.Parse(os.Getenv(s)))
}
