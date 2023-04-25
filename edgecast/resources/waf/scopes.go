// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"context"
	"errors"
	"log"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/scopes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceScopes() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceScopesCreate,
		ReadContext:   ResourceScopesRead,
		UpdateContext: ResourceScopesUpdate,
		DeleteContext: ResourceScopesDelete,
		Importer:      helper.Import(ResourceScopesRead, "account_number", "id"),

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies your account. Find your account number in the upper right-hand corner of the MCC.",
			},
			"scope": {
				Type:     schema.TypeList,
				Required: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "name",
							Description: "Indicates the name assigned to the Security Application Manager configuration.  \n" +
								"**Default Value:** `name`",
						},
						"host": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes a hostname match condition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_case_insensitive": {
										Type:     schema.TypeBool,
										Optional: true,
										Description: "Indicates whether the comparison between the requested hostname and the `values` argument is case-sensitive. Valid values are: \n\n" +
											"        True | False",
									},
									"is_negated": {
										Type:     schema.TypeBool,
										Optional: true,
										Description: "Indicates whether this match condition will be satisfied when the requested hostname matches or does not match the value defined by the `value`/`values` argument. Valid values are: \n\n" +
											"        True | False",
									},
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description: "Indicates how the system will interpret the comparison between the request's hostname and the value defined within the `value`/`values` argument. Valid values are: \n" +
											" * `EM` - Indicates that request's hostname must be an exact match to one of the case-sensitive values specified in the `values` argument. \n" +
											" * `GLOB` - Indicates that the request's hostname must be an exact match to the wildcard pattern defined in the `value` argument. \n" +
											" * `RX` - Indicates that the request's hostname must be an exact match to the regular expression defined in the `value` argument. \n\n" +
											"    ->Apply this Security Application Manager configuration across all hostnames by setting this argument to `GLOB` and setting the `value` argument to `*`. This type of configuration is also known as `Default`.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**host.type=GLOB or RX:** Identifies a value that will be used to identify requests that are eligible for this Security Application Manager configuration.",
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "**host.type=EM:** Identifies one or more values used to identify requests that are eligible for this Security Application Manager configuration.",
									},
								},
							},
						},
						"limit": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Identifies the set of rate rules that will be enforced for this Security Application Manager configuration and the enforcement action that will be applied to rate limited requests.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Indicates the system-defined ID for the rate rule that will be applied to this Security Application Manager configuration.",
									},
									"duration_sec": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description: "Indicates the length of time, in seconds, that the action defined within this object will be applied to a client that violates the rate rule identified by the `id` argument. Valid values are: \n\n" +
											"        10 | 60 | 300",
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description: "Indicates the type of action that will be applied to rate limited requests. Valid values are: \n" +
											" * `ALERT` - Alert only \n" +
											" * `REDIRECT_302` - Redirect (HTTP 302) \n" +
											" * `CUSTOM_RESPONSE` - Custom response \n" +
											" * `DROP_REQUEST` - Drop request (503 Service Unavailable response with a retry-after of 10 seconds)",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "limit action",
										Description: "Indicates the name assigned to this enforcement action.",
									},
									"response_body_base64": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**limit.enf_type=CUSTOM_RESPONSE:** Indicates the response body that will be sent to rate limited requests. This value is Base64 encoded.",
									},
									"response_headers": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "**limit.enf_type=CUSTOM_RESPONSE:** Contains the set of headers that will be included in the response sent to rate limited requests. Set each desired response header as an argument.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "**limit.enf_type=CUSTOM_RESPONSE:** Indicates the HTTP status code (e.g., 404) for the custom response sent to rate limited requests.",
									},
									"url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**limit.enf_type=REDIRECT_302:** Indicates the URL to which rate limited requests will be redirected.",
									},
								},
							},
						},
						"path": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes a URL path match condition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_case_insensitive": {
										Type:     schema.TypeBool,
										Optional: true,
										Description: "**path.type=EM:** Indicates whether the comparison between the requested URL and the `values` argument is case-sensitive. Valid values are: \n\n" +
											"        True | False",
									},
									"is_negated": {
										Type:     schema.TypeBool,
										Optional: true,
										Description: "Indicates whether this match condition will be satisfied when the requested URL matches or does not match the value defined by the `value`/`values` argument. Valid values are: \n\n" +
											"        True | False",
									},
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description: "Indicates how the system will interpret the comparison between the request's URL and the value defined within the `value`/`values` argument. Valid values are: \n" +
											" * `EM` - Indicates that request's URL path must be an exact match to one of the case-sensitive values specified in the `values` argument.\n" +
											" * `GLOB` - Indicates that the request's URL path must be an exact match to the wildcard pattern defined in the `value` argument. \n" +
											" * `RX` - Indicates that the request's URL path must be an exact match to the regular expression defined in the `value` argument. \n\n" +
											"    ->Apply this Security Application Manager configuration across all URLs by setting this argument to `GLOB` and setting the `value` argument to `*`. This type of configuration is also known as `Default`.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**path.type=GLOB|RX:** Identifies a value that will be used to identify requests that are eligible for this Security Application Manager configuration. Specify a URL path pattern that starts directly after the hostname.",
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "**path.type=EM:** Identifies one or more values used to identify requests that are eligible for this Security Application Manager configuration. Specify a URL path pattern that starts directly after the hostname.",
									},
								},
							},
						},
						"acl_audit_action": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes the type of action that will take place when the access rule defined within the `acl_audit_id` argument is violated. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "Alert Only",
										Description: "Indicates the name assigned to this enforcement action configuration.",
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Set to `ALERT`. This indicates that malicious traffic will be audited.",
									},
								},
							},
						},
						"acl_audit_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the system-defined ID for the access rule that will audit production traffic for this Security Application Manager configuration.",
						},
						"acl_prod_action": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes the type of action that will take place when the access rule defined within the `acl_prod_id` argument is violated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:         schema.TypeInt,
										Optional:     true,
										Description:  "Reserved for future use.",
										ValidateFunc: validation.IntAtLeast(0),
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description: "Indicates the enforcement action that will be applied to malicious traffic. Valid values are: \n" +
											" * `BLOCK_REQUEST` - Block request \n" +
											" * `ALERT` - Alert only \n" +
											" * `REDIRECT_302` - Redirect (HTTP 302) \n" +
											" * `CUSTOM_RESPONSE` - Custom response",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "acl action",
										Description: "Indicates the name assigned to this enforcement action configuration.",
									},
									"response_body_base64": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**acl_prod_action.type=CUSTOM_RESPONSE:** Indicates the response body that will be sent to malicious traffic. This value is Base64 encoded.",
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type:        schema.TypeString,
											Description: "**acl_prod_action.type=CUSTOM_RESPONSE:** Indicates the set of response headers that will be sent to malicious traffic. Each response header is specified as a name/value pair. ",
										},
									},
									"status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: "**acl_prod_action.type=CUSTOM_RESPONSE:** Indicates the HTTP status code (e.g., 404) for the custom response that will be sent to malicious traffic.",
									},
									"url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**acl_prod_action.type=REDIRECT_302:** Indicates the URL to which malicious requests will be redirected.",
									},
								},
							},
						},
						"acl_prod_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the system-defined ID for the access rule that will be applied to production traffic for this Security Application Manager configuration.",
						},
						"bots_prod_action": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes the browser challenge that will be applied to requests that satisfy the bot rule set defined within the `bot_prod_id` argument.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Indicates the number of minutes for which our CDN will serve content to a client that solves a browser challenge without requiring an additional browser challenge to be solved. Specify a value between 1 and 1,440 minutes.",
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Set this property to `BROWSER_CHALLENGE`.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "bots action",
										Description: "Indicates the name assigned to this enforcement action configuration.",
									},
									"response_body_base64": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Reserved for future use.",
									},
									"response_headers": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Reserved for future use.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Indicates the HTTP status code (e.g., 404) for the response provided to clients that are being served the browser challenge.",
									},
									"url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Reserved for future use.",
									},
								},
							},
						},
						"bots_prod_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the system-defined ID for the bots rule that will be applied to production traffic for this Security Application Manager configuration.",
						},
						"profile_audit_action": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes the type of action that will take place when the managed rule defined within the `profile_audit_id` property is violated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "Alert Only",
										Description: "Indicates the name assigned to this enforcement action configuration.",
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Set to `ALERT`. This indicates that malicious traffic will be audited.",
									},
								},
							},
						},
						"profile_audit_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the system-defined ID for the managed rule that will audit production traffic for this Security Application Manager configuration.",
						},
						"profile_prod_action": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes the type of action that will take place when the managed rule defined within the `profile_prod_id` property is violated. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Reserved for future use.",
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description: "Indicates the enforcement action that will be applied to malicious traffic. Valid values are: \n" +
											" * `BLOCK_REQUEST` - Block Request \n" +
											" * `ALERT` - Alert Only \n" +
											" * `REDIRECT_302` - Redirect (HTTP 302) \n" +
											" * `CUSTOM_RESPONSE` - Custom Response",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "profile action",
										Description: "Indicates the name assigned to this enforcement action configuration.",
									},
									"response_body_base64": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the response body that will be sent to malicious traffic. This value is Base64 encoded.",
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the set of response headers that will be sent to malicious traffic. \n\n" +
											"    ->Each response header is specified as a name/value pair.",
									},
									"status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the HTTP status code (e.g., 404) for the custom response that will be sent to malicious traffic.",
									},
									"url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the URL to which malicious requests will be redirected.",
									},
								},
							},
						},
						"profile_prod_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the system-defined ID for the managed rule that will be applied to production traffic for this Security Application Manager configuration.",
						},
						"rules_audit_action": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes the type of action that will take place when the custom rule set defined within the `rules_audit_id` property is violated. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "Alert Only",
										Description: "Indicates the name assigned to this enforcement action configuration.",
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Set to `ALERT`. This indicates that malicious traffic will be audited.",
									},
								},
							},
						},
						"rules_audit_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the system-defined ID for the custom rule set that will audit production traffic for this Security Application Manager configuration.",
						},
						"rules_prod_action": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Description: "Describes the type of action that will take place when the custom rule set defined within the `rules_prod_id` property is violated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Reserved for future use.",
									},
									"enf_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description: "Indicates the enforcement action that will be applied to malicious traffic. Valid values are: \n" +
											" * `BLOCK_REQUEST` - Block Request \n" +
											" * `ALERT` - Alert Only \n" +
											" * `REDIRECT_302` - Redirect (HTTP 302) \n" +
											" * `CUSTOM_RESPONSE` - Custom Response",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "rules action",
										Description: "Indicates the name assigned to this enforcement action configuration.",
									},
									"response_body_base64": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the response body that will be sent to malicious traffic. This value is Base64 encoded.",
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the set of response headers that will be sent to malicious traffic. \n\n" +
											"    ->Each response header is specified as a name/value pair.",
									},
									"status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the HTTP status code (e.g., 404) for the custom response that will be sent to malicious traffic.",
									},
									"url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "**enf_type: CUSTOM_RESPONSE Only:** Indicates the URL to which malicious requests will be redirected.",
									},
								},
							},
						},
						"rules_prod_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the system-defined ID for the custom rule set that will be applied to production traffic for this Security Application Manager configuration.",
						},
					},
				},
			},
		},
	}
}

func ResourceScopesCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	scopes, err := readScopes(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	err = modifyAllScopes(ctx, m, accountNumber, scopes)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// Use account number as the entity ID since a customer can only have one
	// set of scopes
	d.SetId(accountNumber)
	return diag.Diagnostics{}
}

func ResourceScopesRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(internal.ProviderConfig)
	wafService, err := buildWAFService(config)

	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("account_number").(string)
	log.Printf("[INFO] Getting WAF Scopes for Account >> %s", accountNumber)
	resp, err := wafService.Scopes.GetAllScopes(scopes.GetAllScopesParams{
		AccountNumber: accountNumber,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully retrieved WAF Scopes: %+v", resp)
	flattenedScopes, err := flattenScopes(resp)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("account_number", accountNumber)
	d.Set("scope", flattenedScopes)

	// Use account number as the entity ID since a customer can only have one
	// set of scopes
	d.SetId(accountNumber)
	return diag.Diagnostics{}
}

func ResourceScopesUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	scopes, err := readScopes(d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = modifyAllScopes(ctx, m, accountNumber, scopes)
	if err != nil {
		return diag.FromErr(err)
	}
	// Use account number as the entity ID since a customer can only have one
	// set of scopes
	d.SetId(accountNumber)
	return diag.Diagnostics{}
}

func ResourceScopesDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	scopes := make([]scopes.Scope, 0)
	err := modifyAllScopes(ctx, m, accountNumber, scopes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}

// expandScopes converts the values read from a Terraform Configuration
// file into the Scope API Model
func expandScopes(flatScopes interface{}) ([]scopes.Scope, error) {
	if flatScopes == nil {
		return make([]scopes.Scope, 0), errors.New("input was nil")
	}
	if scps, ok := flatScopes.([]interface{}); ok {
		expandedScopes := make([]scopes.Scope, len(scps))
		for i, v := range scps {
			if m, ok := v.(map[string]interface{}); ok {
				scope := scopes.Scope{}
				scope.Name = helper.ConvertToString(
					m["name"])
				scope.Host = expandMatchCondition(
					m["host"])
				scope.Path = expandMatchCondition(
					m["path"])
				scope.ACLAuditAction = expandAuditAction(
					m["acl_audit_action"])
				scope.ACLAuditID = helper.ConvertToStringPointer(
					m["acl_audit_id"],
					true)
				scope.ACLProdAction = expandProdAction(
					m["acl_prod_action"])
				scope.ACLProdID = helper.ConvertToStringPointer(
					m["acl_prod_id"],
					true)
				scope.BotsProdAction = expandProdAction(
					m["bots_prod_action"])
				scope.BotsProdID = helper.ConvertToStringPointer(
					m["bots_prod_id"],
					true)
				scope.ProfileAuditAction = expandAuditAction(
					m["profile_audit_action"])
				scope.ProfileAuditID = helper.ConvertToStringPointer(
					m["profile_audit_id"],
					true)
				scope.ProfileProdAction = expandProdAction(
					m["profile_prod_action"])
				scope.ProfileProdID = helper.ConvertToStringPointer(
					m["profile_prod_id"],
					true)
				scope.RuleAuditAction = expandAuditAction(
					m["rules_audit_action"])
				scope.RuleAuditID = helper.ConvertToStringPointer(
					m["rules_audit_id"],
					true)
				scope.RuleProdAction = expandProdAction(
					m["rules_prod_action"])
				scope.RuleProdID = helper.ConvertToStringPointer(
					m["rules_prod_id"],
					true)
				limits, err := expandLimits(m["limit"])
				if err != nil {
					return nil, err
				}
				scope.Limits = limits
				expandedScopes[i] = scope

			} else {
				return nil, errors.New("scope was not a map[string]interface{}")
			}
		}
		return expandedScopes, nil
	}
	return make([]scopes.Scope, 0), errors.New("input was not a []interface{}")
}

// expandAuditAction converts the values read from a Terraform Configuration
// file into the AuditAction API Model
func expandAuditAction(v interface{}) *scopes.AuditAction {
	if v == nil {
		return nil
	}
	m, _ := helper.ConvertSingletonSetToMap(v)
	if len(m) == 0 {
		return nil
	}
	return &scopes.AuditAction{
		Name: helper.ConvertToString(m["name"]),
		Type: helper.ConvertToString(m["enf_type"]),
	}
}

// expandProdAction converts the values read from a Terraform Configuration
// file into the ProdAction API Model
func expandProdAction(v interface{}) *scopes.ProdAction {
	if v == nil {
		return nil
	}
	m, _ := helper.ConvertSingletonSetToMap(v)
	if len(m) == 0 {
		return nil
	}
	return &scopes.ProdAction{
		Name:    helper.ConvertToString(m["name"]),
		ENFType: helper.ConvertToString(m["enf_type"]),
		ResponseBodyBase64: helper.ConvertToStringPointer(
			m["response_body_base64"],
			true),
		ResponseHeaders: helper.ConvertToStringMapPointer(
			m["response_headers"],
			true),
		Status: helper.ConvertToIntPointer(
			m["status"],
			true),
		URL: helper.ConvertToStringPointer(
			m["url"],
			true),
		ValidForSec: helper.ConvertToIntPointer(
			m["valid_for_sec"],
			true),
	}
}

// expandMatchCondition converts the values read from a Terraform Configuration
// file into the MatchCondition API Model
func expandMatchCondition(v interface{}) scopes.MatchCondition {
	m, _ := helper.ConvertSingletonSetToMap(v)
	mc := scopes.MatchCondition{
		Type: helper.ConvertToString(m["type"]),
	}
	if v, ok := m["is_case_insensitive"]; ok {
		mc.IsCaseInsensitive = helper.ConvertToBoolPointer(v)
	}
	if v, ok := m["is_negated"]; ok {
		mc.IsNegated = helper.ConvertToBoolPointer(v)
	}
	if v, ok := m["value"]; ok {
		mc.Value = helper.ConvertToStringPointer(v, true)
	}
	if v, ok := m["values"]; ok {
		mc.Values = helper.ConvertToStringsPointer(v, true)
	}
	return mc
}

// expandLimits converts the values read from a Terraform Configuration file
// into the Limit API Model
func expandLimits(flatLimits interface{}) (*[]scopes.Limit, error) {
	if flatLimits == nil {
		return nil, nil
	}
	if list, ok := flatLimits.([]interface{}); ok {
		if len(list) == 0 {
			return nil, nil
		}
		limits := make([]scopes.Limit, len(list))
		for i, listItem := range list {
			m := listItem.(map[string]interface{})
			limits[i] = scopes.Limit{
				ID: helper.ConvertToString(m["id"]),
				Action: scopes.LimitAction{
					DurationSec: helper.ConvertToInt(
						m["duration_sec"]),
					ENFType: helper.ConvertToString(
						m["enf_type"]),
					Name: helper.ConvertToString(
						m["name"]),
					ResponseBodyBase64: helper.ConvertToStringPointer(
						m["response_body_base64"],
						true),
					ResponseHeaders: helper.ConvertToStringMapPointer(
						m["response_headers"],
						true),
					Status: helper.ConvertToIntPointer(
						m["status"],
						true),
					URL: helper.ConvertToStringPointer(
						m["url"],
						true),
				},
			}
		}
		return &limits, nil

	}
	return nil, errors.New("flatLimits was not a []interface{}")
}

// flattenScopes converts the Scopes API Model
// into a format that Terraform can work with
func flattenScopes(scopes *scopes.Scopes) ([]map[string]interface{}, error) {
	if scopes == nil {
		return nil, errors.New("scopes was nil")
	}
	flattenedScopes := make([]map[string]interface{}, len(scopes.Scopes))
	for i, s := range scopes.Scopes {
		m := make(map[string]interface{})
		m["name"] = s.Name
		m["host"] = flattenMatchCondition(s.Host)
		m["path"] = flattenMatchCondition(s.Path)
		if s.Limits != nil {
			m["limit"] = flattenLimits(*s.Limits)
		}
		if s.ACLAuditID != nil {
			m["acl_audit_id"] = *s.ACLAuditID
		}
		if (s.ACLAuditAction) != nil {
			m["acl_audit_action"] = flattenAuditAction(*s.ACLAuditAction)
		}
		if s.ACLProdID != nil {
			m["acl_prod_id"] = *s.ACLProdID
		}
		if s.ACLProdAction != nil {
			m["acl_prod_action"] = flattenProdAction(*s.ACLProdAction)
		}
		if s.BotsProdID != nil {
			m["bots_prod_id"] = *s.BotsProdID
		}
		if s.BotsProdAction != nil {
			m["bots_prod_action"] = flattenProdAction(*s.BotsProdAction)
		}
		if s.ProfileAuditID != nil {
			m["profile_audit_id"] = *s.ProfileAuditID
		}
		if s.ProfileAuditAction != nil {
			m["profile_audit_action"] =
				flattenAuditAction(*s.ProfileAuditAction)
		}
		if s.ProfileProdID != nil {
			m["profile_prod_id"] = *s.ProfileProdID
		}
		if s.ProfileProdAction != nil {
			m["profile_prod_action"] = flattenProdAction(*s.ProfileProdAction)
		}
		if s.RuleAuditID != nil {
			m["rules_audit_id"] = *s.RuleAuditID
		}
		if s.RuleAuditAction != nil {
			m["rules_audit_action"] = flattenAuditAction(*s.RuleAuditAction)
		}
		if s.RuleProdID != nil {
			m["rules_prod_id"] = *s.RuleProdID
		}
		if s.RuleProdAction != nil {
			m["rules_prod_action"] = flattenProdAction(*s.RuleProdAction)
		}
		flattenedScopes[i] = m
	}
	return flattenedScopes, nil
}

// flattenProdAction converts the ProdAction API Model
// into a format that Terraform can work with
func flattenProdAction(prodAction scopes.ProdAction) []map[string]interface{} {
	m := make(map[string]interface{})
	m["enf_type"] = prodAction.ENFType
	m["name"] = prodAction.Name
	if prodAction.ValidForSec != nil {
		m["valid_for_sec"] = *prodAction.ValidForSec
	}
	if prodAction.ResponseBodyBase64 != nil {
		m["response_body_base64"] = *prodAction.ResponseBodyBase64
	}
	if prodAction.ResponseHeaders != nil {
		m["response_headers"] = *prodAction.ResponseHeaders
	}
	if prodAction.Status != nil {
		m["status"] = *prodAction.Status
	}
	if prodAction.URL != nil {
		m["url"] = *prodAction.URL
	}
	// We return a collection of just 1 item
	// Since we defined ProdActions as 1-item sets in the schema
	return []map[string]interface{}{m}
}

// flattenAuditAction converts the AuditAction API Model
// into a format that Terraform can work with
func flattenAuditAction(auditAction scopes.AuditAction) []map[string]interface{} {
	m := make(map[string]interface{})
	m["enf_type"] = auditAction.Type
	m["name"] = auditAction.Name
	// We return a collection of just 1 item
	// Since we defined AuditActions as 1-item sets in the schema
	return []map[string]interface{}{m}
}

// flattenMatchCondition converts the MatchCondition API Model
// into a format that Terraform can work with
func flattenMatchCondition(
	matchCondition scopes.MatchCondition,
) []map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = matchCondition.Type
	if matchCondition.IsCaseInsensitive != nil {
		m["is_case_insensitive"] = *matchCondition.IsCaseInsensitive
	}
	if matchCondition.IsNegated != nil {
		m["is_negated"] = *matchCondition.IsNegated
	}
	if matchCondition.Value != nil {
		m["value"] = *matchCondition.Value
	}
	if matchCondition.Values != nil {
		m["values"] = *matchCondition.Values
	}
	// We return a collection of just 1 item
	// Since we defined Host and Path as 1-item sets in the schema
	return []map[string]interface{}{m}
}

