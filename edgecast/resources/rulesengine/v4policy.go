// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package rulesengine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/rulesengine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	emptyPolicyFormat string = "{\"@type\":\"policy-create\",\"name\":\"Terraform Placeholder - %s\",\"platform\":\"%s\",\"rules\":[{\"@type\":\"rule-create\",\"description\":\"Placeholder rule created by the Edgecast Terraform Provider\",\"matches\":[{\"features\":[{\"type\":\"feature.comment\",\"value\":\"Empty policy created on %s\"}],\"ordinal\":1,\"type\":\"match.always\"}],\"name\":\"Placeholder Rule\"}],\"state\":\"locked\"}"
	jsonKeyFeatures   string = "features"
	jsonkeyMatches    string = "matches"
)

func ResourceRulesEngineV4Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourcePolicyCreate,
		ReadContext:   ResourcePolicyRead,
		UpdateContext: ResourcePolicyUpdate,
		DeleteContext: ResourcePolicyDelete,
		Importer:      helper.Import(ResourcePolicyRead, "account_number", "id", "portaltypeid", "customeruserid", "ownerid"),

		Schema: map[string]*schema.Schema{
			"customeruserid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reserved for future use.",
			},
			"portaltypeid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reserved for future use.",
			},
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reserved for future use.",
			},
			"ownerid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Required when acting on behalf of a customer and using Wholesaler or Partner credentials. This value should be the customer Account Number in the upper right-hand corner of the MCC.",
			},
			"deploy_to": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Identifies the environment to which the policy will be deployed. Valid values are: \n\n" +
					"        production | staging",
				ValidateFunc: validation.StringInSlice(
					[]string{"production", "staging"},
					false),
			},
			"deploy_request_id": {
				Type:        schema.TypeString,
				Description: "Indicates the system-defined ID for the policy's deploy request.",
				Computed:    true,
			},
			"policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Defines the policy, in JSON format, that will be deployed.",
				StateFunc:   cleanPolicyForTerrafomState,
				ValidateFunc: validation.All(
					validation.StringIsNotWhiteSpace,
					validation.StringIsJSON,
					helper.StringIsNotEmptyJSON,
				),
				DiffSuppressFunc: policyDiffSuppress,
			},
		},
	}
}

// ResourcePolicyCreate - Create a new policy and deploy it to a target platform
func ResourcePolicyCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	policy := d.Get("policy").(string)

	// messy - needs improvement - unmarshalling json, modifying, then
	// marshalling back to string state must always be locked
	policyMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(policy), &policyMap)
	if err != nil {
		return diag.Errorf("error reading policy: %s", err.Error())
	}

	policyMap["state"] = "locked"
	if policyMap["name"] != nil {
		log.Printf(`
		[WARN] Please remove policy name, "%s" in the policy.
		It will be ignored. 
		The policy name would be auto-generated by the provider with following format.
		tf-[customer_account_id]-[deploy_to]-[platform]-[timestamp(utc)]`,
			policyMap["name"])
	}
	policyMap["name"] = fmt.Sprintf("tf-%s-%s-%s-%d",
		d.Get("account_number").(string),
		d.Get("deploy_to").(string),
		policyMap["platform"],
		time.Now().UTC().Unix())

	policyBytes, err := json.Marshal(policyMap)
	if err != nil {
		return diag.FromErr(err)
	}

	policy = string(policyBytes)

	err = addPolicy(policy, false, d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	return ResourcePolicyRead(ctx, d, m)
}

// ResourcePolicyRead reads the current policy
func ResourcePolicyRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	policy, err := getPolicy(m, d)
	log.Printf("[INFO] policy : %+v", policy)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// set id to policy id from body
	d.SetId(policy["id"].(string))

	// Remove unneeded policy and rule metadata - this metadata interferes with
	// terraform diffs
	err = cleanPolicy(policy)
	if err != nil {
		d.SetId("")
		return diag.FromErr(
			fmt.Errorf("error cleaning policy : %w", err))
	}

	// convert to json
	jsonBytes, err := json.Marshal(policy)

	if err != nil {
		d.SetId("")
		return diag.FromErr(
			fmt.Errorf("error marshaling policy to json : %v", err))
	}

	policyAsString := string(jsonBytes)
	log.Printf(
		"[INFO] Successfully retrieved policy %s: %s",
		d.Id(),
		policyAsString)

	d.Set("policy", policyAsString)

	return diag.Diagnostics{}
}

// ResourcePolicyUpdate adds/deletes/updates rules in the current policy and
// deploys the modified policy to a target platform
func ResourcePolicyUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return ResourcePolicyCreate(ctx, d, m)
}

