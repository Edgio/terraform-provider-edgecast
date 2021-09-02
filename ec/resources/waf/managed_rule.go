// Copyright Edgecast, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package waf

import (
	"context"
	"errors"
	"fmt"
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
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies your account by its customer account number.",
			},
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
			"created_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the date and time at which the managed rule was created.",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies your account by its customer account number.",
			},
			"last_modified_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the date and time at which the managed rule was last modified.",
			},
			"last_modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reserved for future use.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reserved for future use.",
			},
			"disabled_rule": {
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
							Type:        schema.TypeInt,
							Description: "Indicates the anomaly score threshold.",
							Required:    true,
						},
						"arg_length": {
							Type:        schema.TypeInt,
							Description: "Indicates the maximum number of characters for any single query string parameter value.",
							Required:    true,
						},
						"arg_name_length": {
							Type:        schema.TypeInt,
							Description: "Indicates the maximum number of characters for any single query string parameter name.",
							Required:    true,
						},
						"combined_file_sizes": {
							Type:        schema.TypeInt,
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
			"rule_target_update": {
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
				Description: "This array describes each target using these properties",
			},
		},
	}
}

func ResourceManagedRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)

	log.Printf("[INFO] Creating WAF Managed Rule for Account >> %s", accountNumber)

	managedRuleRequest := sdkwaf.AddManagedRuleRequest{}

	managedRuleRequest.Name = d.Get("name").(string)
	managedRuleRequest.RulesetID = d.Get("ruleset_id").(string)
	managedRuleRequest.RulesetVersion = d.Get("ruleset_version").(string)

	if policies, ok := helper.ConvertInterfaceToStringArray(d.Get("policies")); ok {
		managedRuleRequest.Policies = *policies
	} else {
		return diag.Errorf("Error reading 'policies''")
	}

	// Disabled Rules
	disabledRules, err := ExpandDisabledRules(d.Get("disabled_rule"))

	if err != nil {
		return diag.Errorf("error parsing disabled_rule: %+v", err)
	}

	managedRuleRequest.DisabledRules = *disabledRules

	// Rule Target Updates
	ruleTargetUpdates, err := ExpandRuleTargetUpdates(d.Get("rule_target_update"))

	if err != nil {
		return diag.Errorf("error parsing rule_target_update: %+v", err)
	}

	managedRuleRequest.RuleTargetUpdates = *ruleTargetUpdates

	// General Settings

	if v, ok := d.GetOk("general_settings"); ok {
		if generalSettings, err := ExpandGeneralSettings(v); err == nil {
			managedRuleRequest.GeneralSettings = *generalSettings
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading General Settings",
				Detail:   err.Error(),
			})
		}
	}

	log.Printf("[DEBUG] Name: %+v\n", managedRuleRequest.Name)
	log.Printf("[DEBUG] RulesetID: %+v\n", managedRuleRequest.RulesetID)
	log.Printf("[DEBUG] Ruleset Version: %+v\n", managedRuleRequest.RulesetVersion)
	log.Printf("[DEBUG] Disabled Rules: %+v\n", managedRuleRequest.DisabledRules)
	log.Printf("[DEBUG] General Settings: %+v\n", managedRuleRequest.GeneralSettings)
	log.Printf("[DEBUG] Policies: %+v\n", managedRuleRequest.Policies)
	log.Printf("[DEBUG] Rule Target Updates: %+v\n", managedRuleRequest.RuleTargetUpdates)

	if diags.HasError() {
		return diags
	}

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := wafService.AddManagedRule(managedRuleRequest, accountNumber)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created WAF Managed Rule: %+v", resp)

	d.SetId(resp.ID)

	return ResourceManagedRuleRead(ctx, d, m)
}

func ResourceManagedRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Retrieve data for Read call
	accountNumber := d.Get("account_number").(string)
	ruleID := d.Id()

	log.Printf("[INFO] Reading WAF Managed Rule ID %s for Account >> %s", ruleID, accountNumber)

	// Initialize WAF Service
	config := m.(**api.ClientConfig)
	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Managed Rule
	resp, err := wafService.GetManagedRule(accountNumber, ruleID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Retrieved Managed Rule: %+v", resp)

	// Store all the values retrieved from the API
	d.SetId(resp.ID)
	d.Set("account_number", accountNumber)
	d.Set("created_date", resp.CreatedDate)
	d.Set("customer_id", resp.CustomerID)
	d.Set("last_modified_date", resp.LastModifiedDate)
	d.Set("last_modified_by", resp.LastModifiedBy)
	d.Set("version", resp.Version)
	d.Set("name", resp.Name)
	d.Set("ruleset_id", resp.RulesetID)
	d.Set("ruleset_version", resp.RulesetVersion)

	disabledRules := FlattenDisabledRules(resp.DisabledRules)
	d.Set("disabled_rule", disabledRules)

	generalSettings := FlattenGeneralSettings(resp.GeneralSettings)
	d.Set("general_settings", generalSettings)

	d.Set("policies", resp.Policies)

	ruleTargetUpdates := FlattenRuleTargetUpdates(resp.RuleTargetUpdates)
	d.Set("rule_target_update", ruleTargetUpdates)

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

// ExpandDisabledRules converts user-provided Terraform configuration data into the Disabled Rules API Model
func ExpandDisabledRules(attr interface{}) (*[]sdkwaf.DisabledRule, error) {

	if set, ok := attr.(*schema.Set); ok {

		items := set.List()
		disabledRules := make([]sdkwaf.DisabledRule, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			disabledRule := sdkwaf.DisabledRule{
				PolicyID: curr["policy_id"].(string),
				RuleID:   curr["rule_id"].(string),
			}

			disabledRules = append(disabledRules, disabledRule)
		}

		return &disabledRules, nil

	} else {
		return nil, errors.New("attr input was not a *schema.Set")
	}
}

// ExpandRuleTargetUpdates converts user-provided Terraform configuration data into the Rule Target Updates API Model
func ExpandRuleTargetUpdates(attr interface{}) (*[]sdkwaf.RuleTargetUpdate, error) {

	if set, ok := attr.(*schema.Set); ok {

		items := set.List()
		ruleTargetUpdates := make([]sdkwaf.RuleTargetUpdate, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			ruleTargetUpdate := sdkwaf.RuleTargetUpdate{
				IsNegated:     curr["is_negated"].(bool),
				IsRegex:       curr["is_regex"].(bool),
				ReplaceTarget: curr["replace_target"].(string),
				RuleID:        curr["rule_id"].(string),
				Target:        curr["target"].(string),
				TargetMatch:   curr["target_match"].(string),
			}

			ruleTargetUpdates = append(ruleTargetUpdates, ruleTargetUpdate)
		}

		return &ruleTargetUpdates, nil

	} else {
		return nil, errors.New("attr input was not a *schema.Set")
	}
}

// Processes all General Settings values from a Terraform configuration file into the General Settings API Model
func ExpandGeneralSettings(attr interface{}) (*sdkwaf.GeneralSettings, error) {
	if attr == nil {
		return nil, fmt.Errorf("attr was nil")
	}

	// The values are stored as a map in a 1-item set
	// So pull it out so we can work with it
	entryMap, err := helper.GetMapFromSet(attr)

	if err != nil {
		return nil, err
	}

	generalSettings := sdkwaf.GeneralSettings{}

	if anomalyThreshold, ok := entryMap["anomaly_threshold"].(int); ok {
		generalSettings.AnomalyThreshold = anomalyThreshold
	} else {
		return nil, fmt.Errorf("%v was not an int, actual: %T", entryMap["anomaly_threshold"], entryMap["anomaly_threshold"])
	}

	if argLength, ok := entryMap["arg_length"].(int); ok {
		generalSettings.ArgLength = argLength
	} else {
		return nil, fmt.Errorf("%v was not an int, actual: %T", entryMap["arg_length"], entryMap["arg_length"])
	}

	if argNameLength, ok := entryMap["arg_name_length"].(int); ok {
		generalSettings.ArgNameLength = argNameLength
	} else {
		return nil, fmt.Errorf("%v was not an int, actual: %T", entryMap["arg_name_length"], entryMap["arg_name_length"])
	}

	if combinedFileSizes, ok := entryMap["combined_file_sizes"].(int); ok {
		generalSettings.CombinedFileSizes = combinedFileSizes
	} else {
		return nil, fmt.Errorf("%v was not an int, actual: %T", entryMap["combined_file_sizes"], entryMap["combined_file_sizes"])
	}

	if ignoreCookieInterface, ok := entryMap["ignore_cookie"].([]interface{}); ok {
		if ignoreCookie, ok := helper.ConvertInterfaceArrayToStringArray(ignoreCookieInterface); ok {
			generalSettings.IgnoreCookie = ignoreCookie
		} else {
			return nil, fmt.Errorf("%v was not successfully converted to []string{}, type: %T", entryMap["ignore_cookie"], entryMap["ignore_cookie"])
		}

	} else {
		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["ignore_cookie"], entryMap["ignore_cookie"])
	}

	if ignoreHeaderInterface, ok := entryMap["ignore_header"].([]interface{}); ok {
		if ignoreHeader, ok := helper.ConvertInterfaceArrayToStringArray(ignoreHeaderInterface); ok {
			generalSettings.IgnoreHeader = ignoreHeader
		} else {
			return nil, fmt.Errorf("%v was not successfully converted to []string{}, type: %T", entryMap["ignore_header"], entryMap["ignore_header"])
		}
	} else {
		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["ignore_header"], entryMap["ignore_header"])
	}

	if ignoreQueryArgsInterface, ok := entryMap["ignore_query_args"].([]interface{}); ok {
		if ignoreQueryArgs, ok := helper.ConvertInterfaceArrayToStringArray(ignoreQueryArgsInterface); ok {
			generalSettings.IgnoreQueryArgs = ignoreQueryArgs
		} else {
			return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["ignore_query_args"], entryMap["ignore_query_args"])
		}
	} else {
		return nil, fmt.Errorf("%v was not a []string{}, actual: %T", entryMap["ignore_query_args"], entryMap["ignore_query_args"])
	}

	if jsonParser, ok := entryMap["json_parser"].(bool); ok {
		generalSettings.JsonParser = jsonParser
	} else {
		return nil, fmt.Errorf("%v was not a boolean, actual: %T", entryMap["json_parser"], entryMap["json_parser"])
	}

	if maxNumArgs, ok := entryMap["max_num_args"].(int); ok {
		generalSettings.MaxNumArgs = maxNumArgs
	} else {
		return nil, fmt.Errorf("%v was not an int, actual: %T", entryMap["max_num_args"], entryMap["max_num_args"])
	}

	if paranoiaLevel, ok := entryMap["paranoia_level"].(int); ok {
		generalSettings.ParanoiaLevel = paranoiaLevel
	} else {
		return nil, fmt.Errorf("%v was not an int, actual: %T", entryMap["paranoia_level"], entryMap["paranoia_level"])
	}

	if processRequestBody, ok := entryMap["process_request_body"].(bool); ok {
		generalSettings.ProcessRequestBody = processRequestBody
	} else {
		return nil, fmt.Errorf("%v was not a boolean, actual: %T", entryMap["process_request_body"], entryMap["process_request_body"])
	}

	if responseHeaderName, ok := entryMap["response_header_name"].(string); ok {
		generalSettings.ResponseHeaderName = responseHeaderName
	} else {
		return nil, fmt.Errorf("%v was not a string{}, actual: %T", entryMap["response_header_name"], entryMap["response_header_name"])
	}

	if totalArgLength, ok := entryMap["total_arg_length"].(int); ok {
		generalSettings.TotalArgLength = totalArgLength
	} else {
		return nil, fmt.Errorf("%v was not an int, actual: %T", entryMap["total_arg_length"], entryMap["total_arg_length"])
	}

	if validateUtf8Encoding, ok := entryMap["validate_utf8_encoding"].(bool); ok {
		generalSettings.ValidateUtf8Encoding = validateUtf8Encoding
	} else {
		return nil, fmt.Errorf("%v was not a boolean, actual: %T", entryMap["validate_utf8_encoding"], entryMap["validate_utf8_encoding"])
	}

	if xmlParser, ok := entryMap["xml_parser"].(bool); ok {
		generalSettings.XmlParser = xmlParser
	} else {
		return nil, fmt.Errorf("%v was not a boolean, actual: %T", entryMap["xml_parser"], entryMap["xml_parser"])
	}

	return &generalSettings, nil
}

