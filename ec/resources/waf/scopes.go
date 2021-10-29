package waf

import (
	"context"
	"errors"
	"log"
	"terraform-provider-ec/ec/api"
	"terraform-provider-ec/ec/helper"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf"
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
		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Identifies your account by its customer " +
					"account number.",
			},
			"scope": {
				Type:     schema.TypeList,
				Required: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the name assigned to the " +
								"Security Application Manager configuration. " +
								"Default Value: 'name'",
						},
						"host": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes a hostname match " +
								"condition. Refer to the URL and Path " +
								"section for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_case_insensitive": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"is_negated": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"limit": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Indicates the " +
											"system-defined ID for the rate " +
											"limit configuration that will " +
											"be applied to this Security " +
											"Application Manager " +
											"configuration.",
									},
									"duration_sec": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description: "Indicates the length " +
											"of time, in seconds, that the " +
											"action defined within this " +
											"object will be applied to a " +
											"client that violates the rate " +
											"rule identified by the id " +
											"property.\\\n\\\nValid values " +
											"are: 10 | 60 | 300",
									},
									"enf_type": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Indicates the type " +
											"of action that will be applied " +
											"to rate limited requests." +
											"\\\n\\\nValid values are:" +
											"ALERT: Alert Only" +
											"REDIRECT_302: Redirect (HTTP 302)" +
											"CUSTOM_RESPONSE: Custom Response" +
											"DROP_REQUEST: Drop Request " +
											"(503 Service Unavailable " +
											"response with a retry-after of " +
											"10 seconds)",
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Indicates the name " +
											"assigned to this enforcement " +
											"action.",
									},
									"response_body_base64": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Note: Only valid when" +
											" ENFType is set to " +
											"CUSTOM_RESPONSE \\\n\\\n" +
											"Indicates the response body that" +
											" will be sent to rate limited " +
											"requests. This value is Base64 " +
											"encoded.",
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Note: Only valid " +
											"when ENFType is set to " +
											"CUSTOM_RESPONSE\\\n\\\n" +
											"Contains the set of headers " +
											"that will be included in " +
											"the response sent to rate " +
											"limited requests.",
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
										Description: "Note: Only valid when " +
											"ENFType is set to " +
											"CUSTOM_RESPONSE\\\n\\\nIndicates" +
											" the HTTP status code " +
											"(e.g., 404) for the custom " +
											"response sent to rate limited " +
											"requests.",
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Note: Only valid when " +
											"ENFType is set to REDIRECT_302" +
											"\\\n\\\nIndicates the URL to " +
											"which rate limited requests " +
											"will be redirected.",
									},
								},
							},
						},
						"path": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes a URL match condition." +
								"Refer to the URL and Path section for " +
								"property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_case_insensitive": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"is_negated": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"acl_audit_action": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describe the type of action that " +
								"will take place when the access rule " +
								"defined within the acl_audit_id property " +
								"is violated. Refer to the Audit Action " +
								"section for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"acl_audit_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the system-defined ID " +
								"for the access rule that will audit " +
								"production traffic for this Security " +
								"Application Manager configuration.",
						},
						"acl_prod_action": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes the type of action " +
								"that will take place when the access rule " +
								"defined within the acl_prod_id property is " +
								"violated. Refer to the Prod Action " +
								"section for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"enf_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_body_base64": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"acl_prod_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the system-defined " +
								"ID for the access rule that will be " +
								"applied to production traffic for this " +
								"Security Application Manager configuration.",
						},
						"bots_prod_action": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes the type of action " +
								"that will take place when the bots rule " +
								"defined within the bots_prod_id property " +
								"is violated. Refer to the Prod Action " +
								"section for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"enf_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_body_base64": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"bots_prod_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the system-defined " +
								"ID for the bots rule that will be applied " +
								"to production traffic for this " +
								"Security Application Manager configuration.",
						},
						"profile_audit_action": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes the type of action " +
								"that will take place when the managed " +
								"rule defined within the profile_audit_id " +
								"property is violated. Refer to the " +
								"Audit Action  section for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"profile_audit_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the system-defined ID " +
								"for the managed rule that will audit " +
								"production traffic for this Security " +
								"Application Manager configuration.",
						},
						"profile_prod_action": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes the type of action " +
								"that will take place when the managed " +
								"rule defined within the profile_prod_id " +
								"property is violated. Refer to the Prod " +
								"Action section for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:     schema.TypeInt,
										Required: true,

										ValidateFunc: validation.IntAtLeast(0),
									},
									"enf_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_body_base64": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,

										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"profile_prod_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the system-defined " +
								"ID for the managed rule that will be applied" +
								" to production traffic for this Security " +
								"Application Manager configuration.",
						},
						"rule_audit_action": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes the type of action that " +
								"will take place when the custom rule set " +
								"defined within the rule_audit_id property " +
								"is violated. Refer to the Audit Action " +
								"section for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"rule_audit_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the system-defined ID " +
								"for the custom rule set that will audit " +
								"production traffic for this Security " +
								"Application Manager configuration.",
						},
						"rule_prod_action": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Description: "Describes the type of action that " +
								"will take place when the custom rule set " +
								"defined within the rule_prod_id property is " +
								"violated. Refer to the Prod Action section " +
								"for property details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"valid_for_sec": {
										Type:     schema.TypeInt,
										Required: true,

										ValidateFunc: validation.IntAtLeast(0),
									},
									"enf_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_body_base64": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"response_headers": {
										Type:     schema.TypeMap,
										Optional: true,

										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"rule_prod_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Indicates the system-defined ID " +
								"for the custom rule set that will be applied" +
								" to production traffic for this Security " +
								"Application Manager configuration.",
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
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	log.Printf("[INFO] Creating WAF Scopes for Account >> %s", accountNumber)
	scopes := waf.Scopes{
		CustomerID: accountNumber,
	}
	if flatScopes, ok := d.GetOk("scope"); ok {
		expandedScopes, err := ExpandScopes(flatScopes)
		if err != nil {
			return diag.FromErr(err)
		}
		scopes.Scopes = expandedScopes
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading Scopes",
			Detail:   "Scopes not found or incorrectly formatted",
		})
		return diags
	}

	logScopes(scopes)
	config := m.(**api.ClientConfig)
	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := wafService.ModifyAllScopes(scopes)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created WAF Scopes: %+v", resp)
	d.SetId(resp.ID)
	return diags
}

func ResourceScopesRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(**api.ClientConfig)
	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("account_number").(string)
	log.Printf("[INFO] Getting WAF Scopes for Account >> %s", accountNumber)
	resp, err := wafService.GetAllScopes(accountNumber)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully retrieved WAF Scopes: %+v", resp)
	flattenedScopes, err := FlattenScopes(resp)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("account_number", accountNumber)
	d.Set("scope", flattenedScopes)
	d.SetId(resp.ID)
	return diags
}

func ResourceScopesUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceScopesDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(**api.ClientConfig)
	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("account_number").(string)
	log.Printf("[INFO] Deleting WAF Scopes for Account >> %s", accountNumber)
	emptyScopes := waf.Scopes{
		CustomerID: accountNumber,
		Scopes:     make([]waf.Scope, 0),
	}
	resp, err := wafService.ModifyAllScopes(emptyScopes)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully deleted WAF Scopes: %+v", resp)
	d.SetId("")
	return diags
}

// ExpandScopes converts the values read from a Terraform Configuration
// file into the Scope API Model
func ExpandScopes(flatScopes interface{}) ([]waf.Scope, error) {
	if flatScopes == nil {
		return make([]waf.Scope, 0), errors.New("input was nil")
	}
	if scopes, ok := flatScopes.([]interface{}); ok {
		expandedScopes := make([]waf.Scope, len(scopes))
		for i, v := range scopes {
			if m, ok := v.(map[string]interface{}); ok {
				scope := waf.Scope{}
				scope.Name = helper.ConvertToString(
					m["name"])
				scope.Host = ExpandMatchCondition(
					m["host"])
				scope.Path = ExpandMatchCondition(
					m["path"])
				scope.ACLAuditAction = ExpandAuditAction(
					m["acl_audit_action"])
				scope.ACLAuditID = helper.ConvertToStringPointer(
					m["acl_audit_id"],
					true)
				scope.ACLProdAction = ExpandProdAction(
					m["acl_prod_action"])
				scope.ACLProdID = helper.ConvertToStringPointer(
					m["acl_prod_id"],
					true)
				scope.BotsProdAction = ExpandProdAction(
					m["bots_prod_action"])
				scope.BotsProdID = helper.ConvertToStringPointer(
					m["bots_prod_id"],
					true)
				scope.ProfileAuditAction = ExpandAuditAction(
					m["profile_audit_action"])
				scope.ProfileAuditID = helper.ConvertToStringPointer(
					m["profile_audit_id"],
					true)
				scope.ProfileProdAction = ExpandProdAction(
					m["profile_prod_action"])
				scope.ProfileProdID = helper.ConvertToStringPointer(
					m["profile_prod_id"],
					true)
				scope.RuleAuditAction = ExpandAuditAction(
					m["rule_audit_action"])
				scope.RuleAuditID = helper.ConvertToStringPointer(
					m["rule_audit_id"],
					true)
				scope.RuleProdAction = ExpandProdAction(
					m["rule_prod_action"])
				scope.RuleProdID = helper.ConvertToStringPointer(
					m["rule_prod_id"],
					true)
				limits, err := ExpandLimits(m["limit"])
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
	return make([]waf.Scope, 0), errors.New("input was not a []interface{}")
}

// ExpandAuditAction converts the values read from a Terraform Configuration
// file into the AuditAction API Model
func ExpandAuditAction(v interface{}) *waf.AuditAction {
	if v == nil {
		return nil
	}
	m, _ := helper.ConvertSingletonSetToMap(v)
	if len(m) == 0 {
		return nil
	}
	return &waf.AuditAction{
		Name: helper.ConvertToString(m["name"]),
		Type: helper.ConvertToString(m["type"]),
	}
}

// ExpandProdAction converts the values read from a Terraform Configuration
// file into the ProdAction API Model
func ExpandProdAction(v interface{}) *waf.ProdAction {
	if v == nil {
		return nil
	}
	m, _ := helper.ConvertSingletonSetToMap(v)
	if len(m) == 0 {
		return nil
	}
	return &waf.ProdAction{
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

// ExpandMatchCondition converts the values read from a Terraform Configuration
// file into the MatchCondition API Model
func ExpandMatchCondition(v interface{}) waf.MatchCondition {
	m, _ := helper.ConvertSingletonSetToMap(v)
	mc := waf.MatchCondition{
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

// ExpandLimits converts the values read from a Terraform Configuration file
// into the Limit API Model
func ExpandLimits(flatLimits interface{}) (*[]waf.Limit, error) {
	if flatLimits == nil {
		return nil, nil
	}
	if list, ok := flatLimits.([]interface{}); ok {
		if len(list) == 0 {
			return nil, nil
		}
		limits := make([]waf.Limit, len(list))
		for i, listItem := range list {
			m := listItem.(map[string]interface{})
			limits[i] = waf.Limit{
				ID: helper.ConvertToString(m["id"]),
				Action: waf.LimitAction{
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

// FlattenScopes converts the Scopes API Model
// into a format that Terraform can work with
func FlattenScopes(scopes *waf.Scopes) ([]map[string]interface{}, error) {
	if scopes == nil {
		return nil, errors.New("scopes was nil")
	}
	flattenedScopes := make([]map[string]interface{}, len(scopes.Scopes))
	for i, s := range scopes.Scopes {
		m := make(map[string]interface{})
		m["name"] = s.Name
		m["host"] = FlattenMatchCondition(s.Host)
		m["path"] = FlattenMatchCondition(s.Path)
		if s.Limits != nil {
			m["limit"] = FlattenLimits(*s.Limits)
		}
		if s.ACLAuditID != nil {
			m["acl_audit_id"] = *s.ACLAuditID
		}
		if (s.ACLAuditAction) != nil {
			m["acl_audit_action"] = FlattenAuditAction(*s.ACLAuditAction)
		}
		if s.ACLProdID != nil {
			m["acl_prod_id"] = *s.ACLProdID
		}
		if s.ACLProdAction != nil {
			m["acl_prod_action"] = FlattenProdAction(*s.ACLProdAction)
		}
		if s.BotsProdID != nil {
			m["bots_prod_id"] = *s.BotsProdID
		}
		if s.BotsProdAction != nil {
			m["bots_prod_action"] = FlattenProdAction(*s.BotsProdAction)
		}
		if s.ProfileAuditID != nil {
			m["profile_audit_id"] = *s.ProfileAuditID
		}
		if s.ProfileAuditAction != nil {
			m["profile_audit_action"] =
				FlattenAuditAction(*s.ProfileAuditAction)
		}
		if s.ProfileProdID != nil {
			m["profile_prod_id"] = *s.ProfileProdID
		}
		if s.ProfileProdAction != nil {
			m["profile_prod_action"] = FlattenProdAction(*s.ProfileProdAction)
		}
		if s.RuleAuditID != nil {
			m["rule_audit_id"] = *s.RuleAuditID
		}
		if s.RuleAuditAction != nil {
			m["rule_audit_action"] = FlattenAuditAction(*s.RuleAuditAction)
		}
		if s.RuleProdID != nil {
			m["rule_prod_id"] = *s.RuleProdID
		}
		if s.RuleProdAction != nil {
			m["rule_prod_action"] = FlattenProdAction(*s.RuleProdAction)
		}
		flattenedScopes[i] = m
	}
	return flattenedScopes, nil
}

// FlattenProdAction converts the ProdAction API Model
// into a format that Terraform can work with
func FlattenProdAction(prodAction waf.ProdAction) []map[string]interface{} {
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

// FlattenAuditAction converts the AuditAction API Model
// into a format that Terraform can work with
func FlattenAuditAction(auditAction waf.AuditAction) []map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = auditAction.Type
	m["name"] = auditAction.Name
	// We return a collection of just 1 item
	// Since we defined AuditActions as 1-item sets in the schema
	return []map[string]interface{}{m}
}

// FlattenMatchCondition converts the MatchCondition API Model
// into a format that Terraform can work with
func FlattenMatchCondition(
	matchCondition waf.MatchCondition,
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

// FlattenLimits converts the Limit API Model
// into a format that Terraform can work with
func FlattenLimits(limits []waf.Limit) []map[string]interface{} {
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

func logScopes(s waf.Scopes) {
	log.Printf("[DEBUG] Customer ID: %+v \n", s.CustomerID)
	log.Printf("[DEBUG] Name: %+v \n", s.Name)
	log.Printf("[DEBUG] Scopes: (%d items) \n\n", len(s.Scopes))
	for _, sc := range s.Scopes {
		logScope(sc)
	}
}

func logScope(s waf.Scope) {
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

func logAuditAction(a waf.AuditAction) {
	log.Printf("[DEBUG] \t\t ID: %+v \n", a.ID)
	log.Printf("[DEBUG] \t\t Name: %+v \n", a.Name)
	log.Printf("[DEBUG] \t\t Type: %+v \n", a.Type)
}

func logProdAction(a waf.ProdAction) {
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

func logLimit(l waf.Limit) {
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

func logMatchCondition(m waf.MatchCondition) {
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
