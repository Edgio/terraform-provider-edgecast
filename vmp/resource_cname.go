// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
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
			"account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"origin_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"origin_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceCnameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(accountNumber)

	if err != nil {
		return diag.FromErr(err)
	}

	addCnameRequest := &api.AddCnameRequest{
		Name:        d.Get("name").(string),
		MediaTypeId: d.Get("type").(int),
		OriginId:    d.Get("origin_id").(int),
		OriginType:  d.Get("origin_type").(int),
	}

	cnameAPIClient := api.NewCnameApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	parsedResponse, err := cnameAPIClient.AddCname(addCnameRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(parsedResponse.CnameId))

	return resourceCnameRead(ctx, d, m)
}

func resourceCnameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	cnameAPIClient := api.NewCnameApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	cnameID, _ := strconv.Atoi(d.Id())

	parsedResponse, err := cnameAPIClient.GetCname(cnameID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("name", parsedResponse.Name)
	d.Set("origin_id", parsedResponse.OriginId)
	d.Set("origin_string", parsedResponse.OriginString)

	return diags
}

func resourceCnameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	addCnameRequest := &api.AddCnameRequest{
		Name:        d.Get("name").(string),
		MediaTypeId: d.Get("type").(int),
		OriginId:    d.Get("origin_id").(int),
		OriginType:  d.Get("origin_type").(int),
	}

	cnameAPIClient := api.NewCnameApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	parsedResponse, err := cnameAPIClient.AddCname(addCnameRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(parsedResponse.CnameId))

	return resourceCnameRead(ctx, d, m)
}

func resourceCnameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	cnameAPIClient := api.NewCnameApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	cnameID, _ := strconv.Atoi(d.Id())

	err = cnameAPIClient.DeleteCname(cnameID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
