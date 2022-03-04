// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Master Server Group
func ResourceZone() *schema.Resource {
	// For debugging purpose
	//time.Sleep(10 * time.Second)

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

func ResourceZoneCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Build Zone Object
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	comment := d.Get("comment").(string)
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

	zone := routedns.Zone{
		DomainName: domainName,
		Status:     status,
		ZoneType:   zoneType,
		Comment:    comment,
	}

	aRecords := toDNSRecords("A", &recordA)
	if len(*aRecords) > 0 {
		zone.Records.A = *aRecords
	}
	aaaaRecords := toDNSRecords("AAAA", &recordAAAA)
	if len(*aaaaRecords) > 0 {
		zone.Records.AAAA = *aaaaRecords
	}
	cnameRecords := toDNSRecords("Cname", &recordCNAME)
	if len(*cnameRecords) > 0 {
		zone.Records.CNAME = *cnameRecords
	}
	mxRecords := toDNSRecords("MX", &recordMX)
	if len(*mxRecords) > 0 {
		zone.Records.MX = *mxRecords
	}
	nsRecords := toDNSRecords("NS", &recordNS)
	if len(*nsRecords) > 0 {
		zone.Records.NS = *nsRecords
	}
	ptrRecords := toDNSRecords("PTR", &recordPTR)
	if len(*ptrRecords) > 0 {
		zone.Records.PTR = *ptrRecords
	}
	soaRecords := toDNSRecords("SOA", &recordSOA)
	if len(*soaRecords) > 0 {
		zone.Records.SOA = *soaRecords
	}
	spfRecords := toDNSRecords("SPF", &recordSPF)
	if len(*spfRecords) > 0 {
		zone.Records.SPF = *spfRecords
	}
	srvRecords := toDNSRecords("SRV", &recordSRV)
	if len(*srvRecords) > 0 {
		zone.Records.SRV = *srvRecords
	}
	txtRecords := toDNSRecords("TXT", &recordTXT)
	if len(*txtRecords) > 0 {
		zone.Records.TXT = *txtRecords
	}
	dnsKeyRecords := toDNSRecords("DNSKEY", &recordDNSKEY)
	if len(*dnsKeyRecords) > 0 {
		zone.Records.DNSKEY = *dnsKeyRecords
	}
	rrsigRecords := toDNSRecords("RRSIG", &recordRRSIG)
	if len(*rrsigRecords) > 0 {
		zone.Records.RRSIG = *rrsigRecords
	}
	dsRecords := toDNSRecords("DS", &recordDS)
	if len(*dsRecords) > 0 {
		zone.Records.DS = *dsRecords
	}
	nsecRecords := toDNSRecords("NSEC", &recordNSEC)
	if len(*nsecRecords) > 0 {
		zone.Records.NSEC = *nsecRecords
	}
	nsec3Records := toDNSRecords("NSEC3", &recordNSEC3)
	if len(*nsec3Records) > 0 {
		zone.Records.NSEC3 = *nsec3Records
	}
	nsec3paramRecords := toDNSRecords("NSEC3PARAM", &recordNSEC3PARAM)
	if len(*nsec3paramRecords) > 0 {
		zone.Records.NSEC3PARAM = *nsec3paramRecords
	}
	dlvRecords := toDNSRecords("DLV", &recordDLV)
	if len(*dlvRecords) > 0 {
		zone.Records.DLV = *dlvRecords
	}
	caaRecords := toDNSRecords("CAA", &recordCAA)
	if len(*caaRecords) > 0 {
		zone.Records.CAA = *caaRecords
	}

	dnsGroups := d.Get("dnsroute_group").([]interface{})

	groups, err := toDNSRouteGroups(&dnsGroups)
	if err != nil {
		d.SetId("")
		diag.FromErr(err)
	}
	if groups != nil && len(*groups) > 0 {
		zone.Groups = *toCreateDNSGroups(groups)
	}

	// Call add Zone API
	params := routedns.NewAddZoneParams()
	params.AccountNumber = accountNumber
	params.Zone = zone
	zoneID, err := routeDNSService.AddZone(*params)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*zoneID))

	return ResourceZoneRead(ctx, d, m)
}

func ResourceZoneRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	zoneID, err := strconv.Atoi(d.Id())
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Call Get Zone API
	params := routedns.NewGetZoneParams()
	params.AccountNumber = accountNumber
	params.ZoneID = zoneID
	zoneObj, err := routeDNSService.GetZone(*params)

	if err != nil {
		return diag.FromErr(err)
	}

	// Update Terraform state with retrieved Zone data
	d.Set("account_number", accountNumber)
	d.Set("domain_name", zoneObj.DomainName)
	d.Set("status", zoneObj.Status)
	d.Set("zone_type", zoneObj.ZoneType)
	d.Set("is_customer_owned", zoneObj.IsCustomerOwned)
	d.Set("comment", zoneObj.Comment)
	d.Set("zone_id", zoneID)
	d.Set("status", zoneObj.Status)
	d.Set("status_name", zoneObj.StatusName)
	d.Set("is_customer_owned", zoneObj.IsCustomerOwned)

	recordAs := flattenDnsRecords(&zoneObj.Records.A)
	if err := d.Set("record_a", recordAs); err != nil {
		return diag.FromErr(err)
	}
	recordAAAAs := flattenDnsRecords(&zoneObj.Records.AAAA)
	if err := d.Set("record_aaaa", recordAAAAs); err != nil {
		return diag.FromErr(err)
	}
	recordCNames := flattenDnsRecords(&zoneObj.Records.CNAME)
	if err := d.Set("record_cname", recordCNames); err != nil {
		return diag.FromErr(err)
	}
	recordDLVs := flattenDnsRecords(&zoneObj.Records.DLV)
	if err := d.Set("record_dlv", recordDLVs); err != nil {
		return diag.FromErr(err)
	}
	recordDNSKEYs := flattenDnsRecords(&zoneObj.Records.DNSKEY)
	if err := d.Set("record_dnskey", recordDNSKEYs); err != nil {
		return diag.FromErr(err)
	}
	recordDSs := flattenDnsRecords(&zoneObj.Records.DS)
	if err := d.Set("record_ds", recordDSs); err != nil {
		return diag.FromErr(err)
	}
	recordMXs := flattenDnsRecords(&zoneObj.Records.MX)
	if err := d.Set("record_mx", recordMXs); err != nil {
		return diag.FromErr(err)
	}
	recordNSs := flattenDnsRecords(&zoneObj.Records.NS)
	if err := d.Set("record_ns", recordNSs); err != nil {
		return diag.FromErr(err)
	}
	recordNSECs := flattenDnsRecords(&zoneObj.Records.NSEC)
	if err := d.Set("record_nsec", recordNSECs); err != nil {
		return diag.FromErr(err)
	}
	recordNSEC3s := flattenDnsRecords(&zoneObj.Records.NSEC3)
	if err := d.Set("record_nsec3", recordNSEC3s); err != nil {
		return diag.FromErr(err)
	}
	recordNSEC3PARAMs := flattenDnsRecords(&zoneObj.Records.NSEC3PARAM)
	if err := d.Set("record_nsec3param", recordNSEC3PARAMs); err != nil {
		return diag.FromErr(err)
	}
	recordPTRs := flattenDnsRecords(&zoneObj.Records.PTR)
	if err := d.Set("record_ptr", recordPTRs); err != nil {
		return diag.FromErr(err)
	}
	recordRRSIGs := flattenDnsRecords(&zoneObj.Records.RRSIG)
	if err := d.Set("record_rrsig", recordRRSIGs); err != nil {
		return diag.FromErr(err)
	}
	recordSOAs := flattenDnsRecords(&zoneObj.Records.SOA)
	if err := d.Set("record_soa", recordSOAs); err != nil {
		return diag.FromErr(err)
	}
	recordSPFs := flattenDnsRecords(&zoneObj.Records.SPF)
	if err := d.Set("record_spf", recordSPFs); err != nil {
		return diag.FromErr(err)
	}
	recordSRVs := flattenDnsRecords(&zoneObj.Records.SRV)
	if err := d.Set("record_srv", recordSRVs); err != nil {
		return diag.FromErr(err)
	}
	recordTXTs := flattenDnsRecords(&zoneObj.Records.TXT)
	if err := d.Set("record_txt", recordTXTs); err != nil {
		return diag.FromErr(err)
	}
	recordCAAs := flattenDnsRecords(&zoneObj.Records.CAA)
	if err := d.Set("record_caa", recordCAAs); err != nil {
		return diag.FromErr(err)
	}

	dnsGroups := flattenDnsGroups(&zoneObj.Groups)

	if dnsGroups != nil {
		if err := d.Set("dnsroute_group", dnsGroups); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(strconv.Itoa(zoneObj.FixedZoneID))
	return diag.Diagnostics{}
}

func ResourceZoneUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	zoneID, err := strconv.Atoi(d.Id())
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)

	// Call Get Zone API
	getParams := routedns.NewGetZoneParams()
	getParams.AccountNumber = accountNumber
	getParams.ZoneID = zoneID
	zoneObj, err := routeDNSService.GetZone(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Build Zone Update Data
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	is_customer_owned := d.Get("is_customer_owned").(bool)
	comment := d.Get("comment").(string)
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
	dnsGroups := d.Get("dnsroute_group").([]interface{})

	groups, err := toDNSRouteGroups(&dnsGroups)
	if err != nil {
		diag.FromErr(err)
	}

	// Creating this object to pass into updateIDs and then to be used for
	// updating the retrieved zone from the API. This should be refactored once
	// we better understand the bug fixed by updateIDs
	comparisonZone := routedns.ZoneGetOK{
		Zone: routedns.Zone{
			DomainName: domainName,
			Status:     status,
			ZoneType:   zoneType,
			Comment:    comment,
			Records: routedns.DNSRecords{
				A:          *toDNSRecords("A", &recordA),
				AAAA:       *toDNSRecords("AAAA", &recordAAAA),
				CNAME:      *toDNSRecords("Cname", &recordCNAME),
				MX:         *toDNSRecords("MX", &recordMX),
				NS:         *toDNSRecords("NS", &recordNS),
				PTR:        *toDNSRecords("PTR", &recordPTR),
				SOA:        *toDNSRecords("SOA", &recordSOA),
				SPF:        *toDNSRecords("SPF", &recordSPF),
				SRV:        *toDNSRecords("SRV", &recordSRV),
				TXT:        *toDNSRecords("TXT", &recordTXT),
				DNSKEY:     *toDNSRecords("DNSKEY", &recordDNSKEY),
				RRSIG:      *toDNSRecords("RRSIG", &recordRRSIG),
				DS:         *toDNSRecords("DS", &recordDS),
				NSEC:       *toDNSRecords("NSEC", &recordNSEC),
				NSEC3:      *toDNSRecords("NSEC3", &recordNSEC3),
				NSEC3PARAM: *toDNSRecords("NSEC3PARAM", &recordNSEC3PARAM),
				DLV:        *toDNSRecords("DLV", &recordDLV),
				CAA:        *toDNSRecords("CAA", &recordCAA),
			},
		},
		Groups: *groups,
	}

	updateIDs(&comparisonZone, zoneObj)

	// Update Zone Object
	zoneObj.DomainName = domainName
	zoneObj.Status = status
	zoneObj.ZoneType = zoneType
	zoneObj.IsCustomerOwned = is_customer_owned
	zoneObj.Comment = comment

	aRecords := toDNSRecords("A", &recordA)
	if len(*aRecords) > 0 {
		zoneObj.Records.A = *aRecords
	}
	aaaaRecords := toDNSRecords("AAAA", &recordAAAA)
	if len(*aaaaRecords) > 0 {
		zoneObj.Records.AAAA = *aaaaRecords
	}
	cnameRecords := toDNSRecords("Cname", &recordCNAME)
	if len(*cnameRecords) > 0 {
		zoneObj.Records.CNAME = *cnameRecords
	}
	mxRecords := toDNSRecords("MX", &recordMX)
	if len(*mxRecords) > 0 {
		zoneObj.Records.MX = *mxRecords
	}
	nsRecords := toDNSRecords("NS", &recordNS)
	if len(*nsRecords) > 0 {
		zoneObj.Records.NS = *nsRecords
	}
	ptrRecords := toDNSRecords("PTR", &recordPTR)
	if len(*ptrRecords) > 0 {
		zoneObj.Records.PTR = *ptrRecords
	}
	soaRecords := toDNSRecords("SOA", &recordSOA)
	if len(*soaRecords) > 0 {
		zoneObj.Records.SOA = *soaRecords
	}
	spfRecords := toDNSRecords("SPF", &recordSPF)
	if len(*spfRecords) > 0 {
		zoneObj.Records.SPF = *spfRecords
	}
	srvRecords := toDNSRecords("SRV", &recordSRV)
	if len(*srvRecords) > 0 {
		zoneObj.Records.SRV = *srvRecords
	}
	txtRecords := toDNSRecords("TXT", &recordTXT)
	if len(*txtRecords) > 0 {
		zoneObj.Records.TXT = *txtRecords
	}
	dnsKeyRecords := toDNSRecords("DNSKEY", &recordDNSKEY)
	if len(*dnsKeyRecords) > 0 {
		zoneObj.Records.DNSKEY = *dnsKeyRecords
	}
	rrsigRecords := toDNSRecords("RRSIG", &recordRRSIG)
	if len(*rrsigRecords) > 0 {
		zoneObj.Records.RRSIG = *rrsigRecords
	}
	dsRecords := toDNSRecords("DS", &recordDS)
	if len(*dsRecords) > 0 {
		zoneObj.Records.DS = *dsRecords
	}
	nsecRecords := toDNSRecords("NSEC", &recordNSEC)
	if len(*nsecRecords) > 0 {
		zoneObj.Records.NSEC = *nsecRecords
	}
	nsec3Records := toDNSRecords("NSEC3", &recordNSEC3)
	if len(*nsec3Records) > 0 {
		zoneObj.Records.NSEC3 = *nsec3Records
	}
	nsec3paramRecords := toDNSRecords("NSEC3PARAM", &recordNSEC3PARAM)
	if len(*nsec3paramRecords) > 0 {
		zoneObj.Records.NSEC3PARAM = *nsec3paramRecords
	}
	dlvRecords := toDNSRecords("DLV", &recordDLV)
	if len(*dlvRecords) > 0 {
		zoneObj.Records.DLV = *dlvRecords
	}
	caaRecords := toDNSRecords("CAA", &recordCAA)
	if len(*caaRecords) > 0 {
		zoneObj.Records.CAA = *caaRecords
	}

	if groups != nil && len(*groups) > 0 {
		zoneObj.Groups = *groups
	}

	// Call update Zone API
	updateParams := routedns.NewUpdateZoneParams()
	updateParams.AccountNumber = accountNumber
	updateParams.Zone = *zoneObj
	err = routeDNSService.UpdateZone(*updateParams)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceZoneRead(ctx, d, m)
}

func ResourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(**api.ClientConfig)
	routeDNSService, err := buildRouteDNSService(**config)
	zoneID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Call Get Zone API
	getParams := routedns.NewGetZoneParams()
	getParams.AccountNumber = accountNumber
	getParams.ZoneID = zoneID
	zoneObj, err := routeDNSService.GetZone(*getParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Delete Zone API
	deleteParams := routedns.NewDeleteZoneParams()
	deleteParams.AccountNumber = accountNumber
	deleteParams.Zone = *zoneObj
	err = routeDNSService.DeleteZone(*deleteParams)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}

//1. to[Func]s: These functions are used to generate Zone request body
func toDNSRecords(recodeType string, input *[]interface{}) *[]routedns.DNSRecord {

	records := make([]routedns.DNSRecord, 0)

	for _, item := range *input {
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
		record := routedns.DNSRecord{
			RecordID:       recordID,
			FixedRecordID:  fixedRecordID,
			FixedGroupID:   fixedGroupID,
			GroupID:        groupID,
			IsDeleted:      isDeleted,
			Name:           name,
			TTL:            ttl,
			Rdata:          rdata,
			Weight:         weight,
			RecordTypeID:   routedns.RecordType(recordTypeID),
			RecordTypeName: recordTypeName,
			VerifyID:       verifyID,
		}

		records = append(records, record)
	}
	return &records
}

func toDNSRouteGroups(input *[]interface{}) (*[]routedns.DnsRouteGroupOK, error) {
	if *input == nil || len(*input) == 0 {
		return nil, nil
	}
	groups := make([]routedns.DnsRouteGroupOK, 0)

	for _, item := range *input {
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
		groupProductTypeID := routedns.NoGroup
		if groupProductType == "failover" {
			groupProductTypeID = routedns.Failover
		} else if groupProductType == "loadbalancing" {
			groupProductTypeID = routedns.LoadBalancing
		} else {
			return nil, fmt.Errorf("invalid group_product_type: %s. It should be failover or loadbalancing.", groupProductType)
		}

		dnsAs := curr["a"].([]interface{})
		arrayAs, err := toGroupRecords(&dnsAs)
		if err != nil {
			return nil, err
		}
		dnsAAAAs := curr["aaaa"].([]interface{})
		arrayAAAAs, err := toGroupRecords(&dnsAAAAs)
		if err != nil {
			return nil, err
		}
		dnsCnames := curr["cname"].([]interface{})
		arrayCnames, err := toGroupRecords(&dnsCnames)
		if err != nil {
			return nil, err
		}

		groupComposition := routedns.DNSGroupRecords{
			A:     *arrayAs,
			AAAA:  *arrayAAAAs,
			CNAME: *arrayCnames,
		}
		group := routedns.DnsRouteGroupOK{
			GroupID:      groupID,
			FixedGroupID: fixedGroupID,
			ZoneID:       zoneID,
			FixedZoneID:  fixedZoneID,
			DnsRouteGroup: routedns.DnsRouteGroup{
				Name:             name,
				GroupTypeID:      routedns.PrimaryZone,
				GroupProductType: groupProductTypeID,
				GroupComposition: groupComposition,
			},
		}
		groups = append(groups, group)
	}

	return &groups, nil
}

