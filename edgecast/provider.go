// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package edgecast

import (
	"context"
	"fmt"

	"terraform-provider-edgecast/edgecast/api"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/resources/cps"
	"terraform-provider-edgecast/edgecast/resources/customer"
	"terraform-provider-edgecast/edgecast/resources/dnsroute"
	"terraform-provider-edgecast/edgecast/resources/edgecname"
	"terraform-provider-edgecast/edgecast/resources/origin"
	"terraform-provider-edgecast/edgecast/resources/rulesengine"
	"terraform-provider-edgecast/edgecast/resources/waf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	apiURLProd       string = "https://api.vdms.io"
	apiURLProdLegacy string = "https://api.edgecast.com"
	idsURLProd       string = "https://id.vdms.io"

	// Version indicates the current version of this provider
	Version string = "0.5.9"

	userAgentFormat = "edgecast/terraform-provider:%s"
)

// Provider creates a new instance of the Edgecast Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:               getProviderSchema(),
		ResourcesMap:         buildResourcesMap(),
		DataSourcesMap:       buildDataSourcesMap(),
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(
	tx context.Context,
	d *schema.ResourceData,
) (interface{}, diag.Diagnostics) {
	config, err := api.ExpandClientConfig(d)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Summary:  "Failed to read provider configuration.",
				Severity: diag.Error,
				Detail:   err.Error(),
			},
		}
	}

	config.UserAgent = fmt.Sprintf(userAgentFormat, Version)
	helper.LogInstanceAsPrettyJson("Provider Configuration", config)

	return &config, nil
}

func getProviderSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "API Token for managing the following resources: Origin, CNAME, Customer, Customer User",
		},
		"ids_client_secret": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "OAuth 2.0 Client Secret for managing the following resources: Rules Engine Policy",
		},
		"ids_client_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "OAuth 2.0 Client ID for managing the following resources: Rules Engine Policy",
		},
		"ids_scope": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "OAuth 2.0 Scopes for managing the following resources: Rules Engine Policy",
		},
		"account_number": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Account Number to use when only managing a single customer's resources. If managing multiple customers, this parameter should be omitted.",
		},
		"partner_user_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Partner User ID to impersonate. If using PCC or MCC credentials, this parameter will be ignored.",
		},
		"partner_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Partner ID to impersonate. If using PCC or MCC credentials, this parameter will be ignored.",
		},
		"api_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The base url of Edgecast resource APIs. Omit to use the default url. For internal testing.",
			Default:     apiURLProd,
		},
		"ids_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The base url of Edgecast identity APIs. Omit to use the default url. For internal testing.",
			Default:     idsURLProd,
		},
		"api_address_legacy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The base url of legacy Edgecast resource APIs. Omit to use the default url. For internal testing.",
			Default:     apiURLProdLegacy,
		},
	}
}

func buildResourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgecast_origin":                 origin.ResourceOrigin(),
		"edgecast_edgecname":              edgecname.ResourceEdgeCname(),
		"edgecast_customer":               customer.ResourceCustomer(),
		"edgecast_customer_user":          customer.ResourceCustomerUser(),
		"edgecast_rules_engine_policy":    rulesengine.ResourceRulesEngineV4Policy(),
		"edgecast_dns_masterservergroup":  dnsroute.ResourceMasterServerGroup(),
		"edgecast_dns_zone":               dnsroute.ResourceZone(),
		"edgecast_dns_group":              dnsroute.ResourceGroup(),
		"edgecast_dns_tsig":               dnsroute.ResourceTsig(),
		"edgecast_dns_secondaryzonegroup": dnsroute.ResourceSecondaryZoneGroup(),
		"edgecast_waf_access_rule":        waf.ResourceAccessRule(),
		"edgecast_waf_rate_rule":          waf.ResourceRateRule(),
		"edgecast_waf_managed_rule":       waf.ResourceManagedRule(),
		"edgecast_waf_custom_rule_set":    waf.ResourceCustomRuleSet(),
		"edgecast_waf_scopes":             waf.ResourceScopes(),
		"edgecast_waf_bot_rule_set":       waf.ResourceBotRuleSet(),
		"edgecast_cps_certificate":        cps.ResourceCertificate(),
	}
}

func buildDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgecast_customer_services":               customer.DataSourceCustomerServices(),
		"edgecast_cps_countrycodes":                cps.DataSourceCountryCodes(),
		"edgecast_cps_dcv_types":                   cps.DataSourceDCVTypes(),
		"edgecast_cps_domain_statuses":             cps.DataSourceDomainStatuses(),
		"edgecast_cps_validation_statuses":         cps.DataSourceValidationStatuses(),
		"edgecast_cps_cert_validation_levels":      cps.DataSourceCertValidationLevels(),
		"edgecast_cps_cert_request_cancel_actions": cps.DataSourceCancelCertReqActions(),
		"edgecast_cps_cert_request_statuses":       cps.DataSourceCertReqStatuses(),
		"edgecast_cps_cert_order_statuses":         cps.DataSourceCertOrderStatuses(),
	}
}
