// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps_test

import (
	"errors"
	"reflect"
	"sort"
	"testing"

	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/resources/cps"

	sdkcps "github.com/EdgeCast/ec-sdk-go/edgecast/cps"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
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
						map[string]interface{}{
							"first_name":   "contact2",
							"last_name":    "lastname2",
							"email":        "first.lastname2@testuser.com",
							"phone":        "111-111-3333",
							"title":        "Developer",
							"contact_type": "TechnicalContact",
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
		actualPtr, err := cps.ExpandOrganization(v.input)

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
		actualPtr, err := cps.ExpandAdditionalContacts(v.input)

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
		actualPtr, err := cps.ExpandDomains(v.input)

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
		actual := cps.FlattenDeployments(c.input)

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
		actual := cps.FlattenActor(c.input)

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

			got, errs := cps.ExpandNotifSettings(tt.args)

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

			got := cps.FlattenNotifSettings(tt.args)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
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
				UpdateOrganization:         true,
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
			want: updaterFlags{
				UpdateDomains:              true,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         true,
			},
		},
		{
			name:       "Status Active",
			expectErr:  false,
			statusFunc: mockStatusFunc("Active"),
			want: updaterFlags{
				UpdateDomains:              true,
				UpdateNotificationSettings: true,
				UpdateDCVMethod:            true,
				UpdateOrganization:         true,
			},
		},
		{
			name:       "Status Active",
			expectErr:  false,
			statusFunc: mockStatusFunc("Active"),
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

			got, err := cps.GetUpdater(mockSvc, cps.CertificateState{})

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

type mockCertificateService struct {
	funcCertificateGetCertificateStatus func(params certificate.CertificateGetCertificateStatusParams) (*certificate.CertificateGetCertificateStatusOK, error)
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
	// default implementation
	return nil, nil
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
