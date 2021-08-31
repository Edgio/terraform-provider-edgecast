package waf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
