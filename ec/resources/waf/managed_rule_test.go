package waf

import (
	"reflect"
	"sort"
	"terraform-provider-ec/ec/helper"
	"testing"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

func TestExpandDisabledRules(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]sdkwaf.DisabledRule
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"policy_id": "myPolicyId1",
					"rule_id":   "myRuleId1",
				},
				map[string]interface{}{
					"policy_id": "myPolicyId2",
					"rule_id":   "myRuleId2",
				},
			}),
			expectedPtr: &[]sdkwaf.DisabledRule{
				{
					PolicyID: "myPolicyId1",
					RuleID:   "myRuleId1",
				},
				{
					PolicyID: "myPolicyId2",
					RuleID:   "myRuleId2",
				},
			},
			expectSuccess: true,
		},
		{
			name:          "Happy path - None Defined",
			input:         helper.NewTerraformSet([]interface{}{}),
			expectedPtr:   &[]sdkwaf.DisabledRule{},
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ExpandDisabledRules(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				// array equality depends on order, so we'll just sort by Policy ID
				sort.Slice(actual, func(i, j int) bool {
					return actual[i].PolicyID < actual[j].PolicyID
				})

				if !reflect.DeepEqual(expected, actual) {
					t.Fatalf("%s: Expected %+v but got %+v", v.name, expected, actual)
				}
			} else {
				t.Fatalf("%s: Encountered error where one was not expected: %+v", v.name, err)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}

func TestExpandRuleTargetUpdates(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]sdkwaf.RuleTargetUpdate
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"is_negated":     true,
					"is_regex":       true,
					"replace_target": "replacementTarget1",
					"rule_id":        "ruleId1",
					"target":         "target1",
					"target_match":   "targetMatch1",
				},
				map[string]interface{}{
					"is_negated":     false,
					"is_regex":       false,
					"replace_target": "replacementTarget2",
					"rule_id":        "ruleId2",
					"target":         "target2",
					"target_match":   "targetMatch2",
				},
			}),
			expectedPtr: &[]sdkwaf.RuleTargetUpdate{
				{
					IsNegated:     true,
					IsRegex:       true,
					ReplaceTarget: "replacementTarget1",
					RuleID:        "ruleId1",
					Target:        "target1",
					TargetMatch:   "targetMatch1",
				},
				{
					IsNegated:     false,
					IsRegex:       false,
					ReplaceTarget: "replacementTarget2",
					RuleID:        "ruleId2",
					Target:        "target2",
					TargetMatch:   "targetMatch2",
				},
			},
			expectSuccess: true,
		},
		{
			name:          "Happy path - None Defined",
			input:         helper.NewTerraformSet([]interface{}{}),
			expectedPtr:   &[]sdkwaf.RuleTargetUpdate{},
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ExpandRuleTargetUpdates(v.input)

		if v.expectSuccess {
			if err != nil {
				t.Fatalf("%s: Encountered error where one was not expected: %+v", v.name, err)
			}

			actual := *actualPtr
			expected := *v.expectedPtr

			// array equality depends on order, so we'll just sort by Rule ID
			sort.Slice(actual, func(i, j int) bool {
				return actual[i].RuleID < actual[j].RuleID
			})

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("%s: Expected %+v but got %+v", v.name, expected, actual)
			}

		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}

func TestConvertInterfaceToGeneralSettings(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *sdkwaf.GeneralSettings
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"anomaly_threshold":      2,
					"arg_length":             1024,
					"arg_name_length":        512,
					"combined_file_sizes":    2048,
					"ignore_cookie":          []interface{}{"cookie1, cookie2"},
					"ignore_header":          []interface{}{"Content-Type", "User-Agent"},
					"ignore_query_args":      []interface{}{"arg1", "arg2"},
					"json_parser":            true,
					"max_num_args":           100,
					"paranoia_level":         3,
					"process_request_body":   false,
					"response_header_name":   "X-Request-Blocked",
					"total_arg_length":       1000,
					"validate_utf8_encoding": true,
					"xml_parser":             false,
				},
			}),
			expectedPtr: &sdkwaf.GeneralSettings{
				AnomalyThreshold:     2,
				ArgLength:            1024,
				ArgNameLength:        512,
				CombinedFileSizes:    2048,
				IgnoreCookie:         []string{"cookie1, cookie2"},
				IgnoreHeader:         []string{"Content-Type", "User-Agent"},
				IgnoreQueryArgs:      []string{"arg1", "arg2"},
				JsonParser:           true,
				MaxNumArgs:           100,
				ParanoiaLevel:        3,
				ProcessRequestBody:   false,
				ResponseHeaderName:   "X-Request-Blocked",
				TotalArgLength:       1000,
				ValidateUtf8Encoding: true,
				XmlParser:            false,
			},
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - GetMapFromSet() error handling",
			input:         helper.NewTerraformSet([]interface{}{}),
			expectedPtr:   &sdkwaf.GeneralSettings{},
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ConvertInterfaceToGeneralSettings(v.input)

		if v.expectSuccess {
			if err != nil {
				t.Fatalf("%s: Encountered error where one was not expected: %+v", v.name, err)
			}

			actual := *actualPtr
			expected := *v.expectedPtr

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("%s: Expected %+v but got %+v", v.name, expected, actual)
			}

		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}
