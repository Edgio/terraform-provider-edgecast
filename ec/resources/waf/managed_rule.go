// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package waf

import (
	"context"
	"log"
	"terraform-provider-ec/ec/api"
	"terraform-provider-ec/ec/helper"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceManagedRule() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceManagedRuleCreate,
		ReadContext:   ResourceManagedRuleRead,
		UpdateContext: ResourceManagedRuleUpdate,
		DeleteContext: ResourceManagedRuleDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies your account by its customer account number.",
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
					},
				},
				Description: "Contains access controls for autonomous system numbers (ASNs).  \\\n" +
					"*Note: ASN access controls are integer values.*",
			},
		},
	}
}

func ResourceManagedRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)

	log.Printf("[INFO] Creating WAF Managed Rule for Account >> %s", accountNumber)

	accessRule := sdkwaf.AccessRule{
		AllowedHTTPMethods:         helper.ConvertInterfaceToStringArray(d.Get("allowed_http_methods")),
		AllowedRequestContentTypes: helper.ConvertInterfaceToStringArray(d.Get("allowed_request_content_types")),
		CustomerID:                 accountNumber,
		DisallowedHeaders:          helper.ConvertInterfaceToStringArray(d.Get("disallowed_headers")),
		DisallowedExtensions:       helper.ConvertInterfaceToStringArray(d.Get("disallowed_extensions")),
		Name:                       d.Get("name").(string),
		ResponseHeaderName:         d.Get("response_header_name").(string),
	}

	if asnAccessControls, err := ConvertInterfaceToAccessControls(d.Get("asn")); err == nil {
		accessRule.ASNAccessControls = asnAccessControls
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading ASN Access Controls",
			Detail:   err.Error(),
		})
	}

	if cookieAccessControls, err := ConvertInterfaceToAccessControls(d.Get("cookie")); err == nil {
		accessRule.CookieAccessControls = cookieAccessControls
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading Cookie Access Controls",
			Detail:   err.Error(),
		})
	}

	if countryAccessControls, err := ConvertInterfaceToAccessControls(d.Get("country")); err == nil {
		accessRule.CountryAccessControls = countryAccessControls
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading Country Access Controls",
			Detail:   err.Error(),
		})
	}

	if ipAccessControls, err := ConvertInterfaceToAccessControls(d.Get("ip")); err == nil {
		accessRule.IPAccessControls = ipAccessControls
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading IP Access Controls",
			Detail:   err.Error(),
		})
	}

	if refererAccessControls, err := ConvertInterfaceToAccessControls(d.Get("referer")); err == nil {
		accessRule.RefererAccessControls = refererAccessControls
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading Referer Access Controls",
			Detail:   err.Error(),
		})
	}

	if urlAccessControls, err := ConvertInterfaceToAccessControls(d.Get("url")); err == nil {
		accessRule.URLAccessControls = urlAccessControls
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading URL Access Controls",
			Detail:   err.Error(),
		})
	}

	if userAgentAccessControls, err := ConvertInterfaceToAccessControls(d.Get("user_agent")); err == nil {
		accessRule.UserAgentAccessControls = userAgentAccessControls
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading User Agent Access Controls",
			Detail:   err.Error(),
		})
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

	log.Printf("[INFO] Successfully created WAF Managed Rule: %+v", resp)

	d.SetId(resp.ID)

	return diags
}

func ResourceManagedRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceManagedRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceManagedRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// func ConvertInterfaceToAccessControls(attr interface{}) (*sdkwaf.AccessControls, error) {
// 	if attr == nil {
// 		return nil, fmt.Errorf("attr was nil")
// 	}

// 	// The values are stored as a map in a 1-item set
// 	// So pull it out so we can work with it
// 	entryMap, err := helper.GetMapFromSet(attr)

// 	if err != nil {
// 		return nil, err
// 	}

// 	accessControls := &sdkwaf.AccessControls{}

// 	if accesslist, ok := entryMap["accesslist"].([]interface{}); ok {
// 		accessControls.Accesslist = accesslist
// 	} else {
// 		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["accesslist"], entryMap["accesslist"])
// 	}

// 	if blacklist, ok := entryMap["blacklist"].([]interface{}); ok {
// 		accessControls.Blacklist = blacklist
// 	} else {
// 		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["blacklist"], entryMap["blacklist"])
// 	}

// 	if whitelist, ok := entryMap["whitelist"].([]interface{}); ok {
// 		accessControls.Whitelist = whitelist
// 	} else {
// 		return nil, fmt.Errorf("%v was not a []interface{}, actual: %T", entryMap["whitelist"], entryMap["whitelist"])
// 	}

// 	return accessControls, nil
// }
