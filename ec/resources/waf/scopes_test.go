// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package waf

import (
	"terraform-provider-edgecast/edgecast/helper"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf"
	"github.com/go-test/deep"
)

func TestExpandScopes(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expected      []waf.Scope
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []interface{}{
				map[string]interface{}{
					"name": "Scope 1",
					"host": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"type":                "EM",
							"is_case_insensitive": true,
							"is_negated":          true,
							"values": []interface{}{
								"site1.com/path1",
								"site2.com",
							},
						},
					}),
					"path": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"type":                "GLOB",
							"is_case_insensitive": false,
							"is_negated":          false,
							"value":               "*",
						},
					}),
					"limit": []interface{}{
						map[string]interface{}{
							"id":           "rateruleid1",
							"duration_sec": 60,
							"enf_type":     "DROP_REQUEST",
							"name":         "limit drop request",
						},
						map[string]interface{}{
							"id":           "rateruleid2",
							"duration_sec": 300,
							"enf_type":     "REDIRECT_302",
							"name":         "limit redirect",
							"url":          "https://mysite.com/redirected",
						},
						map[string]interface{}{
							"id":                   "rateruleid3",
							"duration_sec":         30,
							"enf_type":             "CUSTOM_RESPONSE",
							"status":               404,
							"name":                 "limit custom",
							"response_body_base64": "SGVsbG8sIHdvcmxkIQo=",
							"response_headers": map[string]interface{}{
								"header1": "value1",
								"header2": "value2",
							},
						},
					},
					"acl_audit_id": "accessRuleID",
					"acl_audit_action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name":     "access rule audit action",
							"enf_type": "ALERT",
						},
					}),
					"acl_prod_id": "accessRuleID",
					"acl_prod_action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name":     "access rule prod action",
							"enf_type": "ALERT",
						},
					}),
					"profile_audit_id": "managedRuleID",
					"profile_audit_action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name":     "managed rule audit action",
							"enf_type": "ALERT",
						},
					}),
					"profile_prod_id": "managedRuleID",
					"profile_prod_action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name":                 "managed rule prod action",
							"enf_type":             "CUSTOM_RESPONSE",
							"status":               404,
							"response_body_base64": "SGVsbG8sIHdvcmxkIQo=",
							"response_headers": map[string]interface{}{
								"header1": "value1",
								"header2": "value2",
							},
						},
					}),
					"rules_audit_id": "customRuleID",
					"rules_audit_action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name":     "custom rule audit action",
							"enf_type": "ALERT",
						},
					}),
					"rules_prod_id": "customRuleID",
					"rules_prod_action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name":     "custom rule prod action",
							"enf_type": "BLOCK_REQUEST",
						},
					}),
					"bots_prod_id": "botsRuleID",
					"bots_prod_action": helper.NewTerraformSet([]interface{}{
						map[string]interface{}{
							"name":          "bots rule prod action",
							"valid_for_sec": 60,
							"enf_type":      "BROWSER_CHALLENGE",
						},
					}),
				},
			},
			expected: []waf.Scope{
				{
					Name: "Scope 1",
					Host: waf.MatchCondition{
						Type:              "EM",
						IsCaseInsensitive: wrapBoolInPtr(true),
						IsNegated:         wrapBoolInPtr(true),
						Values: wrapStringsInPtr([]string{
							"site1.com/path1",
							"site2.com",
						}),
					},
					Path: waf.MatchCondition{
						Type:              "GLOB",
						IsCaseInsensitive: wrapBoolInPtr(false),
						IsNegated:         wrapBoolInPtr(false),
						Value:             wrapStringInPtr("*"),
					},
					Limits: &[]waf.Limit{
						{
							ID: "rateruleid1",
							Action: waf.LimitAction{
								DurationSec: 60,
								ENFType:     "DROP_REQUEST",
								Name:        "limit drop request",
							},
						},
						{
							ID: "rateruleid2",
							Action: waf.LimitAction{
								DurationSec: 300,
								ENFType:     "REDIRECT_302",
								Name:        "limit redirect",
								URL:         wrapStringInPtr("https://mysite.com/redirected"),
							},
						},
						{
							ID: "rateruleid3",
							Action: waf.LimitAction{
								DurationSec:        30,
								ENFType:            "CUSTOM_RESPONSE",
								Name:               "limit custom",
								Status:             wrapIntInPtr(404),
								ResponseBodyBase64: wrapStringInPtr("SGVsbG8sIHdvcmxkIQo="),
								ResponseHeaders: &map[string]string{
									"header1": "value1",
									"header2": "value2",
								},
							},
						},
					},
					ACLAuditID: wrapStringInPtr("accessRuleID"),
					ACLAuditAction: &waf.AuditAction{
						Name: "access rule audit action",
						Type: "ALERT",
					},
					ACLProdID: wrapStringInPtr("accessRuleID"),
					ACLProdAction: &waf.ProdAction{
						ENFType: "ALERT",
						Name:    "access rule prod action",
					},

					ProfileAuditID: wrapStringInPtr("managedRuleID"),
					ProfileAuditAction: &waf.AuditAction{
						Name: "managed rule audit action",
						Type: "ALERT",
					},
					ProfileProdID: wrapStringInPtr("managedRuleID"),
					ProfileProdAction: &waf.ProdAction{
						Name:               "managed rule prod action",
						ENFType:            "CUSTOM_RESPONSE",
						Status:             wrapIntInPtr(404),
						ResponseBodyBase64: wrapStringInPtr("SGVsbG8sIHdvcmxkIQo="),
						ResponseHeaders: &map[string]string{
							"header1": "value1",
							"header2": "value2",
						},
					},

					RuleAuditID: wrapStringInPtr("customRuleID"),
					RuleAuditAction: &waf.AuditAction{
						Name: "custom rule audit action",
						Type: "ALERT",
					},
					RuleProdID: wrapStringInPtr("customRuleID"),
					RuleProdAction: &waf.ProdAction{
						ENFType: "BLOCK_REQUEST",
						Name:    "custom rule prod action",
					},

					BotsProdID: wrapStringInPtr("botsRuleID"),
					BotsProdAction: &waf.ProdAction{
						ENFType:     "BROWSER_CHALLENGE",
						ValidForSec: wrapIntInPtr(60),
						Name:        "bots rule prod action",
					},
				},
			},
			expectSuccess: true,
		},
		{
			name:          "Nil input",
			input:         nil,
			expected:      []waf.Scope{},
			expectSuccess: false,
		},
		{
			name:          "Empty input",
			input:         make([]interface{}, 0),
			expected:      []waf.Scope{},
			expectSuccess: true,
		},
	}

	for _, v := range cases {
		actual, err := expandScopes(v.input)
		if v.expectSuccess {
			diffs := deep.Equal(v.expected, actual)
			if len(diffs) > 0 {
				t.Fatalf(
					"Case '%s': \n\nDiffs:\n\n%+v",
					v.name,
					diffs)
			}
		} else {
			if err == nil {
				t.Fatalf("Case '%s': Expected error, but none found", v.name)
			}
		}
	}
}

