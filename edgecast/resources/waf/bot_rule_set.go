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

func ResourceBotRuleSetCreate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	config := m.(**api.ClientConfig)
	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("customer_id").(string)

	log.Printf(
		"[INFO] Creating WAF Bot Rule Set for Account >> %s",
		accountNumber)

	botRuleSet := sdkwaf.BotRuleSet{
		Name: d.Get("name").(string),
	}

	directive, err := expandBotRuleDirectives(d.Get("directive"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing directive: %w", err))
	}

	botRuleSet.Directives = *directive

	log.Printf("[DEBUG] Name: %+v\n", botRuleSet.Name)
	log.Printf("[DEBUG] Directive(s): %+v\n", botRuleSet.Directives)

	params := sdkwaf.NewAddBotRuleSetParams()
	params.AccountNumber = accountNumber
	params.BotRuleSet = botRuleSet
	resp, err := wafService.AddBotRuleSet(params)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] %+v", resp)

	d.SetId(resp)

	return ResourceBotRuleSetRead(ctx, d, m)
}

func ResourceBotRuleSetRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Get("customer_id").(string)
	botRuleSetID := d.Id()

	log.Printf("[INFO] Retrieving Bot Rule Set '%s' for account number %s",
		botRuleSetID,
		accountNumber,
	)

	wafService, err := buildWAFService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	params := sdkwaf.NewGetBotRuleSetParams()
	params.AccountNumber = accountNumber
	params.BotRuleSetID = botRuleSetID
	resp, err := wafService.GetBotRuleSet(params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Successfully retrieved Bot Rule Set '%s': %+v",
		botRuleSetID,
		resp)

	d.SetId(resp.ID)
	d.Set("customer_id", accountNumber)
	d.Set("last_modified_date", resp.LastModifiedDate)
	d.Set("name", resp.Name)

	flattenedDirectives := flattenBotRuleDirectives(resp.Directives)

	d.Set("directive", flattenedDirectives)
	return diags
}

func ResourceBotRuleSetUpdate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	accountNumber := d.Get("customer_id").(string)
	botRuleSetID := d.Id()

	log.Printf("[INFO] Updating WAF Bot Rule '%s' for Account >> %s",
		botRuleSetID,
		accountNumber,
	)

	botRuleSet := sdkwaf.BotRuleSet{}
	botRuleSet.Name = d.Get("name").(string)

	directives, err := expandBotRuleDirectives(d.Get("directive"))
	if err != nil {
		return diag.Errorf("error parsing directives: %+v", err)
	}
	botRuleSet.Directives = *directives

	log.Printf("[DEBUG] Name: %+v\n", botRuleSet.Name)
	log.Printf("[DEBUG] Directives: %+v\n", botRuleSet.Directives)

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	params := sdkwaf.NewUpdateBotRuleSetParams()
	params.AccountNumber = accountNumber
	params.BotRuleSet = botRuleSet
	params.BotRuleSetID = botRuleSetID
	err = wafService.UpdateBotRuleSet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Successfully updated WAF Bot Rule Set: %+v",
		botRuleSet)

	return ResourceBotRuleSetRead(ctx, d, m)
}

func ResourceBotRuleSetDelete(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics

	accountNumber := d.Get("customer_id").(string)
	botRuleSetID := d.Id()

	log.Printf("[INFO] Deleting WAF Bot Rule Set ID %s for Account >> %s",
		botRuleSetID,
		accountNumber,
	)

	config := m.(**api.ClientConfig)

	wafService, err := buildWAFService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	params := sdkwaf.NewDeleteBotRuleSetParams()
	params.AccountNumber = accountNumber
	params.BotRuleSetID = botRuleSetID
	err = wafService.DeleteBotRuleSet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Successfully deleted WAF Bot Rule Set: %+v", botRuleSetID)

	d.SetId("")

	return diags
}

func expandBotRuleDirectives(
	attr interface{},
) (*[]sdkwaf.BotRuleDirective, error) {

	if set, ok := attr.(*schema.Set); ok {

		items := set.List()
		directives := make([]sdkwaf.BotRuleDirective, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			directive := sdkwaf.BotRuleDirective{}

			if secRuleRaw, ok := curr["sec_rule"]; ok {

				secRule, err := expandSecRule(secRuleRaw)

				if err != nil {
					return nil, err
				}

				directive.SecRule = secRule
			}

			if include, ok := curr["include"].(string); ok {
				directive.Include = include
			}

			directives = append(directives, directive)
		}

		return &directives, nil

	} else {
		return nil,
			errors.New(
				"expandCustomRuleSetDirectives: input was not a *schema.Set")
	}
}

// flattenBotRuleDirectives converts the BotRuleDirective API Model
// into a format that Terraform can work with
func flattenBotRuleDirectives(
	directive []sdkwaf.BotRuleDirective,
) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)

	for _, d := range directive {
		m := make(map[string]interface{})

		if d.SecRule != nil {
			m["sec_rule"] = flattenSecRule(*d.SecRule)
		}

		if d.Include != "" {
			m["include"] = d.Include
		}

		flattened = append(flattened, m)
	}

	return flattened
}
