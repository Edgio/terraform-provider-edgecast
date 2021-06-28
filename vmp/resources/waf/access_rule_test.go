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

func TestInterfaceToAccessControls(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected api.AccessControls
	}{
		{
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
		},
		{
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
		},
	}

	for _, v := range cases {
		actual := waf.InterfaceToAccessControls(v.input)

		if !helper.InterfaceSliceEqual(v.expected.Accesslist, actual.Accesslist) {
			t.Fatalf("Expected %q but got %q", v.expected.Accesslist, actual.Accesslist)
		}

		if !helper.InterfaceSliceEqual(v.expected.Blacklist, actual.Blacklist) {
			t.Fatalf("Expected %q but got %q", v.expected.Blacklist, actual.Blacklist)
		}

		if !helper.InterfaceSliceEqual(v.expected.Whitelist, actual.Whitelist) {
			t.Fatalf("Expected %q but got %q", v.expected.Whitelist, actual.Whitelist)
		}
	}
}

func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