func TestFlattenScopes(t *testing.T) {
	cases := []struct {
		name          string
		input         *waf.Scopes
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{
			name:          "Happy path",
			expectSuccess: true,
			input: &waf.Scopes{
				Scopes: []waf.Scope{
					{
						Name: "Scope 1",
						Host: waf.MatchCondition{
							Type:              "EM",
							IsCaseInsensitive: wrapBoolInPtr(true),
							IsNegated:         wrapBoolInPtr(true),
							Values: wrapStringsInPtr([]string{
								"site1.com/path1",
								"site2.com",
							}),
						},
						Path: waf.MatchCondition{
							Type:              "GLOB",
							IsCaseInsensitive: wrapBoolInPtr(false),
							IsNegated:         wrapBoolInPtr(false),
							Value:             wrapStringInPtr("*"),
						},
						Limits: &[]waf.Limit{
							{
								ID: "rateruleid1",
								Action: waf.LimitAction{
									DurationSec: 60,
									ENFType:     "DROP_REQUEST",
									Name:        "limit drop request",
								},
							},
							{
								ID: "rateruleid2",
								Action: waf.LimitAction{
									DurationSec: 300,
									ENFType:     "REDIRECT_302",
									Name:        "limit redirect",
									URL:         wrapStringInPtr("https://mysite.com/redirected"),
								},
							},
							{
								ID: "rateruleid3",
								Action: waf.LimitAction{
									DurationSec:        30,
									ENFType:            "CUSTOM_RESPONSE",
									Name:               "limit custom",
									Status:             wrapIntInPtr(404),
									ResponseBodyBase64: wrapStringInPtr("SGVsbG8sIHdvcmxkIQo="),
									ResponseHeaders: &map[string]string{
										"header1": "value1",
										"header2": "value2",
									},
								},
							},
						},
						ACLAuditID: wrapStringInPtr("accessRuleID"),
						ACLAuditAction: &waf.AuditAction{
							Name: "access rule audit action",
							Type: "ALERT",
						},
						ACLProdID: wrapStringInPtr("accessRuleID"),
						ACLProdAction: &waf.ProdAction{
							ENFType: "ALERT",
							Name:    "access rule prod action",
						},

						ProfileAuditID: wrapStringInPtr("managedRuleID"),
						ProfileAuditAction: &waf.AuditAction{
							Name: "managed rule audit action",
							Type: "ALERT",
						},
						ProfileProdID: wrapStringInPtr("managedRuleID"),
						ProfileProdAction: &waf.ProdAction{
							Name:               "managed rule prod action",
							ENFType:            "CUSTOM_RESPONSE",
							Status:             wrapIntInPtr(404),
							ResponseBodyBase64: wrapStringInPtr("SGVsbG8sIHdvcmxkIQo="),
							ResponseHeaders: &map[string]string{
								"header1": "value1",
								"header2": "value2",
							},
						},

						RuleAuditID: wrapStringInPtr("customRuleID"),
						RuleAuditAction: &waf.AuditAction{
							Name: "custom rule audit action",
							Type: "ALERT",
						},
						RuleProdID: wrapStringInPtr("customRuleID"),
						RuleProdAction: &waf.ProdAction{
							ENFType: "BLOCK_REQUEST",
							Name:    "custom rule prod action",
						},

						BotsProdID: wrapStringInPtr("botsRuleID"),
						BotsProdAction: &waf.ProdAction{
							ENFType: "BROWSER_CHALLENGE",
							Name:    "bots rule prod action",
						},
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"name": "Scope 1",
					"host": []map[string]interface{}{
						{
							"type":                "EM",
							"is_case_insensitive": true,
							"is_negated":          true,
							"values": []string{
								"site1.com/path1",
								"site2.com",
							},
						},
					},
					"path": []map[string]interface{}{
						{
							"type":                "GLOB",
							"is_case_insensitive": false,
							"is_negated":          false,
							"value":               "*",
						},
					},
					"limit": []map[string]interface{}{
						{
							"id":           "rateruleid1",
							"duration_sec": 60,
							"enf_type":     "DROP_REQUEST",
							"name":         "limit drop request",
						},
						{
							"id":           "rateruleid2",
							"duration_sec": 300,
							"enf_type":     "REDIRECT_302",
							"name":         "limit redirect",
							"url":          "https://mysite.com/redirected",
						},
						{
							"id":                   "rateruleid3",
							"duration_sec":         30,
							"enf_type":             "CUSTOM_RESPONSE",
							"status":               404,
							"name":                 "limit custom",
							"response_body_base64": "SGVsbG8sIHdvcmxkIQo=",
							"response_headers": map[string]string{
								"header1": "value1",
								"header2": "value2",
							},
						},
					},
					"acl_audit_id": "accessRuleID",
					"acl_audit_action": []map[string]interface{}{
						{
							"name":     "access rule audit action",
							"enf_type": "ALERT",
						},
					},
					"acl_prod_id": "accessRuleID",
					"acl_prod_action": []map[string]interface{}{
						{
							"name":     "access rule prod action",
							"enf_type": "ALERT",
						},
					},
					"profile_audit_id": "managedRuleID",
					"profile_audit_action": []map[string]interface{}{
						{
							"name":     "managed rule audit action",
							"enf_type": "ALERT",
						},
					},
					"profile_prod_id": "managedRuleID",
					"profile_prod_action": []map[string]interface{}{
						{
							"name":                 "managed rule prod action",
							"enf_type":             "CUSTOM_RESPONSE",
							"status":               404,
							"response_body_base64": "SGVsbG8sIHdvcmxkIQo=",
							"response_headers": map[string]string{
								"header1": "value1",
								"header2": "value2",
							},
						},
					},
					"rules_audit_id": "customRuleID",
					"rules_audit_action": []map[string]interface{}{
						{
							"name":     "custom rule audit action",
							"enf_type": "ALERT",
						},
					},
					"rules_prod_id": "customRuleID",
					"rules_prod_action": []map[string]interface{}{
						{
							"name":     "custom rule prod action",
							"enf_type": "BLOCK_REQUEST",
						},
					},
					"bots_prod_id": "botsRuleID",
					"bots_prod_action": []map[string]interface{}{
						{
							"name":     "bots rule prod action",
							"enf_type": "BROWSER_CHALLENGE",
						},
					},
				},
			},
		},
		{
			name:          "Nil input",
			input:         nil,
			expected:      nil,
			expectSuccess: false,
		},
		{
			name:          "Empty input",
			input:         &waf.Scopes{},
			expected:      make([]map[string]interface{}, 0),
			expectSuccess: true,
		},
	}

	for _, v := range cases {
		actual, err := flattenScopes(v.input)
		if v.expectSuccess {
			diffs := deep.Equal(v.expected, actual)
			if len(diffs) > 0 {
				t.Fatalf(
					"Case '%s': \n\nDiffs:\n\n%+v",
					v.name,
					diffs)
			}
		} else {
			if err == nil {
				t.Fatalf("Case '%s': Expected error, but none found", v.name)
			}
		}
	}
}

func wrapStringInPtr(val string) *string {
	return &val
}
func wrapBoolInPtr(val bool) *bool {
	return &val
}
func wrapIntInPtr(val int) *int {
	return &val
}
func wrapStringsInPtr(val []string) *[]string {
	return &val
}