// ResourcePolicyDelete creates a new empty placeholder policy and deploys it to
// a target platform instead of actual deletion.
func ResourcePolicyDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// We will retrieve a fresh copy of the policy to prevent
	// sending an empty policy to the wrong platform
	policy, err := getPolicy(m, d)

	if err != nil {
		return diag.FromErr(err)
	}

	// pull out platform from existing policy
	platform := policy["platform"].(string)

	// You can't actually delete policies, so we will instead create a
	// placeholder empty policy for the customer for the given platform and
	// policy type
	timestamp := time.Now().Format(time.RFC3339)
	emptyPolicyJSON := fmt.Sprintf(
		emptyPolicyFormat,
		timestamp,
		platform,
		timestamp)

	err = addPolicy(emptyPolicyJSON, true, d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func cleanPolicy(policyMap map[string]interface{}) error {
	delete(policyMap, "id")
	delete(policyMap, "@id")
	delete(policyMap, "@type")
	delete(policyMap, "policy_type")
	delete(policyMap, "state") // will always have the value "locked"
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
			temp, err := cleanMatches(matches)
			if err != nil {
				return fmt.Errorf("error standardizing match: %w", err)
			}
			ruleMap["matches"] = temp
		}

		ruleMaps = append(ruleMaps, ruleMap)

	}

	// replace with cleaned rules
	policyMap["rules"] = ruleMaps
	return nil
}

// recursive function to remove unneeded metadata from matches
func cleanMatches(matches []interface{}) ([]map[string]interface{}, error) {
	cleanedMatches := make([]map[string]interface{}, 0)

	for _, match := range matches {
		cleanedMatch := match.(map[string]interface{})

		// Hacky workaround for a Rules Engine API bug:
		// API returns "select.first-match" as "match.select.first-match.
		// We will adjust type here, but needs to be fixed in the API.
		if cleanedMatch["type"] == "match.select.first-match" {
			cleanedMatch["type"] = "select.first-match"
		}

		delete(cleanedMatch, "ordinal")
		delete(cleanedMatch, "raw_value")

		err := standardizeMatchFeature(cleanedMatch)
		if err != nil {
			return nil, fmt.Errorf("error standardizing match: %w", err)
		}

		// recursively clean child matches
		if childMatches, ok := cleanedMatch["matches"].([]interface{}); ok {
			temp, err := cleanMatches(childMatches)
			if err != nil {
				return nil, fmt.Errorf("error cleaning match: %w", err)
			}
			cleanedMatch["matches"] = temp
		}

		if features, ok := cleanedMatch["features"].([]interface{}); ok {
			cleanedFeatures := make([]map[string]interface{}, 0)
			for _, feature := range features {
				cleanedFeature := feature.(map[string]interface{})

				delete(cleanedFeature, "ordinal")
				delete(cleanedFeature, "raw_source")
				delete(cleanedFeature, "raw_destination")
				delete(cleanedFeature, "raw_value")

				err := standardizeMatchFeature(cleanedFeature)
				if err != nil {
					return nil, fmt.Errorf("error standardizing features: %w", err)
				}
				cleanedFeatures = append(cleanedFeatures, cleanedFeature)
			}

			cleanedMatch["features"] = cleanedFeatures
		}

		cleanedMatches = append(cleanedMatches, cleanedMatch)
	}

	return cleanedMatches, nil
}

// change string arrays to space-separated strings and standardize keys to
// hyperion standard i.e. "-" -> "_"
func standardizeMatchFeature(matchFeatureMap map[string]interface{}) error {
	for k, v := range matchFeatureMap {
		delete(matchFeatureMap, k)

		// the json library unmarshals all arrays into []interface{}
		// so we have to do this roundabout way of converting to []string
		if valArray, ok := v.([]interface{}); ok {
			if k != jsonKeyFeatures && k != jsonkeyMatches {
				stringArray, _ := helper.ConvertSliceToStrings(valArray)
				v = strings.Join(stringArray, " ")
			}
		}
		matchFeatureMap[strings.Replace(k, "-", "_", -1)] = v
	}
	return nil
}

func getDeployRequestData(
	d *schema.ResourceData,
	policyID int,
) *rulesengine.SubmitDeployRequest {
	return &rulesengine.SubmitDeployRequest{
		Message:     "Auto-submitted policy",
		PolicyID:    policyID,
		Environment: d.Get("deploy_to").(string),
	}
}

func getPolicy(
	m interface{},
	d *schema.ResourceData,
) (map[string]interface{}, error) {
	// Retrieve data needed by API call
	config := m.(internal.ProviderConfig)
	accountNumber := d.Get("account_number").(string)
	portalTypeID := d.Get("portaltypeid").(string) // 1=MCC 2=PCC 3=WCC 4=UCC
	customerUserID := d.Get("customeruserid").(string)
	ownerID := d.Get("ownerid").(string)
	policyID, err := strconv.Atoi(d.Id())

	if err != nil {
		return nil,
			fmt.Errorf("error parsing Policy ID from state file: %v", err)
	}

	log.Printf("[INFO] Retrieving policy %d", policyID)

	// Initialize Rules Engine Service
	rulesengineService, err := buildRulesEngineService(config)
	if err != nil {
		d.SetId("")
		return nil, fmt.Errorf("addPolicy: buildRulesEngineService: %v", err)
	}

	// Call Add Policy API
	params := rulesengine.NewGetPolicyParams()
	params.AccountNumber = accountNumber
	params.CustomerUserID = customerUserID
	params.PortalTypeID = portalTypeID
	params.PolicyID = policyID
	params.OwnerID = ownerID

	return rulesengineService.GetPolicy(*params)
}

