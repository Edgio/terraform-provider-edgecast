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
				" **Valid values are:**  \n" +
				" * **1:** Default. Indicates that the IP preference for this customer origin will be system-defined. Currently, this configuration causes hostnames to be resolved to IPv4 only. \n" +
				" * **2:** IPv6 Preferred over IPv4. Indicates that although hostnames for this customer origin can be resolved to either IP version, a preference will be given to IPv6. \n" +
				" * **3:** IPv4 Preferred over IPv6. Indicates that although hostnames for this customer origin can be resolved to either IP version, a preference will be given to IPv4. \n" +
				" * **4:** IPv4 Only. Indicates that hostnames for this customer origin can only be resolved to IPv4. \n" +
				" * **5:** IPv6 Only. Indicates that hostnames for this customer origin can only be resolved to IPv6. \n" +
				" **Default Value:** \n" +
				" * **1**",
		},
		"shield_pops": {
			Type:     schema.TypeList,
			Optional: true,
			Description: "Contains this customer origin group's Origin Shield configuration. \n" +
				" **Key information:** \n" +
				" * A customer origin group is protected by Origin Shield when one or more string value(s) have been defined. \n" +
				" \n" +
				" * Each string value represents an Origin Shield POP or bypass region code that defines how Origin Shield will behave for all regions or a specific region. \n" +
				" \n" +
				" * The number of string values contained within this array determines this customer origin group's Origin Shield configuration: \n" +
				"  	  **Single Value:** If this array contains a single string value, then that string value defines the Origin Shield configuration for all regions. This is known as the Single POP option within the MCC. \n" +
				"	  **Multiple Values:** If this array contains multiple string values, then each string value defines an Origin Shield configuration for the region associated with the specified Origin Shield POP or bypass region code. This is known as the Multiple POPs option within the MCC. \n" +
				"		* Bypass Origin Shield within a specific region through bypass region codes. \n" +
				"		* Define a region's Origin Shield location by passing the code corresponding to the desired Origin Shield POP code. Our service will automatically assign it to a region based off of the POP's location. \n" +
				"		* You may not define duplicate entries or multiple configurations for the same region. \n" +
				"		* If a region has not been assigned an Origin Shield configuration, then requests for that region will be forwarded to the next closest Origin Shield region. This configuration is the equivalent of leaving a region blank from within the MCC. \n" +
				" \n" +
				" * If the strict_pci_certified option has been enabled on this customer origin group, then you may only specify a Payment Card Industry (PCI)-compliant POP. \n" +
				" \n" +
				" * Use the Get Origin Shield POPs endpoint to retrieve a list of regions, the Origin Shield POPs associated with those regions, and each POP's PCI-compliance status (is_pci_certified).",
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
							"**Default Value:** False",
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
		"origin": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Indicates the origin entry's system-defined ID.",
					},
					"host": {
						Type:     schema.TypeString,
						Required: true,
						Description: "Identifies the web server(s) that will be associated with this origin entry through either a hostname or IP address.  \n" +
							"**Key information:**  \n\n" +
							" * If you set the protocol_type_id property to either the HTTPS Only or HTTP Only mode, then you may also define a protocol for edge to origin communication. \n" +
							" * You may not specify a hostname or IP address that points to our network. This type of configuration is disallowed since it may cause your traffic to infinitely loop within our network. \n" +
							" * Use any combination of hostnames and IP addresses when defining a customer origin group's origin entries. \n" +
							" * Our service resolve hostnames to an IP address prior to delivery. The customer origin group's network_type_id property determines how hostnames will be resolved to an IP address. \n" +
							"   **Syntax:** \n" +
							"   * **HTTPS Only and HTTP Only Modes:** {Protocol}{Hostname or IP Address} \n" +
							"     If you choose to assign a protocol, make sure that all origin entries within a customer origin group for a given mode use the same protocol. For example, you may not configure some HTTPS Only origin entries to use the HTTP protocol and others to use the HTTPS protocol. \n" +
							"     Hostname Example: https://www.mydomain.com \n" +
							"     IPv6 Example: [1:2:3:4:5:6:7:8] \n" +
							"   * **Match Client Mode:** {Hostname or IP Address} \n" +
							"     IPv4 Example: 10.10.10.255 \n" +
							"	  Brackets are required when identifying an origin server through the use of IPv6 notation. This is the standard URI convention for IPv6 addresses. \n" +
							"",
					},
					"is_primary": {
						Type:     schema.TypeBool,
						Required: true,
						Description: "Determines whether this origin entry identifies the primary hostname or IP address for the protocol defined within the protocol_type_id property.  \n\n" +
							"You may only enable this property on a single origin entry within a customer origin group per protocol. For the purpose of this restriction, an origin entry that uses the Match Client mode is considered to be assigned both HTTP and HTTPS.  \n\n" +
							"For example, enabling this property on an origin entry that uses the Match Client mode will cause it to be the primary origin entry for both HTTP and HTTPS. \n\n" +
							"This property is critical for determining how requests are load balanced. Setup for both modes are described below. \n" +
							"    * **Primary / Failover:** Our service load balances requests using primary / failover mode when this property is enabled on an origin entry within the desired customer origin group. This load balancing mode is restricted to that origin entry's protocol. If that origin entry has been configured to use the Match Client mode, then requests will be load balanced for both HTTP and HTTPS. \n" +
							"    * **Round-Robin:** Our service load balances requests using round-robin mode when this property is disabled on all origin entries within the desired customer origin group for the desired protocol. \n" +
							"",
					},
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Defines the origin entry's name.",
					},
					"port": {
						Type:     schema.TypeInt,
						Optional: true,
						Description: "Determines the port for communication with your origin servers.  \n" +
							"**Default Value:** \n " +
							"This property's default value varies according to the value defined within the protocol_type_id property. \n " +
							"80 | 443",
					},
					"protocol_type_id": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  1,
						Description: "Determines this origin entry's protocol through its system-defined ID. Valid values are: \n " +
							"* **1:** HTTP Only" +
							"* **2:** HTTPS Only \n " +
							"* **3:** Match Client \n " +
							"",
					},
					"storage_type_id": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  1,
						Description: "Identifies the origin group's type through its system-defined ID. Valid values are: \n\n" +
							"* **1:** Customer origin group ",
					},
					"failover_order": {
						Type:     schema.TypeInt,
						Optional: true,
						Description: "Indicates this origin entry's sort position as a 0-based number. \n\n" +
							"**Key Information** Position 0 is reserved for the primary origin entry as determined by the is_primary field. \n\n" +
							"The primary purpose of this position is to determine the order in which requests are load balanced \n " +
							"for Primary / Failover mode. If a primary origin entry has been defined for this protocol, \n " +
							"then all traffic for that protocol will be directed to the origin entry that has the lowest value. \n " +
							"If the hostname or IP address associated with that origin entry is unreachable, then traffic will be directed \n " +
							"to the next lowest value. This process will continue until our service can establish communication with your origin.",
					},
				},
			},
			Description: "Contains the origin entry(s) associated with this origin group.",
		},
	}
}
