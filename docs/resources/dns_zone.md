---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "edgecast_dns_zone Resource - terraform-provider-edgecast"
subcategory: ""
description: |-
  
---

# edgecast_dns_zone (Resource)
**NOTE: Route DNS feature support via Terraform is currently in Beta status.**

The Managed (Primary) or Secondary DNS module allows the creation and management 
of zones. A zone defines a set of data through which our authoritative name 
servers can provide a response to DNS queries. This data can be found in the 
records associated with your zone. In addition to record administration, zone 
management also allows the definition of load balancing and failover 
configurations for address and CNAME records within that zone.

For more information, please visit the Route Help Center
https://docs.whitecdn.com/dns/index.html#Route/Administration/DNS_Zone_Management.htm

## Example Usage

```terraform
resource "edgecast_dns_zone" "anyl" {
  account_number = "DE0B"
	domain_name = "anyl.com."
  status = 1 # 1: active, 2: inactive
	zone_type = 1 # 1: Primary zone. This value should always be 1.
	is_customer_owned = true # This value should always be true
	comment = "test comment"
	record_a {
    name="mail"
    ttl=3600
    rdata="10.10.10.45"
  }
  record_a {
			name="www"
      ttl=3600
      rdata="10.10.10.114"
  }
  record_a {
			name="news"
      ttl=3600
      rdata="10.10.10.200"
  }
  record_aaaa {
			name="www"
      ttl="3600"
      rdata="10:0:1::0:3"
  }
  record_cname {
			name="www"
      ttl=3600
      rdata="www.cooler.com"
  }
	record_mx {
    name="@"
    ttl=3600
    rdata="10 mail.cooler.com"
  }

  dnsroute_group {
    group_type="zone"
    group_product_type="failover"
    name="fo1"
    a {
      weight=100
      record {
        ttl=300
        rdata="10.10.1.11"
      }
    }
    a {
      weight=0
      record {
        ttl=300
        rdata="10.10.1.12"
      }
    }
  }

  dnsroute_group {
    group_type="zone"
    group_product_type="failover"
    name="fo2"
    a {
      weight=100
      record {
        ttl=300
        rdata="10.10.2.21"
      }
    }
    a {
      weight=0
      record {
        ttl=300
        rdata="10.10.2.22"
      }
    }
  }

  dnsroute_group {
    group_type="zone"
    group_product_type="loadbalancing"
    name="lbg"
    a {
      weight=33
      health_check {
        check_interval=300
        check_type_id=1 # 1: HTTP, 2: HTTPS, 3: TCP Open, 4: TCP SSL
        content_verification="10"
        email_notification_address="notice@glory1.com"
        failed_check_threshold=10
        http_method_id=1 # 1: GET, 2: POST
        # ip_address="" # IP address only required when check_type_id 3,4
        ip_version=1 # 1: IPv4, 2: IPv6
        # port_number=80 # Port only required when check_type_id 3,4
        reintegration_method_id=1 # 1: Automatic, 2: Manual
        uri="www.yahoo.com"
        timeout=100
      }
      record {
        name="lbg1"
        ttl=300
        rdata="10.10.3.1"
      }
    }
    a {
      weight=33
      record {
        name="lbg2"
        ttl=300
        rdata="10.10.3.2"
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_number` (String) Account Number associated with the customer whose 
				resources you wish to manage. This account number may be found 
				in the upper right-hand corner of the MCC.
- `domain_name` (String) Indicates a zone's name.
- `status` (Number) Indicates a zone's status by its system-defined 
				ID. Valid Values: 1 - Active | 2 - Inactive
- `zone_type` (Number) Indicates that a primary zone will be created. Set 
				this request parameter to "1".

### Optional

