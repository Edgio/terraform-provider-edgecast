package data

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/customer"
	"terraform-provider-edgecast/test/integration/cmd/create-import-data/internal"
)

func createCustomerData(cfg edgecast.SDKConfig) (accountNumber string, customerUser int) {
	accountNumber = account()
	customerService := internal.Check(customer.New(cfg)).(*customer.CustomerService)
	customerUser = internal.Check(createCustomerUser(customerService, accountNumber)).(int)
	return
}

func createCustomerUser(service *customer.CustomerService, accountNumber string) (int, error) {
	customerGetOK := internal.Check(service.GetCustomer(customer.GetCustomerParams{
		AccountNumber: accountNumber,
	})).(*customer.CustomerGetOK)

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
			Email:     email(),
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
			AccountID:                 account(),
			Address1:                  "111 main street",
			Address2:                  "",
			City:                      "Beverly Hills",
			State:                     "CA",
			ZIP:                       "90210",
			Country:                   "USA",
			BillingCountry:            "USA",
			BillingState:              "CA",
			BillingZIP:                "90210",
			ContactFirstName:          "Dev",
			ContactLastName:           "Integration",
			ContactMobile:             "555-5555",
			ContactPhone:              "555-5555",
			ContactTitle:              "Chief Dev",
			CompanyName:               "Dev Integration",
			DataTransferredUsageLimit: 10000,
			Notes:                     "",
			PartnerUserID:             1,
			ServiceLevelCode:          "",
			Website:                   "",
			Status:                    1,
		},
	})
}
