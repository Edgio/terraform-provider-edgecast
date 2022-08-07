// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"context"
	"log"

	"terraform-provider-edgecast/edgecast/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/appendix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

func DataSourceValidationStatuses() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceValidationStatusesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relative path to an endpoint through which you may retrieve a list of validation statuses.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returns 'Collection'.",
			},
			"items": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Indicates the name of a validation status.",
							Optional:    true,
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the system-defined ID for a validation status.",
							Optional:    true,
						},
					},
				},
				Description: "Contains a list of validation statuses.",
			},

			"total_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the total number of validation statuses.",
			},
		},
	}
}

func DataSourceValidationStatusesRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize CPS Service
	config := m.(**api.ClientConfig)
	cpsService, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Get Appendix Validation statuses API
	params := appendix.NewAppendixGetValidationStatusesParams()
	validationStatusesObj, err := cpsService.Appendix.AppendixGetValidationStatuses(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Validation statuses: %# v", pretty.Formatter(validationStatusesObj))

	d.SetId(validationStatusesObj.AtID)
	d.SetType(validationStatusesObj.AtType)
	d.Set("total_items", validationStatusesObj.TotalItems)

	flattened := make([]map[string]interface{}, int(validationStatusesObj.TotalItems))
	for key := range validationStatusesObj.Items {
		cc := make(map[string]interface{})
		cc["id"] = validationStatusesObj.Items[key].ID
		cc["name"] = validationStatusesObj.Items[key].Name
		flattened[key] = cc
	}
	d.Set("items", flattened)

	return diag.Diagnostics{}
}
