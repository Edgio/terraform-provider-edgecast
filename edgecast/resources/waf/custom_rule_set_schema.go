// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCustomRuleSet() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceCustomRuleSetCreate,
		ReadContext:   ResourceCustomRuleSetRead,
		UpdateContext: ResourceCustomRuleSetUpdate,
		DeleteContext: ResourceCustomRuleSetDelete,
		Importer:      helper.Import(ResourceCustomRuleSetRead, "account_number", "id"),

		Schema: map[string]*schema.Schema{
			"account_number": {
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
