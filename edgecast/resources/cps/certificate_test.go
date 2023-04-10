// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"errors"
	"log"
	"reflect"
	"sort"
	"testing"
	"time"

	"terraform-provider-edgecast/edgecast/helper"

	sdkcps "github.com/EdgeCast/ec-sdk-go/edgecast/cps"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandCertificate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       map[string]any
		expectedPtr *CertificateState
		expectErrs  bool
		errCount    int
	}{
		{
			name:       "Happy path",
			expectErrs: false,
			input: map[string]any{
				"certificate_label":     "my_cert",
				"auto_renew":            true,
				"certificate_authority": "authorityA",
				"dcv_method":            "methodA",
				"description":           "this is a description",
				"validation_type":       "typeA",
				"domain": []any{
					map[string]any{
						"is_common_name": true,
						"name":           "testdomain1.com",
					},
					map[string]any{
						"is_common_name": false,
						"name":           "testdomain2.com",
					},
				},
				"notification_setting": []any{
					map[string]any{
						"enabled":           true,
						"notification_type": "CertificateRenewal",
						"emails": []any{
							"email1@test.com",
							"email2@test.com",
						},
					},
					map[string]any{
						"enabled":           false,
						"notification_type": "CertificateExpiring",
						"emails":            make([]any, 0),
					},
					map[string]any{
						"enabled":           false,
						"notification_type": "PendingValidations",
						"emails":            make([]any, 0),
					},
				},
				"organization": []any{
					map[string]any{
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
						"additional_contact": []any{
							map[string]any{
								"first_name":   "contact1",
								"last_name":    "lastname",
								"email":        "first.lastname@testuser.com",
								"phone":        "111-111-2222",
								"title":        "Manager",
								"contact_type": "EvApprover",
							},
							map[string]any{
								"first_name":   "contact2",
								"last_name":    "lastname2",
								"email":        "first.lastname2@testuser.com",
								"phone":        "111-111-3333",
								"title":        "Developer",
								"contact_type": "TechnicalContact",
							},
						},
					},
				},
			},
			expectedPtr: &CertificateState{
				CertificateLabel:     "my_cert",
				CertificateAuthority: "authorityA",
				DcvMethod:            "methodA",
				AutoRenew:            true,
				Description:          "this is a description",
				ValidationType:       "typeA",
				NotificationSettings: []*models.EmailNotification{
					{
						Enabled:          true,
						NotificationType: "CertificateRenewal",
						Emails:           []string{"email1@test.com", "email2@test.com"},
					},
					{
						Enabled:          false,
						NotificationType: "CertificateExpiring",
						Emails:           make([]string, 0),
					},
					{
						Enabled:          false,
						NotificationType: "PendingValidations",
						Emails:           make([]string, 0),
					},
				},
				Domains: []*models.DomainCreateUpdate{
					{
						IsCommonName: true,
						Name:         "testdomain1.com",
					},
					{
						IsCommonName: false,
						Name:         "testdomain2.com",
					},
				},
				OrganizationChanged: true, // HasChange == true during create
				Organization: &models.OrganizationDetail{
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
						{
							FirstName:   "contact2",
							LastName:    "lastname2",
							Email:       "first.lastname2@testuser.com",
							Phone:       "111-111-3333",
							Title:       "Developer",
							ContactType: "TechnicalContact",
						},
					},
				},
			},
		},
		{
			name:        "nil input",
			errCount:    1,
			input:       nil,
			expectedPtr: nil,
			expectErrs:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var rd *schema.ResourceData
			if tt.input != nil {
				rd = schema.TestResourceDataRaw(
					t,
					GetCertificateSchema(),
					tt.input)
			}

			actualPtr, errs := ExpandCertificate(rd)

			if !tt.expectErrs && (len(errs) > 0) {
				t.Fatalf("unexpected errors: %v", errs)
			}

			if tt.expectErrs && (len(errs) != tt.errCount) {
				t.Fatalf("expected %d errors but got %d", tt.errCount, len(errs))
			}

			if tt.expectErrs && (len(errs) > 0) {
				return // successful test for error case
			}

			actual := *actualPtr
			expected := *tt.expectedPtr

			// TF sets do not guarantee order, so we must sort before comparing.
			// We will sort on NotificationType.
			got := actual.NotificationSettings
			want := expected.NotificationSettings
			sort.SliceStable(got, func(i, j int) bool {
				return got[i].NotificationType < got[j].NotificationType
			})

			sort.SliceStable(want, func(i, j int) bool {
				return want[i].NotificationType < want[j].NotificationType
			})

			diffs := deep.Equal(got, want)
			if len(diffs) > 0 {
				t.Logf("got notifsettings: %v, want %v", got, want)
				t.Errorf("Differences: %v", diffs)
			}

			// Nil out Notif Settings since we already checked them
			actual.NotificationSettings = nil
			expected.NotificationSettings = nil

			if !reflect.DeepEqual(actual, expected) {
				// deep.Equal doesn't compare pointer values, so we just use it to
				// generate a human friendly diff
				diff := deep.Equal(actual, expected)
				t.Errorf("Diff: %+v", diff)
				t.Fatalf("%s: Expected %+v but got %+v",
					tt.name,
					expected,
					actual,
				)
			}
		})
	}
}

