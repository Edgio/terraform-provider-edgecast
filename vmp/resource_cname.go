// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"log"
	"os"
	"strconv"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func resourceCname() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCnameCreate,
		ReadContext:   resourceCnameRead,
		UpdateContext: resourceCnameUpdate,
		DeleteContext: resourceCnameDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"origin_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"origin_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceCnameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	InfoLogger.Printf("account_number: %s\n", accountNumber)
	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(accountNumber)

	if err != nil {
		ErrorLogger.Println(err)
		return diag.FromErr(err)
	}

	InfoLogger.Printf("origin_id: %d\n", d.Get("origin_id").(int))
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

	updateCnameRequest := &api.UpdateCnameRequest{
		Name:        d.Get("name").(string),
		MediaTypeId: d.Get("type").(int),
		OriginId:    d.Get("origin_id").(int),
		OriginType:  d.Get("origin_type").(int),
	}

	cnameAPIClient := api.NewCnameApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	cnameId, _ := strconv.Atoi(d.Id())
	InfoLogger.Printf("CnameId to be updated %d", cnameId)

	parsedResponse, err := cnameAPIClient.UpdateCname(updateCnameRequest, cnameId)

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
