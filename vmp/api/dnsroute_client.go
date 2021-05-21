package api

import (
	"fmt"
	"log"
	"strconv"
	"terraform-provider-vmp/vmp/helper"
)

//DNSRouteAPIClient -
type DNSRouteAPIClient struct {
	Config        *ClientConfig
	BaseAPIClient *BaseClient
}

//MasterServerGroupRequest -
type MasterServerGroupRequest struct {
	Name    string         `json:"Name"`
	Masters []MasterServer `json:"MasterServers"`
}

//MasterServerGroupUpdateRequest -
type MasterServerGroupUpdateRequest struct {
	Name    string                    `json:"Name"`
	ID      int                       `json:"Id"`
	Masters []MasterServerWithGroupID `json:"Masters"`
}

//MasterServerWithGroupID -
type MasterServerWithGroupID struct {
	ID            int    `json:"Id,omitempty"`
	MasterGroupID int    `json:"MasterGroupId"`
	Name          string `json:"Name"`
	IPAddress     string `json:"IPAddress"`
}

// MasterServerGroupResponse -
type MasterServerGroupResponse struct {
	MasterGroupID int            `json:"MasterGroupId"`
	Name          string         `json:"Name"`
	Masters       []MasterServer `json:"Masters"`
}

// MasterServerGroupUpdateResponse -
type MasterServerGroupUpdateResponse struct {
	MasterGroupID int                               `json:"Id"`
	CustomerID    int                               `json:"CustomerId"`
	LastUpdated   string                            `json:"LastUpdated"`
	Name          string                            `json:"Name"`
	Masters       []MasterServerWithGroupIDResponse `json:"Masters"`
}

//MasterServerWithGroupIDResponse -
type MasterServerWithGroupIDResponse struct {
	ID            int    `json:"Id,omitempty"`
	MasterGroupID int    `json:"MasterGroupId"`
	Name          string `json:"Name"`
	IPAddress     string `json:"IPAddress"`
	LastUpdated   string `json:"LastUpdated"`
}

// MasterServer
type MasterServer struct {
	ID        int    `json:"Id,omitempty"`
	Name      string `json:"Name"`
	IPAddress string `json:"IPAddress"`
}

//ZoneRequest -
type Zone struct {
	FixedZoneID     int             `json:"FixedZoneId,omitempty"`
	ZoneID          int             `json:"ZoneId,omitempty"`
	DomainName      string          `json:"DomainName,omitempty"`
	Status          int             `json:"Status,omitempty"`
	StatusName      string          `json:"StatusName,omitempty"`
	ZoneType        int             `json:"ZoneType,omitempty"`
	IsCustomerOwned bool            `json:"IsCustomerOwned,omitempty"`
	Comment         string          `json:"Comment,omitempty"`
	Records         DNSRecords      `json:"Records"`
	Serial          int             `json:"Serial,omitempty"`
	Groups          []DnsRouteGroup `json:"groups"`
}

// DNSRecords -
type DNSRecords struct {
	A          []DNSRecord `json:"A,omitempty"`
	AAAA       []DNSRecord `json:"AAAA,omitempty"`
	CName      []DNSRecord `json:"CName,omitempty"`
	MX         []DNSRecord `json:"MX,omitempty"`
	NS         []DNSRecord `json:"NS,omitempty"`
	PTR        []DNSRecord `json:"PTR,omitempty"`
	SOA        []DNSRecord `json:"SOA,omitempty"`
	SPF        []DNSRecord `json:"SPF,omitempty"`
	SRV        []DNSRecord `json:"SRV,omitempty"`
	TXT        []DNSRecord `json:"TXT,omitempty"`
	DNSKEY     []DNSRecord `json:"DNSKEY,omitempty"`
	RRSIG      []DNSRecord `json:"RRSIG,omitempty"`
	DS         []DNSRecord `json:"DS,omitempty"`
	NSEC       []DNSRecord `json:"NSEC,omitempty"`
	NSEC3      []DNSRecord `json:"NSEC3,omitempty"`
	NSEC3PARAM []DNSRecord `json:"NSEC3PARAM,omitempty"`
	DLV        []DNSRecord `json:"DLV,omitempty"`
	CAA        []DNSRecord `json:"CAA,omitempty"`
}

