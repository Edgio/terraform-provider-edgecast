// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package data

import (
	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/customer"
)

func createCustomerData(cfg config.Config) CustomerResult {
	// Need to use a PCC token to create a customer.
	pccSDKConfig := cfg.SDKConfig
	pccSDKConfig.APIToken = cfg.APITokenPCC

	svc := internal.Check(customer.New(pccSDKConfig))
	accountNumber := cfg.AccountNumber

	if accountNumber == "" {
		accountNumber = internal.Check(createCustomer(svc))
	}

	customerUserID := internal.Check(
		createCustomerUser(svc, accountNumber, cfg.TestEmail),
	)

	return CustomerResult{accountNumber, customerUserID}
}

func createCustomerUser(
	service *customer.CustomerService,
	accountNumber string,
	testEmail string,
) (int, error) {
	customerGetOK := internal.Check(
		service.GetCustomer(customer.GetCustomerParams{
			AccountNumber: accountNumber,
		}))

	return service.AddCustomerUser(customer.AddCustomerUserParams{
		Customer: *customerGetOK,
		CustomerUser: customer.CustomerUser{
			Address1:  "100 Main St",
			Address2:  "",
			City:      "Beverly Hills",
			State:     "CA",
			ZIP:       "90210",
			Country:   "US",
			Mobile:    "1-555-5555",
			Phone:     "1-555-5555",
			Fax:       "",
			Email:     testEmail,
			Title:     "Tester",
			FirstName: "Dev",
			LastName:  "Integration",
			IsAdmin:   0,
		},
	})
}

func createCustomer(service *customer.CustomerService) (string, error) {
	return service.AddCustomer(customer.AddCustomerParams{
		Customer: customer.Customer{
			Address1:                  "111 main street",
			Address2:                  "",
			City:                      "Beverly Hills",
			State:                     "CA",
			ZIP:                       "90210",
			Country:                   "USA",
			BandwidthUsageLimit:       0,
			BillingAccountTag:         "",
			BillingAddress1:           "",
			BillingAddress2:           "",
			BillingCity:               "",
			BillingContactEmail:       "",
			BillingContactFax:         "",
			BillingContactFirstName:   "",
			BillingContactLastName:    "",
			BillingContactMobile:      "",
			BillingContactPhone:       "",
			BillingContactTitle:       "",
			BillingCountry:            "USA",
			BillingRateInfo:           "",
			BillingState:              "CA",
			BillingZIP:                "90210",
			ContactEmail:              "",
			ContactFax:                "",
			ContactFirstName:          "Dev",
			ContactLastName:           "Integration",
			ContactMobile:             "555-5555",
			ContactPhone:              "555-5555",
			ContactTitle:              "Chief Dev",
			CompanyName:               "Dev Integration",
			DataTransferredUsageLimit: 10000,
			Notes:                     "",
			PartnerUserID:             1,
			ServiceLevelCode:          "STND",
			Website:                   "",
			Status:                    1,
		},
	})
}
