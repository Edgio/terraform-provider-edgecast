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

func DataSourceCancelCertReqActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceCancelCertReqActionsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relative path to an endpoint through which you may retrieve a list of actions.",
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
							Description: "Indicates the name of action.",
							Optional:    true,
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the system-defined ID for an action.",
							Optional:    true,
						},
					},
				},
				Description: "Contains a list of actions.",
			},

			"total_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the total number of actions.",
			},
		},
	}
}

func DataSourceCancelCertReqActionsRead(
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

	// Call Get Appendix Cancel Request Actions API
	params := appendix.NewAppendixGetCancelActionsParams()
	cancelActionsObj, err := cpsService.Appendix.AppendixGetCancelActions(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Cancel Actions: %# v", pretty.Formatter(cancelActionsObj))

	d.SetId(cancelActionsObj.AtID)
	d.SetType(cancelActionsObj.AtType)
	d.Set("total_items", cancelActionsObj.TotalItems)

	flattened := make([]map[string]interface{}, int(cancelActionsObj.TotalItems))
	for key := range cancelActionsObj.Items {
		cc := make(map[string]interface{})
		cc["id"] = cancelActionsObj.Items[key].ID
		cc["name"] = cancelActionsObj.Items[key].Name
		flattened[key] = cc
	}
	d.Set("items", flattened)

	return diag.Diagnostics{}
}
