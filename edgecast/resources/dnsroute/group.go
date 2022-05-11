// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-edgecast/edgecast/helper"

	"terraform-provider-edgecast/edgecast/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Master Server Group
func ResourceGroup() *schema.Resource {
	groupRecord := map[string]*schema.Schema{
		"health_check": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: `Define a record's health check configuration`,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Identifies the health check by its 
						system-defined ID.`,
					},
					"fixed_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Reserved for future use.",
					},
					"check_interval": {
						Type:     schema.TypeInt,
						Required: true,
						Description: `Defines the number of seconds between 
						health checks.`,
					},
					"check_type_id": {
						Type:     schema.TypeInt,
						Required: true,
						Description: `Defines the type of health check by its 
						system-defined ID. The following values are supported: 
						1 - HTTP | 2 - HTTPS | 3 - TCP Open | 4 - TCP SSL. 
						Please refer to the following URL for additional 
						information:
						https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HC_Types.htm`,
					},
					"content_verification": {
						Type:     schema.TypeString,
						Required: true,
						Description: `Defines the text that will be used to 
						verify the success of the health check.`,
					},
					"email_notification_address": {
						Type:     schema.TypeString,
						Required: true,
						Description: `Defines the e-mail address to which 
						health check notifications will be sent.`,
					},
					"failed_check_threshold": {
						Type:     schema.TypeInt,
						Required: true,
						Description: `Defines the number of consecutive 
						times that the same result must be returned before 
						a health check agent will indicate a change in status.`,
					},
					"http_method_id": {
						Type:     schema.TypeInt,
						Optional: true,
						Description: `Defines an HTTP method by its 
						system-defined ID. An HTTP method is only used by 
						HTTP/HTTPs health checks. Supported values are: 
						1 - GET, 2 - POST. Refer to the following URL for 
						additional information:
						https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HTTP_Methods.htm`,
					},
					"record_id": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Defines the DNS record ID this health 
						check is associated with.`,
					},
					"fixed_record_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Reserved for future use.",
					},
					"group_id": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Defines the Group ID this health check 
						is associated with.`,
					},
					"ip_address": {
						Type:     schema.TypeString,
						Optional: true,
						Description: `Defines the IP address (IPv4 or IPv6) to 
						which TCP health checks will be directed. IP address is 
						required when check_type_id is 3 or 4`,
					},
					"ip_version": {
						Type:     schema.TypeInt,
						Optional: true,
						Description: `Defines an IP version by its 
						system-defined ID. This IP version is only used by 
						HTTP/HTTPs health checks. Supported values are: 
						1 - IPv4, 2 - IPv6. Refer to the following URL for 
						additional information:
						https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_IP_Versions_HC.htm`,
					},
					"port_number": {
						Type:     schema.TypeInt,
						Optional: true,
						Description: `Defines the port to which TCP health 
						checks will be directed.`,
					},
					"reintegration_method_id": {
						Type:     schema.TypeInt,
						Required: true,
						Description: `Indicates the method through which an 
						unhealthy server/hostname will be integrated back into a 
						group. Supported values are: 1 - Automatic | 2 - Manual`,
					},
					"status": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Indicates the server/hostname's health 
						check status by its system-defined ID.`,
					},
					"status_name": {
						Type:     schema.TypeString,
						Computed: true,
						Description: `Indicates the server/hostname's health 
						check status.`,
					},
					"timeout": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: `Reserved for future use.`,
					},
					"uri": {
						Type:     schema.TypeString,
						Optional: true,
						Description: `Defines the URI to which HTTP/HTTPs health 
						checks will be directed.`,
					},
				},
			},
		},
		"weight": {
			Type:     schema.TypeInt,
			Required: true,
			Description: `Defines a record's weight. Used to denote preference 
			for a load balancing or failover group.`,
		},
		"record": {
			Type: schema.TypeList,
			Description: `Defines a DNS record that will be associated with the 
			zone.`,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"record_id": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Identifies a DNS Record by its 
						system-defined ID.`,
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Defines a record's name. `,
					},
					"ttl": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: `Defines a record's TTL.`,
					},
					"rdata": {
						Type:        schema.TypeString,
						Required:    true,
						Description: `Defines a record's value.`,
					},
					"verify_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: `Reserved for future use.`,
					},
					"fixed_group_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: `Reserved for future use.`,
					},
					"group_id": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Identifies the group this record is 
						assoicated with by its system-defined ID.`,
					},
					"fixed_record_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: `Reserved for future use.`,
					},
					"zone_id": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: `Reserved for future use.`,
					},
					"fixed_zone_id": {
						Type:     schema.TypeInt,
						Optional: true,
						Description: `Identifies a zone by its system-defined 
						ID.`,
					},
					"record_type_id": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Indicates the system-defined ID assigned 
						to the record type.`,
					},
					"record_type_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Indicates the name of the record type.`,
					},
					"is_delete": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: `Reserved for future use.`,
					},
					"weight": {
						Type:     schema.TypeInt,
						Computed: true,
						Description: `Defines a record's weight. Used to denote 
						preference for a load balancing or failover group.`,
					},
				},
			},
			Required: true,
		},
	}

	return &schema.Resource{
		CreateContext: ResourceGroupCreate,
		ReadContext:   ResourceGroupRead,
		UpdateContext: ResourceGroupUpdate,
		DeleteContext: ResourceGroupDelete,
		Importer:      helper.Import(ResourceGroupRead, "account_number", "id", "group_product_type"),

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Account Number associated with the customer whose 
				resources you wish to manage. This account number may be found 
				in the upper right-hand corner of the MCC.`,
			},
			"group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Identifies the group by its system-defined ID.`,
			},
			"fixed_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Reserved for future use.`,
			},
			"group_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Defines the group type. Valid values are: cname | 
				subdomain`,
			},
			"group_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Defines the group type by its system-defined ID`,
			},
			"group_product_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Defines the group product type. Valid values are:
				loadbalancing | failover`,
			},
			"group_product_type_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `Defines the group product type by its 
				system-defined ID`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Defines the name of the failover or load balancing 
				group.`,
			},
			"zone_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Reserved for future use.`,
			},
			"fixed_zone_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Reserved for future use.`,
			},
			"a": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: `Defines a set of A records associated with this 
				group.`,
				Elem: &schema.Resource{
					Schema: groupRecord,
				},
			},
			"aaaa": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: `Defines a set of AAAA records associated with this 
				group.`,
				Elem: &schema.Resource{
					Schema: groupRecord,
				},
			},
			"cname": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: `Defines a set of CNAME records associated with 
				this group.`,
				Elem: &schema.Resource{
					Schema: groupRecord,
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
				loadbalancing`,
				groupProductType,
			),
		)
	}

	dnsAs := d.Get("a").(*schema.Set).List()
	arrayAs, err := expandGroupRecords(&dnsAs, false)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	dnsAAAAs := d.Get("aaaa").(*schema.Set).List()
	arrayAAAAs, err := expandGroupRecords(&dnsAAAAs, false)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	dnsCnames := d.Get("cname").(*schema.Set).List()
	arrayCnames, err := expandGroupRecords(&dnsCnames, false)
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
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

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
	if err != nil {
		return diag.FromErr(err)
	}

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
	if err != nil {
		return diag.FromErr(err)
	}

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

	if err := d.Set(
		"a", flattenGroupDNSs(&resp.GroupComposition.A),
	); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(
		"aaaa", flattenGroupDNSs(&resp.GroupComposition.AAAA),
	); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(
		"cname", flattenGroupDNSs(&resp.GroupComposition.CNAME),
	); err != nil {
		return diag.FromErr(err)
	}

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
				loadbalancing`,
				rawGroupProductType,
			),
		)
	}

	dnsAs := d.Get("a").(*schema.Set).List()

	arrayAs, err := expandGroupRecords(&dnsAs, false)
	if err != nil {
		return diag.FromErr(err)
	}
	dnsAAAAs := d.Get("aaaa").(*schema.Set).List()
	arrayAAAAs, err := expandGroupRecords(&dnsAAAAs, false)
	if err != nil {
		return diag.FromErr(err)
	}
	dnsCnames := d.Get("cname").(*schema.Set).List()
	arrayCnames, err := expandGroupRecords(&dnsCnames, false)
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
				loadbalancing`,
				rawGroupProductType,
			),
		)
	}

	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	groupID := d.Get("group_id").(int)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
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
