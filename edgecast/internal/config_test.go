// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package internal_test

import (
	"net/url"
	"testing"

	"terraform-provider-edgecast/edgecast"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandProviderConfig(t *testing.T) {
	t.Parallel()

	apiTok := uuid.New().String()
	accountNumber := uuid.New().String()[0:4]
	clientID := uuid.New().String()
	clientSecret := uuid.New().String()
	clientScope := "rules"
	idsAddress := "https://ids.example.com"
	apiAddress := "https://api.example.com"
	apiAddressLegacy := "https://apilegacy.example.com"
	partnerUserID := 999
	partnerID := 888
	idsURL, _ := url.Parse(idsAddress)
	apiURL, _ := url.Parse(apiAddress)
	apiLegacyURL, _ := url.Parse(apiAddressLegacy)

	tests := []struct {
		name        string
		arg         map[string]any
		want        *internal.ProviderConfig
		expectError bool
	}{
		{
			name: "Happy Path",
			arg: map[string]any{
				"api_token":          apiTok,
				"account_number":     accountNumber,
				"ids_client_id":      clientID,
				"ids_client_secret":  clientSecret,
				"ids_scope":          clientScope,
				"ids_address":        idsAddress,
				"api_address":        apiAddress,
				"api_address_legacy": apiAddressLegacy,
				"partner_user_id":    partnerUserID,
				"partner_id":         partnerID,
			},
			want: &internal.ProviderConfig{
				APIToken:         apiTok,
				AccountNumber:    accountNumber,
				IdsClientID:      clientID,
				IdsClientSecret:  clientSecret,
				IdsScope:         clientScope,
				IDSAddress:       idsAddress,
				APIAddress:       apiAddress,
				APIAddressLegacy: apiAddressLegacy,
				IdsURL:           idsURL,
				APIURL:           apiURL,
				APIURLLegacy:     apiLegacyURL,
				PartnerUserID:    partnerUserID,
				PartnerID:        partnerID,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := schema.TestResourceDataRaw(
				t,
				edgecast.GetProviderSchema(),
				tt.arg)
			got, err := internal.ExpandProviderConfig(data)

			if tt.expectError && err != nil {
				return // expected error, successful test
			}

			if tt.expectError && err == nil {
				t.Fatal("expected error, but got none")
			}

			if !tt.expectError && err != nil {
				t.Fatal("unexpected error")
			}

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}
