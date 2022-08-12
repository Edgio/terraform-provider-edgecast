// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"reflect"
	"terraform-provider-edgecast/edgecast/helper"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/go-test/deep"
)

func TestExpandOrganization(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   *models.OrganizationDetail
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: helper.NewTerraformSet([]interface{}{
				map[string]interface{}{
					"city":                "L.A.",
					"company_address":     "111 fantastic way",
					"company_address2":    "111 fantastic way",
					"company_name":        "Test Co.",
					"contact_email":       "user3@test.com",
					"contact_first_name":  "test3",
					"contact_last_name":   "user",
					"contact_phone":       "111-111-1111",
					"contact_title":       "N/A",
					"country":             "US",
					"organizational_unit": "Dept1",
					"state":               "CA",
					"zip_code":            "90001",
					"additional_contact": []interface{}{
						map[string]interface{}{
							"first_name":   "contact1",
							"last_name":    "lastname",
							"email":        "first.lastname@testuser.com",
							"phone":        "111-111-2222",
							"title":        "Manager",
							"contact_type": "EvApprover",
						},
					},
				},
			}),
			expectedPtr: &models.OrganizationDetail{
				City:               "L.A.",
				CompanyAddress:     "111 fantastic way",
				CompanyAddress2:    "111 fantastic way",
				CompanyName:        "Test Co.",
				ContactEmail:       "user3@test.com",
				ContactFirstName:   "test3",
				ContactLastName:    "user",
				ContactPhone:       "111-111-1111",
				ContactTitle:       "N/A",
				Country:            "US",
				OrganizationalUnit: "Dept1",
				State:              "CA",
				ZipCode:            "90001",
				AdditionalContacts: []*models.OrganizationContact{
					{
						FirstName:   "contact1",
						LastName:    "lastname",
						Email:       "first.lastname@testuser.com",
						Phone:       "111-111-2222",
						Title:       "Manager",
						ContactType: "EvApprover",
					},
				},
			},

			expectSuccess: true,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a []interface{}",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {

		actualPtr, err := expandOrganization(v.input)

		if v.expectSuccess {
			if err == nil {

				actual := *actualPtr
				expected := *v.expectedPtr

				if !reflect.DeepEqual(actual, expected) {
					// deep.Equal doesn't compare pointer values, so we just use it to
					// generate a human friendly diff
					diff := deep.Equal(actual, expected)
					t.Errorf("Diff: %+v", diff)
					t.Fatalf("%s: Expected %+v but got %+v",
						v.name,
						expected,
						actual,
					)
				}

			} else {
				t.Fatalf("%s: Encountered error where one was not expected: %+v",
					v.name,
					err,
				)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}
