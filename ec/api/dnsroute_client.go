// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package api

import (
	"fmt"
	"log"
	"strconv"
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
	Name      string `json:"Name,omitempty"`
	IPAddress string `json:"IPAddress,omitempty"`
}

//Zone -
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

// Tsig Algoirthm Ids
const (
	TsigAlgorithm_HMAC_MD5 = iota + 1
	TsigAlgorithm_HMAC_SHA1
	TsigAlgorithm_HMAC_SHA256
	TsigAlgorithm_HMAC_SHA384
	TsigAlgorithm_HMAC_SHA224
	TsigAlgorithm_HMAC_SHA512
)

// DnsRouteTsig -
type DnsRouteTsig struct {
	ID            int    `json:"Id,omitempty"`
	Alias         string `json:"Alias,omitempty"`
	KeyName       string `json:"KeyName,omitempty"`
	KeyValue      string `json:"KeyValue,omitempty"`
	AlgorithmID   int    `json:"AlgorithmId,omitempty"`
	AlgorithmName string `json:"AlgorithmName,omitempty"`
}

// DnsRouteSecondaryGroupRequest -
type DnsRouteSecondaryZoneGroupRequest struct {
	ID              int                       `json:"Id,omitempty"`
	Name            string                    `json:"Name,omitempty"`
	CustomerID      int                       `json:"CustomerId,omitempty"`
	ZoneComposition DnsZoneCompositionRequest `json:"ZoneComposition,omitempty"`
}

//DnsZoneCompositionRequest -
type DnsZoneCompositionRequest struct {
	MasterGroupID     int                       `json:"MasterGroupId,omitempty"`
	Zones             []SecondaryZoneRequest    `json:"Zones,omitempty"`
	MasterServerTsigs []MasterServerTsigRequest `json:"MasterServerTsigs,omitempty"`
}

//SecondaryZoneRequest -
type SecondaryZoneRequest struct {
	DomainName      string `json:"DomainName,omitempty"`
	Status          int    `json:"Status,omitempty"`
	ZoneType        int    `json:"ZoneType,omitempty"`
	Comment         string `json:"Comment,omitempty"`
	IsCustomerOwned bool   `json:"IsCustomerOwned,omitempty"`
}

//MasterServerTsigRequest -
type MasterServerTsigRequest struct {
	MasterServer MasterServerRequest `json:"MasterServer,omitempty"`
	Tsig         TsigRequest         `json:"Tsig,omitempty"`
}

//MasterServerRequest -
type MasterServerRequest struct {
	ID int `json:"Id,omitempty"`
}

//TsigRequest -
type TsigRequest struct {
	ID int `json:"Id,omitempty"`
}

// SecondaryGroupResponse -
type SecondaryGroupResponse struct {
	ID              int                `json:"Id,omitempty"`
	Name            string             `json:"Name,omitempty"`
	ZoneComposition DnsZoneComposition `json:"ZoneComposition,omitempty"`
}

//DnsZoneComposition -
type DnsZoneComposition struct {
	MasterGroupID     int                `json:"MasterGroupId,omitempty"`
	Zones             []ZoneForSZG       `json:"Zones,omitempty"`
	MasterServerTsigs []MasterServerTsig `json:"MasterServerTsigs,omitempty"`
}

//ZoneForSZG -
type ZoneForSZG struct {
	FixedZoneID     int    `json:"FixedZoneId,omitempty"`
	ZoneID          int    `json:"ZoneId,omitempty"`
	DomainName      string `json:"DomainName,omitempty"`
	Status          int    `json:"Status,omitempty"`
	StatusName      string `json:"StatusName,omitempty"`
	ZoneType        int    `json:"ZoneType,omitempty"`
	IsCustomerOwned bool   `json:"IsCustomerOwned,omitempty"`
	Comment         string `json:"Comment,omitempty"`
}

//MasterServerTsig -
type MasterServerTsig struct {
	MasterServer MasterServerForSZG `json:"MasterServer,omitempty"`
	Tsig         Tsig               `json:"Tsig,omitempty"`
}

// MasterServerForSZG
type MasterServerForSZG struct {
	ID        int    `json:"Id,omitempty"`
	Name      string `json:"Name,omitempty"`
	IPAddress string `json:"IpAddress,omitempty"`
}

