// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/eclog"
)

type Set map[string]bool

func (s Set) ToList() []string {
	keys := make([]string, 0)
	for key := range s {
		keys = append(keys, key)
	}

	return keys
}

type Config struct {
	// The account number to run the tests for.
	AccountNumber string

	// Indicates which resources should be populated
	PopulateFlags Set

	// A timestamped test email used for testing purposes.
	TestEmail string

	// The SDK configuration.
	SDKConfig edgecast.SDKConfig

	// The API Token for Customer creation
	APITokenPCC string
}

// NewConfig creates a new instance of Config.
func NewConfig() Config {
	cfg := Config{
		AccountNumber: getEnvRequired(envAccountNumber),
		TestEmail:     getTestEmail(),
		PopulateFlags: make(map[string]bool),
		SDKConfig:     createSDKConfig(),
		APITokenPCC:   getEnvRequired(envAPITokenPCC),
	}

	populateOnly := internal.GetEnvWithDefault(envPopulateOnly, "")
	if len(populateOnly) > 0 {
		for _, s := range strings.Split(populateOnly, ",") {
			cfg.PopulateFlags[s] = true
		}
	}

	return cfg
}

func createSDKConfig() edgecast.SDKConfig {
	config := edgecast.NewSDKConfig()
	config.APIToken = getEnvRequired(envAPIToken)
	config.BaseAPIURL = getEnvURLRequired(envAPIAddress)
	config.BaseAPIURLLegacy = getEnvURLRequired(envAPIAddressLegacy)
	config.BaseIDSURL = getEnvURLRequired(envIDSAddress)
	config.IDSCredentials = edgecast.IDSCredentials{
		ClientID:     getEnvRequired(envIDSClientID),
		ClientSecret: getEnvRequired(envIDSClientSecret),
		Scope:        getEnvRequired(envIDSScope),
	}

	logFile, err := internal.GetEnv(envLogFile)
	if err == nil && len(logFile) > 0 {
		config.Logger = eclog.NewFileLogger(logFile)
	} else {
		config.Logger = eclog.NewStandardLogger()
	}

	return config
}

func getEnvRequired(key string) string {
	return internal.Check(internal.GetEnv(key))
}

func getEnvURLRequired(key string) url.URL {
	return *internal.Check(internal.GetEnvURL(key))
}

func getTestEmail() string {
	return fmt.Sprintf(formatTestEmail, time.Now().Unix())
}
