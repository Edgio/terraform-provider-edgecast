// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package cps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"terraform-provider-edgecast/edgecast/api"
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps"
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
		Schema:        GetCertificateSchema(),
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
		return helper.CreationErrorf(d, "failed to load configuration")
	}

	svc, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Read from TF state.
	certState, errs := ExpandCertificate(d)
	if len(errs) > 0 {
		return helper.DiagsFromErrors("error parsing certificate", errs)
	}

	ns, errs := ExpandNotifSettings(d.Get("notification_setting"))

	if len(errs) > 0 {
		return helper.CreationErrors(d, "error parsing notification_setting", errs)
	}

	// Call APIs.
	cparams := certificate.NewCertificatePostParams()
	cparams.Certificate = &models.CertificateCreate{
		AutoRenew:            certState.AutoRenew,
		CertificateAuthority: certState.CertificateAuthority,
		CertificateLabel:     certState.CertificateLabel,
		DcvMethod:            certState.DcvMethod,
		Description:          certState.Description,
		Domains:              certState.Domains,
		Organization:         certState.Organization,
		ValidationType:       certState.ValidationType,
	}

	cresp, err := svc.Certificate.CertificatePost(cparams)
	if err != nil {
		return helper.CreationError(d, err)
	}

	nparams := certificate.NewCertificateUpdateRequestNotificationsParams()
	nparams.ID = cresp.ID
	nparams.Notifications = ns

	_, err = svc.Certificate.CertificateUpdateRequestNotifications(nparams)
	if err != nil {
		d.SetId("")

		// Cancel the created cert request.
		cnclParams := certificate.NewCertificateCancelParams()
		cnclParams.ID = cresp.ID
		cnclParams.Apply = true
		_, cancelErr := svc.Certificate.CertificateCancel(cnclParams)
		if cancelErr != nil {
			return diag.Errorf(
				"failed to roll back cert request upon error: %v, original err: %v",
				cancelErr.Error(),
				err.Error())
		}

		return diag.Errorf("failed to create certificate: %v", err)
	}

	log.Printf("[INFO] certificate created: %# v\n", pretty.Formatter(cresp))
	log.Printf("[INFO] certificate id: %d\n", cresp.ID)

	d.SetId(strconv.Itoa(int(cresp.ID)))

	return ResourceCertificateRead(ctx, d, m)
}

func ExpandCertificate(
	d *schema.ResourceData,
) (*CertificateState, []error) {
	if d == nil {
		return nil, []error{errors.New("no data to read")}
	}

	errs := make([]error, 0)
	certState := &CertificateState{}

	if v, ok := d.GetOk("auto_renew"); ok {
		if autoRenew, ok := v.(bool); ok {
			certState.AutoRenew = autoRenew
		} else {
			errs = append(errs, errors.New("auto_renew not a bool"))
		}
	}

	if v, ok := d.GetOk("certificate_authority"); ok {
		if certAuthority, ok := v.(string); ok {
			certState.CertificateAuthority = certAuthority
		} else {
			errs = append(errs, errors.New("certificate_authority not a string"))
		}
	}

	if v, ok := d.GetOk("certificate_label"); ok {
		if certLabel, ok := v.(string); ok {
			certState.CertificateLabel = certLabel
		} else {
			errs = append(errs, errors.New("certificate_label not a string"))
		}
	}

	if v, ok := d.GetOk("dcv_method"); ok {
		if dcvMethod, ok := v.(string); ok {
			certState.DcvMethod = dcvMethod
		} else {
			errs = append(errs, errors.New("dcv_method not a string"))
		}
	}

	if v, ok := d.GetOk("description"); ok {
		if desc, ok := v.(string); ok {
			certState.Description = desc
		} else {
			errs = append(errs, errors.New("description not a string"))
		}
	}

	if v, ok := d.GetOk("validation_type"); ok {
		if validationType, ok := v.(string); ok {
			certState.ValidationType = validationType
		} else {
			errs = append(errs, errors.New("validation_type not a string"))
		}
	}

	if v, ok := d.GetOk("domain"); ok {
		if domains, err := ExpandDomains(v); err == nil {
			certState.Domains = domains
		} else {
			errs = append(errs, fmt.Errorf("error parsing domains: %w", err))
		}
	}

	if v, ok := d.GetOk("organization"); ok {
		if org, err := ExpandOrganization(v); err == nil {
			certState.Organization = org
		} else {
			errs = append(errs, fmt.Errorf("error parsing organization: %w", err))
		}
	}

	if v, ok := d.GetOk("notification_setting"); ok {
		if ns, nserrs := ExpandNotifSettings(v); len(errs) == 0 {
			certState.NotificationSettings = ns
		} else {
			errs = append(errs, nserrs...)
		}
	}

	log.Printf("cert state: %# v\n", pretty.Formatter(certState))

	return certState, errs
}

func ResourceCertificateRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config, ok := m.(**api.ClientConfig)
	if !ok {
		return diag.Errorf("failed to load configuration")
	}

	svc, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	certID, err := helper.ParseInt64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Call APIs.
	log.Printf("[INFO] Retrieving certificate : ID: %d\n", certID)

	params := certificate.NewCertificateGetParams()
	params.ID = certID

	resp, err := svc.Certificate.CertificateGet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved certificate: %# v\n", pretty.Formatter(resp))

	nparams := certificate.NewCertificateGetRequestNotificationsParams()
	nparams.ID = certID

	nresp, err := svc.Certificate.CertificateGetRequestNotifications(nparams)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Retrieved certificate notification settings: %# v\n",
		pretty.Formatter(resp))

	// Write TF state.
	err = setCertificateState(d, resp, nresp)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func setCertificateState(
	d *schema.ResourceData,
	resp *certificate.CertificateGetOK,
	nresp *certificate.CertificateGetRequestNotificationsOK,
) error {
	d.Set("certificate_label", resp.CertificateLabel)
	d.Set("description", resp.Description)
	d.Set("dcv_method", resp.DcvMethod)
	d.Set("validation_type", resp.ValidationType)

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

	flattendDomains := FlattenDomains(resp.Domains)
	d.Set("domain", flattendDomains)

	flattendOrganization := FlattenOrganization(resp.Organization)
	d.Set("organization", flattendOrganization)

	if resp.CreatedBy != nil {
		flattenedCreatedBy := FlattenActor(resp.CreatedBy)
		d.Set("created_by", flattenedCreatedBy)
	}

	if resp.ModifiedBy != nil {
		flattenedModifiedBy := FlattenActor(resp.ModifiedBy)
		d.Set("modified_by", flattenedModifiedBy)
	}

	if nresp != nil {
		flattenedNotifSettings := FlattenNotifSettings(nresp.Items)
		d.Set("notification_setting", flattenedNotifSettings)
	}

	return nil
}

func FlattenNotifSettings(
	notifSettings []*models.EmailNotification,
) []map[string]any {
	flattened := make([]map[string]any, len(notifSettings))

	for ix, n := range notifSettings {
		m := make(map[string]any)
		m["notification_type"] = n.NotificationType
		m["enabled"] = n.Enabled

		if len(n.Emails) > 0 {
			m["emails"] = n.Emails
		} else {
			m["emails"] = make([]string, 0)
		}

		flattened[ix] = m
	}

	return flattened
}

func ResourceCertificateUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize CPS Service
	config, ok := m.(**api.ClientConfig)
	if !ok {
		return helper.CreationErrorf(d, "failed to load configuration")
	}

	svc, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	certState, errs := ExpandCertificate(d)
	if len(errs) > 0 {
		return helper.DiagsFromErrors("error parsing certificate", errs)
	}

	certID, err := helper.ParseInt64(d.Id())
	if err != nil {
		return helper.DiagFromError("id was not an int64", err)
	}

	certState.CertificateID = certID

	updater, err := GetUpdater(*svc, *certState)
	if err != nil {
		return helper.DiagFromError("failed to determine update flow", err)
	}

	err = updater.Update()
	if err != nil {
		return helper.DiagFromError("failed to update certificate", err)
	}

	return ResourceCertificateRead(ctx, d, m)
}

func ResourceCertificateDelete(
	ctx context.Context,
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

	//Get certificate status
	statusParams := certificate.NewCertificateGetCertificateStatusParams()
	statusParams.ID = certID
	statusResp, err :=
		cpsService.Certificate.CertificateGetCertificateStatus(statusParams)
	if err != nil {
		return diag.FromErr(err)
	}

	if statusResp.Status == "Processing" &&
		statusResp.OrderValidation == nil {

		//Certificate has not been placed yet.
		cancelParams := certificate.NewCertificateCancelParams()
		cancelParams.ID = certID
		cancelParams.Apply = true
		_, err := cpsService.Certificate.CertificateCancel(cancelParams)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] Canceled Certificate ID: %v", certID)

	} else if (statusResp.Status == "DomainControlValidation" ||
		statusResp.Status == "OtherValidation") &&
		(statusResp.OrderValidation != nil &&
			statusResp.OrderValidation.Status == "Pending") {

		//Certificate has been placed, but not issued yet.
		cancelParams := certificate.NewCertificateCancelParams()
		cancelParams.ID = certID
		cancelParams.Apply = true
		_, err := cpsService.Certificate.CertificateCancel(cancelParams)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] Canceled Certificate ID: %v", certID)

	} else {

		//certificate has been issued.
		deleteParams := certificate.NewCertificateDeleteParams()
		deleteParams.ID = certID
		_, err := cpsService.Certificate.CertificateDelete(deleteParams)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] Deleted Certificate ID: %v", certID)
	}

	d.SetId("")

	return diag.Diagnostics{}
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

