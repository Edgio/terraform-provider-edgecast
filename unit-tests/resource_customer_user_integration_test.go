// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test cases for storage account name conversion logic
var tcCustomerUser = map[string]ResourceNewCustomerUser{
	"terratest.testing.vmp.customer": {
		customerUserInfo: NewCustomerUserInfo{
			accountnumber: "D9127",
			firstname:     "John",
			lastname:      "Doe",
			email:         "admin+1@test20252021-7.com",
			isadmin:       false,
		},
		credential: Credentials{
			apitoken:        "AULdReDoB3gb0D7LNTx857NQvrcIKyvL",
			idsclientsecret: "CDbbMJw7FFJ11a7433ti1l9XgJHKr2Wk",
			idsclientID:     "31ef8e8f-0120-4112-8554-3eb11e83d58b",
			idsscope:        "ec.rules",
			apiaddress:      "http://dev-api.edgecast.com",
			idsaddress:      "https://id-dev.vdms.io",
		},
	},
}

func TestUT_CustomerUser_basic(t *testing.T) {
	t.Parallel()

	for _, input := range tcCustomerUser {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../examples/resource_customer_user",
			Vars: map[string]interface{}{
				"new_admin_user": map[string]interface{}{
					"customer_account_number": input.customerUserInfo.accountnumber,
					"first_name":              input.customerUserInfo.firstname,
					"last_name":               input.customerUserInfo.lastname,
					"email":                   input.customerUserInfo.email,
					"is_admin":                input.customerUserInfo.isadmin,
				},
				"credentials": map[string]interface{}{
					"api_token":         input.credential.apitoken,
					"ids_client_secret": input.credential.idsclientsecret,
					"ids_client_id":     input.credential.idsclientID,
					"ids_scope":         input.credential.idsscope,
					"api_address":       input.credential.apiaddress,
					"ids_address":       input.credential.idsaddress,
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