// DNSRecord -
type DNSRecord struct {
	RecordID       int    `json:"RecordId,omitempty"`
	FixedRecordID  int    `json:"FixedRecordId,omitempty"`
	FixedGroupID   int    `json:"FixedGroupId,omitempty"`
	GroupID        int    `json:"GroupId,omitempty"`
	IsDeleted      bool   `json:"IsDelete,omitempty"`
	Name           string `json:"Name,omitempty"`
	TTL            int    `json:"TTL,omitempty"`
	Rdata          string `json:"Rdata,omitempty"`
	VerifyID       int    `json:"VerifyId,omitemtpy"`
	Weight         int    `json:"Weight,omitempty"`
	RecordTypeID   int    `json:"RecordTypeID,omitempty"`
	RecordTypeName string `json:"RecordTypeName,omitempty"`
}

// const for GroupProductTypeId
const (
	GroupProductType_LoadBalancing = iota + 1
	GroupProductType_Failover
	GroupProductType_NoGroup
)

// const for GroupTypeId
const (
	GroupType_CName = iota + 1
	GroupType_SubDomain
	GroupType_Zone
)

type DnsRouteGroup struct {
	ID                 string          `json:"Id,omitempty"`
	GroupID            int             `json:"GroupId,omitempty"`
	FixedGroupID       int             `json:"FixedGroupId,omitempty"`
	Name               string          `json:"Name,omitempty"`
	GroupTypeID        int             `json:"GroupTypeId,omitempty"`
	ZoneId             int             `json:"ZoneId,omitempty"`
	FixedZoneID        int             `json:"FixedZoneId,omitempty"`
	GroupProductTypeID int             `json:"GroupProductTypeId,omitempty"`
	GroupComposition   DNSGroupRecords `json:"GroupComposition,omitempty"`
}

// DNSGroupRecords -
type DNSGroupRecords struct {
	A          []DnsRouteGroupRecord `json:"A,omitempty"`
	AAAA       []DnsRouteGroupRecord `json:"AAAA,omitempty"`
	CName      []DnsRouteGroupRecord `json:"CName,omitempty"`
	MX         []DnsRouteGroupRecord `json:"MX,omitempty"`
	NS         []DnsRouteGroupRecord `json:"NS,omitempty"`
	PTR        []DnsRouteGroupRecord `json:"PTR,omitempty"`
	SOA        []DnsRouteGroupRecord `json:"SOA,omitempty"`
	SPF        []DnsRouteGroupRecord `json:"SPF,omitempty"`
	SRV        []DnsRouteGroupRecord `json:"SRV,omitempty"`
	TXT        []DnsRouteGroupRecord `json:"TXT,omitempty"`
	DNSKEY     []DnsRouteGroupRecord `json:"DNSKEY,omitempty"`
	RRSIG      []DnsRouteGroupRecord `json:"RRSIG,omitempty"`
	DS         []DnsRouteGroupRecord `json:"DS,omitempty"`
	NSEC       []DnsRouteGroupRecord `json:"NSEC,omitempty"`
	NSEC3      []DnsRouteGroupRecord `json:"NSEC3,omitempty"`
	NSEC3PARAM []DnsRouteGroupRecord `json:"NSEC3PARAM,omitempty"`
	DLV        []DnsRouteGroupRecord `json:"DLV,omitempty"`
	CAA        []DnsRouteGroupRecord `json:"CAA,omitempty"`
}

// DNSGroupRecord -
type DnsRouteGroupRecord struct {
	ID          string       `json:"Id,omitempty"`
	Record      DNSRecord    `json:"Record,omitempty"`
	HealthCheck *HealthCheck `json:"HealthCheck"`
	Weight      int          `json:"Weight"`
}

// RecordTypeID
const (
	RecordType_A = iota + 1
	RecordType_AAAA
	RecordType_CNAME
	RecordType_MX
	RecordType_NS
	RecordType_PTR
	RecordType_SOA
	RecordType_SPF
	RecordType_SRV
	RecordType_TXT
	RecordType_DNSKEY
	RecordType_RRSIG
	RecordType_DS
	RecordType_NSEC
	RecordType_NSEC3
	RecordType_NSEC3PARAM
	RecordType_DLV
	RecordType_CAA
)

