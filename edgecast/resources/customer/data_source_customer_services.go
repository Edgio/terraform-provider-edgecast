// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package customer

import (
	"context"
	"strconv"
	"terraform-provider-edgecast/edgecast/api"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCustomerServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceCustomerServicesRead,
		Schema: map[string]*schema.Schema{
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceCustomerServicesRead(
	ctx context.Context,
	d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)

	customerService, err := buildCustomerService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := customerService.GetAvailableCustomerServices()

	if err != nil {
		return diag.FromErr(err)
	}

	services := []interface{}{}

	for _, s := range *resp {
		service := map[string]interface{}{
			"id":        s.ID,
			"name":      s.Name,
			"parent_id": s.ParentID,
		}

		services = append(services, service)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("services", services); err != nil {
		return diag.FromErr(err)
	}

	// Terraform requires an ID - we will use a unix timestamp so that this is always pulled fresh
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