func TestExpandOrganization(t *testing.T) {
	cases := []struct {
		name          string
		input         any
		expectedPtr   *models.OrganizationDetail
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []any{
				map[string]any{
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
					"additional_contact": []any{
						map[string]any{
							"first_name":   "contact1",
							"last_name":    "lastname",
							"email":        "first.lastname@testuser.com",
							"phone":        "111-111-2222",
							"title":        "Manager",
							"contact_type": "EvApprover",
						},
						map[string]any{
							"first_name":   "contact2",
							"last_name":    "lastname2",
							"email":        "first.lastname2@testuser.com",
							"phone":        "111-111-3333",
							"title":        "Developer",
							"contact_type": "TechnicalContact",
						},
					},
				},
			},
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
					{
						FirstName:   "contact2",
						LastName:    "lastname2",
						Email:       "first.lastname2@testuser.com",
						Phone:       "111-111-3333",
						Title:       "Developer",
						ContactType: "TechnicalContact",
					},
				},
			},
			expectSuccess: true,
		},
		{
			name:          "nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a []interface{}",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ExpandOrganization(v.input)

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

func TestExpandOrganizationContact(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   []*models.OrganizationContact
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []interface{}{
				map[string]interface{}{
					"first_name":   "contact1",
					"last_name":    "lastname",
					"email":        "first.lastname@testuser.com",
					"phone":        "111-111-2222",
					"title":        "Manager",
					"contact_type": "EvApprover",
				},
				map[string]interface{}{
					"first_name":   "contact2",
					"last_name":    "lastname2",
					"email":        "first.lastname2@testuser.com",
					"phone":        "111-111-3333",
					"title":        "Developer",
					"contact_type": "TechnicalContact",
				},
			},
			expectedPtr: []*models.OrganizationContact{
				{
					FirstName:   "contact1",
					LastName:    "lastname",
					Email:       "first.lastname@testuser.com",
					Phone:       "111-111-2222",
					Title:       "Manager",
					ContactType: "EvApprover",
				},
				{
					FirstName:   "contact2",
					LastName:    "lastname2",
					Email:       "first.lastname2@testuser.com",
					Phone:       "111-111-3333",
					Title:       "Developer",
					ContactType: "TechnicalContact",
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
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ExpandAdditionalContacts(v.input)

		if v.expectSuccess {
			if err == nil {
				actual := actualPtr
				expected := v.expectedPtr

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

func TestExpandDomains(t *testing.T) {
	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   []*models.DomainCreateUpdate
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []interface{}{
				map[string]interface{}{
					"is_common_name": true,
					"name":           "testdomain1.com",
				},
				map[string]interface{}{
					"is_common_name": false,
					"name":           "testdomain2.com",
				},
			},
			expectedPtr: []*models.DomainCreateUpdate{
				{
					IsCommonName: true,
					Name:         "testdomain1.com",
				},
				{
					IsCommonName: false,
					Name:         "testdomain2.com",
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
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := ExpandDomains(v.input)

		if v.expectSuccess {
			if err == nil {
				actual := actualPtr
				expected := v.expectedPtr

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

func TestFlattenDeployments(t *testing.T) {
	cases := []struct {
		name     string
		input    []*models.RequestDeployment
		expected []map[string]interface{}
	}{
		{
			name: "Happy path",
			input: []*models.RequestDeployment{
				{
					DeliveryRegion: "Delivery Region 1",
					HexURL:         "hex 1",
					Platform:       "Platform 1",
				},
				{
					DeliveryRegion: "Delivery Region 2",
					HexURL:         "hex 2",
					Platform:       "Platform 2",
				},
			},
			expected: []map[string]interface{}{
				{
					"delivery_region": "Delivery Region 1",
					"hex_url":         "hex 1",
					"platform":        "Platform 1",
				},
				{
					"delivery_region": "Delivery Region 2",
					"hex_url":         "hex 2",
					"platform":        "Platform 2",
				},
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: make([]map[string]interface{}, 0),
		},
		{
			name:     "Empty input",
			input:    make([]*models.RequestDeployment, 0),
			expected: make([]map[string]interface{}, 0),
		},
	}

	for _, c := range cases {
		actual := FlattenDeployments(c.input)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}

func TestFlattenActor(t *testing.T) {
	cases := []struct {
		name          string
		input         *models.Actor
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{
			name:          "Happy path",
			expectSuccess: true,
			input: &models.Actor{
				UserID:       0,
				PortalTypeID: "customer",
				IdentityID:   "abc-xyz",
				IdentityType: "actor",
			},
			expected: []map[string]interface{}{
				{
					"user_id":        0,
					"portal_type_id": "customer",
					"identity_id":    "abc-xyz",
					"identity_type":  "actor",
				},
			},
		},
		{
			name:          "Nil input",
			input:         nil,
			expected:      make([]map[string]interface{}, 0),
			expectSuccess: false,
		},
		{
			name:  "Empty imput",
			input: &models.Actor{},
			expected: []map[string]interface{}{
				{
					"user_id":        0,
					"portal_type_id": "",
					"identity_id":    "",
					"identity_type":  "",
				},
			},
			expectSuccess: true,
		},
	}

	for _, c := range cases {
		actual := FlattenActor(c.input)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}

func TestExpandNotifSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
		args        interface{}
		want        []*models.EmailNotification
	}{
		{
			name:        "Happy Path",
			expectError: false,
			args: helper.NewTerraformSet([]any{
				map[string]any{
					"enabled":           true,
					"notification_type": "CertificateRenewal",
					"emails": []string{
						"email1@test.com",
						"email2@test.com",
					},
				},
				map[string]any{
					"enabled":           false,
					"notification_type": "CertificateExpiring",
					"emails":            make([]string, 0),
				},
				map[string]any{
					"enabled":           false,
					"notification_type": "PendingValidations",
					"emails":            make([]string, 0),
				},
			}),
			want: []*models.EmailNotification{
				{
					Enabled:          true,
					NotificationType: "CertificateRenewal",
					Emails:           []string{"email1@test.com", "email2@test.com"},
				},
				{
					Enabled:          false,
					NotificationType: "CertificateExpiring",
					Emails:           make([]string, 0),
				},
				{
					Enabled:          false,
					NotificationType: "PendingValidations",
					Emails:           make([]string, 0),
				},
			},
		},
		{
			name:        "Empty input results in empty non-nil result",
			expectError: false,
			args:        helper.NewTerraformSet(make([]any, 0)),
			want:        make([]*models.EmailNotification, 0),
		},
		{
			name:        "Nil input results in empty non-nil result",
			expectError: false,
			args:        nil,
			want:        make([]*models.EmailNotification, 0),
		},
		{
			name:        "Error - input is unexpected type",
			expectError: true,
			args:        1,
		},
		{
			name:        "Error - set contains non-map item",
			expectError: true,
			args:        helper.NewTerraformSet([]any{1}),
		},
		{
			name:        "Error - missing attributes",
			expectError: true,
			args: helper.NewTerraformSet([]any{
				map[string]any{
					"enabled": false,
				},
				map[string]any{
					"notification_type": "CertificateExpiring",
				},
				map[string]any{
					"emails": make([]string, 0),
				},
			}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, errs := ExpandNotifSettings(tt.args)

			if tt.expectError && len(errs) > 0 {
				return // successful test
			}

			if tt.expectError && len(errs) == 0 {
				t.Fatal("expected error, but got none")
			}

			if !tt.expectError && len(errs) > 0 {
				t.Fatalf("unexpected errors: %v", errs)
			}

			// TF sets do not guarantee order, so we must sort before comparing.
			// We will sort on NotificationType.
			sort.SliceStable(got, func(i, j int) bool {
				return got[i].NotificationType < got[j].NotificationType
			})

			sort.SliceStable(tt.want, func(i, j int) bool {
				return tt.want[i].NotificationType < tt.want[j].NotificationType
			})

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}

func TestFlattenNotifSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []*models.EmailNotification
		want []map[string]any
	}{
		{
			name: "Happy Path",
			args: []*models.EmailNotification{
				{
					Enabled:          true,
					NotificationType: "CertificateRenewal",
					Emails:           []string{"email1@test.com", "email2@test.com"},
				},
				{
					Enabled:          false,
					NotificationType: "CertificateExpiring",
					Emails:           nil,
				},
				{
					Enabled:          false,
					NotificationType: "PendingValidations",
					Emails:           make([]string, 0),
				},
			},
			want: []map[string]any{
				{
					"enabled":           true,
					"notification_type": "CertificateRenewal",
					"emails": []string{
						"email1@test.com",
						"email2@test.com",
					},
				},
				{
					"enabled":           false,
					"notification_type": "CertificateExpiring",
					"emails":            make([]string, 0),
				},
				{
					"enabled":           false,
					"notification_type": "PendingValidations",
					"emails":            make([]string, 0),
				},
			},
		},
		{
			name: "Empty input results in empty non-nil result",
			args: make([]*models.EmailNotification, 0),
			want: make([]map[string]any, 0),
		},
		{
			name: "Nil input results in empty non-nil result",
			args: nil,
			want: make([]map[string]any, 0),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FlattenNotifSettings(tt.args)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}

func TestFlattenDomainsDV(t *testing.T) {
	activeDate := strfmt.DateTime(time.Now())
	createdDate := strfmt.DateTime(time.Now())

	cases := []struct {
		name          string
		inputDomains  []*models.Domain
		inputMetadata []*models.DomainDcvFull
		input         string
		expected      []map[string]interface{}
	}{
		{
			name: "Happy path",
			inputDomains: []*models.Domain{
				{
					ID:           1,
					Name:         "domain 1",
					Status:       "Active",
					IsCommonName: true,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
				{
					ID:           2,
					Name:         "domain 2",
					Status:       "Active",
					IsCommonName: false,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
			},
			inputMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DV",
					DcvToken:  &models.DcvToken{Token: "token"},
					DomainID:  0,
					Emails:    []string{"email1@test.com", "email2@test.com"},
				},
			},
			input: "DV",
			expected: []map[string]interface{}{
				{
					"id":             int64(1),
					"name":           "domain 1",
					"status":         "Active",
					"is_common_name": true,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
					"emails": []string{
						"email1@test.com",
						"email2@test.com",
					},
					"dcv_token": "token",
				},
				{
					"id":             int64(2),
					"name":           "domain 2",
					"status":         "Active",
					"is_common_name": false,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
					"emails": []string{
						"email1@test.com",
						"email2@test.com",
					},
					"dcv_token": "token",
				},
			},
		},
		{
			name: "Happy path - no metadata",
			inputDomains: []*models.Domain{
				{
					ID:           1,
					Name:         "domain 1",
					Status:       "Active",
					IsCommonName: true,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
				{
					ID:           2,
					Name:         "domain 2",
					Status:       "Active",
					IsCommonName: false,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
			},
			inputMetadata: nil,
			input:         "DV",
			expected: []map[string]interface{}{
				{
					"id":             int64(1),
					"name":           "domain 1",
					"status":         "Active",
					"is_common_name": true,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
				},
				{
					"id":             int64(2),
					"name":           "domain 2",
					"status":         "Active",
					"is_common_name": false,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
				},
			},
		},
		{
			name:          "Nil input",
			inputDomains:  nil,
			inputMetadata: nil,
			input:         "",
			expected:      make([]map[string]interface{}, 0),
		},
		{
			name:          "Empty input",
			inputDomains:  make([]*models.Domain, 0),
			inputMetadata: make([]*models.DomainDcvFull, 0),
			input:         "",
			expected:      make([]map[string]interface{}, 0),
		},
	}

	for _, c := range cases {
		actual, _ := FlattenDomains(c.inputDomains, c.inputMetadata, c.input)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}

func TestFlattenDomainsOV(t *testing.T) {
	activeDate := strfmt.DateTime(time.Now())
	createdDate := strfmt.DateTime(time.Now())

	cases := []struct {
		name          string
		inputDomains  []*models.Domain
		inputMetadata []*models.DomainDcvFull
		input         string
		expected      []map[string]interface{}
	}{
		{
			name: "Happy path",
			inputDomains: []*models.Domain{
				{
					ID:           1,
					Name:         "domain 1",
					Status:       "Active",
					IsCommonName: true,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
				{
					ID:           2,
					Name:         "domain 2",
					Status:       "Active",
					IsCommonName: false,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
			},
			inputMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "OV",
					DcvToken:  &models.DcvToken{Token: "token 1"},
					DomainID:  1,
					Emails:    []string{"email1@test.com", "email2@test.com"},
				},
				{
					DcvMethod: "OV",
					DcvToken:  &models.DcvToken{Token: "token 2"},
					DomainID:  2,
					Emails:    []string{"email3@test.com", "email4@test.com"},
				},
			},
			input: "OV",
			expected: []map[string]interface{}{
				{
					"id":             int64(1),
					"name":           "domain 1",
					"status":         "Active",
					"is_common_name": true,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
					"emails": []string{
						"email1@test.com",
						"email2@test.com",
					},
					"dcv_token": "token 1",
				},
				{
					"id":             int64(2),
					"name":           "domain 2",
					"status":         "Active",
					"is_common_name": false,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
					"emails": []string{
						"email3@test.com",
						"email4@test.com",
					},
					"dcv_token": "token 2",
				},
			},
		},
		{
			name: "missing id in metadata",
			inputDomains: []*models.Domain{
				{
					ID:           1,
					Name:         "domain 1",
					Status:       "Active",
					IsCommonName: true,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
				{
					ID:           2,
					Name:         "domain 2",
					Status:       "Active",
					IsCommonName: false,
					ActiveDate:   activeDate,
					Created:      createdDate,
				},
			},
			inputMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "OV",
					DcvToken:  &models.DcvToken{Token: "token 2"},
					DomainID:  2,
					Emails:    []string{"email1@test.com", "email2@test.com"},
				},
				{
					DcvMethod: "OV",
					DcvToken:  &models.DcvToken{Token: "token 4"},
					DomainID:  4,
					Emails:    []string{"email3@test.com", "email4@test.com"},
				},
			},
			input: "OV",
			expected: []map[string]interface{}{
				{
					"id":             int64(1),
					"name":           "domain 1",
					"status":         "Active",
					"is_common_name": true,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
				},
				{
					"id":             int64(2),
					"name":           "domain 2",
					"status":         "Active",
					"is_common_name": false,
					"active_date":    activeDate.String(),
					"created":        createdDate.String(),
					"emails": []string{
						"email1@test.com",
						"email2@test.com",
					},
					"dcv_token": "token 2",
				},
			},
		},
		{
			name:          "Nil input",
			inputDomains:  nil,
			inputMetadata: nil,
			input:         "",
			expected:      make([]map[string]interface{}, 0),
		},
		{
			name:          "Empty input",
			inputDomains:  make([]*models.Domain, 0),
			inputMetadata: make([]*models.DomainDcvFull, 0),
			input:         "",
			expected:      make([]map[string]interface{}, 0),
		},
	}

	for _, c := range cases {
		actual, _ := FlattenDomains(c.inputDomains, c.inputMetadata, c.input)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}

func TestFlattenOrganization(t *testing.T) {
	cases := []struct {
		name          string
		input         *models.OrganizationDetail
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{
			name:          "Happy path",
			expectSuccess: true,
			input: &models.OrganizationDetail{
				ID:                 0,
				Country:            "US",
				City:               "L.A.",
				CompanyAddress:     "111 fantastic way",
				CompanyAddress2:    "111 fantastic way",
				CompanyName:        "Test Co.",
				ContactEmail:       "user3@test.com",
				ContactFirstName:   "test3",
				ContactLastName:    "user",
				ContactPhone:       "111-111-1111",
				ContactTitle:       "N/A",
				OrganizationalUnit: "Dept1",
				State:              "CA",
				ZipCode:            "90001",
				AdditionalContacts: []*models.OrganizationContact{
					{
						ID:          0,
						FirstName:   "contact1",
						LastName:    "lastname",
						Email:       "first.lastname@testuser.com",
						Phone:       "111-111-2222",
						Title:       "Manager",
						ContactType: "EvApprover",
					},
					{
						ID:          1,
						FirstName:   "contact2",
						LastName:    "lastname2",
						Email:       "first.lastname2@testuser.com",
						Phone:       "111-111-3333",
						Title:       "Developer",
						ContactType: "TechnicalContact",
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"id":                  int64(0),
					"country":             "US",
					"city":                "L.A.",
					"company_address":     "111 fantastic way",
					"company_address2":    "111 fantastic way",
					"company_name":        "Test Co.",
					"contact_email":       "user3@test.com",
					"contact_first_name":  "test3",
					"contact_last_name":   "user",
					"contact_phone":       "111-111-1111",
					"contact_title":       "N/A",
					"organizational_unit": "Dept1",
					"state":               "CA",
					"zip_code":            "90001",
					"additional_contact": []map[string]interface{}{
						{
							"id":           int64(0),
							"first_name":   "contact1",
							"last_name":    "lastname",
							"email":        "first.lastname@testuser.com",
							"phone":        "111-111-2222",
							"title":        "Manager",
							"contact_type": "EvApprover",
						},
						{
							"id":           int64(1),
							"first_name":   "contact2",
							"last_name":    "lastname2",
							"email":        "first.lastname2@testuser.com",
							"phone":        "111-111-3333",
							"title":        "Developer",
							"contact_type": "TechnicalContact",
						},
					},
				},
			},
		},
		{
			name:          "Nil input",
			input:         nil,
			expected:      make([]map[string]interface{}, 0),
			expectSuccess: false,
		},
	}

	for _, c := range cases {
		actual := FlattenOrganization(c.input)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}

func TestFlattenRequestStatus(t *testing.T) {
	cases := []struct {
		name          string
		input         *models.CertificateStatus
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{
			name:          "Happy path",
			expectSuccess: true,
			input: &models.CertificateStatus{
				Step:              0,
				Status:            "Step 0",
				RequiresAttention: false,
				ErrorMessage:      "Error Message",
				OrderValidation: &models.OrderValidation{
					Status: "Active",
					OrganizationValidation: &models.OrganizationValidation{
						ValidationType: "Validation Type 1",
						Status:         "Active",
					},
					DomainValidations: []*models.DomainValidation{
						{
							Status: "Active",
							DomainNames: []string{
								"domain1",
								"domain2",
							},
						},
						{
							Status: "Active",
							DomainNames: []string{
								"domain3",
								"domain4",
							},
						},
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"step":               int32(0),
					"status":             "Step 0",
					"requires_attention": false,
					"error_message":      "Error Message",
					"order_validation": []map[string]interface{}{
						{
							"status": "Active",

							"organization_validation": []map[string]interface{}{
								{
									"validation_type": "Validation Type 1",
									"status":          "Active",
								},
							},
							"domain_validation": []map[string]interface{}{
								{
									"status":       "Active",
									"domain_names": []string{"domain1", "domain2"},
								},
								{
									"status":       "Active",
									"domain_names": []string{"domain3", "domain4"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:          "Nil input",
			input:         nil,
			expected:      make([]map[string]interface{}, 0),
			expectSuccess: false,
		},
	}

	for _, c := range cases {
		actual := FlattenRequestStatus(c.input)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}

func TestGroupDomains(t *testing.T) {
	cases := []struct {
		name         string
		inputDomains []*models.Domain
		expected     map[int64][]*models.Domain
	}{
		{
			name: "parent child relationships",
			inputDomains: []*models.Domain{
				{
					ID:           1,
					Name:         "testparentdomain.com",
					Status:       "Active",
					IsCommonName: true,
				},
				{
					ID:           2,
					Name:         "subdomain.testparentdomain.com",
					Status:       "Active",
					IsCommonName: false,
				},
				{
					ID:           3,
					Name:         "subdomain2.testparentdomain.com",
					Status:       "Active",
					IsCommonName: false,
				},
				{
					ID:           4,
					Name:         "test.anotherparentdomain.com",
					Status:       "Active",
					IsCommonName: false,
				},
				{
					ID:           5,
					Name:         "Subdomain.test.anotherparentdomain.com",
					Status:       "Active",
					IsCommonName: false,
				},
				{
					ID:           6,
					Name:         "parentdomain.com",
					Status:       "Active",
					IsCommonName: false,
				},
				{
					ID:           7,
					Name:         "someparentdomain.com",
					Status:       "Active",
					IsCommonName: false,
				},
			},
			expected: map[int64][]*models.Domain{
				int64(1): {
					{
						ID:           1,
						Name:         "testparentdomain.com",
						Status:       "Active",
						IsCommonName: true,
					},
					{
						ID:           2,
						Name:         "subdomain.testparentdomain.com",
						Status:       "Active",
						IsCommonName: false,
					},
					{
						ID:           3,
						Name:         "subdomain2.testparentdomain.com",
						Status:       "Active",
						IsCommonName: false,
					},
				},
				int64(4): {
					{
						ID:           4,
						Name:         "test.anotherparentdomain.com",
						Status:       "Active",
						IsCommonName: false,
					},
					{
						ID:           5,
						Name:         "Subdomain.test.anotherparentdomain.com",
						Status:       "Active",
						IsCommonName: false,
					},
				},
				int64(6): {
					{
						ID:           6,
						Name:         "parentdomain.com",
						Status:       "Active",
						IsCommonName: false,
					},
				},
				int64(7): {
					{
						ID:           7,
						Name:         "someparentdomain.com",
						Status:       "Active",
						IsCommonName: false,
					},
				},
			},
		},
		{
			name: "no parent child relationships",
			inputDomains: []*models.Domain{
				{
					ID:           1,
					Name:         "testdomain.com",
					Status:       "Active",
					IsCommonName: true,
				},
				{
					ID:           2,
					Name:         "subdomain.test.com",
					Status:       "Active",
					IsCommonName: false,
				},
				{
					ID:           3,
					Name:         "subdomain2.domain.com",
					Status:       "Active",
					IsCommonName: false,
				},
			},
			expected: map[int64][]*models.Domain{
				int64(1): {
					{
						ID:           1,
						Name:         "testdomain.com",
						Status:       "Active",
						IsCommonName: true,
					},
				},
				int64(2): {
					{
						ID:           2,
						Name:         "subdomain.test.com",
						Status:       "Active",
						IsCommonName: false,
					},
				},
				int64(3): {
					{
						ID:           3,
						Name:         "subdomain2.domain.com",
						Status:       "Active",
						IsCommonName: false,
					},
				},
			},
		},
		{
			name:         "Empty input",
			inputDomains: make([]*models.Domain, 0),
			expected:     make(map[int64][]*models.Domain, 0),
		},
	}

	for _, c := range cases {
		actual := getDomainGroups(c.inputDomains)

		if !reflect.DeepEqual(actual, c.expected) {
			// deep.Equal doesn't compare pointer values, so we just use it to
			// generate a human friendly diff
			diff := deep.Equal(actual, c.expected)
			t.Errorf("Diff: %+v", diff)
			t.Fatalf("%s: Expected %+v but got %+v",
				c.name,
				c.expected,
				actual,
			)
		}
	}
}

type updaterFlags struct {
	UpdateDomains              bool
	UpdateNotificationSettings bool
	UpdateDCVMethod            bool
	UpdateOrganization         bool
}

func TestGetUpdater(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		expectErr  bool
		statusFunc func(params certificate.CertificateGetCertificateStatusParams) (*certificate.CertificateGetCertificateStatusOK, error)
		state      CertificateState
		want       updaterFlags
	}{
		{
			name:       "Status Processing",
			expectErr:  false,
			statusFunc: mockStatusFunc("Processing"),
			want: updaterFlags{
				UpdateDomains:              false,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         false,
			},
		},
		{
			name:       "Status DomainControlValidation",
			expectErr:  false,
			statusFunc: mockStatusFunc("DomainControlValidation"),
			want: updaterFlags{
				UpdateDomains:              false,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         false,
			},
		},
		{
			name:       "Status OtherValidation",
			expectErr:  false,
			statusFunc: mockStatusFunc("OtherValidation"),
			want: updaterFlags{
				UpdateDomains:              false,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         false,
			},
		},
		{
			name:       "Status Deployment",
			expectErr:  false,
			statusFunc: mockStatusFunc("Deployment"),
			state: CertificateState{
				OrganizationChanged: false,
			},
			want: updaterFlags{
				UpdateDomains:              true,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         false,
			},
		},
		{
			name:       "Status Active",
			expectErr:  false,
			statusFunc: mockStatusFunc("Active"),
			state: CertificateState{
				OrganizationChanged: true,
			},
			want: updaterFlags{
				UpdateDomains:              true,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         true,
			},
		},
		{
			name:      "Error path API Error",
			expectErr: true,
			statusFunc: func(params certificate.CertificateGetCertificateStatusParams) (*certificate.CertificateGetCertificateStatusOK, error) {
				return nil, errors.New("some api error")
			},
		},
		{
			name:       "Error path deleted certificate",
			expectErr:  true,
			statusFunc: mockStatusFunc("Deleted"),
		},
		{
			name:       "Error path unknown status",
			expectErr:  true,
			statusFunc: mockStatusFunc("unknown status"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSvc := sdkcps.CpsService{
				Certificate: mockCertificateService{
					funcCertificateGetCertificateStatus: tt.statusFunc,
				},
			}

			got, err := GetUpdater(mockSvc, tt.state)

			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectErr && err == nil {
				t.Fatal("expected error but got none")
			}

			if tt.expectErr && err != nil {
				t.Logf("got error as expected (test is OK): %v", err)
				return // successful error case
			}

			if got.UpdateDCVMethod != tt.want.UpdateDCVMethod {
				t.Fatalf("UpdateDCVMethod: expected %t but got %t", tt.want.UpdateDCVMethod, got.UpdateDCVMethod)
			}

			if got.UpdateDomains != tt.want.UpdateDomains {
				t.Fatalf("UpdateDomains: expected %t but got %t", tt.want.UpdateDomains, got.UpdateDomains)
			}

			if got.UpdateNotificationSettings != tt.want.UpdateNotificationSettings {
				t.Fatalf("UpdateNotificationSettings: expected %t but got %t", tt.want.UpdateNotificationSettings, got.UpdateNotificationSettings)
			}

			if got.UpdateOrganization != tt.want.UpdateOrganization {
				t.Fatalf("UpdateOrganization: expected %t but got %t", tt.want.UpdateOrganization, got.UpdateOrganization)
			}
		})
	}
}

func TestUpdaterUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args updaterFlags
	}{
		{
			name: "Status Processing",
			args: updaterFlags{
				UpdateDomains:              false,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         true,
			},
		},
		{
			name: "Status DomainControlValidation",
			args: updaterFlags{
				UpdateDomains:              false,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         false,
			},
		},
		{
			name: "Status OtherValidation",
			args: updaterFlags{
				UpdateDomains:              false,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         false,
			},
		},
		{
			name: "Status Deployment",
			args: updaterFlags{
				UpdateDomains:              true,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         true,
			},
		},
		{
			name: "Status Active",
			args: updaterFlags{
				UpdateDomains:              true,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         true,
			},
		},
		{
			name: "Null path no flags",
			args: updaterFlags{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockNotifFunc := NewMockUpdateNotificationsFunc()
			mockOrgFunc := NewMockUpdateOrgFunc()

			mockSvc := sdkcps.CpsService{
				Certificate: mockCertificateService{
					funcCertificateUpdateRequestNotifications: mockNotifFunc,
					funcCertificatePutOrganizationDetails:     mockOrgFunc,
				},
			}

			updater := CertUpdater{
				Svc: mockSvc,
				State: CertificateState{
					CertificateID:        1,
					NotificationSettings: make([]*models.EmailNotification, 0),
					Organization:         &models.OrganizationDetail{},
				},
				UpdateDomains:              tt.args.UpdateDomains,
				UpdateNotificationSettings: tt.args.UpdateNotificationSettings,
				UpdateDCVMethod:            tt.args.UpdateDCVMethod,
				UpdateOrganization:         tt.args.UpdateOrganization,
			}

			updater.Update()

			// If UpdateNotificationSettings, assert that call did occurred.
			if tt.args.UpdateNotificationSettings && len(mockNotifFunc.ParamsPassed) == 0 {
				t.Fatal("expected UpdateNotificationSettings to be called")
			}

			// If !UpdateNotificationSettings, assert call did not occurr.
			if !tt.args.UpdateNotificationSettings && len(mockNotifFunc.ParamsPassed) > 0 {
				t.Fatal("expected no call to UpdateNotificationSettings")
			}

			// If UpdateOrganiztion, assert that call did occurred.
			if tt.args.UpdateOrganization && len(mockOrgFunc.ParamsPassed) == 0 {
				t.Fatal("expected UpdateOrganization to be called")
			}

			// If !UpdateOrganiztion, assert call did not occurr.
			if !tt.args.UpdateOrganization && len(mockOrgFunc.ParamsPassed) > 0 {
				t.Fatal("expected no call to UpdateOrganization")
			}
		})
	}
}

type mockCertificateService struct {
	funcCertificateGetCertificateStatus       func(params certificate.CertificateGetCertificateStatusParams) (*certificate.CertificateGetCertificateStatusOK, error)
	funcCertificateUpdateRequestNotifications *MockUpdateNotificationsFunc
	funcCertificatePutOrganizationDetails     *MockUpdateOrgFunc
}

func (svc mockCertificateService) CertificateCancel(params certificate.CertificateCancelParams) (*certificate.CertificateCancelNoContent, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificateDelete(params certificate.CertificateDeleteParams) (*certificate.CertificateDeleteNoContent, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificateFind(params certificate.CertificateFindParams) (*certificate.CertificateFindOK, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificateGet(params certificate.CertificateGetParams) (*certificate.CertificateGetOK, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificateGetCertificateStatus(params certificate.CertificateGetCertificateStatusParams) (*certificate.CertificateGetCertificateStatusOK, error) {
	if svc.funcCertificateGetCertificateStatus != nil {
		return svc.funcCertificateGetCertificateStatus(params)
	}

	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificateGetRequestNotifications(params certificate.CertificateGetRequestNotificationsParams) (*certificate.CertificateGetRequestNotificationsOK, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificatePatch(params certificate.CertificatePatchParams) (*certificate.CertificatePatchOK, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificatePost(params certificate.CertificatePostParams) (*certificate.CertificatePostCreated, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificatePutOrganizationDetails(params certificate.CertificatePutOrganizationDetailsParams) (*certificate.CertificatePutOrganizationDetailsOK, error) {
	if svc.funcCertificatePutOrganizationDetails != nil {
		return svc.funcCertificatePutOrganizationDetails.Func(params)
	}
	// default implementation
	return nil, nil
}

type MockUpdateOrgFunc struct {
	ParamsPassed []certificate.CertificatePutOrganizationDetailsParams
	Func         func(
		params certificate.CertificatePutOrganizationDetailsParams,
	) (*certificate.CertificatePutOrganizationDetailsOK, error)
}

func NewMockUpdateOrgFunc() *MockUpdateOrgFunc {
	mf := &MockUpdateOrgFunc{
		ParamsPassed: make([]certificate.CertificatePutOrganizationDetailsParams, 0),
	}

	mf.Func = func(
		params certificate.CertificatePutOrganizationDetailsParams,
	) (*certificate.CertificatePutOrganizationDetailsOK, error) {
		log.Printf("MockUpdateOrgFunc called with %v", params)
		mf.ParamsPassed = append(mf.ParamsPassed, params)
		log.Printf("MockUpdateOrgFunc calls: %d", len(mf.ParamsPassed))

		return &certificate.CertificatePutOrganizationDetailsOK{
			OrganizationDetail: *params.OrgDetails,
		}, nil
	}

	return mf
}

func (svc mockCertificateService) CertificatePutRenewal(params certificate.CertificatePutRenewalParams) (*certificate.CertificatePutRenewalNoContent, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificatePutRetrigger(params certificate.CertificatePutRetriggerParams) (*certificate.CertificatePutRetriggerNoContent, error) {
	// default implementation
	return nil, nil
}

func (svc mockCertificateService) CertificateUpdateRequestNotifications(params certificate.CertificateUpdateRequestNotificationsParams) (*certificate.CertificateUpdateRequestNotificationsOK, error) {
	if svc.funcCertificateUpdateRequestNotifications != nil {
		return svc.funcCertificateUpdateRequestNotifications.Func(params)
	}
	// default implementation
	return nil, nil
}

func mockStatusFunc(status string) func(params certificate.CertificateGetCertificateStatusParams) (*certificate.CertificateGetCertificateStatusOK, error) {
	return func(params certificate.CertificateGetCertificateStatusParams) (*certificate.CertificateGetCertificateStatusOK, error) {
		return &certificate.CertificateGetCertificateStatusOK{
			CertificateStatus: models.CertificateStatus{Status: status},
		}, nil
	}
}

type MockUpdateNotificationsFunc struct {
	ParamsPassed []certificate.CertificateUpdateRequestNotificationsParams
	Func         func(
		params certificate.CertificateUpdateRequestNotificationsParams,
	) (*certificate.CertificateUpdateRequestNotificationsOK, error)
}

func NewMockUpdateNotificationsFunc() *MockUpdateNotificationsFunc {
	mf := &MockUpdateNotificationsFunc{
		ParamsPassed: make([]certificate.CertificateUpdateRequestNotificationsParams, 0),
	}

	mf.Func = func(
		params certificate.CertificateUpdateRequestNotificationsParams,
	) (*certificate.CertificateUpdateRequestNotificationsOK, error) {
		log.Printf("MockUpdateNotificationsFunc called with %v", params)
		mf.ParamsPassed = append(mf.ParamsPassed, params)
		log.Printf("MockUpdateNotificationsFunc calls: %d", len(mf.ParamsPassed))

		return &certificate.CertificateUpdateRequestNotificationsOK{
			HyperionCollectionEmailNotification: models.HyperionCollectionEmailNotification{
				Items:      params.Notifications,
				TotalItems: int32(len(params.Notifications)),
			},
		}, nil
	}

	return mf
}
