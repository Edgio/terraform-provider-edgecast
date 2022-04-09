// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"log"
	"strconv"

	"terraform-provider-edgecast/ec/api"

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
				Type:     schema.TypeString,
				Required: true,
				Description: `Account Number associated with the customer whose 
				resources you wish to manage. This account number may be found 
				in the upper right-hand corner of the MCC.`},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Indicates the name assigned to the new secondary 
				zone group.`},
			"zone_composition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Description: `ZoneCompositionResponse defines parameters of the 
				secondary zone group.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zones": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_name": {
										Type:     schema.TypeString,
										Required: true,
										Description: `Identifies a secondary 
										zone by its zone name 
										(e.g., example.com). Edgecast name 
										servers will request a zone transfer for 
										this zone. This name must match the one 
										defined on the master name server(s) 
										associated with this secondary zone 
										group.`,
									},
									"status": {
										Type:     schema.TypeInt,
										Required: true,
										Description: `Defines whether the zone 
										is enabled or disabled. Valid values 
										are: 1 - Enabled, 2 - Disabled`,
									},
									"zone_type": {
										Type:     schema.TypeInt,
										Required: true,
										Description: `This parameter is reserved 
										for future use. The only supported value 
										for this parameter is "2".`,
									},
									"comment": {
										Type:     schema.TypeString,
										Required: true,
										Description: `Comment about this 
										secondary zone.`,
									},
								},
							},
						},
						"master_group_id": {
							Type:     schema.TypeInt,
							Required: true,
							Description: `Associates a master server group, as 
							identified by its system-defined ID, with the 
							secondary zone group.`,
						},
						"master_server_tsigs": {
							Type:     schema.TypeList,
							Required: true,
							Description: `Defines TSIG keys to the desired 
							master name servers in the master server group.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_server": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"master_server_id": {
													Type:     schema.TypeInt,
													Required: true,
													Description: `Identifies the 
													master name server to which 
													a TSIG key will be assigned.`,
												},
											},
										},
									},
									"tsig": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tsig_id": {
													Type:     schema.TypeInt,
													Required: true,
													Description: `Identifies the 
													TSIG key that will be 
													assigned to the master name 
													server identified by the 
													MasterServer object.`,
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
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// Construct Secondary Zone Group Object
	name := d.Get("name").(string)
	zoneComposition, err := expandZoneCompositionCreate(
		d.Get("zone_composition"),
	)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	secondaryZoneGroup := routedns.SecondaryZoneGroup{
		Name:            name,
		ZoneComposition: *zoneComposition,
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
		resp.ID,
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
	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

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
	zoneComposition := flattenZoneComposition(resp.ZoneComposition)
	d.Set("zone_composition", zoneComposition)
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
	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Construct Secondary Zone Group Update Object
	name := d.Get("name").(string)
	zoneComposition, err := expandZoneCompositionUpdate(
		d.Get("zone_composition"),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Existing Secondary Zone Group Object
	getParams := routedns.NewGetSecondaryZoneGroupParams()
	getParams.AccountNumber = accountNumber
	getParams.ID = secondaryZoneGroupID
	groupObj, err := routeDNSService.GetSecondaryZoneGroup(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update Secondary Zone Group Object
	groupObj.Name = name
	groupObj.ZoneComposition = *zoneComposition

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
	if err != nil {
		return diag.FromErr(err)
	}

	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

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

func expandZoneCompositionCreate(zoneCompositionList interface{},
) (*routedns.ZoneComposition, error) {
	// This is a list of length one, restricted in the schema definition
	// TODO: Review if this should be a set of type one
	zcl := zoneCompositionList.([]interface{})
	zoneComposition := zcl[0].(map[string]interface{})

	masterGroupID := zoneComposition["master_group_id"].(int)
	zones := expandZones(zoneComposition["zones"].([]interface{}))
	masterServerTSIGs, err := expandMasterServerTSIGs(
		zoneComposition["master_server_tsigs"].([]interface{}),
	)
	if err != nil {
		return nil, err
	}

	zc := routedns.ZoneComposition{
		MasterGroupID:     masterGroupID,
		MasterServerTSIGs: *masterServerTSIGs,
		Zones:             zones,
	}
	return &zc, nil
}

// Due to differences in the payload requirements for Secondary Zone API Add vs
// Update requests, it is necessary to construct a different object. Using the
// same expand methods for reusability and transforming the data for updates.
func expandZoneCompositionUpdate(zoneCompositionList interface{},
) (*routedns.ZoneCompositionResponse, error) {
	// This is a list of length one, restricted in the schema definition
	// TODO: Review if this should be a set of type one
	zcl := zoneCompositionList.([]interface{})
	zoneComposition := zcl[0].(map[string]interface{})

	masterGroupID := zoneComposition["master_group_id"].(int)
	zones := expandZones(zoneComposition["zones"].([]interface{}))
	masterServerTSIGs, err := expandMasterServerTSIGs(
		zoneComposition["master_server_tsigs"].([]interface{}),
	)
	if err != nil {
		return nil, err
	}

	// Convert routedns.SecondaryZone to routedns.SecondaryZoneResponse
	// Update operation API requires a different payload
	updateZones := make([]routedns.SecondaryZoneResponse, 0)
	for _, zone := range zones {
		secondaryZone := routedns.SecondaryZoneResponse{
			SecondaryZone: routedns.SecondaryZone{
				Comment:    zone.Comment,
				DomainName: zone.DomainName,
				Status:     zone.Status,
			},
		}
		updateZones = append(updateZones, secondaryZone)
	}

	// Convert routedns.MasterServerTSIGIDs to routedns.MasterServerTSIG
	// Update operation API requires a different payload
	updateTSIGs := make([]routedns.MasterServerTSIG, 0)
	for _, tsig := range *masterServerTSIGs {
		tsig := routedns.MasterServerTSIG{
			MasterServer: routedns.MasterServer{
				ID: tsig.MasterServer.ID,
			},
			TSIG: routedns.TSIGGetOK{
				ID: tsig.TSIG.ID,
			},
		}
		updateTSIGs = append(updateTSIGs, tsig)
	}

	zc := routedns.ZoneCompositionResponse{
		MasterGroupID:     masterGroupID,
		MasterServerTsigs: updateTSIGs,
		Zones:             updateZones,
	}
	return &zc, nil
}

func expandZones(zones []interface{}) []routedns.SecondaryZone {
	secondaryZones := []routedns.SecondaryZone{}
	for _, item := range zones {
		curr := item.(map[string]interface{})

		domainName := curr["domain_name"].(string)
		status := curr["status"].(int)
		comment := curr["comment"].(string)
		// TODO: Determine if this logic is required
		if len(comment) == 0 {
			comment = ""
		}
		secondaryZone := routedns.SecondaryZone{
			DomainName: domainName,
			Status:     status,
			Comment:    comment,
		}
		secondaryZones = append(secondaryZones, secondaryZone)
	}

	return secondaryZones
}

func expandMasterServerTSIGs(
	tsigs []interface{},
) (*[]routedns.MasterServerTSIGIDs, error) {
	masterServerTSIGs := []routedns.MasterServerTSIGIDs{}
	for _, item := range tsigs {
		curr := item.(map[string]interface{})

		// This is a list of length one, restricted in the schema definition
		// TODO: Review if this should be a set of type one
		s := curr["master_server"].([]interface{})
		server := s[0].(map[string]interface{})

		serverID := server["master_server_id"].(int)
		ms := routedns.MasterServerID{
			ID: serverID,
		}

		// This is a list of length one, restricted in the schema definition
		// TODO: Review if this should be a set of type one
		t := curr["tsig"].([]interface{})
		tsig := t[0].(map[string]interface{})
		tsigID := tsig["tsig_id"].(int)
		ts := routedns.TSIGID{
			ID: tsigID,
		}

		masterServerTSIGIDs := routedns.MasterServerTSIGIDs{
			MasterServer: ms,
			TSIG:         ts,
		}
		masterServerTSIGs = append(masterServerTSIGs, masterServerTSIGIDs)
	}
	return &masterServerTSIGs, nil
}

func flattenZones(
	secondaryZones []routedns.SecondaryZoneResponse,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, zone := range secondaryZones {
		z := make(map[string]interface{}, 0)
		z["comment"] = zone.Comment
		z["domain_name"] = zone.DomainName
		z["status"] = zone.Status
		flattened = append(flattened, z)
	}
	return flattened
}

func flattenMasterServerTSIGs(
	masterServerTSIGs []routedns.MasterServerTSIG,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, tsig := range masterServerTSIGs {
		masterServerTSIG := make(map[string]interface{}, 0)

		msID := make(map[string]interface{}, 0)
		msID["master_server_id"] = tsig.MasterServer.ID
		masterServerTSIG["master_server"] = msID

		tsID := make(map[string]interface{}, 0)
		tsID["tsig_id"] = tsig.TSIG.ID
		masterServerTSIG["tsig"] = tsID

		flattened = append(flattened, masterServerTSIG)
	}

	return flattened
}

func flattenZoneComposition(
	secondaryZoneGroup routedns.ZoneCompositionResponse,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)
	zc := make(map[string]interface{})
	zc["master_group_id"] = secondaryZoneGroup.MasterGroupID
	zc["zones"] = flattenZones(secondaryZoneGroup.Zones)
	zc["master_server_tsigs"] = flattenMasterServerTSIGs(
		secondaryZoneGroup.MasterServerTsigs,
	)

	flattened = append(flattened, zc)
	return flattened
}
