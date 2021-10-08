// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license.
//See LICENSE file in project root for terms.

package waf

import (
	"context"
	"errors"
	"log"
	"terraform-provider-ec/ec/api"
	"terraform-provider-ec/ec/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

func ResourceCustomRuleSet() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceCustomRuleSetCreate,
		ReadContext:   ResourceCustomRuleSetRead,
		UpdateContext: ResourceCustomRuleSetUpdate,
		DeleteContext: ResourceCustomRuleSetDelete,

		Schema: map[string]*schema.Schema{
			"customer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies your account by its customer account number.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the system-defined ID for the custom rule set.",
			},
			"last_modified_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the date and time at which the custom rule was last modified.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Indicates the name of the custom rule.",
			},
			"directive": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sec_rule": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeSet,
										Description: "Determines whether the string identified in a variable object will be transformed and the metadata that will be assigned to malicious traffic.",
										Required:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Optional: true,
													Type:     schema.TypeString,
													Description: "Determines the custom ID that will be assigned to this custom rule. This custom ID is exposed via the Threats Dashboard. \\\n" +
														"    Valid values fall within this range: 66000000 - 66999999 \\\n" +
														"    *Note: This field is only applicable for the action object that resides in the root of the sec_rule object.* \\\n" +
														"    Default Value: Random number",
												},
												"msg": {
													Optional: true,
													Type:     schema.TypeString,
													Description: "Determines the rule message that will be assigned to this custom rule. This message is exposed via the Threats Dashboard. \\\n" +
														"    *Note: This field is only applicable for the action object that resides in the root of the sec_rule object.* \\\n" +
														"    Default Value: Blank",
												},
												"transformations": {
													Type:     schema.TypeList,
													Optional: true,
													Description: "Determines the set of transformations that will be applied to the value derived from the request element identified in a variable object (i.e., source value). Transformations are always applied to the source value, regardless of the number of transformations that have been defined.  \\\n" +
														"Valid Values are: \\\n" +
														"*NONE*: Indicates that the source value should not be modified. \\\n" +
														"*LOWERCASE*: Indicates that the source value should be converted to lowercase characters. \\\n" +
														"*URLDECODE*: Indicates that the source value should be URL decoded. This transformation is useful when the source value has been URL encoded twice. \\\n" +
														"*REMOVENULLS*: Indicates that null values should be removed from the source value. \\\n" +
														"*Note: A criterion is satisfied if the source value or any of the modified string values meet the conditions defined by the operator object.*",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"chained_rule": {
										Type:        schema.TypeList,
										Description: "Each object within the chained_rule array describes an additional set of criteria that must be satisfied in order to identify a malicious request.",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeSet,
													Description: "Determines whether the string identified in a variable object will be transformed and the metadata that will be assigned to malicious traffic.",
													Required:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Optional: true,
																Type:     schema.TypeString,
																Description: "Determines the custom ID that will be assigned to this custom rule. This custom ID is exposed via the Threats Dashboard. \\\n" +
																	"    Valid values fall within this range: 66000000 - 66999999 \\\n" +
																	"    *Note: This field is only applicable for the action object that resides in the root of the sec_rule object.* \\\n" +
																	"    Default Value: Random number",
															},
															"msg": {
																Optional: true,
																Type:     schema.TypeString,
																Description: "Determines the rule message that will be assigned to this custom rule. This message is exposed via the Threats Dashboard. \\\n" +
																	"    *Note: This field is only applicable for the action object that resides in the root of the sec_rule object.* \\\n" +
																	"    Default Value: Blank",
															},
															"transformations": {
																Type:     schema.TypeList,
																Optional: true,
																Description: "Determines the set of transformations that will be applied to the value derived from the request element identified in a variable object (i.e., source value). Transformations are always applied to the source value, regardless of the number of transformations that have been defined.  \\\n" +
																	"Valid Values are: \\\n" +
																	"*NONE*: Indicates that the source value should not be modified. \\\n" +
																	"*LOWERCASE*: Indicates that the source value should be converted to lowercase characters. \\\n" +
																	"*URLDECODE*: Indicates that the source value should be URL decoded. This transformation is useful when the source value has been URL encoded twice. \\\n" +
																	"*REMOVENULLS*: Indicates that null values should be removed from the source value. \\\n" +
																	"*Note: A criterion is satisfied if the source value or any of the modified string values meet the conditions defined by the operator object.*",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"operator": {
													Type:        schema.TypeSet,
													Description: "Indicates the comparison that will be performed against the request element(s) identified within a variable object.",
													Required:    true,
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Indicates whether a condition will be satisfied when the value derived from the request element defined within a variable object matches or does not match the value property.  \\\n" +
																	"    Valid values are: \\\n" +
																	"	 *True*: Does not match \\\n" +
																	"    *False*: Matches",
															},
															"type": {
																Required: true,
																Type:     schema.TypeString,
																Description: "Indicates how the system will interpret the comparison between the value property and the value derived from the request element defined within a variable object.  \\\n" +
																	"Valid Values are: \\\n" +
																	"*RX*: Indicates that the string value derived from the request element must satisfy the regular expression defined in the value property. \\\n" +
																	"*STREQ*: Indicates that the string value derived from the request element must be an exact match to the value property. \\\n" +
																	"*CONTAINS*: Indicates that the value property must contain the string value derived from the request element. \\\n" +
																	"*BEGINSWITH*: Indicates that the value property must start with the string value derived from the request element. \\\n" +
																	"*ENDSWITH*: Indicates that the value property must end with the string value derived from the request element. \\\n" +
																	"*EQ*: Indicates that the number derived from the variable object must be an exact match to the value property. \\\n" +
																	"*Note: You should only use EQ when the is_count property has been enabled.* \\\n" +
																	"*IPMATCH*: Requires that the request's IP address either be contained by an IP block or be an exact match to an IP address defined in the values property. Only use IPMATCH with the REMOTE_ADDR variable.",
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Description: "Indicates a value that will be compared against the string or number value derived from the request element defined within a variable object. \\\n" +
																	"*Note*: If you are identifying traffic via a URL path (REQUEST_URI), then you should specify a URL path pattern that starts directly after the hostname. Exclude a protocol or a hostname when defining this property. \\\n" +
																	"Sample Values: " +
																	"/marketing \\\n" +
																	"/800001/mycustomerorigin",
															},
														},
													},
												},
												"variable": {
													Type:        schema.TypeList,
													Description: "Contains criteria that identifies a request element.",
													Optional:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Required: true,
																Type:     schema.TypeString,
																Description: "Determines the request element that will be assessed. \\\n" +
																	"    Valid values are: `ARGS_POST | GEO | QUERY_STRING | REMOTE_ADDR | REQUEST_BODY | REQUEST_COOKIES | REQUEST_HEADERS | REQUEST_METHOD | REQUEST_URI` \\\n" +
																	"    *Note: If a request element consists of one or more key-value pairs, then you may identify a key via a match object. If is_count has been disabled, then you may identify a specific value via the operator object.*",
															},
															"is_count": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Determines whether a comparison will be performed between the operator object and a string value or the number of matches found. \\\n" +
																	"    Valid values are: \\\n" +
																	"    *true*: A counter will increment whenever the request element defined by this variable object is found. The operator object will perform a comparison against this number. \\\n" +
																	"    *Note: If you enable is_count, then you must also set the type property to EQ.* \\\n" +
																	"    *false*: The operator object will perform a comparison against the string value derived from the request element defined by this variable object.",
															},
															"match": {
																Type:        schema.TypeList,
																Description: "The match array determines the comparison conditions for the request element identified by the type property.",
																Optional:    true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_negated": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: "Determines whether this condition is satisfied when the request element identified by the variable object is found or not found. \\\n" +
																				"    Valid values are: \\\n" +
																				"    *True*: Not found \\\n" +
																				"    *False*: Found",
																		},
																		"is_regex": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: "Determines whether the value property will be interpreted as a regular expression. \\\n" +
																				"    Valid values are: \\\n" +
																				"    *true*: Regular expression \\\n" +
																				"    *false*: Default value. Literal value.",
																		},
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Description: "Restricts the match condition defined by the type property to the specified value. \\\n" +
																				"Example: \\\n" +
																				"If the type property is set to REQUEST_HEADERS and this property is set to User-Agent, then this match condition is restricted to the User-Agent request header. If the value property is omitted, then this match condition applies to all request headers.",
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
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Indicates the name assigned to this custom rule.",
									},
									"operator": {
										Type:        schema.TypeSet,
										Description: "Indicates the comparison that will be performed against the request element(s) identified within a variable object.",
										Required:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_negated": {
													Optional: true,
													Type:     schema.TypeBool,
													Description: "Indicates whether a condition will be satisfied when the value derived from the request element defined within a variable object matches or does not match the value property.  \\\n" +
														"    Valid values are: \\\n" +
														"	 *True*: Does not match \\\n" +
														"    *False*: Matches",
												},
												"type": {
													Required: true,
													Type:     schema.TypeString,
													Description: "Indicates how the system will interpret the comparison between the value property and the value derived from the request element defined within a variable object.  \\\n" +
														"Valid Values are: \\\n" +
														"*RX*: Indicates that the string value derived from the request element must satisfy the regular expression defined in the value property. \\\n" +
														"*STREQ*: Indicates that the string value derived from the request element must be an exact match to the value property. \\\n" +
														"*CONTAINS*: Indicates that the value property must contain the string value derived from the request element. \\\n" +
														"*BEGINSWITH*: Indicates that the value property must start with the string value derived from the request element. \\\n" +
														"*ENDSWITH*: Indicates that the value property must end with the string value derived from the request element. \\\n" +
														"*EQ*: Indicates that the number derived from the variable object must be an exact match to the value property. \\\n" +
														"*Note: You should only use EQ when the is_count property has been enabled.* \\\n" +
														"*IPMATCH*: Requires that the request's IP address either be contained by an IP block or be an exact match to an IP address defined in the values property. Only use IPMATCH with the REMOTE_ADDR variable.",
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
													Description: "Indicates a value that will be compared against the string or number value derived from the request element defined within a variable object. \\\n" +
														"*Note*: If you are identifying traffic via a URL path (REQUEST_URI), then you should specify a URL path pattern that starts directly after the hostname. Exclude a protocol or a hostname when defining this property. \\\n" +
														"Sample Values: " +
														"/marketing \\\n" +
														"/800001/mycustomerorigin",
												},
											},
										},
									},
									"variable": {
										Type:        schema.TypeList,
										Description: "Contains criteria that identifies a request element.",
										Required:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Required: true,
													Type:     schema.TypeString,
													Description: "Determines the request element that will be assessed. \\\n" +
														"    Valid values are: `ARGS_POST | GEO | QUERY_STRING | REMOTE_ADDR | REQUEST_BODY | REQUEST_COOKIES | REQUEST_HEADERS | REQUEST_METHOD | REQUEST_URI` \\\n" +
														"    *Note: If a request element consists of one or more key-value pairs, then you may identify a key via a match object. If is_count has been disabled, then you may identify a specific value via the operator object.*",
												},
												"is_count": {
													Optional: true,
													Type:     schema.TypeBool,
													Description: "Determines whether a comparison will be performed between the operator object and a string value or the number of matches found. \\\n" +
														"    Valid values are: \\\n" +
														"    *true*: A counter will increment whenever the request element defined by this variable object is found. The operator object will perform a comparison against this number. \\\n" +
														"    *Note: If you enable is_count, then you must also set the type property to EQ.* \\\n" +
														"    *false*: The operator object will perform a comparison against the string value derived from the request element defined by this variable object.",
												},
												"match": {
													Type:        schema.TypeList,
													Description: "The match array determines the comparison conditions for the request element identified by the type property.",
													Optional:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Determines whether this condition is satisfied when the request element identified by the variable object is found or not found. \\\n" +
																	"    Valid values are: \\\n" +
																	"    *True*: Not found \\\n" +
																	"    *False*: Found",
															},
															"is_regex": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Determines whether the value property will be interpreted as a regular expression. \\\n" +
																	"    Valid values are: \\\n" +
																	"    *true*: Regular expression \\\n" +
																	"    *false*: Default value. Literal value.",
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Description: "Restricts the match condition defined by the type property to the specified value. \\\n" +
																	"Example: \\\n" +
																	"If the type property is set to REQUEST_HEADERS and this property is set to User-Agent, then this match condition is restricted to the User-Agent request header. If the value property is omitted, then this match condition applies to all request headers.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
							Description: "sec_rule object describes a custom rule",
						},
					},
				},
				Description: "Contains custom rules. Each directive object defines a custom rule via the sec_rule object \\\n" +
					"    Note: You may create up to 10 custom rules.",
			},
		},
	}
}

