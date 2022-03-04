// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package test

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-edgecast/unit-tests/model"
	"testing"

	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test cases for storage account name conversion logic
var tsRulesEngine = map[string]model.ResourceREV4{
	"terratest.testing.edgecast.rulesengine": {
		Policy: `
		        {
		            "name": "test policy-$UUID$",
		            "description": "This is a test policy of PolicyCreate.",
		            "platform": "adn",
		            "rules": [
		                {
		                    "name": "rule $UUID$",
		                    "description": "This is a test rule3.",
		                    "matches": [
		                        {
		                            "type": "match.always",
		                            "features": [
		                                {
		                                    "type": "feature.comment",
		                                    "value": "test $UUID$"
		                                }
		                            ]
		                        }
		                    ]
		                }
		            ]
		        }  
                `,
		RulesEngineEnvironment: "staging",
		Credentials: model.Credentials{
			ApiToken:         "<apitoken>",
			IdsClientSecret:  "<idsclientsecret>",
			IdsClientID:      "<idssclientID>",
			IdsScope:         "<scope>",
			ApiAddress:       "<apiUrl>",
			IdsAddress:       "<idsaddress>",
			ApiAddressLegacy: "<apiaddresslegacy>",
		},
		TestCustomerInfo: model.CustomerInfo{
			AccountNumber:  "C1B6",
			CustomerUserID: "133172",
			PortalTypeID:   1,
		},
	},
}

func TestUT_RulesEngine_basic(t *testing.T) {
	t.Parallel()

	for expected, input := range tsRulesEngine {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../examples/resources/edgecast_rules_engine_policy",
			Vars: map[string]interface{}{
				"policy": strings.Replace(input.Policy, "$UUID$", uuid.New().String(), -1),
				"test_customer_info": map[string]interface{}{
					"account_number": input.TestCustomerInfo.AccountNumber,
					"customeruserid": input.TestCustomerInfo.CustomerUserID,
					"portaltypeid":   input.TestCustomerInfo.PortalTypeID,
				},
				"credentials": map[string]interface{}{
					"api_token":          input.Credentials.ApiToken,
					"ids_client_secret":  input.Credentials.IdsClientSecret,
					"ids_client_id":      input.Credentials.IdsClientID,
					"ids_scope":          input.Credentials.IdsScope,
					"api_address":        input.Credentials.ApiAddress,
					"ids_address":        input.Credentials.IdsAddress,
					"api_address_legacy": input.Credentials.ApiAddressLegacy,
				},
				"rulesEngineEnvironment": input.RulesEngineEnvironment,
			},
		}

		// Construct the terraform options with default retryable errors to handle the most common
		// retryable errors in terraform testing.
		terraformOptions := terraform.WithDefaultRetryableErrors(t, tfOptions)
		// At the end of the test, run `terraform destroy` to clean up any resources that were created.
		defer terraform.Destroy(t, terraformOptions)

		// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
		terraform.InitAndApply(t, terraformOptions)

		// Run `terraform output` to get the IP of the instance
		policyID := terraform.Output(t, terraformOptions, "policy_id")
		fmt.Print("\n=========================================================\n")
		fmt.Printf("Checking %s...\n", expected)
		fmt.Print("\n---------------------------------------------------------\n")
		fmt.Printf("Policy ID created: %s\n", policyID)
		fmt.Print("\n=========================================================\n")
		iPolicyID, err := strconv.Atoi(strings.Replace(strings.Replace(policyID, "[", "", -1), "]", "", -1))
		if err != nil || iPolicyID == 0 {
			t.Fatalf("invalid policyID >> %s", err)
		}
	}
}
