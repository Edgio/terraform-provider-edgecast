package helper_test

import (
	"terraform-provider-vmp/vmp/helper"
	"testing"
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
