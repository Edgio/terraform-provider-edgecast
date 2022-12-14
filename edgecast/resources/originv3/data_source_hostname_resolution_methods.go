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

func DataSourceHostnameResolutionMethods() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceHostnameResolutionMethodsRead,
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the hostname resolution method's system-defined ID.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies a hostname resolution method by name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the hostname resolution method's code name.",
						},
					},
				},
			},
		},
	}
}

func DataSourceHostnameResolutionMethodsRead(
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

	// Call Get Available Hostname Resolution Methods
	resp, err := svc.Phase3.GetAvailableHostnameResolutionMethods()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Retrieved Available Hostname Resolution Methods: %# v",
		pretty.Formatter(resp))

	flattened := flattenHostnameResolutionMethods(resp)
	d.Set("items", flattened)

	// always run
	d.SetId(helper.GetUnixTimeStamp())

	return diag.Diagnostics{}
}

func flattenHostnameResolutionMethods(
	networkTypes []originv3.NetworkType,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range networkTypes {
		m := make(map[string]interface{})

		m["id"] = v.Id
		m["name"] = v.Name
		m["value"] = v.Value
		flattened = append(flattened, m)
	}

	return flattened
}
