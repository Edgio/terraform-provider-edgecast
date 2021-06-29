// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package helper

// InterfaceSliceEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func IsInterfaceSliceEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// StringSliceEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func IsStringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// IsInterfaceArray deterimins if an interface{} is actually an []interface{}
func IsInterfaceArray(input interface{}) bool {
	switch input.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

// MapEqual tells whether a and b contain the same key-value pairs.
// A nil argument is equivalent to an empty map.
func IsMapEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if v != b[k] {
			return false
		}
	}
	return true
}
