// Copyright Edgecast, Licensed under the terms of the Apache 2.0 license. See LICENSE file in project root for terms.
package helper

import (
	"errors"
	"fmt"
	"math"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
	ConvertInterfaceToStringArray takes the following two types and converts them to []string:
		- []interface{} where every item in the slice is a string
		- (Terraform) *schema.Set whose elements are strings
*/
func ConvertInterfaceToStringArray(attr interface{}) (*[]string, bool) {
	if attr == nil {
		return nil, false
	}

	// Terraform's schema.TypeList stores values as []interface{}
	if interfaceArray, ok := attr.([]interface{}); ok {
		if values, ok := ConvertInterfaceArrayToStringArray(interfaceArray); ok {
			return &values, true
		}
	}

	// Terraform's schema.TypeSet stores values as *schema.Set, which holds the list internally
	if set, ok := attr.(*schema.Set); ok {

		items := set.List()
		strings := make([]string, len(items))

		for i := range items {
			if v, ok := items[i].(string); ok {
				strings[i] = v
			} else {
				return nil, false
			}
		}

		return &strings, true
	}

	return nil, false
}

// ConvertInterfaceArrayToStringArray converts []interface{} to []string. Note that this only works if the underlying items are strings.
func ConvertInterfaceArrayToStringArray(interfaces []interface{}) ([]string, bool) {
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

// GetMapFromSet converts a interface{} that is actually a Terraform schema.TypeSet into a Map
// This is useful when using the schem.TypeSet with MaxItems=1 workaround
func GetMapFromSet(attr interface{}) (map[string]interface{}, error) {
	if attr == nil {
		return nil, fmt.Errorf("set was nil")
	}

	// Must convert the set to a slice/array to work with the values
	if set, ok := attr.(*schema.Set); ok {
		arr := set.List()

		// Terraform will normally not allow more than one item because we specify MaxItems=1, but we'll check again here
		if len(arr) != 1 {
			return nil, errors.New("set must contain exactly one item")
		}

		// The single item in the list is expected to be a map[string][]interface{}
		if entryMap, ok := arr[0].(map[string]interface{}); ok {
			return entryMap, nil
		} else {
			return nil, fmt.Errorf("%v must be of type map[string]interface{}, actual: %T", arr[0], arr[0])
		}
	}

	return nil, fmt.Errorf("attr was not a *schema.Set")
}

// NewTerraformSet is useful for unit testing when mocking Terraform
func NewTerraformSet(items []interface{}) *schema.Set {
	return schema.NewSet(dummySetFunc, items)
}

// dummySetFunc is to be used when imitating Terraform in unit tests by using schema.NewSet
func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
