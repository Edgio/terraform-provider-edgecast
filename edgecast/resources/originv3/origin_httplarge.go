// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package originv3

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/originv3"
	"github.com/EdgeCast/ec-sdk-go/edgecast/shared/enums"
	"github.com/ahmetalpbalkan/go-linq"
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

	// Initialize Service
	config, ok := m.(internal.ProviderConfig)
	if !ok {
		return helper.CreationErrorf(d, "failed to load configuration")
	}

	svc, err := buildOriginV3Service(config)
	if err != nil {
		return diag.FromErr(err)
	}

	originGroupState, errs := expandHttpLargeOriginGroup(d)
	if len(errs) > 0 {
		return helper.DiagsFromErrors("error parsing origin group", errs)
	}
	// Call APIs.
	cparams := originv3.NewAddHttpLargeGroupParams()
	cparams.CustomerOriginGroupHTTPRequest =
		originv3.CustomerOriginGroupHTTPRequest{
			Name:               originGroupState.Name,
			HostHeader:         originv3.NewNullableString(originGroupState.HostHeader),
			ShieldPops:         originGroupState.ShieldPops,
			NetworkTypeId:      originv3.NewNullableInt32(originGroupState.NetworkTypeID),
			StrictPciCertified: originv3.NewNullableBool(originGroupState.StrictPCICertified),
			TlsSettings:        originGroupState.TLSSettings,
		}

	cresp, err := svc.HttpLargeOnly.AddHttpLargeGroup(cparams)
	if err != nil {
		return helper.CreationError(d, err)
	}
	log.Printf("[INFO] origin group created: %# v\n", pretty.Formatter(cresp))
	log.Printf("[INFO] origin group id: %d\n", cresp.Id)

	grpID := cresp.Id

	mlock := &sync.Mutex{}
	wg := sync.WaitGroup{}

	if len(originGroupState.Origins) > 0 {
		errs := make([]error, 0)
		for _, origin := range originGroupState.Origins {

			params := originv3.NewAddOriginParams()
			params.MediaType = enums.HttpLarge.String()
			params.CustomerOriginRequest = originv3.CustomerOriginRequest{
				GroupId:        *grpID,
				Name:           originv3.NewNullableString(origin.Name),
				Host:           origin.Host,
				Port:           &origin.Port,
				IsPrimary:      origin.IsPrimary,
				StorageTypeId:  originv3.NewNullableInt32(origin.StorageTypeID),
				ProtocolTypeId: originv3.NewNullableInt32(origin.ProtocolTypeID),
			}

			// Spin up a worker to call the api.
			wg.Add(1)

			go func(params originv3.AddOriginParams) {
				defer wg.Done()

				resp, err := svc.Common.AddOrigin(params)
				if err == nil {
					log.Printf("[INFO] origin created: %# v\n", pretty.Formatter(resp))
				} else {
					mlock.Lock()
					errs = append(errs, err)
					mlock.Unlock()
				}
			}(params)
		}
		// Wait for all api workers to finish.
		wg.Wait()

		if len(errs) > 0 {
			d.SetId("")

			//delete the created origin group
			deleteparams := originv3.NewDeleteGroupParams()
			deleteparams.GroupId = *grpID
			deleteparams.MediaType = enums.HttpLarge.String()

			deleteErr := svc.Common.DeleteGroup(deleteparams)
			if deleteErr != nil {
				errs = append(errs, deleteErr)
			}
			return helper.DiagsFromErrors("error updating origin group", errs)
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

	originsParams := originv3.NewGetOriginsByGroupParams()
	originsParams.GroupId = int32(grpID)
	originsParams.MediaType = enums.HttpLarge.String()

	originsResp, err := svc.Common.GetOriginsByGroup(originsParams)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Retrieved origins: %# v\n", pretty.Formatter(originsResp))

	// Write TF state.
	err = setHttpLargeOriginGroupState(d, resp, originsResp)
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

	originGroupState, errors := expandHttpLargeOriginGroup(d)
	if len(errors) > 0 {
		return helper.DiagsFromErrors("error parsing origin group", errors)
	}
	originGroupState.ID = int32(grpID)

	log.Printf("[INFO] Updating origin group : ID: %d\n", grpID)
	errs := make([]error, 0)
	mlock := &sync.Mutex{}
	wg := sync.WaitGroup{}

	updateParams := originv3.NewUpdateHttpLargeGroupParams()
	updateParams.GroupId = int32(grpID)
	updateParams.CustomerOriginGroupHTTPRequest = originv3.CustomerOriginGroupHTTPRequest{
		Name:               originGroupState.Name,
		HostHeader:         originv3.NewNullableString(originGroupState.HostHeader),
		ShieldPops:         originGroupState.ShieldPops,
		NetworkTypeId:      originv3.NewNullableInt32(originGroupState.NetworkTypeID),
		StrictPciCertified: originv3.NewNullableBool(originGroupState.StrictPCICertified),
		TlsSettings:        originGroupState.TLSSettings,
	}

	// Spin up a worker to call the api.
	wg.Add(1)
	go func(updateParams originv3.UpdateHttpLargeGroupParams) {
		defer wg.Done()

		_, err = svc.HttpLargeOnly.UpdateHttpLargeGroup(updateParams)
		if err == nil {
			log.Printf("[INFO] Updated origin group info: ID: %d\n", grpID)
		} else {
			mlock.Lock()
			errs = append(errs, err)
			mlock.Unlock()
		}
	}(updateParams)

	if d.HasChange("origin") {
		toadd, todelete, toupdate := getOriginsDelta(d)
		log.Printf("[INFO] toadd: %# v\n", pretty.Formatter(toadd))
		log.Printf("[INFO] todelete: %# v\n", pretty.Formatter(todelete))
		log.Printf("[INFO] toupdate: %# v\n", pretty.Formatter(toupdate))

		//add new origins
		if len(toadd) > 0 {
			for _, v := range toadd {
				params := originv3.NewAddOriginParams()
				params.MediaType = enums.HttpLarge.String()
				params.CustomerOriginRequest = originv3.CustomerOriginRequest{
					GroupId:        int32(grpID),
					Name:           originv3.NewNullableString(v.Name),
					Host:           v.Host,
					Port:           &v.Port,
					IsPrimary:      v.IsPrimary,
					StorageTypeId:  originv3.NewNullableInt32(v.StorageTypeID),
					ProtocolTypeId: originv3.NewNullableInt32(v.ProtocolTypeID),
				}

				wg.Add(1)
				go func(params originv3.AddOriginParams) {
					defer wg.Done()

					resp, err := svc.Common.AddOrigin(params)
					if err == nil {
						log.Printf("[INFO] Added origin: ID: %d\n", *resp.Id)
					} else {
						mlock.Lock()
						errs = append(errs, err)
						mlock.Unlock()
					}
				}(params)
			}
		}

		if len(todelete) > 0 {
			for _, v := range todelete {
				params := originv3.NewDeleteOriginParams()
				params.Id = v.ID
				params.MediaType = enums.HttpLarge.String()

				wg.Add(1)
				go func(params originv3.DeleteOriginParams) {
					defer wg.Done()

					err := svc.Common.DeleteOrigin(params)
					if err == nil {
						log.Printf("[INFO] Deleted origin: ID: %d\n", params.Id)
					} else {
						mlock.Lock()
						errs = append(errs, err)
						mlock.Unlock()
					}
				}(params)
			}
		}

		if len(toupdate) > 0 {
			for _, v := range toupdate {
				params := originv3.NewUpdateOriginParams()
				params.Id = v.ID
				params.MediaType = enums.HttpLarge.String()
				params.CustomerOriginRequest = originv3.CustomerOriginRequest{
					GroupId:        int32(grpID),
					Name:           originv3.NewNullableString(v.Name),
					Host:           v.Host,
					Port:           &v.Port,
					IsPrimary:      v.IsPrimary,
					StorageTypeId:  originv3.NewNullableInt32(v.StorageTypeID),
					ProtocolTypeId: originv3.NewNullableInt32(v.ProtocolTypeID),
				}

				wg.Add(1)
				go func(params originv3.UpdateOriginParams) {
					defer wg.Done()

					_, err := svc.Common.UpdateOrigin(params)
					if err == nil {
						log.Printf("[INFO] Updated origin: ID: %d\n", params.Id)
					} else {
						mlock.Lock()
						errs = append(errs, err)
						mlock.Unlock()
					}
				}(params)
			}
		}
	}
	// Wait for all api workers to finish.
	wg.Wait()

	if len(errs) > 0 {
		return helper.DiagsFromErrors("error updating origin group", errs)
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
) (*OriginGroupState, []error) {
	if d == nil {
		return nil, []error{errors.New("no data to read")}
	}

	errs := make([]error, 0)

	originGrpState := &OriginGroupState{
		Origins: make([]*OriginState, 0),
	}

	if v, ok := d.GetOk("name"); ok {
		if name, ok := v.(string); ok {
			originGrpState.Name = name
		} else {
			errs = append(errs, errors.New("name not a string"))
		}
	}

	if v, ok := d.GetOk("host_header"); ok {
		if hostHeader, ok := v.(string); ok {
			originGrpState.HostHeader = hostHeader
		} else {
			errs = append(errs, errors.New("host_header not a string"))
		}
	}

	if v, ok := d.GetOk("network_type_id"); ok {
		if networkTypeID, ok := v.(int); ok {
			originGrpState.NetworkTypeID = int32(networkTypeID)
		} else {
			errs = append(errs, errors.New("network_type_id not a int32"))
		}
	}

	if v, ok := d.GetOk("strict_pci_certified"); ok {
		if strictPCICertified, ok := v.(bool); ok {
			originGrpState.StrictPCICertified = strictPCICertified
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
			originGrpState.TLSSettings = tlsSettings
		} else {
			errs = append(errs, fmt.Errorf("error parsing tls_settings: %w", err))
		}
	}

	if v, ok := d.GetOk("origin"); ok {
		if origins, err := expandOrigins(v); err == nil {
			originGrpState.Origins = origins
		} else {
			errs = append(errs, fmt.Errorf("error parsing origins: %w", err))
		}
	}

	return originGrpState, errs
}

func setHttpLargeOriginGroupState(
	d *schema.ResourceData,
	resp *originv3.CustomerOriginGroupHTTP,
	originsResp []originv3.CustomerOriginFailoverOrder,
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

	flattenedOrigins := flattenOrigins(originsResp)
	d.Set("origin", flattenedOrigins)

	return nil
}

func getOriginsDelta(
	d *schema.ResourceData,
) ([]*OriginState, []*OriginState, []*OriginState) {
	old, new := d.GetChange("origin")
	// Represents current resource, state prior to latest Terraform apply
	oldOrigins, _ := expandOrigins(old)
	// Repesents desired resource state
	newOrigins, _ := expandOrigins(new)

	toAdd := make([]*OriginState, 0)
	toUpdate := make([]*OriginState, 0)
	toDelete := make([]*OriginState, 0)

	// To add
	linq.From(newOrigins).Where(func(c interface{}) bool {
		return c.(*OriginState).ID == 0
	}).Select(func(c interface{}) interface{} {
		return (c.(*OriginState))
	}).ToSlice(&toAdd)

	// To delete
	linq.From(oldOrigins).
		ExceptBy(linq.From(newOrigins),
			func(c interface{}) interface{} { return c.(*OriginState).ID },
		).ToSlice(&toDelete)

	// To update
	linq.From(newOrigins).Where(func(c interface{}) bool {
		return c.(*OriginState).ID != 0
	}).Select(func(c interface{}) interface{} {
		return (c.(*OriginState))
	}).ToSlice(&toUpdate)

	return toAdd, toDelete, toUpdate
}

// OriginGroupState represents the state of a Origin Group as it exists in the
// TF state file. This is an intermediate model before being translated to API
// models.
type OriginGroupState struct {
	ID int32

	HostHeader string

	Name string

	NetworkTypeID int32

	ShieldPops []*string

	StrictPCICertified bool

	TLSSettings *originv3.TlsSettings

	Origins []*OriginState
}

// OriginState represents the state of a Origin as it exists in the
// TF state file. This is an intermediate model before being translated to API
// models.
type OriginState struct {
	ID int32

	GroupID int32

	Host string

	IsPrimary bool

	Name string

	Port int32

	ProtocolTypeID int32

	StorageTypeID int32

	HostHeader string
}
