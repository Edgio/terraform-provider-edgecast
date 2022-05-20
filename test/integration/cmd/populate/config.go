package main

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"net/url"
	"os"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

func createConfig() edgecast.SDKConfig {
	config := edgecast.NewSDKConfig()
	config.APIToken = os.Getenv("API_TOKEN")
	config.BaseAPIURL = getEnvURL("API_ADDRESS")
	config.BaseAPIURLLegacy = getEnvURL("API_ADDRESS_LEGACY")
	config.BaseIDSURL = getEnvURL("IDS_ADDRESS")
	config.IDSCredentials = edgecast.IDSCredentials{
		ClientID:     os.Getenv("IDS_CLIENT_ID"),
		ClientSecret: os.Getenv("IDS_CLIENT_SECRET"),
		Scope:        os.Getenv("IDS_SCOPE"),
	}
	return config
}

func getEnvURL(s string) url.URL {
	u := internal.Check(url.Parse(os.Getenv(s))).(*url.URL)
	return *u
}
