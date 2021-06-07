// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"terraform-provider-vmp/vmp/api"

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration."},
			"master_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Master Server GroupID."},
			"master_server_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Master Server Group Name"},
			"masters": {
				Type:     schema.TypeList,
				Computed: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ipaddress": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required:    true,
				Description: "List of Master Server in the Master Server Group"},
		},
	}
}

func ResourceMSGCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	newMSGName := d.Get("master_server_group_name").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber
	masters := d.Get("masters").([]interface{})
	masterServers := []api.MasterServer{}

	for _, item := range masters {
		curr := item.(map[string]interface{})

		id := curr["id"].(int)
		name := curr["name"].(string)
		ipaddress := curr["ipaddress"].(string)

		master := api.MasterServer{
			ID:        id,
			Name:      name,
			IPAddress: ipaddress,
		}

		masterServers = append(masterServers, master)
	}

	masterServerGroupRequest := &api.MasterServerGroupRequest{
		Name:    newMSGName,
		Masters: masterServers,
	}

	log.Printf("[INFO] Creating a new Master Server Group for Account '%s': %+v", accountNumber, masterServerGroupRequest)

	dnsrouteClient := api.NewDNSRouteAPIClient(*config)

	resp, err := dnsrouteClient.AddMasterServerGroup(masterServerGroupRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Create successful - New Master Server Group ID: %d", resp[0].MasterGroupID)
	d.SetId(strconv.Itoa(resp[0].MasterGroupID))

	ResourceMSGRead(ctx, d, m)

	return diags
}

func ResourceMSGRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)

	log.Printf("[INFO] Retrieving Master Server Group by GroupID: %d", groupId)

	resp, err := dnsRouteClient.GetMasterServerGroup(groupId)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Master Server Group: %+v", resp[0])
	msg := flattenMasterServerGroupData(*resp[0])
	newId := strconv.Itoa(resp[0].MasterGroupID)

	d.SetId(newId)
	d.Set("account_number", accountNumber)
	d.Set("master_group_id", resp[0].MasterGroupID)
	d.Set("master_server_group_name", resp[0].Name)
	jsonMsg, _ := json.Marshal(msg)
	log.Printf("master_server_group>>ResourceMSGRead>>msg:%s", jsonMsg)
	d.Set("masters", msg)
	return diags
}

func ResourceMSGUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	accountNumber := d.Get("account_number").(string)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)

	masters := d.Get("masters").([]interface{})
	masterGroupName := d.Get("master_server_group_name").(string)

	masterServers := []api.MasterServerWithGroupID{}

	for _, item := range masters {
		curr := item.(map[string]interface{})

		id := curr["id"].(int)
		name := curr["name"].(string)
		ipaddress := curr["ipaddress"].(string)

		master := api.MasterServerWithGroupID{
			ID:            id,
			Name:          name,
			IPAddress:     ipaddress,
			MasterGroupID: groupID,
		}

		masterServers = append(masterServers, master)
	}

	masterServerGroupUpdateRequest := &api.MasterServerGroupUpdateRequest{
		Name:    masterGroupName,
		ID:      groupID,
		Masters: masterServers,
	}
	err = dnsRouteClient.UpdateMasterServerGroup(masterServerGroupUpdateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceMSGRead(ctx, d, m)
}

func ResourceMSGDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteAPIClient := api.NewDNSRouteAPIClient(*config)

	msgID, _ := strconv.Atoi(d.Id())
	log.Printf("[INFO] Retrieving Master Server Group by GroupID for deletion: %d", msgID)

	err := dnsRouteAPIClient.DeleteMasterServerGroup(msgID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenMasterServerGroupData(mstGroup api.MasterServerGroupResponse) []interface{} {
	if mstGroup.Masters != nil {
		svgs := make([]interface{}, len(mstGroup.Masters), len(mstGroup.Masters))

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
