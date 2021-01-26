// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"encoding/json"
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
			"account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	policy := d.Get("policy").(string)
	// customerid := d.Get("customerid").(string)
	// portaltypeid := d.Get("portaltypeid").(string) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	// customeruserid := d.Get("customeruserid").(string)
	InfoLogger.Printf("policy: %s\n", policy)

	// config := m.(*ProviderConfiguration)

	// reClient := api.NewRulesEngineApiClient(config.APIClient)

	// parsedResponse, err := reClient.AddPolicy(policy, customerid, portaltypeid, customeruserid)

	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// //InfoLogger.Printf("policy: %s\n", res)
	// d.SetId(parsedResponse.Id)
	d.SetId("20")
	d.Set("policy", policy)

	return diags
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	InfoLogger.Printf("account number: %s\n", d.Get("account_number").(string))

	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	policy := d.Get("policy").(string)
	portalTypeID := d.Get("portaltypeid").(string) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	customerUserID := d.Get("customeruserid").(string)
	InfoLogger.Printf("policy: %s\n", policy)

	rulesEngineAPIClient := api.NewRulesEngineApiClient(providerConfiguration.APIClient)

	policyID, _ := strconv.Atoi(d.Id())
	InfoLogger.Printf("Policy ID is %d \n", policyID)

	var customerID int

	parsedCustomerID, parseErr := strconv.ParseInt(providerConfiguration.AccountNumber, 16, 32)

	if parseErr == nil {
		customerID = int(parsedCustomerID)
	} else {
		d.SetId("")
		return diag.FromErr(parseErr)
	}

	policyMap, err := rulesEngineAPIClient.GetPolicy(customerID, customerUserID, portalTypeID, policyID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// set id to policy id from body
	d.SetId(policyMap["id"].(string))

	// Remove unneeded policy metadata
	delete(policyMap, "id")
	delete(policyMap, "@id")
	delete(policyMap, "@type")
	delete(policyMap, "history")
	delete(policyMap, "created_at")
	delete(policyMap, "updated_at")

	rules := policyMap["rules"].([]interface{})
	ruleMaps := make([]map[string]interface{}, 0)

	// Remove unneeded rule metadata
	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		delete(ruleMap, "id")
		delete(ruleMap, "@id")
		delete(ruleMap, "@type")
		delete(ruleMap, "created_at")
		delete(ruleMap, "updated_at")
		ruleMaps = append(ruleMaps, ruleMap)
	}

	// replace rules with cleaned rules
	policyMap["rules"] = ruleMaps

	// convert to json
	jsonBytes, marshalErr := json.Marshal(policyMap)

	if marshalErr != nil {
		d.SetId("")
		return diag.FromErr(marshalErr)
	}

	json := string(jsonBytes)
	InfoLogger.Printf("Retrieved policy: %s\n", json)
	d.Set("policy", json)

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
