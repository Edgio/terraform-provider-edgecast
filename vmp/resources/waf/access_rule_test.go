package waf_test

import (
	"math"
	"terraform-provider-vmp/vmp/api"
	"terraform-provider-vmp/vmp/helper"
	"terraform-provider-vmp/vmp/resources/waf"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestConvertInterfaceToAccessControls(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expected      api.AccessControls
		expectSuccess bool
	}{
		{
			name: "Happy path - strings",
			input: schema.NewSet(dummySetFunc, []interface{}{map[string]interface{}{
				"accesslist": []interface{}{"val1", "val2", "val3"},
				"blacklist":  []interface{}{"val4", "val5", "val6"},
				"whitelist":  []interface{}{"val7", "val8", "val9"},
			}}),
			expected: api.AccessControls{
				Accesslist: []interface{}{"val1", "val2", "val3"},
				Blacklist:  []interface{}{"val4", "val5", "val6"},
				Whitelist:  []interface{}{"val7", "val8", "val9"},
			},
			expectSuccess: true,
		},
		{
			name: "Happy path - integers",
			input: schema.NewSet(dummySetFunc, []interface{}{map[string]interface{}{
				"accesslist": []interface{}{1, 2, 3},
				"blacklist":  []interface{}{4, 5, 6},
				"whitelist":  []interface{}{7, 8, 9},
			}}),
			expected: api.AccessControls{
				Accesslist: []interface{}{1, 2, 3},
				Blacklist:  []interface{}{4, 5, 6},
				Whitelist:  []interface{}{7, 8, 9},
			},
			expectSuccess: true,
		},
		{
			name:          "Error path - more than one item defined",
			input:         schema.NewSet(dummySetFunc, []interface{}{2, 3}),
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
		actual, err := waf.ConvertInterfaceToAccessControls(v.input)

		if v.expectSuccess {
			if err == nil {
				if !helper.IsInterfaceSliceEqual(v.expected.Accesslist, actual.Accesslist) {
					t.Fatalf("%s: Expected %q but got %q", v.name, v.expected.Accesslist, actual.Accesslist)
				}

				if !helper.IsInterfaceSliceEqual(v.expected.Blacklist, actual.Blacklist) {
					t.Fatalf("%s: Expected %q but got %q", v.name, v.expected.Blacklist, actual.Blacklist)
				}

				if !helper.IsInterfaceSliceEqual(v.expected.Whitelist, actual.Whitelist) {
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

func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
