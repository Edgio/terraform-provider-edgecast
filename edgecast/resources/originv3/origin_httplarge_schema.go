// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package originv3

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetOriginGrpHttpLargeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Indicates the customer origin group's system-defined ID.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Indicates the customer origin group's name.",
		},
		"host_header": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Indicates the value that will be assigned to the Host header for all requests to this customer origin configuration.",
		},
		"network_type_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
			Description: "Indicates the method for resolving hostnames through its system-defined ID. \n" +
				"Valid values are:  \n" +
				" 1: Default. Indicates that the IP preference for this customer origin will be system-defined. Currently, this configuration causes hostnames to be resolved to IPv4 only. \n" +
				" 2: IPv6 Preferred over IPv4. Indicates that although hostnames for this customer origin can be resolved to either IP version, a preference will be given to IPv6. \n" +
				" 3: IPv4 Preferred over IPv6. Indicates that although hostnames for this customer origin can be resolved to either IP version, a preference will be given to IPv4. \n" +
				" 4: IPv4 Only. Indicates that hostnames for this customer origin can only be resolved to IPv4. \n" +
				" 5: IPv6 Only. Indicates that hostnames for this customer origin can only be resolved to IPv6. \n" +
				" \n" +
				"Default Value: \n" +
				"1",
		},
		"shield_pops": {
			Type:     schema.TypeList,
			Optional: true,
			Description: "Contains this customer origin group's Origin Shield configuration. \n" +
				"Key information: \n" +
				"A customer origin group is protected by Origin Shield when one or more string value(s) have been defined. \n" +
				" \n" +
				"Each string value represents an Origin Shield POP or bypass region code that defines how Origin Shield will behave for all regions or a specific region. \n" +
				" \n" +
				"The number of string values contained within this array determines this customer origin group's Origin Shield configuration: \n" +
				"Single Value: If this array contains a single string value, then that string value defines the Origin Shield configuration for all regions. This is known as the Single POP option within the MCC. \n" +
				"Multiple Values: If this array contains multiple string values, then each string value defines an Origin Shield configuration for the region associated with the specified Origin Shield POP or bypass region code. This is known as the Multiple POPs option within the MCC. \n" +
				"Bypass Origin Shield within a specific region through bypass region codes. \n" +
				"Define a region's Origin Shield location by passing the code corresponding to the desired Origin Shield POP code. Our service will automatically assign it to a region based off of the POP's location. \n" +
				"You may not define duplicate entries or multiple configurations for the same region. \n" +
				"If a region has not been assigned an Origin Shield configuration, then requests for that region will be forwarded to the next closest Origin Shield region. This configuration is the equivalent of leaving a region blank from within the MCC. \n" +
				" \n" +
				"If the strict_pci_certified option has been enabled on this customer origin group, then you may only specify a Payment Card Industry (PCI)-compliant POP. \n" +
				" \n" +
				"Use the Get Origin Shield POPs endpoint to retrieve a list of regions, the Origin Shield POPs associated with those regions, and each POP's PCI-compliance status (is_pci_certified).",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"strict_pci_certified": {
			Type:     schema.TypeBool,
			Optional: true,
			Description: "Indicates whether this customer origin group is restricted to Payment Card Industry (PCI)-compliant Origin Shield POPs.  \n" +
				"Valid values are: true | false",
		},
		"tls_settings": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Contains settings that define TLS behavior.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"allow_self_signed": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
						Description: "Indicates whether our CDN will allow delivery when an edge server detects a self-signed certificate from the origin server during the TLS handshake. \n" +
							"The corresponding MCC setting is Disallow Self-Signed. If you set the allow_self_signed property to True, then the Disallow Self-Signed option will be set to False.  \n" +
							"Default Value: False",
					},
					"public_keys_to_verify": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Contains a list of SHA-1 digests for the public key of your end-entity (i.e., leaf) certificate.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"sni_hostname": {
						Type:     schema.TypeString,
						Optional: true,
						Description: "Indicates the hostname that will be sent as a SNI hint during the TLS handshake for this customer origin group. \n" +
							"Our service will not use Server Name Indication (SNI) during the TLS handshake for this customer origin group when this property is set to null or an empty value.",
					},
				},
			},
		},
	}
}
