// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"reflect"
	"terraform-provider-edgecast/edgecast/helper"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules"
	"github.com/go-test/deep"
)

func TestExpandChainedRules(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]rules.ChainedRule
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
			expectedPtr: &[]rules.ChainedRule{
				{
					Action: rules.Action{
						ID:      "66000001",
						Message: "Invalid user agent - chainedRule 1.",
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
				{
					Action: rules.Action{
						ID:      "66000002",
						Message: "Invalid user agent - chainedRule 2.",
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

		actualPtr, err := expandChainedRules(v.input)

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

func TestExpandSecRule(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *rules.SecRule
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
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
			expectedPtr: &rules.SecRule{
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
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: true,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a *schema.Set",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {

		actualPtr, err := expandSecRule(v.input)

		if v.expectSuccess {
			if err == nil {
				if actualPtr != nil {
					actual := *actualPtr
					expected := *v.expectedPtr

					if !reflect.DeepEqual(actual, expected) {
						// deep.Equal doesn't compare pointer values, so we just use
						// it to generate a human friendly diff
						diff := deep.Equal(actual, expected)
						t.Errorf("Diff: %+v", diff)
						t.Fatalf("%s: Expected %+v but got %+v",
							v.name,
							expected,
							actual,
						)
					}
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