func ResourceCustomRuleSetCreate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("customer_id").(string)

	log.Printf("[INFO] Creating WAF Rate Rule for Account >> %s", accountNumber)

	customRuleSet := sdkwaf.CustomRuleSetDetail{
		Name: d.Get("name").(string),
	}

	directive, err := ExpandDirectives(d.Get("directive"))
	if err != nil {
		return diag.Errorf("error parsing directive: %+v", err)
	}

	customRuleSet.Directives = *directive

	log.Printf("[DEBUG] Name: %+v\n", customRuleSet.Name)
	log.Printf("[DEBUG] Directive(s): %+v\n", customRuleSet.Directives)

	resp, err := wafService.AddCustomRuleSet(customRuleSet, accountNumber)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] %+v", resp)

	d.SetId(resp.ID)

	return ResourceCustomRuleSetRead(ctx, d, m)
}

func ResourceCustomRuleSetRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Get("customer_id").(string)
	ruleID := d.Id()

	log.Printf("[INFO] Retrieving custom rule %s for account number %s",
		ruleID,
		accountNumber,
	)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := wafService.GetCustomRuleSet(accountNumber, ruleID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully retrieved rate rule %s: %+v", ruleID, resp)

	d.SetId(resp.ID)
	d.Set("customer_id", accountNumber)
	d.Set("last_modified_date", resp.LastModifiedDate)
	d.Set("name", resp.Name)

	flattenDirectiveGroups := FlattenDirectives(resp.Directives)

	d.Set("directive", flattenDirectiveGroups)
	return diags
}

