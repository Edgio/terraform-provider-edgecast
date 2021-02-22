// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"terraform-provider-vmp/vmp/api"
	"terraform-provider-vmp/vmp/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const emptyPolicyFormat string = "{\"@type\":\"policy-create\",\"name\":\"Terraform Placeholder - %s\",\"platform\":\"%s\",\"rules\":[{\"@type\":\"rule-create\",\"description\":\"Placeholder rule created by the Verizon Media Terraform Provider\",\"matches\":[{\"features\":[{\"type\":\"feature.comment\",\"value\":\"Empty policy created on %s\"}],\"ordinal\":1,\"type\":\"match.always\"}],\"name\":\"Placeholder Rule\"}],\"state\":\"locked\"}"

func resourceRulesEngineV4Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"customeruserid":    {Type: schema.TypeString, Optional: true},
			"portaltypeid":      {Type: schema.TypeString, Optional: true},
			"account_number":    {Type: schema.TypeString, Optional: true},
			"deploy_to":         {Type: schema.TypeString, Required: true},
			"deploy_request_id": {Type: schema.TypeString, Computed: true},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					policyMap := make(map[string]interface{})
					json.Unmarshal([]byte(val.(string)), &policyMap)
					// remove unneeded metadata the user may have input
					cleanPolicy(policyMap)
					jsonBytes, err := json.Marshal(policyMap)
					if err != nil {
						panic(fmt.Errorf("policy StateFunc: %v", err))
					}
					return string(jsonBytes)
				},
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy := d.Get("policy").(string)

	// messy - needs improvement - unmarshalling json, modifying, then marshalling back to string
	// state must always be locked
	policyMap := make(map[string]interface{})
	json.Unmarshal([]byte(policy), &policyMap)
	policyMap["state"] = "locked"

	policyBytes, err := json.Marshal(policyMap)
	if err != nil {
		return diag.FromErr(err)
	}

	policy = string(policyBytes)

	err = addPolicy(policy, false, d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = d.Get("account_number").(string)
	portalTypeID := d.Get("portaltypeid").(string) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	customerUserID := d.Get("customeruserid").(string)
	policyID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Retrieving policy %d", policyID)
	rulesEngineAPIClient := api.NewRulesEngineAPIClient(*config)
	policyFromAPIMap, err := rulesEngineAPIClient.GetPolicy((**config).AccountNumber, customerUserID, portalTypeID, policyID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// set id to policy id from body
	d.SetId(policyFromAPIMap["id"].(string))

	// Remove unneeded policy and rule  metadata - this metadata interferes with terraform diffs
	cleanPolicy(policyFromAPIMap)

	// convert to json
	jsonBytes, err := json.Marshal(policyFromAPIMap)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	policyFromAPI := string(jsonBytes)

	log.Printf("[INFO] Successfully retrieved policy %d: %s", policyID, policyFromAPI)
	d.Set("policy", policyFromAPI)

	return diags
}

func cleanPolicy(policyMap map[string]interface{}) {
	delete(policyMap, "id")
	delete(policyMap, "@id")
	delete(policyMap, "@type")
	delete(policyMap, "policy_type")
	delete(policyMap, "state") // will always be "locked"
	delete(policyMap, "history")
	delete(policyMap, "created_at")
	delete(policyMap, "updated_at")

	rules := policyMap["rules"].([]interface{})
	ruleMaps := make([]map[string]interface{}, 0)

	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		delete(ruleMap, "id")
		delete(ruleMap, "@id")
		delete(ruleMap, "@type")
		delete(ruleMap, "ordinal")
		delete(ruleMap, "created_at")
		delete(ruleMap, "updated_at")

		if matches, ok := ruleMap["matches"].([]interface{}); ok {
			// replace with cleaned matches
			ruleMap["matches"] = cleanMatches(matches)
		}

		ruleMaps = append(ruleMaps, ruleMap)
	}

	// replace with cleaned rules
	policyMap["rules"] = ruleMaps
}

// recursive function to remove unneeded metadata from matches
func cleanMatches(matches []interface{}) []map[string]interface{} {
	cleanedMatches := make([]map[string]interface{}, 0)

	for _, match := range matches {
		cleanedMatch := match.(map[string]interface{})
		delete(cleanedMatch, "ordinal")
		standardizeMatchFeature(cleanedMatch)

		// recursively clean child matches
		if childMatches, ok := cleanedMatch["matches"].([]interface{}); ok {
			cleanedMatch["matches"] = cleanMatches(childMatches)
		}

		if features, ok := cleanedMatch["features"].([]interface{}); ok {
			cleanedFeatures := make([]map[string]interface{}, 0)

			for _, feature := range features {
				cleanedFeature := feature.(map[string]interface{})
				delete(cleanedFeature, "ordinal")
				standardizeMatchFeature(cleanedFeature)
				cleanedFeatures = append(cleanedFeatures, cleanedFeature)
			}

			cleanedMatch["features"] = cleanedFeatures
		}

		cleanedMatches = append(cleanedMatches, cleanedMatch)
	}

	return cleanedMatches
}

// change string arrays to space-separated strings and standardize keys to hyperion standard i.e. "-" -> "_"
func standardizeMatchFeature(matchFeatureMap map[string]interface{}) {
	for k, v := range matchFeatureMap {
		delete(matchFeatureMap, k)
		// the json library unmarshals all arrays into []interface{}
		// so we have to do this roundabout way of converting to []string
		if valArray, ok := v.([]interface{}); ok {
			if stringArray, ok := helper.InterfaceArrayToStringArray(valArray); ok {
				v = strings.Join(stringArray, " ")
			}
		}
		matchFeatureMap[strings.Replace(k, "-", "_", -1)] = v
	}
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

func getDeployRequestData(d *schema.ResourceData, policyID int) *api.AddDeployRequest {
	return &api.AddDeployRequest{
		Message:     "Auto-submitted policy",
		PolicyID:    policyID,
		Environment: d.Get("deploy_to").(string),
	}
}

func addPolicy(policy string, isEmptyPolicy bool, d *schema.ResourceData, m interface{}) error {
	config := m.(**api.ClientConfig)

	customerid := d.Get("account_number").(string)
	portaltypeid := d.Get("portaltypeid").(string) //1:mcc 2:pcc 3:whole 4:uber 5:opencdn
	customeruserid := d.Get("customeruserid").(string)

	reClient := api.NewRulesEngineAPIClient(*config)

	parsedResponse, err := reClient.AddPolicy(policy, customerid, portaltypeid, customeruserid)
	if err != nil {
		return fmt.Errorf("addPolicy: %v", err)
	}

	policyID, err := strconv.Atoi(parsedResponse.ID)
	if err != nil {
		return fmt.Errorf("addPolicy: parsing policy ID: %v", err)
	}

	if !isEmptyPolicy {
		d.SetId(parsedResponse.ID)
		d.Set("policy", policy)
	}

	deployRequest := getDeployRequestData(d, policyID)
	log.Printf("[INFO] Deploying new policy for Account %s: %+v", customerid, deployRequest)

	deployResponse, deployErr := reClient.DeployPolicy(deployRequest, customerid, portaltypeid, customeruserid)

	if deployErr != nil {
		log.Printf("[WARN] Deploying new policy for Account %s failed", customerid)
		return fmt.Errorf("addPolicy: %v", deployErr)
	}

	log.Printf("[INFO] Successfully deployed new policy for Account %s: %+v", customerid, deployResponse)

	if isEmptyPolicy {
		d.SetId("") // indicates "delete" happened
	} else {
		d.Set("deploy_request_id", deployResponse.ID)
	}

	return nil
}
