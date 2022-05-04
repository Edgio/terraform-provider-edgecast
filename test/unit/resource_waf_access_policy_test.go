// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package test

import (
	"terraform-provider-edgecast/test/unit/helper"
	"terraform-provider-edgecast/test/unit/model"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestUT_AccessRule_basic(t *testing.T) {
	t.Parallel()
	t.Skip("test is not ready for unit testing")

	// // Test cases for storage account name conversion logic
	tc, err := getTestCases()
	if err != nil {
		t.Errorf("Reading credential_ucc.json file error:%s", err)
	}

	for _, input := range *tc {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../../examples/resources/ec_waf_access_rule",
			Vars: map[string]interface{}{
				"credentials": map[string]interface{}{
					"api_token":          input.ApiToken,
					"ids_client_secret":  input.IdsClientSecret,
					"ids_client_id":      input.IdsClientID,
					"ids_scope":          input.IdsScope,
					"api_address":        input.ApiAddress,
					"api_address_legacy": input.ApiAddressLegacy,
					"ids_address":        input.IdsAddress,
				},
			},
		}

		// Construct the terraform options with default retryable errors to handle the most common
		// retryable errors in terraform testing.
		terraformOptions := terraform.WithDefaultRetryableErrors(t, tfOptions)
		// At the end of the test, run `terraform destroy` to clean up any resources that were created.
		//defer terraform.Destroy(t, terraformOptions)

		// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
		terraform.InitAndApply(t, terraformOptions)

	}
}

func getTestCases() (*map[string]model.Credentials, error) {
	tc := make(map[string]model.Credentials)
	credential := model.Credentials{}
	err := helper.ReadCredentialJsonfile("credential_ucc.json", &credential)
	if err != nil {
		return nil, err
	}

	tc["terratest.testing.ec.acl1"] = credential

	return &tc, nil
}
