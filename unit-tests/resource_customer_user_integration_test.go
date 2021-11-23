// Copyright 2021 Edgecast Inc. Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package test

import (
	"terraform-provider-ec/unit-tests/model"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test cases for storage account name conversion logic
var tcCustomerUser = map[string]model.ResourceNewCustomerUser{
	"terratest.testing.ec.customeruser": {
		CustomerUserInfo: model.NewCustomerUserInfo{
			AccountNumber: "D9127",
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "admin+1@test20252021-7.com",
			IsAdmin:       false,
		},
		Credential: model.Credentials{
			ApiToken:         "<apitoken>",
			IdsClientSecret:  "<idsclientsecret>",
			IdsClientID:      "<idssclientID>",
			IdsScope:         "<scope>",
			ApiAddress:       "<apiUrl>",
			IdsAddress:       "<idsaddress>",
			ApiAddressLegacy: "<apiaddresslegacy>",
		},
	},
}

func TestUT_CustomerUser_basic(t *testing.T) {
	t.Parallel()

	for _, input := range tcCustomerUser {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../examples/resources/ec_customer_user",
			Vars: map[string]interface{}{
				"new_admin_user": map[string]interface{}{
					"customer_account_number": input.CustomerUserInfo.AccountNumber,
					"first_name":              input.CustomerUserInfo.FirstName,
					"last_name":               input.CustomerUserInfo.LastName,
					"email":                   input.CustomerUserInfo.Email,
					"is_admin":                input.CustomerUserInfo.IsAdmin,
				},
				"credentials": map[string]interface{}{
					"api_token":          input.Credential.ApiToken,
					"ids_client_secret":  input.Credential.IdsClientSecret,
					"ids_client_id":      input.Credential.IdsClientID,
					"ids_scope":          input.Credential.IdsScope,
					"api_address":        input.Credential.ApiAddress,
					"ids_address":        input.Credential.IdsAddress,
					"api_address_legacy": input.Credential.ApiAddressLegacy,
				},
			},
		}

		// Construct the terraform options with default retryable errors to handle the most common
		// retryable errors in terraform testing.
		terraformOptions := terraform.WithDefaultRetryableErrors(t, tfOptions)
		// At the end of the test, run `terraform destroy` to clean up any resources that were created.
		defer terraform.Destroy(t, terraformOptions)

		// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
		terraform.InitAndApply(t, terraformOptions)
	}
}
