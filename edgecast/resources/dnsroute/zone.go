// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Master Server Group
func ResourceZone() *schema.Resource {
	groupRecordSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"health_check": {
				// We use a 1-item TypeSet as a workaround since TypeMap
				// doesn't support schema.Resource as a child element type (yet)
				Type:        schema.TypeSet,
				MaxItems:    1,
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
							Description: `Defines the type of health check by 
							its system-defined ID. The following values are 
							supported: 1 - HTTP | 2 - HTTPS | 3 - TCP Open | 
							4 - TCP SSL. Please refer to the following URL for 
							additional information: 
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
							a health check agent will indicate a change in 
							status.`,
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
							Description: `Defines the IP address (IPv4 or IPv6) 
							to which TCP health checks will be directed. IP 
							address is required when check_type_id is 3 or 4`,
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
							unhealthy server/hostname will be integrated back 
							into a group. Supported values are: 1 - Automatic | 
							2 - Manual`,
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
							Description: `Defines the URI to which HTTP/HTTPs 
							health checks will be directed.`,
						},
					},
				},
			},
			"weight": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.`,
			},
			"record": {
				Type: schema.TypeList,
				Description: `Defines a DNS record that will be associated with 
				the zone.`,
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
							Description: `Defines a record's name.`,
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
							Description: `Identifies a zone by its 
							system-defined ID.`,
						},
						"record_type_id": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `Indicates the system-defined ID 
							assigned to the record type.`,
						},
						"record_type_name": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `Indicates the name of the record 
							type.`,
						},
						"is_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Reserved for future use.`,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `Defines a record's weight. Used to 
							denote preference for a load balancing or failover 
							group.`,
						},
					},
				},
				Required: true,
			},
		},
	}
	singleRecordSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"record_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Identifies a DNS Record by its system-defined ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Defines a record's name.`,
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
			},
			"fixed_record_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Reserved for future use.`,
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
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.`,
			},
			"record_type_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `Indicates the system-defined ID assigned to the 
				record type.`,
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
		},
	}

	return &schema.Resource{
		CreateContext: ResourceZoneCreate,
		ReadContext:   ResourceZoneRead,
		UpdateContext: ResourceZoneUpdate,
		DeleteContext: ResourceZoneDelete,
		Importer:      helper.Import(ResourceZoneRead, "account_number", "id"),

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Account Number associated with the customer whose 
				resources you wish to manage. This account number may be found 
				in the upper right-hand corner of the MCC.`,
			},
			"zone_type": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Indicates that a primary zone will be created. Set 
				this request parameter to "1".`,
			},
			"zone_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Reserved for future use.`,
			},
			"fixed_zone_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Identifies a zone by its system-defined ID.",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates a zone's name.",
			},
			"status": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Indicates a zone's status by its system-defined 
				ID. Valid Values: 1 - Active | 2 - Inactive`,
			},
			"status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates a zone's status by its name.",
			},
			"is_customer_owned": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `This parameter is reserved for future use. The 
				only supported value for this parameter is "true."`,
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates the comment associated with a zone.",
			},

			"record_a": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        singleRecordSchema,
				Description: "List of A records",
			},
			"record_aaaa": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of AAAA records",
			},
			"record_cname": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of CNAME records",
			},
			"record_mx": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of MX records",
			},
			"record_ns": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of NS records",
			},
			"record_ptr": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of PTR records",
			},
			"record_soa": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of SOA records",
			},
			"record_spf": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of SPF records",
			},
			"record_srv": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of SRV records",
			},
			"record_txt": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of TXT records",
			},
			"record_dnskey": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of DNSKEY records",
			},
			"record_rrsig": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of RRSIG records",
			},
			"record_ds": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of DS records",
			},
			"record_nsec": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of NSEC records",
			},
			"record_nsec3": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of NSEC3 records",
			},
			"record_nsec3param": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of NSEC3PARAM records",
			},
			"record_dlv": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of DLV records",
			},
			"record_caa": {
				Type:        schema.TypeSet,
				Elem:        singleRecordSchema,
				Optional:    true,
				Description: "List of CAA records",
			},
			"dnsroute_group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `Identifies the group by its 
							system-defined ID.`,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `Identifies the group by its 
							system-defined ID.`,
						},
						"fixed_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Reserved for future use.`,
						},
						"group_type": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Defines the group type. Valid values 
							are: zone`,
						},
						"group_type_id": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `Defines the group type by its 
							system-defined ID`,
						},
						"group_product_type": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Defines the group product type. Valid 
							values are: loadbalancing | failover`,
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
							Description: `Defines the name of the failover or 
							load balancing group.`,
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
							Type:     schema.TypeList,
							Optional: true,
							Elem:     groupRecordSchema,
							Description: `Defines a set of A records associated 
							with this group.`,
						},
						"aaaa": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     groupRecordSchema,
							Description: `Defines a set of AAAA records 
							associated with this group.`,
						},
						"cname": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     groupRecordSchema,
							Description: `Defines a set of CNAME records 
							associated with this group.`,
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
	config := m.(internal.ProviderConfig)
	routeDNSService, err := buildRouteDNSService(config)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// Retrieve Zone data from Resource file
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	comment := d.Get("comment").(string)
	recordsA := d.Get("record_a").(*schema.Set).List()
	recordsAAAA := d.Get("record_aaaa").(*schema.Set).List()
	recordsCNAME := d.Get("record_cname").(*schema.Set).List()
	recordsMX := d.Get("record_mx").(*schema.Set).List()
	recordsNS := d.Get("record_ns").(*schema.Set).List()
	recordsPTR := d.Get("record_ptr").(*schema.Set).List()
	recordsSOA := d.Get("record_soa").(*schema.Set).List()
	recordsSPF := d.Get("record_spf").(*schema.Set).List()
	recordsSRV := d.Get("record_srv").(*schema.Set).List()
	recordsTXT := d.Get("record_txt").(*schema.Set).List()
	recordsDNSKEY := d.Get("record_dnskey").(*schema.Set).List()
	recordsRRSIG := d.Get("record_rrsig").(*schema.Set).List()
	recordsDS := d.Get("record_ds").(*schema.Set).List()
	recordsNSEC := d.Get("record_nsec").(*schema.Set).List()
	recordsNSEC3 := d.Get("record_nsec3").(*schema.Set).List()
	recordsNSEC3PARAM := d.Get("record_nsec3param").(*schema.Set).List()
	recordsDLV := d.Get("record_dlv").(*schema.Set).List()
	recordsCAA := d.Get("record_caa").(*schema.Set).List()
	dnsGroups := d.Get("dnsroute_group").(*schema.Set).List()

	// Build Zone Object
	zoneObj := routedns.Zone{
		DomainName: domainName,
		Status:     status,
		ZoneType:   zoneType,
		Comment:    comment,
	}

	zoneObj.Records.A = expandDNSRecords("A", &recordsA, nil)
	zoneObj.Records.AAAA = expandDNSRecords("AAAA", &recordsAAAA, nil)
	zoneObj.Records.CNAME = expandDNSRecords("Cname", &recordsCNAME, nil)
	zoneObj.Records.MX = expandDNSRecords("MX", &recordsMX, nil)
	zoneObj.Records.NS = expandDNSRecords("NS", &recordsNS, nil)
	zoneObj.Records.PTR = expandDNSRecords("PTR", &recordsPTR, nil)
	zoneObj.Records.SOA = expandDNSRecords("SOA", &recordsSOA, nil)
	zoneObj.Records.SPF = expandDNSRecords("SPF", &recordsSPF, nil)
	zoneObj.Records.SRV = expandDNSRecords("SRV", &recordsSRV, nil)
	zoneObj.Records.TXT = expandDNSRecords("TXT", &recordsTXT, nil)
	zoneObj.Records.DNSKEY = expandDNSRecords("DNSKEY", &recordsDNSKEY, nil)
	zoneObj.Records.RRSIG = expandDNSRecords("RRSIG", &recordsRRSIG, nil)
	zoneObj.Records.DS = expandDNSRecords("DS", &recordsDS, nil)
	zoneObj.Records.NSEC = expandDNSRecords("NSEC", &recordsNSEC, nil)
	zoneObj.Records.NSEC3 = expandDNSRecords("NSEC3", &recordsNSEC3, nil)
	zoneObj.Records.NSEC3PARAM = expandDNSRecords(
		"NSEC3PARAM", &recordsNSEC3PARAM, nil)
	zoneObj.Records.DLV = expandDNSRecords("DLV", &recordsDLV, nil)
	zoneObj.Records.CAA = expandDNSRecords("CAA", &recordsCAA, nil)

	groups, err := expandDNSRouteGroups(&dnsGroups, nil)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if groups != nil {
		zoneObj.Groups = expandCreateDNSGroups(groups)
	}

	// Call add Zone API
	params := routedns.NewAddZoneParams()
	params.AccountNumber = accountNumber
	params.Zone = zoneObj
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
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(internal.ProviderConfig)
	routeDNSService, err := buildRouteDNSService(config)
	if err != nil {
		return diag.FromErr(err)
	}

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
	d.Set("zone_type", zoneObj.ZoneType)
	d.Set("is_customer_owned", zoneObj.IsCustomerOwned)
	d.Set("comment", zoneObj.Comment)
	d.Set("zone_id", zoneID)
	d.Set("status", zoneObj.Status)
	d.Set("status_name", zoneObj.StatusName)

	recordAs := flattenDNSRecords(&zoneObj.Records.A)
	if err := d.Set("record_a", recordAs); err != nil {
		return diag.FromErr(err)
	}
	recordAAAAs := flattenDNSRecords(&zoneObj.Records.AAAA)
	if err := d.Set("record_aaaa", recordAAAAs); err != nil {
		return diag.FromErr(err)
	}
	recordCNames := flattenDNSRecords(&zoneObj.Records.CNAME)
	if err := d.Set("record_cname", recordCNames); err != nil {
		return diag.FromErr(err)
	}
	recordDLVs := flattenDNSRecords(&zoneObj.Records.DLV)
	if err := d.Set("record_dlv", recordDLVs); err != nil {
		return diag.FromErr(err)
	}
	recordDNSKEYs := flattenDNSRecords(&zoneObj.Records.DNSKEY)
	if err := d.Set("record_dnskey", recordDNSKEYs); err != nil {
		return diag.FromErr(err)
	}
	recordDSs := flattenDNSRecords(&zoneObj.Records.DS)
	if err := d.Set("record_ds", recordDSs); err != nil {
		return diag.FromErr(err)
	}
	recordMXs := flattenDNSRecords(&zoneObj.Records.MX)
	if err := d.Set("record_mx", recordMXs); err != nil {
		return diag.FromErr(err)
	}
	recordNSs := flattenDNSRecords(&zoneObj.Records.NS)
	if err := d.Set("record_ns", recordNSs); err != nil {
		return diag.FromErr(err)
	}
	recordNSECs := flattenDNSRecords(&zoneObj.Records.NSEC)
	if err := d.Set("record_nsec", recordNSECs); err != nil {
		return diag.FromErr(err)
	}
	recordNSEC3s := flattenDNSRecords(&zoneObj.Records.NSEC3)
	if err := d.Set("record_nsec3", recordNSEC3s); err != nil {
		return diag.FromErr(err)
	}
	recordNSEC3PARAMs := flattenDNSRecords(&zoneObj.Records.NSEC3PARAM)
	if err := d.Set("record_nsec3param", recordNSEC3PARAMs); err != nil {
		return diag.FromErr(err)
	}
	recordPTRs := flattenDNSRecords(&zoneObj.Records.PTR)
	if err := d.Set("record_ptr", recordPTRs); err != nil {
		return diag.FromErr(err)
	}
	recordRRSIGs := flattenDNSRecords(&zoneObj.Records.RRSIG)
	if err := d.Set("record_rrsig", recordRRSIGs); err != nil {
		return diag.FromErr(err)
	}
	recordSOAs := flattenDNSRecords(&zoneObj.Records.SOA)
	if err := d.Set("record_soa", recordSOAs); err != nil {
		return diag.FromErr(err)
	}
	recordSPFs := flattenDNSRecords(&zoneObj.Records.SPF)
	if err := d.Set("record_spf", recordSPFs); err != nil {
		return diag.FromErr(err)
	}
	recordSRVs := flattenDNSRecords(&zoneObj.Records.SRV)
	if err := d.Set("record_srv", recordSRVs); err != nil {
		return diag.FromErr(err)
	}
	recordTXTs := flattenDNSRecords(&zoneObj.Records.TXT)
	if err := d.Set("record_txt", recordTXTs); err != nil {
		return diag.FromErr(err)
	}
	recordCAAs := flattenDNSRecords(&zoneObj.Records.CAA)
	if err := d.Set("record_caa", recordCAAs); err != nil {
		return diag.FromErr(err)
	}

	dnsGroups := flattenDNSGroups(&zoneObj.Groups)

	if dnsGroups != nil {
		if err := d.Set("dnsroute_group", dnsGroups); err != nil {
			return diag.FromErr(err)
		}
	}

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
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(internal.ProviderConfig)
	routeDNSService, err := buildRouteDNSService(config)
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

	// Build Zone Update Data
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	is_customer_owned := d.Get("is_customer_owned").(bool)
	comment := d.Get("comment").(string)

	// Not using HasChange before GetChange because we need to send the whole
	// Zone for any change, not just the changed records. This should be
	// reworked when we have PATCH ability in the API. findDeletedRecords
	// handles handles empty sets.
	recordsA, deletesA := findDeletedRecords(d.GetChange("record_a"))
	recordsAAAA, deletesAAAA := findDeletedRecords(
		d.GetChange("record_aaaa"))
	recordsCNAME, deletesCNAME := findDeletedRecords(
		d.GetChange("record_cname"))
	recordsMX, deletesMX := findDeletedRecords(d.GetChange("record_mx"))
	recordsNS, deletesNS := findDeletedRecords(d.GetChange("record_ns"))
	recordsPTR, deletesPTR := findDeletedRecords(d.GetChange("record_ptr"))
	recordsSOA, deletesSOA := findDeletedRecords(d.GetChange("record_soa"))
	recordsSPF, deletesSPF := findDeletedRecords(d.GetChange("record_spf"))
	recordsSRV, deletesSRV := findDeletedRecords(d.GetChange("record_srv"))
	recordsTXT, deletesTXT := findDeletedRecords(d.GetChange("record_txt"))
	recordsDNSKEY, deletesDNSKEY := findDeletedRecords(
		d.GetChange("record_dnskey"))
	recordsRRSIG, deletesRRSIG := findDeletedRecords(
		d.GetChange("record_rrsig"))
	recordsDS, deletesDS := findDeletedRecords(d.GetChange("record_ds"))
	recordsNSEC, deletesNSEC := findDeletedRecords(
		d.GetChange("record_nsec"))
	recordsNSEC3, deletesNSEC3 := findDeletedRecords(
		d.GetChange("record_nsec3"))
	recordsNSEC3PARAM, deletesNSEC3PARAM := findDeletedRecords(
		d.GetChange("record_nsec3param"))
	recordsDLV, deletesDLV := findDeletedRecords(d.GetChange("record_dlv"))
	recordsCAA, deletesCAA := findDeletedRecords(d.GetChange("record_caa"))
	recordsGroups, deletesGroups := findDeletedRecords(d.GetChange(
		"dnsroute_group"))
	groups, err := expandDNSRouteGroups(&recordsGroups, deletesGroups)
	if err != nil {
		diag.FromErr(err)
	}
	// Update Zone Object
	zoneObj.DomainName = domainName
	zoneObj.Status = status
	zoneObj.ZoneType = zoneType
	zoneObj.IsCustomerOwned = is_customer_owned
	zoneObj.Comment = comment
	zoneObj.Records.A = expandDNSRecords(
		"A", &recordsA, deletesA)
	zoneObj.Records.AAAA = expandDNSRecords(
		"AAAA", &recordsAAAA, deletesAAAA)
	zoneObj.Records.CNAME = expandDNSRecords(
		"Cname", &recordsCNAME, deletesCNAME)
	zoneObj.Records.MX = expandDNSRecords(
		"MX", &recordsMX, deletesMX)
	zoneObj.Records.NS = expandDNSRecords(
		"NS", &recordsNS, deletesNS)
	zoneObj.Records.PTR = expandDNSRecords(
		"PTR", &recordsPTR, deletesPTR)
	zoneObj.Records.SOA = expandDNSRecords(
		"SOA", &recordsSOA, deletesSOA)
	zoneObj.Records.SPF = expandDNSRecords(
		"SPF", &recordsSPF, deletesSPF)
	zoneObj.Records.SRV = expandDNSRecords(
		"SRV", &recordsSRV, deletesSRV)
	zoneObj.Records.TXT = expandDNSRecords(
		"TXT", &recordsTXT, deletesTXT)
	zoneObj.Records.DNSKEY = expandDNSRecords(
		"DNSKEY", &recordsDNSKEY, deletesDNSKEY)
	zoneObj.Records.RRSIG = expandDNSRecords(
		"RRSIG", &recordsRRSIG, deletesRRSIG)
	zoneObj.Records.DS = expandDNSRecords(
		"DS", &recordsDS, deletesDS)
	zoneObj.Records.NSEC = expandDNSRecords(
		"NSEC", &recordsNSEC, deletesNSEC)
	zoneObj.Records.NSEC3 = expandDNSRecords(
		"NSEC3", &recordsNSEC3, deletesNSEC3)
	zoneObj.Records.NSEC3PARAM = expandDNSRecords(
		"NSEC3PARAM", &recordsNSEC3PARAM, deletesNSEC3PARAM)
	zoneObj.Records.DLV = expandDNSRecords(
		"DLV", &recordsDLV, deletesDLV)
	zoneObj.Records.CAA = expandDNSRecords(
		"CAA", &recordsCAA, deletesCAA)
	zoneObj.Groups = groups

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

func ResourceZoneDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Initialize Route DNS Service
	accountNumber := d.Get("account_number").(string)
	config := m.(internal.ProviderConfig)
	routeDNSService, err := buildRouteDNSService(config)
	if err != nil {
		return diag.FromErr(err)
	}

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

// findDeletedRecords identifies records deleted from the resource file and
// returns two arrays. One array contains new and old (all) records. One array
// contains only the deleted records. Route APIs require that deleted records
// are submitted with is_deleted=true instead of simply removing deleted records
// from the paylod.
func findDeletedRecords(
	old interface{},
	new interface{},
) ([]interface{}, []interface{}) {
	// Represents current resource, state prior to latest Terraform apply
	os := old.(*schema.Set)
	// Repesents desired resource state
	ns := new.(*schema.Set)
	// Set with records only present in the old set (now deleted)
	toDelete := os.Difference(ns).List()
	// All records, new and old, to allow upstream function to process changes
	allRecords := os.Union(ns).List()

	return allRecords, toDelete
}

func expandDNSRecords(
	recodeType string,
	input *[]interface{},
	toDelete []interface{},
) []routedns.DNSRecord {
	records := make([]routedns.DNSRecord, 0)

	for _, item := range *input {
		curr := item.(map[string]interface{})

		name := curr["name"].(string)

		isDeleted := false
		for _, deleteItem := range toDelete {
			deleteMap := deleteItem.(map[string]interface{})
			if name == deleteMap["name"] {
				isDeleted = true
			}
		}

		ttl := curr["ttl"].(int)
		rdata := curr["rdata"].(string)
		verifyID := curr["verify_id"].(int)
		recordID := curr["record_id"].(int)
		fixedRecordID := curr["fixed_record_id"].(int)
		fixedGroupID := curr["fixed_group_id"].(int)
		groupID := curr["group_id"].(int)
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
	return records
}

func expandDNSRouteGroups(
	input *[]interface{},
	toDelete []interface{},
) ([]routedns.DnsRouteGroupOK, error) {
	if *input == nil || len(*input) == 0 {
		return nil, nil
	}
	groups := make([]routedns.DnsRouteGroupOK, 0)

	for _, item := range *input {

		curr := item.(map[string]interface{})

		groupID := curr["group_id"].(int)
		if groupID == 0 {
			groupID = -1
		}

		deleteGroup := false
		for _, delItem := range toDelete {
			delMap := delItem.(map[string]interface{})
			delID := delMap["group_id"].(int)

			if groupID == delID {
				deleteGroup = true
			}
		}

		name := curr["name"].(string)
		groupType := curr["group_type"].(string)
		if groupType != "zone" {
			return nil, fmt.Errorf(
				"invalid group_type: %s. It should be zone", groupType,
			)
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
			return nil, fmt.Errorf(
				`invalid group_product_type: %s. It should be failover or 
				loadbalancing`, groupProductType,
			)
		}

		dnsAs := curr["a"].([]interface{})
		arrayAs, err := expandGroupRecords(&dnsAs, deleteGroup)
		if err != nil {
			return nil, err
		}
		dnsAAAAs := curr["aaaa"].([]interface{})
		arrayAAAAs, err := expandGroupRecords(&dnsAAAAs, deleteGroup)
		if err != nil {
			return nil, err
		}
		dnsCnames := curr["cname"].([]interface{})
		arrayCnames, err := expandGroupRecords(&dnsCnames, deleteGroup)
		if err != nil {
			return nil, err
		}

		groupComposition := routedns.DNSGroupRecords{
			A:     arrayAs,
			AAAA:  arrayAAAAs,
			CNAME: arrayCnames,
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

	return groups, nil
}

func expandCreateDNSGroups(
	groups []routedns.DnsRouteGroupOK,
) []routedns.DnsRouteGroup {
	groupsArr := make([]routedns.DnsRouteGroup, 0)
	for _, group := range groups {
		g := routedns.DnsRouteGroup{
			Name:             group.Name,
			GroupTypeID:      group.GroupTypeID,
			GroupProductType: group.GroupProductType,
			GroupComposition: group.GroupComposition,
		}
		groupsArr = append(groupsArr, g)
	}

	return groupsArr
}

func expandGroupRecords(
	input *[]interface{},
	toDelete bool,
) ([]routedns.DNSGroupRecord, error) {
	records := make([]routedns.DNSGroupRecord, 0)

	for _, item := range *input {
		curr := item.(map[string]interface{})
		var healthCheck *routedns.HealthCheck
		var err error

		if val, ok := curr["health_check"]; ok {
			hc := val.(*schema.Set).List()
			if len(hc) != 0 {
				healthCheck, err = expandHealthCheck(&hc)
				if err != nil {
					return nil, err
				}
			}

			if healthCheck == nil || (routedns.HealthCheck{}) == *healthCheck {
				healthCheck = nil
			}
		}

		weight := curr["weight"].(int)

		var dnsRecord *routedns.DNSRecord = nil

		if val, ok := curr["record"]; ok {
			rc := val.([]interface{})
			dnsRecord, err = expandDNSRouteZoneRecord(&rc, weight, toDelete)
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
	return records, nil
}

func expandHealthCheck(
	healthChecks *[]interface{},
) (*routedns.HealthCheck, error) {
	if len(*healthChecks) > 1 {
		return nil, fmt.Errorf(
			"only one health_check may be defined per record",
		)
	}

	healthCheckObj := routedns.HealthCheck{}
	healthCheckMap := (*healthChecks)[0].(map[string]interface{})

	if healthCheckMap["check_interval"] != nil {
		healthCheckObj.CheckInterval = healthCheckMap["check_interval"].(int)
	}
	if healthCheckMap["content_verification"] != nil {
		healthCheckObj.ContentVerification = healthCheckMap["content_verification"].(string)
	}
	if healthCheckMap["email_notification_address"] != nil {
		healthCheckObj.EmailNotificationAddress = healthCheckMap["email_notification_address"].(string)
	}
	if healthCheckMap["check_type_id"] != nil {
		healthCheckObj.CheckTypeID = healthCheckMap["check_type_id"].(int)
	}
	if healthCheckMap["failed_check_threshold"] != nil {
		healthCheckObj.FailedCheckThreshold = healthCheckMap["failed_check_threshold"].(int)
	}
	if healthCheckMap["http_method_id"] != nil {
		healthCheckObj.HTTPMethodID = healthCheckMap["http_method_id"].(int)
	}
	if healthCheckMap["ip_address"] != nil {
		healthCheckObj.IPAddress = healthCheckMap["ip_address"].(string)
	}
	if healthCheckMap["ip_version"] != nil {
		healthCheckObj.IPVersion = healthCheckMap["ip_version"].(int)
	}
	if healthCheckMap["port_number"] != nil {
		healthCheckObj.PortNumber = healthCheckMap["port_number"].(int)
	}
	if healthCheckMap["reintegration_method_id"] != nil {
		methodID := healthCheckMap["reintegration_method_id"].(int)
		healthCheckObj.ReintegrationMethodID = methodID
	}
	if healthCheckMap["status"] != nil {
		healthCheckObj.Status = healthCheckMap["status"].(int)
	}
	if healthCheckMap["status_name"] != nil {
		healthCheckObj.StatusName = healthCheckMap["status_name"].(string)
	}
	if healthCheckMap["uri"] != nil {
		healthCheckObj.Uri = healthCheckMap["uri"].(string)
	}
	if healthCheckMap["timeout"] != nil {
		healthCheckObj.Timeout = healthCheckMap["timeout"].(int)
	}

	return &healthCheckObj, nil
}

func expandDNSRouteZoneRecord(
	items *[]interface{},
	weight int,
	toDelete bool,
) (*routedns.DNSRecord, error) {
	if *items != nil && len(*items) > 0 {
		for _, element := range *items {
			item := element.(map[string]interface{})
			record := routedns.DNSRecord{
				FixedRecordID: item["fixed_record_id"].(int),
				FixedGroupID:  item["fixed_group_id"].(int),
				RecordID:      item["record_id"].(int),
				Weight:        weight,
				RecordTypeID: routedns.RecordType(
					item["record_type_id"].(int)),
				RecordTypeName: item["record_type_name"].(string),
				GroupID:        item["group_id"].(int),
				IsDeleted:      toDelete,
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

func flattenDNSRecords(recordItems *[]routedns.DNSRecord) []interface{} {
	if *recordItems != nil && len(*recordItems) > 0 {
		dnsRecords := make([]interface{}, len(*recordItems))

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
			item["record_id"] = dns.RecordID

			dnsRecords[i] = item
		}
		return dnsRecords
	}
	return nil
}

func flattenDNSGroups(groupItems *[]routedns.DnsRouteGroupOK) []interface{} {
	if *groupItems != nil && len(*groupItems) > 0 {
		groupArr := make([]interface{}, len(*groupItems))

		for i, group := range *groupItems {
			item := make(map[string]interface{})

			// Setting ID to GroupID as the API only cares about GroupID
			item["id"] = group.GroupID
			item["group_id"] = group.GroupID
			item["fixed_group_id"] = group.FixedGroupID
			item["fixed_zone_id"] = group.FixedZoneID
			item["group_product_type_id"] = group.GroupProductType
			if group.GroupProductType == routedns.Failover {
				item["group_product_type"] = "failover"
			} else if group.GroupProductType == routedns.LoadBalancing {
				item["group_product_type"] = "loadbalancing"
			}
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

		dnsArr := make([]interface{}, len(*dnsItems))

		for i, dns := range *dnsItems {
			item := make(map[string]interface{})

			item["weight"] = dns.Record.Weight
			record := flattenGroupDNSRecord(&dns.Record)

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

func flattenGroupDNSRecord(dns *routedns.DNSRecord) []interface{} {
	if dns != nil && (*dns != routedns.DNSRecord{}) {
		record := make([]interface{}, 1)
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

func flattenHealthCheck(hc *routedns.HealthCheck) interface{} {
	if hc != nil && (*hc != routedns.HealthCheck{}) {
		record := make([]interface{}, 1)
		healthCheck := make(map[string]interface{})

		healthCheck["check_interval"] = hc.CheckInterval
		healthCheck["check_type_id"] = hc.CheckTypeID
		healthCheck["content_verification"] = hc.ContentVerification
		healthCheck["email_notification_address"] = hc.EmailNotificationAddress
		healthCheck["failed_check_threshold"] = hc.FailedCheckThreshold
		healthCheck["fixed_id"] = hc.FixedID
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

		record[0] = healthCheck
		return record
	}
	return nil
}
