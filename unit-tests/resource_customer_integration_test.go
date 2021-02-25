// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test cases for storage account name conversion logic
var tcCustomer = map[string]ResourceNewCustomer{
	"terratest.testing.vmp.customer": {
		customerInfo: NewCustomerInfo{
			companyname:      "Terraform test customer02192021-4",
			servicelevelcode: "STND",
			services:         []int{1, 9, 15, 19},
			deliveryregion:   1,
			accessmodules:    []int{1, 4, 5, 7, 8, 21, 22, 25, 26, 27, 29, 30, 32, 40, 46, 53, 56, 71, 72, 73, 74, 75, 76, 77, 78, 79, 81, 138, 139, 140, 144, 145, 146, 149, 153, 157, 159, 160, 161, 162, 163, 164, 166, 168, 169, 170, 171, 172, 174, 175, 176, 177, 178, 179, 180, 182, 183, 184, 185, 186, 187, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 204, 386, 387, 409, 410, 411, 412, 413, 414, 415, 416, 479, 501, 502},
		},
		credential: Credentials{
			apitoken:        "<apitoken>",
			idsclientsecret: "<idsclientsecret>",
			idsclientID:     "<idssclientID>",
			idsscope:        "<scope>",
			apiaddress:      "<apiUrl>",
			idsaddress:      "<idsaddress>",
		},
	},
}

func TestUT_Customer_basic(t *testing.T) {
	t.Parallel()

	for _, input := range tcCustomer {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../examples/resource_customer",
			Vars: map[string]interface{}{
				"new_customer_info": map[string]interface{}{
					"company_name":       input.customerInfo.companyname,
					"service_level_code": input.customerInfo.servicelevelcode,
					"services":           input.customerInfo.services,
					"delivery_region":    input.customerInfo.deliveryregion,
					"access_modules":     input.customerInfo.accessmodules,
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
