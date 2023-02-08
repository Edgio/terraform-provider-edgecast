// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package cps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	readCNAMEDefaultTimeout = "4h"
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
			"wait_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          readCNAMEDefaultTimeout,
				ValidateDiagFunc: internal.ValidateDuration,
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
	timeoutRaw := d.Get("wait_timeout").(string)
	timeout, err := time.ParseDuration(timeoutRaw)
	if err != nil {
		return diag.Errorf("invalid wait_timeout: %v", err)
	}

	log.Printf("timeout: %v\n", timeout)

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

			// if workflow error, return error
			if len(resp.WorkflowErrorMessage) > 0 {
				return resource.NonRetryableError(
					fmt.Errorf(
						"error in workflow: %s",
						resp.WorkflowErrorMessage))
			}

			// There should be exactly one deployment - for HTTP Large.
			var deployment *models.RequestDeployment

			for _, d := range resp.Deployments {
				if strings.EqualFold(d.Platform, "HttpLarge") {
					deployment = d
					break
				}
			}

			// No target cname found.
			retryErr := CheckForCNAMERetry(retry, deployment)
			if retryErr == nil {
				d.Set("value", deployment.HexURL)
				d.SetId(helper.GetUnixTimeStamp())
			}

			return retryErr
		})

	return diag.FromErr(err)
}

// CheckForCNAMERetry determines whether the provider should check for a target
// CNAME again.
func CheckForCNAMERetry(
	doRetry bool,
	deployment *models.RequestDeployment,
) *resource.RetryError {
	// if a deployment is available with a target cname, no retry is needed.
	if deployment != nil && len(deployment.HexURL) > 0 {
		return nil
	}

	log.Println("target cname not availale")

	if doRetry {
		log.Println("retrying")
		return resource.RetryableError(
			errors.New("target cname not available"))
	}

	// Just exit if retry is not desired.
	// The user will need to run refresh to try again.
	log.Println("not retrying")
	return nil
}
