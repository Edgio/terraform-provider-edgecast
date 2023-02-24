// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package internal_test

import (
	"testing"

	"terraform-provider-edgecast/edgecast/internal"

	"github.com/hashicorp/go-cty/cty"
)

func TestValidateDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		arg         any
		expectError bool
	}{
		{
			name:        "Happy Path - string",
			arg:         "6h",
			expectError: false,
		},
		{
			name:        "Happy Path - int64",
			arg:         int64(10000),
			expectError: false,
		},
		{
			name:        "Happy Path - float64",
			arg:         float64(90000),
			expectError: false,
		},
		{
			name:        "Error - invalid type",
			arg:         []string{"not a valid duration"},
			expectError: true,
		},
		{
			name:        "Error - invalid string duration format",
			arg:         "8fsadfs",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := internal.ValidateDuration(tt.arg, cty.Path{})

			if tt.expectError && !got.HasError() {
				t.Fatal("expected error but got not none")
			}

			if !tt.expectError && got.HasError() {
				t.Fatal("unexpected error")
			}

			if tt.expectError && got.HasError() {
				return // successful test, got expected error
			}

			if !tt.expectError && !got.HasError() {
				return // successful test, got no error and none expected
			}
		})
	}
}
