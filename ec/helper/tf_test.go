package helper_test

import (
	"math"
	"reflect"
	"terraform-provider-ec/ec/helper"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandTerraformStrings(t *testing.T) {
	cases := []struct {
		input      interface{}
		expected   []string
		expectedOk bool
	}{
		{
			input:      []interface{}{"val1", "val2", "val3"},
			expected:   []string{"val1", "val2", "val3"},
			expectedOk: true,
		},
		{
			input:      make([]interface{}, 0),
			expected:   make([]string, 0),
			expectedOk: true,
		},
		{
			input:      nil,
			expected:   nil,
			expectedOk: false,
		},
	}
	for _, v := range cases {
		actual, ok := helper.ExpandTerraformStrings(v.input)
		if ok == v.expectedOk {
			if ok {
				if !helper.IsStringSliceEqual(v.expected, actual) {
					t.Fatalf("Expected %q but got %q", v.expected, actual)
				}
			}
		} else {
			t.Fatalf("Expected ok result of %t but got %t", v.expectedOk, ok)
		}
	}
}

func TestExpandStrings(t *testing.T) {
	cases := []struct {
		input    []interface{}
		expected ExpandStringsResult
	}{
		{
			input: []interface{}{"val1", "val2", "val3"},
			expected: ExpandStringsResult{
				Array: []string{"val1", "val2", "val3"},
				Ok:    true,
			},
		},
		{
			input: []interface{}{"val1", 1, "val3"},
			expected: ExpandStringsResult{
				Array: nil,
				Ok:    false,
			},
		},
		{
			input: make([]interface{}, 0),
			expected: ExpandStringsResult{
				Array: make([]string, 0),
				Ok:    true,
			},
		},
		{
			input: nil,
			expected: ExpandStringsResult{
				Array: nil,
				Ok:    false,
			},
		},
	}
	for _, v := range cases {
		actual, ok := helper.ExpandStrings(v.input)
		if !helper.IsStringSliceEqual(v.expected.Array, actual) {
			t.Fatalf("Expected %q but got %q", v.expected.Array, actual)
		}
		if v.expected.Ok != ok {
			t.Fatalf("Expected %t but got %t", v.expected.Ok, ok)
		}
	}
}

type ExpandStringsResult struct {
	Array []string
	Ok    bool
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
		actual, err := helper.ExpandSingletonSet(v.input)
		if v.expectSuccess {
			if err == nil {
				if !reflect.DeepEqual(*v.expected, actual) {
					t.Fatalf(
						"%s: Expected %v but got %v",
						v.name,
						*v.expected,
						actual)
				}
			} else {
				t.Fatalf("%s: Encountered error: %v", v.name, err)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}

func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
