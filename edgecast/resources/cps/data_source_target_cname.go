// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package cps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"
	"time"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTargetCNAME() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceTargetCNAMERead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"wait_until_available": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"retry_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the CDN domain through which requests for this certificate will be routed.",
			},
		},
	}
}

func DataSourceTargetCNAMERead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config, ok := m.(internal.ProviderConfig)
	if !ok {
		return diag.Errorf("failed to load configuration")
	}

	svc, err := buildCPSService(config)
	if err != nil {
		return diag.Errorf("failed to build CPS Service: %v", err)
	}

	certID, err := helper.ParseInt64(d.Get("certificate_id").(string))
	if err != nil {
		return diag.Errorf("failed to parse certificate ID: %v", err)
	}

	// Call APIs.
	log.Printf("[INFO] Retrieving certificate : ID: %d\n", certID)

	params := certificate.NewCertificateGetParams()
	params.ID = certID

	retry := d.Get("wait_until_available").(bool)
	log.Printf("wait_until_available: %t\n", retry)
	log.Printf("timeout: %v\n", d.Timeout(schema.TimeoutRead))
	// var retryTime time.Duration
	// if retryTimeRaw, ok := d.GetOk("retry_time"); ok {
	// 	retryTimeString := retryTimeRaw.(string)
	// 	retryTime, err = time.ParseDuration(retryTimeString)
	// 	if err != nil {
	// 		return errors.New("retry_time format invalid")
	// 	}
	// }

	err = resource.RetryContext(
		ctx,
		d.Timeout(schema.TimeoutRead)-time.Minute,
		func() *resource.RetryError {

			// 1. Call API
			resp, err := svc.Certificate.CertificateGet(params)
			if err != nil {
				return resource.NonRetryableError(
					fmt.Errorf(
						"error while retrieving certificate details: %w",
						err))
			}

			// only http large - there should be exactly one
			var deployment *models.RequestDeployment

			for _, d := range resp.Deployments {
				if strings.EqualFold(d.Platform, "HttpLarge") {
					deployment = d
				}
			}

			// No target cname found.
			if deployment == nil || len(deployment.HexURL) == 0 {
				log.Println("target cname not availale")
				if retry {
					log.Println("retrying")
					return resource.RetryableError(errors.New("target cname not available"))
				} else {
					// Just exit if retry is not desired.
					// The user will need to run refresh to try again.
					log.Println("not retrying")
					return nil
				}
			}

			d.Set("value", deployment.HexURL)
			d.SetId(helper.GetUnixTimeStamp())
			return nil
		})

	return diag.FromErr(err)
}

func setDomainsState(
	d *schema.ResourceData,
	resp *certificate.CertificateGetOK,
	dcvresp []*models.DomainDcvFull,
) error {
	return nil
}
