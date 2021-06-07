// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package test

import (
	"terraform-provider-vmp/unit-tests/model"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test cases for storage account name conversion logic
var tcCustomerUser = map[string]model.ResourceNewCustomerUser{
	"terratest.testing.vmp.customer": {
		CustomerUserInfo: model.NewCustomerUserInfo{
			AccountNumber: "D9127",
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "admin+1@test20252021-7.com",
			IsAdmin:       false,
		},
		Credential: model.Credentials{
			ApiToken:        "AULdReDoB3gb0D7LNTx857NQvrcIKyvL",
			IdsClientSecret: "CDbbMJw7FFJ11a7433ti1l9XgJHKr2Wk",
			IdsClientID:     "31ef8e8f-0120-4112-8554-3eb11e83d58b",
			IdsScope:        "ec.rules",
			ApiAddress:      "http://dev-api.edgecast.com",
			IdsAddress:      "https://id-dev.vdms.io",
		},
	},
}

func TestUT_CustomerUser_basic(t *testing.T) {
	t.Parallel()

	for _, input := range tcCustomerUser {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../examples/resources/vmp_customer_user",
			Vars: map[string]interface{}{
				"new_admin_user": map[string]interface{}{
					"customer_account_number": input.CustomerUserInfo.AccountNumber,
					"first_name":              input.CustomerUserInfo.FirstName,
					"last_name":               input.CustomerUserInfo.LastName,
					"Email":                   input.CustomerUserInfo.Email,
					"is_admin":                input.CustomerUserInfo.IsAdmin,
				},
				"credentials": map[string]interface{}{
					"api_token":         input.Credential.ApiAddress,
					"ids_client_secret": input.Credential.IdsClientSecret,
					"ids_client_id":     input.Credential.IdsClientID,
					"ids_scope":         input.Credential.IdsScope,
					"api_address":       input.Credential.ApiAddress,
					"ids_address":       input.Credential.IdsAddress,
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
