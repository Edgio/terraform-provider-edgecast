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
	"github.com/EdgeCast/ec-sdk-go/edgecast/shared/enums"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

func ResourceOriginGrpHttpLarge() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOriginGroupCreate,
		ReadContext:   ResourceOriginGroupRead,
		UpdateContext: ResourceOriginGroupUpdate,
		DeleteContext: ResourceOriginGroupDelete,
		Importer:      helper.Import(ResourceOriginGroupRead, "id"),
		Schema:        GetOriginGrpHttpLargeSchema(),
	}
}

func ResourceOriginGroupCreate(
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

	originGroupState, originsState, errs := expandHttpLargeOriginGroup(d)
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

	grpID := cresp.Id

	if len(originsState) > 0 {
		for key, origin := range originsState {
			originsState[key].GroupId = *grpID

			params := originv3.NewAddOriginParams()
			params.MediaType = enums.HttpLarge.String()
			params.CustomerOriginRequest = *origin

			resp, err := svc.Common.AddOrigin(params)
			if err != nil {
				d.SetId("")

				//delete the created origin group
				deleteparams := originv3.NewDeleteGroupParams()
				deleteparams.GroupId = *grpID
				params.MediaType = enums.HttpLarge.String()

				deleteErr := svc.Common.DeleteGroup(deleteparams)
				if deleteErr != nil {
					return diag.Errorf(
						"failed to roll back origin group request upon error: %v, original err: %v",
						deleteErr.Error(),
						err.Error())
				}
				return diag.Errorf("failed to create origin group: %v", err)
			}

			log.Printf("[INFO] origin created: %# v\n", pretty.Formatter(resp))
		}
	}

	d.SetId(strconv.Itoa(int(*grpID)))

	return ResourceOriginGroupRead(ctx, d, m)
}

func ResourceOriginGroupRead(
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

	grpID, err := helper.ParseInt64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieving origin group : ID: %d\n", grpID)
	// call APIs
	params := originv3.NewGetHttpLargeGroupParams()
	params.GroupId = int32(grpID)

	resp, err := svc.HttpLargeOnly.GetHttpLargeGroup(params)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Retrieved origin group: %# v\n", pretty.Formatter(resp))

	// Write TF state.
	err = setHttpLargeOriginGroupState(d, resp)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func ResourceOriginGroupUpdate(
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

	grpID, err := helper.ParseInt64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	originGroupState, _, errs := expandHttpLargeOriginGroup(d)
	if len(errs) > 0 {
		return helper.DiagsFromErrors("error parsing origin group", errs)
	}

	log.Printf("[INFO] Updating origin group : ID: %d\n", grpID)
	updateParams := originv3.NewUpdateHttpLargeGroupParams()
	updateParams.GroupId = int32(grpID)
	updateParams.CustomerOriginGroupHTTPRequest = *originGroupState

	_, err = svc.HttpLargeOnly.UpdateHttpLargeGroup(updateParams)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updated origin group ID: %v", grpID)

	return ResourceOriginGroupRead(ctx, d, m)
}

func ResourceOriginGroupDelete(
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

	grpID, err := helper.ParseInt64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := originv3.NewDeleteGroupParams()
	params.GroupId = int32(grpID)
	params.MediaType = enums.HttpLarge.String()

	err = svc.Common.DeleteGroup(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleted origin group id: %v", grpID)
	d.SetId("")

	return diag.Diagnostics{}
}

func expandHttpLargeOriginGroup(
	d *schema.ResourceData,
) (*originv3.CustomerOriginGroupHTTPRequest, []*originv3.CustomerOriginRequest, []error) {
	if d == nil {
		return nil, make([]*originv3.CustomerOriginRequest, 0), []error{errors.New("no data to read")}
	}

	errs := make([]error, 0)

	originGrpState := &originv3.CustomerOriginGroupHTTPRequest{}
	originsState := make([]*originv3.CustomerOriginRequest, 0)

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

	if v, ok := d.GetOk("origin"); ok {
		if origins, err := expandOrigins(v); err == nil {
			originsState = origins
		} else {
			errs = append(errs, fmt.Errorf("error parsing origins: %w", err))
		}
	}

	log.Printf("[INFO] Origins: %# v\n", pretty.Formatter(originsState))

	return originGrpState, originsState, errs
}

func setHttpLargeOriginGroupState(
	d *schema.ResourceData,
	resp *originv3.CustomerOriginGroupHTTP,
) error {

	d.Set("name", resp.Name)
	d.Set("host_header", resp.HostHeader)
	d.Set("network_type_id", resp.NetworkTypeId)
	d.Set("strict_pci_certified", resp.StrictPciCertified)

	if len(resp.ShieldPops) > 0 {
		d.Set("shield_pops", resp.ShieldPops)
	} else {
		d.Set("shield_pops", make([]string, 0))
	}

	flattenedTLSSettings := flattenTLSSettings(resp.TlsSettings)
	d.Set("tls_settings", flattenedTLSSettings)

	return nil
}