func ResourceCustomRuleSetUpdate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics
	return diags
}

func ResourceCustomRuleSetDelete(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics

	accountNumber := d.Get("customer_id").(string)
	customRuleID := d.Id()

	log.Printf("[INFO] Deleting WAF Custom Rule Set ID %s for Account >> %s",
		customRuleID,
		accountNumber,
	)

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := wafService.DeleteCustomRuleSet(accountNumber, customRuleID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully deleted WAF Custom Rule Set: %+v", resp)

	d.SetId("")

	return diags
}

func ExpandDirectives(attr interface{}) (*[]sdkwaf.Directive, error) {

	if set, ok := attr.(*schema.Set); ok {

		items := set.List()

		directives := make([]sdkwaf.Directive, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			directive := sdkwaf.Directive{}

			secRule, err := ExpandSecRule(curr["sec_rule"])
			if err != nil {
				return nil, err
			}

			directive.SecRule = *secRule

			directives = append(directives, directive)
		}

		return &directives, nil

	} else {
		return nil,
			errors.New("ExpandDirectives: attr input was not a *schema.Set")
	}

}

func ExpandSecRule(attr interface{}) (*sdkwaf.SecRule, error) {

	if attr == nil {
		return nil, errors.New("sec rule attr was nil")
	}

	curr, err := helper.GetMapFromSet(attr)
	if err != nil {
		return nil, err
	}

	secRule := sdkwaf.SecRule{
		Name: curr["name"].(string),
	}

	actionMap, err := helper.GetMapFromSet(curr["action"])
	if err != nil {
		return nil, err
	}

	operatorMap, err := helper.GetMapFromSet(curr["operator"])
	if err != nil {
		return nil, err
	}

	chainedRule, err := ExpandChainedRules(curr["chained_rule"])
	if err != nil {
		return nil, err
	}
	secRule.ChainedRules = *chainedRule

	variables, err := ExpandVariables(curr["variable"])
	if err != nil {
		return nil, err
	}
	secRule.Variables = *variables

	if actionId, ok := actionMap["id"]; ok {
		secRule.Action.ID = actionId.(string)
	}

	if actionMsg, ok := actionMap["msg"]; ok {
		secRule.Action.Message = actionMsg.(string)
	}

	if actionT, ok := actionMap["transformations"]; ok {
		if arr, ok := helper.ConvertInterfaceToStringArray(actionT); ok {
			secRule.Action.Transformations = *arr
		}
	}

	if v, ok := operatorMap["is_negated"]; ok {
		secRule.Operator.IsNegated = v.(bool)
	}

	if operatorType, ok := operatorMap["type"]; ok {
		secRule.Operator.Type = operatorType.(string)
	}

	if operatorValue, ok := operatorMap["value"]; ok {
		secRule.Operator.Value = operatorValue.(string)
	}

	return &secRule, nil
}

func ExpandChainedRules(attr interface{}) (*[]sdkwaf.ChainedRule, error) {

	if items, ok := attr.([]interface{}); ok {
		chainedRules := make([]sdkwaf.ChainedRule, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			chainedRule := sdkwaf.ChainedRule{}

			actionMap, err := helper.GetMapFromSet(curr["action"])
			if err != nil {
				return nil, err
			}

			operatorMap, err := helper.GetMapFromSet(curr["operator"])
			if err != nil {
				return nil, err
			}

			variables, err := ExpandVariables(curr["variable"])
			if err != nil {
				return nil, err
			}
			chainedRule.Variables = *variables

			if actionId, ok := actionMap["id"]; ok {
				chainedRule.Action.ID = actionId.(string)
			}

			if actionMsg, ok := actionMap["msg"]; ok {
				chainedRule.Action.Message = actionMsg.(string)
			}

			if actionT, ok := actionMap["transformations"]; ok {
				if arr, ok := helper.ConvertInterfaceToStringArray(actionT); ok {
					chainedRule.Action.Transformations = *arr
				}
			}

			if v, ok := operatorMap["is_negated"]; ok {
				chainedRule.Operator.IsNegated = v.(bool)
			}

			if operatorType, ok := operatorMap["type"]; ok {
				chainedRule.Operator.Type = operatorType.(string)
			}

			if operatorValue, ok := operatorMap["value"]; ok {
				chainedRule.Operator.Value = operatorValue.(string)
			}

			chainedRules = append(chainedRules, chainedRule)
		}

		return &chainedRules, nil

	} else {
		return nil,
			errors.New("ExpandChainedRules: attr input was not a []interface{}")
	}
}

func ExpandVariables(attr interface{}) (*[]sdkwaf.Variable, error) {

	if items, ok := attr.([]interface{}); ok {

		variables := make([]sdkwaf.Variable, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			variable := sdkwaf.Variable{
				Type:    curr["type"].(string),
				IsCount: curr["is_count"].(bool),
			}

			matches, err := ExpandMatches(curr["match"])
			if err != nil {
				return nil, err
			}
			variable.Matches = *matches

			variables = append(variables, variable)
		}

		return &variables, nil

	} else {
		return nil,
			errors.New("ExpandVariables: attr input was not a []interface{}")
	}
}

func ExpandMatches(attr interface{}) (*[]sdkwaf.Match, error) {

	if items, ok := attr.([]interface{}); ok {
		matches := make([]sdkwaf.Match, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			match := sdkwaf.Match{
				IsNegated: curr["is_negated"].(bool),
				IsRegex:   curr["is_regex"].(bool),
				Value:     curr["value"].(string),
			}
			matches = append(matches, match)
		}

		return &matches, nil

	} else {
		return nil,
			errors.New("ExpandMatches: attr input was not a []interface{}")
	}
}

// FlattenDirectives converts the ConditionGroup API Model
// into a format that Terraform can work with
func FlattenDirectives(directive []sdkwaf.Directive) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)

	for _, cg := range directive {
		m := make(map[string]interface{})

		m["sec_rule"] = FlattenSecRule(cg.SecRule)

		flattened = append(flattened, m)
	}
	return flattened
}

