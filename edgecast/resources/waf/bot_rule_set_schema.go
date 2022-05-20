// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceBotRuleSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceBotRuleSetCreate,
		ReadContext:   ResourceBotRuleSetRead,
		UpdateContext: ResourceBotRuleSetUpdate,
		DeleteContext: ResourceBotRuleSetDelete,

		Schema: map[string]*schema.Schema{
			"customer_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Identifies your account by its customer account 
				number.`,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the system-defined ID for this Bot Rule 
				Set.`,
			},
			"last_modified_date": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the date and time at which the bot rule 
				set was last modified.`,
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
				Description: `Contains the bot rules associated with this bot rule set. 
				You may create up to 10 bot rules per bot rule set.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `Identifies a bot rule that uses our reputation database. This 
							type of rule is satisfied when the client's IP address matches an IP address 
							defined within our reputation database. Our reputation database contains a 
							list of IP addresses known to be used by bots. Set this property to the following 
							value to include a bot rule that uses our reputation database: 
							r3010_ec_bot_challenge_reputation.conf.json`,
						},
						"sec_rule": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Description: `Identifies a bot rule that uses custom match conditions. This
							type of rule is satisfied when a match is found for each of its conditions.
							A condition determines request identification by defining what will be 
							matched (i.e., variable), how it will be matched (i.e., operator), and a 
							match value.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type: schema.TypeSet,
										Description: `Determines whether the string identified in a variable object 
										will be transformed and metadata through which you may identify traffic to 
										which this bot rule was applied.`,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Optional: true,
													Type:     schema.TypeString,
													Description: `Determines the custom ID that will be assigned to this rule 
													set. This custom ID is exposed via the Browser Challenges 
													Dashboard. 
													 
													Valid values fall within this range: 77000000 - 77999999
													
													Default Value: Random number`,
												},
												"msg": {
													Optional: true,
													Type:     schema.TypeString,
													Description: `Determines the rule message that will be assigned to this rule.
													Thismessage is exposed via the Browser Challenges Dashboard.
													 
													Default Value: Blank`,
												},
												"transformations": {
													Type:     schema.TypeList,
													Optional: true,
													Description: `Determines the set of transformations that will be applied to
													the value derived from the request element identified in a variable object
													(i.e., source value). Transformations are always applied to the source value,
													regardless of the number of transformations that have been defined.

													 Valid Values are: 
													 *NONE*: Indicates that the source value should not be modified.
													 *LOWERCASE*: Indicates that the source value should be converted to 
													 lowercase characters.
													 *REMOVENULLS*: Indicates that null values should be removed from the source
													 value.

													 *Note: A criterion is satisfied if the source value or any of the modified
													 string values meet the conditions defined by the operator object.*`,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"chained_rule": {
										Type: schema.TypeList,
										Description: `Contains additional criteria that must be satisfied to 
										identify traffic to which this bot  rule will be applied. You may add up to
										5 chained_rule objects per bot rule.`,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type: schema.TypeSet,
													Description: `Determines whether the string value derived from the request element identified
													in a variable object will be transformed and metadata through which you may identify traffic 
													to which this bot rule was applied.`,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"transformations": {
																Type:     schema.TypeList,
																Optional: true,
																Description: `Determines the set of transformations that will be applied to 
																the value derived from the request element identified in a variable object 
																(i.e., source value). Transformations are always applied to the source value,
																regardless of the number of transformations that have been defined.

																Valid Values are: 

																*NONE*: Indicates that the source value should not be modified.
																*LOWERCASE*: Indicates that the source value should be converted to lowercase 
																characters.
																*REMOVENULLS*: Indicates that null values should be removed from the source  value.
																
																*Note: A criterion is satisfied if the source value or any of the modified 
																string values meet the conditions defined by the operator object.*`,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"operator": {
													Type: schema.TypeSet,
													Description: `Indicates the comparison that will be performed against the 
													request element(s) identified within a variable object.`,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: `Indicates whether a condition will be satisfied when the value
																derived from the request element defined within a variable object matches or
																does not match the value property. 
																
																Valid values are:

																*True*: Does not match
																*False*: Matches`,
															},
															"type": {
																Required: true,
																Type:     schema.TypeString,
																Description: `Indicates how the system will interpret the comparison between
																the value property and the value derived from the request element defined 
																within a variable object.
																
																Valid Values are:

																*RX*: Indicates that the string value derived from the request element must
																satisfy the regular expression defined in the value property.
																*STREQ*: Indicates that the string value derived from the request element 
																must be an exact match to the value property.
																*CONTAINS*: Indicates that the value property must contain the string value
																derived from the request element.
																*BEGINSWITH*: Indicates that the value property must start with the string 
																value derived from the request element.
																*ENDSWITH*: Indicates that the value property must end with the string value
																derived from the request element.
																*EQ*: Indicates that the number derived from the variable object must be an
																exact match to the value property.
																*Note: You should only use EQ when the is_count property has been enabled.*
																*IPMATCH*: Requires that the request's IP address either be contained by an
																IP block or be an exact match to an IP address defined in the values 
																property. Only use IPMATCH with the REMOTE_ADDR variable.`,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Description: `Indicates a value that will be compared against the string or
																number value derived from the request element defined within a variable 
																object.`,
															},
														},
													},
												},
												"variable": {
													Type: schema.TypeList,
													Description: `Identifies each request element for which a comparison will be
													made.`,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Required: true,
																Type:     schema.TypeString,
																Description: `Determines the request element that will be assessed.
																
																Valid values are: ARGS_POST | GEO | QUERY_STRING | REMOTE_ADDR | REQUEST_BODY | REQUEST_COOKIES | REQUEST_HEADERS | REQUEST_METHOD | REQUEST_URI
																
																*Note: If a request element consists of one or more key-value pairs, then 
																you may identify a key via a match object. If is_count has been disabled, 
																then you may identify a specific value via the operator object.*`,
															},
															"is_count": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: `Determines whether a comparison will be performed between the
																operator object and a string value or the number of matches found.
																
																Valid values are:

																*true*: A counter will increment whenever the request element defined by 
																this variable object is found. The operator object will perform a comparison
																against this number. *Note: If you enable is_count, then you must also set
																the type property to EQ.*
																*false*: The operator object will perform a comparison against the string 
																value derived from the request element defined by this variable object.`,
															},
															"match": {
																Type: schema.TypeList,
																Description: `Determines the comparison conditions for the request element
																identified by the type property.`,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_negated": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: `Determines whether this condition is satisfied when the request 
																			element identified by the variable object is found or not found.
																			
																			Valid values are:
																			*True*: Not found. If this property has been enabled, then the match array 
																			should contain an initial object that sets both the is_negated and is_regex
																			properties to False.
																			*False*: Found`,
																		},
																		"is_regex": {
																			Optional: true,
																			Type:     schema.TypeBool,
																			Description: `Determines whether the value property will be interpreted as 
																			a regular expression.
																			
																			Valid values are:
																			*true*: Regular expression
																			*false*: Default value. Literal value.`,
																		},
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Description: `Restricts the match condition defined by the type property to
																			the specified value.

																			Example:
																			If the type property is set to REQUEST_HEADERS and this property is set to
																			User-Agent, then this match condition is restricted to the User-Agent 
																			request header. If the value property is omitted, then this match condition
																			applies to all request headers.`,
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
										Description:  `Indicates the name assigned to this bot rule.`,
									},
									"operator": {
										Type: schema.TypeSet,
										Description: `Indicates the comparison that will be performed against the 
										request element(s) identified within a variable object.`,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_negated": {
													Optional: true,
													Type:     schema.TypeBool,
													Description: `Indicates whether a condition will be satisfied when the 
													value derived from the request element defined within a variable object 
													matches or does not match the value property.
													
													Valid values are:
													*True*: Does not match
													*False*: Matches`,
												},
												"type": {
													Required: true,
													Type:     schema.TypeString,
													Description: `Indicates how the system will interpret the comparison between
													the value property and the value derived from the request element defined 
													within a variable object.
													
													Valid Values are:
													*RX*: Indicates that the string value derived from the request element must
													satisfy the regular expression defined in the value property. 
													*STREQ*: Indicates that the string value derived from the request element 
													must be an exact match to the value property.  
													*CONTAINS*: Indicates that the value property must contain the string value
													derived from the request element.  
													*BEGINSWITH*: Indicates that the value property must start with the string 
													value derived from the request element.  
													*ENDSWITH*: Indicates that the value property must end with the string 
													value derived from the request element.  
													*EQ*: Indicates that the number derived from the variable object must be 
													an exact match to the value property.  
													*Note: You should only use EQ when the is_count property has been enabled.*  
													*IPMATCH*: Requires that the request's IP address either be contained by 
													an IP block or be an exact match to an IP address defined in the values 
													property. Only use IPMATCH with the REMOTE_ADDR variable.`,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
													Description: `Indicates a value that will be compared against the string 
													or number value derived from the request element defined within a variable
													object.`,
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
													Description: `Determines the request element that will be assessed.
													
													Valid values are: ARGS_POST | GEO | QUERY_STRING | REMOTE_ADDR | REQUEST_BODY | REQUEST_COOKIES | REQUEST_HEADERS | REQUEST_METHOD | REQUEST_URI 
													
													*Note: If a request element consists of one or more key-value pairs, 
													then you may identify a key via a match object. If is_count has been 
													disabled, then you may identify a specific value via the operator object.*`,
												},
												"is_count": {
													Optional: true,
													Type:     schema.TypeBool,
													Description: `Determines whether a comparison will be performed between the
													operator object and a string value or the number of matches found.
													
													Valid values are:
													*true*: A counter will increment whenever the request element defined by 
													this variable object is found. The operator object will perform a comparison
													against this number. *Note: If you enable is_count, then you must also set
													the type property to EQ.*
													*false*: The operator object will perform a comparison against the string
													value derived from the request element defined by this variable object.`,
												},
												"match": {
													Type: schema.TypeList,
													Description: `Contains comparison settings for the request element 
													identified by the type property.`,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_negated": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: `Determines whether this condition is satisfied when the
																request element identified by the variable object is found or not found.  
																
																Valid values are:
																*True*: Not found
																*False*: Found`,
															},
															"is_regex": {
																Optional: true,
																Type:     schema.TypeBool,
																Description: `Determines whether the value property will be interpreted as
																a regular expression.  
																
																Valid values are:
																*true*: Regular expression
																*false*: Default value. Literal value.`,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Description: `Restricts the match condition defined by the type property 
																to the specified value.  
																
																Example:
																If the type property is set to REQUEST_HEADERS and this property is set 
																to User-Agent, then this match condition is restricted to the User-Agent 
																request header. If the value property is omitted, then this match 
																condition applies to all request headers.`,
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
