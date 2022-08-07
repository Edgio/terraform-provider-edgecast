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

func DataSourceDomainStatuses() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceDomainStatusesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relative path to an endpoint through which you may retrieve a list of domain statuses.",
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
							Description: "Indicates the name of a domain status.",
							Optional:    true,
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the system-defined ID for a domain status.",
							Optional:    true,
						},
					},
				},
				Description: "Contains a list of domain statuses.",
			},

			"total_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the total number of domain statuses.",
			},
		},
	}
}

func DataSourceDomainStatusesRead(
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

	// Call Get Appendix Domain Statuses API
	params := appendix.NewAppendixGetDomainStatusesParams()
	domainStatusesObj, err := cpsService.Appendix.AppendixGetDomainStatuses(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Domain Statuses: %# v", pretty.Formatter(domainStatusesObj))

	d.SetId(domainStatusesObj.AtID)
	d.SetType(domainStatusesObj.AtType)
	d.Set("total_items", domainStatusesObj.TotalItems)

	flattened := make([]map[string]interface{}, int(domainStatusesObj.TotalItems))
	for key := range domainStatusesObj.Items {
		cc := make(map[string]interface{})
		cc["id"] = domainStatusesObj.Items[key].ID
		cc["name"] = domainStatusesObj.Items[key].Name
		flattened[key] = cc
	}
	d.Set("items", flattened)

	return diag.Diagnostics{}
}
