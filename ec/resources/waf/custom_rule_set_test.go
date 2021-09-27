package waf

import (
	"reflect"
	"terraform-provider-ec/ec/helper"
	"testing"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

func TestExpandDirectives(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]sdkwaf.Directive
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"sec_rule": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name": "REQUEST_HEADERS",
							"action": helper.NewTerraformSet([]interface{}{
								map[string]interface{}{
									"id":              "66000000",
									"msg":             "Invalid user agent.",
									"transformations": []interface{}{"NONE"},
								},
							}),
							"operator": helper.NewTerraformSet([]interface{}{
								map[string]interface{}{
									"is_negated": false,
									"type":       "CONTAINS",
									"value":      "bot",
								},
							}),
							"variable": []interface{}{
								map[string]interface{}{
									"is_count": false,
									"type":     "REQUEST_HEADERS",
									"match": []interface{}{
										map[string]interface{}{
											"is_negated": false,
											"is_regex":   false,
											"value":      "User-Agent",
										},
									},
								},
							},
							"chained_rule": []interface{}{
								map[string]interface{}{
									"action": helper.NewTerraformSet([]interface{}{
										map[string]interface{}{
											"id":              "66000001",
											"msg":             "Invalid user agent - chained.",
											"transformations": []interface{}{"NONE"},
										},
									}),
									"operator": helper.NewTerraformSet([]interface{}{
										map[string]interface{}{
											"is_negated": false,
											"type":       "CONTAINS",
											"value":      "bot",
										},
									}),
									"variable": []interface{}{
										map[string]interface{}{
											"is_count": false,
											"type":     "REQUEST_HEADERS",
											"match": []interface{}{
												map[string]interface{}{
													"is_negated": false,
													"is_regex":   false,
													"value":      "User-Agent",
												},
											},
										},
									},
								},
							},
						},
					}),
				},
			}),
			expectedPtr: &[]sdkwaf.Directive{
				{
					SecRule: sdkwaf.SecRule{
						Name: "REQUEST_HEADERS",
						Action: sdkwaf.Action{
							ID:              "66000000",
							Message:         "Invalid user agent.",
							Transformations: []string{"NONE"},
						},
						Operator: sdkwaf.Operator{
							IsNegated: false,
							Type:      "CONTAINS",
							Value:     "bot",
						},
						Variables: []sdkwaf.Variable{
							{
								IsCount: false,
								Type:    "REQUEST_HEADERS",
								Matches: []sdkwaf.Match{
									{
										IsRegex:   false,
										IsNegated: false,
										Value:     "User-Agent",
									},
								},
							},
						},
						ChainedRules: []sdkwaf.ChainedRule{
							{
								Action: sdkwaf.Action{
									ID:              "66000001",
									Message:         "Invalid user agent - chained.",
									Transformations: []string{"NONE"},
								},
								Operator: sdkwaf.Operator{
									IsNegated: false,
									Type:      "CONTAINS",
									Value:     "bot",
								},
								Variables: []sdkwaf.Variable{
									{
										IsCount: false,
										Type:    "REQUEST_HEADERS",
										Matches: []sdkwaf.Match{
											{
												IsRegex:   false,
												IsNegated: false,
												Value:     "User-Agent",
											},
										},
									},
								},
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

		actualPtr, err := ExpandDirectives(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				if !reflect.DeepEqual(expected, actual) {
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

func TestExpandChainedRules(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]sdkwaf.ChainedRule
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []interface{}{
				map[string]interface{}{
					"action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"id":              "66000001",
							"msg":             "Invalid user agent - chainedRule 1.",
							"transformations": []interface{}{"NONE"},
						},
					}),
					"operator": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"is_negated": false,
							"type":       "CONTAINS",
							"value":      "bot",
						},
					}),
					"variable": []interface{}{
						map[string]interface{}{
							"is_count": false,
							"type":     "REQUEST_HEADERS",
							"match": []interface{}{
								map[string]interface{}{
									"is_negated": false,
									"is_regex":   false,
									"value":      "User-Agent",
								},
							},
						},
					},
				},
				map[string]interface{}{
					"action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"id":              "66000002",
							"msg":             "Invalid user agent - chainedRule 2.",
							"transformations": []interface{}{"NONE"},
						},
					}),
					"operator": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"is_negated": false,
							"type":       "CONTAINS",
							"value":      "bot",
						},
					}),
					"variable": []interface{}{
						map[string]interface{}{
							"is_count": false,
							"type":     "REQUEST_HEADERS",
							"match": []interface{}{
								map[string]interface{}{
									"is_negated": false,
									"is_regex":   false,
									"value":      "User-Agent",
								},
							},
						},
					},
				},
			},
			expectedPtr: &[]sdkwaf.ChainedRule{
				{
					Action: sdkwaf.Action{
						ID:              "66000001",
						Message:         "Invalid user agent - chainedRule 1.",
						Transformations: []string{"NONE"},
					},
					Operator: sdkwaf.Operator{
						IsNegated: false,
						Type:      "CONTAINS",
						Value:     "bot",
					},
					Variables: []sdkwaf.Variable{
						{
							IsCount: false,
							Type:    "REQUEST_HEADERS",
							Matches: []sdkwaf.Match{
								{
									IsRegex:   false,
									IsNegated: false,
									Value:     "User-Agent",
								},
							},
						},
					},
				},
				{
					Action: sdkwaf.Action{
						ID:              "66000002",
						Message:         "Invalid user agent - chainedRule 2.",
						Transformations: []string{"NONE"},
					},
					Operator: sdkwaf.Operator{
						IsNegated: false,
						Type:      "CONTAINS",
						Value:     "bot",
					},
					Variables: []sdkwaf.Variable{
						{
							IsCount: false,
							Type:    "REQUEST_HEADERS",
							Matches: []sdkwaf.Match{
								{
									IsRegex:   false,
									IsNegated: false,
									Value:     "User-Agent",
								},
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

		actualPtr, err := ExpandChainedRules(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				if !reflect.DeepEqual(expected, actual) {
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
