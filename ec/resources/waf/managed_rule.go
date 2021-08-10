// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package waf

import (
	"context"
	"log"
	"terraform-provider-ec/ec/api"
	"terraform-provider-ec/ec/helper"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceManagedRule() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceManagedRuleCreate,
		ReadContext:   ResourceManagedRuleRead,
		UpdateContext: ResourceManagedRuleUpdate,
		DeleteContext: ResourceManagedRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates the name of the custom rule.",
			},
			"ruleset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the ID for the rule set associated with this managed rule.",
			},
			"ruleset_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the version of the rule set associated with this managed rule.",
			},
			"disabled_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Description: "Identifies a policy from which a rule will be disabled by its system-defined ID.",
							Optional:    true,
						},
						"rule_id": {
							Type:        schema.TypeString,
							Description: "Identifies a rule that will be disabled by its system-defined ID.",
							Optional:    true,
						},
					},
				},
				Description: "This array identifies each rule that has been disabled using these properties",
			},
			"general_settings": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"anomaly_threshold": {
							Type:        schema.TypeString,
							Description: "Indicates the anomaly score threshold.",
							Required:    true,
						},
						"arg_length": {
							Type:        schema.TypeString,
							Description: "Indicates the maximum number of characters for any single query string parameter value.",
							Required:    true,
						},
						"arg_name_length": {
							Type:        schema.TypeString,
							Description: "Indicates the maximum number of characters for any single query string parameter name.",
							Required:    true,
						},
						"combined_file_sizes": {
							Type:        schema.TypeString,
							Description: "Indicates the total file size for multipart message lengths.",
							Optional:    true,
						},
						"ignore_cookie": {
							Type:        schema.TypeList,
							Description: "Identifies each cookie that will be ignored for the purpose of determining whether a request is malicious traffic.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ignore_header": {
							Type:        schema.TypeList,
							Description: "Identifies each request header that will be ignored for the purpose of determining whether a request is malicious traffic.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ignore_query_args": {
							Type:        schema.TypeList,
							Description: "Identifies each query string argument that will be ignored for the purpose of determining whether a request is malicious traffic.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"json_parser": {
							Type:        schema.TypeBool,
							Description: "Determines whether JSON payloads will be inspected. Valid values are: true | false",
							Optional:    true,
						},
						"max_num_args": {
							Type:        schema.TypeInt,
							Description: "Indicates the maximum number of query string parameters.",
							Required:    true,
						},
						"paranoia_level": {
							Type:        schema.TypeInt,
							Description: "Indicates the balance between the level of protection and false positives. Valid values are: 1 | 2 | 3 | 4",
							Optional:    true,
						},
						"process_request_body": {
							Type:        schema.TypeBool,
							Description: "Determines whether JSON payloads will be inspected.",
							Optional:    true,
						},
						"response_header_name": {
							Type:        schema.TypeString,
							Description: "Determines the name of the response header that will be included with blocked requests.",
							Optional:    true,
						},
						"total_arg_length": {
							Type:        schema.TypeInt,
							Description: "Indicates the maximum number of characters for the query string value.",
							Required:    true,
						},
						"validate_utf8_encoding": {
							Type:        schema.TypeBool,
							Description: "Indicates whether WAF may check whether a request variable (e.g., ARGS, ARGS_NAMES, and REQUEST_FILENAME) is a valid UTF-8 string.",
							Optional:    true,
						},
						"xml_parser": {
							Type:        schema.TypeBool,
							Description: "Determines whether XML payloads will be inspected.",
							Optional:    true,
						},
					},
				},
				Description: "Contains settings that define the profile for a valid request.",
			},
			"policies": {
				Type:        schema.TypeList,
				Description: "Contains a list of policies that have been enabled on this managed rule.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rule_target_updates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_negated": {
							Type:        schema.TypeBool,
							Description: "Determines whether the current target, as defined within this object, will be ignored when identifying threats.",
							Optional:    true,
						},
						"is_regex": {
							Type:        schema.TypeBool,
							Description: "Identifies a rule that will be disabled by its system-defined ID.",
							Required:    true,
						},
						"replace_target": {
							Type:        schema.TypeString,
							Description: "Defines the data source (e.g., REQUEST_COOKIES, ARGS, GEO, etc.) that will be used instead of the one defined in the target parameter.",
							Optional:    true,
						},
						"rule_id": {
							Type:        schema.TypeString,
							Description: "Identifies a rule by its system-defined ID.",
							Required:    true,
						},
						"target": {
							Type:        schema.TypeString,
							Description: "Identifies the type of data source (e.g., REQUEST_COOKIES, ARGS, GEO, etc.) for which a target will be created.",
							Required:    true,
						},
						"target_match": {
							Type:        schema.TypeString,
							Description: "Identifies a name or category (e.g., cookie name, query string name, country code, etc.) for the data source defined in the target parameter.",
							Required:    true,
						},
					},
				},
				Description: "This array identifies each rule that has been disabled using these properties",
			},
		},
	}
}

func ResourceManagedRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)

	log.Printf("[INFO] Creating WAF Managed Rule for Account >> %s", accountNumber)

	managedRule := sdkwaf.ManagedRule{
		Name:              d.Get("name").(string),
		RulesetID:         d.Get("ruleset_id").(string),
		RulesetVersion:    d.Get("ruleset_version").(string),
		DisabledRules:     d.Get("disabled_rules").List(),
		Policies:          helper.ConvertInterfaceToStringArray(d.Get("polices")),
		RuleTargetUpdates: d.Get("rule_target_updates").List(),
	}

	generalSettings, err := helper.GetMapFromSet(d.Get("general_settings"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading General Settings",
			Detail:   err.Error(),
		})
	}

	managedRule.GeneralSettings = generalSettings

	log.Printf("[DEBUG] Allowed HTTP Methods: %+v\n", accessRule.AllowedHTTPMethods)
	log.Printf("[DEBUG] Allowed Request Content Types: %+v\n", accessRule.AllowedRequestContentTypes)
	log.Printf("[DEBUG] Disallowed Headers: %+v\n", accessRule.DisallowedHeaders)
	log.Printf("[DEBUG] Disallowed Extensions: %+v\n", accessRule.DisallowedExtensions)
	log.Printf("[DEBUG] ASN: %+v\n", accessRule.ASNAccessControls)
	log.Printf("[DEBUG] Cookie: %+v\n", accessRule.CookieAccessControls)
	log.Printf("[DEBUG] Country: %+v\n", accessRule.CountryAccessControls)
	log.Printf("[DEBUG] IP: %+v\n", accessRule.IPAccessControls)
	log.Printf("[DEBUG] Referer: %+v\n", accessRule.RefererAccessControls)
	log.Printf("[DEBUG] URL: %+v\n", accessRule.URLAccessControls)
	log.Printf("[DEBUG] User Agent: %+v\n", accessRule.UserAgentAccessControls)

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := wafService.AddAccessRule(accessRule)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created WAF Managed Rule: %+v", resp)

	d.SetId(resp.ID)

	return diags
}

func ResourceManagedRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceManagedRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceManagedRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// func ConvertInterfaceToAccessControls(attr interface{}) (*sdkwaf.AccessControls, error) {
// 	if attr == nil {
// 		return nil, fmt.Errorf("attr was nil")
// 	}

// 	// The values are stored as a map in a 1-item set
// 	// So pull it out so we can work with it
// 	entryMap, err := helper.GetMapFromSet(attr)

// 	if err != nil {
// 		return nil, err
// 	}

// 	accessControls := &sdkwaf.AccessControls{}

// 	if accesslist, ok := entryMap["accesslist"].([]interface{}); ok {
// 		accessControls.Accesslist = accesslist
// 	} else {
// 		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["accesslist"], entryMap["accesslist"])
// 	}

// 	if blacklist, ok := entryMap["blacklist"].([]interface{}); ok {
// 		accessControls.Blacklist = blacklist
// 	} else {
// 		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["blacklist"], entryMap["blacklist"])
// 	}

// 	if whitelist, ok := entryMap["whitelist"].([]interface{}); ok {
// 		accessControls.Whitelist = whitelist
// 	} else {
// 		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["whitelist"], entryMap["whitelist"])
// 	}

// 	return accessControls, nil
// }
