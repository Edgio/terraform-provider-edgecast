// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"terraform-provider-edgecast/edgecast/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Master Server Group
func ResourceMasterServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceMSGCreate,
		ReadContext:   ResourceMSGRead,
		UpdateContext: ResourceMSGUpdate,
		DeleteContext: ResourceMSGDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Account Number associated with the customer whose 
				resources you wish to manage. This account number may be found 
				in the upper right-hand corner of the MCC.`},
			"master_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `Indicates the system-defined ID assigned to a 
				master server group.`},
			"master_server_group_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Indicates the name that will be assigned to the 
				new master server group.`},
			"masters": {
				Type:     schema.TypeList,
				Computed: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `Indicates the system-defined ID 
							assigned to an existing master name server that will 
							be associated with the master server group being 
							created.`,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Indicates the name that will be 
							assigned to a new master name server that will be 
							associated with the master server group being 
							created.`,
						},
						"ipaddress": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Indicates the IP address that will be 
							assigned to a new master name server that will be 
							associated with the master server group being 
							created.`,
						},
					},
				},
				Required: true,
				Description: `Contains the master name servers associated with 
				a master server group.`,
			},
		},
	}
}

func ResourceMSGCreate(
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

	// Construct Master Server Group Object
	masterServerGroupName := d.Get("master_server_group_name").(string)
	masterServers := expandMasterServers(d.Get("masters").([]interface{}))

	masterServerGroup := &routedns.MasterServerGroupAddRequest{
		Name:    masterServerGroupName,
		Masters: masterServers,
	}

	// Call Add Master Server Group API
	log.Printf(
		"[INFO] Creating a new Master Server Group for Account '%s': %+v",
		accountNumber,
		masterServerGroup,
	)

	params := routedns.NewAddMasterServerGroupParams()
	params.AccountNumber = accountNumber
	params.MasterServerGroup = *masterServerGroup
	resp, err := routeDNSService.AddMasterServerGroup(*params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Create successful - New Master Server Group ID: %d",
		resp.MasterGroupID,
	)
	d.SetId(strconv.Itoa(resp.MasterGroupID))

	return ResourceMSGRead(ctx, d, m)
}

func ResourceMSGRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)

	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Get Master Server Group API
	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := routedns.NewGetMasterServerGroupParams()
	params.AccountNumber = accountNumber
	params.MasterServerGroupID = groupID

	log.Printf("[INFO] Retrieving Master Server Group by GroupID: %d", groupID)

	resp, err := routeDNSService.GetMasterServerGroup(*params)

	if err != nil {
		return diag.FromErr(err)
	}

	// Update Terraform state with retrieved Master Server Group data
	log.Printf("[INFO] Retrieved Master Server Group: %+v", resp)
	msg := flattenMasterServers(*resp)
	newId := strconv.Itoa(resp.MasterGroupID)

	d.SetId(newId)
	d.Set("account_number", accountNumber)
	d.Set("master_group_id", resp.MasterGroupID)
	d.Set("master_server_group_name", resp.Name)
	jsonMsg, _ := json.Marshal(msg)
	log.Printf("master_server_group>>ResourceMSGRead>>msg:%s", jsonMsg)
	d.Set("masters", msg)

	return diag.Diagnostics{}
}

func ResourceMSGUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)

	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Construct Master Server Group Update Data
	masterServerGroupName := d.Get("master_server_group_name").(string)
	masterServers := expandMasterServers(d.Get("masters").([]interface{}))

	masterServerGroupUpdateRequest := routedns.MasterServerGroupUpdateRequest{
		MasterGroupID: groupID,
		MasterServerGroup: routedns.MasterServerGroup{
			Name:    masterServerGroupName,
			Masters: masterServers,
		},
	}

	// Call Update Group API
	params := routedns.NewUpdateMasterServerGroupParams()
	params.AccountNumber = accountNumber
	params.MasterServerGroup = masterServerGroupUpdateRequest
	err = routeDNSService.UpdateMasterServerGroup(*params)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceMSGRead(ctx, d, m)
}

func ResourceMSGDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	msgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Get Master Server Group API
	log.Printf(
		"[INFO] Retrieving Master Server Group by GroupID for deletion: %d",
		msgID,
	)

	getParams := routedns.NewGetMasterServerGroupParams()
	getParams.AccountNumber = accountNumber
	getParams.MasterServerGroupID = msgID
	groupObj, err := routeDNSService.GetMasterServerGroup(*getParams)

	if err != nil {
		return diag.FromErr(err)
	}

	// Call Delete Master Server Group API
	deleteParams := routedns.NewDeleteMasterServerGroupParams()
	deleteParams.AccountNumber = accountNumber
	deleteParams.MasterServerGroup = *groupObj
	err = routeDNSService.DeleteMasterServerGroup(*deleteParams)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}

func expandMasterServers(masters []interface{}) []routedns.MasterServer {
	masterServers := []routedns.MasterServer{}

	for _, item := range masters {
		curr := item.(map[string]interface{})

		id := curr["id"].(int)
		name := curr["name"].(string)
		ipaddress := curr["ipaddress"].(string)

		master := routedns.MasterServer{
			ID:        id,
			Name:      name,
			IPAddress: ipaddress,
		}

		masterServers = append(masterServers, master)
	}
	return masterServers
}

func flattenMasterServers(
	mstGroup routedns.MasterServerGroupAddGetOK,
) []interface{} {
	if mstGroup.Masters != nil {
		svgs := make([]interface{}, len(mstGroup.Masters))

		for i, masterServer := range mstGroup.Masters {
			oi := make(map[string]interface{})

			oi["id"] = masterServer.ID
			oi["name"] = masterServer.Name
			oi["ipaddress"] = masterServer.IPAddress

			svgs[i] = oi
		}

		return svgs
	}

	return make([]interface{}, 0)
}