- `comment` (String) Indicates the comment associated with a zone.
- `dnsroute_group` (Block Set) (see [below for nested schema](#nestedblock--dnsroute_group))
- `is_customer_owned` (Boolean) This parameter is reserved for future use. The 
				only supported value for this parameter is "true."
- `record_a` (Block Set) List of A records (see [below for nested schema](#nestedblock--record_a))
- `record_aaaa` (Block Set) List of AAAA records (see [below for nested schema](#nestedblock--record_aaaa))
- `record_caa` (Block Set) List of CAA records (see [below for nested schema](#nestedblock--record_caa))
- `record_cname` (Block Set) List of CNAME records (see [below for nested schema](#nestedblock--record_cname))
- `record_dlv` (Block Set) List of DLV records (see [below for nested schema](#nestedblock--record_dlv))
- `record_dnskey` (Block Set) List of DNSKEY records (see [below for nested schema](#nestedblock--record_dnskey))
- `record_ds` (Block Set) List of DS records (see [below for nested schema](#nestedblock--record_ds))
- `record_mx` (Block Set) List of MX records (see [below for nested schema](#nestedblock--record_mx))
- `record_ns` (Block Set) List of NS records (see [below for nested schema](#nestedblock--record_ns))
- `record_nsec` (Block Set) List of NSEC records (see [below for nested schema](#nestedblock--record_nsec))
- `record_nsec3` (Block Set) List of NSEC3 records (see [below for nested schema](#nestedblock--record_nsec3))
- `record_nsec3param` (Block Set) List of NSEC3PARAM records (see [below for nested schema](#nestedblock--record_nsec3param))
- `record_ptr` (Block Set) List of PTR records (see [below for nested schema](#nestedblock--record_ptr))
- `record_rrsig` (Block Set) List of RRSIG records (see [below for nested schema](#nestedblock--record_rrsig))
- `record_soa` (Block Set) List of SOA records (see [below for nested schema](#nestedblock--record_soa))
- `record_spf` (Block Set) List of SPF records (see [below for nested schema](#nestedblock--record_spf))
- `record_srv` (Block Set) List of SRV records (see [below for nested schema](#nestedblock--record_srv))
- `record_txt` (Block Set) List of TXT records (see [below for nested schema](#nestedblock--record_txt))

### Read-Only

- `fixed_zone_id` (Number) Identifies a zone by its system-defined ID.
- `id` (String) The ID of this resource.
- `status_name` (String) Indicates a zone's status by its name.
- `zone_id` (Number) Reserved for future use.

<a id="nestedblock--dnsroute_group"></a>
### Nested Schema for `dnsroute_group`

Required:

- `group_product_type` (String) Defines the group product type. Valid 
							values are: loadbalancing | failover
- `group_type` (String) Defines the group type. Valid values 
							are: zone
- `name` (String) Defines the name of the failover or 
							load balancing group.

Optional:

- `a` (Block List) Defines a set of A records associated 
							with this group. (see [below for nested schema](#nestedblock--dnsroute_group--a))
- `aaaa` (Block List) Defines a set of AAAA records 
							associated with this group. (see [below for nested schema](#nestedblock--dnsroute_group--aaaa))
- `cname` (Block List) Defines a set of CNAME records 
							associated with this group. (see [below for nested schema](#nestedblock--dnsroute_group--cname))

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number) Identifies the group by its 
							system-defined ID.
- `group_product_type_id` (Number) Defines the group product type by its 
							system-defined ID
- `group_type_id` (Number) Defines the group type by its 
							system-defined ID
- `id` (Number) Identifies the group by its 
							system-defined ID.
- `zone_id` (Number) Reserved for future use.

<a id="nestedblock--dnsroute_group--a"></a>
### Nested Schema for `dnsroute_group.a`

Required:

- `record` (Block List, Min: 1) Defines a DNS record that will be associated with 
				the zone. (see [below for nested schema](#nestedblock--dnsroute_group--a--record))
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Optional:

- `health_check` (Block Set, Max: 1) Define a record's health check configuration (see [below for nested schema](#nestedblock--dnsroute_group--a--health_check))

<a id="nestedblock--dnsroute_group--a--record"></a>
### Nested Schema for `dnsroute_group.a.record`

Required:

- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `fixed_zone_id` (Number) Identifies a zone by its 
							system-defined ID.
- `is_delete` (Boolean) Reserved for future use.
- `zone_id` (Number) Reserved for future use.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `group_id` (Number) Identifies the group this record is 
							assoicated with by its system-defined ID.
- `name` (String) Defines a record's name.
- `record_id` (Number) Identifies a DNS Record by its 
							system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID 
							assigned to the record type.
- `record_type_name` (String) Indicates the name of the record 
							type.
- `verify_id` (Number) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to 
							denote preference for a load balancing or failover 
							group.


<a id="nestedblock--dnsroute_group--a--health_check"></a>
### Nested Schema for `dnsroute_group.a.health_check`

Required:

- `check_interval` (Number) Defines the number of seconds between 
							health checks.
- `check_type_id` (Number) Defines the type of health check by 
							its system-defined ID. The following values are 
							supported: 1 - HTTP | 2 - HTTPS | 3 - TCP Open | 
							4 - TCP SSL. Please refer to the following URL for 
							additional information: 
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HC_Types.htm
- `email_notification_address` (String) Defines the e-mail address to which 
							health check notifications will be sent.
- `failed_check_threshold` (Number) Defines the number of consecutive 
							times that the same result must be returned before 
							a health check agent will indicate a change in 
							status.
- `reintegration_method_id` (Number) Indicates the method through which an 
							unhealthy server/hostname will be integrated back 
							into a group. Supported values are: 1 - Automatic | 
							2 - Manual

Optional:

- `content_verification` (String) Defines the text that will be used to 
							verify the success of the health check.
- `http_method_id` (Number) Defines an HTTP method by its 
							system-defined ID. An HTTP method is only used by 
							HTTP/HTTPs health checks. Supported values are: 
							1 - GET, 2 - POST. Refer to the following URL for 
							additional information: 
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HTTP_Methods.htm
- `ip_address` (String) Defines the IP address (IPv4 or IPv6) 
							to which TCP health checks will be directed. IP 
							address is required when check_type_id is 3 or 4
- `ip_version` (Number) Defines an IP version by its 
							system-defined ID. This IP version is only used by 
							HTTP/HTTPs health checks. Supported values are: 
							1 - IPv4, 2 - IPv6. Refer to the following URL for 
							additional information:
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_IP_Versions_HC.htm
- `port_number` (Number) Defines the port to which TCP health 
							checks will be directed.
- `timeout` (Number) Reserved for future use.
- `uri` (String) Defines the URI to which HTTP/HTTPs 
							health checks will be directed.

Read-Only:

- `fixed_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `group_id` (Number) Defines the Group ID this health check 
							is associated with.
- `id` (Number) Identifies the health check by its 
							system-defined ID.
- `record_id` (Number) Defines the DNS record ID this health 
							check is associated with.
- `status` (Number) Indicates the server/hostname's health 
							check status by its system-defined ID.
- `status_name` (String) Indicates the server/hostname's health 
							check status.



<a id="nestedblock--dnsroute_group--aaaa"></a>
### Nested Schema for `dnsroute_group.aaaa`

Required:

- `record` (Block List, Min: 1) Defines a DNS record that will be associated with 
				the zone. (see [below for nested schema](#nestedblock--dnsroute_group--aaaa--record))
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Optional:

- `health_check` (Block Set, Max: 1) Define a record's health check configuration (see [below for nested schema](#nestedblock--dnsroute_group--aaaa--health_check))

<a id="nestedblock--dnsroute_group--aaaa--record"></a>
### Nested Schema for `dnsroute_group.aaaa.record`

Required:

- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `fixed_zone_id` (Number) Identifies a zone by its 
							system-defined ID.
- `is_delete` (Boolean) Reserved for future use.
- `zone_id` (Number) Reserved for future use.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `group_id` (Number) Identifies the group this record is 
							assoicated with by its system-defined ID.
- `name` (String) Defines a record's name.
- `record_id` (Number) Identifies a DNS Record by its 
							system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID 
							assigned to the record type.
- `record_type_name` (String) Indicates the name of the record 
							type.
- `verify_id` (Number) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to 
							denote preference for a load balancing or failover 
							group.


<a id="nestedblock--dnsroute_group--aaaa--health_check"></a>
### Nested Schema for `dnsroute_group.aaaa.health_check`

Required:

- `check_interval` (Number) Defines the number of seconds between 
							health checks.
- `check_type_id` (Number) Defines the type of health check by 
							its system-defined ID. The following values are 
							supported: 1 - HTTP | 2 - HTTPS | 3 - TCP Open | 
							4 - TCP SSL. Please refer to the following URL for 
							additional information: 
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HC_Types.htm
- `email_notification_address` (String) Defines the e-mail address to which 
							health check notifications will be sent.
- `failed_check_threshold` (Number) Defines the number of consecutive 
							times that the same result must be returned before 
							a health check agent will indicate a change in 
							status.
- `reintegration_method_id` (Number) Indicates the method through which an 
							unhealthy server/hostname will be integrated back 
							into a group. Supported values are: 1 - Automatic | 
							2 - Manual

Optional:

- `content_verification` (String) Defines the text that will be used to 
							verify the success of the health check.
- `http_method_id` (Number) Defines an HTTP method by its 
							system-defined ID. An HTTP method is only used by 
							HTTP/HTTPs health checks. Supported values are: 
							1 - GET, 2 - POST. Refer to the following URL for 
							additional information: 
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HTTP_Methods.htm
- `ip_address` (String) Defines the IP address (IPv4 or IPv6) 
							to which TCP health checks will be directed. IP 
							address is required when check_type_id is 3 or 4
- `ip_version` (Number) Defines an IP version by its 
							system-defined ID. This IP version is only used by 
							HTTP/HTTPs health checks. Supported values are: 
							1 - IPv4, 2 - IPv6. Refer to the following URL for 
							additional information:
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_IP_Versions_HC.htm
- `port_number` (Number) Defines the port to which TCP health 
							checks will be directed.
- `timeout` (Number) Reserved for future use.
- `uri` (String) Defines the URI to which HTTP/HTTPs 
							health checks will be directed.

Read-Only:

- `fixed_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `group_id` (Number) Defines the Group ID this health check 
							is associated with.
- `id` (Number) Identifies the health check by its 
							system-defined ID.
- `record_id` (Number) Defines the DNS record ID this health 
							check is associated with.
- `status` (Number) Indicates the server/hostname's health 
							check status by its system-defined ID.
- `status_name` (String) Indicates the server/hostname's health 
							check status.



<a id="nestedblock--dnsroute_group--cname"></a>
### Nested Schema for `dnsroute_group.cname`

Required:

- `record` (Block List, Min: 1) Defines a DNS record that will be associated with 
				the zone. (see [below for nested schema](#nestedblock--dnsroute_group--cname--record))
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Optional:

- `health_check` (Block Set, Max: 1) Define a record's health check configuration (see [below for nested schema](#nestedblock--dnsroute_group--cname--health_check))

<a id="nestedblock--dnsroute_group--cname--record"></a>
### Nested Schema for `dnsroute_group.cname.record`

Required:

- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `fixed_zone_id` (Number) Identifies a zone by its 
							system-defined ID.
- `is_delete` (Boolean) Reserved for future use.
- `zone_id` (Number) Reserved for future use.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `group_id` (Number) Identifies the group this record is 
							assoicated with by its system-defined ID.
- `name` (String) Defines a record's name.
- `record_id` (Number) Identifies a DNS Record by its 
							system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID 
							assigned to the record type.
- `record_type_name` (String) Indicates the name of the record 
							type.
- `verify_id` (Number) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to 
							denote preference for a load balancing or failover 
							group.


<a id="nestedblock--dnsroute_group--cname--health_check"></a>
### Nested Schema for `dnsroute_group.cname.health_check`

Required:

- `check_interval` (Number) Defines the number of seconds between 
							health checks.
- `check_type_id` (Number) Defines the type of health check by 
							its system-defined ID. The following values are 
							supported: 1 - HTTP | 2 - HTTPS | 3 - TCP Open | 
							4 - TCP SSL. Please refer to the following URL for 
							additional information: 
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HC_Types.htm
- `email_notification_address` (String) Defines the e-mail address to which 
							health check notifications will be sent.
- `failed_check_threshold` (Number) Defines the number of consecutive 
							times that the same result must be returned before 
							a health check agent will indicate a change in 
							status.
- `reintegration_method_id` (Number) Indicates the method through which an 
							unhealthy server/hostname will be integrated back 
							into a group. Supported values are: 1 - Automatic | 
							2 - Manual

Optional:

- `content_verification` (String) Defines the text that will be used to 
							verify the success of the health check.
- `http_method_id` (Number) Defines an HTTP method by its 
							system-defined ID. An HTTP method is only used by 
							HTTP/HTTPs health checks. Supported values are: 
							1 - GET, 2 - POST. Refer to the following URL for 
							additional information: 
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_HTTP_Methods.htm
- `ip_address` (String) Defines the IP address (IPv4 or IPv6) 
							to which TCP health checks will be directed. IP 
							address is required when check_type_id is 3 or 4
- `ip_version` (Number) Defines an IP version by its 
							system-defined ID. This IP version is only used by 
							HTTP/HTTPs health checks. Supported values are: 
							1 - IPv4, 2 - IPv6. Refer to the following URL for 
							additional information:
							https://developer.edgecast.com/cdn/api/Content/Media_Management/DNS/Get_A_IP_Versions_HC.htm
- `port_number` (Number) Defines the port to which TCP health 
							checks will be directed.
- `timeout` (Number) Reserved for future use.
- `uri` (String) Defines the URI to which HTTP/HTTPs 
							health checks will be directed.

Read-Only:

- `fixed_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `group_id` (Number) Defines the Group ID this health check 
							is associated with.
- `id` (Number) Identifies the health check by its 
							system-defined ID.
- `record_id` (Number) Defines the DNS record ID this health 
							check is associated with.
- `status` (Number) Indicates the server/hostname's health 
							check status by its system-defined ID.
- `status_name` (String) Indicates the server/hostname's health 
							check status.




<a id="nestedblock--record_a"></a>
### Nested Schema for `record_a`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_aaaa"></a>
### Nested Schema for `record_aaaa`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_caa"></a>
### Nested Schema for `record_caa`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_cname"></a>
### Nested Schema for `record_cname`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_dlv"></a>
### Nested Schema for `record_dlv`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_dnskey"></a>
### Nested Schema for `record_dnskey`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_ds"></a>
### Nested Schema for `record_ds`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_mx"></a>
### Nested Schema for `record_mx`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_ns"></a>
### Nested Schema for `record_ns`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_nsec"></a>
### Nested Schema for `record_nsec`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_nsec3"></a>
### Nested Schema for `record_nsec3`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_nsec3param"></a>
### Nested Schema for `record_nsec3param`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_ptr"></a>
### Nested Schema for `record_ptr`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_rrsig"></a>
### Nested Schema for `record_rrsig`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_soa"></a>
### Nested Schema for `record_soa`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_spf"></a>
### Nested Schema for `record_spf`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_srv"></a>
### Nested Schema for `record_srv`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.


<a id="nestedblock--record_txt"></a>
### Nested Schema for `record_txt`

Required:

- `name` (String) Defines a record's name.
- `rdata` (String) Defines a record's value.
- `ttl` (Number) Defines a record's TTL.

Optional:

- `is_delete` (Boolean) Reserved for future use.
- `weight` (Number) Defines a record's weight. Used to denote 
				preference for a load balancing or failover group.

Read-Only:

- `fixed_group_id` (Number) Reserved for future use.
- `fixed_record_id` (Number) Reserved for future use.
- `fixed_zone_id` (Number) Reserved for future use.
- `group_id` (Number)
- `record_id` (Number) Identifies a DNS Record by its system-defined ID.
- `record_type_id` (Number) Indicates the system-defined ID assigned to the 
				record type.
- `record_type_name` (String) Indicates the name of the record type.
- `verify_id` (Number) Reserved for future use.
- `zone_id` (Number) Reserved for future use.




## Import

To import a resource, create a resource block for it in your configuration:

```terraform
resource "edgecast_dns_zone" "example" {

}
```

Now run terraform import to attach an existing instance to the resource configuration:

```shell
terraform import edgecast_dns_zone.example ACCOUNT_NUMBER:ID
```
|                 |                                                                   |
|:----------------|-------------------------------------------------------------------|
| `ACCOUNT_NUMBER`  | The account number the DNS zone ID is associated with. |
| `ID` | The DNS zone ID to import.                                        |

As a result of the above command, the resource is recorded in the state file.

