// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license. See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	apiURLProd string = "https://api.vdms.io"
	idsURLProd string = "https://id.vdms.io"
)

// Provider creates a new instance of the Verizon Media Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token":         {Type: schema.TypeString, Optional: true},
			"ids_client_secret": {Type: schema.TypeString, Optional: true},
			"ids_client_id":     {Type: schema.TypeString, Optional: true},
			"ids_scope":         {Type: schema.TypeString, Optional: true},
			"account_number":    {Type: schema.TypeString, Optional: true},
			"partner_user_id":   {Type: schema.TypeInt, Optional: true},
			"partner_id":        {Type: schema.TypeInt, Optional: true},
			"api_address":       {Type: schema.TypeString, Optional: true},
			"ids_address":       {Type: schema.TypeString, Optional: true},
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

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err error
	config, err := api.NewClientConfig(
		d.Get("api_token").(string),
		d.Get("account_number").(string),
		d.Get("ids_client_id").(string),
		d.Get("ids_client_secret").(string),
		d.Get("ids_scope").(string),
		d.Get("api_address").(string),
		d.Get("ids_address").(string),
	)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Failed to read vmp Provider configuration data: %v", err))
	}

	config.BaseClient = api.NewBaseClient(config)
	log.Printf("config[api_token]:%s", config.APIToken)
	log.Printf("config[ids_client_id]:%s", config.IdsClientID)
	log.Printf("config[ids_client_secret]:%s", config.IdsClientSecret)
	log.Printf("config[ids_scope]:%s", config.IdsScope)
	log.Printf("config[api_address]:%s", config.APIURL)
	log.Printf("config[ids_address]:%s", config.IdsURL)
	log.Printf("config[account_number]:%s", config.AccountNumber)
	log.Printf("config[BaseClient]:%v", config.BaseClient)

	if partnerUserIDValue, ok := d.GetOk("partner_user_id"); ok {
		partnerUserID := partnerUserIDValue.(int)
		config.PartnerUserID = partnerUserID
	}

	if partnerIDValue, ok := d.GetOk("partner_id"); ok {
		partnerID := partnerIDValue.(int)
		config.PartnerID = partnerID
	}
	return &config, diags
}
