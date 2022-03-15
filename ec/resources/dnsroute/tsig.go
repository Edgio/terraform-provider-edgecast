// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"log"
	"strconv"
	"strings"
	"terraform-provider-edgecast/ec/helper"

	"terraform-provider-edgecast/ec/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTsig() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceTsigCreate,
		ReadContext:   ResourceTsigRead,
		UpdateContext: ResourceTsigUpdate,
		DeleteContext: ResourceTsigDelete,
		Importer:      helper.Import(ResourceTsigRead, "account_number", "id"),

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration."},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alias."},
			"key_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "tsig key name"},
			"key_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "tsig value"},
			"algorithm_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "tsig encryption type:[HMAC-MD5,HMAC-SHA1,HMAC-SHA256,HMAC-SHA384,HMAC-SHA224,HMAC-SHA512]"},
		},
	}
}

func ResourceTsigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	alias := d.Get("alias").(string)
	keyName := d.Get("key_name").(string)
	keyValue := d.Get("key_value").(string)
	algorithm := strings.ToLower(d.Get("algorithm_name").(string))
	algorithmID := 0

	switch algorithm {
	case "hmac-md5":
		algorithmID = api.TsigAlgorithm_HMAC_MD5
	case "hmac-sha1":
		algorithmID = api.TsigAlgorithm_HMAC_SHA1
	case "hmac-sha256":
		algorithmID = api.TsigAlgorithm_HMAC_SHA256
	case "hmac-sha384":
		algorithmID = api.TsigAlgorithm_HMAC_SHA384
	case "hmac-sha224":
		algorithmID = api.TsigAlgorithm_HMAC_SHA224
	case "hmac-sha512":
		algorithmID = api.TsigAlgorithm_HMAC_SHA512
	}
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	tsigRequest := &api.DnsRouteTsig{
		ID:            -1,
		Alias:         alias,
		KeyName:       keyName,
		KeyValue:      keyValue,
		AlgorithmID:   algorithmID,
		AlgorithmName: algorithm,
	}

	log.Printf("[INFO] Creating a new TSig for Account '%s': %+v", accountNumber, tsigRequest)

	dnsrouteClient := api.NewDNSRouteAPIClient(*config)

	resp, err := dnsrouteClient.AddTsig(tsigRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Create successful - New TSig ID: %d", resp)
	d.SetId(strconv.Itoa(resp))

	ResourceTsigRead(ctx, d, m)

	return diags
}

func ResourceTsigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tsigID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)

	log.Printf("[INFO] Retrieving Master Server Group by tsigID: %d", tsigID)

	resp, err := dnsRouteClient.GetTsig(tsigID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	newId := strconv.Itoa(resp.ID)

	d.SetId(newId)
	return diags
}

func ResourceTsigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tsigID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)
	alias := d.Get("alias").(string)
	keyName := d.Get("key_name").(string)
	keyValue := d.Get("key_value").(string)
	algorithm := strings.ToLower(d.Get("algorithm_name").(string))
	algorithmID := 0

	switch algorithm {
	case "hmac-md5":
		algorithmID = api.TsigAlgorithm_HMAC_MD5
	case "hmac-sha1":
		algorithmID = api.TsigAlgorithm_HMAC_SHA1
	case "hmac-sha256":
		algorithmID = api.TsigAlgorithm_HMAC_SHA256
	case "hmac-sha384":
		algorithmID = api.TsigAlgorithm_HMAC_SHA384
	case "hmac-sha224":
		algorithmID = api.TsigAlgorithm_HMAC_SHA224
	case "hmac-sha512":
		algorithmID = api.TsigAlgorithm_HMAC_SHA512
	}
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	tsigRequest := &api.DnsRouteTsig{
		ID:            tsigID,
		Alias:         alias,
		KeyName:       keyName,
		KeyValue:      keyValue,
		AlgorithmID:   algorithmID,
		AlgorithmName: algorithm,
	}

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)

	err = dnsRouteClient.UpdateTsig(tsigRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceTsigRead(ctx, d, m)
}

func ResourceTsigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteAPIClient := api.NewDNSRouteAPIClient(*config)

	tsigID, _ := strconv.Atoi(d.Id())

	err := dnsRouteAPIClient.DeleteTsig(tsigID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
