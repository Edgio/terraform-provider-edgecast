// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package internal

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderConfig holds configuration values for the provider.
type ProviderConfig struct {
	APIToken         string `json:"-"` // sensitive.
	AccountNumber    string
	IdsClientID      string `json:"-"` // sensitive.
	IdsClientSecret  string `json:"-"` // sensitive.
	IdsScope         string `json:"-"` // sensitive.
	IDSAddress       string
	APIAddress       string
	APIAddressLegacy string
	APIURL           *url.URL
	IdsURL           *url.URL
	APIURLLegacy     *url.URL
	PartnerID        int
	PartnerUserID    int
	UserAgent        string
}

// ExpandProviderConfig reads ProviderConfig using the TF Resource Data.
func ExpandProviderConfig(d *schema.ResourceData) (*ProviderConfig, error) {
	config := &ProviderConfig{
		APIToken:         d.Get("api_token").(string),
		AccountNumber:    d.Get("account_number").(string),
		IdsClientID:      d.Get("ids_client_id").(string),
		IdsClientSecret:  d.Get("ids_client_secret").(string),
		IdsScope:         d.Get("ids_scope").(string),
		IDSAddress:       d.Get("ids_address").(string),
		APIAddress:       d.Get("api_address").(string),
		APIAddressLegacy: d.Get("api_address_legacy").(string),
	}

	var err error

	config.IdsURL, err = url.Parse(config.IDSAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IDS URL: %w", err)
	}

	config.APIURL, err = url.Parse(config.APIAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	config.APIURLLegacy, err = url.Parse(config.APIAddressLegacy)
	if err != nil {
		return nil, fmt.Errorf("failed to parse legacy API URL: %w", err)
	}

	if partnerUserIDValue, ok := d.GetOk("partner_user_id"); ok {
		config.PartnerUserID = partnerUserIDValue.(int)
	}

	if partnerIDValue, ok := d.GetOk("partner_id"); ok {
		config.PartnerID = partnerIDValue.(int)
	}

	return config, nil
}
