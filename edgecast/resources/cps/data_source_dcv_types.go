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

func DataSourceDCVTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceDCVTypesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relative path to an endpoint through which you may retrieve a list of DCV types.",
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
							Description: "Indicates a type of DCV.",
							Optional:    true,
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the system-defined ID of a DCV type.",
							Optional:    true,
						},
					},
				},
				Description: "Contains a list of DCV types.",
			},

			"total_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the total number of DCV types.",
			},
		},
	}
}

func DataSourceDCVTypesRead(
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

	// Call Get Appendix DCV Types API
	params := appendix.NewAppendixGetDcvTypesParams()
	dcvTypesObj, err := cpsService.Appendix.AppendixGetDcvTypes(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved DCV Types: %# v", pretty.Formatter(dcvTypesObj))

	d.SetId(dcvTypesObj.AtID)
	d.SetType(dcvTypesObj.AtType)
	d.Set("total_items", dcvTypesObj.TotalItems)

	flattened := make([]map[string]interface{}, int(dcvTypesObj.TotalItems))
	for key := range dcvTypesObj.Items {
		cc := make(map[string]interface{})
		cc["id"] = dcvTypesObj.Items[key].ID
		cc["name"] = dcvTypesObj.Items[key].Name
		flattened[key] = cc
	}
	d.Set("items", flattened)

	return diag.Diagnostics{}
}
