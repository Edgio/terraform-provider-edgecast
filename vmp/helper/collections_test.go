package helper_test

import (
	"terraform-provider-vmp/vmp/helper"
	"testing"
)

func TestIsInterfaceSliceEqual(t *testing.T) {
	cases := []struct {
		input    IsInterfaceSliceEqualInput
		expected bool
	}{
		{
			input: IsInterfaceSliceEqualInput{
				A: []interface{}{"val1", "val2", 1},
				B: []interface{}{"val1", "val2", 1},
			},
			expected: true,
		},
		{
			input: IsInterfaceSliceEqualInput{
				A: []interface{}{"val1", "val2", 1},
				B: []interface{}{"val1", "val2", "val3"},
			},
			expected: false,
		},
		{
			input: IsInterfaceSliceEqualInput{
				A: []interface{}{"val1", "val2", 1, 2},
				B: []interface{}{"val1", "val2", 1},
			},
			expected: false,
		},
		{
			input: IsInterfaceSliceEqualInput{
				A: make([]interface{}, 0),
				B: make([]interface{}, 0),
			},
			expected: true,
		},
	}

	for _, v := range cases {
		actual := helper.IsInterfaceSliceEqual(v.input.A, v.input.B)

		if v.expected != actual {
			t.Fatalf("Failed for case: %+v. Expected %t but got %t", v.input, v.expected, actual)
		}
	}
}

func TestIsStringSliceEqual(t *testing.T) {

	cases := []struct {
		input    IsStringSliceEqualInput
		expected bool
	}{
		{
			input: IsStringSliceEqualInput{
				A: []string{"val1", "val2"},
				B: []string{"val1", "val2"},
			},
			expected: true,
		},
		{
			input: IsStringSliceEqualInput{
				A: []string{"val1", "val2"},
				B: []string{"val1", "val2", "val3"},
			},
			expected: false,
		},
		{
			input: IsStringSliceEqualInput{
				A: []string{"val1", "val2"},
				B: []string{"val1", "val3"},
			},
			expected: false,
		},
		{
			input: IsStringSliceEqualInput{
				A: make([]string, 0),
				B: make([]string, 0),
			},
			expected: true,
		},
	}

	for _, v := range cases {
		actual := helper.IsStringSliceEqual(v.input.A, v.input.B)

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

type IsInterfaceSliceEqualInput struct {
	A []interface{}
	B []interface{}
}

type IsStringSliceEqualInput struct {
	A []string
	B []string
}
