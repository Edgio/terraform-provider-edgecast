// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"reflect"
	"sort"
	"terraform-provider-edgecast/edgecast/helper"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/bot"
	"github.com/go-test/deep"
)

func TestFlattenBotRuleDirectives(t *testing.T) {
	cases := []struct {
		name     string
		input    []bot.BotRuleDirective
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []bot.BotRuleDirective{
				{
					Include: "r3010_ec_bot_challenge_reputation.conf.json",
				},
				{
					SecRule: &rules.SecRule{
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
				},
			},
			expected: []map[string]interface{}{
				{
					"include": "r3010_ec_bot_challenge_reputation.conf.json",
				},
				{
					"sec_rule": []map[string]interface{}{
						{
							"name": "REQUEST_HEADERS",
							"action": []map[string]interface{}{
								{
									"id":              "66000000",
									"msg":             "Invalid user agent.",
									"transformations": []string{"NONE"},
								},
							},
							"operator": []map[string]interface{}{
								{
									"is_negated": false,
									"type":       "CONTAINS",
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
							"chained_rule": []map[string]interface{}{
								{
									"action": []map[string]interface{}{
										{
											"transformations": []string{"NONE"},
										},
									},
									"operator": []map[string]interface{}{
										{
											"is_negated": false,
											"type":       "CONTAINS",
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
			input:    make([]bot.BotRuleDirective, 0),
			expected: make([]map[string]interface{}, 0),
		},
	}

	for _, c := range cases {
		actual := flattenBotRuleDirectives(c.input)

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

func TestExpandBotRuleDirectives(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *[]bot.BotRuleDirective
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

				map[string]interface{}{
					"include": "r3010_ec_bot_challenge_reputation.conf.json",
				},
			}),
			expectedPtr: &[]bot.BotRuleDirective{
				{
					SecRule: &rules.SecRule{
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
				},
				{
					Include: "r3010_ec_bot_challenge_reputation.conf.json",
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

		actualPtr, err := expandBotRuleDirectives(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				/*
					Terraform's schema.Set.List() function claims to be
					deterministic but that seems to be the case for the same
					instance of an array, not across arrays with the same value.
					Therefore, the ordering is inconsistent. Arrays have to be
					in the same order to be deeply equal.

					To account for this, we will sort the two arrays before
					comparing.
				*/

				sort.SliceStable(actual, func(i, j int) bool {
					return actual[i].Include < actual[j].Include
				})

				sort.SliceStable(expected, func(i, j int) bool {
					return expected[i].Include < expected[j].Include
				})

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
