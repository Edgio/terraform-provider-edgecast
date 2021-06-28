package helper_test

import (
	"terraform-provider-vmp/vmp/helper"
	"testing"
)

func TestInterfaceSliceEqual(t *testing.T) {
	cases := []struct {
		input    InterfaceSliceEqualInput
		expected bool
	}{
		{
			input: InterfaceSliceEqualInput{
				A: []interface{}{"val1", "val2", 1},
				B: []interface{}{"val1", "val2", 1},
			},
			expected: true,
		},
		{
			input: InterfaceSliceEqualInput{
				A: []interface{}{"val1", "val2", 1},
				B: []interface{}{"val1", "val2", "val3"},
			},
			expected: false,
		},
		{
			input: InterfaceSliceEqualInput{
				A: []interface{}{"val1", "val2", 1, 2},
				B: []interface{}{"val1", "val2", 1},
			},
			expected: false,
		},
		{
			input: InterfaceSliceEqualInput{
				A: make([]interface{}, 0),
				B: make([]interface{}, 0),
			},
			expected: true,
		},
	}

	for _, v := range cases {
		actual := helper.InterfaceSliceEqual(v.input.A, v.input.B)

		if v.expected != actual {
			t.Fatalf("Failed for case: %+v. Expected %t but got %t", v.input, v.expected, actual)
		}
	}
}

func TestStringSliceEqual(t *testing.T) {

	cases := []struct {
		input    StringSliceEqualInput
		expected bool
	}{
		{
			input: StringSliceEqualInput{
				A: []string{"val1", "val2"},
				B: []string{"val1", "val2"},
			},
			expected: true,
		},
		{
			input: StringSliceEqualInput{
				A: []string{"val1", "val2"},
				B: []string{"val1", "val2", "val3"},
			},
			expected: false,
		},
		{
			input: StringSliceEqualInput{
				A: []string{"val1", "val2"},
				B: []string{"val1", "val3"},
			},
			expected: false,
		},
		{
			input: StringSliceEqualInput{
				A: make([]string, 0),
				B: make([]string, 0),
			},
			expected: true,
		},
	}

	for _, v := range cases {
		actual := helper.StringSliceEqual(v.input.A, v.input.B)

		if v.expected != actual {
			t.Fatalf("Failed for case: %+v. Expected %t but got %t", v.input, v.expected, actual)
		}
	}
}

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

type InterfaceSliceEqualInput struct {
	A []interface{}
	B []interface{}
}

type StringSliceEqualInput struct {
	A []string
	B []string
}
