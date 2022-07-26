// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"reflect"
	"terraform-provider-edgecast/edgecast/helper"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/custom"
	"github.com/go-test/deep"
)

func TestFlattenCustomRuleDirectives(t *testing.T) {
	cases := []struct {
		name     string
		input    []custom.CustomRuleDirective
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []custom.CustomRuleDirective{
				{
					SecRule: rules.SecRule{
						Name: "Test Rule 1",
						Action: rules.Action{
							ID:      "66000000",
							Message: "Invalid user agent.",
							Transformations: []rules.Transformation{
								rules.TransformRemoveNulls,
							},
						},
						Operator: rules.Operator{
							IsNegated: false,
							Type:      rules.OpBeginsWith,
							Value:     "bot",
						},
						Variables: []rules.Variable{
							{
								IsCount: false,
								Type:    rules.VarRequestCookies,
								Matches: []rules.Match{
									{
										IsRegex:   false,
										IsNegated: false,
										Value:     "mycookie",
									},
								},
							},
						},
						ChainedRules: []rules.ChainedRule{
							{
								Action: rules.Action{
									ID:      "66000001",
									Message: "Invalid user agent - chained.",
									Transformations: []rules.Transformation{
										rules.TransformLowerCase,
									},
								},
								Operator: rules.Operator{
									IsNegated: false,
									Type:      rules.OpEndsWith,
									Value:     "bot",
								},
								Variables: []rules.Variable{
									{
										IsCount: false,
										Type:    rules.VarRequestHeaders,
										Matches: []rules.Match{
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
			expected: []map[string]interface{}{
				{
					"sec_rule": []map[string]interface{}{
						{
							"name": "Test Rule 1",
							"action": []map[string]interface{}{
								{
									"id":              "66000000",
									"msg":             "Invalid user agent.",
									"transformations": []string{"REMOVENULLS"},
								},
							},
							"operator": []map[string]interface{}{
								{
									"is_negated": false,
									"type":       "BEGINSWITH",
									"value":      "bot",
								},
							},
							"variable": []map[string]interface{}{
								{
									"is_count": false,
									"type":     "REQUEST_COOKIES",
									"match": []map[string]interface{}{
										{
											"is_negated": false,
											"is_regex":   false,
											"value":      "mycookie",
										},
									},
								},
							},
							"chained_rule": []map[string]interface{}{
								{
									"action": []map[string]interface{}{
										{
											"id":              "66000001",
											"msg":             "Invalid user agent - chained.",
											"transformations": []string{"LOWERCASE"},
										},
									},
									"operator": []map[string]interface{}{
										{
											"is_negated": false,
											"type":       "ENDSWITH",
											"value":      "bot",
										},
									},
									"variable": []map[string]interface{}{
										{
											"is_count": false,
											"type":     "REQUEST_HEADERS",
											"match": []map[string]interface{}{
												{
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
					},
				},
			},
		},
		{
			name:     "Nil path",
			input:    nil,
			expected: make([]map[string]interface{}, 0),
		},
		{
			name:     "Empty path",
			input:    make([]custom.CustomRuleDirective, 0),
			expected: make([]map[string]interface{}, 0),
		},
	}

	for _, c := range cases {
		actual := flattenCustomRuleDirectives(c.input)

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

func TestExpandCustomRuleDirectives(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]custom.CustomRuleDirective
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
			expectedPtr: &[]custom.CustomRuleDirective{
				{
					SecRule: rules.SecRule{
						Name: "REQUEST_HEADERS",
						Action: rules.Action{
							ID:      "66000000",
							Message: "Invalid user agent.",
							Transformations: []rules.Transformation{
								rules.TransformNone,
							},
						},
						Operator: rules.Operator{
							IsNegated: false,
							Type:      rules.OpContains,
							Value:     "bot",
						},
						Variables: []rules.Variable{
							{
								IsCount: false,
								Type:    rules.VarRequestHeaders,
								Matches: []rules.Match{
									{
										IsRegex:   false,
										IsNegated: false,
										Value:     "User-Agent",
									},
								},
							},
						},
						ChainedRules: []rules.ChainedRule{
							{
								Action: rules.Action{
									ID:      "66000001",
									Message: "Invalid user agent - chained.",
									Transformations: []rules.Transformation{
										rules.TransformNone,
									},
								},
								Operator: rules.Operator{
									IsNegated: false,
									Type:      rules.OpContains,
									Value:     "bot",
								},
								Variables: []rules.Variable{
									{
										IsCount: false,
										Type:    rules.VarRequestHeaders,
										Matches: []rules.Match{
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

		actualPtr, err := expandCustomRuleDirectives(v.input)

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
