// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package cname_origin

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"terraform-provider-ec/ec/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOrigin() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOriginCreate,
		ReadContext:   ResourceOriginRead,
		UpdateContext: ResourceOriginUpdate,
		DeleteContext: ResourceOriginDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration.",
			},
			"directory_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Identifies the directory name that will be assigned to the customer origin configuration. " +
					"This alphanumeric value is appended to the end of the base CDN URL that points to the customer origin " +
					"server. Note: A protocol should not be specified when setting this parameter. " +
					"Examples: `www.example.com:80`,`10.10.10.255:80`,`[1:2:3:4:5:6:7:8]:80`",
			},
			"host_header": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Defines the value that will be assigned to the Host header for all requests to this " +
					"customer origin configuration. A host header is especially useful when there are multiple virtual " +
					"hostnames hosted on a single physical server or load-balanced set of servers.",
			},
			"media_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the Delivery Platform to use. Valid values are `httplarge`,`httpsmall`, and `adn`",
			},
			"http": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancing": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Determines how HTTP requests will be load balanced across the specified " +
								"hostnames/IP addresses. Valid values: `RR` for Round Robin and `PF` for Primary and Failover.",
						},
						"hostnames": {
							Type:     schema.TypeList,
							Required: true,
							Description: "This request parameter contains the set of hostnames/IP addresses to which " +
								"HTTP requests for this customer origin configuration may be fulfilled. ",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func ResourceOriginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = d.Get("account_number").(string)

	httpConfiguration := d.Get("http").(*schema.Set).List()[0].(map[string]interface{})
	log.Printf("resourceOriginCreate>>httpConfiguration: %v", httpConfiguration)
	addOriginRequest := &api.AddOriginRequest{
		DirectoryName:     d.Get("directory_name").(string),
		HostHeader:        d.Get("host_header").(string),
		HTTPLoadBalancing: httpConfiguration["load_balancing"].(string),
	}

	rawHTTPHostnames := httpConfiguration["hostnames"].([]interface{})

	httpHostnames := make([]api.AddOriginRequestHostname, len(rawHTTPHostnames))

	for i := range rawHTTPHostnames {
		httpHostnames[i] = api.AddOriginRequestHostname{Name: rawHTTPHostnames[i].(string)}
	}

	addOriginRequest.HTTPHostnames = httpHostnames
	log.Printf("resourceOriginCreate>>addOriginRequest: %v", addOriginRequest)
	originAPIClient := api.NewOriginAPIClient(*config)

	mediaType := d.Get("media_type").(string)

	parsedResponse, err := originAPIClient.AddOrigin(addOriginRequest, mediaType)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(parsedResponse.CustomerOriginID))

	return ResourceOriginRead(ctx, d, m)
}

func ResourceOriginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = d.Get("account_number").(string)

	originAPIClient := api.NewOriginAPIClient(*config)

	originID, _ := strconv.Atoi(d.Id())
	mediaType := d.Get("media_type").(string)

	parsedResponse, err := originAPIClient.GetOrigin(originID, mediaType)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("directory_name", parsedResponse.DirectoryName)
	d.Set("host_header", parsedResponse.HostHeader)
	d.Set("http_load_balancing", parsedResponse.HttpLoadBalancing)
	d.Set("http_hostnames", parsedResponse.HttpHostnames)

	return diags
}

func ResourceOriginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = d.Get("account_number").(string)
	httpConfiguration := d.Get("http").(*schema.Set).List()[0].(map[string]interface{})
	fmt.Printf("Origin>>Update[load_balancing]:%s", httpConfiguration["load_balancing"])
	updateOriginRequest := &api.UpdateOriginRequest{
		DirectoryName:     d.Get("directory_name").(string),
		HostHeader:        d.Get("host_header").(string),
		HTTPLoadBalancing: httpConfiguration["load_balancing"].(string),
	}

	rawHTTPHostnames := httpConfiguration["hostnames"].([]interface{})

	httpUpdateHostnames := make([]api.UpdateOriginRequestHostname, len(rawHTTPHostnames))

	for i := range rawHTTPHostnames {
		httpUpdateHostnames[i] = api.UpdateOriginRequestHostname{Name: rawHTTPHostnames[i].(string)}
	}

	updateOriginRequest.HTTPHostnames = httpUpdateHostnames

	originAPIClient := api.NewOriginAPIClient(*config)

	mediaType := d.Get("media_type").(string)
	originID, _ := strconv.Atoi(d.Id())
	parsedResponse, err := originAPIClient.UpdateOrigin(updateOriginRequest, originID, mediaType)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(parsedResponse.CustomerOriginID))
	return ResourceOriginRead(ctx, d, m)
}

func ResourceOriginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = d.Get("account_number").(string)
	originAPIClient := api.NewOriginAPIClient(*config)

	originID, _ := strconv.Atoi(d.Id())

	err := originAPIClient.DeleteOrigin(originID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
