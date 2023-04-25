// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package edgecname

import (
	"context"
	"errors"
	"log"
	"strconv"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/edgecname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceEdgeCname() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceEdgeCnameCreate,
		ReadContext:   ResourceEdgeCnameRead,
		UpdateContext: ResourceEdgeCnameUpdate,
		DeleteContext: ResourceEdgeCnameDelete,
		Importer:      helper.Import(ResourceEdgeCnameRead, "account_number", "id"),

		CustomizeDiff: customdiff.ValidateChange(
			"media_type_id",
			func(
				ctx context.Context,
				oldValue,
				newValue,
				meta interface{},
			) error {
				if oldValue.(int) == 0 {
					// stop the check if this is a new resource
					return nil
				}

				// throw an error when changing media_type_id
				if oldValue.(int) != newValue.(int) {
					return errors.New(
						"media_type_id cannot be modified after creation")
				}

				// No change to media_type_id
				return nil
			},
		),

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Identifies your account. Find your account number in the upper right-hand corner of the MCC.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Identifies a hostname through which your content will be delivered. It should only consist of lower-case alphanumeric characters, dashes, and periods. From your DNS service provider, configure a CNAME record for this hostname.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"dir_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies a location on the origin server. Specify a relative path from the root folder of the origin server to the desired location. Set this argument to an empty string to point the edge CNAME to the root folder of the origin server.",
			},
			"enable_custom_reports": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Determines whether hits and data transferred statistics will be tracked for this edge CNAME. View this data through the Custom Reports module. Valid values are: \n" +
					" * `0` - Disabled (Default Value). \n" +
					" * `1` - Enabled. CDN activity on this edge CNAME will be logged.",
			},
			"media_type_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "Identifies the delivery platform on which the edge CNAME will be created. Valid values are: \n" +
					" * `3` - Http Large \n" +
					" * `8` - HTTP Small \n" +
					" * `14` - ADN",
			},
			"origin_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "Determines whether this edge CNAME will point to CDN storage or a customer origin group. Valid values are: \n" +
					" * `-1` - CDN storage \n" +
					" * `<Customer Origin Group ID>` - Customer Origin Group. Specify the system-defined ID for the desired customer origin group.",
			},
			"origin_string": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the origin identifier, the account number, and the relative path associated with the edge CNAME.",
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
	config := m.(internal.ProviderConfig)
	edgecnameService, err := buildEdgeCnameService(config)
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
	config := m.(internal.ProviderConfig)
	edgecnameService, err := buildEdgeCnameService(config)
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
	config := m.(internal.ProviderConfig)
	edgecnameService, err := buildEdgeCnameService(config)
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
	config := m.(internal.ProviderConfig)
	edgecnameService, err := buildEdgeCnameService(config)
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
