// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf_bot_manager

import (
	"context"

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

	return diags
}