// FlattenDisabledRules converts the Disabled Rules API Model into a format that Terraform can work with
func FlattenDisabledRules(disabledRules []sdkwaf.DisabledRule) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)

	for _, rule := range disabledRules {
		m := make(map[string]interface{})
		m["policy_id"] = rule.PolicyID
		m["rule_id"] = rule.RuleID
		flattened = append(flattened, m)
	}
	return flattened
}

// FlattenRuleTargetUpdates converts the Rule Target Update API Model into a format that Terraform can work with
func FlattenRuleTargetUpdates(ruleTargetUpdate []sdkwaf.RuleTargetUpdate) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)

	for _, rule := range ruleTargetUpdate {
		m := make(map[string]interface{})
		m["is_negated"] = rule.IsNegated
		m["is_regex"] = rule.IsRegex
		m["replace_target"] = rule.ReplaceTarget
		m["rule_id"] = rule.RuleID
		m["target"] = rule.Target
		m["target_match"] = rule.TargetMatch
		flattened = append(flattened, m)
	}

	return flattened
}

// FlattenGeneralSettings converts the General Settings API Model into a format that Terraform can work with
func FlattenGeneralSettings(generalSettings sdkwaf.GeneralSettings) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)

	m := make(map[string]interface{})
	m["anomaly_threshold"] = generalSettings.AnomalyThreshold
	m["arg_length"] = generalSettings.ArgLength
	m["arg_name_length"] = generalSettings.ArgNameLength
	m["combined_file_sizes"] = generalSettings.CombinedFileSizes
	m["ignore_cookie"] = generalSettings.IgnoreCookie
	m["ignore_header"] = generalSettings.IgnoreHeader
	m["ignore_query_args"] = generalSettings.IgnoreQueryArgs
	m["json_parser"] = generalSettings.JsonParser
	m["max_num_args"] = generalSettings.MaxNumArgs
	m["paranoia_level"] = generalSettings.ParanoiaLevel
	m["process_request_body"] = generalSettings.ProcessRequestBody
	m["response_header_name"] = generalSettings.ResponseHeaderName
	m["total_arg_length"] = generalSettings.TotalArgLength
	m["validate_utf8_encoding"] = generalSettings.ValidateUtf8Encoding
	m["xml_parser"] = generalSettings.XmlParser
	flattened = append(flattened, m)

	return flattened
}
