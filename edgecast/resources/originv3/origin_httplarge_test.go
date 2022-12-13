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
		expectedPtr *OriginGroupState
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
				"origin": []any{
					map[string]any{
						"name":             "marketing-origin-entry-a",
						"host":             "https://cdn-la.example.com",
						"port":             443,
						"is_primary":       true,
						"storage_type_id":  1,
						"protocol_type_id": 2,
					},
				},
			},
			expectedPtr: &OriginGroupState{
				Name:          "my_origin_group",
				HostHeader:    "myhost",
				NetworkTypeID: 2,
				ShieldPops:    []*string{&pop1, &pop2},
				TLSSettings: &originv3.TlsSettings{
					SniHostname:        originv3.NewNullableString("mysnihost"),
					AllowSelfSigned:    &isAllowSelfSigned,
					PublicKeysToVerify: []string{"key1", "key2"},
				},
				Origins: []*OriginState{
					{
						ID:             0,
						GroupID:        0,
						Name:           "marketing-origin-entry-a",
						Host:           "https://cdn-la.example.com",
						Port:           443,
						IsPrimary:      true,
						StorageTypeID:  1,
						ProtocolTypeID: 2,
					},
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

			actualGrpPtr, errs := expandHttpLargeOriginGroup(rd)

			if !tt.expectErrs && (len(errs) > 0) {
				t.Fatalf("unexpected errors: %v", errs)
			}

			if tt.expectErrs && (len(errs) != tt.errCount) {
				t.Fatalf("expected %d errors but got %d", tt.errCount, len(errs))
			}

			if tt.expectErrs && (len(errs) > 0) {
				return // successful test for error case
			}

			//Group
			actualGrp := *actualGrpPtr
			expectedGrp := *tt.expectedPtr

			if !reflect.DeepEqual(actualGrp, expectedGrp) {
				// deep.Equal doesn't compare pointer values, so we just use it to
				// generate a human friendly diff
				diff := deep.Equal(actualGrp, expectedGrp)
				t.Errorf("Diff: %+v", diff)
				t.Fatalf("%s: Expected %+v but got %+v",
					tt.name,
					expectedGrp,
					actualGrp,
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

func TestExpandOrigins(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   []*OriginState
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []any{
				map[string]any{
					"id":               1,
					"name":             "marketing-origin-entry-a",
					"host":             "https://cdn-la.example.com",
					"port":             443,
					"is_primary":       true,
					"storage_type_id":  1,
					"protocol_type_id": 2,
				},
				map[string]any{
					"id":               2,
					"name":             "marketing-origin-entry-b",
					"host":             "https://cdn-lb.example.com",
					"port":             443,
					"is_primary":       true,
					"storage_type_id":  1,
					"protocol_type_id": 2,
				},
			},
			expectedPtr: []*OriginState{
				{
					ID:             1,
					GroupID:        0,
					Name:           "marketing-origin-entry-a",
					Host:           "https://cdn-la.example.com",
					Port:           443,
					IsPrimary:      true,
					StorageTypeID:  1,
					ProtocolTypeID: 2,
				},
				{
					ID:             2,
					GroupID:        0,
					Name:           "marketing-origin-entry-b",
					Host:           "https://cdn-lb.example.com",
					Port:           443,
					IsPrimary:      true,
					StorageTypeID:  1,
					ProtocolTypeID: 2,
				},
			},
			expectSuccess: true,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a []interface{}",
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "nil input",
			input:         nil,
			expectedPtr:   make([]*OriginState, 0),
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := expandOrigins(v.input)

		if v.expectSuccess {
			if err == nil {
				actual := actualPtr
				expected := v.expectedPtr

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
					"sni_hostname":          sniHost.Get(),
					"allow_self_signed":     false,
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

func TestFlattenOrigins(t *testing.T) {
	origin1ID := int32(1)
	origin1Name := "origin1"
	host1 := "host1"
	isOrigin1Primary := true
	failoverOrder1 := int32(0)

	origin2ID := int32(2)
	origin2Name := "origin2"
	host2 := "host2"
	isOrigin2Primary := false
	failoverOrder2 := int32(1)

	port := int32(443)
	storogeTypeID := int32(1)
	protocolTypeID := int32(2)

	cases := []struct {
		name          string
		input         []originv3.CustomerOriginFailoverOrder
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{
			name:          "Happy path",
			expectSuccess: true,
			input: []originv3.CustomerOriginFailoverOrder{
				{
					Id:             &origin1ID,
					Name:           &origin1Name,
					Host:           &host1,
					Port:           &port,
					IsPrimary:      &isOrigin1Primary,
					StorageTypeId:  &storogeTypeID,
					ProtocolTypeId: &protocolTypeID,
					FailoverOrder:  &failoverOrder1,
				},
				{
					Id:             &origin2ID,
					Name:           &origin2Name,
					Host:           &host2,
					Port:           &port,
					IsPrimary:      &isOrigin2Primary,
					StorageTypeId:  &storogeTypeID,
					ProtocolTypeId: &protocolTypeID,
					FailoverOrder:  &failoverOrder2,
				},
			},
			expected: []map[string]interface{}{
				{
					"id":               1,
					"name":             "origin1",
					"host":             "host1",
					"port":             443,
					"is_primary":       true,
					"storage_type_id":  1,
					"protocol_type_id": 2,
					"failover_order":   0,
				},
				{
					"id":               2,
					"name":             "origin2",
					"host":             "host2",
					"port":             443,
					"is_primary":       false,
					"storage_type_id":  1,
					"protocol_type_id": 2,
					"failover_order":   1,
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
		actual := flattenOrigins(c.input)

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
