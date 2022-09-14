// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper_test

import (
	"reflect"
	"testing"

	"terraform-provider-edgecast/edgecast/helper"
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
				Array: nil,
				Ok:    true,
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

func TestGetStringFromMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
		argMap      map[string]any
		argKey      string
		want        string
	}{
		{
			name:        "Happy Path",
			expectError: false,
			argMap: map[string]any{
				"k1": "v1",
			},
			argKey: "k1",
			want:   "v1",
		},
		{
			name:        "Error - value not in map",
			expectError: true,
			argMap: map[string]any{
				"k1": "v1",
			},
			argKey: "k2",
			want:   "",
		},
		{
			name:        "Error - value is not a string",
			expectError: true,
			argMap: map[string]any{
				"k1": 1,
			},
			argKey: "k1",
			want:   "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := helper.GetStringFromMap(tt.argMap, tt.argKey)

			if tt.expectError && ok {
				t.Fatal("expected error but got not none")
			}

			if !tt.expectError && !ok {
				t.Fatal("unexpected error")
			}

			if tt.expectError && !ok {
				return // successful test, got expected error
			}

			if got != tt.want {
				t.Logf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestGetBoolFromMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
		argMap      map[string]any
		argKey      string
		want        bool
	}{
		{
			name:        "Happy Path",
			expectError: false,
			argMap: map[string]any{
				"k1": true,
			},
			argKey: "k1",
			want:   true,
		},
		{
			name:        "Error - value not in map",
			expectError: true,
			argMap: map[string]any{
				"k1": "v1",
			},
			argKey: "k2",
		},
		{
			name:        "Error - value is not a bool",
			expectError: true,
			argMap: map[string]any{
				"k1": 1,
			},
			argKey: "k1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := helper.GetBoolFromMap(tt.argMap, tt.argKey)

			if tt.expectError && ok {
				t.Fatal("expected error but got not none")
			}

			if !tt.expectError && !ok {
				t.Fatal("unexpected error")
			}

			if tt.expectError && !ok {
				return // successful test, got expected error
			}

			// expect success, check results
			if got != tt.want {
				t.Logf("got %t, want %t", got, tt.want)
			}
		})
	}
}
