// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"strconv"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRulesEngineV4Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"policyid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"customerid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"customeruserid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"portaltypeid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	policy := d.Get("policy").(string)
	customerid := d.Get("customerid").(string)
	portaltypeid := d.Get("portaltypeid").(string) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	customeruserid := d.Get("customeruserid").(string)
	InfoLogger.Printf("policy: %s\n", policy)

	config := m.(*ProviderConfiguration)

	reClient := api.NewRulesEngineApiClient(config.APIClient)

	parsedResponse, err := reClient.AddPolicy(policy, customerid, portaltypeid, customeruserid)

	if err != nil {
		return diag.FromErr(err)
	}
	//InfoLogger.Printf("policy: %s\n", res)
	d.SetId(parsedResponse.Id)

	return diags
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfiguration := m.(*ProviderConfiguration)
	policy := d.Get("policy").(string)
	portaltypeid := d.Get("portaltypeid").(int) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	customeruserid := d.Get("customeruserid").(int)
	InfoLogger.Printf("policy: %s\n", policy)

	rulesEngineAPIClient := api.NewRulesEngineApiClient(providerConfiguration.APIClient)

	policyID, _ := strconv.Atoi(d.Id())

	var customerID int

	parsedCustomerID, parseErr := strconv.ParseInt(providerConfiguration.AccountNumber, 16, 32)

	if parseErr == nil {
		customerID = int(parsedCustomerID)
	} else {
		d.SetId("")
		return diag.FromErr(parseErr)
	}

	parsedResponse, err := rulesEngineAPIClient.GetPolicy(customerID, policyID)
	InfoLogger.Printf("%+v\n", parsedResponse)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// Set properties
	d.SetId(parsedResponse.Id)

	return diags
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