//Tsig -
type Tsig struct {
	ID            int    `json:"Id,omitempty"`
	Alias         string `json:"Alias,omitempty"`
	KeyName       string `json:"KeyName,omitempty"`
	KeyValue      string `json:"KeyValue,omitempty"`
	AlgorithmID   int    `json:"AlgorithmId,omitempty"`
	AlgorithmName string `json:"AlgorithmName,omitempty"`
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

	request, err := c.BaseAPIClient.BuildRequest("GET", apiURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetMasterServerGroup: %v", err)
	}

	parsedResponse := []*MasterServerGroupResponse{}

	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("GetMasterServerGroup: %v", err)
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

// GetTsig -
func (c *DNSRouteAPIClient) GetTsig(id int) (*DnsRouteTsig, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/tsigs/%d", c.Config.AccountNumber, id)
	log.Printf("apiURL:%s", apiURL)
	request, err := c.BaseAPIClient.BuildRequest("GET", apiURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetTsig: %v", err)
	}

	parsedResponse := &DnsRouteTsig{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetTsig: %v", err)
	}

	return parsedResponse, nil
}

// AddTsig -
func (c *DNSRouteAPIClient) AddTsig(tsig *DnsRouteTsig) (int, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/tsig", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, tsig, false)
	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddTsig->BuildRequest: %v", err)
	}

	resp, err := c.BaseAPIClient.SendRequestWithStringResponse(request)

	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddTsig->SendRequestWithStringResponse: %v", err)
	}

	tsigID, err := strconv.Atoi(*resp)
	if err != nil {
		return -1, fmt.Errorf("dnsroute_client>>AddTsig->API Response Error: %v", err)
	}
	return tsigID, nil
}

// UpdateTsig -
func (c *DNSRouteAPIClient) UpdateTsig(tsig *DnsRouteTsig) error {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/tsigs/%d", c.Config.AccountNumber, tsig.ID)
	request, err := c.BaseAPIClient.BuildRequest("PUT", apiURL, tsig, false)
	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateTsig: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("dnsroute_client>>UpdateTsig->API Response Error: %v", err)
	}

	return nil
}

// DeleteTsig -
func (c *DNSRouteAPIClient) DeleteTsig(tsigID int) error {
	// TODO: support custom ids for accounts
	apiURL := fmt.Sprintf("v2/mcc/customers/%s/dns/tsigs/%d", c.Config.AccountNumber, tsigID)

	request, err := c.BaseAPIClient.BuildRequest("DELETE", apiURL, nil, false)

	if err != nil {
		return fmt.Errorf("DeleteTsig: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteTsig: %v", err)
	}

	return nil
}

// GetSecondaryZoneGroup -
func (c *DNSRouteAPIClient) GetSecondaryZoneGroup(id int) (*SecondaryGroupResponse, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/secondarygroup?id=%d", c.Config.AccountNumber, id)
	log.Printf("apiURL:%s", apiURL)
	request, err := c.BaseAPIClient.BuildRequest("GET", apiURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetSecondaryZoneGroup: %v", err)
	}

	parsedResponse := []*SecondaryGroupResponse{}
	resp, err := c.BaseAPIClient.SendRequest(request, &parsedResponse)
	log.Printf("GetSecondaryZoneGroup:%v", resp)
	if err != nil {
		return nil, fmt.Errorf("GetSecondaryZoneGroup: %v", err)
	}

	return parsedResponse[0], nil
}

// AddSecondaryZoneGroup -
func (c *DNSRouteAPIClient) AddSecondaryZoneGroup(secondaryGroup *DnsRouteSecondaryZoneGroupRequest) (int, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/secondarygroup", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, *secondaryGroup, false)
	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddSecondaryZone->BuildRequest: %v", err)
	}

	parsedResponse := SecondaryGroupResponse{}
	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddSecondaryZoneGroup->SendRequest: %v", err)
	}

	return parsedResponse.ID, nil
}

// UpdateZone -
func (c *DNSRouteAPIClient) UpdateSecondaryZoneGroup(secondaryGroup *DnsRouteSecondaryZoneGroupRequest) error {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/secondarygroup", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("PUT", apiURL, *secondaryGroup, false)
	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateSecondaryZoneGroup: %v", err)
	}
	_, err = c.BaseAPIClient.SendRequestWithStringResponse(request)

	if err != nil {
		return fmt.Errorf("api>>dnsroute_client>>UpdateSecondaryZoneGroup: %v", err)
	}

	return nil
}

// DeleteZone -
func (c *DNSRouteAPIClient) DeleteSecondaryZoneGroup(zoneID int) error {
	// TODO: support custom ids for accounts
	apiURL := fmt.Sprintf("v2/mcc/customers/%s/dns/secondarygroup?id=%d", c.Config.AccountNumber, zoneID)

	request, err := c.BaseAPIClient.BuildRequest("DELETE", apiURL, nil, false)

	if err != nil {
		return fmt.Errorf("DeleteSecondaryZoneGroup: %v", err)
	}

	_, err = c.BaseAPIClient.SendRequest(request, nil)

	if err != nil {
		return fmt.Errorf("DeleteSecondaryZoneGroup: %v", err)
	}

	return nil
}