// HealthCheck -
type HealthCheck struct {
	ID                       int    `json:"Id,omitempty"`
	FixedID                  int    `json:"FixedId,omitempty"`
	CheckInterval            int    `json:"CheckInterval,omitempty"`
	CheckTypeID              int    `json:"CheckTypeId,omitempty"`
	ContentVerification      string `json:"ContentVerification,omitempty"`
	EmailNotificationAddress string `json:"EmailNotificationAddress,omitempty"`
	FailedCheckThreshold     int    `json:"FailedCheckThreshold,omitempty"`
	HTTPMethodID             int    `json:"HTTPMethodId,omitempty"`
	RecordID                 int    `json:"RecordId,omitempty"`
	FixedRecordID            int    `json:"FixedRecordId,omitempty"`
	GroupID                  int    `json:"GroupId,omitempty"`
	FixedGroupID             int    `json:"GroupFixedId,omitempty"`
	IPAddress                string `json:"IPAddress,omitempty"`
	IPVersion                int    `json:"IPVersion,omitempty"`
	PortNumber               string `json:"PortNumber,omitempty"`
	ReintegrationMethodID    int    `json:"ReintegrationMethodId,omitempty"`
	Status                   int    `json:"Status,omitempty"`
	UserID                   int    `json:"UserId,omitempty"`
	TimeOut                  int    `json:"Timeout,omitempty"`
	StatusName               string `json:"StatusName,omitempty"`
	Uri                      string `json:"Uri,omitempty"`
	WhiteListedHc            int    `json:"WhiteListedHc,omitempty"`
}

// LiteralResponse -
type LiteralResponse struct {
	Value interface{}
}

// FOGroup - child element of FailoverGroups in a zone
type FOGroup struct {
	Group FOGWrapper `json:"Group,omitempty"`
}

// FOGWrapper - child of FailoverGroup
type FOGWrapper struct {
	GroupTypeID int              `json:"GroupTypeId,omitempty"`
	Name        string           `json:"Name,omitempty"`
	A           []FailoverRecord `json:"A,omitempty"`
	AAAA        []FailoverRecord `json:"AAAA,omitempty"`
	CName       []FailoverRecord `json:"CNAME,omitempty"`
}

// FailoverRecord -
type FailoverRecord struct {
	HealthCheck HealthCheck `json:"HealthCheck,omitempty"`
	IsPrimary   bool        `json:"IsPrimary,omitempty"`
	Record      DNSRecord   `json:"Record,omitempty"`
}

type LBGroup struct {
	Group LBGroupWrapper `json:"Group,omitempty"`
}

type LBGroupWrapper struct {
	GroupTypeID int        `json:"GroupTypeId,omitempty"`
	Name        string     `json:"Name,omitempty"`
	A           []LBRecord `json:"A,omitempty"`
	AAAA        []LBRecord `json:"AAAA,omitempty"`
	CName       []LBRecord `json:"CNAME,omitempty"`
}

type LBRecord struct {
	HealthCheck HealthCheck `json:"HealthCheck,omitempty"`
	Weight      int         `json:"Weight,omitempty"`
	Record      DNSRecord   `json:"Record,omitempty"`
}

//NewRDNSRouteAPIClient -
func NewDNSRouteAPIClient(config *ClientConfig) *DNSRouteAPIClient {
	APIClient := &DNSRouteAPIClient{
		Config:        config,
		BaseAPIClient: config.BaseClient,
	}

	return APIClient
}

// GetMasterServerGroup -
func (c *DNSRouteAPIClient) GetMasterServerGroup(id int) ([]*MasterServerGroupResponse, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/mastergroups?id=%d", c.Config.AccountNumber, id)
	log.Printf("apiURL:%s", apiURL)
	request, err := c.BaseAPIClient.BuildRequest("GET", apiURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetMasterServerGroup: %v", err)
	}

	parsedResponse := []*MasterServerGroupResponse{}
	log.Printf("dnsroute_client>>GetMasterServerGroup>>parsedResponse:%v", parsedResponse)

	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetOrigin: %v", err)
	}

	return parsedResponse, nil
}

// AddMasterServerGroup -
func (c *DNSRouteAPIClient) AddMasterServerGroup(svg *MasterServerGroupRequest) ([]*MasterServerGroupResponse, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/mastergroup", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, svg, false)
	if err != nil {
		return nil, fmt.Errorf("api>>dnsroute_client>>AddMasterServerGroup: %v", err)
	}

	parsedResponse := []*MasterServerGroupResponse{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("dnsroute_client>>AddMasterServerGroup->API Response Error: %v", err)
	}

	return parsedResponse, nil
}

// UpdateMasterServerGroup -
func (c *DNSRouteAPIClient) UpdateMasterServerGroup(svg *MasterServerGroupUpdateRequest) error {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/mastergroup", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("PUT", apiURL, svg, false)
	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateMasterServerGroup: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("dnsroute_client>>UpdateMasterServerGroup->API Response Error: %v", err)
	}

	return nil
}

// DeleteMasterServerGroup -
func (c *DNSRouteAPIClient) DeleteMasterServerGroup(msgID int) error {
	// TODO: support custom ids for accounts
	apiURL := fmt.Sprintf("v2/mcc/customers/%s/dns/mastergroup/%d", c.Config.AccountNumber, msgID)

	request, err := c.BaseAPIClient.BuildRequest("DELETE", apiURL, nil, false)

	if err != nil {
		return fmt.Errorf("DeleteMasterServerGroupr: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteMasterServerGroup: %v", err)
	}

	return nil
}

