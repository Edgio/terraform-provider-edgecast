// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrigin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOriginCreate,
		ReadContext:   resourceOriginRead,
		UpdateContext: resourceOriginUpdate,
		DeleteContext: resourceOriginDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_header": {
				Type:     schema.TypeString,
				Required: true,
			},
			"media_type": {
				Type:     schema.TypeString,
				Required: true,
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
						},
						"hostnames": {
							Type:     schema.TypeList,
							Required: true,
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

func resourceOriginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	return resourceOriginRead(ctx, d, m)
}

func resourceOriginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceOriginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	return resourceOriginRead(ctx, d, m)
}

func resourceOriginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
