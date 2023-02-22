// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf_bot_manager

import (
	"context"
	"log"

	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf_bot_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBotManagerCreate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	/*
		config := m.(internal.ProviderConfig)

		botManagerService, err := buildBotManagerService(config)

		if err != nil {
			return diag.FromErr(err)
		}

		customerID := d.Get("customer_id").(string)

		log.Printf("[INFO] Creating Bot Manager for Account >> %s", customerID)

		botManager := botmanager.BotManager{
			//Name:       d.Get("name").(string),
			CustomerId: &customerID,
			//BotsProdId: d.Get("bots_prod_id").(string),

		}
	*/
	return ResourceBotManagerRead(ctx, d, m)
}

func ResourceBotManagerRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceBotManagerUpdate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return ResourceBotManagerRead(ctx, d, m)
}

func ResourceBotManagerDelete(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	customerID := d.Get("customer_id").(string)
	botManagerID := d.Id()

	log.Printf("[INFO] Deleting WAF Bot Manager ID %s for Account >> %s",
		botManagerID,
		customerID,
	)

	config := m.(internal.ProviderConfig)

	svc, err := buildBotManagerService(config)
	if err != nil {
		return diag.FromErr(err)
	}

	params := waf_bot_manager.NewDeleteBotManagerParams()
	params.CustId = customerID
	params.BotManagerId = botManagerID

	err = svc.BotManagers.DeleteBotManager(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully deleted WAF Bot Manager: %s", botManagerID)
	d.SetId("")

	return diags
}
