// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"log"
	"strconv"
	"strings"

	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTsig() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceTsigCreate,
		ReadContext:   ResourceTsigRead,
		UpdateContext: ResourceTsigUpdate,
		DeleteContext: ResourceTsigDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Account Number for the customer if not already 
				specified in the provider configuration.`},
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
				Type:     schema.TypeString,
				Required: true,
				Description: `tsig encryption type:[HMAC-MD5,HMAC-SHA1,
				HMAC-SHA256,HMAC-SHA384,HMAC-SHA224,HMAC-SHA512]`},
		},
	}
}

func ResourceTsigCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Construct TSIG Object
	alias := d.Get("alias").(string)
	keyName := d.Get("key_name").(string)
	keyValue := d.Get("key_value").(string)
	rawAlgorithm := strings.ToLower(d.Get("algorithm_name").(string))
	var algorithm routedns.TSIGAlgorithmType

	switch rawAlgorithm {
	case "hmac-md5":
		algorithm = routedns.HMAC_MD5
	case "hmac-sha1":
		algorithm = routedns.HMAC_SHA1
	case "hmac-sha256":
		algorithm = routedns.HMAC_SHA256
	case "hmac-sha384":
		algorithm = routedns.HMAC_SHA384
	case "hmac-sha224":
		algorithm = routedns.HMAC_SHA224
	case "hmac-sha512":
		algorithm = routedns.HMAC_SHA512
	}

	tsig := routedns.TSIG{
		Alias:       alias,
		KeyName:     keyName,
		KeyValue:    keyValue,
		AlgorithmID: algorithm,
	}

	log.Printf(
		"[INFO] Creating a new TSIG for Account '%s': %+v",
		accountNumber,
		tsig,
	)

	// Call add TSIG API
	params := routedns.NewAddTSIGParams()
	params.AccountNumber = accountNumber
	params.TSIG = tsig

	tsigID, err := routeDNSService.AddTSIG(*params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Create successful - New TSIG ID: %d", tsigID)
	d.SetId(strconv.Itoa(*tsigID))

	return ResourceTsigRead(ctx, d, m)
}

func ResourceTsigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	tsigID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Call get TSIG API
	log.Printf("[INFO] Retrieving TSIG by TSIGID: %d", tsigID)
	params := routedns.NewGetTSIGParams()
	params.AccountNumber = accountNumber
	params.TSIGID = tsigID

	tsigObj, err := routeDNSService.GetTSIG(*params)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved TSIG %+v", tsigObj)

	// TODO: Add handling to process TSIG data into terraform state file

	return diag.Diagnostics{}
}

func ResourceTsigUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	tsigID, err := strconv.Atoi(d.Id())
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Construct TSIG Update Object
	alias := d.Get("alias").(string)
	keyName := d.Get("key_name").(string)
	keyValue := d.Get("key_value").(string)
	rawAlgorithm := strings.ToLower(d.Get("algorithm_name").(string))
	var algorithm routedns.TSIGAlgorithmType

	switch rawAlgorithm {
	case "hmac-md5":
		algorithm = routedns.HMAC_MD5
	case "hmac-sha1":
		algorithm = routedns.HMAC_SHA1
	case "hmac-sha256":
		algorithm = routedns.HMAC_SHA256
	case "hmac-sha384":
		algorithm = routedns.HMAC_SHA384
	case "hmac-sha224":
		algorithm = routedns.HMAC_SHA224
	case "hmac-sha512":
		algorithm = routedns.HMAC_SHA512
	}

	// Get Existing TSIG Object
	getParams := routedns.NewGetTSIGParams()
	getParams.AccountNumber = accountNumber
	getParams.TSIGID = tsigID

	tsigObj, err := routeDNSService.GetTSIG(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Apply updated TSIG data
	tsigObj.Alias = alias
	tsigObj.KeyName = keyName
	tsigObj.KeyValue = keyValue
	tsigObj.AlgorithmID = algorithm

	// Call Update TSIG API
	updateParams := routedns.NewUpdateTSIGParams()
	updateParams.AccountNumber = accountNumber
	updateParams.TSIG = *tsigObj
	err = routeDNSService.UpdateTSIG(*updateParams)

	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceTsigRead(ctx, d, m)
}

func ResourceTsigDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	tsigID, err := strconv.Atoi(d.Id())
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Get Existing TSIG Object
	getParams := routedns.NewGetTSIGParams()
	getParams.AccountNumber = accountNumber
	getParams.TSIGID = tsigID

	tsigObj, err := routeDNSService.GetTSIG(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Delete Existing TSIG Object
	deleteParams := routedns.NewDeleteTSIGParams()
	deleteParams.AccountNumber = accountNumber
	deleteParams.TSIG = *tsigObj

	err = routeDNSService.DeleteTSIG(*deleteParams)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}
