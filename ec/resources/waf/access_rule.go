// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package waf

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-ec/ec/api"
	"terraform-provider-ec/ec/helper"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"

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
				Description: "Contains access controls for autonomous system numbers (ASNs).  \\\n" +
					"*Note: ASN access controls are integer values.*",
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
				Description: "Contains access controls for referrers.  \\\n" +
					"*Note: All referrers defined within a whitelist, accesslist, or blacklist are regular expressions.*",
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
					"Exclude a protocol or a hostname when defining value | values. Sample values: /marketing, /800001/mycustomerorigin.  \\\n" +
					"*Note: All URL paths defined within a whitelist, accesslist, or blacklist are regular expressions.*",
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
				Description: "Contains access controls for user agents.  \\\n" +
					"*Note: All user agents defined within a whitelist, accesslist, or blacklist are regular expressions.*",
			},
		},
	}
}

func ResourceAccessRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)

	log.Printf("[INFO] Creating WAF Access Rule for Account >> %s", accountNumber)

	accessRule := sdkwaf.AccessRule{
		CustomerID:         accountNumber,
		Name:               d.Get("name").(string),
		ResponseHeaderName: d.Get("response_header_name").(string),
	}

	if v, ok := d.GetOk("allowed_http_methods"); ok {
		if allowedHttpControlsPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.AllowedHTTPMethods = *allowedHttpControlsPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Allowed HTTP Methods",
			})
		}
	}

	if v, ok := d.GetOk("allowed_request_content_types"); ok {
		if allowedRequestContentTypesPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.AllowedRequestContentTypes = *allowedRequestContentTypesPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Allowed Request Content Types",
			})
		}
	}

	if v, ok := d.GetOk("disallowed_headers"); ok {
		if disallowedHeadersPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.DisallowedHeaders = *disallowedHeadersPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Disallowed Headers",
			})
		}
	}

	if v, ok := d.GetOk("disallowed_extensions"); ok {
		if disallowedExtensionsPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.DisallowedExtensions = *disallowedExtensionsPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Disallowed Extensions",
			})
		}
	}

	if v, ok := d.GetOk("asn"); ok {
		if asnAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.ASNAccessControls = asnAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading ASN Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("cookie"); ok {
		if cookieAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.CookieAccessControls = cookieAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Cookie Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("country"); ok {
		if countryAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.CountryAccessControls = countryAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Country Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("ip"); ok {
		if ipAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.IPAccessControls = ipAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading IP Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("referer"); ok {
		if refererAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.RefererAccessControls = refererAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Referer Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("url"); ok {
		if urlAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.URLAccessControls = urlAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading URL Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("user_agent"); ok {
		if userAgentAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.UserAgentAccessControls = userAgentAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading User Agent Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	helper.LogInstanceAsPrettyJson("[DEBUG] ACCESSRULE", accessRule)
	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := wafService.AddAccessRule(accessRule)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created WAF Access Rule: %+v", resp)

	d.SetId(resp.ID)
	ResourceAccessRuleRead(ctx, d, m)

	return ResourceAccessRuleRead(ctx, d, m)
}

func ResourceAccessRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	ruleID := d.Id()
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	//change the order
	log.Printf("[INFO] Reading WAF Access Rule Id %s for Account >> %s", ruleID, accountNumber)

	resp, err := wafService.GetAccessRuleByID(accountNumber, ruleID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	//change %v with printing funcation.
	helper.LogInstanceAsPrettyJson("[INFO] Retrieved Rule", resp)

	d.SetId(resp.ID)
	d.Set("account_number", accountNumber)
	d.Set("allowed_http_methods", resp.AllowedHTTPMethods)
	d.Set("allowed_request_content_types", resp.AllowedRequestContentTypes)
	flattenedASN := FlattenAccessControls(*resp.ASNAccessControls)
	d.Set("asn", flattenedASN)
	flattenedCookie := FlattenAccessControls(*resp.CookieAccessControls)
	d.Set("cookie", flattenedCookie)
	flattenedCountry := FlattenAccessControls(*resp.CountryAccessControls)
	d.Set("country", flattenedCountry)
	d.Set("disallowed_extensions", resp.DisallowedExtensions)
	d.Set("disallowed_headers", resp.DisallowedHeaders)
	flattenedIp := FlattenAccessControls(*resp.IPAccessControls)
	d.Set("ip", flattenedIp)
	d.Set("name", resp.Name)
	flattenedReferer := FlattenAccessControls(*resp.RefererAccessControls)
	d.Set("referer", flattenedReferer)
	d.Set("response_header_name", resp.ResponseHeaderName)
	flattenedUrl := FlattenAccessControls(*resp.URLAccessControls)
	d.Set("url", flattenedUrl)
	flattenedUserAgent := FlattenAccessControls(*resp.UserAgentAccessControls)
	d.Set("user_agent", flattenedUserAgent)

	return diags
}

func ResourceAccessRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	ruleID := d.Id()
	log.Printf("[INFO] Updating WAF Access Rule for Account >> %s", accountNumber)

	accessRule := sdkwaf.AccessRule{
		CustomerID:         accountNumber,
		Name:               d.Get("name").(string),
		ResponseHeaderName: d.Get("response_header_name").(string),
	}

	if v, ok := d.GetOk("allowed_http_methods"); ok {
		if allowedHttpControlsPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.AllowedHTTPMethods = *allowedHttpControlsPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Allowed HTTP Methods",
			})
		}
	}

	if v, ok := d.GetOk("allowed_request_content_types"); ok {
		if allowedRequestContentTypesPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.AllowedRequestContentTypes = *allowedRequestContentTypesPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Allowed Request Content Types",
			})
		}
	}

	if v, ok := d.GetOk("disallowed_headers"); ok {
		if disallowedHeadersPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.DisallowedHeaders = *disallowedHeadersPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Disallowed Headers",
			})
		}
	}

	if v, ok := d.GetOk("disallowed_extensions"); ok {
		if disallowedExtensionsPtr, ok := helper.ConvertInterfaceToStringArray(v); ok {
			accessRule.DisallowedExtensions = *disallowedExtensionsPtr
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Disallowed Extensions",
			})
		}
	}

	if v, ok := d.GetOk("asn"); ok {
		if asnAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.ASNAccessControls = asnAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading ASN Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("cookie"); ok {
		if cookieAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.CookieAccessControls = cookieAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Cookie Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("country"); ok {
		if countryAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.CountryAccessControls = countryAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Country Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("ip"); ok {
		if ipAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.IPAccessControls = ipAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading IP Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("referer"); ok {
		if refererAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.RefererAccessControls = refererAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Referer Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("url"); ok {
		if urlAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.URLAccessControls = urlAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading URL Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	if v, ok := d.GetOk("user_agent"); ok {
		if userAgentAccessControls, err := ConvertInterfaceToAccessControls(v); err == nil {
			accessRule.UserAgentAccessControls = userAgentAccessControls
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading User Agent Access Controls",
				Detail:   err.Error(),
			})
		}
	}

	helper.ConvertInterfaceToStringArray(accessRule)
	helper.LogInstanceAsPrettyJson("[DEBUG] ACCESSRULE", accessRule)
	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := wafService.UpdateAccessRule(accessRule, ruleID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully updated WAF Access Rule: %+v", resp)

	return ResourceAccessRuleRead(ctx, d, m)

}

func ResourceAccessRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ConvertInterfaceToAccessControls(attr interface{}) (*sdkwaf.AccessControls, error) {
	if attr == nil {
		return nil, fmt.Errorf("attr was nil")
	}

	// The values are stored as a map in a 1-item set
	// So pull it out so we can work with it
	entryMap, err := helper.GetMapFromSet(attr)

	if err != nil {
		return nil, err
	}

	accessControls := &sdkwaf.AccessControls{}

	if accesslist, ok := entryMap["accesslist"].([]interface{}); ok {
		accessControls.Accesslist = accesslist
	} else {
		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["accesslist"], entryMap["accesslist"])
	}

	if blacklist, ok := entryMap["blacklist"].([]interface{}); ok {
		accessControls.Blacklist = blacklist
	} else {
		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["blacklist"], entryMap["blacklist"])
	}

	if whitelist, ok := entryMap["whitelist"].([]interface{}); ok {
		accessControls.Whitelist = whitelist
	} else {
		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["whitelist"], entryMap["whitelist"])
	}

	return accessControls, nil
}

// FlattenAccessControls converts the AccessControls API Model
// into a format that Terraform can work with
func FlattenAccessControls(accessControlsGroups sdkwaf.AccessControls) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	m["accesslist"] = accessControlsGroups.Accesslist
	m["blacklist"] = accessControlsGroups.Blacklist
	m["whitelist"] = accessControlsGroups.Whitelist

	flattened = append(flattened, m)

	return flattened
}
