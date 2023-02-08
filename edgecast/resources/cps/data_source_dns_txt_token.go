// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
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
	readTokenDefaultTimeout = "20m"
)

func DataSourceDNSTXTToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceDNSTXTTokenRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the system-defined ID assigned to a certificate.",
			},
			"wait_until_available": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether Terraform should wait until the token is available.",
			},
			"wait_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          readTokenDefaultTimeout,
				ValidateDiagFunc: internal.ValidateDuration,
				Description:      "Indicates the maximum time Terraform will wait (e.g. `60m` for 60 minutes, `10s` for ten seconds, or `2h` for two hours). If `wait_until_available` is not set, this value is ignored.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The token value through which you may prove control over your certificate request's domains.",
			},
		},
	}
}

func DataSourceDNSTXTTokenRead(
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

	log.Printf("[INFO] Retrieving certificate Status: ID: %d\n", certID)
	statusparams := certificate.NewCertificateGetCertificateStatusParams()
	statusparams.ID = certID

	retry := d.Get("wait_until_available").(bool)
	log.Printf("wait_until_available: %t\n", retry)
	log.Printf("timeout: %v\n", d.Timeout(schema.TimeoutRead))

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

			statusresp, err := svc.Certificate.CertificateGetCertificateStatus(statusparams)
			if err != nil {
				return resource.NonRetryableError(
					fmt.Errorf(
						"error while retrieving certificate details: %w",
						err))
			}

			// If cert is not DV, return error.
			if resp.ValidationType != models.CdnProvidedCertificateValidationTypeDV {
				return resource.NonRetryableError(errors.New("certificate must have validation type DV"))
			}

			// If workflow error, return error.
			if len(resp.WorkflowErrorMessage) > 0 {
				return resource.NonRetryableError(
					fmt.Errorf(
						"error in workflow: %s",
						resp.WorkflowErrorMessage))
			}

			metadata := GetDomainMetadata(resp, svc)

			// No token found.
			retryErr := CheckForDCVTokenRetry(retry, metadata, statusresp)
			if retryErr == nil {
				// All of the domains have the same token, so take the first.
				log.Printf("setting token to %s", metadata[0].DcvToken.Token)
				d.Set("value", metadata[0].DcvToken.Token)
				d.SetId(helper.GetUnixTimeStamp())
			}

			return retryErr
		})

	return diag.FromErr(err)
}

// CheckForTokenRetry determines whether the provider should check for a dcv
// token again.
func CheckForDCVTokenRetry(
	doRetry bool,
	metadata []*models.DomainDcvFull,
	statusresp *certificate.CertificateGetCertificateStatusOK,
) *resource.RetryError {
	// if token is available, and status is not processing, no retry is needed.
	if statusresp.Status != "" &&
		strings.ToLower(statusresp.Status) != "processing" &&
		metadata != nil && len(metadata) > 0 &&
		metadata[0].DcvToken != nil && len(metadata[0].DcvToken.Token) > 0 {
		return nil
	}

	log.Println("token not availale")

	if doRetry {
		log.Println("retrying")
		return resource.RetryableError(
			errors.New("token not available"))
	}

	// Just exit if retry is not desired.
	// The user will need to run refresh to try again.
	log.Println("not retrying")
	return nil
}
