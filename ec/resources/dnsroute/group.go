// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"terraform-provider-edgecast/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Master Server Group
func ResourceGroup() *schema.Resource {

	return &schema.Resource{
		CreateContext: ResourceGroupCreate,
		ReadContext:   ResourceGroupRead,
		UpdateContext: ResourceGroupUpdate,
		DeleteContext: ResourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration.",
			},
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fixed_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_product_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_product_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fixed_zone_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"a": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fixed_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"check_interval": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"check_type_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"content_verification": {
										Type:     schema.TypeString,
										Required: true,
									},
									"email_notification_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"failed_check_threshold": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"http_method_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fixed_record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ip_version": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"port_number": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"reintegration_method_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"status_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"record": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"rdata": {
										Type:     schema.TypeString,
										Required: true,
									},
									"verify_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_group_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_record_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"zone_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_zone_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"record_type_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"record_type_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"is_delete": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
							Required: true,
						},
					},
				},
			},
			"aaaa": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fixed_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"check_interval": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"check_type_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"content_verification": {
										Type:     schema.TypeString,
										Required: true,
									},
									"email_notification_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"failed_check_threshold": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"http_method_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fixed_record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ip_version": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"port_number": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"reintegration_method_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"status_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"record": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"rdata": {
										Type:     schema.TypeString,
										Required: true,
									},
									"verify_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_group_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_record_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"zone_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_zone_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"record_type_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"record_type_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"is_delete": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
							Required: true,
						},
					},
				},
			},
			"cname": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fixed_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"check_interval": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"check_type_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"content_verification": {
										Type:     schema.TypeString,
										Required: true,
									},
									"email_notification_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"failed_check_threshold": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"http_method_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fixed_record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ip_version": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"port_number": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"reintegration_method_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"status_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"record": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"record_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"rdata": {
										Type:     schema.TypeString,
										Required: true,
									},
									"verify_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_group_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_record_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"zone_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fixed_zone_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"record_type_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"record_type_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"is_delete": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
							Required: true,
						},
					},
				},
			},
		},
	}
}

func ResourceGroupCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Construct Group Object
	name := d.Get("name").(string)

	groupProductType := d.Get("group_product_type").(string)
	groupProductTypeID := routedns.NoGroup
	if groupProductType == "failover" {
		groupProductTypeID = routedns.Failover
	} else if groupProductType == "loadbalancing" {
		groupProductTypeID = routedns.LoadBalancing
	} else {
		d.SetId("")
		return diag.FromErr(
			fmt.Errorf(
				`invalid group_product_type: %s. It should be failover or 
				loadbalancing.`,
				groupProductType,
			),
		)
	}

	dnsAs := d.Get("a").([]interface{})
	arrayAs, err := toGroupRecords(&dnsAs)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	dnsAAAAs := d.Get("aaaa").([]interface{})
	arrayAAAAs, err := toGroupRecords(&dnsAAAAs)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	dnsCnames := d.Get("cname").([]interface{})
	arrayCnames, err := toGroupRecords(&dnsCnames)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	groupType := d.Get("group_type").(string)
	groupTypeID := routedns.PrimaryZone
	if strings.ToLower(groupType) == "cname" {
		groupTypeID = routedns.CName
	} else if strings.ToLower(groupType) == "subdomain" {
		groupTypeID = routedns.SubDomain
	}

	groupComposition := routedns.DNSGroupRecords{
		A:     *arrayAs,
		AAAA:  *arrayAAAAs,
		CNAME: *arrayCnames,
	}

	group := routedns.DnsRouteGroup{
		Name:             name,
		GroupTypeID:      groupTypeID,
		GroupProductType: groupProductTypeID,
		GroupComposition: groupComposition,
	}
	group.Name = name
	group.GroupTypeID = groupTypeID
	group.GroupProductType = groupProductTypeID
	group.GroupComposition = groupComposition

	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Call Add Group API
	params := routedns.NewAddGroupParams()
	params.AccountNumber = accountNumber
	params.Group = group

	groupID, err := routeDNSService.AddGroup(*params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*groupID))

	return ResourceGroupRead(ctx, d, m)
}

func ResourceGroupRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Construct Group Get Object
	groupID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)
	rawGroupProductType := d.Get("group_product_type").(string)
	groupProductType := routedns.NoGroup
	if strings.ToLower(rawGroupProductType) == "failover" {
		groupProductType = routedns.Failover
	} else if strings.ToLower(rawGroupProductType) == "loadbalancing" {
		groupProductType = routedns.LoadBalancing
	}
	config := m.(**api.ClientConfig)

	// Initialize Route DNS Service
	routeDNSService, err := buildRouteDNSService(**config)

	// Call Get Group API
	params := routedns.NewGetGroupParams()
	params.AccountNumber = accountNumber
	params.GroupID = groupID
	params.GroupProductType = groupProductType
	resp, err := routeDNSService.GetGroup(*params)

	if err != nil {
		return diag.FromErr(err)
	}

	// Update Terraform state with retrieved Group data
	d.Set("group_id", groupID)
	d.Set("fixed_group_id", resp.FixedGroupID)
	d.Set("fixed_zone_id", resp.FixedZoneID)
	d.Set("group_product_type_id", resp.GroupProductType)
	if resp.GroupProductType == routedns.Failover {
		d.Set("group_product_type", "failover")
	} else if resp.GroupProductType == routedns.LoadBalancing {
		d.Set("group_product_type", "loadbalancing")
	}
	d.Set("group_type_id", resp.GroupTypeID)
	if resp.GroupTypeID == routedns.CName {
		d.Set("group_type", "cname")
	} else if resp.GroupTypeID == routedns.SubDomain {
		d.Set("group_type", "subdomain")
	}

	d.Set("name", resp.Name)
	d.Set("zone_id", resp.ZoneID)

	d.Set("a", flattenGroupDNSs(&resp.GroupComposition.A))
	d.Set("aaaa", flattenGroupDNSs(&resp.GroupComposition.AAAA))
	d.Set("cname", flattenGroupDNSs(&resp.GroupComposition.CNAME))

	return diag.Diagnostics{}
}

func ResourceGroupUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	if err != nil {
		return diag.FromErr(err)
	}

	// Construct Group Update Data
	name := d.Get("name").(string)

	groupID := d.Get("group_id").(int)
	if groupID == 0 {
		groupID = -1
	}
	rawGroupType := d.Get("group_type").(string)
	groupType := routedns.PrimaryZone
	if strings.ToLower(rawGroupType) == "cname" {
		groupType = routedns.CName
	} else if strings.ToLower(rawGroupType) == "subdomain" {
		groupType = routedns.SubDomain
	}
	rawGroupProductType := d.Get("group_product_type").(string)
	groupProductType := routedns.NoGroup
	if rawGroupProductType == "failover" {
		groupProductType = routedns.Failover
	} else if rawGroupProductType == "loadbalancing" {
		groupProductType = routedns.LoadBalancing
	} else {
		return diag.FromErr(
			fmt.Errorf(
				`invalid group_product_type: %s. It should be failover or 
				loadbalancing.`,
				rawGroupProductType,
			),
		)
	}

	dnsAs := d.Get("a").([]interface{})
	arrayAs, err := toGroupRecords(&dnsAs)
	if err != nil {
		return diag.FromErr(err)
	}
	dnsAAAAs := d.Get("aaaa").([]interface{})
	arrayAAAAs, err := toGroupRecords(&dnsAAAAs)
	if err != nil {
		return diag.FromErr(err)
	}
	dnsCnames := d.Get("cname").([]interface{})
	arrayCnames, err := toGroupRecords(&dnsCnames)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Get Group API
	getParams := routedns.NewGetGroupParams()
	getParams.AccountNumber = accountNumber
	getParams.GroupID = groupID
	getParams.GroupProductType = groupProductType
	groupObj, err := routeDNSService.GetGroup(*getParams)

	if err != nil {
		return diag.FromErr(err)
	}

	// Update retrieved Group object
	groupObj.Name = name
	groupObj.GroupTypeID = groupType
	groupObj.GroupProductType = groupProductType
	groupObj.GroupComposition.A = *arrayAs
	groupObj.GroupComposition.AAAA = *arrayAAAAs
	groupObj.GroupComposition.CNAME = *arrayCnames

	// Call Update Group API
	updateParams := routedns.NewUpdateGroupParams()
	updateParams.AccountNumber = accountNumber
	updateParams.Group = groupObj
	err = routeDNSService.UpdateGroup(updateParams)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(groupObj.GroupID)) // Group ID changes on update

	return ResourceGroupRead(ctx, d, m)
}

func ResourceGroupDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Obtain group product type
	rawGroupProductType := d.Get("group_product_type").(string)
	groupProductType := routedns.NoGroup
	if rawGroupProductType == "failover" {
		groupProductType = routedns.Failover
	} else if rawGroupProductType == "loadbalancing" {
		groupProductType = routedns.LoadBalancing
	} else {
		return diag.FromErr(
			fmt.Errorf(
				`invalid group_product_type: %s. It should be failover or 
				loadbalancing.`,
				rawGroupProductType,
			),
		)
	}

	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	groupID := d.Get("group_id").(int)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Call Get Group API
	getParams := routedns.NewGetGroupParams()
	getParams.AccountNumber = accountNumber
	getParams.GroupID = groupID
	getParams.GroupProductType = groupProductType
	groupObj, err := routeDNSService.GetGroup(*getParams)

	if err != nil {
		return diag.FromErr(err)
	}

	// Call Delete Group API
	deleteParams := routedns.NewDeleteGroupParams()
	deleteParams.AccountNumber = accountNumber
	deleteParams.Group = *groupObj
	err = routeDNSService.DeleteGroup(*deleteParams)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}
