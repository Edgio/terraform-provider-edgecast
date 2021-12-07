// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper

// IsInterfaceArray deterimins if an interface{} is actually an []interface{}
func IsInterfaceArray(input interface{}) bool {
	switch input.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}
