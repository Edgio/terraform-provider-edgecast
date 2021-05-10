// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"terraform-provider-vmp/vmp/api"
	"terraform-provider-vmp/vmp/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Master Server Group
func ResourceZone() *schema.Resource {
	// For debugging purpose
	time.Sleep(10 * time.Second)

	return &schema.Resource{
		CreateContext: ResourceZoneCreate,
		ReadContext:   ResourceZoneRead,
		UpdateContext: ResourceZoneUpdate,
		DeleteContext: ResourceZoneDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration.",
			},
			"zone_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Zone Type.",
			},
			"zone_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ZoneID",
			},
			"fixed_zone_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "FixedZoneID",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain Name.",
			},
			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Status.",
			},
			"status_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status.",
			},
			"is_customer_owned": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "is customer owned",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment",
			},

			"record_a": {
				Type:     schema.TypeList,
				Optional: true,
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
				Description: "List of A records",
			},
			"record_aaaa": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_cname": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_mx": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_ns": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_ptr": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_soa": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_spf": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_srv": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_txt": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_dnskey": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_rrsig": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_ds": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_nsec": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_nsec3": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_nsec3param": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_dlv": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"record_caa": {
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
				Optional:    true,
				Description: "List of A records",
			},
			"dnsroute_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
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
				},
			},
		},
	}
}

func ResourceZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	is_customer_owned := d.Get("is_customer_owned").(bool)
	comment := d.Get("comment").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber
	recordA := d.Get("record_a").([]interface{})
	recordAAAA := d.Get("record_aaaa").([]interface{})
	recordCNAME := d.Get("record_cname").([]interface{})
	recordMX := d.Get("record_mx").([]interface{})
	recordNS := d.Get("record_ns").([]interface{})
	recordPTR := d.Get("record_ptr").([]interface{})
	recordSOA := d.Get("record_soa").([]interface{})
	recordSPF := d.Get("record_spf").([]interface{})
	recordSRV := d.Get("record_srv").([]interface{})
	recordTXT := d.Get("record_txt").([]interface{})
	recordDNSKEY := d.Get("record_dnskey").([]interface{})
	recordRRSIG := d.Get("record_rrsig").([]interface{})
	recordDS := d.Get("record_ds").([]interface{})
	recordNSEC := d.Get("record_nsec").([]interface{})
	recordNSEC3 := d.Get("record_nsec3").([]interface{})
	recordNSEC3PARAM := d.Get("record_nsec3param").([]interface{})
	recordDLV := d.Get("record_dlv").([]interface{})
	recordCAA := d.Get("record_caa").([]interface{})

	zoneRequest := api.Zone{
		FixedZoneID:     -1,
		ZoneID:          -1,
		DomainName:      domainName,
		Status:          status,
		ZoneType:        zoneType,
		IsCustomerOwned: is_customer_owned,
		Comment:         comment,
	}

	aRecords := toDNSRecords("A", recordA)
	if len(*aRecords) > 0 {
		zoneRequest.Records.A = *aRecords
	}
	aaaaRecords := toDNSRecords("AAAA", recordAAAA)
	if len(*aaaaRecords) > 0 {
		zoneRequest.Records.AAAA = *aaaaRecords
	}
	cnameRecords := toDNSRecords("Cname", recordCNAME)
	if len(*cnameRecords) > 0 {
		zoneRequest.Records.CName = *cnameRecords
	}
	mxRecords := toDNSRecords("MX", recordMX)
	if len(*mxRecords) > 0 {
		zoneRequest.Records.MX = *mxRecords
	}
	nsRecords := toDNSRecords("NS", recordNS)
	if len(*nsRecords) > 0 {
		zoneRequest.Records.NS = *nsRecords
	}
	ptrRecords := toDNSRecords("PTR", recordPTR)
	if len(*ptrRecords) > 0 {
		zoneRequest.Records.PTR = *ptrRecords
	}
	soaRecords := toDNSRecords("SOA", recordSOA)
	if len(*soaRecords) > 0 {
		zoneRequest.Records.SOA = *soaRecords
	}
	spfRecords := toDNSRecords("SPF", recordSPF)
	if len(*spfRecords) > 0 {
		zoneRequest.Records.SPF = *spfRecords
	}
	srvRecords := toDNSRecords("SRV", recordSRV)
	if len(*srvRecords) > 0 {
		zoneRequest.Records.SRV = *srvRecords
	}
	txtRecords := toDNSRecords("TXT", recordTXT)
	if len(*txtRecords) > 0 {
		zoneRequest.Records.TXT = *txtRecords
	}
	dnsKeyRecords := toDNSRecords("DNSKEY", recordDNSKEY)
	if len(*dnsKeyRecords) > 0 {
		zoneRequest.Records.DNSKEY = *dnsKeyRecords
	}
	rrsigRecords := toDNSRecords("RRSIG", recordRRSIG)
	if len(*rrsigRecords) > 0 {
		zoneRequest.Records.RRSIG = *rrsigRecords
	}
	dsRecords := toDNSRecords("DS", recordDS)
	if len(*dsRecords) > 0 {
		zoneRequest.Records.DS = *dsRecords
	}
	nsecRecords := toDNSRecords("NSEC", recordNSEC)
	if len(*nsecRecords) > 0 {
		zoneRequest.Records.NSEC = *nsecRecords
	}
	nsec3Records := toDNSRecords("NSEC3", recordNSEC3)
	if len(*nsec3Records) > 0 {
		zoneRequest.Records.NSEC3 = *nsec3Records
	}
	nsec3paramRecords := toDNSRecords("NSEC3PARAM", recordNSEC3PARAM)
	if len(*nsec3paramRecords) > 0 {
		zoneRequest.Records.NSEC3PARAM = *nsec3paramRecords
	}
	dlvRecords := toDNSRecords("DLV", recordDLV)
	if len(*dlvRecords) > 0 {
		zoneRequest.Records.DLV = *dlvRecords
	}
	caaRecords := toDNSRecords("CAA", recordCAA)
	if len(*caaRecords) > 0 {
		zoneRequest.Records.CAA = *caaRecords
	}

	dnsGroups := d.Get("dnsroute_group").([]interface{})

	groups, err := toDNSRouteGroups(dnsGroups)
	if err != nil {
		diag.FromErr(err)
	}
	if groups != nil && len(*groups) > 0 {
		zoneRequest.Groups = *groups
	}
	dnsrouteClient := api.NewDNSRouteAPIClient(*config)

	helper.LogInstanceToPrettyJson("Zone Create Request body", zoneRequest)
	newZoneID, err := dnsrouteClient.AddZone(&zoneRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(newZoneID))

	return ResourceZoneRead(ctx, d, m)
}

func ResourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	zoneID, err := strconv.Atoi(d.Id())
	accountNumber := d.Get("account_number").(string)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteClient := api.NewDNSRouteAPIClient(*config)
	resp, err := dnsRouteClient.GetZone(zoneID)

	helper.LogInstanceToPrettyJson("ZoneRead Response body", resp)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("account_number", accountNumber)
	d.Set("domain_name", resp.DomainName)
	d.Set("status", resp.Status)
	d.Set("zone_type", resp.ZoneType)
	d.Set("is_customer_owned", resp.IsCustomerOwned)
	d.Set("comment", resp.Comment)
	d.Set("zone_id", zoneID)
	d.Set("status", resp.Status)
	d.Set("status_name", resp.StatusName)
	d.Set("is_customer_owned", resp.IsCustomerOwned)

	recordAs := flattenDnsRecords(&resp.Records.A)
	if err := d.Set("record_a", recordAs); err != nil {
		return diag.FromErr(err)
	}
	recordAAAAs := flattenDnsRecords(&resp.Records.AAAA)
	if err := d.Set("record_aaaa", recordAAAAs); err != nil {
		return diag.FromErr(err)
	}
	recordCNames := flattenDnsRecords(&resp.Records.CName)
	if err := d.Set("record_cname", recordCNames); err != nil {
		return diag.FromErr(err)
	}
	recordDLVs := flattenDnsRecords(&resp.Records.DLV)
	if err := d.Set("record_dlv", recordDLVs); err != nil {
		return diag.FromErr(err)
	}
	recordDNSKEYs := flattenDnsRecords(&resp.Records.DNSKEY)
	if err := d.Set("record_dnskey", recordDNSKEYs); err != nil {
		return diag.FromErr(err)
	}
	recordDSs := flattenDnsRecords(&resp.Records.DS)
	if err := d.Set("record_ds", recordDSs); err != nil {
		return diag.FromErr(err)
	}
	recordMXs := flattenDnsRecords(&resp.Records.MX)
	if err := d.Set("record_mx", recordMXs); err != nil {
		return diag.FromErr(err)
	}
	recordNSs := flattenDnsRecords(&resp.Records.NS)
	if err := d.Set("record_ns", recordNSs); err != nil {
		return diag.FromErr(err)
	}
	recordNSECs := flattenDnsRecords(&resp.Records.NSEC)
	if err := d.Set("record_nsec", recordNSECs); err != nil {
		return diag.FromErr(err)
	}
	recordNSEC3s := flattenDnsRecords(&resp.Records.NSEC3)
	if err := d.Set("record_nsec3", recordNSEC3s); err != nil {
		return diag.FromErr(err)
	}
	recordNSEC3PARAMs := flattenDnsRecords(&resp.Records.NSEC3PARAM)
	if err := d.Set("record_nsec3param", recordNSEC3PARAMs); err != nil {
		return diag.FromErr(err)
	}
	recordPTRs := flattenDnsRecords(&resp.Records.PTR)
	if err := d.Set("record_ptr", recordPTRs); err != nil {
		return diag.FromErr(err)
	}
	recordRRSIGs := flattenDnsRecords(&resp.Records.RRSIG)
	if err := d.Set("record_rrsig", recordRRSIGs); err != nil {
		return diag.FromErr(err)
	}
	recordSOAs := flattenDnsRecords(&resp.Records.SOA)
	if err := d.Set("record_soa", recordSOAs); err != nil {
		return diag.FromErr(err)
	}
	recordSPFs := flattenDnsRecords(&resp.Records.SPF)
	if err := d.Set("record_spf", recordSPFs); err != nil {
		return diag.FromErr(err)
	}
	recordSRVs := flattenDnsRecords(&resp.Records.SRV)
	if err := d.Set("record_srv", recordSRVs); err != nil {
		return diag.FromErr(err)
	}
	recordTXTs := flattenDnsRecords(&resp.Records.TXT)
	if err := d.Set("record_txt", recordTXTs); err != nil {
		return diag.FromErr(err)
	}
	recordCAAs := flattenDnsRecords(&resp.Records.CAA)
	if err := d.Set("record_caa", recordCAAs); err != nil {
		return diag.FromErr(err)
	}

	dnsGroups := flattenDnsGroups(&resp.Groups)

	if dnsGroups != nil {
		if err := d.Set("dnsroute_group", dnsGroups); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(strconv.Itoa(resp.FixedZoneID))
	return diags
}

func ResourceZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	is_customer_owned := d.Get("is_customer_owned").(bool)
	comment := d.Get("comment").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber
	recordA := d.Get("record_a").([]interface{})
	recordAAAA := d.Get("record_aaaa").([]interface{})
	recordCNAME := d.Get("record_cname").([]interface{})
	recordMX := d.Get("record_mx").([]interface{})
	recordNS := d.Get("record_ns").([]interface{})
	recordPTR := d.Get("record_ptr").([]interface{})
	recordSOA := d.Get("record_soa").([]interface{})
	recordSPF := d.Get("record_spf").([]interface{})
	recordSRV := d.Get("record_srv").([]interface{})
	recordTXT := d.Get("record_txt").([]interface{})
	recordDNSKEY := d.Get("record_dnskey").([]interface{})
	recordRRSIG := d.Get("record_rrsig").([]interface{})
	recordDS := d.Get("record_ds").([]interface{})
	recordNSEC := d.Get("record_nsec").([]interface{})
	recordNSEC3 := d.Get("record_nsec3").([]interface{})
	recordNSEC3PARAM := d.Get("record_nsec3param").([]interface{})
	recordDLV := d.Get("record_dlv").([]interface{})
	recordCAA := d.Get("record_caa").([]interface{})

	zoneRequest := api.Zone{
		FixedZoneID:     -1,
		ZoneID:          -1,
		DomainName:      domainName,
		Status:          status,
		ZoneType:        zoneType,
		IsCustomerOwned: is_customer_owned,
		Comment:         comment,
	}
	zoneID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	zoneRequest.FixedZoneID = zoneID

	aRecords := toDNSRecords("A", recordA)
	if len(*aRecords) > 0 {
		zoneRequest.Records.A = *aRecords
	}
	aaaaRecords := toDNSRecords("AAAA", recordAAAA)
	if len(*aaaaRecords) > 0 {
		zoneRequest.Records.AAAA = *aaaaRecords
	}
	cnameRecords := toDNSRecords("Cname", recordCNAME)
	if len(*cnameRecords) > 0 {
		zoneRequest.Records.CName = *cnameRecords
	}
	mxRecords := toDNSRecords("MX", recordMX)
	if len(*mxRecords) > 0 {
		zoneRequest.Records.MX = *mxRecords
	}
	nsRecords := toDNSRecords("NS", recordNS)
	if len(*nsRecords) > 0 {
		zoneRequest.Records.NS = *nsRecords
	}
	ptrRecords := toDNSRecords("PTR", recordPTR)
	if len(*ptrRecords) > 0 {
		zoneRequest.Records.PTR = *ptrRecords
	}
	soaRecords := toDNSRecords("SOA", recordSOA)
	if len(*soaRecords) > 0 {
		zoneRequest.Records.SOA = *soaRecords
	}
	spfRecords := toDNSRecords("SPF", recordSPF)
	if len(*spfRecords) > 0 {
		zoneRequest.Records.SPF = *spfRecords
	}
	srvRecords := toDNSRecords("SRV", recordSRV)
	if len(*srvRecords) > 0 {
		zoneRequest.Records.SRV = *srvRecords
	}
	txtRecords := toDNSRecords("TXT", recordTXT)
	if len(*txtRecords) > 0 {
		zoneRequest.Records.TXT = *txtRecords
	}
	dnsKeyRecords := toDNSRecords("DNSKEY", recordDNSKEY)
	if len(*dnsKeyRecords) > 0 {
		zoneRequest.Records.DNSKEY = *dnsKeyRecords
	}
	rrsigRecords := toDNSRecords("RRSIG", recordRRSIG)
	if len(*rrsigRecords) > 0 {
		zoneRequest.Records.RRSIG = *rrsigRecords
	}
	dsRecords := toDNSRecords("DS", recordDS)
	if len(*dsRecords) > 0 {
		zoneRequest.Records.DS = *dsRecords
	}
	nsecRecords := toDNSRecords("NSEC", recordNSEC)
	if len(*nsecRecords) > 0 {
		zoneRequest.Records.NSEC = *nsecRecords
	}
	nsec3Records := toDNSRecords("NSEC3", recordNSEC3)
	if len(*nsec3Records) > 0 {
		zoneRequest.Records.NSEC3 = *nsec3Records
	}
	nsec3paramRecords := toDNSRecords("NSEC3PARAM", recordNSEC3PARAM)
	if len(*nsec3paramRecords) > 0 {
		zoneRequest.Records.NSEC3PARAM = *nsec3paramRecords
	}
	dlvRecords := toDNSRecords("DLV", recordDLV)
	if len(*dlvRecords) > 0 {
		zoneRequest.Records.DLV = *dlvRecords
	}
	caaRecords := toDNSRecords("CAA", recordCAA)
	if len(*caaRecords) > 0 {
		zoneRequest.Records.CAA = *caaRecords
	}

	dnsGroups := d.Get("dnsroute_group").([]interface{})

	groups, err := toDNSRouteGroups(dnsGroups)
	if err != nil {
		diag.FromErr(err)
	}
	if groups != nil && len(*groups) > 0 {
		zoneRequest.Groups = *groups
	}
	helper.LogInstanceToPrettyJson("Update request body:", zoneRequest)
	dnsrouteClient := api.NewDNSRouteAPIClient(*config)
	err = dnsrouteClient.UpdateZone(&zoneRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceZoneRead(ctx, d, m)
}

func ResourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	dnsRouteAPIClient := api.NewDNSRouteAPIClient(*config)

	zoneID, _ := strconv.Atoi(d.Id())

	err := dnsRouteAPIClient.DeleteZone(zoneID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenDnsRecords(recordItems *[]api.DNSRecord) []interface{} {
	if recordItems != nil {
		dnsRecords := make([]interface{}, len(*recordItems), len(*recordItems))

		for i, dns := range *recordItems {
			item := make(map[string]interface{})

			item["name"] = dns.Name
			item["ttl"] = dns.TTL
			item["rdata"] = dns.Rdata
			item["verify_id"] = dns.VerifyID
			item["fixed_record_id"] = dns.FixedRecordID
			item["fixed_group_id"] = dns.FixedGroupID
			item["weight"] = dns.Weight
			item["record_type_id"] = dns.RecordTypeID
			item["record_type_name"] = dns.RecordTypeName
			item["group_id"] = dns.GroupID
			item["is_delete"] = dns.IsDeleted

			dnsRecords[i] = item
		}
		return dnsRecords
	}
	return make([]interface{}, 0)
}

func flattenDnsGroups(groupItems *[]api.DnsRouteGroup) []interface{} {
	if groupItems != nil && len(*groupItems) > 0 {
		groupArr := make([]interface{}, len(*groupItems), len(*groupItems))

		for i, group := range *groupItems {
			item := make(map[string]interface{})

			item["group_id"] = group.GroupID
			item["id"] = group.ID
			item["fixed_group_id"] = group.FixedGroupID
			item["fixed_zone_id"] = group.FixedZoneID
			item["group_product_type_id"] = group.GroupProductTypeID
			if group.GroupProductTypeID == api.GroupProductType_Failover {
				item["group_product_type"] = "failover"
			} else if group.GroupProductTypeID == api.GroupProductType_LoadBalancing {
				item["group_product_type"] = "loadbalancing"
			}
			item["group_type_id"] = group.GroupTypeID
			if group.GroupTypeID == api.GroupType_Zone {
				item["group_type"] = "zone"
			}
			item["name"] = group.Name
			item["zone_id"] = group.ZoneId

			item["a"] = flattenGroupDNSs(&group.GroupComposition.RA)
			item["aaaa"] = flattenGroupDNSs(&group.GroupComposition.RAAAA)
			item["cname"] = flattenGroupDNSs(&group.GroupComposition.RCName)

			groupArr[i] = item
		}
		return groupArr
	}
	return nil
}

func flattenGroupDNSs(dnsItems *[]api.DnsRouteGroupRecord) []interface{} {
	if dnsItems != nil && len(*dnsItems) > 0 {

		dnsArr := make([]interface{}, len(*dnsItems), len(*dnsItems))

		for i, dns := range *dnsItems {
			item := make(map[string]interface{})

			item["weight"] = dns.Record.Weight
			//item["id"] = dns.ID
			record := flattenGroupDnsRecord(&dns.Record)

			if record != nil {
				item["record"] = record
			}

			healthCheck := flattenHealthCheck(dns.HealthCheck)

			if healthCheck != nil {
				item["health_check"] = healthCheck
			}

			dnsArr[i] = item
		}
		return dnsArr
	}
	return make([]interface{}, 0)
}

func flattenGroupDnsRecord(dns *api.DNSRecord) []interface{} {
	if dns != nil && (*dns != api.DNSRecord{}) {
		record := make([]interface{}, 1, 1)
		m := make(map[string]interface{})
		m["fixed_record_id"] = dns.FixedRecordID
		m["fixed_group_id"] = dns.FixedGroupID
		m["group_id"] = dns.GroupID
		m["is_delete"] = dns.IsDeleted
		m["name"] = dns.Name
		m["rdata"] = dns.Rdata
		m["record_id"] = dns.RecordID
		m["record_type_id"] = dns.RecordTypeID
		m["record_type_name"] = dns.RecordTypeName
		m["ttl"] = dns.TTL
		m["verify_id"] = dns.VerifyID
		m["weight"] = dns.Weight
		record[0] = m
		return record
	}
	return nil
}

func flattenHealthCheck(hc *api.HealthCheck) []interface{} {
	if hc != nil && (*hc != api.HealthCheck{}) {
		record := make([]interface{}, 1, 1)
		healthCheck := make(map[string]interface{})

		healthCheck["check_interval"] = hc.CheckInterval
		healthCheck["check_type_id"] = hc.CheckTypeID
		healthCheck["content_verification"] = hc.ContentVerification
		healthCheck["email_notification_address"] = hc.EmailNotificationAddress
		healthCheck["failed_check_threshold"] = hc.FailedCheckThreshold
		healthCheck["fixed_id"] = hc.FixedID
		//healthCheck["fixed_group_id"] = hc.FixedGroupID
		healthCheck["fixed_record_id"] = hc.FixedRecordID
		healthCheck["group_id"] = hc.GroupID
		healthCheck["id"] = hc.ID
		healthCheck["http_method_id"] = hc.HTTPMethodID
		healthCheck["ip_address"] = hc.IPAddress
		healthCheck["ip_version"] = hc.IPVersion
		healthCheck["port_number"] = hc.PortNumber
		healthCheck["record_id"] = hc.RecordID
		healthCheck["reintegration_method_id"] = hc.ReintegrationMethodID
		healthCheck["status"] = hc.Status
		healthCheck["status_name"] = hc.StatusName
		healthCheck["timeout"] = hc.TimeOut
		healthCheck["uri"] = hc.Uri
		//healthCheck["user_id"] = hc.UserID
		//healthCheck[""] = hc.WhiteListedHc

		record[0] = healthCheck
		return record
	}
	return nil
}

func toDNSRecords(recodeType string, input []interface{}) *[]api.DNSRecord {

	records := make([]api.DNSRecord, 0)

	for _, item := range input {
		curr := item.(map[string]interface{})

		name := curr["name"].(string)
		ttl := curr["ttl"].(int)
		rdata := curr["rdata"].(string)
		verifyID := curr["verify_id"].(int)
		recordID := curr["record_id"].(int)
		fixedRecordID := curr["fixed_record_id"].(int)
		fixedGroupID := curr["fixed_group_id"].(int)
		groupID := curr["group_id"].(int)
		isDeleted := curr["is_delete"].(bool)
		weight := curr["weight"].(int)
		recordTypeID := curr["record_type_id"].(int)
		recordTypeName := curr["record_type_name"].(string)
		record := api.DNSRecord{
			RecordID:       recordID,
			FixedRecordID:  fixedRecordID,
			FixedGroupID:   fixedGroupID,
			GroupID:        groupID,
			IsDeleted:      isDeleted,
			Name:           name,
			TTL:            ttl,
			Rdata:          rdata,
			Weight:         weight,
			RecordTypeID:   recordTypeID,
			RecordTypeName: recordTypeName,
			VerifyID:       verifyID,
		}

		records = append(records, record)
	}
	return &records
}

func toDNSRouteGroups(input []interface{}) (*[]api.DnsRouteGroup, error) {
	if input == nil || len(input) == 0 {
		return nil, nil
	}
	groups := make([]api.DnsRouteGroup, 0)

	for _, item := range input {
		curr := item.(map[string]interface{})

		name := curr["name"].(string)
		groupType := curr["group_type"].(string)
		if groupType != "zone" {
			return nil, fmt.Errorf("invalid group_type: %s. It should be zone.", groupType)
		}
		groupID := curr["group_id"].(int)
		if groupID == 0 {
			groupID = -1
		}
		fixedGroupID := curr["fixed_group_id"].(int)
		if fixedGroupID == 0 {
			fixedGroupID = -1
		}
		zoneID := curr["zone_id"].(int)
		if zoneID == 0 {
			zoneID = -1
		}
		fixedZoneID := curr["fixed_zone_id"].(int)
		if fixedZoneID == 0 {
			fixedZoneID = -1
		}
		groupProductType := curr["group_product_type"].(string)
		groupProductTypeID := api.GroupProductType_NoGroup
		if groupProductType == "failover" {
			groupProductTypeID = api.GroupProductType_Failover
		} else if groupProductType == "loadbalancing" {
			groupProductTypeID = api.GroupProductType_LoadBalancing
		} else {
			return nil, fmt.Errorf("invalid group_product_type: %s. It should be failover or loadbalancing.", groupProductType)
		}

		dnsAs := curr["a"].([]interface{})
		arrayAs, err := toGroupRecords(dnsAs)
		if err != nil {
			return nil, err
		}
		dnsAAAAs := curr["aaaa"].([]interface{})
		arrayAAAAs, err := toGroupRecords(dnsAAAAs)
		if err != nil {
			return nil, err
		}
		dnsCnames := curr["cname"].([]interface{})
		arrayCnames, err := toGroupRecords(dnsCnames)
		if err != nil {
			return nil, err
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
			GroupTypeID:        api.GroupType_Zone,
			GroupProductTypeID: groupProductTypeID,
			GroupComposition:   groupComposition,
		}
		groups = append(groups, group)
	}

	return &groups, nil
}

func toGroupRecords(input []interface{}) (*[]api.DnsRouteGroupRecord, error) {
	records := make([]api.DnsRouteGroupRecord, 0)

	for _, item := range input {
		curr := item.(map[string]interface{})
		var healthCheck *api.HealthCheck = nil
		var err error
		if val, ok := curr["health_check"]; ok {

			hc := val.([]interface{})
			healthCheck, err = toHealthCheck(hc)
			if err != nil {
				return nil, err
			}
			if healthCheck == nil || (api.HealthCheck{}) == *healthCheck {
				healthCheck = nil
			}
		}

		weight := curr["weight"].(int)

		var dnsRecord *api.DNSRecord = nil

		if val, ok := curr["record"]; ok {
			rc := val.([]interface{})
			dnsRecord, err = toDNSRouteZoneRecord(rc, weight)
			if err != nil {
				return nil, err
			}
			if (api.DNSRecord{}) == *dnsRecord {
				dnsRecord = nil
			}
		}

		record := api.DnsRouteGroupRecord{
			HealthCheck: healthCheck,
			Record:      *dnsRecord,
		}

		records = append(records, record)
	}
	return &records, nil
}

func toHealthCheck(input []interface{}) (*api.HealthCheck, error) {
	if input != nil && len(input) > 0 {
		for _, element := range input {
			item := element.(map[string]interface{})
			healthCheck := api.HealthCheck{}
			if item["check_interval"] != nil {
				healthCheck.CheckInterval = item["check_interval"].(int)
			}
			if item["content_verification"] != nil {
				healthCheck.ContentVerification = item["content_verification"].(string)
			}
			if item["email_notification_address"] != nil {
				healthCheck.EmailNotificationAddress = item["email_notification_address"].(string)
			}
			if item["check_type_id"] != nil {
				healthCheck.CheckTypeID = item["check_type_id"].(int)
			}
			if item["failed_check_threshold"] != nil {
				healthCheck.FailedCheckThreshold = item["failed_check_threshold"].(int)
			}
			if item["http_method_id"] != nil {
				healthCheck.HTTPMethodID = item["http_method_id"].(int)
			}
			if item["ip_address"] != nil {
				healthCheck.IPAddress = item["ip_address"].(string)
			}
			if item["ip_version"] != nil {
				healthCheck.IPVersion = item["ip_version"].(int)
			}
			if item["port_number"] != nil {
				healthCheck.PortNumber = item["port_number"].(string)
			}
			if item["reintegration_method_id"] != nil {
				methodID := item["reintegration_method_id"].(int)
				healthCheck.ReintegrationMethodID = methodID
			}
			if item["status"] != nil {
				healthCheck.Status = item["status"].(int)
			}
			if item["status_name"] != nil {
				healthCheck.StatusName = item["status_name"].(string)
			}
			if item["uri"] != nil {
				healthCheck.Uri = item["uri"].(string)
			}
			return &healthCheck, nil
		}
	}
	return nil, nil
}

func toDNSRouteZoneRecord(items []interface{}, weight int) (*api.DNSRecord, error) {
	if items != nil && len(items) > 0 {
		for _, element := range items {
			item := element.(map[string]interface{})
			record := api.DNSRecord{
				FixedRecordID:  item["fixed_record_id"].(int),
				FixedGroupID:   item["fixed_group_id"].(int),
				RecordID:       item["record_id"].(int),
				Weight:         weight,
				RecordTypeID:   item["record_type_id"].(int),
				RecordTypeName: item["record_type_name"].(string),
				GroupID:        item["group_id"].(int),
				IsDeleted:      item["is_delete"].(bool),
				Name:           item["name"].(string),
				TTL:            item["ttl"].(int),
				Rdata:          item["rdata"].(string),
				VerifyID:       item["verify_id"].(int),
			}

			return &record, nil
		}
	}
	return nil, nil
}
