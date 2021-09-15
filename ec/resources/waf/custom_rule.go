package waf

import (
	"context"
	"log"
	"terraform-provider-ec/ec/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

func ResourceCustomRule() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceCustomRuleCreate,
		ReadContext:   ResourceCustomRuleRead,
		UpdateContext: ResourceCustomRuleUpdate,
		DeleteContext: ResourceCustomRuleDelete,

		Schema: map[string]*schema.Schema{
			"customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the system-defined ID for the custom rule set.",
			},
			"directive": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sec_rule": {
							Type:        schema.TypeSet,
							Description: "Contains list of custom rule set.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeSet,
										Description: "Determines whether the string identified in a variable object will be transformed and the metadata that will be assigned to malicious traffic.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Computed: true,
													Type:     schema.TypeString,
													Description: "Determines the custom ID that will be assigned to this custom rule. \\\n" +
														"	This custom ID is exposed via the Threats Dashboard. Valid values fall within this range: `66000000 - 66999999`",
												},
												"msg": {
													Optional:    true,
													Type:        schema.TypeString,
													Description: "Determines the rule message that will be assigned to this custom rule. This message is exposed via the Threats Dashboard.",
												},
												"t": {
													Optional:    true,
													Type:        schema.TypeSet,
													Description: "Determines the set of transformations that will be applied to the value derived from the request element identified in a variable object (i.e., source value).",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"chained_rule": {
										Type:        schema.TypeSet,
										Description: "Contains additional criteria that must be satisfied to identify a malicious request.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeSet,
													Description: "Determines whether the string identified in a variable object will be transformed and the metadata that will be assigned to malicious traffic.",
													Optional:    true,
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Computed: true,
																Type:     schema.TypeString,
																Description: "Determines the custom ID that will be assigned to this custom rule. \\\n" +
																	"	This custom ID is exposed via the Threats Dashboard. Valid values fall within this range: `66000000 - 66999999`",
															},
															"msg": {
																Optional:    true,
																Type:        schema.TypeString,
																Description: "Determines the rule message that will be assigned to this custom rule. This message is exposed via the Threats Dashboard.",
															},
															"t": {
																Optional:    true,
																Type:        schema.TypeSet,
																Description: "Determines the set of transformations that will be applied to the value derived from the request element identified in a variable object (i.e., source value).",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"operator": {
													Type: schema.TypeSet,
													Description: "describes the comparison that will be performed on the request \\\n" +
														"	element(s) defined within a variable object.",
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Indicates whether a condition will be satisfied when \\\n" +
																	"	the value derived from the request element defined within a variable object matches or does not match the value property.",
															},
															"type": {
																Optional: true,
																Type:     schema.TypeString,
																Description: "Indicates how the system will interpret the \\\n" +
																	"	comparison between the value property and the value derived from the request element defined within a variable object.",
															},
															"value": {
																Optional: true,
																Type:     schema.TypeString,
																Description: "Indicates a value that will be compared against \\\n" +
																	"	the string or number value derived from the request element defined within a variable object.",
															},
														},
													},
												},
												"variable": {
													Type: schema.TypeSet,
													Description: "describes the comparison that will be performed on the request \\\n" +
														"	element(s) defined within a variable object.",
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Optional:    true,
																Type:        schema.TypeString,
																Description: "Determines the request element that will be assessed.",
															}, "match": {
																Optional:    true,
																Type:        schema.TypeBool,
																Description: "Contains comparison settings for the request element identified by the type property.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_negated": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: "Indicates whether a condition will be satisfied when \\\n" +
																				"	the value derived from the request element defined within a variable object matches or does not match the value property.",
																		},
																		"is_regex": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: "Indicates a value that will be compared against \\\n" +
																				"	the string or number value derived from the request element defined within a variable object.",
																		},
																		"value": {
																			Optional: true,
																			Type:     schema.TypeString,
																			Description: "Indicates a value that will be compared against \\\n" +
																				"	the string or number value derived from the request element defined within a variable object.",
																		},
																	},
																},
															},

															"is_count": {
																Optional:    true,
																Type:        schema.TypeBool,
																Description: "Determines whether a comparison will be performed between the operator object and a string value or the number of matches found.",
															},
														},
													},
												},
											},
										},
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the system-defined ID for this custom rule.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the name assigned to this custom rule.",
									},
									"operator": {
										Type: schema.TypeSet,
										Description: "describes the comparison that will be performed on the request \\\n" +
											"	element(s) defined within a variable object.",
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_negated": {
													Optional: true,
													Type:     schema.TypeBool,
													Description: "Indicates whether a condition will be satisfied when \\\n" +
														"	the value derived from the request element defined within a variable object matches or does not match the value property.",
												},
												"type": {
													Optional: true,
													Type:     schema.TypeString,
													Description: "Indicates how the system will interpret the \\\n" +
														"	comparison between the value property and the value derived from the request element defined within a variable object.",
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
													Description: "Indicates a value that will be compared against \\\n" +
														"	the string or number value derived from the request element defined within a variable object.",
												},
											},
										},
									},
									"variable": {
										Type: schema.TypeSet,
										Description: "describes the comparison that will be performed on the request \\\n" +
											"	element(s) defined within a variable object.",
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Optional:    true,
													Type:        schema.TypeString,
													Description: "Determines the request element that will be assessed.",
												}, "match": {
													Optional:    true,
													Type:        schema.TypeSet,
													Description: "Contains comparison settings for the request element identified by the type property.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Indicates whether a condition will be satisfied when \\\n" +
																	"	the value derived from the request element defined within a variable object matches or does not match the value property.",
															},
															"is_regex": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Indicates a value that will be compared against \\\n" +
																	"	the string or number value derived from the request element defined within a variable object.",
															},
															"value": {
																Optional: true,
																Type:     schema.TypeString,
																Description: "Indicates a value that will be compared against \\\n" +
																	"	the string or number value derived from the request element defined within a variable object.",
															},
														},
													},
												},

												"is_count": {
													Optional:    true,
													Type:        schema.TypeBool,
													Description: "Determines whether a comparison will be performed between the operator object and a string value or the number of matches found.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Description: "Contains the set of condition groups associated with a rule.",
			},
		},
	}
}
func ResourceCustomRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceCustomRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	ruleID := d.Id()

	log.Printf("[INFO] Retrieving custom rule %s for account number %s", ruleID, accountNumber)

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

	flattenDirectiveGroups := FlattenDirectiveGroups(resp.Directives)

	d.Set("directive", flattenDirectiveGroups)

	return diags
}

func ResourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceCustomRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// FlattenDirectiveGroups converts the ConditionGroup API Model
// into a format that Terraform can work with
func FlattenDirectiveGroups(directive []sdkwaf.Directive) []map[string]interface{} {

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
	m["t"] = action.Transformations

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
