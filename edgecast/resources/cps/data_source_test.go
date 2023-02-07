// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package cps_test

import (
	"strconv"
	"testing"

	"terraform-provider-edgecast/edgecast/resources/cps"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/appendix"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/go-test/deep"
)

func Test_FlattenCountryCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args *appendix.AppendixGetOK
		want []map[string]interface{}
	}{
		{
			name: "Happy Path",
			args: &appendix.AppendixGetOK{
				HyperionCollectionCountryCode: models.HyperionCollectionCountryCode{
					Items: []*models.CountryCode{
						{
							Country:       "United States of America",
							TwoLetterCode: "US",
						},
						{
							Country:       "Switzerland",
							TwoLetterCode: "CH",
						},
					},
				},
			},
			want: []map[string]interface{}{
				{
					"country":         "United States of America",
					"two_letter_code": "US",
				},
				{
					"country":         "Switzerland",
					"two_letter_code": "CH",
				},
			},
		},
		{
			name: "No countries",
			args: &appendix.AppendixGetOK{},
			want: make([]map[string]interface{}, 0),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := cps.FlattenCountries(tt.args)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}

func Test_FlattenNamedEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args *models.HyperionCollectionNamedEntity
		want []map[string]interface{}
	}{
		{
			name: "Happy Path",
			args: createNamedEntities(2),
			want: createNamedTFEntities(2),
		},
		{
			name: "No entities",
			args: &models.HyperionCollectionNamedEntity{},
			want: make([]map[string]interface{}, 0),
		},
		{
			name: "Nil entities",
			args: nil,
			want: make([]map[string]interface{}, 0),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := cps.FlattenNamedEntities(tt.args)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}

func Test_CheckForRetry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		argsMetadata []*models.DomainDcvFull
		argsStatus   *certificate.CertificateGetCertificateStatusOK
		want         bool
	}{
		{
			name: "Happy Path - status is not processing, dcv token is available",
			argsMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: "token"},
					DomainID:  1,
					Emails:    nil,
				},
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: "token"},
					DomainID:  2,
					Emails:    nil,
				},
			},
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "DomainControlValidation"},
			},
			want: false,
		},
		{
			name:         "status is processing, metadata is empty",
			argsMetadata: nil,
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "Processing"},
			},
			want: true,
		},
		{
			name: "status is processing, dcv token is empty",
			argsMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  nil,
					DomainID:  1,
					Emails:    nil,
				},
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  nil,
					DomainID:  2,
					Emails:    nil,
				},
			},
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "Processing"},
			},
			want: true,
		},
		{
			name: "status is processing, dcv token is empty string",
			argsMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: ""},
					DomainID:  1,
					Emails:    nil,
				},
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: ""},
					DomainID:  2,
					Emails:    nil,
				},
			},
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "Processing"},
			},
			want: true,
		},
		{
			name: "status is processing, dcv token is available",
			argsMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: "token"},
					DomainID:  1,
					Emails:    nil,
				},
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: "token"},
					DomainID:  2,
					Emails:    nil,
				},
			},
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "Processing"},
			},
			want: true,
		},
		{
			name:         "status is not processing, metadata is empty",
			argsMetadata: nil,
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "DomainControlValidation"},
			},
			want: true,
		},
		{
			name: "status is not processing, dcv token is empty",
			argsMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  nil,
					DomainID:  1,
					Emails:    nil,
				},
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  nil,
					DomainID:  2,
					Emails:    nil,
				},
			},
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "DomainControlValidation"},
			},
			want: true,
		},
		{
			name: "status is not processing, dcv token is empty string",
			argsMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: ""},
					DomainID:  1,
					Emails:    nil,
				},
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: ""},
					DomainID:  2,
					Emails:    nil,
				},
			},
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: "DomainControlValidation"},
			},
			want: true,
		},
		{
			name: "status is empty, dcv token is available",
			argsMetadata: []*models.DomainDcvFull{
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: "token"},
					DomainID:  1,
					Emails:    nil,
				},
				{
					DcvMethod: "DnsCnameToken",
					DcvToken:  &models.DcvToken{Token: "token"},
					DomainID:  2,
					Emails:    nil,
				},
			},
			argsStatus: &certificate.CertificateGetCertificateStatusOK{
				CertificateStatus: models.CertificateStatus{
					Status: ""},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := cps.CheckForRetry(tt.argsMetadata, tt.argsStatus)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}

func createNamedEntities(
	count int,
) *models.HyperionCollectionNamedEntity {
	items := make([]*models.NamedEntity, count)
	for i := 0; i < count; i++ {
		items[i] = &models.NamedEntity{
			ID:   int64(i),
			Name: "item " + strconv.Itoa(count),
		}
	}

	return &models.HyperionCollectionNamedEntity{
		Items: items,
	}
}

func createNamedTFEntities(count int) []map[string]interface{} {
	items := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = map[string]interface{}{
			"id":   int64(i),
			"name": "item " + strconv.Itoa(count),
		}
	}

	return items
}