// ExpandNotifSettings converts the Terraform representation of organization
// into the EmailNotification API Model.
func ExpandNotifSettings(
	attr interface{},
) ([]*models.EmailNotification, []error) {
	if attr == nil {
		return make([]*models.EmailNotification, 0), nil
	}

	tfSet, ok := attr.(*schema.Set)
	if !ok {
		return nil, []error{errors.New("error parsing notification settings")}
	}

	if tfSet == nil {
		return make([]*models.EmailNotification, 0), nil
	}

	maps := tfSet.List()

	// Empty map
	if len(maps) == 0 {
		return make([]*models.EmailNotification, 0), nil
	}

	errs := make([]error, 0)
	emailNotifs := make([]*models.EmailNotification, 0)

	for ix, v := range maps {
		m, ok := v.(map[string]any)
		if !ok {
			errs = append(
				errs,
				fmt.Errorf("error parsing notification_setting[%d]", ix))

			continue
		}

		emailNotif := &models.EmailNotification{}

		if notifType, ok := helper.GetStringFromMap(m, "notification_type"); ok {
			emailNotif.NotificationType = notifType
		} else {
			err := fmt.Errorf(
				"error parsing notification_setting[%d].notification_type",
				ix)
			errs = append(errs, err)
		}

		if enabled, ok := helper.GetBoolFromMap(m, "enabled"); ok {
			emailNotif.Enabled = enabled
		} else {
			err := fmt.Errorf(
				"error parsing notification_setting[%d].enabled",
				ix)
			errs = append(errs, err)
		}

		emails, err := helper.ConvertTFCollectionToStrings(m["emails"])
		if err == nil {
			emailNotif.Emails = emails
		} else {
			err := fmt.Errorf(
				"error parsing notification_setting[%d].emails: %w",
				ix,
				err)
			errs = append(errs, err)
		}

		emailNotifs = append(emailNotifs, emailNotif)
	}

	return emailNotifs, errs
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

func FlattenDomains(
	domains []*models.Domain,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range domains {
		m := make(map[string]interface{})

		m["id"] = v.ID
		m["is_common_name"] = v.IsCommonName
		m["name"] = v.Name
		m["status"] = v.Status

		tActive, _ := time.Parse(datetimeFormat, v.ActiveDate.String())
		m["active_date"] = tActive.Format(datetimeFormat)

		tCreated, _ := time.Parse(datetimeFormat, v.Created.String())
		m["created"] = tCreated.Format(datetimeFormat)

		flattened = append(flattened, m)
	}

	return flattened
}

func FlattenOrganization(
	organization *models.OrganizationDetail,
) []map[string]interface{} {
	if organization == nil {
		return make([]map[string]interface{}, 0)
	}
	flattened := make([]map[string]interface{}, 0)

	m := make(map[string]interface{})

	m["city"] = organization.City
	m["company_address"] = organization.CompanyAddress
	m["company_address2"] = organization.CompanyAddress2
	m["company_name"] = organization.CompanyName
	m["contact_email"] = organization.ContactEmail
	m["contact_first_name"] = organization.ContactFirstName
	m["contact_last_name"] = organization.ContactLastName
	m["contact_phone"] = organization.ContactPhone
	m["contact_title"] = organization.ContactTitle
	m["country"] = organization.Country
	m["id"] = organization.ID
	m["organizational_unit"] = organization.OrganizationalUnit
	m["state"] = organization.State
	m["zip_code"] = organization.ZipCode
	if organization.AdditionalContacts != nil {
		m["additional_contact"] = flattenAdditionalContacts(organization.AdditionalContacts)
	}

	flattened = append(flattened, m)
	return flattened
}

func flattenAdditionalContacts(
	additionalcontacts []*models.OrganizationContact,
) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range additionalcontacts {
		m := make(map[string]interface{})

		m["contact_type"] = v.ContactType
		m["email"] = v.Email
		m["first_name"] = v.FirstName
		m["id"] = v.ID
		m["last_name"] = v.LastName
		m["phone"] = v.Phone
		m["title"] = v.Title

		flattened = append(flattened, m)
	}

	return flattened
}

