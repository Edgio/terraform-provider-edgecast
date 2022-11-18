// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package originv3

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/originv3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

func ResourceOriginV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOriginV3Create,
		ReadContext:   ResourceOriginV3Read,
		UpdateContext: ResourceOriginV3Update,
		DeleteContext: ResourceOriginV3Delete,
		Importer:      helper.Import(ResourceOriginV3Read, "id"),
		Schema:        GetOriginV3Schema(),
	}
}

func ResourceOriginV3Create(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	// Initialize Servicee
	config, ok := m.(internal.ProviderConfig)
	if !ok {
		return helper.CreationErrorf(d, "failed to load configuration")
	}

	svc, err := buildOriginV3Service(config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Read from TF state.
	originGroupState, errs := expandOriginGroup(d)
	if len(errs) > 0 {
		return helper.DiagsFromErrors("error parsing origin group", errs)
	}

	// Call APIs.
	cparams := originv3.NewAddHttpLargeGroupParams()
	cparams.CustomerOriginGroupHTTPRequest =
		originv3.CustomerOriginGroupHTTPRequest{
			Name:               originGroupState.Name,
			HostHeader:         originGroupState.HostHeader,
			ShieldPops:         originGroupState.ShieldPops,
			NetworkTypeId:      originGroupState.NetworkTypeId,
			StrictPciCertified: originGroupState.StrictPciCertified,
			TlsSettings:        originGroupState.TlsSettings,
		}

	cresp, err := svc.HttpLargeOnly.AddHttpLargeGroup(cparams)
	if err != nil {
		return helper.CreationError(d, err)
	}
	log.Printf("[INFO] origin group created: %# v\n", pretty.Formatter(cresp))
	log.Printf("[INFO] origin group id: %d\n", cresp.Id)

	d.SetId(strconv.Itoa(int(*cresp.Id)))

	return diag.Diagnostics{}
}

func ResourceOriginV3Read(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	//Not implemented
	return diag.Diagnostics{}
}

func ResourceOriginV3Update(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	//Not implemented
	return diag.Diagnostics{}
}

func ResourceOriginV3Delete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	//Not implemented
	return diag.Diagnostics{}
}

func expandOriginGroup(
	d *schema.ResourceData,
) (*originv3.CustomerOriginGroupHTTPRequest, []error) {
	if d == nil {
		return nil, []error{errors.New("no data to read")}
	}

	errs := make([]error, 0)

	originGrpState := &originv3.CustomerOriginGroupHTTPRequest{}

	if v, ok := d.GetOk("name"); ok {
		if name, ok := v.(string); ok {
			originGrpState.Name = name
		} else {
			errs = append(errs, errors.New("name not a string"))
		}
	}

	if v, ok := d.GetOk("host_header"); ok {
		if hostHeader, ok := v.(string); ok {
			originGrpState.SetHostHeader(hostHeader)
		} else {
			errs = append(errs, errors.New("host_header not a string"))
		}
	}

	if v, ok := d.GetOk("network_type_id"); ok {
		if networkTypeID, ok := v.(int); ok {
			originGrpState.SetNetworkTypeId(int32(networkTypeID))
		} else {
			errs = append(errs, errors.New("network_type_id not a int32"))
		}
	}

	if v, ok := d.GetOk("strict_pci_certified"); ok {
		if strictPCICertified, ok := v.(bool); ok {
			originGrpState.SetStrictPciCertified(strictPCICertified)
		} else {
			errs = append(errs, errors.New("strict_pci_certified not a bool"))
		}
	}

	if v, ok := d.GetOk("shield_pops"); ok {
		shieldPOPs, err := helper.ConvertTFCollectionToPtrStrings(v)
		if err == nil {
			originGrpState.ShieldPops = shieldPOPs
		} else {
			errs = append(errs, fmt.Errorf("error parsing shield_pops: %w", err))
		}
	}

	if v, ok := d.GetOk("tls_settings"); ok {
		if tlsSettings, err := expandTLSSettings(v); err == nil {
			originGrpState.TlsSettings = tlsSettings
		} else {
			errs = append(errs, fmt.Errorf("error parsing tls_settings: %w", err))
		}
	}

	return originGrpState, errs
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