// FlattenSecRule converts the Condition API Model
// into a format that Terraform can work with
func FlattenSecRule(secrule sdkwaf.SecRule) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	m["action"] = FlattenAction(secrule.Action)
	m["chained_rule"] = FlattenChainedrule(secrule.ChainedRules)
	m["operator"] = FlattenOperator(secrule.Operator)
	m["variable"] = FlattenVariable(secrule.Variables)

	flattened = append(flattened, m)
	return flattened
}

// FlattenAction converts the Condition API Model
// into a format that Terraform can work with
func FlattenAction(action sdkwaf.Action) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	m["id"] = action.ID
	m["msg"] = action.Message
	m["transformations"] = action.Transformations

	flattened = append(flattened, m)
	return flattened
}

// FlattenChainrule converts the Condition API Model
// into a format that Terraform can work with
func FlattenChainedrule(chain []sdkwaf.ChainedRule) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)

	for _, cg := range chain {
		m := make(map[string]interface{})

		m["action"] = FlattenAction(cg.Action)
		m["operator"] = FlattenOperator(cg.Operator)
		m["variable"] = FlattenVariable(cg.Variables)
		flattened = append(flattened, m)
	}
	return flattened
}

// FlattenAction converts the Condition API Model
// into a format that Terraform can work with
func FlattenOperator(operator sdkwaf.Operator) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	m["is_negated"] = operator.IsNegated
	m["type"] = operator.Type
	m["value"] = operator.Value

	flattened = append(flattened, m)
	return flattened
}

// FlattenVariable converts the Condition API Model
// into a format that Terraform can work with
func FlattenVariable(val []sdkwaf.Variable) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)
	for _, cg := range val {
		m := make(map[string]interface{})

		m["type"] = cg.Type
		m["match"] = FlattenMatch(cg.Matches)
		m["is_count"] = cg.IsCount

		flattened = append(flattened, m)
	}
	return flattened
}

// FlattenMatch converts the Condition API Model
// into a format that Terraform can work with
func FlattenMatch(match []sdkwaf.Match) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)
	for _, cg := range match {
		m := make(map[string]interface{})

		m["is_negated"] = cg.IsNegated
		m["is_regex"] = cg.IsRegex
		m["value"] = cg.Value

		flattened = append(flattened, m)
	}
	return flattened
}