// flattenLimits converts the Limit API Model
// into a format that Terraform can work with
func flattenLimits(limits []scopes.Limit) []map[string]interface{} {
	maps := make([]map[string]interface{}, len(limits))
	for i, l := range limits {
		m := make(map[string]interface{})
		m["id"] = l.ID
		m["duration_sec"] = l.Action.DurationSec
		m["enf_type"] = l.Action.ENFType
		m["name"] = l.Action.Name
		if l.Action.ResponseBodyBase64 != nil {
			m["response_body_base64"] = *l.Action.ResponseBodyBase64
		}
		if l.Action.ResponseHeaders != nil {
			m["response_headers"] = *l.Action.ResponseHeaders
		}
		if l.Action.Status != nil {
			m["status"] = *l.Action.Status
		}
		if l.Action.URL != nil {
			m["url"] = *l.Action.URL
		}
		maps[i] = m
	}
	return maps
}

func modifyAllScopes(
	ctx context.Context,
	m interface{},
	accountNumber string,
	scps []scopes.Scope,
) error {
	log.Printf("[INFO] Modifying WAF Scopes for Account >> %s", accountNumber)
	payload := scopes.Scopes{
		CustomerID: accountNumber,
		Scopes:     scps,
	}
	logScopes(payload)
	config := m.(internal.ProviderConfig)
	wafService, err := buildWAFService(config)

	if err != nil {
		return err
	}

	resp, err := wafService.Scopes.ModifyAllScopes(payload)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Successfully modified WAF Scopes: %+v", resp)
	return nil
}

