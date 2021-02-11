// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license. See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"errors"
	"fmt"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	apiURLProd string = "https://api.edgecast.com"
	idsURLProd string = "https://id.vdms.io"
)

// Provider creates a new instance of the Verizon Media Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids_client_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"partner_user_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"partner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"api_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  apiURLProd,
			},
			"ids_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  idsURLProd,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vmp_origin":              resourceOrigin(),
			"vmp_cname":               resourceCname(),
			"vmp_customer":            resourceCustomer(),
			"vmp_customer_user":       resourceCustomerUser(),
			"vmp_rules_engine_policy": resourceRulesEngineV4Policy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"vmp_customer_services": dataSourceCustomerServices(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

// ProviderConfiguration contains configuration values for this provider
type ProviderConfiguration struct {
	AccountNumber string
	APIClient     *api.ApiClient
	PartnerUserID *int
	PartnerID     *int
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiToken := d.Get("api_token").(string)
	accountNumber := d.Get("account_number").(string)
	idsClientID := d.Get("ids_client_id").(string)
	idsClientSecret := d.Get("ids_client_secret").(string)
	idsScope := d.Get("ids_scope").(string)

	apiURL := d.Get("api_address").(string)
	idsURL := d.Get("ids_address").(string)

	apiClient, err := api.NewApiClient(apiURL, idsURL, apiToken, idsClientID, idsClientSecret, idsScope)

	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Failed to create API Client: %v", err))
	}

	config := ProviderConfiguration{
		APIClient:     apiClient,
		AccountNumber: accountNumber,
	}

	if partnerUserIDValue, ok := d.GetOk("partner_user_id"); ok {
		partnerUserID := partnerUserIDValue.(int)
		config.PartnerUserID = &partnerUserID
	}

	if partnerIDValue, ok := d.GetOk("partner_id"); ok {
		partnerID := partnerIDValue.(int)
		config.PartnerID = &partnerID
	}

	return &config, diags
}

// ApplyAccountNumberOverride creates a new ProviderConfiguration when an override is desired for Account Number
func (c *ProviderConfiguration) ApplyAccountNumberOverride(accountNumberOverride string) (*ProviderConfiguration, error) {
	if accountNumberOverride == "" {
		accountNumberOverride = c.AccountNumber
	}

	if accountNumberOverride == "" {
		return nil, errors.New("invalid account number")
	}

	return &ProviderConfiguration{
		APIClient:     c.APIClient,
		AccountNumber: accountNumberOverride,
		PartnerUserID: c.PartnerUserID,
		PartnerID:     c.PartnerID,
	}, nil
}
