// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package cps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"
	"time"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	readTokenDefaultTimeout = 30 * time.Minute
)

func DataSourceDNSTXTToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceDNSTXTTokenRead,
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
			"wait_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          readTokenDefaultTimeout,
				ValidateDiagFunc: internal.ValidateDuration,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(readTokenDefaultTimeout),
		},
	}
}

func DataSourceDNSTXTTokenRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	timeout := d.Timeout(schema.TimeoutRead)
	log.Printf("read timeout: %v\n", timeout)

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

	err = resource.RetryContext(
		ctx,
		timeout,
		func() *resource.RetryError {
			// 1. Call API
			resp, err := svc.Certificate.CertificateGet(params)
			if err != nil {
				return resource.NonRetryableError(
					fmt.Errorf(
						"error while retrieving certificate details: %w",
						err))
			}

			if resp.ValidationType != models.CdnProvidedCertificateValidationTypeDV {
				return resource.NonRetryableError(errors.New("certificate must have validation type DV"))
			}

			metadata := GetDomainMetadata(resp, svc)

			// No token found.
			if len(metadata) == 0 || metadata[0].DcvToken == nil || len(metadata[0].DcvToken.Token) == 0 {
				log.Println("token not availale")
				if retry {
					log.Println("retrying")
					return resource.RetryableError(errors.New("token not available"))
				} else {
					// Just exit if retry is not desired.
					// The user will need to run refresh to try again.
					log.Println("not retrying")
					return nil
				}
			}

			// All of the domains have the same token, so take the first.
			log.Printf("setting token to %s", metadata[0].DcvToken.Token)
			helper.LogInstanceAsPrettyJson("metadata output", metadata)
			d.Set("value", metadata[0].DcvToken.Token)

			// always run
			d.SetId(helper.GetUnixTimeStamp())
			return nil
		})

	return diag.FromErr(err)
}
