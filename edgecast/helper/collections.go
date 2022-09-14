// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper

import (
	"errors"
)

// IsInterfaceArray deterimins if an interface{} is actually an []interface{}.
func IsInterfaceArray(input interface{}) bool {
	switch input.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

// ConvertSliceToStrings converts a []interface{} to []string.
func ConvertSliceToStrings(v []interface{}) ([]string, error) {
	if v == nil {
		return nil, nil
	}

	strings := make([]string, len(v))

	for i, val := range v {
		if s, ok := val.(string); ok {
			strings[i] = s
		} else {
			return nil,
				errors.New("slice contained a non-string value")
		}
	}

	return strings, nil
}

// GetStringFromMap returns a string value from the provided map using the
// provided key. If the item is not in the map or it is not a string, this
// function will return an empty string and a false 'ok' value.
func GetStringFromMap(m map[string]any, key string) (string, bool) {
	raw, ok := m[key]
	if !ok {
		return "", false
	}

	val, ok := raw.(string)
	return val, ok
}

// GetBoolFromMap returns a bool value from the provided map using the
// provided key. If the item is not in the map or it is not a bool, this
// function will return a false 'ok' value.
func GetBoolFromMap(m map[string]any, key string) (bool, bool) {
	raw, ok := m[key]
	if !ok {
		return false, false
	}

	val, ok := raw.(bool)
	return val, ok
}
