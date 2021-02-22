// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"log"
	"strconv"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCname() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCnameCreate,
		ReadContext:   resourceCnameRead,
		UpdateContext: resourceCnameUpdate,
		DeleteContext: resourceCnameDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {Type: schema.TypeString, Optional: true},
			"name":           {Type: schema.TypeString, Required: true},
			"type":           {Type: schema.TypeInt, Required: true},
			"origin_id":      {Type: schema.TypeInt, Required: true},
			"origin_type":    {Type: schema.TypeInt, Required: true},
		},
	}
}

func resourceCnameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	return resourceCnameRead(ctx, d, m)
}

func resourceCnameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceCnameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	return resourceCnameRead(ctx, d, m)
}

func resourceCnameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
