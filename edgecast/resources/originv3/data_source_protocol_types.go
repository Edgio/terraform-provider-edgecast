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

func DataSourceProtocolTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceProtocolTypesRead,
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "Indicates the protocol's system-defined ID.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies a protocol by name.",
						},
					},
				},
			},
		},
	}
}

func DataSourceProtocolTypesRead(
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

	// Call Get Protocol Types
	resp, err := svc.Phase3.GetAvailableProtocols()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Retrieved Protocol Types: %# v",
		pretty.Formatter(resp))

	flattened := flattenProtocolTypes(resp)
	d.Set("items", flattened)

	// always run
	d.SetId(helper.GetUnixTimeStamp())

	return diag.Diagnostics{}
}

func flattenProtocolTypes(
	protocols []originv3.ProtocolType,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range protocols {
		m := make(map[string]interface{})

		m["id"] = v.Id
		m["name"] = v.Name
		flattened = append(flattened, m)
	}

	return flattened
}
