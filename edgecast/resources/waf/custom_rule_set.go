// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"context"
	"errors"
	"fmt"
	"log"

	"terraform-provider-edgecast/edgecast/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
)

func ResourceCustomRuleSetCreate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("customer_id").(string)

	log.Printf("[INFO] Creating WAF Rate Rule for Account >> %s", accountNumber)

	customRuleSet := sdkwaf.CustomRuleSet{
		Name: d.Get("name").(string),
	}

	directive, err := expandCustomRuleDirectives(d.Get("directive"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing directive: %w", err))
	}

	customRuleSet.Directives = *directive

	log.Printf("[DEBUG] Name: %+v\n", customRuleSet.Name)
	log.Printf("[DEBUG] Directive(s): %+v\n", customRuleSet.Directives)

	params := sdkwaf.NewAddCustomRuleSetParams()
	params.AccountNumber = accountNumber
	params.CustomRuleSet = customRuleSet
	resp, err := wafService.AddCustomRuleSet(params)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] %+v", resp)

	d.SetId(resp)

	return ResourceCustomRuleSetRead(ctx, d, m)
}

func ResourceCustomRuleSetRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Get("customer_id").(string)
	ruleID := d.Id()

	log.Printf("[INFO] Retrieving custom rule %s for account number %s",
		ruleID,
		accountNumber,
	)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	params := sdkwaf.NewGetCustomRuleSetParams()
	params.AccountNumber = accountNumber
	params.CustomRuleSetID = ruleID
	resp, err := wafService.GetCustomRuleSet(params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully retrieved rate rule %s: %+v", ruleID, resp)

	d.SetId(resp.ID)
	d.Set("customer_id", accountNumber)
	d.Set("last_modified_date", resp.LastModifiedDate)
	d.Set("name", resp.Name)

	flattenDirectiveGroups := flattenCustomRuleDirectives(resp.Directives)

	d.Set("directive", flattenDirectiveGroups)
	return diags
}

func ResourceCustomRuleSetUpdate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	accountNumber := d.Get("customer_id").(string)
	customRuleSetID := d.Id()

	log.Printf("[INFO] Updating WAF Custom Rule Set ID %s for Account >> %s",
		customRuleSetID,
		accountNumber,
	)

	customRuleSetRequest := sdkwaf.CustomRuleSet{}
	customRuleSetRequest.Name = d.Get("name").(string)

	directives, err := expandCustomRuleDirectives(d.Get("directive"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing directive: %w", err))
	}
	customRuleSetRequest.Directives = *directives

	log.Printf("[DEBUG] Name: %+v\n", customRuleSetRequest.Name)
	log.Printf("[DEBUG] Directives: %+v\n", customRuleSetRequest.Directives)

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	params := sdkwaf.NewUpdateCustomRuleSetParams()
	params.AccountNumber = accountNumber
	params.CustomRuleSet = customRuleSetRequest
	params.CustomRuleSetID = customRuleSetID
	err = wafService.UpdateCustomRuleSet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Successfully updated WAF Custom Rule Set: %+v",
		customRuleSetRequest)

	return ResourceCustomRuleSetRead(ctx, d, m)
}

func ResourceCustomRuleSetDelete(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics

	accountNumber := d.Get("customer_id").(string)
	customRuleID := d.Id()

	log.Printf("[INFO] Deleting WAF Custom Rule Set ID %s for Account >> %s",
		customRuleID,
		accountNumber,
	)

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	params := sdkwaf.NewDeleteCustomRuleSetParams()
	params.AccountNumber = accountNumber
	params.CustomRuleSetID = customRuleID
	err = wafService.DeleteCustomRuleSet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Successfully deleted WAF Custom Rule Set: %+v", customRuleID)

	d.SetId("")

	return diags
}

func expandCustomRuleDirectives(
	attr interface{},
) (*[]sdkwaf.CustomRuleDirective, error) {

	if set, ok := attr.(*schema.Set); ok {

		items := set.List()
		directives := make([]sdkwaf.CustomRuleDirective, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			directive := sdkwaf.CustomRuleDirective{}

			secRule, err := expandSecRule(curr["sec_rule"])
			if err != nil {
				return nil, err
			}

			directive.SecRule = *secRule

			directives = append(directives, directive)
		}

		return &directives, nil

	} else {
		return nil,
			errors.New(
				"expandCustomRuleSetDirectives: input was not a *schema.Set")
	}
}

// flattenCustomRuleDirectives converts the CustomRuleDirective API Model
// into a format that Terraform can work with
func flattenCustomRuleDirectives(
	directive []sdkwaf.CustomRuleDirective,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, d := range directive {
		m := make(map[string]interface{})
		m["sec_rule"] = flattenSecRule(d.SecRule)
		flattened = append(flattened, m)
	}

	return flattened
}
