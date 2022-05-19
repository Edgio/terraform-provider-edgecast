package helper_test

import (
	"math"
	"reflect"
	"terraform-provider-edgecast/ec/helper"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestConvertTFCollectionToStrings(t *testing.T) {
	cases := []struct {
		name       string
		input      interface{}
		expected   []string
		expectedOk bool
	}{
		{
			name:       "Happy path - all strings",
			input:      []interface{}{"val1", "val2", "val3"},
			expected:   []string{"val1", "val2", "val3"},
			expectedOk: true,
		},
		{
			name:       "Empty list input",
			input:      make([]interface{}, 0),
			expected:   make([]string, 0),
			expectedOk: true,
		},
		{
			name:       "Nil input",
			input:      nil,
			expected:   nil,
			expectedOk: false,
		},
	}
	for _, v := range cases {
		actual, err := helper.ConvertTFCollectionToStrings(v.input)

		if v.expectedOk {
			if err == nil {
				if !reflect.DeepEqual(v.expected, actual) {
					t.Fatalf(
						"Case '%s': Expected %q but got %q",
						v.name,
						v.expected,
						actual)
				}
			} else {
				t.Fatalf("unexpected error: %w", err)
			}
		} else {
			if err == nil {
				t.Fatalf("Case '%s': Expected an error but got none", v.name)
			}
		}
	}
}

func TestExpandSingletonSet(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expected      *map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: schema.NewSet(
				dummySetFunc,
				[]interface{}{map[string]interface{}{
					"stringProperty": "string value",
					"intProperty":    1,
					"arrayProperty":  []string{"one", "two"},
				}}),
			expected: &map[string]interface{}{
				"stringProperty": "string value",
				"intProperty":    1,
				"arrayProperty":  []string{"one", "two"},
			},
			expectSuccess: true,
		},
		{
			name: "Error path - set has more than one map",
			input: schema.NewSet(dummySetFunc, []interface{}{
				map[string]interface{}{
					"map1prop2": "string value",
				},
				map[string]interface{}{
					"map2prop1": "string value",
				}}),
			expected:      nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - nil input",
			input:         nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "a string",
			expectSuccess: false,
		},
	}
	for _, v := range cases {
		actual, err := helper.ConvertSingletonSetToMap(v.input)
		if v.expectSuccess {
			if err == nil {
				if !reflect.DeepEqual(*v.expected, actual) {
					t.Fatalf(
						"Case '%s': Expected %v but got %v",
						v.name,
						*v.expected,
						actual)
				}
			} else {
				t.Fatalf("Case '%s': Encountered error: %v", v.name, err)
			}
		} else {
			if err == nil {
				t.Fatalf("Case '%s': Expected error, but got no error", v.name)
			}
		}
	}
}

func TestConvertToInt(t *testing.T) {
	cases := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name:     "Happy path",
			input:    1,
			expected: 1,
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: 0,
		},
		{
			name:     "Non-Integer input",
			input:    "not a number",
			expected: 0,
		},
		{
			name:     "Max Int value",
			input:    math.MaxInt32,
			expected: math.MaxInt32,
		},
		{
			name:     "Min Int value",
			input:    math.MinInt32,
			expected: math.MinInt32,
		},
	}
	for _, v := range cases {
		actual := helper.ConvertToInt(v.input)
		if v.expected != actual {
			t.Fatalf(
				"Case '%s': Expected %q but got %q",
				v.name,
				v.expected,
				actual)
		}
	}
}

func TestConvertToIntPointer(t *testing.T) {
	cases := []struct {
		name                string
		input               interface{}
		expectedPtrValue    int
		expectNilPtr        bool
		assertNaturalNumber bool
	}{
		{
			name:             "Happy path",
			input:            1,
			expectedPtrValue: 1,
			expectNilPtr:     false,
		},
		{
			name:         "Nil input",
			input:        nil,
			expectNilPtr: true,
		},
		{
			name:                "Non-natural number error",
			input:               -1,
			expectNilPtr:        true,
			assertNaturalNumber: true,
		},
		{
			name:         "Non-Integer input",
			input:        "not a number",
			expectNilPtr: true,
		},
		{
			name:             "Max Int value",
			input:            math.MaxInt32,
			expectedPtrValue: math.MaxInt32,
			expectNilPtr:     false,
		},
		{
			name:             "Min Int value",
			input:            math.MinInt32,
			expectedPtrValue: math.MinInt32,
			expectNilPtr:     false,
		},
	}
	for _, v := range cases {
		actual := helper.ConvertToIntPointer(v.input, v.assertNaturalNumber)
		if v.expectNilPtr {
			if actual != nil {
				t.Fatalf("Case '%s': Expected nil pointer", v.name)
			}
		} else {
			if actual == nil {
				t.Fatalf("Case '%s': Expected non-nil pointer", v.name)
			} else {
				if v.expectedPtrValue != *actual {
					t.Fatalf(
						"Case '%s': Expected %d but got %d",
						v.name,
						v.expectedPtrValue,
						*actual)
				}
			}
		}
	}
}

func TestConvertToBoolPointer(t *testing.T) {
	cases := []struct {
		name             string
		input            interface{}
		expectedPtrValue bool
		expectNilPtr     bool
	}{
		{
			name:             "Happy path",
			input:            true,
			expectedPtrValue: true,
			expectNilPtr:     false,
		},
		{
			name:         "Nil input",
			input:        nil,
			expectNilPtr: true,
		},
		{
			name:         "Non-Boolean input",
			input:        "not a bool",
			expectNilPtr: true,
		},
	}
	for _, v := range cases {
		actual := helper.ConvertToBoolPointer(v.input)
		if v.expectNilPtr {
			if actual != nil {
				t.Fatalf("Case '%s': Expected nil pointer", v.name)
			}
		} else {
			if actual == nil {
				t.Fatalf("Case '%s': Expected non-nil pointer", v.name)
			} else {
				if v.expectedPtrValue != *actual {
					t.Fatalf(
						"Case '%s': Expected %t but got %t",
						v.name,
						v.expectedPtrValue,
						*actual)
				}
			}
		}
	}
}

