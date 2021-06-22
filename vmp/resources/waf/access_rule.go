// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package waf

import (
	"context"
	"log"
	"terraform-provider-vmp/vmp/api"
	"terraform-provider-vmp/vmp/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAccessRule() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceAccessRuleCreate,
		ReadContext:   ResourceAccessRuleRead,
		UpdateContext: ResourceAccessRuleUpdate,
		DeleteContext: ResourceAccessRuleDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies your account by its customer account number.",
			},
			"allowed_http_methods": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Identifies each allowed HTTP method (e.g., GET).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_request_content_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Identifies each allowed media type (e.g., application/json).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"asn": {
				// We use a 1-item TypeSet as a workaround since TypeMap doesn't support schema.Resource as a child element type (yet)
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accesslist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content upon passing a threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"blacklist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that will be blocked or for which an alert will be generated.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"whitelist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content without undergoing threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
				Description: "Contains access controls for autonomous system numbers (ASNs). Note: ASN access controls are integer values.",
			},
			"cookie": {
				// We use a 1-item TypeSet as a workaround since TypeMap doesn't support schema.Resource as a child element type (yet)
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accesslist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content upon passing a threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"blacklist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that will be blocked or for which an alert will be generated.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"whitelist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content without undergoing threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Description: "Contains access controls for cookies.",
			},
			"country": {
				// We use a 1-item TypeSet as a workaround since TypeMap doesn't support schema.Resource as a child element type (yet)
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accesslist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content upon passing a threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"blacklist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that will be blocked or for which an alert will be generated.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"whitelist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content without undergoing threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Description: "Contains access controls for countries. Specify each desired country using its country code.",
			},
			"disallowed_extensions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Indicates each file extension for which WAF will send an alert or block the request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"disallowed_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Indicates each request header for which WAF will send an alert or block the request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip": {
				// We use a 1-item TypeSet as a workaround since TypeMap doesn't support schema.Resource as a child element type (yet)
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accesslist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content upon passing a threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"blacklist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that will be blocked or for which an alert will be generated.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"whitelist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content without undergoing threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Description: "Contains access controls for IPv4 and/or IPv6 addresses. Specify each desired IP address using standard IPv4/IPv6 and CIDR notation.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Assigns a name to this access rule.",
			},
			"referer": {
				// We use a 1-item TypeSet as a workaround since TypeMap doesn't support schema.Resource as a child element type (yet)
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accesslist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content upon passing a threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"blacklist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that will be blocked or for which an alert will be generated.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"whitelist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content without undergoing threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Description: "Contains access controls for referrers. Note: All referrers defined within a whitelist, accesslist, or blacklist are regular expressions.",
			},
			"response_header_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Determines the name of the response header that will be included with blocked requests.",
			},
			"url": {
				// We use a 1-item TypeSet as a workaround since TypeMap doesn't support schema.Resource as a child element type (yet)
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accesslist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content upon passing a threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"blacklist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that will be blocked or for which an alert will be generated.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"whitelist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content without undergoing threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Description: "Contains access controls for URL paths. Specify a URL path pattern that starts directly after the hostname. " +
					"Exclude a protocol or a hostname when defining value | values. Sample values: /marketing, /800001/mycustomerorigin. " +
					"Note: All URL paths defined within a whitelist, accesslist, or blacklist are regular expressions.",
			},
			"user_agent": {
				// We use a 1-item TypeSet as a workaround since TypeMap doesn't support schema.Resource as a child element type (yet)
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accesslist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content upon passing a threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"blacklist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that will be blocked or for which an alert will be generated.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"whitelist": {
							Type:        schema.TypeList,
							Description: "Contains entries that identify traffic that may access your content without undergoing threat assessment.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Description: "Contains access controls for user agents. Note: All user agents defined within a whitelist, accesslist, or blacklist are regular expressions.",
			},
		},
	}
}

func ResourceAccessRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	(*config).AccountNumber = accountNumber

	log.Printf("[INFO] Creating WAF Access Rule for Account >> %s", accountNumber)

	accessRule := api.AccessRule{
		AllowedHTTPMethods:         getStringArray(d, "allowed_http_methods"),
		AllowedRequestContentTypes: getStringArray(d, "allowed_request_content_types"),
		ASNAccessControls:          getAccessControls(d, "asn"),
		CookieAccessControls:       getAccessControls(d, "cookie"),
		CountryAccessControls:      getAccessControls(d, "country"),
		CustomerID:                 accountNumber,
		DisallowedHeaders:          getStringArray(d, "disallowed_headers"),
		DisallowedExtensions:       getStringArray(d, "disallowed_extensions"),
		IPAccessControls:           getAccessControls(d, "ip"),
		Name:                       d.Get("name").(string),
		RefererAccessControls:      getAccessControls(d, "referer"),
		ResponseHeaderName:         d.Get("response_header_name").(string),
		URLAccessControls:          getAccessControls(d, "url"),
		UserAgentAccessControls:    getAccessControls(d, "user_agent"),
	}

	log.Printf("[DEBUG] Allowed HTTP Methods: %+v\n", accessRule.AllowedHTTPMethods)
	log.Printf("[DEBUG] Allowed Request Content Types: %+v\n", accessRule.AllowedRequestContentTypes)
	log.Printf("[DEBUG] Disallowed Headers: %+v\n", accessRule.DisallowedHeaders)
	log.Printf("[DEBUG] Disallowed Extensions: %+v\n", accessRule.DisallowedExtensions)
	log.Printf("[DEBUG] ASN: %+v\n", accessRule.ASNAccessControls)
	log.Printf("[DEBUG] Cookie: %+v\n", accessRule.CookieAccessControls)
	log.Printf("[DEBUG] Country: %+v\n", accessRule.CountryAccessControls)
	log.Printf("[DEBUG] IP: %+v\n", accessRule.IPAccessControls)
	log.Printf("[DEBUG] Referer: %+v\n", accessRule.RefererAccessControls)
	log.Printf("[DEBUG] URL: %+v\n", accessRule.URLAccessControls)
	log.Printf("[DEBUG] User Agent: %+v\n", accessRule.UserAgentAccessControls)

	apiClient := api.NewWAFAPIClient(*config)

	accessRuleId, err := apiClient.AddAccessRule(accessRule)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created WAF Access Rule, ID=%s", accessRuleId)

	d.SetId(accessRuleId)

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

func getStringArray(d *schema.ResourceData, key string) []string {
	var values []string

	if attr, ok := d.GetOk(key); ok {
		interfaceArray := attr.([]interface{})
		values, _ = helper.InterfaceArrayToStringArray(interfaceArray)
	}

	return values
}

func getAccessControls(d *schema.ResourceData, key string) *api.AccessControls {
	var accessControls *api.AccessControls

	if attr, ok := d.GetOk(key); ok {
		// Must convert the set to a slice/array to work with the values
		set := attr.(*schema.Set)
		arr := set.List()

		// The single item in the list a map[string][]interface{}
		entryMap := arr[0].(map[string]interface{})

		// Each interface{} in the map is a []interface{}
		accessControls = &api.AccessControls{
			Accesslist: entryMap["accesslist"].([]interface{}),
			Blacklist:  entryMap["blacklist"].([]interface{}),
			Whitelist:  entryMap["whitelist"].([]interface{}),
		}
	}

	return accessControls
}
