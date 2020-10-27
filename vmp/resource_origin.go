// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
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
			"account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host_header": &schema.Schema{
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
	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	httpConfiguration := d.Get("http").(*schema.Set).List()[0].(map[string]interface{})

	addOriginRequest := &api.AddOriginRequest{
		DirectoryName:     d.Get("directory_name").(string),
		HostHeader:        d.Get("host_header").(string),
		HttpLoadBalancing: httpConfiguration["load_balancing"].(string),
	}

	rawHTTPHostnames := httpConfiguration["hostnames"].([]interface{})

	httpHostnames := make([]api.AddOriginRequestHostname, len(rawHTTPHostnames))

	for i := range rawHTTPHostnames {
		httpHostnames[i] = api.AddOriginRequestHostname{Name: rawHTTPHostnames[i].(string)}
	}

	addOriginRequest.HttpHostnames = httpHostnames

	originAPIClient := api.NewOriginApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	parsedResponse, err := originAPIClient.AddHttpLargeOrigin(addOriginRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(parsedResponse.CustomerOriginId))

	return resourceOriginRead(ctx, d, m)
}

func resourceOriginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	originAPIClient := api.NewOriginApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	originID, _ := strconv.Atoi(d.Id())

	parsedResponse, err := originAPIClient.GetHttpLargeOrigin(originID)

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
	return resourceOriginRead(ctx, d, m)
}

func resourceOriginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfiguration, err := m.(*ProviderConfiguration).ApplyAccountNumberOverride(d.Get("account_number").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	originAPIClient := api.NewOriginApiClient(providerConfiguration.APIClient, providerConfiguration.AccountNumber)

	originID, _ := strconv.Atoi(d.Id())

	err = originAPIClient.DeleteOrigin(originID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