func readScopes(d *schema.ResourceData) ([]scopes.Scope, error) {
	if flatScopes, ok := d.GetOk("scope"); ok {
		expandedScopes, err := expandScopes(flatScopes)
		if err != nil {
			return nil, err
		}
		return expandedScopes, nil
	}
	return nil, errors.New("scopes not found or incorrectly formatted")
}

func logScopes(s scopes.Scopes) {
	log.Printf("[DEBUG] Customer ID: %+v \n", s.CustomerID)
	log.Printf("[DEBUG] Name: %+v \n", s.Name)
	log.Printf("[DEBUG] Scopes: (%d items) \n\n", len(s.Scopes))
	for _, sc := range s.Scopes {
		logScope(sc)
	}
}

func logScope(s scopes.Scope) {
	log.Printf("[DEBUG] \tID: %+v \n", s.ID)
	log.Printf("[DEBUG] \tName: %+v \n", s.Name)
	log.Print("[DEBUG] \tHost: \n\n")
	logMatchCondition(s.Host)
	log.Print("[DEBUG] \tPath: \n\n")
	logMatchCondition(s.Path)
	log.Print("[DEBUG] \tLimits: \n\n")
	if s.Limits != nil {
		for _, l := range *s.Limits {
			logLimit(l)
		}
	}

	if s.ACLAuditID != nil {
		log.Printf("[DEBUG] \tACL Audit ID: %+v \n", *s.ACLAuditID)
	}
	if s.ACLAuditAction != nil {
		log.Print("[DEBUG] \tACL Audit Action: \n\n")
		logAuditAction(*s.ACLAuditAction)
	}
	if s.ACLProdID != nil {
		log.Printf("[DEBUG] \tACL Prod ID: %+v \n", *s.ACLProdID)
	}
	if s.ACLProdAction != nil {
		log.Print("[DEBUG] \tACL Prod Action: \n\n")
		logProdAction(*s.ACLProdAction)
	}

	if s.BotsProdID != nil {
		log.Printf("[DEBUG] \tBots Prod ID: %+v \n", *s.BotsProdID)
	}
	if s.BotsProdAction != nil {
		log.Print("[DEBUG] \tBots Action: \n\n")
		logProdAction(*s.BotsProdAction)
	}

	if s.ProfileAuditID != nil {
		log.Printf("[DEBUG] \tProfile Audit ID: %+v \n", *s.ProfileAuditID)
	}
	if s.ProfileAuditAction != nil {
		log.Print("[DEBUG] \tProfile Audit Action: \n\n")
		logAuditAction(*s.ProfileAuditAction)
	}
	if s.ProfileProdID != nil {
		log.Printf("[DEBUG] \tProfile Prod ID: %+v \n", *s.ProfileProdID)
	}
	if s.ProfileProdAction != nil {
		log.Print("[DEBUG] \tProfile Prod Action: \n\n")
		logProdAction(*s.ProfileProdAction)
	}

	if s.RuleAuditID != nil {
		log.Printf("[DEBUG] \tRule Audit ID: %+v \n", *s.RuleAuditID)
	}
	if s.RuleAuditAction != nil {
		log.Print("[DEBUG] \tRule Audit Action: \n\n")
		logAuditAction(*s.RuleAuditAction)
	}
	if s.RuleProdID != nil {
		log.Printf("[DEBUG] \tRule Prod ID: %+v \n", *s.RuleProdID)
	}
	if s.RuleProdAction != nil {
		log.Print("[DEBUG] \tRule Prod Action: \n\n")
		logProdAction(*s.RuleProdAction)
	}
}

