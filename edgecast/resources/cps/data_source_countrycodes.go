// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"context"
	"log"

	"terraform-provider-edgecast/edgecast/api"
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/appendix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

func DataSourceCountryCodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceCountryCodesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The unix timestamp when the data source was refreshed.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a specific country. If provided, only that country will be present in `items`.",
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country": {
							Type:        schema.TypeString,
							Description: "Identifies a country by its name.",
							Computed:    true,
						},
						"two_letter_code": {
							Type:        schema.TypeString,
							Description: "Identifies a country by its two-letter country code.",
							Computed:    true,
						},
					},
				},
				Description: "Contains a list of countries.",
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

	// Call Get Country Codes
	params := appendix.NewAppendixGetParams()

	if attr, ok := d.GetOk("name"); ok {
		name := attr.(string)
		params.Name = &name
	}

	resp, err := cpsService.Appendix.AppendixGet(params)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Retrieved Country Codes: %# v",
		pretty.Formatter(resp))

	flattened := FlattenCountries(resp)
	d.Set("items", flattened)

	// always run
	d.SetId(helper.GetUnixTimeStamp())

	return diag.Diagnostics{}
}

func FlattenCountries(
	countries *appendix.AppendixGetOK,
) []map[string]interface{} {
	if countries != nil {
		flattened := make([]map[string]interface{}, len(countries.Items), len(countries.Items))

		for ix := range countries.Items {
			cc := make(map[string]interface{})
			cc["country"] = countries.Items[ix].Country
			cc["two_letter_code"] = countries.Items[ix].TwoLetterCode
			flattened[ix] = cc
		}

		return flattened
	}

	return make([]map[string]interface{}, 0)
}
