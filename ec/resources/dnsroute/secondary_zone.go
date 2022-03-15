// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"log"
	"strconv"
	"terraform-provider-edgecast/ec/helper"

	"terraform-provider-edgecast/ec/api"

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
		Importer:      helper.Import(ResourceSecondaryZoneGroupRead, "account_number", "id"),

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

func ResourceSecondaryZoneGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	name := d.Get("name").(string)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	zoneComposition := d.Get("zone_composition").([]interface{})
	zc := zoneComposition[0].(map[string]interface{})
	masterGroupID := zc["master_group_id"].(int)
	zones := zc["zones"].([]interface{})
	zoneArr := []api.SecondaryZoneRequest{}
	for _, item := range zones {
		curr := item.(map[string]interface{})

		domainName := curr["domain_name"].(string)
		status := curr["status"].(int)
		zoneType := curr["zone_type"].(int)
		comment := curr["comment"].(string)
		if len(comment) == 0 {
			comment = ""
		}
		secondaryZone := api.SecondaryZoneRequest{
			DomainName:      domainName,
			Status:          status,
			ZoneType:        zoneType,
			Comment:         comment,
			IsCustomerOwned: true,
		}
		zoneArr = append(zoneArr, secondaryZone)
	}

	mst := zc["master_server_tsigs"].([]interface{})
	mstArr := []api.MasterServerTsigRequest{}
	for _, item := range mst {
		curr := item.(map[string]interface{})

		server := curr["master_server"].([]interface{})
		tsig := curr["tsig"].([]interface{})

		s := server[0].(map[string]interface{})
		sID := s["master_server_id"].(int)

		ms := api.MasterServerRequest{
			ID: sID,
		}

		t := tsig[0].(map[string]interface{})
		tID := t["tsig_id"].(int)
		ts := api.TsigRequest{
			ID: tID,
		}

		mstItem := api.MasterServerTsigRequest{
			MasterServer: ms,
			Tsig:         ts,
		}
		mstArr = append(mstArr, mstItem)
	}

	zcr := api.DnsZoneCompositionRequest{
		MasterGroupID:     masterGroupID,
		MasterServerTsigs: mstArr,
		Zones:             zoneArr,
	}

	secondaryZoneGroupRequest := &api.DnsRouteSecondaryZoneGroupRequest{
		ID:              -1,
		Name:            name,
		ZoneComposition: zcr,
		CustomerID:      -1,
	}

	log.Printf("[INFO] Creating a new Secondary Zone Group for Account '%s': %+v", accountNumber, secondaryZoneGroupRequest)

	dnsrouteClient := api.NewDNSRouteAPIClient(*config)

	resp, err := dnsrouteClient.AddSecondaryZoneGroup(secondaryZoneGroupRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Create successful - New Secondary Zone Group ID: %d", resp)
	d.SetId(strconv.Itoa(resp))

	ResourceSecondaryZoneGroupRead(ctx, d, m)

	return diags
}

func ResourceSecondaryZoneGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	secondaryZoneGroupId, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)

	log.Printf("[INFO] Retrieving Secondary Zone Group by zoneID: %d", secondaryZoneGroupId)

	resp, err := dnsRouteClient.GetSecondaryZoneGroup(secondaryZoneGroupId)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Secondary Zone Group: %+v", resp)
	newID := strconv.Itoa(resp.ID)
	d.Set("name", resp.Name)
	d.Set("account_number", accountNumber)
	//d.Set("ZoneComposition", resp.ZoneComposition)
	d.SetId(newID)
	return diags
}

func ResourceSecondaryZoneGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	secondaryZoneGroupID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)
	name := d.Get("name").(string)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	zoneComposition := d.Get("zone_composition").([]interface{})
	zc := zoneComposition[0].(map[string]interface{})
	masterGroupID := zc["master_group_id"].(int)
	zones := zc["zones"].([]interface{})
	zoneArr := []api.SecondaryZoneRequest{}
	for _, item := range zones {
		curr := item.(map[string]interface{})

		domainName := curr["domain_name"].(string)
		status := curr["status"].(int)
		zoneType := curr["zone_type"].(int)
		comment := curr["comment"].(string)
		if len(comment) == 0 {
			comment = " "
		}
		secondaryZone := api.SecondaryZoneRequest{
			DomainName:      domainName,
			Status:          status,
			ZoneType:        zoneType,
			Comment:         comment,
			IsCustomerOwned: true,
		}
		zoneArr = append(zoneArr, secondaryZone)
	}

	mst := zc["master_server_tsigs"].([]interface{})
	mstArr := []api.MasterServerTsigRequest{}
	for _, item := range mst {
		curr := item.(map[string]interface{})

		server := curr["master_server"].([]interface{})
		tsig := curr["tsig"].([]interface{})

		s := server[0].(map[string]interface{})
		sID := s["master_server_id"].(int)

		ms := api.MasterServerRequest{
			ID: sID,
		}

		t := tsig[0].(map[string]interface{})
		tID := t["tsig_id"].(int)
		ts := api.TsigRequest{
			ID: tID,
		}

		mstItem := api.MasterServerTsigRequest{
			MasterServer: ms,
			Tsig:         ts,
		}
		mstArr = append(mstArr, mstItem)
	}

	zcr := api.DnsZoneCompositionRequest{
		MasterGroupID:     masterGroupID,
		MasterServerTsigs: mstArr,
		Zones:             zoneArr,
	}

	secondaryZoneGroupRequest := &api.DnsRouteSecondaryZoneGroupRequest{
		ID:              secondaryZoneGroupID,
		Name:            name,
		ZoneComposition: zcr,
		CustomerID:      -1,
	}

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)
	err = dnsRouteClient.UpdateSecondaryZoneGroup(secondaryZoneGroupRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceSecondaryZoneGroupRead(ctx, d, m)
}

func ResourceSecondaryZoneGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteAPIClient := api.NewDNSRouteAPIClient(*config)

	szgID, _ := strconv.Atoi(d.Id())

	err := dnsRouteAPIClient.DeleteSecondaryZoneGroup(szgID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
