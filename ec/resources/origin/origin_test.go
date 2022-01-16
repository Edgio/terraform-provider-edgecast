package origin

import (
	"reflect"
	"sort"
	"terraform-provider-ec/ec/helper"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/origin"
)

func TestExpandHostname(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]origin.Hostname
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"name":       "http://origin1.customer.com",
					"is_primary": 1,
					"ordinal":    0,
				},
				map[string]interface{}{
					"name":       "http://origin2.customer.com",
					"is_primary": 0,
					"ordinal":    1,
				},
			}),
			expectedPtr: &[]origin.Hostname{
				{
					Name:      "http://origin1.customer.com",
					IsPrimary: 1,
					Ordinal:   0,
				},
				{
					Name:      "http://origin2.customer.com",
					IsPrimary: 0,
					Ordinal:   1,
				},
			},
			expectSuccess: true,
		},
		{
			name:          "Happy path - None Defined",
			input:         helper.NewTerraformSet([]interface{}{}),
			expectedPtr:   &[]origin.Hostname{},
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ExpandHostname(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				// array equality depends on order, sort by Name
				sort.Slice(actual, func(i, j int) bool {
					return actual[i].Name < actual[j].Name
				})

				if !reflect.DeepEqual(expected, actual) {
					t.Fatalf(
						"%s: Expected %+v but got %+v",
						v.name, expected, actual,
					)
				}
			} else {
				t.Fatalf(
					"%s: Encountered error where one was not expected: %+v",
					v.name,
					err,
				)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}

func TestExpandShieldPOPs(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]origin.ShieldPOP
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"pop_code": "ABC",
				},
				map[string]interface{}{
					"pop_code": "XYZ",
				},
			}),
			expectedPtr: &[]origin.ShieldPOP{
				{
					POPCode: "ABC",
				},
				{
					POPCode: "XYZ",
				},
			},
			expectSuccess: true,
		},
		{
			name:          "Happy path - None Defined",
			input:         helper.NewTerraformSet([]interface{}{}),
			expectedPtr:   &[]origin.ShieldPOP{},
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ExpandShieldPOPs(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				// array equality depends on order, sort by POPCode
				sort.Slice(actual, func(i, j int) bool {
					return actual[i].POPCode < actual[j].POPCode
				})

				if !reflect.DeepEqual(expected, actual) {
					t.Fatalf(
						"%s: Expected %+v but got %+v",
						v.name,
						expected,
						actual,
					)
				}
			} else {
				t.Fatalf(
					"%s: Encountered error where one was not expected: %+v",
					v.name,
					err,
				)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}

func TestFlattenHostname(t *testing.T) {

	cases := []struct {
		name     string
		input    []origin.Hostname
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []origin.Hostname{
				{
					Name:      "http://origin1.customer.com",
					IsPrimary: 1,
					Ordinal:   0,
				},
				{
					Name:      "http://origin2.customer.com",
					IsPrimary: 0,
					Ordinal:   1,
				},
			},
			expected: []map[string]interface{}{
				{
					"name":       "http://origin1.customer.com",
					"is_primary": 1,
					"ordinal":    0,
				},
				{
					"name":       "http://origin2.customer.com",
					"is_primary": 0,
					"ordinal":    1,
				},
			},
		},
		{
			name:     "Empty collection",
			input:    []origin.Hostname{},
			expected: []map[string]interface{}{},
		},
	}

	for _, v := range cases {

		actual := FlattenHostname(v.input)

		if !reflect.DeepEqual(actual, v.expected) {
			t.Fatalf("%s: Expected %+v but got %+v", v.name, actual, v.expected)
		}
	}
}

func TestFlattenShieldPOP(t *testing.T) {

	cases := []struct {
		name     string
		input    []origin.ShieldPOP
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []origin.ShieldPOP{
				{
					POPCode: "ABC",
				},
				{
					POPCode: "XYZ",
				},
			},
			expected: []map[string]interface{}{
				{
					"pop_code": "ABC",
				},
				{
					"pop_code": "XYZ",
				},
			},
		},
		{
			name:     "Empty collection",
			input:    []origin.ShieldPOP{},
			expected: []map[string]interface{}{},
		},
	}

	for _, v := range cases {

		actual := FlattenShieldPOP(v.input)

		if !reflect.DeepEqual(actual, v.expected) {
			t.Fatalf("%s: Expected %+v but got %+v", v.name, actual, v.expected)
		}
	}
}
