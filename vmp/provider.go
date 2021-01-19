// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license. See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"errors"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider creates a new instance of the Verizon Media Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.vdms.io/v2/",
			},
			"api_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids_client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids_client_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"partner_user_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"partner_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vmp_origin":        resourceOrigin(),
			"vmp_cname":         resourceCname(),
			"vmp_customer":      resourceCustomer(),
			"vmp_customer_user": resourceCustomerUser(),
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

	apiBaseURI := d.Get("api_address").(string)
	apiToken := d.Get("api_token").(string)
	accountNumber := d.Get("account_number").(string)
	idsClientID := d.Get("ids_client_id").(string)
	idsClientSecret := d.Get("ids_client_secret").(string)
	idsScope := d.Get("ids_scope").(string)

	// Must use either API Token or IDS credentials, but not both
	idsProvided := len(idsClientID) > 0 && len(idsClientSecret) > 0 && len(idsScope) > 0
	apiTokenProvided := len(apiToken) > 0

	if (apiTokenProvided && idsProvided) || (!apiTokenProvided && !idsProvided) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Please provide either an API Token or IDS credentials, but not both",
		})

		return nil, diags
	}

	apiClient, err := api.NewApiClient(apiBaseURI, apiToken, idsClientID, idsClientSecret, idsScope)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create API Client",
		})

		return nil, diags
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
