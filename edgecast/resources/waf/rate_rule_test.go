// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package waf

import (
	"reflect"
	"sort"
	"terraform-provider-edgecast/edgecast/helper"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/rate"
)

var (
	// For pointer use
	trueBool  = true
	falseBool = false
)

func TestExpandConditions(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]rate.Condition
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
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
			}),
			expectedPtr: &[]rate.Condition{
				{
					Target: rate.Target{
						Type:  "REQUEST_HEADERS",
						Value: "Host",
					},
					OP: rate.OP{
						IsCaseInsensitive: &trueBool,
						IsNegated:         &trueBool,
						Type:              "EM",
						Values:            []string{"val1", "val2"},
					},
				},
				{
					Target: rate.Target{
						Type: "REQUEST_URI",
					},
					OP: rate.OP{
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
			input: helper.NewTerraformSet([]interface{}{
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
			}),
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
		actualPtr, err := expandConditions(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				// array equality depends on order, so we'll just sort by target type
				sort.Slice(actual, func(i, j int) bool {
					return actual[i].Target.Type < actual[j].Target.Type
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

func TestExpandConditionGroups(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]rate.ConditionGroup
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"id":   "group1",
					"name": "Group 1",
					"condition": helper.NewTerraformSet([]interface{}{
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
					}),
				},
				map[string]interface{}{
					"id":   "group2",
					"name": "Group 2",
					"condition": helper.NewTerraformSet([]interface{}{
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
					}),
				},
			}),
			expectedPtr: &[]rate.ConditionGroup{
				{
					ID:   "group1",
					Name: "Group 1",
					Conditions: []rate.Condition{
						{
							Target: rate.Target{
								Type:  "REQUEST_HEADERS",
								Value: "Host",
							},
							OP: rate.OP{
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
					Conditions: []rate.Condition{
						{
							Target: rate.Target{
								Type: "REQUEST_URI",
							},
							OP: rate.OP{
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

		actualPtr, err := expandConditionGroups(v.input)

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

func TestFlattenConditionGroups(t *testing.T) {

	cases := []struct {
		name     string
		input    []rate.ConditionGroup
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []rate.ConditionGroup{
				{
					ID:   "group1",
					Name: "Group 1",
					Conditions: []rate.Condition{
						{
							Target: rate.Target{
								Type:  "REQUEST_HEADERS",
								Value: "Host",
							},
							OP: rate.OP{
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
					Conditions: []rate.Condition{
						{
							Target: rate.Target{
								Type: "REQUEST_URI",
							},
							OP: rate.OP{
								IsCaseInsensitive: &(falseBool),
								IsNegated:         &falseBool,
								Type:              "RX",
								Value:             "someregex",
							},
						},
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"id":   "group1",
					"name": "Group 1",
					"condition": []map[string]interface{}{
						{
							"target": []map[string]interface{}{
								{
									"type":  "REQUEST_HEADERS",
									"value": "Host",
								},
							},
							"op": []map[string]interface{}{
								{
									"type":                "EM",
									"values":              []string{"val1", "val2"},
									"is_case_insensitive": true,
									"is_negated":          true,
								},
							},
						},
					},
				},
				{
					"id":   "group2",
					"name": "Group 2",
					"condition": []map[string]interface{}{
						{
							"target": []map[string]interface{}{
								{
									"type": "REQUEST_URI",
								},
							},
							"op": []map[string]interface{}{
								{
									"type":                "RX",
									"value":               "someregex",
									"is_case_insensitive": false,
									"is_negated":          false,
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "Empty collection",
			input:    []rate.ConditionGroup{},
			expected: []map[string]interface{}{},
		},
	}

	for _, v := range cases {

		actualGroups := flattenConditionGroups(v.input)

		for i, actual := range actualGroups {

			expected := v.expected[i]

			actualName, ok := actual["name"].(string)

			if !ok || actualName != expected["name"].(string) {
				t.Fatalf("group[%d].name does not match expected name", i)
				return
			}

			id, ok := actual["id"].(string)

			if !ok || id != expected["id"].(string) {
				t.Fatalf("group[%d].id does not match expected ID", i)
				return
			}

			actualConditions, ok := actual["condition"].([]map[string]interface{})

			if !ok {
				t.Fatalf("expected group[%d].conditions to be []interface", i)
				return
			}

			expectedConditions := expected["condition"].([]map[string]interface{})

			checkConditions(actualConditions, expectedConditions, t)
		}
	}
}

func checkConditions(conditions []map[string]interface{}, expectedConditions []map[string]interface{}, t *testing.T) {

	for j, actual := range conditions {

		expected := expectedConditions[j]

		opSetActual := actual["op"].([]map[string]interface{})
		opSetExpected := expected["op"].([]map[string]interface{})

		// there should only be one
		opActual := opSetActual[0]
		opExpected := opSetExpected[0]

		if expectedOpType, doCheck := opExpected["type"].(string); doCheck {

			if opActual["type"].(string) != expectedOpType {
				t.Fatalf("conditions[%d].op.type does not match expected", j)
				return
			}
		}

		if expectedOpValue, doCheck := opExpected["value"].(string); doCheck {

			if opActual["value"].(string) != expectedOpValue {
				t.Fatalf("conditions[%d].op.value does not match expected", j)
				return
			}
		}

		if expectedOpValues, doCheck := opExpected["values"].([]string); doCheck {

			if !reflect.DeepEqual(opActual["values"].([]string), expectedOpValues) {
				t.Fatalf("conditions[%d].op.values does not match expected", j)
				return
			}
		}

		if expectedOpIsNegated, doCheck := opExpected["is_negated"].(bool); doCheck {

			if opActual["is_negated"].(bool) != expectedOpIsNegated {
				t.Fatalf("conditions[%d].op.is_negated does not match expected", j)
				return
			}
		}

		if expectedOpIsCaseInsensitive, doCheck := opExpected["is_case_insensitive"].(bool); doCheck {

			if opActual["is_case_insensitive"].(bool) != expectedOpIsCaseInsensitive {
				t.Fatalf("conditions[%d].op.is_case_insensitive does not match expected", j)
				return
			}
		}

		targetSetActual := actual["target"].([]map[string]interface{})
		targetSetExpected := expected["target"].([]map[string]interface{})

		// there should only be one
		targetActual := targetSetActual[0]
		targetExpected := targetSetExpected[0]

		if expectedTargetType, doCheck := targetExpected["type"].(string); doCheck {

			if targetActual["type"].(string) != expectedTargetType {
				t.Fatalf("conditions[%d].target.type does not match expected", j)
				return
			}
		}

		if expectedTargetValue, doCheck := targetExpected["value"].(string); doCheck {

			if targetActual["value"].(string) != expectedTargetValue {
				t.Fatalf("conditions[%d].target.value does not match expected", j)
				return
			}
		}
	}
}
