package waf

import (
	"reflect"
	"sort"
	"terraform-provider-ec/ec/helper"
	"testing"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

var (
	// For pointer use
	trueBool  = true
	falseBool = false
)

func TestConvertInterfaceToConditions(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]sdkwaf.Condition
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []interface{}{
				map[string]interface{}{
					"target": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"type":  "REQUEST_HEADERS",
							"value": "Host",
						},
					}),
					"op": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"is_case_insensitive": true,
							"is_negated":          true,
							"type":                "EM",
							"values":              []interface{}{"val1", "val2"},
						},
					}),
				},
				map[string]interface{}{
					"target": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"type": "REQUEST_URI",
						},
					}),
					"op": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"is_case_insensitive": false,
							"is_negated":          false,
							"type":                "RX",
							"value":               "someregex",
						},
					}),
				},
			},
			expectedPtr: &[]sdkwaf.Condition{
				{
					Target: sdkwaf.Target{
						Type:  "REQUEST_HEADERS",
						Value: "Host",
					},
					OP: sdkwaf.OP{
						IsCaseInsensitive: &trueBool,
						IsNegated:         &trueBool,
						Type:              "EM",
						Values:            []string{"val1", "val2"},
					},
				},
				{
					Target: sdkwaf.Target{
						Type: "REQUEST_URI",
					},
					OP: sdkwaf.OP{
						IsCaseInsensitive: &falseBool,
						IsNegated:         &falseBool,
						Type:              "RX",
						Value:             "someregex",
					},
				},
			},
			expectSuccess: true,
		},
		{
			name: "Error path - more than one target defined",
			input: []interface{}{
				map[string]interface{}{
					"target": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"type":  "REQUEST_HEADERS",
							"value": "Host",
						},
						// Error here - can't have multiple targets
						map[string]interface{}{
							"type": "REQUEST_URI",
						},
					}),
					"op": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"is_case_insensitive": true,
							"is_negated":          true,
							"type":                "EM",
							"values":              []interface{}{"val1", "val2"},
						},
					}),
				},
			},
			expectedPtr:   nil,
			expectSuccess: false,
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
		actualPtr, err := ConvertInterfaceToConditions(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

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

func TestConvertInterfaceToConditionGroups(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]sdkwaf.ConditionGroup
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"id":   "group1",
					"name": "Group 1",
					"condition": []interface{}{
						map[string]interface{}{
							"target": helper.NewTerraformSet([]interface{}{
								map[string]interface{}{
									"type":  "REQUEST_HEADERS",
									"value": "Host",
								},
							}),
							"op": helper.NewTerraformSet([]interface{}{
								map[string]interface{}{
									"is_case_insensitive": true,
									"is_negated":          true,
									"type":                "EM",
									"values":              []interface{}{"val1", "val2"},
								},
							}),
						},
					},
				},
				map[string]interface{}{
					"id":   "group2",
					"name": "Group 2",
					"condition": []interface{}{
						map[string]interface{}{
							"target": helper.NewTerraformSet([]interface{}{
								map[string]interface{}{
									"type": "REQUEST_URI",
								},
							}),
							"op": helper.NewTerraformSet([]interface{}{
								map[string]interface{}{
									"is_case_insensitive": false,
									"is_negated":          false,
									"type":                "RX",
									"value":               "someregex",
								},
							}),
						},
					},
				},
			}),
			expectedPtr: &[]sdkwaf.ConditionGroup{
				{
					ID:   "group1",
					Name: "Group 1",
					Conditions: []sdkwaf.Condition{
						{
							Target: sdkwaf.Target{
								Type:  "REQUEST_HEADERS",
								Value: "Host",
							},
							OP: sdkwaf.OP{
								IsCaseInsensitive: &trueBool,
								IsNegated:         &trueBool,
								Type:              "EM",
								Values:            []string{"val1", "val2"},
							},
						},
					},
				},
				{
					ID:   "group2",
					Name: "Group 2",
					Conditions: []sdkwaf.Condition{
						{
							Target: sdkwaf.Target{
								Type: "REQUEST_URI",
							},
							OP: sdkwaf.OP{
								IsCaseInsensitive: &falseBool,
								IsNegated:         &falseBool,
								Type:              "RX",
								Value:             "someregex",
							},
						},
					},
				},
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
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {

		actualPtr, err := ConvertInterfaceToConditionGroups(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				// array equality depends on order
				sort.Slice(actual, func(i, j int) bool {
					return actual[i].ID < actual[j].ID
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
