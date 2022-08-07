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

func DataSourceCertReqStatuses() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceCertReqStatusesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relative path to an endpoint through which you may retrieve a list of certificate request statuses.",
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
							Description: "Indicates the name of a certificate request status.",
							Optional:    true,
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the system-defined ID for a certificate request status.",
							Optional:    true,
						},
					},
				},
				Description: "Contains a list of certificate request statuses.",
			},

			"total_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the total number of certificate request statuses.",
			},
		},
	}
}

func DataSourceCertReqStatusesRead(
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

	// Call Get Appendix certificate request statuses API
	params := appendix.NewAppendixGetCertificateStatusesParams()
	certReqStatusesObj, err := cpsService.Appendix.AppendixGetCertificateStatuses(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Certificate Request Statuses: %# v", pretty.Formatter(certReqStatusesObj))

	d.SetId(certReqStatusesObj.AtID)
	d.SetType(certReqStatusesObj.AtType)
	d.Set("total_items", certReqStatusesObj.TotalItems)

	flattened := make([]map[string]interface{}, int(certReqStatusesObj.TotalItems))
	for key := range certReqStatusesObj.Items {
		cc := make(map[string]interface{})
		cc["id"] = certReqStatusesObj.Items[key].ID
		cc["name"] = certReqStatusesObj.Items[key].Name
		flattened[key] = cc
	}
	d.Set("items", flattened)

	return diag.Diagnostics{}
}
