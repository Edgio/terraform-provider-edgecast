// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package origin

import (
	"context"
	"errors"
	"log"
	"strconv"

	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/origin"
	"github.com/EdgeCast/ec-sdk-go/edgecast/shared/enums"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOrigin() *schema.Resource {
	hostnameSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Identifies the hostname/IP 
				address of an origin server that will be 
				associated with the customer origin 
				configuration being created.`,
			},
			"is_primary": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Indicates whether a 
				particular hostname/IP address is the 
				primary one for HTTP requests.
				Valid values are:
				0: Is not the primary one.
				1: Is the primary one.`,
			},
			"ordinal": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Indicates the position in 
				the ordered list for the current 
				hostname/IP address. This position is 
				primarily used by "Primary and Failover" 
				load balancing mode to determine which 
				hostname/IP address will take over when 
				a hostname/IP address higher on the list 
				is unreachable.`,
			},
		},
	}

	return &schema.Resource{
		CreateContext: ResourceOriginCreate,
		ReadContext:   ResourceOriginRead,
		UpdateContext: ResourceOriginUpdate,
		DeleteContext: ResourceOriginDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Account Number associated with the customer whose 
				origins you wish to manage. This account number may be found in 
				the upper right-hand corner of the MCC.`,
			},
			"directory_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Identifies the directory name that will be 
				assigned to the customer origin configuration. This alphanumeric 
				value is appended to the end of the base CDN URL that points to 
				the customer origin server. Note: A protocol should not be 
				specified when setting this parameter. Examples: 
				"www.example.com:80","10.10.10.255:80","[1:2:3:4:5:6:7:8]:80"`,
			},
			"host_header": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Defines the value that will be assigned to the 
				Host header for all requests to this customer origin 
				configuration. A host header is especially useful when there are 
				multiple virtual hostnames hosted on a single physical server or 
				load-balanced set of servers.`,
			},
			"media_type_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Identifies the platform on which the customer 
				origin configuration resides. Valid values are:
				3: HTTP Large, 8: HTTP Small, 14: ADN`,
			},
			"network_configuration": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Determines how hostnames associated with this 
				customer origin configuration will be resolved to an IP 
				address.`,
			},
			"load_balancing_scheme_http": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Determines how HTTP requests will be 
				load balanced across the specified hostnames/IP 
				addresses. Valid values: "RR" for Round Robin and 
				"PF" for Primary and Failover.`,
			},
			"origin_hostname_http": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     hostnameSchema,
				Description: `Identifies the HTTP origin servers, load 
				balancing configuration, and origin precedence if applicable to 
				the load balancing type specified.`,
			},
			"load_balancing_scheme_https": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Determines how HTTPS requests will be 
				load balanced across the specified hostnames/IP 
				addresses. Valid values: "RR" for Round Robin and 
				"PF" for Primary and Failover.`,
			},
			"origin_hostname_https": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     hostnameSchema,
				Description: `Identifies the HTTPS origin servers, load 
				balancing configuration, and origin precedence if applicable to 
				the load balancing type specified.`,
			},
			"shield_pop": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pop_code": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `Defines an Origin Shield configuration 
							for this customer origin. This configuration is 
							defined through a three or four-letter code.`,
						},
					},
				},
			},
			"follow_redirects": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Indicates whether our edge servers will respect a 
				URL redirect when validating the set of optimal ADN gateway 
				servers for your customer origin configuration.`,
			},
			"validation_url": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Indicates the URL to a sample asset. A set of 
				optimal ADN gateway servers for your customer origin server is 
				determined through the delivery of this sample asset.`,
			},
			"http_full_url": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the CDN URL for HTTP requests to this 
				customer origin server.`,
			},
			"https_full_url": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the CDN URL for HTTPS requests to this 
				customer origin server.`,
			},
			"use_origin_shield": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `Indicates whether Origin Shield has been activated 
				on the customer origin. Valid values are:
				0: Disabled, 1: Enabled`,
			},
		},
	}
}

func ResourceOriginCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	mediaTypeID := d.Get("media_type_id").(int)

	// Build Origin object
	originObj := &origin.Origin{
		DirectoryName:        d.Get("directory_name").(string),
		HostHeader:           d.Get("host_header").(string),
		NetworkConfiguration: d.Get("network_configuration").(int),
		HTTPLoadBalancing:    d.Get("load_balancing_scheme_http").(string),
		HTTPSLoadBalancing:   d.Get("load_balancing_scheme_https").(string),
		FollowRedirects:      d.Get("follow_redirects").(bool),
		ValidationURL:        d.Get("validation_url").(string),
	}

	if attr, ok := d.GetOk("origin_hostname_http"); ok {
		hostnamesHTTP, err := expandHostname(attr)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		originObj.HTTPHostnames = *hostnamesHTTP
	}

	if attr, ok := d.GetOk("origin_hostname_https"); ok {
		hostnamesHTTPS, err := expandHostname(attr)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		originObj.HTTPSHostnames = *hostnamesHTTPS
	}

	if attr, ok := d.GetOk("shield_pop"); ok {
		shieldPOPs, err := expandShieldPOPs(attr)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		originObj.ShieldPOPs = *shieldPOPs
	}

	log.Printf("resourceOriginCreate>>origin object:%v\n", originObj)

	// Initialize Origin Service
	originService, err := buildOriginService(**config)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// Call Create Origin API
	params := origin.NewAddOriginParams()
	params.AccountNumber = accountNumber
	params.MediaTypeID = enums.Platform(mediaTypeID)
	params.Origin = *originObj
	originID, err := originService.AddOrigin(*params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*originID))

	return ResourceOriginRead(ctx, d, m)
}

func ResourceOriginRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	mediaTypeID := d.Get("media_type_id").(int)
	customerOriginID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Initialize Origin Service
	originService, err := buildOriginService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Origin from API
	params := origin.NewGetOriginParams()
	params.AccountNumber = accountNumber
	params.CustomerOriginID = customerOriginID
	params.MediaTypeID = enums.Platform(mediaTypeID)
	parsedResponse, err := originService.GetOrigin(*params)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("directory_name", parsedResponse.DirectoryName)
	d.Set("host_header", parsedResponse.HostHeader)
	d.Set("network_configuration", parsedResponse.NetworkConfiguration)
	d.Set("load_balancing_scheme_http", parsedResponse.HTTPLoadBalancing)
	d.Set("load_balancing_scheme_https", parsedResponse.HTTPSLoadBalancing)
	d.Set("follow_redirects", parsedResponse.FollowRedirects)
	d.Set("validation_url", parsedResponse.ValidationURL)
	d.Set("http_full_url", parsedResponse.HTTPFullURL)
	d.Set("https_full_url", parsedResponse.HTTPSFullURL)
	d.Set("use_origin_shield", parsedResponse.UseOriginShield)

	httpHostnames := flattenHostname(parsedResponse.HTTPHostnames)
	d.Set("origin_hostname_http", httpHostnames)

	httpsHostnames := flattenHostname(parsedResponse.HTTPSHostnames)
	d.Set("origin_hostname_https", httpsHostnames)

	shieldPOPs := flattenShieldPOP(parsedResponse.ShieldPOPs)
	d.Set("shield_pop", shieldPOPs)

	return diag.Diagnostics{}
}

func ResourceOriginUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	mediaTypeID := d.Get("media_type_id").(int)
	customerOriginID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Initialize Origin Service
	originService, err := buildOriginService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve existing Origin object from API
	getParams := origin.NewGetOriginParams()
	getParams.AccountNumber = accountNumber
	getParams.CustomerOriginID = customerOriginID
	getParams.MediaTypeID = enums.Platform(mediaTypeID)

	originObj, err := originService.GetOrigin(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Prepare updated Origin object
	originObj.DirectoryName = d.Get("directory_name").(string)
	originObj.HostHeader = d.Get("host_header").(string)
	originObj.NetworkConfiguration = d.Get("network_configuration").(int)
	originObj.HTTPLoadBalancing = d.Get("load_balancing_scheme_http").(string)
	originObj.HTTPSLoadBalancing = d.Get("load_balancing_scheme_https").(string)
	originObj.FollowRedirects = d.Get("follow_redirects").(bool)
	originObj.ValidationURL = d.Get("validation_url").(string)

	if attr, ok := d.GetOk("origin_hostname_http"); ok {
		hostnamesHTTP, err := expandHostname(attr)
		if err != nil {
			return diag.FromErr(err)
		}
		originObj.HTTPHostnames = *hostnamesHTTP
	}

	if attr, ok := d.GetOk("origin_hostname_https"); ok {
		hostnamesHTTPS, err := expandHostname(attr)
		if err != nil {
			return diag.FromErr(err)
		}
		originObj.HTTPSHostnames = *hostnamesHTTPS
	}

	if attr, ok := d.GetOk("shield_pop"); ok {
		shieldPOPs, err := expandShieldPOPs(attr)
		if err != nil {
			return diag.FromErr(err)
		}
		originObj.ShieldPOPs = *shieldPOPs
	}

	log.Printf("resourceOriginUpdate>>origin object:%v\n", originObj)

	// Call Create Origin API
	updateParams := origin.NewUpdateOriginParams()
	updateParams.AccountNumber = accountNumber
	updateParams.Origin = *originObj

	_, err = originService.UpdateOrigin(*updateParams)

	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceOriginRead(ctx, d, m)
}

func ResourceOriginDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	mediaTypeID := d.Get("media_type_id").(int)
	customerOriginID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Initialize Origin Service
	originService, err := buildOriginService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve existing Origin object from API
	getParams := origin.NewGetOriginParams()
	getParams.AccountNumber = accountNumber
	getParams.CustomerOriginID = customerOriginID
	getParams.MediaTypeID = enums.Platform(mediaTypeID)

	originObj, err := originService.GetOrigin(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Delete Origin API
	deleteParams := origin.NewDeleteOriginParams()
	deleteParams.AccountNumber = accountNumber
	deleteParams.Origin = *originObj

	err = originService.DeleteOrigin(*deleteParams)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}

func expandHostname(attr interface{}) (*[]origin.Hostname, error) {
	if set, ok := attr.(*schema.Set); ok {

		items := set.List()
		hostgroups := make([]origin.Hostname, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			hostgroup := origin.Hostname{
				Name:      curr["name"].(string),
				IsPrimary: curr["is_primary"].(int),
				Ordinal:   curr["ordinal"].(int),
			}

			hostgroups = append(hostgroups, hostgroup)
		}

		return &hostgroups, nil

	} else {
		return nil, errors.New("attr input was not a *schema.Set")
	}
}

func expandShieldPOPs(attr interface{}) (*[]origin.ShieldPOP, error) {
	if set, ok := attr.(*schema.Set); ok {

		items := set.List()
		shieldPOPs := make([]origin.ShieldPOP, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			shieldPOP := origin.ShieldPOP{
				POPCode: curr["pop_code"].(string),
			}

			shieldPOPs = append(shieldPOPs, shieldPOP)
		}

		return &shieldPOPs, nil

	} else {
		return nil, errors.New("attr input was not a *schema.Set")
	}
}

func flattenHostname(hostgroups []origin.Hostname) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, hostgroup := range hostgroups {
		m := make(map[string]interface{})
		m["name"] = hostgroup.Name
		m["is_primary"] = hostgroup.IsPrimary
		m["ordinal"] = hostgroup.Ordinal
		flattened = append(flattened, m)
	}
	return flattened
}

func flattenShieldPOP(shieldPOPs []origin.ShieldPOP) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, shieldPOP := range shieldPOPs {
		m := make(map[string]interface{})
		m["pop_code"] = shieldPOP.POPCode
		flattened = append(flattened, m)
	}
	return flattened
}
