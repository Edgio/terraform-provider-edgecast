// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ConvertToStrings converts Terraform's
// TypeList and TypeSet collections into a []string.
func ConvertToStrings(v interface{}) ([]string, bool) {
	if listItems, ok := v.([]interface{}); ok {
		return ConvertSliceToStrings(listItems)
	}
	if set, ok := v.(*schema.Set); ok {
		setItems := set.List()
		return ConvertSliceToStrings(setItems)
	}
	return nil, false
}

// ConvertSliceToStrings converts a []interface{} to []string.
// Note that any items in v that are not strings will be excluded.
func ConvertSliceToStrings(v []interface{}) ([]string, bool) {
	if v == nil {
		return nil, false
	}
	strings := make([]string, len(v))
	for i, val := range v {
		if s, ok := val.(string); ok {
			strings[i] = s
		} else {
			return nil, false
		}
	}
	return strings, true
}

// ConvertSingletonSetToMap converts a interface{} that is actually a Terraform
// schema.TypeSet into a Map. This is useful when using the schema.TypeSet
// with MaxItems=1 workaround.
func ConvertSingletonSetToMap(attr interface{}) (map[string]interface{}, error) {
	if attr == nil {
		return nil, fmt.Errorf("set was nil")
	}

	// Must convert the set to a slice/array to work with the values
	if set, ok := attr.(*schema.Set); ok {
		arr := set.List()

		// Terraform will normally not allow more than one item because we
		// specify MaxItems=1, but we'll check again here
		if len(arr) != 1 {
			return nil, errors.New("set must contain exactly one item")
		}

		// The single item in the list
		// is expected to be a map[string][]interface{}
		if entryMap, ok := arr[0].(map[string]interface{}); ok {
			return entryMap, nil
		} else {
			return nil, fmt.Errorf(
				"%v must be of type map[string]interface{}, actual: %T",
				arr[0],
				arr[0])
		}
	}

	return nil, fmt.Errorf("attr was not a *schema.Set")
}

// NewTerraformSet is useful for unit testing when mocking Terraform
func NewTerraformSet(items []interface{}) *schema.Set {
	return schema.NewSet(dummySetFunc, items)
}

// ConvertToInt converts an interface{} to int
// If v is nil or not an int, 0 is returned
func ConvertToInt(v interface{}) int {
	i, _ := v.(int)
	return i
}

// ConvertToIntPointer converts an interface{} to *int
// If v is nil or not an integer, a nil pointer is returned
func ConvertToIntPointer(v interface{}, assertNaturalNumber bool) *int {
	if v == nil {
		return nil
	}
	if i, ok := v.(int); ok {
		if assertNaturalNumber && i <= 0 {
			return nil
		}
		return &i
	}
	return nil
}

// ConvertToBoolPointer converts an interface{} to *bool
// If v is nil or not a bool, a nil pointer is returned
func ConvertToBoolPointer(v interface{}) *bool {
	if v == nil {
		return nil
	}
	if b, ok := v.(bool); ok {
		return &b
	}
	return nil
}

// ConvertToString converts an interface{} to string
// If v is nil or not a string, an empty string is returned
func ConvertToString(v interface{}) string {
	s, _ := v.(string)
	return s
}

// ConvertToStringPointer converts an interface{} to *string
// If v is nil or not a string, a nil pointer is returned
func ConvertToStringPointer(v interface{}, excludeWhiteSpace bool) *string {
	if s, ok := v.(string); ok {
		if excludeWhiteSpace && len(strings.TrimSpace(s)) == 0 {
			return nil
		}
		return &s
	}
	return nil
}

// ConvertToStringSlicePointer converts an interface{} to *[]string
// If v is nil or not a string, a nil pointer is returned
func ConvertToStringsPointer(v interface{}, excludeEmpty bool) *[]string {
	if v == nil {
		return nil
	}
	if l, ok := v.([]interface{}); ok {
		if excludeEmpty && len(l) == 0 {
			return nil
		}
		s := make([]string, len(l))
		for i, v := range l {
			s[i] = v.(string)
		}
		return &s
	}
	return nil
}

// ConvertToMapPointer converts an interface{} to *map[string]string{}
// If v is nil or not a string, a nil pointer is returned
func ConvertToStringMapPointer(v interface{}, excludeEmpty bool) *map[string]string {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]interface{}); ok {
		if excludeEmpty && len(m) == 0 {
			return nil
		}
		result := make(map[string]string)
		for key, value := range m {
			stringValue, _ := value.(string)
			result[key] = stringValue
		}
		return &result
	}
	return nil
}

// StringIsNotEmptyJSON is a SchemaValidateFunc which tests to make sure the
// supplied string is not an empty JSON object i.e. "{}"
func StringIsNotEmptyJSON(
	i interface{},
	k string,
) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(
			errors,
			fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(v), &data); err != nil {
		errors = append(
			errors,
			fmt.Errorf("%q contains invalid JSON: %s", k, err))
		return warnings, errors
	}

	if len(data) == 0 {
		errors = append(
			errors,
			fmt.Errorf("%q contains empty JSON: '%s'", k, v))
	}

	return warnings, errors
}

// dummySetFunc is to be used when imitating Terraform
// in unit tests by using schema.NewSet
func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
