package waf

import (
	"reflect"
	"sort"
	"terraform-provider-edgecast/ec/helper"
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

func TestExpandGeneralSettings(t *testing.T) {

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
		actualPtr, err := ExpandGeneralSettings(v.input)

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

func TestFlattenDisabledRules(t *testing.T) {

	cases := []struct {
		name     string
		input    []sdkwaf.DisabledRule
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []sdkwaf.DisabledRule{
				{
					PolicyID: "policyID1",
					RuleID:   "ruleID1",
				},
				{
					PolicyID: "policyID2",
					RuleID:   "ruleID2",
				},
			},
			expected: []map[string]interface{}{
				{
					"policy_id": "policyID1",
					"rule_id":   "ruleID1",
				},
				{
					"policy_id": "policyID2",
					"rule_id":   "ruleID2",
				},
			},
		},
		{
			name:     "Empty collection",
			input:    []sdkwaf.DisabledRule{},
			expected: []map[string]interface{}{},
		},
	}

	for _, v := range cases {

		actualRules := FlattenDisabledRules(v.input)

		if !reflect.DeepEqual(actualRules, v.expected) {
			t.Fatalf("%s: Expected %+v but got %+v", v.name, actualRules, v.expected)
		}
	}
}

func TestFlattenRuleTargetUpdates(t *testing.T) {

	cases := []struct {
		name     string
		input    []sdkwaf.RuleTargetUpdate
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []sdkwaf.RuleTargetUpdate{
				{
					IsNegated:     true,
					IsRegex:       true,
					ReplaceTarget: "ARGS",
					RuleID:        "ruleID1",
					Target:        "GEO",
					TargetMatch:   "name",
				},
				{
					IsNegated:     false,
					IsRegex:       false,
					ReplaceTarget: "REQUEST_COOKIES",
					RuleID:        "ruleID2",
					Target:        "REQUEST_HEADERS",
					TargetMatch:   "value",
				},
			},
			expected: []map[string]interface{}{
				{
					"is_negated":     true,
					"is_regex":       true,
					"replace_target": "ARGS",
					"rule_id":        "ruleID1",
					"target":         "GEO",
					"target_match":   "name",
				},
				{
					"is_negated":     false,
					"is_regex":       false,
					"replace_target": "REQUEST_COOKIES",
					"rule_id":        "ruleID2",
					"target":         "REQUEST_HEADERS",
					"target_match":   "value",
				},
			},
		},
		{
			name:     "Empty collection",
			input:    []sdkwaf.RuleTargetUpdate{},
			expected: []map[string]interface{}{},
		},
	}

	for _, v := range cases {

		actualRules := FlattenRuleTargetUpdates(v.input)

		if !reflect.DeepEqual(actualRules, v.expected) {
			t.Fatalf("%s: Expected %+v but got %+v", v.name, actualRules, v.expected)
		}
	}
}

func TestFlattenGeneralSettings(t *testing.T) {

	cases := []struct {
		name     string
		input    sdkwaf.GeneralSettings
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: sdkwaf.GeneralSettings{
				AnomalyThreshold:     10,
				ArgLength:            1024,
				ArgNameLength:        2048,
				CombinedFileSizes:    4092,
				IgnoreCookie:         []string{"iCookie1", "iCookie2"},
				IgnoreHeader:         []string{"iHeader1", "iHeader2"},
				IgnoreQueryArgs:      []string{"iQueryArgs1", "iQueryArgs2"},
				JsonParser:           true,
				MaxNumArgs:           100,
				ParanoiaLevel:        3,
				ProcessRequestBody:   false,
				ResponseHeaderName:   "response_header",
				TotalArgLength:       256,
				ValidateUtf8Encoding: true,
				XmlParser:            false,
			},
			expected: []map[string]interface{}{
				{
					"anomaly_threshold":      10,
					"arg_length":             1024,
					"arg_name_length":        2048,
					"combined_file_sizes":    4092,
					"ignore_cookie":          []string{"iCookie1", "iCookie2"},
					"ignore_header":          []string{"iHeader1", "iHeader2"},
					"ignore_query_args":      []string{"iQueryArgs1", "iQueryArgs2"},
					"json_parser":            true,
					"max_num_args":           100,
					"paranoia_level":         3,
					"process_request_body":   false,
					"response_header_name":   "response_header",
					"total_arg_length":       256,
					"validate_utf8_encoding": true,
					"xml_parser":             false,
				},
			},
		},
		{
			name:  "Empty collection",
			input: sdkwaf.GeneralSettings{},
			expected: []map[string]interface{}{
				{
					"anomaly_threshold":      0,
					"arg_length":             0,
					"arg_name_length":        0,
					"combined_file_sizes":    0,
					"ignore_cookie":          []string(nil),
					"ignore_header":          []string(nil),
					"ignore_query_args":      []string(nil),
					"json_parser":            false,
					"max_num_args":           0,
					"paranoia_level":         0,
					"process_request_body":   false,
					"response_header_name":   "",
					"total_arg_length":       0,
					"validate_utf8_encoding": false,
					"xml_parser":             false,
				},
			},
		},
	}

	for _, v := range cases {

		actualGeneralSettings := FlattenGeneralSettings(v.input)

		actualGeneralSettingsLength := len(actualGeneralSettings)
		if actualGeneralSettingsLength != 1 {
			t.Fatalf("GeneralSettings map of length 1 expected, actual length: %v", actualGeneralSettingsLength)
			return
		}

		if !reflect.DeepEqual(actualGeneralSettings, v.expected) {
			t.Fatalf("%s: Expected %+v but got %+v", v.name, actualGeneralSettings, v.expected)
		}
	}
}
