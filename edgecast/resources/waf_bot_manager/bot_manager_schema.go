// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf_bot_manager

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetBotManagerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Description: "The unique name by which this Bot Manager configuration will be identified. \n" +
				"This name should be sufficiently descriptive to identify it when setting up a Security Application Manager configuration.",
		},
		"bots_prod_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Indicates the system-defined ID assigned to an existing Bot Rule.",
		},
		"actions": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"alert": {
						Type: schema.TypeList,
						Description: "Configuration for generating an alert. \n" +
							"Use this mode to track detected threats through the Bots dashboard without impacting production traffic.",
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Computed:    true,
									Type:        schema.TypeString,
									Description: "Indicates the system-defined ID assigned to this alert.",
								},
								"name": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "The name by which this alert configuration will be identified.",
								},
							},
						},
					},
					"custom_response": {
						Type:        schema.TypeList,
						Description: "Configuration for returning a custom response.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Computed:    true,
									Type:        schema.TypeString,
									Description: "Indicates the system-defined ID assigned to this custom response.",
								},
								"name": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "The name by which this custom response configuration will be identified.",
								},
								"response_body_base64": {
									Optional: true,
									Type:     schema.TypeString,
									Description: "Defines the payload that will be delivered to the client. \n" +
										" * This option supports the use of event variables to customize the response. \n\n" +
										"**Sample payload for a HTML file:** \n" +
										"```<!DOCTYPE html><html> \n" +
										"<head><title>Page Not Found</title></head> \n" +
										"<body>Page not found.</body> \n" +
										"</html>```",
								},
								"status": {
									Optional: true,
									Type:     schema.TypeInt,
									Description: "Defines the HTTP status code that will be sent to the client. \n" +
										" --> Value must be of format 'uint32'",
								},
								"response_headers": {
									Type:     schema.TypeMap,
									Optional: true,
									Description: "Defines one or more response headers that will be sent to the client. Define each custom response header on a separate line. \n" +
										"Syntax: {Name}:{Value} \n" +
										"Example: MyCustomHeader: True \n" +
										" * This option supports the use of event variables to customize the response. \n" +
										" * All characters, including spaces, defined before or after the colon will be treated as a part of the specified header name or value, respectively.",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
					"block_request": {
						Type:        schema.TypeList,
						Description: "Configuration for dropping the request and providing the client with a 403 Forbidden response.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Computed:    true,
									Type:        schema.TypeString,
									Description: "Indicates the system-defined ID assigned to this block request.",
								},
								"name": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "The name by which this block request configuration will be identified.",
								},
							},
						},
					},
					"redirect_302": {
						Type:        schema.TypeList,
						Description: "Configuration for Redirecting requests to the specified URL. The HTTP status code for this response will be a 302 Found.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Computed:    true,
									Type:        schema.TypeString,
									Description: "Indicates the system-defined ID assigned to this redirect.",
								},
								"name": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "The name by which this redirect configuration will be identified.",
								},
								"url": {
									Optional: true,
									Type:     schema.TypeString,
									Description: "The full URL to which requests will be redirected. \n" +
										"Example: http://cdn.mydomain.com/marketing/busy.html",
								},
							},
						},
					},
					"browser_challenge": {
						Type: schema.TypeList,
						Description: "Configuration for Sending a browser challenge to the client. The client must solve this challenge within a few seconds. \n" +
							"**Key Information:** \n" +
							" * Solving a challenge requires a JavaScript-enabled client. Users that have disabled JavaScript on their browsing session will be unable to access content protected by browser challenges. \n" +
							" * We strongly recommend that you avoid applying browser challenges to machine-to-machine interactions. For example, applying browser challenges to API traffic will disrupt your API workflow. \n" +
							" * The HTTP Status Code option determines the HTTP status code for the response provided to clients that are being served the browser challenge. \n" +
							"   - Setting this option to certain status codes (e.g., 204) may prevent clients that successfully solve a browser challenge from properly displaying your site. \n" +
							" * You may define a custom payload for the browser challenge by enabling the Custom Browser Challenge Page option and then setting the Browser Challenge Page Template option to the desired payload.",
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Computed:    true,
									Type:        schema.TypeString,
									Description: "Indicates the system-defined ID assigned to this browser challenge.",
								},
								"name": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "The name by which this browser challenge will be identified.",
								},
								"is_custom_challenge": {
									Optional:    true,
									Type:        schema.TypeBool,
									Description: "Valid Values: True | False",
								},
								"response_body_base64": {
									Optional: true,
									Type:     schema.TypeString,
									Description: "Describes the results of the above browser challenge determines what happens next. \n" +
										" * Solved: If the client is able to solve the challenge, then our CDN serves the requested content. Additionally, a cookie will be added to the user's session. This cookie instructs our CDN to serve content to the user without requiring a browser challenge. Once the cookie expires, new requests for content protected by Bot Manager will once again require the client to solve a challenge. \n" +
										" * Unsolved: If the client is unable to solve the challenge, then our CDN responds with a new browser challenge.",
								},
								"valid_for_sec": {
									Optional:    true,
									Type:        schema.TypeInt,
									Description: "Defines the duration for the cookie. Value must be of format 'uint32'",
								},
								"status": {
									Optional: true,
									Type:     schema.TypeInt,
									Description: "Determines the HTTP status code for the response provided to clients that are being served the browser challenge. \n " +
										"**Note:** Setting this option to certain status codes (e.g., 204) may prevent clients that successfully solve a browser challenge from properly displaying your site. \n" +
										"Value must be of format 'uint32'",
								},
							},
						},
					},
				},
			},
			Description: "Contains the type of actions that can be applied to bot traffic, and can include browser challenge, custom response, or redirect that can be applied to known bots, spoofed bots, and bots detected through rules. \n\n" +
				"    --> Unlike other actions, alert and block actions do not require configuration before they can be applied to bot traffic.",
		},
		"exception_cookie": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Bypass the above bot detection measures by creating an exception for one or more cookie(s).",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"exception_ja3": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Bypass the above bot detection measures by creating an exception for one or more JA3 fingerprint(s).",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"exception_url": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Bypass the above bot detection measures by creating an exception for one or more URL(s).",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"exception_user_agent": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Bypass the above bot detection measures by creating an exception for one or more user agent(s).",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"inspect_known_bots": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Valid Values: True | False.",
		},
		"known_bot": {
			Type:        schema.TypeList,
			Description: "List of known bots.",
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
						Description: "Must be one of available bot tokens",
					},
				},
			},
		},
		"last_modified_date": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies the date and time when the Bot Manager was most recently modiified.",
		},
		"last_modified_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Identifies the user account which most recently modified the Bot Manager.",
		},
		"spoof_bot_action_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Valid Values : ALERT, BLOCK_REQUEST, CUSTOM_RESPONSE, BROWSER_CHALLENGE, REDIRECT_302",
		},
	}
}
