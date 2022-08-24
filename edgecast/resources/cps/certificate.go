// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package cps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"terraform-provider-edgecast/edgecast/api"
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

func ResourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceCertificateCreate,
		ReadContext:   ResourceCertificateRead,
		UpdateContext: ResourceCertificateUpdate,
		DeleteContext: ResourceCertificateDelete,
		//Importer:      helper.Import(ResourceCertificateRead, "id"),

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the system-defined ID assigned to this certificate.",
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Indicates the timestamp at which this request for a certificate was initially submitted. \n" +
					"Syntax: \n" +
					"{YYYY}-{MM}-{DD}T{hh}:{mm}:{ss}.{ffffff}Z",
			},
			"created_by": {
				Type:     schema.TypeSet,
				Computed: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"portal_type_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "[ Customer, Partner, Wholesaler, Uber, OpenCdn ]",
						},
						"identity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "[ User, Client ]",
						},
					},
				},
				Description: "Describes the user that submitted this certificate request.",
			},
			"deployments": {
				Type:     schema.TypeSet,
				Computed: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delivery_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "[ GlobalPremiumPlusAsia, NorthAmericaAndEurope, GlobalStandard, Internal, GlobalPremiumAsiaPlusChina, GlobalPremiumAsiaPlusIndia, GlobalPremiumAsiaPlusChinaAndIndia, GlobalPremiumAsiaPlusLatam, GlobalPremiumAsiaPremiumChinaPlusLatam ]",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "[ HttpLarge, HttpSmall, Adn ]",
						},
						"hex_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "Returns a null value.",
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Indicates the timestamp at which this certificate will expire. \n" +
					"Syntax: \n" +
					"{YYYY}-{MM}-{DD}T{hh}:{mm}:{ss}.{ffffff}Z \n" +
					"If the Certificate Authority (CA) is still processing the certificate request, then this property returns the following timestamp: \n" +
					"0001-01-01T00:00:00Z",
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Indicates the timestamp at which this request for a certificate was last modified. \n" +
					"Syntax: \n" +
					"{YYYY}-{MM}-{DD}T{hh}:{mm}:{ss}.{ffffff}Z ",
			},
			"modified_by": {
				Type:     schema.TypeSet,
				Computed: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"portal_type_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "[ Customer, Partner, Wholesaler, Uber, OpenCdn ]",
						},
						"identity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "[ User, Client ]",
						},
					},
				},
				Description: "Returns a null value.",
			},
			"request_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returns Enterprise.",
			},
			"thumbprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returns a null value.",
			},
			"workflow_error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returns a null value.",
			},

			"auto_renew": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether this certificate will automatically renew prior to its expiration date.",
			},
			"certificate_authority": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set to DigiCert.",
			},
			"certificate_label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sets the certificate's name. Specify a name that solely consists of alphanumeric characters, underscores, and dashes.",
			},
			"dcv_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Determines the method through which your control over the domains associated with this certificate will be validated.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the certificate's description.",
			},
			"domain": {
				Type:     schema.TypeList,
				Required: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_common_name": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "Determines whether this domain corresponds to the certificate's common name. \n" +
								"Note: You may only designate a single domain as the common name.  \n" +
								"Default Value:  \n" +
								"If you do not designate a domain as the common name, then our system will assign it to one of your domains.",
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Sets the domain name. \n" +
								"Example  \n" +
								"cdn.example.com",
						},
					},
				},
				Description: "Contains the certificate's domain(s).",
			},
			"organization": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_contact": {
							Type:     schema.TypeList,
							Optional: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"contact_type": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Required for EV certificates. \n" +
											"Set to EVApprover.",
									},
									"email": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Required for EV certificates. \n" +
											"Sets the email address for the current contact.",
									},
									"first_name": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Required for EV certificates. \n" +
											"Sets the first name for the current contact.",
									},
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Reserved for future use.",
									},
									"last_name": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Required for EV certificates. \n" +
											"Sets the last name for the current contact.",
									},
									"phone": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Required for EV certificates. \n" +
											"Sets the phone number for the current contact.",
									},
									"title": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Required for EV certificates. \n" +
											"Sets the title of the current contact.",
									},
								},
							},
							Description: "Required for EV certificates. \n" +
								"Contains additional contacts that are also responsible for validating certificates for this organization.",
						},
						"city": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the organization's city.",
						},
						"company_address": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the organization's address.",
						},
						"company_address2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sets the organization's secondary address information (e.g., suite number).",
						},
						"company_name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the organization's name.\n" +
								"Note: If we are unable to identify an organization through the id property, then we will compare the name specified in this property to all of your organizations. If an exact match is found, then the certificate request will be associated with that organization. Additionally, all other properties defined within this object will be ignored. \n" +
								"Note: If we cannot identify an existing organization through either the id or company_name properties, then we will create a new organization using the information supplied in this object.",
						},
						"contact_email": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the email address for the individual responsible for validating certificates for this organization.",
						},
						"contact_first_name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the first name for the individual responsible for validating certificates for this organization.",
						},
						"contact_last_name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the last name for the individual responsible for validating certificates for this organization.",
						},
						"contact_phone": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the phone number for the individual responsible for validating certificates for this organization.",
						},
						"contact_title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sets the title of the individual responsible for validating certificates for this organization.",
						},
						"country": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Required for OV and EV certificates. \n" +
								"Sets the organization's country by its ISO 3166 country code.",
						},
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "Identifies an organization by its system-defined ID. \n" +
								"Key information: \n" +
								"Specify an existing organization by passing either of the following values: \n" +
								"ID: Set this property to the system-defined ID for the desired organization. \n" +
								"0: Set this property to 0 if the desired organization has only been registered with Digicert. Additionally, you must set the company_name property to your organization's exact name as defined within Digicert. \n" +
								"You cannot modify an existing organization. If you assign an existing organization to this certificate, then other properties defined within this object will be ignored.",
						},
						"organizational_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sets the name of the organizational unit responsible for validating certificates for this organization.",
						},
						"state": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "United States Only: Required for OV and EV certificates. \n" +
								"Sets the organization's state by its abbreviation.",
						},
						"zip_code": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "United States Only: Required for OV and EV certificates. \n" +
								"Sets the organization's zip code.",
						},
					},
				},
				Description: "Required for OV and EV certificates. \n" +
					"Describes the certificate request's organization. \n" +
					"Note: Do not specify an organization for DV certificates.",
			},
			"validation_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Determines the certificate's level of validation.",
			},
		},
	}
}

func ResourceCertificateCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize CPS Service
	config := m.(**api.ClientConfig)
	cpsService, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	//create certificate object
	certificateObj := models.CertificateCreate{
		AutoRenew:            d.Get("auto_renew").(bool),
		CertificateAuthority: d.Get("certificate_authority").(string),
		CertificateLabel:     d.Get("certificate_label").(string),
		DcvMethod:            d.Get("dcv_method").(string),
		Description:          d.Get("description").(string),
		ValidationType:       d.Get("validation_type").(string),
	}

	domains, err := expandDomains(d.Get("domain"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing domains: %w", err))
	}

	certificateObj.Domains = domains

	organization, err := expandOrganization(d.Get("organization"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing organization: %w", err))
	}

	certificateObj.Organization = organization

	// Call create certificate API
	params := certificate.NewCertificatePostParams()
	params.Certificate = &certificateObj

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

func ResourceCertificateRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	cpsService, err := buildCPSService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	certID, _ := strconv.ParseInt(d.Id(), 10, 64)

	log.Printf(
		"[INFO] Retriving certificate : ID: %v",
		certID,
	)

	params := certificate.NewCertificateGetParams()
	params.ID = certID
	resp, err := cpsService.Certificate.CertificateGet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved certificate: %# v", pretty.Formatter(resp))

	d.SetId(strconv.Itoa(int(resp.ID)))
	d.Set("certificate_label", resp.CertificateLabel)
	d.Set("description", resp.Description)
	d.Set("last_modified", resp.LastModified.String())
	d.Set("created", resp.Created.String())
	d.Set("expiration_date", resp.ExpirationDate.String())
	d.Set("request_type", resp.RequestType)
	d.Set("thumbprint", resp.Thumbprint)
	d.Set("workflow_error_message", resp.WorkflowErrorMessage)
	d.Set("auto_renew", resp.AutoRenew)

	flattenedDeployments := flattenDeployments(resp.Deployments)
	d.Set("deployments", flattenedDeployments)

	if resp.CreatedBy != nil {
		flattenedCreatedBy := flattenActor(resp.CreatedBy)
		d.Set("created_by", flattenedCreatedBy)
	}

	if resp.ModifiedBy != nil {
		flattenedModifiedBy := flattenActor(resp.ModifiedBy)
		d.Set("modified_by", flattenedModifiedBy)
	}

	return diag.Diagnostics{}
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

// expandDomains converts the Terraform representation of Domains into
// the Domains API Model
func expandDomains(attr interface{}) ([]*models.DomainCreateUpdate, error) {
	if items, ok := attr.([]interface{}); ok {

		domains := make([]*models.DomainCreateUpdate, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			domain := models.DomainCreateUpdate{
				Name:         curr["name"].(string),
				IsCommonName: curr["is_common_name"].(bool),
			}

			domains = append(domains, &domain)
		}

		return domains, nil

	} else {
		return nil,
			errors.New("ExpandDomains: attr input was not a []interface{}")
	}
}

// expandOrganization converts the Terraform representation of organization
// into the Organization API Model
func expandOrganization(attr interface{}) (*models.OrganizationDetail, error) {
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
		additionalContacts, err :=
			expandAdditionalContacts(curr["additional_contact"])
		if err != nil {
			return nil, err
		}
		organization.AdditionalContacts = additionalContacts
	}

	return &organization, nil
}

// expandAdditionalContacts converts the Terraform representation of
// organization contacts into the OrganizationContact API Model
func expandAdditionalContacts(attr interface{}) ([]*models.OrganizationContact, error) {
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
		return nil,
			errors.New("expandAdditionalContacts: attr input was not a []interface{}")
	}
}

// FlattenActor converts the Actor API Model
// into a format that Terraform can work with
func flattenActor(actor *models.Actor) []map[string]interface{} {
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
// into a format that Terraform can work with
func flattenDeployments(deployments []*models.RequestDeployment) []map[string]interface{} {
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
