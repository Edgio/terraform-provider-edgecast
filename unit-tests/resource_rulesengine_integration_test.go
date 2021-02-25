// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test cases for storage account name conversion logic
var tsRulesEngine = map[string]ResourceREV4{
	"terratest.testing.vmp.rulesengine": {
		policy: `
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
		rulesEngineEnvironment: "staging",
		credentials: Credentials{
			apitoken:        "",
			idsclientsecret: "CDbbMJw7FFJ11a7433ti1l9XgJHKr2Wk",
			idsclientID:     "31ef8e8f-0120-4112-8554-3eb11e83d58b",
			idsscope:        "ec.rules",
			apiaddress:      "http://dev-api.edgecast.com",
			idsaddress:      "https://id-dev.vdms.io",
		},
		testcustomerinfo: CustomerInfo{
			accountnumber:  "C1B6",
			customeruserID: "133172",
			portaltypeID:   1,
		},
	},
}

func TestUT_RulesEngine_basic(t *testing.T) {
	t.Parallel()

	for expected, input := range tsRulesEngine {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../examples/resource_rulesengine",
			Vars: map[string]interface{}{
				"policy": strings.Replace(input.policy, "$UUID$", uuid.New().String(), -1),
				"test_customer_info": map[string]interface{}{
					"account_number": input.testcustomerinfo.accountnumber,
					"customeruserid": input.testcustomerinfo.customeruserID,
					"portaltypeid":   input.testcustomerinfo.portaltypeID,
				},
				"credentials": map[string]interface{}{
					"api_token":         input.credentials.apitoken,
					"ids_client_secret": input.credentials.idsclientsecret,
					"ids_client_id":     input.credentials.idsclientID,
					"ids_scope":         input.credentials.idsscope,
					"api_address":       input.credentials.apiaddress,
					"ids_address":       input.credentials.idsaddress,
				},
				"rulesEngineEnvironment": input.rulesEngineEnvironment,
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
