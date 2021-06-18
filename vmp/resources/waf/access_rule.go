// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package waf

import (
	"context"
	"log"
	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAccessRule() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceAccessRuleCreate,
		ReadContext:   ResourceAccessRuleRead,
		UpdateContext: ResourceAccessRuleUpdate,
		DeleteContext: ResourceAccessRuleDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Identifies your account by its customer account number.",
			},
			"allowed_http_methods": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Identifies each allowed HTTP method (e.g., GET).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_request_content_types": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Identifies each allowed media type (e.g., application/json).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"asn": {
				Type:        schema.TypeMap,
				Description: "Contains access controls for autonomous system numbers (ASNs).",
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
			"cookie": {
				Type:        schema.TypeMap,
				Description: "Contains access controls for cookies.",
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"country": {
				Type:        schema.TypeMap,
				Description: "Contains access controls for countries. Specify each desired country using its country code.",
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"disallowed_extensions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Indicates each file extension for which WAF will send an alert or block the request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"disallowed_headers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Indicates each request header for which WAF will send an alert or block the request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip": {
				Type:        schema.TypeMap,
				Description: "Contains access controls for IPv4 and/or IPv6 addresses. Specify each desired IP address using standard IPv4/IPv6 and CIDR notation.",
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Assigns a name to this access rule.",
			},
			"referer": {
				Type:        schema.TypeMap,
				Description: "Contains access controls for referrers. All referrers defined within a whitelist, accesslist, or blacklist are regular expressions.",
				// TODO: validate func for regex?
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"response_header_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Determines the name of the response header that will be included with blocked requests.",
			},
			"url": {
				Type: schema.TypeMap,
				Description: "Contains access controls for URL paths. Specify a URL path pattern that starts directly after the hostname. " +
					"Exclude a protocol or a hostname when defining value | values. Sample values: /marketing, /800001/mycustomerorigin " +
					"All URL paths defined within a whitelist, accesslist, or blacklist are regular expressions.",
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"user_agent": {
				Type:        schema.TypeMap,
				Description: "Contains access controls for user agents.  All user agents defined within a whitelist, accesslist, or blacklist are regular expressions.",
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func ResourceAccessRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	(*config).AccountNumber = accountNumber
	log.Printf("[INFO] Creating WAF Access Rule for Account >> [AccountNumber]: %s", accountNumber)

	accessRule := api.AccessRule{
		AllowedHTTPMethods:         d.Get("allowed_http_methods").([]string),
		AllowedRequestContentTypes: d.Get("allowed_request_content_types").([]string),
		ASNAccessControls:          getAccessControls(d, "asn"),
		CookieAccessControls:       getAccessControls(d, "cookie"),
		CountryAccessControls:      getAccessControls(d, "country"),
		CustomerID:                 accountNumber,
		DisallowedExtensions:       d.Get("disallowed_exensions").([]string),
		DisallowedHeaders:          d.Get("disallowed_headers").([]string),
		IPAccessControls:           getAccessControls(d, "ip"),
		Name:                       d.Get("name").(string),
		RefererAccessControls:      getAccessControls(d, "referer"),
		ResponseHeaderName:         d.Get("response_header_name").(string),
		URLAccessControls:          getAccessControls(d, "url"),
	}

	log.Printf("[DEBUG] Access Rule %+v\n", accessRule)

	// apiClient := api.NewWAFAPIClient(*config)

	// accessRuleId, err := apiClient.AddAccessRule(accessRule)

	// if err != nil {
	// 	d.SetId("")
	// 	return diag.FromErr(err)
	// }

	// log.Printf("[INFO] Successfully created WAF Access Rule, ID=%d", accessRuleId)

	// d.SetId(strconv.Itoa(accessRuleId))

	return diags
}

func ResourceAccessRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceAccessRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceAccessRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func getAccessControls(d *schema.ResourceData, key string) api.AccessControls {
	var accessControls api.AccessControls

	if attr, getOk := d.GetOk(key); getOk {
		log.Printf("[DEBUG] %s Access Controls raw: %+v\n", key, attr)

		if m, mapCastOk := attr.(map[string]interface{}); mapCastOk {
			log.Printf("[DEBUG] %s Access Controls map: %+v\n", key, m)

			accessControls.AccessList = m["accesslist"].([]interface{})
			accessControls.Blacklist = m["blacklist"].([]interface{})
			accessControls.Whitelist = m["whitelist"].([]interface{})
		}
	}

	return accessControls
}