func logAuditAction(a scopes.AuditAction) {
	log.Printf("[DEBUG] \t\t ID: %+v \n", a.ID)
	log.Printf("[DEBUG] \t\t Name: %+v \n", a.Name)
	log.Printf("[DEBUG] \t\t Type: %+v \n", a.Type)
}

func logProdAction(a scopes.ProdAction) {
	log.Printf("[DEBUG] \t\t ID: %+v \n", a.ID)
	log.Printf("[DEBUG] \t\t Name: %+v \n", a.Name)
	log.Printf("[DEBUG] \t\t ENFType: %+v \n", a.ENFType)
	if a.ResponseBodyBase64 != nil {
		log.Printf(
			"[DEBUG] \t\t ResponseBodyBase64: %+v \n",
			*a.ResponseBodyBase64)
	}
	if a.ResponseHeaders != nil {
		log.Printf("[DEBUG] \t\t ResponseHeaders: %+v \n", *a.ResponseHeaders)
	}
	if a.Status != nil {
		log.Printf("[DEBUG] \t\t Status: %+v \n", *a.Status)
	}
	if a.URL != nil {
		log.Printf("[DEBUG] \t\t URL: %+v \n", *a.URL)
	}
	if a.ValidForSec != nil {
		log.Printf("[DEBUG] \t\t ValidForSec: %+v \n", *a.ValidForSec)
	}
}

