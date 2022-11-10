// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

// buildCPSService builds the SDK CPS service to manage CPS
// resources.
func buildCPSService(
	config internal.ProviderConfig,
) (*cps.CpsService, error) {
	idsCredentials := edgecast.IDSCredentials{
		ClientID:     config.IdsClientID,
		ClientSecret: config.IdsClientSecret,
		Scope:        config.IdsScope,
	}

	sdkConfig := edgecast.NewSDKConfig()
	sdkConfig.APIToken = config.APIToken
	sdkConfig.IDSCredentials = idsCredentials
	sdkConfig.BaseAPIURL = *config.APIURL
	sdkConfig.BaseAPIURLLegacy = *config.APIURLLegacy
	sdkConfig.BaseIDSURL = *config.IdsURL
	sdkConfig.UserAgent = config.UserAgent

	return cps.New(sdkConfig)
}

func DataSourceNamedEntityRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	readFunc func(svc *cps.CpsService, d *schema.ResourceData) (*models.HyperionCollectionNamedEntity, error),
) diag.Diagnostics {
	// Initialize CPS Service
	config := m.(internal.ProviderConfig)

	cpsService, err := buildCPSService(config)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := readFunc(cpsService, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved: %# v\n", pretty.Formatter(resp))

	flattened := FlattenNamedEntities(resp)
	d.Set("items", flattened)

	// always run
	d.SetId(helper.GetUnixTimeStamp())

	return diag.Diagnostics{}
}

func FlattenNamedEntities(
	entities *models.HyperionCollectionNamedEntity,
) []map[string]interface{} {
	if entities == nil {
		return make([]map[string]interface{}, 0)
	}

	flattened := make([]map[string]interface{}, len(entities.Items))

	for ix := range entities.Items {
		cc := make(map[string]interface{})
		cc["id"] = entities.Items[ix].ID
		cc["name"] = entities.Items[ix].Name
		flattened[ix] = cc
	}

	return flattened
}

func namedEntitySchema(resource string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "Indicates the Unix timestamp at which the data source was refreshed.",
			Computed:    true,
		},
		"items": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: fmt.Sprintf("Indicates the name for a resource of type %s.", resource),
						Computed:    true,
					},
					"id": {
						Type:        schema.TypeInt,
						Description: fmt.Sprintf("Indicates the system-defined ID for a resource of type %s.", resource),
						Computed:    true,
					},
				},
			},
			Description: fmt.Sprintf("Contains a list of %s objects.", resource),
		},
	}
}