// CertificateState represents the state of a certificate as it exists in the
// TF state file. This is an intermediate model before being translated to API
// models.
type CertificateState struct {
	// Certificate ID
	CertificateID int64

	// auto renew
	AutoRenew bool

	// certificate authority
	CertificateAuthority string

	// certificate label
	CertificateLabel string

	// dcv method.
	// Enum: [Email DnsCnameToken DnsTxtToken]
	DcvMethod string

	// description
	Description string

	// domains
	Domains []*models.DomainCreateUpdate

	// organization
	Organization *models.OrganizationDetail

	// validation type
	// Enum: [None DV OV EV]
	ValidationType string

	// notification settings
	NotificationSettings []*models.EmailNotification
}

func GetUpdater(
	svc cps.CpsService,
	state CertificateState,
) (*CertUpdater, error) {
	params := certificate.NewCertificateGetCertificateStatusParams()
	params.ID = state.CertificateID

	resp, err := svc.Certificate.CertificateGetCertificateStatus(params)
	if err != nil {
		return nil, fmt.Errorf("error retreiving certificate status: %w", err)
	}

	if strings.EqualFold(resp.Status, "Deleted") {
		return nil, errors.New("attempted to update a certificate that is deleted")
	}

	if strings.EqualFold(resp.Status, "Processing") {
		return &CertUpdater{
			svc:                        svc,
			state:                      state,
			UpdateDomains:              false,
			UpdateNotificationSettings: true,
			UpdateDCVMethod:            true,
			UpdateOrganization:         true,
		}, nil
	}

	if strings.EqualFold(resp.Status, "DomainControlValidation") ||
		strings.EqualFold(resp.Status, "OtherValidation") {
		return &CertUpdater{
			svc:                        svc,
			state:                      state,
			UpdateDomains:              false,
			UpdateNotificationSettings: true,
			UpdateDCVMethod:            true,
			UpdateOrganization:         false,
		}, nil
	}

	if strings.EqualFold(resp.Status, "Deployment") ||
		strings.EqualFold(resp.Status, "Active") {
		return &CertUpdater{
			svc:                        svc,
			state:                      state,
			UpdateDomains:              true,
			UpdateNotificationSettings: true,
			UpdateDCVMethod:            true,
			UpdateOrganization:         true,
		}, nil
	}

	return nil, errors.New("unknown update flow")
}

type CertUpdater struct {
	svc                        cps.CpsService
	state                      CertificateState
	UpdateDomains              bool
	UpdateNotificationSettings bool
	UpdateDCVMethod            bool
	UpdateOrganization         bool
}

func (u CertUpdater) Update() error {
	if err := u.updateBasicSettings(); err != nil {
		return err
	}

	if err := u.updateNotificationSettings(); err != nil {
		return err
	}

	if err := u.updateDCVMethod(); err != nil {
		return err
	}

	if err := u.updateOrganization(); err != nil {
		return err
	}

	return nil
}

func (u CertUpdater) updateBasicSettings() error {
	params := certificate.NewCertificatePatchParams()
	params.ID = u.state.CertificateID
	params.CertificateRequest = &models.CertificateUpdate{
		AutoRenew:        u.state.AutoRenew,
		CertificateLabel: u.state.CertificateLabel,
		Description:      u.state.Description,
		DcvMethod:        u.state.DcvMethod,
	}

	if u.UpdateDomains {
		params.CertificateRequest.Domains = u.state.Domains
	}

	resp, err := u.svc.Certificate.CertificatePatch(params)
	if err != nil {
		return fmt.Errorf("failed to update certificate: %w", err)
	}

	log.Printf("[INFO] certificate updated: %# v\n", pretty.Formatter(resp))

	return nil
}

func (u CertUpdater) updateNotificationSettings() error {
	if u.UpdateNotificationSettings {
		// not yet implemeted
	} else {
		log.Printf("[INFO] Skipped updating notification settings")
	}

	return nil
}

func (u CertUpdater) updateDCVMethod() error {
	if u.UpdateNotificationSettings {
		// not yet implemeted
	} else {
		log.Printf("[INFO] Skipped updating DCV method")
	}

	return nil
}

func (u CertUpdater) updateOrganization() error {
	if u.UpdateNotificationSettings {
		// not yet implemeted
	} else {
		log.Printf("[INFO] Skipped updating organization")
	}

	return nil
}
