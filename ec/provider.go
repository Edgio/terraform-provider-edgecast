// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package ec

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-edgecast/ec/api"
	"terraform-provider-edgecast/ec/resources/customer"
	"terraform-provider-edgecast/ec/resources/dnsroute"
	"terraform-provider-edgecast/ec/resources/edgecname"
	"terraform-provider-edgecast/ec/resources/origin"
	"terraform-provider-edgecast/ec/resources/rulesengine"
	"terraform-provider-edgecast/ec/resources/waf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	apiURLProd       string = "https://api.vdms.io"
	apiURLProdLegacy string = "https://api.edgecast.com"
	idsURLProd       string = "https://id.vdms.io"

	// Version indicates the current version of this provider
	Version string = "0.4.4"

	userAgentFormat = "edgecast/terraform-provider:%s"
)

// TODO Platforms should be a data source retrieved via an API call, not a local collection
var (
	// Platforms specifies the corresponding Media Type IDs for EdgeCast CDN Delivery Platforms
	Platforms = map[string]int{
		"httplarge": 3,
		"httpsmall": 8,
		"adn":       14,
	}
)

// Provider creates a new instance of the Edgecast Terraform Provider
func Provider() *schema.Provider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API Token for managing the following resources: Origin, CNAME, Customer, Customer User"},
			"ids_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAuth 2.0 Client Secret for managing the following resources: Rules Engine Policy",
			},
			"ids_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAuth 2.0 Client ID for managing the following resources: Rules Engine Policy"},
			"ids_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAuth 2.0 Scopes for managing the following resources: Rules Engine Policy"},
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account Number to use when only managing a single customer's resources. If managing multiple customers, this parameter should be omitted.",
			},
			"partner_user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Partner User ID to impersonate. If using PCC or MCC credentials, this parameter will be ignored."},
			"partner_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Partner ID to impersonate. If using PCC or MCC credentials, this parameter will be ignored."},
			"api_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The base url of Edgecast resource APIs. Omit to use the default url. For internal testing.",
				Default:     apiURLProd},
			"ids_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The base url of Edgecast identity APIs. Omit to use the default url. For internal testing.",
				Default:     idsURLProd},
			"api_address_legacy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The base url of legacy Edgecast resource APIs. Omit to use the default url. For internal testing.",
				Default:     apiURLProdLegacy},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ec_origin":                 origin.ResourceOrigin(),
			"ec_edgecname":              edgecname.ResourceEdgeCname(),
			"ec_customer":               customer.ResourceCustomer(),
			"ec_customer_user":          customer.ResourceCustomerUser(),
			"ec_rules_engine_policy":    rulesengine.ResourceRulesEngineV4Policy(),
			"ec_dns_masterservergroup":  dnsroute.ResourceMasterServerGroup(),
			"ec_dns_zone":               dnsroute.ResourceZone(),
			"ec_dns_group":              dnsroute.ResourceGroup(),
			"ec_dns_tsig":               dnsroute.ResourceTsig(),
			"ec_dns_secondaryzonegroup": dnsroute.ResourceSecondaryZoneGroup(),
			"ec_waf_access_rule":        waf.ResourceAccessRule(),
			"ec_waf_rate_rule":          waf.ResourceRateRule(),
			"ec_waf_managed_rule":       waf.ResourceManagedRule(),
			"ec_waf_custom_rule_set":    waf.ResourceCustomRuleSet(),
			"ec_waf_scopes":             waf.ResourceScopes(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ec_customer_services": customer.DataSourceCustomerServices(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// For debugging purpose
	//time.Sleep(10 * time.Second)

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
		d.Get("api_address_legacy").(string),
	)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to read ec Provider configuration data: %v", err))
	}

	config.BaseClient = api.NewBaseClient(config)
	config.BaseClientLegacy = api.NewLegacyBaseClient(config)

	log.Printf("config[api_token]:%s", config.APIToken)
	log.Printf("config[ids_client_id]:%s", config.IdsClientID)
	log.Printf("config[ids_client_secret]:%s", config.IdsClientSecret)
	log.Printf("config[ids_scope]:%s", config.IdsScope)
	log.Printf("config[api_address]:%s", config.APIURL)
	log.Printf("config[ids_address]:%s", config.IdsURL)
	log.Printf("config[api_address_legacy]:%s", config.APIURLLegacy)
	log.Printf("config[account_number]:%s", config.AccountNumber)
	log.Printf("config[BaseClient]:%v", config.BaseClient)
	log.Printf("config[BaseClientLegacy]:%v", config.BaseClientLegacy)

	if partnerUserIDValue, ok := d.GetOk("partner_user_id"); ok {
		partnerUserID := partnerUserIDValue.(int)
		config.PartnerUserID = partnerUserID
	}

	if partnerIDValue, ok := d.GetOk("partner_id"); ok {
		partnerID := partnerIDValue.(int)
		config.PartnerID = partnerID
	}

	config.UserAgent = fmt.Sprintf(userAgentFormat, Version)

	return &config, diags
}
