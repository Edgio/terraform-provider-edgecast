package helper_test

import (
	"terraform-provider-ec/ec/helper"
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
