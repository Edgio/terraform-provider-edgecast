// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package test

import (
	"terraform-provider-edgecast/test/unit/model"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test cases for storage account name conversion logic
var tcCustomer = map[string]model.ResourceNewCustomer{
	"terratest.testing.ec.customer": {
		CustomerInfo: model.NewCustomerInfo{
			CompanyName:      "Terraform test customer02192021-4",
			ServiceLevelCode: "STND",
			Services:         []int{1, 9, 15, 19},
			DeliveryRegion:   1,
			AccessModules:    []int{1, 4, 5, 7, 8, 21, 22, 25, 26, 27, 29, 30, 32, 40, 46, 53, 56, 71, 72, 73, 74, 75, 76, 77, 78, 79, 81, 138, 139, 140, 144, 145, 146, 149, 153, 157, 159, 160, 161, 162, 163, 164, 166, 168, 169, 170, 171, 172, 174, 175, 176, 177, 178, 179, 180, 182, 183, 184, 185, 186, 187, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 204, 386, 387, 409, 410, 411, 412, 413, 414, 415, 416, 479, 501, 502},
		},
		Credential: model.Credentials{
			ApiToken:         "<apitoken>",
			IdsClientSecret:  "<idsclientsecret>",
			IdsClientID:      "<idssclientID>",
			IdsScope:         "<scope>",
			ApiAddress:       "<apiUrl>",
			IdsAddress:       "<idsaddress>",
			ApiAddressLegacy: "<apiAddressLegacy",
		},
	},
}

func TestUT_Customer_basic(t *testing.T) {
	t.Parallel()
	t.Skip("test is not ready for unit testing")

	for _, input := range tcCustomer {
		// Specify the test case folder and "-var" options
		tfOptions := &terraform.Options{
			TerraformDir: "../examples/resources/ec_customer",
			Vars: map[string]interface{}{
				"new_customer_info": map[string]interface{}{
					"company_name":       input.CustomerInfo.CompanyName,
					"service_level_code": input.CustomerInfo.ServiceLevelCode,
					"services":           input.CustomerInfo.Services,
					"delivery_region":    input.CustomerInfo.DeliveryRegion,
					"access_modules":     input.CustomerInfo.AccessModules,
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
