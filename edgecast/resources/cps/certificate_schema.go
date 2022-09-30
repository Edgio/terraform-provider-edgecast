// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package cps

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getCertificateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Indicates the system-defined ID assigned to this certificate.",
		},
		"notification_setting": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Determine the conditions under which notifications will be sent and to whom they will be sent for a specific certificate request.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"notification_type": {
						Type:     schema.TypeString,
						Required: true,
						Description: "Identifies the type of notification that will be configured. Valid values are:\n" +
							"CertificateRenewal | CertificateExpiring | PendingValidations",
					},
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Determines whether emails for this type of notification will be sent.",
					},
					"emails": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Required when enabled=true. Defines one or more email addresses to which a notification will be sent. Set this parameter to an email address associated with a MCC user in your account. Your account manager may also define an email address associated with a partner user. Our service returns a 400 Bad Request when this parameter is set to any other email address.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
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
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Determines whether this certificate will automatically renew prior to its expiration date.\n" +
				"Default Value: true",
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
							"Default Value: If you do not designate a domain as the common name, then our system will assign it to one of your domains.",
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
	}
}
