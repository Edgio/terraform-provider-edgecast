// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package internal

import (
	"fmt"
	"net/url"
	"os"
)

// GetEnv retrieves an environment variable. If it does not exist or is empty,
// it returns an empty string and an error.
func GetEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}

	return "", fmt.Errorf("env variable not found or empty: %s", key)
}

// GetEnvWithDefault retrives an environment variable. If it does not exist,
// it returns the provided default.
func GetEnvWithDefault(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// GetEnvWithDefault retrives a URL set as an environment variable. If it does
// not exist, it returns nil and an error.
func GetEnvURL(key string) (*url.URL, error) {
	v, err := GetEnv(key)
	if err != nil {
		return nil, err
	}

	return url.Parse(v)
}
