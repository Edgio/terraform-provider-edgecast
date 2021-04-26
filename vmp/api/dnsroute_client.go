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
	Name      string `json:"Name"`
	IPAddress string `json:"IPAddress"`
}

//ZoneRequest -
type ZoneRequest struct {
	FixedZoneID     int        `json:"FixedZoneId,omitempty"`
	ZoneID          int        `json:"ZoneId,omitempty"`
	DomainName      string     `json:"DomainName,omitempty"`
	Status          int        `json:"Status,omitempty"`
	ZoneType        int        `json:"ZoneType,omitempty"`
	IsCustomerOwned int        `json:"IsCustomerOwned,omitempty"`
	Comment         string     `json:"Comment,omitempty"`
	Records         DNSRecords `json:"Records"`
}

// ZoneResponse -
type ZoneResponse struct {
	FixedZoneID         int            `json:"FixedZoneId,omitempty"`
	ZoneID              int            `json:"ZoneId,omitempty"`
	DomainName          string         `json:"DomainName,omitempty"`
	Status              int            `json:"Status,omitempty"`
	ZoneType            int            `json:"ZoneType,omitempty"`
	IsCustomerOwned     int            `json:"IsCustomerOwned,omitempty"`
	Comment             string         `json:"Comment,omitempty"`
	Version             string         `json:"Version,omitempty"`
	Records             DNSRecords     `json:"Records"`
	FailoverGroups      []MasterServer `json:"FailoverGroups:omitempty"`
	LoadBalancingGroups []MasterServer `json:"LoadBalancingGroups:omitempty"`
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
	Name     string `json:"Name,omitempty"`
	TTL      string `json:"TTL,omitempty"`
	Rdata    string `json:"Rdata,omitempty"`
	VerifyID int    `json:"VerifyId,omitemtpy"`
}

// LiteralResponse -
type LiteralResponse struct {
	Value interface{}
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
func (c *DNSRouteAPIClient) GetZone(id int) (*ZoneResponse, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/routezone?id=%d", c.Config.AccountNumber, id)
	log.Printf("apiURL:%s", apiURL)
	request, err := c.BaseAPIClient.BuildRequest("GET", apiURL, nil, false)

	if err != nil {
		return nil, fmt.Errorf("GetZone: %v", err)
	}

	parsedResponse := ZoneResponse{}
	log.Printf("dnsroute_client>>GetZone>>parsedResponse:%v", parsedResponse)

	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetOrigin: %v", err)
	}

	return &parsedResponse, nil
}

// AddZone -
func (c *DNSRouteAPIClient) AddZone(zone *ZoneRequest) (int, error) {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/zone", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, zone, false)
	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddZone: %v", err)
	}
	resp, err := c.BaseAPIClient.SendRequestWithStringResponse(request)

	if err != nil {
		return -1, fmt.Errorf("api>>dnsroute_client>>AddZone: %v", err)
	}

	zoneID, err := strconv.Atoi(*resp)
	if err != nil {
		return -1, fmt.Errorf("dnsroute_client>>AddZone->API Response Error: %v", err)
	}
	return zoneID, nil
}

// UpdateZone -
func (c *DNSRouteAPIClient) UpdateZone(zone *ZoneRequest) error {
	apiURL := fmt.Sprintf("/v2/mcc/customers/%s/dns/zone", c.Config.AccountNumber)
	request, err := c.BaseAPIClient.BuildRequest("POST", apiURL, zone, false)
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
