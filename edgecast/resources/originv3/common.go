// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package originv3

import (
	"errors"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/originv3"
)

// buildOriginV3Service builds the SDK OriginV3 service to manage Origin Groups
// and Origin resources.
func buildOriginV3Service(
	config internal.ProviderConfig,
) (*originv3.Service, error) {
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

	return originv3.New(sdkConfig)
}

// expandTLSSettings converts the Terraform representation of TLS Settings
// into the TLSSetting API Model.
func expandTLSSettings(attr interface{}) (*originv3.TlsSettings, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}
	if len(raw) > 1 {
		return nil, errors.New("only one tls setting is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	tls := originv3.TlsSettings{}

	if sniHostName, ok := curr["sni_hostname"].(string); ok {
		tls.SetSniHostname(sniHostName)
	}

	if allowSelfSigned, ok := curr["allow_self_signed"].(bool); ok {
		tls.AllowSelfSigned = &allowSelfSigned
	}

	if v, ok := curr["public_keys_to_verify"]; ok {
		keys, err := helper.ConvertTFCollectionToStrings(v)
		if err == nil {
			tls.PublicKeysToVerify = keys
		}
	}

	return &tls, nil
}

// expandOrigins converts the Terraform representation of Origins
// into the OriginState Model.
func expandOrigins(attr interface{}) ([]*OriginState, error) {
	if items, ok := attr.([]interface{}); ok {
		origins := make([]*OriginState, 0)

		for _, item := range items {
			curr, ok := item.(map[string]interface{})
			if !ok {
				return nil, errors.New("origin was not map[string]interface{}")
			}

			id, ok := curr["id"].(int)
			if !ok {
				id = 0
			}

			host, ok := curr["host"].(string)
			if !ok {
				return nil, errors.New("origin.host was not a string")
			}

			name, ok := curr["name"].(string)
			if !ok {
				return nil, errors.New("origin.name was not a string")
			}

			port, ok := curr["port"].(int)
			if !ok {
				return nil, errors.New("origin.port was not an int")
			}
			port32 := int32(port)

			isPrimary, ok := curr["is_primary"].(bool)
			if !ok {
				return nil, errors.New("domain.is_primary was not a bool")
			}

			storageTypeID, ok := curr["storage_type_id"].(int)
			if !ok {
				return nil, errors.New("origin.storage_type_id was not an int")
			}

			protocolTypeID, ok := curr["protocol_type_id"].(int)
			if !ok {
				return nil, errors.New("origin.protocol_type_id was not an int")
			}

			failoverOrder, ok := curr["failover_order"].(int)
			if !ok {
				return nil, errors.New("origin.failover_order was not an int")
			}

			origin := OriginState{
				ID:             int32(id),
				Name:           name,
				Host:           host,
				Port:           port32,
				IsPrimary:      isPrimary,
				StorageTypeID:  int32(storageTypeID),
				ProtocolTypeID: int32(protocolTypeID),
				FailoverOrder: int32(failoverOrder),
			}

			origins = append(origins, &origin)
		}

		return origins, nil
	} else {
		return nil,
			errors.New("ExpandOrigins: attr input was not a []interface{}")
	}
}

func flattenTLSSettings(
	settings *originv3.TlsSettings,
) []map[string]interface{} {
	if settings == nil {
		return make([]map[string]interface{}, 0)
	}

	flattened := make([]map[string]interface{}, 0)

	m := make(map[string]interface{})

	m["allow_self_signed"] = settings.GetAllowSelfSigned()
	m["sni_hostname"] = settings.SniHostname.Get()
	if len(settings.PublicKeysToVerify) > 0 {
		m["public_keys_to_verify"] = settings.PublicKeysToVerify
	} else {
		m["public_keys_to_verify"] = make([]string, 0)
	}

	flattened = append(flattened, m)
	return flattened
}

func flattenOrigins(
	origins []originv3.CustomerOriginFailoverOrder,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range origins {
		m := make(map[string]interface{})

		m["id"] = int(v.GetId())
		m["host"] = v.GetHost()
		m["is_primary"] = v.GetIsPrimary()
		m["name"] = v.GetName()
		m["port"] = int(v.GetPort())
		m["protocol_type_id"] = int(v.GetProtocolTypeId())
		m["storage_type_id"] = int(v.GetStorageTypeId())
		m["failover_order"] = int(v.GetFailoverOrder())

		flattened = append(flattened, m)
	}
	return flattened
}
