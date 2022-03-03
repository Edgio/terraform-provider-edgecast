// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package edgecname

import (
	"context"
	"log"
	"strconv"
	"terraform-provider-ec/ec/helper"

	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/edgecname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEdgeCname() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceEdgeCnameCreate,
		ReadContext:   ResourceEdgeCnameRead,
		UpdateContext: ResourceEdgeCnameUpdate,
		DeleteContext: ResourceEdgeCnameDelete,

		Importer: helper.AccountIDImporter(ResourceEdgeCnameRead),
		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Account Number for the customer if not already
					specified in the provider configuration.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Sets the name that will be assigned to the edge
					CNAME. It should only contain lower-case alphanumeric
					characters, dashes, and periods. The name specified for
					this parameter should also be defined as a CNAME record
					on a DNS server. The CNAME record defined on the DNS server 
					should point to the CDN hostname
					(e.g., wpc.0001.edgecastcdn.net) for the platform
					identified by the "platform" parameter`,
			},
			"dir_path": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Identifies a location on the origin server. This
					string should specify the relative path from the root
					folder of the origin server to the desired location. Set
					this parameter to blank to point the edge CNAME to the
					root folder of the origin server.`,
			},
			"enable_custom_reports": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Determines whether hits and data transferred
					statistics will be tracked for this edge CNAME. Logged
					data can be viewed through the Custom Reports module.
					Valid values are:
					0: Disabled (Default Value).
					1: Enabled. CDN activity on this edge CNAME will be logged.`,
			},
			"media_type_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Identifies the Delivery Platform on which the
					edge CNAME will be created. 
					3:Http Large, 8:HTTP Small, 14: ADN`,
			},
			"origin_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Identifies whether an edge CNAME will be created
					for a CDN origin server or a customer origin server. 
					Valid values: 
					-1: Indicates that you would like to create an
					edge CNAME for our CDN storage service,
					CustomerOriginID: Specifying an ID for an existing
					customer origin configuration indicates that you would
					like to create an edge CNAME for that customer origin`,
			},
			"origin_string": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the origin identifier, the account
					number, and the relative path associated with the edge CNAME.`,
			},
		},
	}
}

func ResourceEdgeCnameCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	// Create Edge CNAME object
	edgecnameObj := &edgecname.EdgeCname{
		Name:                d.Get("name").(string),
		DirPath:             d.Get("dir_path").(string),
		EnableCustomReports: d.Get("enable_custom_reports").(int),
		MediaTypeID:         d.Get("media_type_id").(int),
		OriginID:            d.Get("origin_id").(int),
	}

	// Initialize Edge CNAME Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	edgecnameService, err := buildEdgeCnameService(**config)
	if err != nil {
		d.SetId("") // Terraform requires an empty ID for failed creation
		return diag.FromErr(err)
	}

	// Call Add Edge CNAME API
	log.Printf(
		"[INFO] Creating Edge CNAME for Account '%s': %+v",
		accountNumber,
		edgecnameObj,
	)

	params := edgecname.NewAddEdgeCnameParams()
	params.AccountNumber = accountNumber
	params.EdgeCname = *edgecnameObj

	edgeCnameID, err := edgecnameService.AddEdgeCname(*params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Create successful - New Edge CNAME ID: %d", *edgeCnameID)
	d.SetId(strconv.Itoa(*edgeCnameID))

	return ResourceEdgeCnameRead(ctx, d, m)
}

func ResourceEdgeCnameRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Edge CNAME Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	edgecnameService, err := buildEdgeCnameService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Get EDGE CNAME API
	edgecnameID, _ := strconv.Atoi(d.Id())

	log.Printf(
		"[INFO] Retriving Edge CNAME for Account '%s': CNAME ID: %v",
		accountNumber,
		edgecnameID,
	)

	params := edgecname.NewGetEdgeCnameParams()
	params.AccountNumber = accountNumber
	params.EdgeCnameID = edgecnameID
	edgecnameObj, err := edgecnameService.GetEdgeCname(*params)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Edge CNAME: %+v", edgecnameObj)

	// Update Terraform state with retrieved Edge CNAME data
	d.Set("name", edgecnameObj.Name)
	d.Set("dir_path", edgecnameObj.DirPath)
	d.Set("enable_custom_reports", edgecnameObj.EnableCustomReports)
	d.Set("media_type_id", edgecnameObj.MediaTypeID)
	d.Set("origin_id", edgecnameObj.OriginID)
	d.Set("origin_string", edgecnameObj.OriginString)

	return diag.Diagnostics{}
}

func ResourceEdgeCnameUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	// Initialize Edge CNAME Service
	accountNumber := d.Get("account_number").(string)
	edgecnameID, _ := strconv.Atoi(d.Id())
	config := m.(**api.ClientConfig)
	edgecnameService, err := buildEdgeCnameService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Edge CNAME object from API
	getEdgecnameParams := edgecname.NewGetEdgeCnameParams()
	getEdgecnameParams.AccountNumber = accountNumber
	getEdgecnameParams.EdgeCnameID = edgecnameID
	edgecnameObj, err := edgecnameService.GetEdgeCname(*getEdgecnameParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update Edge CNAME object
	edgecnameObj.Name = d.Get("name").(string)
	edgecnameObj.DirPath = d.Get("dir_path").(string)
	edgecnameObj.EnableCustomReports = d.Get("enable_custom_reports").(int)
	edgecnameObj.MediaTypeID = d.Get("media_type_id").(int)
	edgecnameObj.OriginID = d.Get("origin_id").(int)

	// Call Update Edge CNAME API
	log.Printf(
		"[INFO] Updating Edge CNAME for Account '%s': Edge CNAME ID %d, Body %+v",
		accountNumber,
		edgecnameID,
		edgecnameObj,
	)

	updateEdgecnameParams := edgecname.NewUpdateEdgeCnameParams()
	updateEdgecnameParams.AccountNumber = accountNumber
	updateEdgecnameParams.EdgeCname = *edgecnameObj

	_, err = edgecnameService.UpdateEdgeCname(*updateEdgecnameParams)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updated Edge CNAME ID: %v", edgecnameID)

	return ResourceEdgeCnameRead(ctx, d, m)
}

func ResourceEdgeCnameDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Edge CNAME Service
	accountNumber := d.Get("account_number").(string)
	edgecnameID, _ := strconv.Atoi(d.Id())
	config := m.(**api.ClientConfig)
	edgecnameService, err := buildEdgeCnameService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Edge CNAME object from API
	getEdgecnameParams := edgecname.NewGetEdgeCnameParams()
	getEdgecnameParams.AccountNumber = accountNumber
	getEdgecnameParams.EdgeCnameID = edgecnameID
	edgecnameObj, err := edgecnameService.GetEdgeCname(*getEdgecnameParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Delete EDGE CNAME API
	log.Printf(
		"[INFO] Deleting Edge CNAME for Account '%s': CNAME ID: %v",
		accountNumber,
		edgecnameID,
	)

	deleteEdgecnameParams := edgecname.NewDeleteEdgeCnameParams()
	deleteEdgecnameParams.AccountNumber = accountNumber
	deleteEdgecnameParams.EdgeCname = *edgecnameObj
	err = edgecnameService.DeleteEdgeCname(*deleteEdgecnameParams)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleted Edge CNAME ID: %v", edgecnameID)

	d.SetId("")

	return diag.Diagnostics{}
}
