// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package originv3

import (
	"reflect"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/originv3"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandHttpLargeOriginGrp(t *testing.T) {
	t.Parallel()

	isAllowSelfSigned := false
	pop1 := "pop1"
	pop2 := "pop2"

	tests := []struct {
		name        string
		input       map[string]any
		expectedPtr *originv3.CustomerOriginGroupHTTPRequest
		expectErrs  bool
		errCount    int
	}{
		{
			name:       "Happy path",
			expectErrs: false,
			input: map[string]any{
				"name":            "my_origin_group",
				"platform":        "http-large",
				"host_header":     "myhost",
				"network_type_id": 2,
				"shield_pops":     []any{"pop1", "pop2"},
				"tls_settings": []any{
					map[string]any{
						"sni_hostname":          "mysnihost",
						"allow_self_signed":     false,
						"public_keys_to_verify": []any{"key1", "key2"},
					},
				},
			},
			expectedPtr: &originv3.CustomerOriginGroupHTTPRequest{
				Name:          "my_origin_group",
				HostHeader:    originv3.NewNullableString("myhost"),
				NetworkTypeId: originv3.NewNullableInt32(2),
				ShieldPops:    []*string{&pop1, &pop2},
				TlsSettings: &originv3.TlsSettings{
					SniHostname:        originv3.NewNullableString("mysnihost"),
					AllowSelfSigned:    &isAllowSelfSigned,
					PublicKeysToVerify: []string{"key1", "key2"},
				},
			},
		},
		{
			name:        "nil input",
			errCount:    1,
			input:       nil,
			expectedPtr: nil,
			expectErrs:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var rd *schema.ResourceData
			if tt.input != nil {
				rd = schema.TestResourceDataRaw(
					t,
					GetOriginGrpHttpLargeSchema(),
					tt.input)
			}

			actualPtr, errs := expandHttpLargeOriginGroup(rd)

			if !tt.expectErrs && (len(errs) > 0) {
				t.Fatalf("unexpected errors: %v", errs)
			}

			if tt.expectErrs && (len(errs) != tt.errCount) {
				t.Fatalf("expected %d errors but got %d", tt.errCount, len(errs))
			}

			if tt.expectErrs && (len(errs) > 0) {
				return // successful test for error case
			}

			actual := *actualPtr
			expected := *tt.expectedPtr

			if !reflect.DeepEqual(actual, expected) {
				// deep.Equal doesn't compare pointer values, so we just use it to
				// generate a human friendly diff
				diff := deep.Equal(actual, expected)
				t.Errorf("Diff: %+v", diff)
				t.Fatalf("%s: Expected %+v but got %+v",
					tt.name,
					expected,
					actual,
				)
			}
		})
	}
}

func TestExpandTLSSettings(t *testing.T) {
	isAllowSelfSigned := false

	cases := []struct {
		name          string
		input         any
		expectedPtr   *originv3.TlsSettings
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []any{
				map[string]any{
					"sni_hostname":          "mysnihost",
					"allow_self_signed":     false,
					"public_keys_to_verify": []any{"key1", "key2"},
				},
			},
			expectedPtr: &originv3.TlsSettings{
				SniHostname:        originv3.NewNullableString("mysnihost"),
				AllowSelfSigned:    &isAllowSelfSigned,
				PublicKeysToVerify: []string{"key1", "key2"},
			},
			expectSuccess: true,
		},
		{
			name:          "nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a []interface{}",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := expandTLSSettings(v.input)

		if v.expectSuccess {
			if err == nil {
				actual := *actualPtr
				expected := *v.expectedPtr

				if !reflect.DeepEqual(actual, expected) {
					// deep.Equal doesn't compare pointer values, so we just use it to
					// generate a human friendly diff
					diff := deep.Equal(actual, expected)
					t.Errorf("Diff: %+v", diff)
					t.Fatalf("%s: Expected %+v but got %+v",
						v.name,
						expected,
						actual,
					)
				}
			} else {
				t.Fatalf("%s: Encountered error where one was not expected: %+v",
					v.name,
					err,
				)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}

func TestFlattenTLSSettings(t *testing.T) {

	isAllowSelfSigned := false
	sniHost := originv3.NewNullableString("mysnihost")
	cases := []struct {
		name          string
		input         *originv3.TlsSettings
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{
			name:          "Happy path",
			expectSuccess: true,
			input: &originv3.TlsSettings{
				SniHostname:        sniHost,
				AllowSelfSigned:    &isAllowSelfSigned,
				PublicKeysToVerify: []string{"key1", "key2"},
			},
			expected: []map[string]interface{}{
				{
					"sni_hostname":          sniHost,
					"allow_self_signed":     &isAllowSelfSigned,
					"public_keys_to_verify": []string{"key1", "key2"},
				},
			},
		},
		{
			name:          "Nil input",
			input:         nil,
			expected:      make([]map[string]interface{}, 0),
			expectSuccess: false,
		},
	}

	for _, c := range cases {
		actual := flattenTLSSettings(c.input)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}
