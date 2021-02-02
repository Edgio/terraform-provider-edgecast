// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const emptyPolicyFormat string = "{\"@type\":\"policy-create\",\"name\":\"Terraform Placeholder - %s\",\"platform\":\"%s\",\"rules\":[{\"@type\":\"rule-create\",\"description\":\"placeholder rule created by the Verizon Media Terraform Provider\",\"matches\":[{\"features\":[{\"type\":\"feature.comment\",\"value\":\"empty policy created on %s\"}],\"ordinal\":1,\"type\":\"match.always\"}],\"name\":\"placeholder rule\"}],\"state\":\"locked\"}"

func resourceRulesEngineV4Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"customeruserid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"portaltypeid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"deploy_to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"deploy_request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	policy := d.Get("policy").(string)

	err := addPolicy(policy, false, d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	policy := d.Get("policy").(string)
	portalTypeID := d.Get("portaltypeid").(string) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	customerUserID := d.Get("customeruserid").(string)
	InfoLogger.Printf("user input policy: %s\n", policy)

	rulesEngineAPIClient := api.NewRulesEngineApiClient(providerConfiguration.APIClient)

	policyID, _ := strconv.Atoi(d.Id())
	InfoLogger.Printf("Policy ID is %d \n", policyID)

	customerID, err := parseCustomerID(providerConfiguration.AccountNumber)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
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

	InfoLogger.Printf("Retrieved policy from API: %s\n", json)
	d.Set("policy", json)

	return diags
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePolicyCreate(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	policy := d.Get("policy").(string)

	// pull out platform from existing policy
	policyMap := make(map[string]interface{})
	json.Unmarshal([]byte(policy), &policyMap)

	platform := policyMap["platform"].(string)

	// You can't actually delete policies, so we will instead
	// create a placeholder empty policy for the customer for the given platform and policy type
	timestamp := time.Now().Format(time.RFC3339)
	emptyPolicyJSON := fmt.Sprintf(emptyPolicyFormat, timestamp, platform, timestamp)

	err := addPolicy(emptyPolicyJSON, true, d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func parseCustomerID(accountNumber string) (int, error) {
	parsedCustomerID, parseErr := strconv.ParseInt(accountNumber, 16, 32)

	if parseErr == nil {
		return int(parsedCustomerID), nil
	}

	return 0, parseErr
}

func getDeployRequestData(d *schema.ResourceData, policyId int) *api.AddDeployRequest {

	return &api.AddDeployRequest{
		Message:     "Auto-submitted policy",
		PolicyId:    policyId,
		Environment: d.Get("deploy_to").(string),
	}
}

func addPolicy(policy string, isEmptyPolicy bool, d *schema.ResourceData, m interface{}) error {
	customerid := d.Get("account_number").(string)
	portaltypeid := d.Get("portaltypeid").(string) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	customeruserid := d.Get("customeruserid").(string)

	InfoLogger.Printf("addPolicy >> policy: %s\n", policy)

	config := m.(*ProviderConfiguration)

	reClient := api.NewRulesEngineApiClient(config.APIClient)

	parsedResponse, addPolicyErr := reClient.AddPolicy(policy, customerid, portaltypeid, customeruserid)
	if addPolicyErr != nil {
		return addPolicyErr
	}

	policyId, parseErr := strconv.Atoi(parsedResponse.Id)
	if addPolicyErr != nil {
		return parseErr
	}

	if !isEmptyPolicy {
		d.SetId(parsedResponse.Id)
		d.Set("policy", policy)
	}

	InfoLogger.Printf("addPolicy >> PolicyCreateResponse >> policyId: %d\n", policyId)

	payload := getDeployRequestData(d, policyId)
	InfoLogger.Printf("addPolicy >> PolicyCreateResponse >> DeployRequest >> payload: %+v\n", payload)
	deployResponse, deployErr := reClient.DeployPolicy(payload, customerid, portaltypeid, customeruserid)

	if deployErr != nil {
		return deployErr
	}

	InfoLogger.Printf("resourcePolicyCreate >> PolicyCreateResponse >> DeployRequest >> deployrequestId: %+v\n", deployResponse)

	if isEmptyPolicy {
		d.SetId("") // indicates "delete" happened
	} else {
		d.Set("deploy_request_id", deployResponse.Id)
	}

	return nil
}
