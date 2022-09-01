// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"fmt"
	"strconv"
	"time"

	"terraform-provider-edgecast/edgecast/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// buildCPSService builds the SDK CPS service to manage CPS
// resources
func buildCPSService(
	config api.ClientConfig,
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

	return cps.New(sdkConfig)
}

func FlattenNamedEntities(
	entities models.HyperionCollectionNamedEntity,
) []map[string]interface{} {
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
			Description: "The unix timestamp when the data source was refreshed.",
			Computed:    true,
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: fmt.Sprintf("Indicates the name of the %s.", resource),
						Computed:    true,
					},
					"id": {
						Type:        schema.TypeInt,
						Description: fmt.Sprintf("Indicates the system-defined ID for the %s.", resource),
						Computed:    true,
					},
				},
			},
			Description: fmt.Sprintf("Contains a list of %s", resource),
		},
	}
}

func getTimeStamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
