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

func DataSourceCertValidationLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceCertValidationLevelsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relative path to an endpoint through which you may retrieve a list of certificate validation levels.",
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
							Description: "Indicates a certificate validation level.",
							Optional:    true,
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the system-defined ID for a certificate validation level.",
							Optional:    true,
						},
					},
				},
				Description: "Contains a list of certificate validation levels.",
			},

			"total_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the total number of certificate validation levels.",
			},
		},
	}
}

func DataSourceCertValidationLevelsRead(
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

	// Call Get Appendix Validation Levels API
	params := appendix.NewAppendixGetValidationTypesParams()
	validationLevelsObj, err := cpsService.Appendix.AppendixGetValidationTypes(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Validation Levels: %# v", pretty.Formatter(validationLevelsObj))

	d.SetId(validationLevelsObj.AtID)
	d.SetType(validationLevelsObj.AtType)
	d.Set("total_items", validationLevelsObj.TotalItems)

	flattened := make([]map[string]interface{}, int(validationLevelsObj.TotalItems))
	for key := range validationLevelsObj.Items {
		cc := make(map[string]interface{})
		cc["id"] = validationLevelsObj.Items[key].ID
		cc["name"] = validationLevelsObj.Items[key].Name
		flattened[key] = cc
	}
	d.Set("items", flattened)

	return diag.Diagnostics{}
}