func TestConvertToString(t *testing.T) {
	cases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Happy path",
			input:    "hello world!",
			expected: "hello world!",
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "Non-String input",
			input:    -1,
			expected: "",
		},
	}
	for _, v := range cases {
		actual := helper.ConvertToString(v.input)
		if v.expected != actual {
			t.Fatalf(
				"Case '%s': Expected %s but got %s",
				v.name,
				v.expected,
				actual)
		}
	}
}

func TestConvertToStringPointer(t *testing.T) {
	cases := []struct {
		name              string
		input             interface{}
		expectedPtrValue  string
		expectNilPtr      bool
		excludeWhiteSpace bool
	}{
		{
			name:             "Happy path",
			input:            "hello world!",
			expectedPtrValue: "hello world!",
			expectNilPtr:     false,
		},
		{
			name:         "Nil input",
			input:        nil,
			expectNilPtr: true,
		},
		{
			name:         "Non-String input",
			input:        -1,
			expectNilPtr: true,
		},
		{
			name:              "Whitespace allowed",
			input:             "    ",
			expectedPtrValue:  "    ",
			expectNilPtr:      false,
			excludeWhiteSpace: false,
		},
		{
			name:              "Whitespace not allowed",
			input:             "    ",
			expectNilPtr:      true,
			excludeWhiteSpace: true,
		},
	}
	for _, v := range cases {
		actual := helper.ConvertToStringPointer(v.input, v.excludeWhiteSpace)
		if v.expectNilPtr {
			if actual != nil {
				t.Fatalf(
					"Case '%s': Expected nil pointer, but contains '%s'",
					v.name,
					*actual)
			}
		} else {
			if actual == nil {
				t.Fatalf("Case '%s': Expected non-nil pointer", v.name)
			} else {
				if v.expectedPtrValue != *actual {
					t.Fatalf(
						"Case '%s': Expected %s but got %s",
						v.name,
						v.expectedPtrValue,
						*actual)
				}
			}
		}
	}
}

func TestConvertToStringsPointer(t *testing.T) {
	cases := []struct {
		name             string
		input            interface{}
		expectedPtrValue []string
		expectNilPtr     bool
		excludeEmpty     bool
	}{
		{
			name:             "Happy path",
			input:            []interface{}{"hello", "world"},
			expectedPtrValue: []string{"hello", "world"},
			expectNilPtr:     false,
		},
		{
			name:         "Nil input",
			input:        nil,
			expectNilPtr: true,
		},
		{
			name:         "Non-Slice input",
			input:        -1,
			expectNilPtr: true,
		},
		{
			name:             "Empty slice allowed",
			input:            make([]interface{}, 0),
			expectedPtrValue: make([]string, 0),
			expectNilPtr:     false,
			excludeEmpty:     false,
		},
		{
			name:         "Empty slice not allowed",
			input:        make([]interface{}, 0),
			expectNilPtr: true,
			excludeEmpty: true,
		},
	}
	for _, v := range cases {
		actual := helper.ConvertToStringsPointer(v.input, v.excludeEmpty)
		if v.expectNilPtr {
			if actual != nil {
				t.Fatalf("Case '%s': Expected nil pointer", v.name)
			}
		} else {
			if actual == nil {
				t.Fatalf("Case '%s': Expected non-nil pointer", v.name)
			} else {
				if !reflect.DeepEqual(v.expectedPtrValue, *actual) {
					t.Fatalf(
						"Case '%s': Expected %v but got %v",
						v.name,
						v.expectedPtrValue,
						*actual)
				}
			}
		}
	}
}

func TestConvertToStringMapPointer(t *testing.T) {
	cases := []struct {
		name             string
		input            interface{}
		expectedPtrValue map[string]string
		expectNilPtr     bool
		excludeEmpty     bool
	}{
		{
			name:             "Happy path",
			input:            map[string]interface{}{"hello": "world"},
			expectedPtrValue: map[string]string{"hello": "world"},
			expectNilPtr:     false,
		},
		{
			name:         "Nil input",
			input:        nil,
			expectNilPtr: true,
		},
		{
			name:         "Non-Map input",
			input:        -1,
			expectNilPtr: true,
		},
		{
			name:             "Empty map allowed",
			input:            make(map[string]interface{}),
			expectedPtrValue: make(map[string]string),
			expectNilPtr:     false,
			excludeEmpty:     false,
		},
		{
			name:         "Empty map not allowed",
			input:        make(map[string]interface{}),
			expectNilPtr: true,
			excludeEmpty: true,
		},
	}
	for _, v := range cases {
		actual := helper.ConvertToStringMapPointer(v.input, v.excludeEmpty)
		if v.expectNilPtr {
			if actual != nil {
				t.Fatalf("Case '%s': Expected nil pointer", v.name)
			}
		} else {
			if actual == nil {
				t.Fatalf("Case '%s': Expected non-nil pointer", v.name)
			} else {
				if !reflect.DeepEqual(v.expectedPtrValue, *actual) {
					t.Fatalf(
						"Case '%s': Expected %v but got %v",
						v.name,
						v.expectedPtrValue,
						*actual)
				}
			}
		}
	}
}

// dummySetFunc is to be used when imitating Terraform
// in unit tests by using schema.NewSet
func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
