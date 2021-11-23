// Copyright 2021 Edgecast Inc. Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package cname_origin

import (
	"context"
	"log"
	"strconv"

	"terraform-provider-ec/ec/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCname() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceCnameCreate,
		ReadContext:   ResourceCnameRead,
		UpdateContext: ResourceCnameUpdate,
		DeleteContext: ResourceCnameDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration."},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Sets the name that will be assigned to the edge CNAME. It should only contain lower-case " +
					"alphanumeric characters, dashes, and periods. The name specified for this parameter should also be " +
					"defined as a CNAME record on a DNS server. The CNAME record defined on the DNS server should point " +
					"to the CDN hostname (e.g., wpc.0001.edgecastcdn.net) for the platform identified by the `platform` " +
					"parameter"},
			// TODO: 'type' parameter should be changed to 'media_type' to be consistent with resource_origin
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifies the Delivery Platform on which the edge CNAME will be created. 3:Http Large, 8:HTTP Small, 14: ADN",
			},
			"origin_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "Identifies whether an edge CNAME will be created for a CDN origin server " +
					"or a customer origin server. Valid values: -1: Indicates that you would like to create " +
					"an edge CNAME for our CDN storage service, CustomerOriginID: Specifying an ID for an " +
					"existing customer origin configuration indicates that you would like to create an " +
					"edge CNAME for that customer origin",
			},
			"origin_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     80,
				Description: "Indicates the type of origin server the CNAME is created on. Default: `80` to indicate a Customer Origin",
			},
		},
	}
}

func ResourceCnameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	addCnameRequest := &api.AddCnameRequest{
		Name:        d.Get("name").(string),
		MediaTypeID: d.Get("type").(int),
		OriginID:    d.Get("origin_id").(int),
		OriginType:  d.Get("origin_type").(int),
	}

	log.Printf("[INFO] Creating CNAME for Account '%s': %+v", accountNumber, addCnameRequest)

	cnameAPIClient := api.NewCnameAPIClient(*config)

	resp, err := cnameAPIClient.AddCname(addCnameRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Create successful - New CNAME ID: %d", resp.CnameID)
	d.SetId(strconv.Itoa(resp.CnameID))

	return ResourceCnameRead(ctx, d, m)
}

func ResourceCnameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	cnameAPIClient := api.NewCnameAPIClient(*config)

	cnameID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Retrieving CNAME ID: %d", cnameID)

	resp, err := cnameAPIClient.GetCname(cnameID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved CNAME: %+v", resp)

	d.Set("name", resp.Name)
	d.Set("origin_id", resp.OriginID)
	d.Set("origin_string", resp.OriginString)

	return diags
}

func ResourceCnameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	updateCnameRequest := &api.UpdateCnameRequest{
		Name:        d.Get("name").(string),
		MediaTypeID: d.Get("type").(int),
		OriginID:    d.Get("origin_id").(int),
		OriginType:  d.Get("origin_type").(int),
	}

	cnameAPIClient := api.NewCnameAPIClient(*config)

	cnameID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Updating CNAME ID=%d: %+v", cnameID, updateCnameRequest)

	_, err := cnameAPIClient.UpdateCname(updateCnameRequest, cnameID)

	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceCnameRead(ctx, d, m)
}

func ResourceCnameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	cnameAPIClient := api.NewCnameAPIClient(*config)

	cnameID, _ := strconv.Atoi(d.Id())

	err := cnameAPIClient.DeleteCname(cnameID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
