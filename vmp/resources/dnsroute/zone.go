// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package dnsroute

import (
	"context"
	"log"
	"strconv"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DNS Master Server Group
func ResourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceZoneCreate,
		ReadContext:   ResourceZoneRead,
		UpdateContext: ResourceZoneUpdate,
		DeleteContext: ResourceZoneDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account Number for the customer if not already specified in the provider configuration."},
			"zone_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Zone Type."},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain Name."},
			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Status."},
			"is_customer_owned": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "is customer owned"},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment"},
			"record_a": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_aaaa": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_cname": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_mx": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_ns": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_ptr": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_soa": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_spf": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_srv": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_txt": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_dnskey": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_rrsig": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_ds": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_nsec": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_nsec3": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_nsec3param": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_dlv": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
			"record_caa": {
				Type:     schema.TypeList,
				Computed: false,
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
							Type:     schema.TypeString,
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
					},
				},
				Optional:    true,
				Description: "List of A records",
			},
		},
	}
}

func ResourceZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	is_customer_owned := d.Get("is_customer_owned").(int)
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

	zoneRequest := api.ZoneRequest{
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

	log.Printf("[INFO] Creating a new Zone for Account '%s': %+v", accountNumber, zoneRequest)

	dnsrouteClient := api.NewDNSRouteAPIClient(*config)

	newZoneID, err := dnsrouteClient.AddZone(&zoneRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Create successful - New Zone ID: %d", newZoneID)
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

	log.Printf("[INFO] Retrieving Zone by GroupID: %d", zoneID)

	resp, err := dnsRouteClient.GetZone(zoneID)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("zone", resp)
	return diags
}

func ResourceZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	domainName := d.Get("domain_name").(string)
	status := d.Get("status").(int)
	zoneType := d.Get("zone_type").(int)
	is_customer_owned := d.Get("is_customer_owned").(int)
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

	zoneRequest := api.ZoneRequest{
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

	log.Printf("[INFO] Updating a Zone, %s for Account '%s': %+v", zoneRequest.DomainName, accountNumber, zoneRequest)

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

func flattenZoneData(mstGroup api.MasterServerGroupResponse) []interface{} {
	if mstGroup.Masters != nil {
		svgs := make([]interface{}, len(mstGroup.Masters), len(mstGroup.Masters))

		for i, masterServer := range mstGroup.Masters {
			oi := make(map[string]interface{})

			oi["id"] = masterServer.ID
			oi["name"] = masterServer.Name
			oi["ipaddress"] = masterServer.IPAddress

			svgs[i] = oi
		}

		return svgs
	}

	return make([]interface{}, 0)
}

func toDNSRecords(recodeType string, input []interface{}) *[]api.DNSRecord {
	records := make([]api.DNSRecord, 0)

	for _, item := range input {
		curr := item.(map[string]interface{})

		name := curr["name"].(string)
		ttl := curr["ttl"].(string)
		rdata := curr["rdata"].(string)
		verifyid := curr["verify_id"].(int)

		record := api.DNSRecord{
			Name:     name,
			TTL:      ttl,
			Rdata:    rdata,
			VerifyID: verifyid,
		}

		records = append(records, record)
	}
	return &records
}
