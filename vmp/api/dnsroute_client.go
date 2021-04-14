package api

import (
	"fmt"
	"log"
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
