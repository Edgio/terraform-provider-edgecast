// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package waf

import (
	"reflect"
	"terraform-provider-edgecast/edgecast/helper"
	"testing"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

func TestExpandAccessControls(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expected      sdkwaf.AccessControls
		expectSuccess bool
	}{
		{
			name: "Happy path - strings",
			input: helper.NewTerraformSet([]interface{}{map[string]interface{}{
				"accesslist": []interface{}{"val1", "val2", "val3"},
				"blacklist":  []interface{}{"val4", "val5", "val6"},
				"whitelist":  []interface{}{"val7", "val8", "val9"},
			}}),
			expected: sdkwaf.AccessControls{
				Accesslist: []interface{}{"val1", "val2", "val3"},
				Blacklist:  []interface{}{"val4", "val5", "val6"},
				Whitelist:  []interface{}{"val7", "val8", "val9"},
			},
			expectSuccess: true,
		},
		{
			name: "Happy path - integers",
			input: helper.NewTerraformSet([]interface{}{map[string]interface{}{
				"accesslist": []interface{}{1, 2, 3},
				"blacklist":  []interface{}{4, 5, 6},
				"whitelist":  []interface{}{7, 8, 9},
			}}),
			expected: sdkwaf.AccessControls{
				Accesslist: []interface{}{1, 2, 3},
				Blacklist:  []interface{}{4, 5, 6},
				Whitelist:  []interface{}{7, 8, 9},
			},
			expectSuccess: true,
		},
		{
			name:          "Error path - more than one item defined",
			input:         helper.NewTerraformSet([]interface{}{2, 3}),
			expectSuccess: false,
		},
		{
			name:          "Error path - nil input",
			input:         nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actual, err := ExpandAccessControls(v.input)

		if v.expectSuccess {
			if err == nil {
				if !reflect.DeepEqual(v.expected.Accesslist, actual.Accesslist) {
					t.Fatalf("%s: Expected %q but got %q", v.name, v.expected.Accesslist, actual.Accesslist)
				}

				if !reflect.DeepEqual(v.expected.Blacklist, actual.Blacklist) {
					t.Fatalf("%s: Expected %q but got %q", v.name, v.expected.Blacklist, actual.Blacklist)
				}

				if !reflect.DeepEqual(v.expected.Whitelist, actual.Whitelist) {
					t.Fatalf("%s: Expected %q but got %q", v.name, v.expected.Whitelist, actual.Whitelist)
				}
			} else {
				t.Fatalf("%s: Encountered error: %v", v.name, err)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}
func TestFlattenAccessControls(t *testing.T) {

	cases := []struct {
		name          string
		input         sdkwaf.AccessControls
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: sdkwaf.AccessControls{
				Accesslist: []interface{}{"10.10.10.4"},
				Blacklist:  []interface{}{"10.10.10.3"},
				Whitelist:  []interface{}{"10.10.10.2"},
			},
			expected: []map[string]interface{}{
				{
					"accesslist": []string{
						"10.10.10.4",
					},
					"blacklist": []string{
						"10.10.10.3",
					},
					"whitelist": []string{
						"10.10.10.2",
					},
				},
			},
			expectSuccess: true,
		},
		{
			name: "Nil path",
			input: sdkwaf.AccessControls{
				Accesslist: nil,
				Blacklist:  nil,
				Whitelist:  nil,
			},
			expected:      nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {

		actualGroups := FlattenAccessControls(v.input)

		for i, actual := range actualGroups {

			if v.expectSuccess {
				expected := v.expected[i]

				actualAccessList := actual["accesslist"].([]interface{})
				exceptedAccessList := expected["accesslist"].([]string)
				if actualAccessList[0] != exceptedAccessList[0] {
					t.Fatalf("accesslist[%d] does not match with expected", actualAccessList[0])
					return
				}

				actualBlackList := actual["blacklist"].([]interface{})
				exceptedBlackList := expected["blacklist"].([]string)
				if actualBlackList[0] != exceptedBlackList[0] {
					t.Fatalf("blacklist[%d] does not match with expected", actualBlackList[0])
					return
				}
				actualWhiteList := actual["whitelist"].([]interface{})
				exceptedWhiteList := expected["whitelist"].([]string)
				if actualWhiteList[0] != exceptedWhiteList[0] {
					t.Fatalf("Whitelist[%d] does not match with expected", actualWhiteList[0])
					return
				}
			} else {
				actualAccessList := actual["accesslist"].([]interface{})
				if actualAccessList != nil {
					t.Fatalf("accesslist[%d] does not match with expected", actualAccessList[0])
					return
				}
				actualBlackList := actual["blacklist"].([]interface{})

				if actualBlackList != nil {
					t.Fatalf("blacklist[%d] does not match with expected", actualBlackList[0])
					return
				}
				actualWhiteList := actual["whitelist"].([]interface{})

				if actualWhiteList != nil {
					t.Fatalf("Whitelist[%d] does not match with expected", actualWhiteList[0])
					return
				}
			}
		}
	}
}
