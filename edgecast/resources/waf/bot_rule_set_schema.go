// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceBotRuleSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceBotRuleSetCreate,
		ReadContext:   ResourceBotRuleSetRead,
		UpdateContext: ResourceBotRuleSetUpdate,
		DeleteContext: ResourceBotRuleSetDelete,
		Importer:      helper.Import(ResourceBotRuleSetRead, "account_number", "id"),

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Identifies your account. Find your account number in the upper right-hand corner of the MCC.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the system-defined ID for this Bot Rule Set.",
			},
			"last_modified_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the date and time at which the bot rule set was last modified.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  `Indicates the name of the bot rule set.`,
			},
			"directive": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "Contains the bot rules associated with this bot rule set. \n\n" +
					"    ->You may create up to 10 bot rules per bot rule set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Identifies a bot rule that uses our reputation database. This type of rule is satisfied when the client's IP address matches an IP address defined within our reputation database. Our reputation database contains a list of IP addresses known to be used by bots.  \n" +
								"Set this argument to the following value to include a bot rule that uses our reputation database: \n\n" +
								"        r3010_ec_bot_challenge_reputation.conf.json",
						},
						"sec_rule": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Identifies a bot rule that uses custom match conditions. This type of rule is satisfied when a match is found for each of its conditions. A condition determines request identification by defining what will be matched (i.e., variable), how it will be matched (i.e., operator), and a match value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeSet,
										Description: "Determines whether the string identified in a `variable` block will be transformed and defines metadata through which you may identify traffic to which this bot rule was applied.",
										Required:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Optional: true,
													Type:     schema.TypeString,
													Description: "Determines the custom ID that will be assigned to this rule set. This custom ID is exposed via the Browser Challenges Dashboard. Valid values fall within this range: `77000000 - 77999999`.  \n" +
														"**Default Value:** Random number \n\n" +
														"    ->This field is only applicable for the `action` block that resides in the root of the `sec_rule` block. ",
												},
												"msg": {
													Optional: true,
													Type:     schema.TypeString,
													Description: "Determines the rule message that will be assigned to this rule. This message is exposed via the Browser Challenges Dashboard.  \n" +
														"**Default Value:** Blank",
												},
												"transformations": {
													Type:     schema.TypeList,
													Optional: true,
													Description: "Determines the set of transformations that will be applied to the value  derived from the request element identified in a `variable` block (i.e., source value). Transformations are always applied to the source value, regardless of the number of transformations that have been defined.	Valid values are: \n" +
														" * `NONE` - Indicates that the source value should not be modified. \n" +
														" * `LOWERCASE` - Indicates that the source value should be converted to lowercase characters. \n" +
														" * `REMOVENULLS` - Indicates that null values should be removed from the source value. \n\n" +
														"    ->A criterion is satisfied if the source value or any of the modified string values meet the conditions defined by the `operator` block.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"chained_rule": {
										Type: schema.TypeList,
										Description: "Contains additional criteria that must be satisfied to identify traffic to which this bot  rule will be applied. \n\n" +
											"    ->You may add up to 5 `chained_rule` blocks per bot rule.",
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeSet,
													Description: "Determines whether the string value derived from the request element identified in a `variable` block will be transformed and metadata through which you may identify traffic to which this bot rule was applied.",
													Required:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"transformations": {
																Type:     schema.TypeList,
																Optional: true,
																Description: "Determines the set of transformations that will be applied to the value derived from the request element identified in a `variable` block (i.e., source value). Transformations are always applied to the source value, regardless of the number of transformations that have been defined. Valid values are: \n" +
																	" * `NONE` - Indicates that the source value should not be modified. \n" +
																	" * `LOWERCASE` - Indicates that the source value should be converted to lowercase characters. \n" +
																	" * `REMOVENULLS` - Indicates that null values should be removed from the source value. \n\n" +
																	"    ->A criterion is satisfied if the source value or any of the modified string values meet the conditions defined by the `operator` block.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"operator": {
													Type:        schema.TypeSet,
													Description: "Indicates the comparison that will be performed against the request element(s) identified within a `variable` block.",
													Required:    true,
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Indicates whether a condition will be satisfied when the value derived from the request element defined within a `variable` block matches or does not match the `value` argument. Valid values are: \n" +
																	" * `True` - Does not match \n" +
																	" * `False` - Matches ",
															},
															"type": {
																Required:     true,
																Type:         schema.TypeString,
																ValidateFunc: validation.StringIsNotWhiteSpace,
																Description: "Indicates how the system will interpret the comparison between the `value` argument and the value derived from the request element defined within a `variable` block. Valid values are: \n" +
																	" * `RX` - Indicates that the string value derived from the request element must satisfy the regular expression defined in the `value` argument. \n" +
																	" * `STREQ` - Indicates that the string value derived from the request element must be an exact match to the `value` argument. \n" +
																	" * `CONTAINS` - Indicates that the `value` argument must contain the string value derived from the request element. \n" +
																	" * `BEGINSWITH` - Indicates that the `value` argument must start with the string value derived from the request element. \n" +
																	" * `ENDSWITH` - Indicates that the `value` argument must end with the string value	derived from the request element. \n" +
																	" * `EQ` - Indicates that the number derived from the `variable` block must be an exact match to the `value` argument. \n\n" + "        ->You should only use `EQ` when the `is_count` argument has been enabled. \n" +
																	" * `IPMATCH` - Requires that the request's IP address either be contained by an IP block or be an exact match to an IP address defined in the `values` argument. Only use `IPMATCH` with the `REMOTE_ADDR` variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Indicates a value that will be compared against the string or number value derived from the request element defined within a `variable` block.",
															},
														},
													},
												},
												"variable": {
													Type:        schema.TypeList,
													Description: "Identifies each request element for which a comparison will be made.",
													Optional:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Required:     true,
																Type:         schema.TypeString,
																ValidateFunc: validation.StringIsNotWhiteSpace,
																Description: "Determines the request element that will be assessed.	Valid values are: \n\n" +
																	"        GEO | QUERY_STRING | REMOTE_ADDR | REMOTE_ASN | REQUEST_COOKIES | REQUEST_FILENAME | REQUEST_HEADERS | REQUEST_METHOD | REQUEST_URI \n\n" +
																	"    ->If a request element consists of one or more key-value pairs, then  you may identify a key via a `match` block. If `is_count` has been disabled, then you may identify a specific value via the `operator` block.",
															},
															"is_count": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Determines whether a comparison will be performed between the `operator` block and a string value or the number of matches found. Valid values are: \n" +
																	" * `True`  - A counter will increment whenever the request element defined by this `variable` block is found. The `operator` block will perform a comparison against this number. \n\n" +
																	"        ->If you enable `is_count`, then you must also set the `type` argument to `EQ`.\n" +
																	" * `False` - The `operator` block will perform a comparison against the string value derived from the request element defined by this `variable` block.",
															},
															"match": {
																Type:        schema.TypeList,
																Description: "Determines the comparison conditions for the request element identified by the `type` argument.",
																Optional:    true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_negated": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: "Determines whether this condition is satisfied when the request element identified by the `variable` block is found or not found. Valid values are: \n" +
																				" * `True` - Not found. If this argument has been enabled, then the `match` block should contain an initial block that sets both the `is_negated` and `is_regex` properties to `False`. \n" +
																				" * `False` - Found",
																		},
																		"is_regex": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: "Determines whether the `value` argument will be interpreted as a regular expression. Valid values are: \n" +
																				" * `True` - Regular expression \n" +
																				" * `False` - Default value. Literal value.",
																		},
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Description: "Restricts the match condition defined by the `type` argument to the specified value.  \n" +
																				"  \n" +
																				"**Example:**  \n" +
																				"If the `type` argument is set to `REQUEST_HEADERS` and this argument is set to	`User-Agent`, then this match condition is restricted to the `User-Agent` request header. If the `value` argument is omitted, then this match condition	applies to all request headers.",
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
										Description:  "Indicates the name assigned to this bot rule.",
									},
									"operator": {
										Type:        schema.TypeSet,
										Description: "Indicates the comparison that will be performed against the request element(s) identified within a `variable` block.",
										Required:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_negated": {
													Optional: true,
													Type:     schema.TypeBool,
													Description: "Indicates whether a condition will be satisfied when the value derived from the request element defined within a `variable` block matches or does not match the `value` argument. Valid values are: \n" +
														" * `True` - Does not match \n" +
														" * `False` - Matches",
												},
												"type": {
													Required:     true,
													Type:         schema.TypeString,
													ValidateFunc: validation.StringIsNotWhiteSpace,
													Description: "Indicates how the system will interpret the comparison between the `value` argument and the value derived from the request element defined within a `variable` block. Valid values are: \n" +
														" * `RX` - Indicates that the string value derived from the request element must satisfy the regular expression defined in the `value` argument. \n" +
														" * `STREQ` - Indicates that the string value derived from the request element must be an exact match to the `value` argument. \n" +
														" * `CONTAINS` - Indicates that the `value` argument must contain the string value derived from the request element. \n" +
														" * `BEGINSWITH` - Indicates that the `value` argument must start with the string value derived from the request element. \n" +
														" * `ENDSWITH` - Indicates that the `value` argument must end with the string value	derived from the request element. \n" +
														" * `EQ` - Indicates that the number derived from the `variable` block must be an exact match to the `value` argument. \n\n" +
														"        ->You should only use `EQ` when the `is_count` argument has been enabled. \n" +
														" * `IPMATCH` - Requires that the request's IP address either be contained by an IP block or be an exact match to an IP address defined in the `values` argument. Only use `IPMATCH` with the `REMOTE_ADDR` variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Indicates a value that will be compared against the string or number value derived from the request element defined within a `variable` block.",
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
													Required:     true,
													Type:         schema.TypeString,
													ValidateFunc: validation.StringIsNotWhiteSpace,
													Description: "Determines the request element that will be assessed.	Valid values are: \n\n" +
														"        GEO | QUERY_STRING | REMOTE_ADDR | REMOTE_ASN | REQUEST_COOKIES | REQUEST_FILENAME | REQUEST_HEADERS | REQUEST_METHOD | REQUEST_URI \n\n" +
														"    ->If a request element consists of one or more key-value pairs, then  you may identify a key via a `match` block. If `is_count` has been disabled, then you may identify a specific value via the `operator` block.",
												},
												"is_count": {
													Optional: true,
													Type:     schema.TypeBool,
													Description: "Determines whether a comparison will be performed between the `operator` block and a string value or the number of matches found. Valid values are: \n" +
														" * `True` - A counter will increment whenever the request element defined by this `variable` block is found. The `operator` block will perform a comparison against this number. \n\n" +
														"    ->If you enable `is_count`, then you must also set the `type` argument to `EQ`. \n" +
														" * `False` - The `operator` block will perform a comparison against the string value derived from the request element defined by this `variable` block.",
												},
												"match": {
													Type:        schema.TypeList,
													Description: "Contains comparison settings for the request element identified by the `type` argument.",
													Optional:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Determines whether this condition is satisfied when the request element identified by the `variable` block is found or not found. Valid values are: \n" +
																	" * `True` - Not found.  \n" +
																	" * `False` - Found",
															},
															"is_regex": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: "Determines whether the `value` argument will be interpreted as a regular expression. Valid values are: \n" +
																	" * `True` - Regular expression \n" +
																	" * `False` - Default value. Literal value.",
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Description: "Restricts the match condition defined by the `type` argument to the specified value.  \n" +
																	"  \n" +
																	"**Example:**  \n" +
																	"If the `type` argument is set to `REQUEST_HEADERS` and this argument is set to `User-Agent`, then this match condition is restricted to the `User-Agent` request header. If the `value` argument is omitted, then this match condition applies to all request headers.",
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
				},
			},
		},
	}
}
