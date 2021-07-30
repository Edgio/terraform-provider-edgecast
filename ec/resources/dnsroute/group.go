// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"terraform-provider-ec/ec/api"

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

func ResourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	name := d.Get("name").(string)

	groupID := d.Get("group_id").(int)
	if groupID == 0 {
		groupID = -1
	}
	fixedGroupID := d.Get("fixed_group_id").(int)
	if fixedGroupID == 0 {
		fixedGroupID = -1
	}
	zoneID := d.Get("zone_id").(int)
	if zoneID == 0 {
		zoneID = -1
	}
	fixedZoneID := d.Get("fixed_zone_id").(int)
	if fixedZoneID == 0 {
		fixedZoneID = -1
	}
	groupProductType := d.Get("group_product_type").(string)
	groupProductTypeID := api.GroupProductType_NoGroup
	if groupProductType == "failover" {
		groupProductTypeID = api.GroupProductType_Failover
	} else if groupProductType == "loadbalancing" {
		groupProductTypeID = api.GroupProductType_LoadBalancing
	} else {
		return diag.FromErr(fmt.Errorf("invalid group_product_type: %s. It should be failover or loadbalancing.", groupProductType))
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

	groupType := d.Get("group_type").(string)
	groupTypeID := api.GroupType_Zone
	if strings.ToLower(groupType) == "cname" {
		groupTypeID = api.GroupType_CName
	} else if strings.ToLower(groupType) == "subdomain" {
		groupTypeID = api.GroupType_SubDomain
	}

	groupComposition := api.DNSGroupRecords{
		A:     *arrayAs,
		AAAA:  *arrayAAAAs,
		CName: *arrayCnames,
	}

	group := api.DnsRouteGroup{
		GroupID:            groupID,
		FixedGroupID:       fixedGroupID,
		ZoneId:             zoneID,
		FixedZoneID:        fixedGroupID,
		Name:               name,
		GroupTypeID:        groupTypeID,
		GroupProductTypeID: groupProductTypeID,
		GroupComposition:   groupComposition,
	}
	dnsrouteClient := api.NewDNSRouteAPIClient(*config)

	newGroupID, err := dnsrouteClient.AddGroup(&group)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(newGroupID))

	return ResourceGroupRead(ctx, d, m)
}

func ResourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)
	groupProductType := d.Get("group_product_type").(string)
	gType := "nogroup"
	if strings.ToLower(groupProductType) == "failover" {
		gType = "fo"
	} else if strings.ToLower(groupProductType) == "loadbalancing" {
		gType = "lb"
	}

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)
	resp, err := dnsRouteClient.GetGroup(groupID, gType)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("account_number", accountNumber)
	d.Set("group_id", groupID)

	d.Set("fixed_group_id", resp.FixedGroupID)
	d.Set("fixed_zone_id", resp.FixedZoneID)
	d.Set("group_product_type_id", resp.GroupProductTypeID)
	if resp.GroupProductTypeID == api.GroupProductType_Failover {
		d.Set("group_product_type", "failover")
	} else if resp.GroupProductTypeID == api.GroupProductType_LoadBalancing {
		d.Set("group_product_type", "loadbalancing")
	}
	d.Set("group_type_id", resp.GroupTypeID)
	if resp.GroupTypeID == api.GroupType_CName {
		d.Set("group_type", "cname")
	} else if resp.GroupTypeID == api.GroupType_SubDomain {
		d.Set("group_type", "subdomain")
	}

	d.Set("name", resp.Name)
	d.Set("zone_id", resp.ZoneId)

	d.Set("a", flattenGroupDNSs(&resp.GroupComposition.A))
	d.Set("aaaa", flattenGroupDNSs(&resp.GroupComposition.AAAA))
	d.Set("cname", flattenGroupDNSs(&resp.GroupComposition.CName))

	return diags
}

func ResourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	name := d.Get("name").(string)

	groupID := d.Get("group_id").(int)
	if groupID == 0 {
		groupID = -1
	}
	fixedGroupID := d.Get("fixed_group_id").(int)
	if fixedGroupID == 0 {
		fixedGroupID = -1
	}
	zoneID := d.Get("zone_id").(int)
	if zoneID == 0 {
		zoneID = -1
	}
	fixedZoneID := d.Get("fixed_zone_id").(int)
	if fixedZoneID == 0 {
		fixedZoneID = -1
	}
	groupType := d.Get("group_type").(string)
	groupTypeID := api.GroupType_Zone
	if strings.ToLower(groupType) == "cname" {
		groupTypeID = api.GroupType_CName
	} else if strings.ToLower(groupType) == "subdomain" {
		groupTypeID = api.GroupType_SubDomain
	}
	groupProductType := d.Get("group_product_type").(string)
	groupProductTypeID := api.GroupProductType_NoGroup
	if groupProductType == "failover" {
		groupProductTypeID = api.GroupProductType_Failover
	} else if groupProductType == "loadbalancing" {
		groupProductTypeID = api.GroupProductType_LoadBalancing
	} else {
		return diag.FromErr(fmt.Errorf("invalid group_product_type: %s. It should be failover or loadbalancing.", groupProductType))
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

	groupComposition := api.DNSGroupRecords{
		A:     *arrayAs,
		AAAA:  *arrayAAAAs,
		CName: *arrayCnames,
	}

	group := api.DnsRouteGroup{
		GroupID:            groupID,
		FixedGroupID:       fixedGroupID,
		ZoneId:             zoneID,
		FixedZoneID:        fixedGroupID,
		Name:               name,
		GroupTypeID:        groupTypeID,
		GroupProductTypeID: groupProductTypeID,
		GroupComposition:   groupComposition,
	}

	dnsrouteClient := api.NewDNSRouteAPIClient(*config)

	err = dnsrouteClient.UpdateGroup(&group)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(groupID))

	return ResourceGroupRead(ctx, d, m)
}

func ResourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	groupProductTypeID := d.Get("group_product_type_id").(int)
	groupProductType := ""
	if groupProductTypeID == api.GroupProductType_Failover {
		groupProductType = "fo"
	} else if groupProductTypeID == api.GroupProductType_LoadBalancing {
		groupProductType = "lb"
	}
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteAPIClient := api.NewDNSRouteAPIClient(*config)

	groupID, _ := strconv.Atoi(d.Id())

	err := dnsRouteAPIClient.DeleteGroup(groupID, groupProductType)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
