// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package originv3

import (
	"context"
	"log"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/originv3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

func DataSourceOriginShieldPops() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOriginShieldPopsRead,
		Schema: map[string]*schema.Schema{
			"code": {
				Type: schema.TypeString,
				Description: "Replace this variable with either of the following codes:  \n" +
					" --> POP Code: Pass an Origin Shield's POP code to retrieve information about that specific Origin Shield POP. \n" +
					" --> Bypass Region Code: Pass one of the following bypass code to retrieve information about that region including all of its Origin Shield POPs: \n" +
					"     - BYAP: Asia \n" +
					"     - BYEC: US East \n" +
					"     - BYEU: Europe \n" +
					"     - BYWC: US West \n" +
					"\n" +
					"",
				Optional: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bypass_code": {
							Type:        schema.TypeString,
							Description: "Indicates the four-letter abbreviation that corresponds to this region's bypass code.",
							Computed:    true,
						},
						"bypass_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the name of this region's bypass configuration.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the region's system-defined ID.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the name of the current region.",
						},
						"pops": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicates the system-defined ID assigned to this Origin Shield POP.",
									},
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the three-letter abbreviation corresponding to this POP location.",
									},
									"city": {
										Type:        schema.TypeString,
										Description: "Indicates the Origin Shield's city.",
										Computed:    true,
									},
									"is_pci_certified": {
										Type:        schema.TypeBool,
										Description: "Indicates whether this customer origin group is restricted to Payment Card Industry (PCI)-compliant Origin Shield POPs. Valid values are: true|false",
										Computed:    true,
									},
								},
							},
							Description: "Contains a list of countries.",
						},
					},
				},
			},
		},
	}
}

func DataSourceOriginShieldPopsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config, ok := m.(internal.ProviderConfig)
	if !ok {
		return diag.Errorf("failed to load configuration")
	}

	svc, err := buildOriginV3Service(config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Get Origin Shield Pops
	params := originv3.NewGetOriginShieldPopsParams()

	if attr, ok := d.GetOk("code"); ok {
		code := attr.(string)
		params.Findcode = code
	}

	resp, err := svc.HttpLargeOnly.GetOriginShieldPops(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Retrieved Origin Shield Pops: %# v",
		pretty.Formatter(resp))

	flattened := flattenOriginShieldEdgeNodes(resp)
	d.Set("items", flattened)

	// always run
	d.SetId(helper.GetUnixTimeStamp())

	return diag.Diagnostics{}
}

func flattenOriginShieldEdgeNodes(
	shieldEdgeNodes []originv3.OriginShieldEdgeNode,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range shieldEdgeNodes {
		m := make(map[string]interface{})

		m["bypass_code"] = v.BypassCode
		m["bypass_name"] = v.BypassName
		m["region_id"] = v.RegionId
		m["region_name"] = v.RegionName
		m["pops"] = flattenOriginShieldPops(v.Pops)

		flattened = append(flattened, m)
	}

	return flattened
}

func flattenOriginShieldPops(
	shieldPops []originv3.OriginShieldPop,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range shieldPops {
		m := make(map[string]interface{})

		m["id"] = v.Id
		m["code"] = v.Code
		m["city"] = v.City
		m["is_pci_certified"] = v.IsPciCertified

		flattened = append(flattened, m)
	}

	return flattened
}
