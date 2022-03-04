// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"log"
	"strconv"

	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Secondary Group
func ResourceSecondaryZoneGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceSecondaryZoneGroupCreate,
		ReadContext:   ResourceSecondaryZoneGroupRead,
		UpdateContext: ResourceSecondaryZoneGroupUpdate,
		DeleteContext: ResourceSecondaryZoneGroupDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration."},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secondary Group Name."},
			"zone_composition": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Collection of Secondary Zone Info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zones": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "secondary domain zone name",
									},
									"status": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "1:Active",
									},
									"zone_type": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "2:SecondaryZone",
									},
									"comment": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Comment, at least provide an empty string.",
									},
								},
							},
						},
						"master_group_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "master group id",
						},
						"master_server_tsigs": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_server": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"master_server_id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Referenced master server id",
												},
											},
										},
									},
									"tsig": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tsig_id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Referenced tsig id",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ResourceSecondaryZoneGroupCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Construct Secondary Zone Group Object
	name := d.Get("name").(string)
	zoneComposition := d.Get("zone_composition").([]interface{})
	zc := zoneComposition[0].(map[string]interface{})
	masterGroupID := zc["master_group_id"].(int)
	zones := zc["zones"].([]interface{})
	zoneArr := []routedns.SecondaryZone{}
	for _, item := range zones {
		curr := item.(map[string]interface{})

		domainName := curr["domain_name"].(string)
		status := curr["status"].(int)
		comment := curr["comment"].(string)
		if len(comment) == 0 {
			comment = ""
		}
		secondaryZone := routedns.SecondaryZone{
			DomainName: domainName,
			Status:     status,
			Comment:    comment,
		}
		zoneArr = append(zoneArr, secondaryZone)
	}

	mst := zc["master_server_tsigs"].([]interface{})
	mstArr := []routedns.MasterServerTSIGIDs{}
	for _, item := range mst {
		curr := item.(map[string]interface{})

		server := curr["master_server"].([]interface{})
		tsig := curr["tsig"].([]interface{})

		s := server[0].(map[string]interface{})
		sID := s["master_server_id"].(int)

		ms := routedns.MasterServerID{
			ID: sID,
		}

		t := tsig[0].(map[string]interface{})
		tID := t["tsig_id"].(int)
		ts := routedns.TSIGID{
			ID: tID,
		}

		mstItem := routedns.MasterServerTSIGIDs{
			MasterServer: ms,
			TSIG:         ts,
		}
		mstArr = append(mstArr, mstItem)
	}

	zcr := routedns.ZoneComposition{
		MasterGroupID:     masterGroupID,
		MasterServerTSIGs: mstArr,
		Zones:             zoneArr,
	}

	secondaryZoneGroup := routedns.SecondaryZoneGroup{
		Name:            name,
		ZoneComposition: zcr,
	}

	// Call Create Secondary Zone Group API
	log.Printf(
		"[INFO] Creating a new Secondary Zone Group for Account '%s': %+v",
		accountNumber,
		secondaryZoneGroup,
	)
	params := routedns.NewAddSecondaryZoneGroupParams()
	params.AccountNumber = accountNumber
	params.SecondaryZoneGroup = secondaryZoneGroup

	resp, err := routeDNSService.AddSecondaryZoneGroup(*params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Create successful - New Secondary Zone Group ID: %d",
		resp,
	)

	d.SetId(strconv.Itoa(resp.ID))

	return ResourceSecondaryZoneGroupRead(ctx, d, m)
}

func ResourceSecondaryZoneGroupRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	secondaryZoneGroupID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Call Get Secondary Zone Group API
	params := routedns.NewGetSecondaryZoneGroupParams()
	params.AccountNumber = accountNumber
	params.ID = secondaryZoneGroupID

	log.Printf(
		"[INFO] Retrieving Secondary Zone Group by zoneID: %d",
		secondaryZoneGroupID,
	)

	resp, err := routeDNSService.GetSecondaryZoneGroup(*params)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Secondary Zone Group: %+v", resp)

	// Update Terraform state with retrieved Secondary Zone Group data
	newID := strconv.Itoa(resp.ID)
	d.Set("name", resp.Name)
	d.Set("account_number", accountNumber)
	//TODO: Properly process zone composition into terraform state
	//d.Set("ZoneComposition", resp.ZoneComposition)
	d.SetId(newID)
	return diag.Diagnostics{}
}

func ResourceSecondaryZoneGroupUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	secondaryZoneGroupID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Construct Secondary Zone Group Update Object
	name := d.Get("name").(string)
	zoneComposition := d.Get("zone_composition").([]interface{})
	zc := zoneComposition[0].(map[string]interface{})
	masterGroupID := zc["master_group_id"].(int)
	zones := zc["zones"].([]interface{})
	zoneArr := []routedns.SecondaryZoneResponse{}
	for _, item := range zones {
		curr := item.(map[string]interface{})

		domainName := curr["domain_name"].(string)
		status := curr["status"].(int)
		comment := curr["comment"].(string)
		if len(comment) == 0 {
			comment = " "
		}
		secondaryZone := routedns.SecondaryZoneResponse{}
		secondaryZone.DomainName = domainName
		secondaryZone.Status = status
		secondaryZone.Comment = comment
		zoneArr = append(zoneArr, secondaryZone)
	}

	mst := zc["master_server_tsigs"].([]interface{})
	mstArr := []routedns.MasterServerTSIG{}
	for _, item := range mst {
		curr := item.(map[string]interface{})

		server := curr["master_server"].([]interface{})
		tsig := curr["tsig"].([]interface{})

		s := server[0].(map[string]interface{})
		sID := s["master_server_id"].(int)

		ms := routedns.MasterServer{
			ID: sID,
		}

		t := tsig[0].(map[string]interface{})
		tID := t["tsig_id"].(int)
		ts := routedns.TSIGGetOK{
			ID: tID,
		}

		mstItem := routedns.MasterServerTSIG{
			MasterServer: ms,
			TSIG:         ts,
		}
		mstArr = append(mstArr, mstItem)
	}

	zcr := routedns.ZoneCompositionResponse{
		MasterGroupID:     masterGroupID,
		MasterServerTsigs: mstArr,
		Zones:             zoneArr,
	}

	// Retrieve Existing Secondary Zone Group Object
	getParams := routedns.NewGetSecondaryZoneGroupParams()
	getParams.AccountNumber = accountNumber
	getParams.ID = secondaryZoneGroupID
	groupObj, err := routeDNSService.GetSecondaryZoneGroup(*getParams)

	// Update Secondary Zone Group Object
	groupObj.Name = name
	groupObj.ZoneComposition = zcr

	// Call Update Secondary Zone Group API
	updateParams := routedns.NewUpdateSecondaryZoneGroupParams()
	updateParams.AccountNumber = accountNumber
	updateParams.SecondaryZoneGroup = *groupObj

	err = routeDNSService.UpdateSecondaryZoneGroup(*updateParams)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceSecondaryZoneGroupRead(ctx, d, m)
}

func ResourceSecondaryZoneGroupDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	secondaryZoneGroupID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Call Get Secondary Zone Group API
	getParams := routedns.NewGetSecondaryZoneGroupParams()
	getParams.AccountNumber = accountNumber
	getParams.ID = secondaryZoneGroupID
	groupObj, err := routeDNSService.GetSecondaryZoneGroup(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Delete Secondary Zone Group API
	deleteParams := routedns.NewDeleteSecondaryZoneGroupParams()
	deleteParams.AccountNumber = accountNumber
	deleteParams.SecondaryZoneGroup = *groupObj
	err = routeDNSService.DeleteSecondaryZoneGroup(*deleteParams)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}