func logLimit(l scopes.Limit) {
	log.Printf("[DEBUG] \t\t ID: %+v \n", l.ID)
	log.Printf("[DEBUG] \t\t DurationSec: %+v \n", l.Action.DurationSec)
	log.Printf("[DEBUG] \t\t ENFType: %+v \n", l.Action.ENFType)
	log.Printf("[DEBUG] \t\t Name: %+v \n", l.Action.Name)
	if l.Action.ResponseBodyBase64 != nil {
		log.Printf(
			"[DEBUG] \t\t ResponseBodyBase64: %+v \n",
			*l.Action.ResponseBodyBase64)
	}
	if l.Action.ResponseHeaders != nil {
		log.Printf(
			"[DEBUG] \t\t ResponseHeaders: %+v \n",
			*l.Action.ResponseHeaders)
	}
	if l.Action.Status != nil {
		log.Printf("[DEBUG] \t\t Status: %+v \n", *l.Action.Status)
	}
	if l.Action.URL != nil {
		log.Printf("[DEBUG] \t\t URL: %+v \n", *l.Action.URL)
	}
}

func logMatchCondition(m scopes.MatchCondition) {
	if m.IsCaseInsensitive != nil {
		log.Printf(
			"[DEBUG] \t\t IsCaseInsensitive: %+v \n",
			*m.IsCaseInsensitive)
	}
	if m.IsNegated != nil {
		log.Printf("[DEBUG] \t\t IsNegated: %+v \n", *m.IsNegated)
	}
	if m.Value != nil {
		log.Printf("[DEBUG] \t\t Value: %+v \n", *m.Value)
	}
	if m.Values != nil {
		log.Printf("[DEBUG] \t\t Values: %+v \n", *m.Values)
	}
	log.Printf("[DEBUG] \t\t Type: %+v \n", m.Type)
}
