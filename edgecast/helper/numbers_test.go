// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper_test

import (
	"terraform-provider-edgecast/edgecast/helper"
	"testing"
)

func TestParseInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		arg         string
		want        int64
		expectError bool
	}{
		{
			name:        "Happy Path - valid int64",
			arg:         "9223372036854775807",
			want:        9223372036854775807,
			expectError: false,
		},
		{
			name:        "not a valid int64",
			arg:         "not an int64",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := helper.ParseInt64(tt.arg)

			if tt.expectError && err == nil {
				t.Fatal("expected error but got not none")
			}

			if !tt.expectError && err != nil {
				t.Fatal("unexpected error")
			}

			if tt.expectError && err != nil {
				return // successful test, got expected error
			}

			if got != tt.want {
				t.Fatalf("expected %d but got %d", got, tt.want)
			}
		})
	}
}
