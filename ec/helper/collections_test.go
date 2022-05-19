package helper_test

import (
	"reflect"
	"terraform-provider-edgecast/ec/helper"
	"testing"
)

func TestIsInterfaceArray(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    1,
			expected: false,
		},
		{
			input:    []string{"val1", "val2"},
			expected: false,
		},
		{
			input:    []int{1, 2},
			expected: false,
		},
		{
			input:    []interface{}{"val1", "val2", 1},
			expected: true,
		},
		{
			input:    make([]interface{}, 0),
			expected: true,
		},
	}

	for _, v := range cases {
		actual := helper.IsInterfaceArray(v.input)

		if v.expected != actual {
			t.Fatalf("Failed for case: %+v. Expected %t but got %t", v.input, v.expected, actual)
		}
	}
}

func TestConvertSliceToStrings(t *testing.T) {
	cases := []struct {
		name     string
		input    []interface{}
		expected ConvertSliceToStringsResult
	}{
		{
			name:  "Happy path - all strings",
			input: []interface{}{"val1", "val2", "val3"},
			expected: ConvertSliceToStringsResult{
				Array: []string{"val1", "val2", "val3"},
				Ok:    true,
			},
		},
		{
			name:  "Error path - one value is not an int",
			input: []interface{}{"val1", 1, "val3"},
			expected: ConvertSliceToStringsResult{
				Array: nil,
				Ok:    false,
			},
		},
		{
			name:  "Empty list",
			input: make([]interface{}, 0),
			expected: ConvertSliceToStringsResult{
				Array: make([]string, 0),
				Ok:    true,
			},
		},
		{
			name:  "Nil input",
			input: nil,
			expected: ConvertSliceToStringsResult{
				Array: make([]string, 0),
				Ok:    false,
			},
		},
	}
	for _, v := range cases {
		actual, err := helper.ConvertSliceToStrings(v.input)
		if v.expected.Ok {
			if err == nil {
				if !reflect.DeepEqual(v.expected.Array, actual) {
					t.Fatalf(
						"Case '%s': Expected %q but got %q",
						v.name,
						v.expected.Array,
						actual)
				}
			} else {
				t.Fatalf(
					"Case '%s': Encountered error when none were expected",
					v.name)
			}
		} else {
			if err == nil {
				t.Fatalf(
					"Case '%s': No error when one was expected",
					v.name)
			}
		}
	}
}

type ConvertSliceToStringsResult struct {
	Array []string
	Ok    bool
}