func addPolicy(
	policy string,
	isEmptyPolicy bool,
	d *schema.ResourceData,
	m interface{},
) error {
	// Retrieve data needed by API calls
	config := m.(internal.ProviderConfig)
	accountNumber := d.Get("account_number").(string)
	customerUserID := d.Get("customeruserid").(string)
	portalTypeID := d.Get("portaltypeid").(string) // 1=MCC 2=PCC 3=WCC 4=UCC
	ownerID := d.Get("ownerid").(string)

	// Initialize Rules Engine Service
	rulesengineService, err := buildRulesEngineService(config)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("addPolicy: buildRulesEngineService: %v", err)
	}
	// Call Add Policy API
	policyParams := rulesengine.NewAddPolicyParams()
	policyParams.AccountNumber = accountNumber
	policyParams.CustomerUserID = customerUserID
	policyParams.PortalTypeID = portalTypeID
	policyParams.PolicyAsString = policy
	policyParams.OwnerID = ownerID

	parsedResponse, err := rulesengineService.AddPolicy(*policyParams)
	if err != nil {
		return fmt.Errorf("addPolicy: %v\n%v", policyParams, err)
	}

	// Process response data and prepare Deploy Request
	policyID, err := strconv.Atoi(parsedResponse.ID)
	if err != nil {
		return fmt.Errorf("addPolicy: parsing policy ID: %v", err)
	}

	if !isEmptyPolicy {
		d.SetId(parsedResponse.ID)
		d.Set("policy", policy)
	}

	deployRequest := getDeployRequestData(d, policyID)
	log.Printf(
		"[INFO] Deploying new policy for Account %s: %+v",
		accountNumber,
		deployRequest,
	)

	// Call Submit Deploy Request API
	deployRequestParams := rulesengine.NewSubmitDeployRequestParams()
	deployRequestParams.AccountNumber = accountNumber
	deployRequestParams.CustomerUserID = customerUserID
	deployRequestParams.PortalTypeID = portalTypeID
	deployRequestParams.DeployRequest = *deployRequest
	deployRequestParams.OwnerID = ownerID

	deployResponse, deployErr := rulesengineService.SubmitDeployRequest(
		*deployRequestParams,
	)

	if deployErr != nil {
		log.Printf(
			"[WARN] Deploying new policy for Account %s failed",
			accountNumber)
		return fmt.Errorf("addPolicy: %v\n%v", policyParams, deployErr)
	}

	log.Printf(
		"[INFO] Successfully deployed new policy for Account %s: %+v",
		accountNumber,
		deployResponse)

	if isEmptyPolicy {
		d.SetId("") // indicates "delete" happened
	} else {
		d.Set("deploy_request_id", deployResponse.ID)
	}

	return nil
}

func cleanPolicyForTerrafomState(val interface{}) string {
	policy := val.(string)
	if len(policy) == 0 {
		return policy
	}
	policyMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(policy), &policyMap)
	if err != nil {
		panic(fmt.Errorf("cleanPolicyForTerrafomState: %w", err))
	}

	// remove unneeded metadata the user may have input
	err = cleanPolicy(policyMap)
	if err != nil {
		panic(fmt.Errorf("cleanPolicyForTerrafomState: %w", err))
	}

	jsonBytes, err := json.Marshal(policyMap)
	if err != nil {
		panic(fmt.Errorf("cleanPolicyForTerrafomState: %w", err))
	}

	return string(jsonBytes)
}

func policyDiffSuppress(k, old, new string, _ *schema.ResourceData) bool {
	oldPolicy := make(map[string]interface{})
	newPolicy := make(map[string]interface{})

	_ = json.Unmarshal([]byte(old), &oldPolicy)
	_ = json.Unmarshal([]byte(new), &newPolicy)

	// Ignore policy name changes.
	delete(oldPolicy, "name")
	delete(newPolicy, "name")

	// Empty and null rule names are equal - just clear them out.
	deleteEmptyRuleNames(oldPolicy)
	deleteEmptyRuleNames(newPolicy)

	return reflect.DeepEqual(oldPolicy, newPolicy)
}

// deleteEmptyRuleNames is used to remove empty or null rule names before a
// policy diff comparison. They interfere with terraform diffs.
func deleteEmptyRuleNames(policy map[string]interface{}) {
	rulesRaw, ok := policy["rules"]

	if !ok || rulesRaw == nil {
		return
	}

	rules := rulesRaw.([]interface{})
	cleanedRules := make([]map[string]interface{}, 0)

	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})

		val, ok := ruleMap["name"]

		if ok && len(val.(string)) == 0 {
			delete(ruleMap, "name")
		}

		cleanedRules = append(cleanedRules, ruleMap)
	}

	// replace with cleaned rules
	policy["rules"] = cleanedRules
}
