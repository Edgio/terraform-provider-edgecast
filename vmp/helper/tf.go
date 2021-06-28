package helper

// InterfaceToStringArray converts a interface{} whose underlying type is []string. Useful for parsing Terraform resource values.
func InterfaceToStringArray(attr interface{}) []string {
	if attr == nil {
		return nil
	}

	// Terraform's schema.TypeList stores values as []interface{}
	if interfaceArray, ok := attr.([]interface{}); ok {
		if values, ok := InterfaceArrayToStringArray(interfaceArray); ok {
			return values
		}
	}

	return nil
}

// InterfaceArrayToStringArray converts []interface{} to []string. Note that this only works if the underlying items are strings.
func InterfaceArrayToStringArray(interfaces []interface{}) ([]string, bool) {
	if interfaces == nil {
		return nil, false
	}

	strings := make([]string, len(interfaces))

	for i, v := range interfaces {
		if s, ok := v.(string); ok {
			strings[i] = s
		} else {
			return nil, false
		}
	}

	return strings, true
}
