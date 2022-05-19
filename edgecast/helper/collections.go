// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper

import "errors"

// IsInterfaceArray deterimins if an interface{} is actually an []interface{}
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
