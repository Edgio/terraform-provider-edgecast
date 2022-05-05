// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"log"
	"strconv"
	"strings"

	"terraform-provider-edgecast/edgecast/api"

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
				Description: `Account Number associated with the customer whose 
				resources you wish to manage. This account number may be found 
				in the upper right-hand corner of the MCC.`},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates a brief description for the TSIG key."},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Identifies the key on the master name server and 
				our Route name servers. This name must be unique.`},
			"key_value": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Identifies a hash value through which our name 
				servers will be authenticated to a master name server.`},
			"algorithm_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Identifies a cryptographic hash function name. 
				Options: HMAC-MD5 | HMAC-SHA1 | HMAC-SHA256 | HMAC-SHA384 | 
				HMAC-SHA224 | HMAC-SHA512`},
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
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// Construct TSIG Object
	tsig := expandTSIG(d)

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
	if err != nil {
		return diag.FromErr(err)
	}

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
	d.Set("alias", tsigObj.Alias)
	d.Set("key_name", tsigObj.KeyName)
	d.Set("key_value", tsigObj.KeyValue)
	d.Set("algorithm_name", getAlgorithmNameFromID(*tsigObj))

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
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Construct TSIG Update Object
	updatedTsigObj := expandTSIG(d)

	// Get Existing TSIG Object
	getParams := routedns.NewGetTSIGParams()
	getParams.AccountNumber = accountNumber
	getParams.TSIGID = tsigID

	tsigObj, err := routeDNSService.GetTSIG(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Apply updated TSIG data
	tsigObj.Alias = updatedTsigObj.Alias
	tsigObj.KeyName = updatedTsigObj.KeyName
	tsigObj.KeyValue = updatedTsigObj.KeyValue
	tsigObj.AlgorithmID = updatedTsigObj.AlgorithmID

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
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

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

func expandTSIG(d *schema.ResourceData) routedns.TSIG {
	alias := d.Get("alias").(string)
	keyName := d.Get("key_name").(string)
	keyValue := d.Get("key_value").(string)

	// Convert Algorithm Name to Algorithm ID needed by API
	rawAlgorithm := strings.ToUpper(d.Get("algorithm_name").(string))
	var algorithm routedns.TSIGAlgorithmType

	switch rawAlgorithm {
	case "HMAC-MD5":
		algorithm = routedns.HMAC_MD5
	case "HMAC-SHA1":
		algorithm = routedns.HMAC_SHA1
	case "HMAC-SHA256":
		algorithm = routedns.HMAC_SHA256
	case "HMAC-SHA384":
		algorithm = routedns.HMAC_SHA384
	case "HMAC-SHA224":
		algorithm = routedns.HMAC_SHA224
	case "HMAC-SHA512":
		algorithm = routedns.HMAC_SHA512
	}

	tsig := routedns.TSIG{
		Alias:       alias,
		KeyName:     keyName,
		KeyValue:    keyValue,
		AlgorithmID: algorithm,
	}

	return tsig
}

func getAlgorithmNameFromID(tsig routedns.TSIGGetOK) string {
	// Convert Algorithm ID to Algorithm Name used in resource file
	var algorithmName string

	switch tsig.AlgorithmID {
	case routedns.HMAC_MD5:
		algorithmName = "HMAC-MD5"
	case routedns.HMAC_SHA1:
		algorithmName = "HMAC-SHA1"
	case routedns.HMAC_SHA256:
		algorithmName = "HMAC-SHA256"
	case routedns.HMAC_SHA384:
		algorithmName = "HMAC-SHA384"
	case routedns.HMAC_SHA224:
		algorithmName = "HMAC-SHA224"
	case routedns.HMAC_SHA512:
		algorithmName = "HMAC-SHA512"
	}

	return algorithmName
}