func toCreateDNSGroups(groups *[]routedns.DnsRouteGroupOK) *[]routedns.DnsRouteGroup {
	newgroups := make([]routedns.DnsRouteGroup, 0)
	for _, group := range *groups {
		newgroup := routedns.DnsRouteGroup{
			Name:             group.Name,
			GroupTypeID:      group.GroupTypeID,
			GroupProductType: group.GroupProductType,
			GroupComposition: group.GroupComposition,
		}
		newgroups = append(newgroups, newgroup)
	}

	return &newgroups
}

func toGroupRecords(input *[]interface{}) (*[]routedns.DNSGroupRecord, error) {
	records := make([]routedns.DNSGroupRecord, 0)

	for _, item := range *input {
		curr := item.(map[string]interface{})
		var healthCheck *routedns.HealthCheck = nil
		var err error
		if val, ok := curr["health_check"]; ok {

			hc := val.([]interface{})
			healthCheck, err = toHealthCheck(&hc)
			if err != nil {
				return nil, err
			}
			if healthCheck == nil || (routedns.HealthCheck{}) == *healthCheck {
				healthCheck = nil
			}
		}

		weight := curr["weight"].(int)

		var dnsRecord *routedns.DNSRecord = nil

		if val, ok := curr["record"]; ok {
			rc := val.([]interface{})
			dnsRecord, err = toDNSRouteZoneRecord(&rc, weight)
			if err != nil {
				return nil, err
			}
			if (routedns.DNSRecord{}) == *dnsRecord {
				dnsRecord = nil
			}
		}

		record := routedns.DNSGroupRecord{
			HealthCheck: healthCheck,
			Record:      *dnsRecord,
		}

		records = append(records, record)
	}
	return &records, nil
}

