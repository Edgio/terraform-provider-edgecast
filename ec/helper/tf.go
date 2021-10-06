// Copyright Edgecast, Licensed under the terms of the Apache 2.0 license. See LICENSE file in project root for terms.
package helper

import (
	"errors"
	"fmt"
	"math"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ExpandTerraformStrings converts Terraform's
// TypeList and TypeSet collections into a []string.
func ExpandTerraformStrings(v interface{}) ([]string, bool) {
	if listItems, ok := v.([]interface{}); ok {
		return ExpandStrings(listItems)
	}
	if set, ok := v.(*schema.Set); ok {
		setItems := set.List()
		return ExpandStrings(setItems)
	}
	return nil, false
}

// ExpandStrings converts a []interface{} to []string.
// Note that any items in v that are not strings will be excluded.
func ExpandStrings(v []interface{}) ([]string, bool) {
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

// ExpandSet converts a interface{} that is actually a Terraform
// schema.TypeSet into a Map. This is useful when using the schema.TypeSet
// with MaxItems=1 workaround.
func ExpandSingletonSet(attr interface{}) (map[string]interface{}, error) {
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

// ExpandInt converts an interface{} to int
// If v is nil or not an int, 0 is returned
func ExpandInt(v interface{}) int {
	i, _ := v.(int)
	return i
}

// ExpandIntPointer converts an interface{} to *int
// If v is nil or not an integer, a nil pointer is returned
func ExpandIntPointer(v interface{}, assertNaturalNumber bool) *int {
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

// ExpandBoolPointer converts an interface{} to *bool
// If v is nil or not a bool, a nil pointer is returned
func ExpandBoolPointer(v interface{}) *bool {
	if v == nil {
		return nil
	}
	if b, ok := v.(bool); ok {
		return &b
	}
	return nil
}

// ExpandString converts an interface{} to string
// If v is nil or not a string, an empty string is returned
func ExpandString(v interface{}) string {
	s, _ := v.(string)
	return s
}

// ExpandStringPointer converts an interface{} to *string
// If v is nil or not a string, a nil pointer is returned
func ExpandStringPointer(v interface{}, excludeWhiteSpace bool) *string {
	if s, ok := v.(string); ok {
		if excludeWhiteSpace && len(s) == 0 {
			return nil
		}
		return &s
	}
	return nil
}

// ExpandStringSlicePointer converts an interface{} to *[]string
// If v is nil or not a string, a nil pointer is returned
func ExpandStringsPointer(v interface{}, excludeEmpty bool) *[]string {
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

// ExpandMapPointer converts an interface{} to *map[string]string{}
// If v is nil or not a string, a nil pointer is returned
func ExpandStringMapPointer(v interface{}, excludeEmpty bool) *map[string]string {
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

// dummySetFunc is to be used when imitating Terraform
// in unit tests by using schema.NewSet
func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
