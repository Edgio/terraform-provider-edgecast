// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package cps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"terraform-provider-edgecast/edgecast/api"
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

const (
	datetimeFormat string = "2006-01-02T15:04:05.000Z07:00"
)

func ResourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceCertificateCreate,
		ReadContext:   ResourceCertificateRead,
		UpdateContext: ResourceCertificateUpdate,
		DeleteContext: ResourceCertificateDelete,
		Importer:      helper.Import(ResourceCertificateRead, "id"),
		Schema:        getCertificateSchema(),
	}
}

func ResourceCertificateCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize CPS Service
	config, ok := m.(**api.ClientConfig)
	if !ok {
		return diag.Errorf("failed to load configuration")
	}

	cpsService, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateObj, errs := ExpandCertificate(d)
	if len(errs) > 0 {
		return helper.DiagsFromErrors(errs)
	}

	// Call create certificate API
	params := certificate.NewCertificatePostParams()
	params.Certificate = certificateObj

	resp, err := cpsService.Certificate.CertificatePost(params)
	if err != nil {
		d.SetId("")

		return diag.FromErr(err)
	}

	log.Printf("[INFO] certificate created: %# v", pretty.Formatter(resp))
	log.Printf("[INFO] certificate id: %d", resp.ID)

	d.SetId(strconv.Itoa(int(resp.ID)))

	return ResourceCertificateRead(ctx, d, m)
}

func ExpandCertificate(
	d *schema.ResourceData,
) (*models.CertificateCreate, []error) {
	errs := make([]error, 0)

	autoRenew, ok := d.Get("auto_renew").(bool)
	if !ok {
		errs = append(errs, errors.New("auto_renew not a bool"))
	}

	certAuthority, ok := d.Get("certificate_authority").(string)
	if !ok {
		errs = append(errs, errors.New("certificate_authority not a string"))
	}

	certLabel, ok := d.Get("certificate_label").(string)
	if !ok {
		errs = append(errs, errors.New("certificate_label not a string"))
	}

	dvcMethod, ok := d.Get("dcv_method").(string)
	if !ok {
		errs = append(errs, errors.New("dcv_method not a string"))
	}

	desc, ok := d.Get("description").(string)
	if !ok {
		errs = append(errs, errors.New("description not a string"))
	}

	validationType, ok := d.Get("validation_type").(string)
	if !ok {
		errs = append(errs, errors.New("validation_type not a string"))
	}

	domains, err := ExpandDomains(d.Get("domain"))
	if err != nil {
		errs = append(errs, fmt.Errorf("error parsing domains: %w", err))
	}

	organization, err := ExpandOrganization(d.Get("organization"))
	if err != nil {
		errs = append(errs, fmt.Errorf("error parsing organization: %w", err))
	}

	return &models.CertificateCreate{
		AutoRenew:            autoRenew,
		CertificateAuthority: certAuthority,
		CertificateLabel:     certLabel,
		DcvMethod:            dvcMethod,
		Description:          desc,
		ValidationType:       validationType,
		Domains:              domains,
		Organization:         organization,
	}, errs
}

func ResourceCertificateRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config, ok := m.(**api.ClientConfig)
	if !ok {
		return diag.Errorf("failed to load configuration")
	}

	cpsService, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	certID, err := helper.ParseInt64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieving certificate : ID: %d", certID)

	params := certificate.NewCertificateGetParams()
	params.ID = certID

	resp, err := cpsService.Certificate.CertificateGet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved certificate: %# v", pretty.Formatter(resp))

	err = setCertificateState(d, resp)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func setCertificateState(
	d *schema.ResourceData,
	resp *certificate.CertificateGetOK,
) error {
	d.Set("certificate_label", resp.CertificateLabel)
	d.Set("description", resp.Description)

	tLastModified, err := time.Parse(datetimeFormat, resp.LastModified.String())
	if err != nil {
		return fmt.Errorf("error parsing cert last modified date: %w", err)
	}

	d.Set("last_modified", tLastModified.Format(datetimeFormat))

	tCreated, err := time.Parse(datetimeFormat, resp.Created.String())
	if err != nil {
		return fmt.Errorf("error parsing cert created date: %w", err)
	}

	d.Set("created", tCreated.Format(datetimeFormat))

	tExpiration, err := time.Parse(datetimeFormat, resp.ExpirationDate.String())
	if err != nil {
		return fmt.Errorf("error parsing cert expiration date: %w", err)
	}

	d.Set("expiration_date", tExpiration.Format(datetimeFormat))

	d.Set("request_type", resp.RequestType)
	d.Set("thumbprint", resp.Thumbprint)

	d.Set("workflow_error_message", resp.WorkflowErrorMessage)
	d.Set("auto_renew", resp.AutoRenew)

	flattenedDeployments := FlattenDeployments(resp.Deployments)
	d.Set("deployments", flattenedDeployments)

	if resp.CreatedBy != nil {
		flattenedCreatedBy := FlattenActor(resp.CreatedBy)
		d.Set("created_by", flattenedCreatedBy)
	}

	if resp.ModifiedBy != nil {
		flattenedModifiedBy := FlattenActor(resp.ModifiedBy)
		d.Set("modified_by", flattenedModifiedBy)
	}

	return nil
}

func ResourceCertificateUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Not Yet Implemented
	return ResourceCertificateRead(ctx, d, m)
}

func ResourceCertificateDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Not Yet Implemented
	var diags diag.Diagnostics

	return diags
}

// ExpandDomains converts the Terraform representation of Domains into
// the Domains API Model.
func ExpandDomains(attr interface{}) ([]*models.DomainCreateUpdate, error) {
	if items, ok := attr.([]interface{}); ok {
		domains := make([]*models.DomainCreateUpdate, 0)

		for _, item := range items {
			curr, ok := item.(map[string]interface{})
			if !ok {
				return nil, errors.New("domain was not map[string]interface{}")
			}

			name, ok := curr["name"].(string)
			if !ok {
				return nil, errors.New("domain.name was not a string")
			}

			isCommonName, ok := curr["is_common_name"].(bool)
			if !ok {
				return nil, errors.New("domain.is_common_name was not a bool")
			}

			domain := models.DomainCreateUpdate{
				Name:         name,
				IsCommonName: isCommonName,
			}

			domains = append(domains, &domain)
		}

		return domains, nil
	} else {
		return nil,
			errors.New("ExpandDomains: attr input was not a []interface{}")
	}
}

// ExpandOrganization converts the Terraform representation of organization
// into the Organization API Model.
func ExpandOrganization(attr interface{}) (*models.OrganizationDetail, error) {
	curr, err := helper.ConvertSingletonSetToMap(attr)
	if err != nil {
		return nil, fmt.Errorf("error expanding orgnization detail: %w", err)
	}

	// Empty map
	if len(curr) == 0 {
		return nil, nil
	}

	organization := models.OrganizationDetail{
		Country:            curr["country"].(string),
		State:              curr["state"].(string),
		ZipCode:            curr["zip_code"].(string),
		City:               curr["city"].(string),
		CompanyName:        curr["company_name"].(string),
		CompanyAddress:     curr["company_address"].(string),
		CompanyAddress2:    curr["company_address2"].(string),
		OrganizationalUnit: curr["organizational_unit"].(string),
		ContactFirstName:   curr["contact_first_name"].(string),
		ContactLastName:    curr["contact_last_name"].(string),
		ContactEmail:       curr["contact_email"].(string),
		ContactPhone:       curr["contact_phone"].(string),
		ContactTitle:       curr["contact_title"].(string),
	}

	if orgID, ok := curr["id"].(int); ok {
		organization.ID = int64(orgID)
	}

	if curr["additional_contact"] != nil {
		additionalContacts, err := ExpandAdditionalContacts(curr["additional_contact"])
		if err != nil {
			return nil, err
		}

		organization.AdditionalContacts = additionalContacts
	}

	return &organization, nil
}

// ExpandAdditionalContacts converts the Terraform representation of
// organization contacts into the OrganizationContact API Model.
func ExpandAdditionalContacts(
	attr interface{},
) ([]*models.OrganizationContact, error) {
	if items, ok := attr.([]interface{}); ok {
		contacts := make([]*models.OrganizationContact, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			contact := models.OrganizationContact{
				ContactType: curr["contact_type"].(string),
				Email:       curr["email"].(string),
				FirstName:   curr["first_name"].(string),
				LastName:    curr["last_name"].(string),
				Phone:       curr["phone"].(string),
				Title:       curr["title"].(string),
			}
			if contactID, ok := curr["id"].(int); ok {
				contact.ID = int64(contactID)
			}

			contacts = append(contacts, &contact)
		}

		return contacts, nil
	} else {
		return nil, errors.New(
			"ExpandAdditionalContacts: attr input was not a []interface{}")
	}
}

// FlattenActor converts the Actor API Model
// into a format that Terraform can work with.
func FlattenActor(actor *models.Actor) []map[string]interface{} {
	if actor == nil {
		return make([]map[string]interface{}, 0)
	}

	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	m["user_id"] = int(actor.UserID)
	m["portal_type_id"] = actor.PortalTypeID
	m["identity_id"] = actor.IdentityID
	m["identity_type"] = actor.IdentityType

	flattened = append(flattened, m)

	return flattened
}

// FlattenDeployments converts the Deployment API Model
// into a format that Terraform can work with.
func FlattenDeployments(
	deployments []*models.RequestDeployment,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range deployments {
		m := make(map[string]interface{})

		m["delivery_region"] = v.DeliveryRegion
		m["hex_url"] = v.HexURL
		m["platform"] = v.Platform

		flattened = append(flattened, m)
	}

	return flattened
}