func toHealthCheck(input *[]interface{}) (*routedns.HealthCheck, error) {
	if *input != nil && len(*input) > 0 {
		for _, element := range *input {
			item := element.(map[string]interface{})
			healthCheck := routedns.HealthCheck{}
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

func toDNSRouteZoneRecord(items *[]interface{}, weight int) (*routedns.DNSRecord, error) {
	if *items != nil && len(*items) > 0 {
		for _, element := range *items {
			item := element.(map[string]interface{})
			record := routedns.DNSRecord{
				FixedRecordID:  item["fixed_record_id"].(int),
				FixedGroupID:   item["fixed_group_id"].(int),
				RecordID:       item["record_id"].(int),
				Weight:         weight,
				RecordTypeID:   routedns.RecordType(item["record_type_id"].(int)),
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

//end 1.____________________________________________________________________________________

//2. flatten[func]s are used to save Zone State from API READ API reponse
func flattenDnsRecords(recordItems *[]routedns.DNSRecord) []interface{} {
	if *recordItems != nil && len(*recordItems) > 0 {
		dnsRecords := make([]interface{}, len(*recordItems), len(*recordItems))

		for i, dns := range *recordItems {
			item := make(map[string]interface{})
			item["fixed_record_id"] = dns.FixedRecordID
			item["fixed_group_id"] = dns.FixedGroupID
			item["weight"] = dns.Weight
			item["record_type_id"] = dns.RecordTypeID
			item["record_type_name"] = dns.RecordTypeName
			item["group_id"] = dns.GroupID
			item["is_delete"] = dns.IsDeleted
			item["name"] = dns.Name
			item["ttl"] = dns.TTL
			item["rdata"] = dns.Rdata
			item["verify_id"] = dns.VerifyID

			dnsRecords[i] = item
		}
		return dnsRecords
	}
	return nil
}

func flattenDnsGroups(groupItems *[]routedns.DnsRouteGroupOK) []interface{} {
	if *groupItems != nil && len(*groupItems) > 0 {
		groupArr := make([]interface{}, len(*groupItems), len(*groupItems))

		for i, group := range *groupItems {
			item := make(map[string]interface{})

			item["group_id"] = group.GroupID
			item["id"] = group.GroupID
			item["fixed_group_id"] = group.FixedGroupID
			item["fixed_zone_id"] = group.FixedZoneID
			item["group_product_type_id"] = group.GroupProductType.String()
			item["group_type_id"] = group.GroupTypeID
			if group.GroupTypeID == routedns.PrimaryZone {
				item["group_type"] = "zone"
			}
			item["name"] = group.Name
			item["zone_id"] = group.ZoneID

			item["a"] = flattenGroupDNSs(&group.GroupComposition.A)
			item["aaaa"] = flattenGroupDNSs(&group.GroupComposition.AAAA)
			item["cname"] = flattenGroupDNSs(&group.GroupComposition.CNAME)

			groupArr[i] = item
		}
		return groupArr
	}
	return nil
}

func flattenGroupDNSs(dnsItems *[]routedns.DNSGroupRecord) []interface{} {
	if *dnsItems != nil && len(*dnsItems) > 0 {

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
	return nil
}

func flattenGroupDnsRecord(dns *routedns.DNSRecord) []interface{} {
	if dns != nil && (*dns != routedns.DNSRecord{}) {
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

func flattenHealthCheck(hc *routedns.HealthCheck) []interface{} {
	if hc != nil && (*hc != routedns.HealthCheck{}) {
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
		healthCheck["timeout"] = hc.Timeout
		healthCheck["uri"] = hc.Uri
		//healthCheck["user_id"] = hc.UserID
		//healthCheck[""] = hc.WhiteListedHc

		record[0] = healthCheck
		return record
	}
	return nil
}

//end 2.____________________________________________________________________________________

/*
	When terraform tracks state of zone, following bug exists:
	Scenario: two failover groups and a loadbalancing group
	Initial state:
	fo1: id=1 A[{firstA[id=1], content{computedid=1, dataFromResource=1}}, {secondA[id=2],content{computedid=2, dataFromResource=2}}]
	fo2: id=2 A[{firstA[id=3], content{computedid=3, dataFromResource=3}}, {secondA[id=4],content{computedid=4, dataFromResource=4}}]
	lb1: id=3 A[{firstA[id=5], content{computedid=5, dataFromResource=5}}, {secondA[id=6],content{computedid=6, dataFromResource=6}}]
	When fo2 group was removed, following unexpected result was recorded in tfstate.
	TFState:
	fo1: id=1 A[{firstA[id=1], content{computedid=1, dataFromResource=1}}, {secondA[id=2],content{computedid=2, dataFromResource=2}}]
	fo2: id=2 A[{firstA[id=3], content{computedid=3, dataFromResource=5}}, {secondA[id=4],content{computedid=4, dataFromResource=6}}]

	Note that even though fo2 was deleted, lb1 was gone but its dataFromResource exists in fo2.

	So,
	1. if the number of groups in the resource file  == the number of groups from API
	   => normal update operation
	2. otherwise,
	   => delete all existing groups in db and create new groups from resource file
*/
func updateIDs(local *routedns.ZoneGetOK, remote *routedns.ZoneGetOK) {
	if local == nil || remote == nil {
		return
	}
	local.Status = remote.Status
	records := applyDnsReordChanges(local.Records.A, remote.Records.A)
	remote.Records.A = records
	records = applyDnsReordChanges(local.Records.AAAA, remote.Records.AAAA)
	remote.Records.AAAA = records
	records = applyDnsReordChanges(local.Records.CNAME, remote.Records.CNAME)
	remote.Records.CNAME = records
	records = applyDnsReordChanges(local.Records.CAA, remote.Records.CAA)
	remote.Records.CAA = records
	records = applyDnsReordChanges(local.Records.DLV, remote.Records.DLV)
	remote.Records.DLV = records
	records = applyDnsReordChanges(local.Records.DNSKEY, remote.Records.DNSKEY)
	remote.Records.DNSKEY = records
	records = applyDnsReordChanges(local.Records.DS, remote.Records.DS)
	remote.Records.DS = records
	records = applyDnsReordChanges(local.Records.MX, remote.Records.MX)
	remote.Records.MX = records
	records = applyDnsReordChanges(local.Records.NS, remote.Records.NS)
	remote.Records.NS = records
	records = applyDnsReordChanges(local.Records.NSEC, remote.Records.NSEC)
	remote.Records.NSEC = records
	records = applyDnsReordChanges(local.Records.NSEC3, remote.Records.NSEC3)
	remote.Records.NSEC3 = records
	records = applyDnsReordChanges(local.Records.NSEC3PARAM, remote.Records.NSEC3PARAM)
	remote.Records.NSEC3PARAM = records
	records = applyDnsReordChanges(local.Records.PTR, remote.Records.PTR)
	remote.Records.PTR = records
	records = applyDnsReordChanges(local.Records.RRSIG, remote.Records.RRSIG)
	remote.Records.RRSIG = records
	records = applyDnsReordChanges(local.Records.SOA, remote.Records.SOA)
	remote.Records.SOA = records
	records = applyDnsReordChanges(local.Records.SPF, remote.Records.SPF)
	remote.Records.SPF = records
	records = applyDnsReordChanges(local.Records.SRV, remote.Records.SRV)
	remote.Records.SRV = records
	records = applyDnsReordChanges(local.Records.TXT, remote.Records.TXT)
	remote.Records.TXT = records

	//GROUP-------------------------------------------------------
	if local.Groups == nil || len(local.Groups) == 0 {
		// if db has records but update one has empty, delete db records
		for i := 0; i < len(remote.Groups); i++ {
			markItAsDeleted(remote.Groups[i].GroupComposition.A)
			markItAsDeleted(remote.Groups[i].GroupComposition.AAAA)
			markItAsDeleted(remote.Groups[i].GroupComposition.CNAME)
		}
	} else if remote.Groups == nil || len(remote.Groups) == 0 {
		// if db doesn't have anything...
		remote.Groups = append(remote.Groups, local.Groups...)
	} else {
		// else both local and remote has groups.
		// 1. if found the same group, copy content from update-ones to db values
		if len(local.Groups) != len(remote.Groups) {
			swapGroupsFromLocalToRemote(remote, local)
		} else {
			for i := 0; i < len(local.Groups); i++ {
				for j := 0; j < len(remote.Groups); j++ {
					if local.Groups[i].FixedGroupID == remote.Groups[j].FixedGroupID {
						copyDnsRouteGroupRecordAllIDs(local.Groups[i].GroupComposition.A, remote.Groups[j].GroupComposition.A)
						copyDnsRouteGroupRecordAllIDs(local.Groups[i].GroupComposition.AAAA, remote.Groups[j].GroupComposition.AAAA)
						copyDnsRouteGroupRecordAllIDs(local.Groups[i].GroupComposition.CNAME, remote.Groups[j].GroupComposition.CNAME)
						break
					}
				}
			}
		}
	}
}

func applyDnsReordChanges(local []routedns.DNSRecord, remote []routedns.DNSRecord) []routedns.DNSRecord {

	// modify
	for i := 0; i < len(remote); i++ {

	}
	for i, a1 := range remote {
		for j, a2 := range local {
			if a1.Name == a2.Name {
				a1 = copyDnsRecordContent(&remote[i], &local[j])
				break
			}
		}
	}

	// delete
	for i, _ := range remote {
		isDeleted := true
		for j, _ := range local {
			if remote[i].Name == local[j].Name {

				isDeleted = false
				remote[i].IsDeleted = false
				break
			}
		}
		if isDeleted {
			remote[i].IsDeleted = true
		}
	}

	// add
	addList := make([]routedns.DNSRecord, len(local))
	i := 0
	for _, a1 := range local {
		isFound := false
		for _, a2 := range remote {
			if a1.Name == a2.Name {
				isFound = true
				break
			}
		}
		if !isFound {
			addList[i] = a1
			i++
		}
	}
	if i > 0 {
		addList = addList[:i]
		remote = append(remote, addList...)
	}
	return remote
}

func markItAsDeleted(dnsArr []routedns.DNSGroupRecord) {
	if dnsArr != nil && len(dnsArr) > 0 {
		for i := 0; i < len(dnsArr); i++ {
			dnsArr[i].Record.IsDeleted = true
		}
	}
}

func swapGroupsFromLocalToRemote(remote *routedns.ZoneGetOK, local *routedns.ZoneGetOK) {
	for i := 0; i < len(remote.Groups); i++ {
		if len(remote.Groups[i].GroupComposition.A) > 0 {
			for j := 0; j < len(remote.Groups[i].GroupComposition.A); j++ {
				remote.Groups[i].GroupComposition.A[j].Record.IsDeleted = true
			}
		}
		if len(remote.Groups[i].GroupComposition.AAAA) > 0 {
			for j := 0; j < len(remote.Groups[i].GroupComposition.A); j++ {
				remote.Groups[i].GroupComposition.AAAA[j].Record.IsDeleted = true
			}
		}
		if len(remote.Groups[i].GroupComposition.CNAME) > 0 {
			for j := 0; j < len(remote.Groups[i].GroupComposition.CNAME); j++ {
				remote.Groups[i].GroupComposition.CNAME[j].Record.IsDeleted = true
			}
		}
	}

	for i := 0; i < len(local.Groups); i++ {
		local.Groups[i].FixedGroupID = 0
		local.Groups[i].FixedZoneID = 0
		local.Groups[i].GroupID = 0
		if len(local.Groups[i].GroupComposition.A) > 0 {
			for j := 0; j < len(local.Groups[i].GroupComposition.A); j++ {
				local.Groups[i].GroupComposition.A[j].Record.FixedGroupID = 0
				local.Groups[i].GroupComposition.A[j].Record.FixedRecordID = 0
				local.Groups[i].GroupComposition.A[j].Record.RecordID = 0
			}
		}
		if len(local.Groups[i].GroupComposition.AAAA) > 0 {
			for j := 0; j < len(local.Groups[i].GroupComposition.A); j++ {
				local.Groups[i].GroupComposition.AAAA[j].Record.FixedGroupID = 0
				local.Groups[i].GroupComposition.AAAA[j].Record.FixedRecordID = 0
				local.Groups[i].GroupComposition.AAAA[j].Record.RecordID = 0
			}
		}
		if len(local.Groups[i].GroupComposition.CNAME) > 0 {
			for j := 0; j < len(local.Groups[i].GroupComposition.CNAME); j++ {
				local.Groups[i].GroupComposition.CNAME[j].Record.FixedGroupID = 0
				local.Groups[i].GroupComposition.CNAME[j].Record.FixedRecordID = 0
				local.Groups[i].GroupComposition.CNAME[j].Record.RecordID = 0
			}
		}
	}
	remote.Groups = append(remote.Groups, local.Groups...)
}

func copyDnsRouteGroupRecordAllIDs(from []routedns.DNSGroupRecord, to []routedns.DNSGroupRecord) []routedns.DNSGroupRecord {
	// modify
	for i := 0; i < len(to); i++ {
		for j := 0; j < len(from); j++ {
			if to[i].Record.RecordID == from[j].Record.RecordID {
				to[i].Record = copyDnsRecordContent(&to[i].Record, &from[j].Record)
				if to[i].HealthCheck != nil && from[j].HealthCheck != nil {
					copyHealthCheckIDs(to[i].HealthCheck, from[j].HealthCheck)
				} else if from[i].HealthCheck != nil {
					to[i].HealthCheck = from[j].HealthCheck
				}

				break
			}
		}
	}
	// delete
	for i := 0; i < len(to); i++ {
		isDelete := true
		for j := 0; j < len(from); j++ {
			if to[i].Record.RecordID == from[j].Record.RecordID {
				to[i].Record.IsDeleted = false
				isDelete = false
				break
			}
		}
		if isDelete {
			to[i].Record.IsDeleted = true
		}
	}

	// add
	for _, a1 := range from {
		isFound := false
		for _, a2 := range to {
			if a1.Record.RecordID == a2.Record.RecordID {
				isFound = true
				break
			}
		}
		if !isFound {
			to = append(to, a1)
		}
	}

	return to
}

func copyDnsRecordContent(a1 *routedns.DNSRecord, a2 *routedns.DNSRecord) routedns.DNSRecord {
	a1.Weight = a2.Weight
	a1.Rdata = a2.Rdata
	a1.TTL = a2.TTL
	a1.Name = a2.Name
	return *a1
}

func copyHealthCheckIDs(hc1 *routedns.HealthCheck, hc2 *routedns.HealthCheck) {
	if hc1 == nil {
		hc1 = hc2
	} else if hc1 != nil && hc2 != nil {
		hc1.PortNumber = hc2.PortNumber
		hc1.Timeout = hc2.Timeout
	}
}