// GetZone - Get Zone information of the provided ZoneID which include all dns records, failover servers, and loadbalancing servers if any exists.
func (c *DNSRouteAPIClient) GetZone(id int) (*Zone, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/zone/%d", c.Config.AccountNumber, id)
	log.Printf("apiURL:%s", apiURL)
	request, err := c.BaseAPIClient.BuildRequest("GET", apiURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetZone: %v", err)
	}

	parsedResponse := Zone{}

	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetZone: %v", err)
	}
	log.Printf("dnsroute_client>>GetZone>>parsedResponse:%v", parsedResponse)

	return &parsedResponse, nil
}

// AddZone -
func (c *DNSRouteAPIClient) AddZone(zone *Zone) (int, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/zone", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, *zone, false)
	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddZone->BuildRequest: %v", err)
	}

	resp, err := c.BaseAPIClient.SendRequestWithStringResponse(request)

	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddZone->SendRequestWithStringResponse: %v", err)
	}

	zoneID, err := strconv.Atoi(*resp)
	if err != nil {
		return -1, fmt.Errorf("dnsroute_client>>AddZone->API Response Error: %v", err)
	}
	return zoneID, nil
}

// UpdateZone -
func (c *DNSRouteAPIClient) UpdateZone(zone *Zone) error {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/zone", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, *zone, false)
	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateZone: %v", err)
	}
	_, err = c.BaseAPIClient.SendRequestWithStringResponse(request)

	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateZone: %v", err)
	}

	return nil
}

// DeleteZone -
func (c *DNSRouteAPIClient) DeleteZone(zoneID int) error {
	// TODO: support custom ids for accounts
	apiURL := fmt.Sprintf("v2/mcc/customers/%s/dns/routezone/%d", c.Config.AccountNumber, zoneID)

	request, err := c.BaseAPIClient.BuildRequest("DELETE", apiURL, nil, false)

	if err != nil {
		return fmt.Errorf("DeleteZone: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteZone: %v", err)
	}

	return nil
}

// GetGroup - Get Group information of the provided groupID.
// groupID is a groupID not FixedGroupID
func (c *DNSRouteAPIClient) GetGroup(groupID int, groupProductType string) (*DnsRouteGroup, error) {

	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/group?id=%d&groupType=%s", c.Config.AccountNumber, groupID, groupProductType)
	log.Printf("apiURL:%s", apiURL)
	request, err := c.BaseAPIClient.BuildRequest("GET", apiURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetZone: %v", err)
	}

	parsedResponse := DnsRouteGroup{}

	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)
	helper.LogInstanceToPrettyJson("GET: "+apiURL, parsedResponse, "GetGroup.log")
	if err != nil {
		return nil, fmt.Errorf("GetGroup: %v", err)
	}
	log.Printf("dnsroute_client>>GetGroup>>parsedResponse:%v", parsedResponse)

	return &parsedResponse, nil
}

// AddGroup -
func (c *DNSRouteAPIClient) AddGroup(group *DnsRouteGroup) (int, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/group", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, *group, false)
	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddGroup->BuildRequest: %v", err)
	}

	resp, err := c.BaseAPIClient.SendRequestWithStringResponse(request)

	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddGroup->SendRequestWithStringResponse: %v", err)
	}

	groupID, err := strconv.Atoi(*resp)
	if err != nil {
		return -1, fmt.Errorf("dnsroute_client>>AddGroup->API Response Error: %v", err)
	}
	return groupID, nil
}

// UpdateGroup -
func (c *DNSRouteAPIClient) UpdateGroup(group *DnsRouteGroup) error {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/group", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, *group, false)
	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateGroup: %v", err)
	}
	_, err = c.BaseAPIClient.SendRequestWithStringResponse(request)

	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateGroup: %v", err)
	}

	return nil
}

// DeleteGroup -
func (c *DNSRouteAPIClient) DeleteGroup(groupID int, groupType string) error {
	// TODO: support custom ids for accounts
	apiURL := fmt.Sprintf("v2/mcc/customers/%s/dns/group?id=%d&groupType=%s", c.Config.AccountNumber, groupID, groupType)

	request, err := c.BaseAPIClient.BuildRequest("DELETE", apiURL, nil, false)

	if err != nil {
		return fmt.Errorf("DeleteGroup: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteGroup: %v", err)
	}

	return nil
}
