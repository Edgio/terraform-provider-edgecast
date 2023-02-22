// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf_bot_manager

import (
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBotManagerAdvanced() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceBotManagerCreate,
		ReadContext:   ResourceBotManagerRead,
		UpdateContext: ResourceBotManagerUpdate,
		DeleteContext: ResourceBotManagerDelete,
		Importer:      helper.Import(ResourceBotManagerRead, "account_number", "id"),

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Indicates the system-defined ID assigned to this Bot Manager.",
				Computed:    true,
			},
			"customer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the customer id.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "the unique name by which this Bot Manager configuration will be identified. /n" +
					"This name should be sufficiently descriptive to identify it when setting up a Security Application Manager configuration.",
			},
			"bots_prod_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bots Production ID.",
			},
			"action": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert": {
							Type:        schema.TypeList,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"name": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"enf_type": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "enum. Value = ALERT",
									},
								},
							},
						},
						"custom_response": {
							Type:        schema.TypeList,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"name": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"enf_type": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "enum. Default value = CUSTOM_RESPONSE",
									},
									"response_body_base64": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"status": {
										Optional:    true,
										Type:        schema.TypeInt,
										Description: "Value must be of format 'uint32'",
									},
									"response_headers": {
										Optional:    true,
										Type:        schema.TypeSet,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Optional:    true,
													Type:        schema.TypeString,
													Description: "",
												},
												"value": {
													Optional:    true,
													Type:        schema.TypeString,
													Description: "",
												},
											},
										},
									},
								},
							},
						},
						"block_request": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"name": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"enf_type": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "enum. Default value = BLOCK_REQUEST",
									},
								},
							},
						},
						"redirect_302": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"name": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"enf_type": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "enum. Default value = REDIRECT_302",
									},
									"url": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
								},
							},
						},
						"browser_challenge": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"name": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"enf_type": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "enum. Default value = BROWSER_CHALLENGE",
									},
									"is_custom_challenge": {
										Optional:    true,
										Type:        schema.TypeBool,
										Description: "",
									},
									"response_body_base64": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "",
									},
									"valid_for_sec": {
										Optional:    true,
										Type:        schema.TypeInt,
										Description: "Value must be of format 'uint32'",
									},
									"status": {
										Optional:    true,
										Type:        schema.TypeInt,
										Description: "Value must be of format 'uint32'",
									},
								},
							},
						},
					},
				},
				Description: "Contains browser challenge, custom response, or redirect that can be applied to known bots, spoofed bots, and bots detected through rules. \n\n" +
					"    --> Unlike other actions, alert and block actions do not require configuration before they can be applied to bot traffic.",
			},
			"exception_cookie": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exception_ja3": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exception_url": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exception_user_agent": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"inspect_known_bots": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
			},
			"known_bots": {
				Type:        schema.TypeList,
				Description: "",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_type": {
							Optional:    true,
							Type:        schema.TypeString,
							Description: "Valid Values : ALERT, BLOCK_REQUEST, CUSTOM_RESPONSE, BROWSER_CHALLENGE, REDIRECT_302",
						},
						"bot_token": {
							Optional:    true,
							Type:        schema.TypeString,
							Description: "must be one of available bot tokens",
						},
					},
				},
			},
			"last_modified_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"last_modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"spoof_bot_action_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Valid Values : ALERT, BLOCK_REQUEST, CUSTOM_RESPONSE, BROWSER_CHALLENGE, REDIRECT_302",
			},
		},
	}
}
