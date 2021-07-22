package helper_test

import (
	"math"
	"reflect"
	"terraform-provider-vmp/vmp/helper"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestInterfaceToStringArray(t *testing.T) {

	cases := []struct {
		input    interface{}
		expected []string
	}{
		{
			input:    []interface{}{"val1", "val2", "val3"},
			expected: []string{"val1", "val2", "val3"},
		},
		{
			input:    make([]interface{}, 0),
			expected: make([]string, 0),
		},
		{
			input:    nil,
			expected: nil,
		},
	}

	for _, v := range cases {
		actual := helper.ConvertInterfaceToStringArray(v.input)

		if !helper.IsStringSliceEqual(v.expected, actual) {
			t.Fatalf("Expected %q but got %q", v.expected, actual)
		}
	}
}

func TestInterfaceArrayToStringArray(t *testing.T) {

	cases := []struct {
		input    []interface{}
		expected InterfaceArrayToStringArrayResult
	}{
		{
			input: []interface{}{"val1", "val2", "val3"},
			expected: InterfaceArrayToStringArrayResult{
				Array: []string{"val1", "val2", "val3"},
				Ok:    true,
			},
		},
		{
			input: []interface{}{"val1", 1, "val3"},
			expected: InterfaceArrayToStringArrayResult{
				Array: nil,
				Ok:    false,
			},
		},
		{
			input: make([]interface{}, 0),
			expected: InterfaceArrayToStringArrayResult{
				Array: make([]string, 0),
				Ok:    true,
			},
		},
		{
			input: nil,
			expected: InterfaceArrayToStringArrayResult{
				Array: nil,
				Ok:    false,
			},
		},
	}

	for _, v := range cases {
		actual, ok := helper.ConvertInterfaceArrayToStringArray(v.input)

		if !helper.IsStringSliceEqual(v.expected.Array, actual) {
			t.Fatalf("Expected %q but got %q", v.expected.Array, actual)
		}

		if v.expected.Ok != ok {
			t.Fatalf("Expected %t but got %t", v.expected.Ok, ok)
		}
	}
}

type InterfaceArrayToStringArrayResult struct {
	Array []string
	Ok    bool
}

func TestGetMapFromSet(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expected      *map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: schema.NewSet(dummySetFunc, []interface{}{map[string]interface{}{
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
		actual, err := helper.GetMapFromSet(v.input)

		if v.expectSuccess {
			if err == nil {
				if !reflect.DeepEqual(*v.expected, actual) {
					t.Fatalf("%s: Expected %v but got %v", v.name, *v.expected, actual)
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
