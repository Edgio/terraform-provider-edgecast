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

func DataSourceCountryCodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceCountryCodesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query by name parameter",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relative path to an endpoint through which you may retrieve a list of countries.",
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
						"country": {
							Type:        schema.TypeString,
							Description: "Identifies a country by its name.",
							Optional:    true,
						},
						"two_letter_code": {
							Type:        schema.TypeString,
							Description: "Identifies a country by its country code.",
							Optional:    true,
						},
					},
				},
				Description: "Contains a list of countries.",
			},

			"total_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the total number of countries.",
			},
		},
	}
}

func DataSourceCountryCodesRead(
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

	// Call Get Appendix API
	params := appendix.NewAppendixGetParams()
	query := d.Get("name").(string)
	params.Name = &query
	countryCodeObj, err := cpsService.Appendix.AppendixGet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Country Codes: %# v", pretty.Formatter(countryCodeObj))

	d.SetId(countryCodeObj.AtID)
	d.SetType(countryCodeObj.AtType)
	d.Set("total_items", countryCodeObj.TotalItems)

	flattened := make([]map[string]interface{}, int(countryCodeObj.TotalItems))
	for key := range countryCodeObj.Items {
		cc := make(map[string]interface{})
		cc["country"] = countryCodeObj.Items[key].Country
		cc["two_letter_code"] = countryCodeObj.Items[key].TwoLetterCode
		flattened[key] = cc
	}
	d.Set("items", flattened)

	return diag.Diagnostics{}
}
