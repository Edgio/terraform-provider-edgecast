package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const rulesEngineRelURLFormat = "rules-engine/v1.1/%s"

//RulesEngineAPIClient -
type RulesEngineAPIClient struct {
	Config        *ClientConfig
	BaseAPIClient *BaseClient
}

//AddDeployPolicyResponse -
type AddDeployPolicyResponse struct {
	ID          string                   `json:"id,omitempty"`
	AtID        string                   `json:"@id,omitempty"`
	Type        string                   `json:"@type,omitempty"`
	Links       []map[string]interface{} `json:"@links,omitempty"`
	State       string                   `json:"state,omitempty"`
	Environment string                   `json:"environment,omitempty"`
	CustomerID  string                   `json:"customer_id"`
	CreatedAt   time.Time                `json:"created_at,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at,omitempty"`
	IsVisible   bool                     `json:"is_visible,omitempty"`
	Policies    AddPolicyResponse        `json:"policies,omitempty"`
	History     []map[string]interface{} `json:"history,omitempty"`
	User        User                     `json:"user,omitempty"`
}

//UpdateDeployPolicyStateResponse -
type UpdateDeployPolicyStateResponse struct {
	ID    string `json:"id,omitempty"`
	State string `json:"state,omitempty"`
}

//AddPolicyResponse -
type AddPolicyResponse struct {
	ID          string    `json:"id,omitempty"`
	Type        string    `json:"@type,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	PolicyType  string    `json:"policy_type,omitempty"`
	State       string    `json:"state,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Rules       []Rule    `json:"rules,omitempty"`
}

//UpdatePolicyResponse -
type UpdatePolicyResponse struct {
	ID          string    `json:"id,omitempty"`
	Type        string    `json:"@type,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	PolicyType  string    `json:"policy_type,omitempty"`
	State       string    `json:"state,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Rules       []Rule    `json:"rules,omitempty"`
}

//Rule -
type Rule struct {
	ID          string                   `json:"id,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Ordinal     int                      `json:"ordinal,omitempty"`
	CreatedAt   time.Time                `json:"created_at,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at,omitempty"`
	Matches     []map[string]interface{} `json:"matches,omitempty"`
}

//Match -
type Match struct {
	ID         int    `json:"id"`
	Type       string `json:"@type"`
	Ordinal    int    `json:"ordinal,omitempty"`
	Value      string `json:"value,omitempty"`
	Codes      string `json:"codes,omitempty"`
	Compare    string `json:"compare,omitempty"`
	Encoded    bool   `json:"encoded,omitempty"`
	Hostnames  string `json:"hostnames,omitempty"`
	IgnoreCase bool   `json:"ignore-case,omitempty"`
	Name       string `json:"name,omitempty"`
	RelativeTo string `json:"relative-to,omitempty"`
	Result     string `json:"result,omitempty"`
	Matches    []map[string]interface{}
	Features   []map[string]interface{}
}

//Feature -
type Feature struct {
	Action          string   `json:"action,omitempty"`
	Code            string   `json:"code,omitempty"`
	Destination     string   `json:"destination,omitempty"`
	Enabled         bool     `json:"enabled,omitempty"`
	Expires         int      `json:"expires,omitempty"`
	Extensions      string   `json:"extensions,omitempty"`
	Format          string   `json:"format,omitempty"`
	HeaderName      string   `json:"header-name,omitempty"`
	HeaderValue     string   `json:"header-value,omitempty"`
	Instance        string   `json:"instance,omitempty"`
	KbytesPerSecond int      `json:"kbytes-per-second,omitempty"`
	MediaTypes      []string `json:"mediaTypes,omitempty"`
	Methods         string   `json:"methods,omitempty"`
	Milliseconds    int      `json:"milliseconds,omitempty"`
	Mode            string   `json:"mode,omitempty"`
	Name            string   `json:"name,omitempty"`
	Names           []string `json:"names,omitempty"`
	Parameters      string   `json:"parameters,omitempty"`
	PrebufSeconds   int      `json:"prebuf-seconds,omitempty"`
	Requests        int      `json:"requests,omitempty"`
	Seconds         int      `json:"seconds,omitempty"`
	SeekEnd         string   `json:"seekEnd,omitempty"`
	SeekStart       string   `json:"seekStart,omitempty"`
	Site            string   `json:"site,omitempty"`
	Source          string   `json:"source,omitempty"`
	Status          string   `json:"status,omitempty"`
	Tags            string   `json:"tags,omitempty"`
	Treatment       string   `json:"treatment,omitempty"`
	Units           string   `json:"units,omitempty"`
	Value           string   `json:"value,omitempty"`
}

//AddDeployRequest -
type AddDeployRequest struct {
	PolicyID    int    `json:"policy_id"`
	Environment string `json:"environment,omitempty"`
	Message     string `json:"message"`
}

//User -
type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

//NewRulesEngineAPIClient -
func NewRulesEngineAPIClient(config *ClientConfig) *RulesEngineAPIClient {
	APIClient := &RulesEngineAPIClient{
		Config:        config,
		BaseAPIClient: config.BaseClient,
	}

	return APIClient
}

// GetPolicy -
func (APIClient *RulesEngineAPIClient) GetPolicy(accountNumber string, customerUserID string, portalTypeID string, policyID int) (map[string]interface{}, error) {
	relURL := formatRulesEngineRelURL("policies/%d", policyID)
	request, err := APIClient.BaseAPIClient.BuildRequest("GET", relURL, nil, true)

	if err != nil {
		return nil, fmt.Errorf("GetPolicy: %v", err)
	}

	// account number hex string -> customer ID
	customerID, err := strconv.ParseInt(accountNumber, 16, 64)

	if err != nil {
		return nil, fmt.Errorf("GetPolicy: ParseInt: %v", err)
	}

	request.Header.Set("Portals_CustomerId", strconv.FormatInt(customerID, 10))
	request.Header.Set("Portals_UserId", customerUserID)
	request.Header.Set("Portals_PortalTypeId", portalTypeID)

	parsedResponse := make(map[string]interface{})

	_, err = APIClient.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("GetPolicy: %v", err)
	}

	return parsedResponse, nil
}

//AddPolicy -
func (c *RulesEngineAPIClient) AddPolicy(policy string, accountNumber string, portalTypeID string, customerUserID string) (*AddPolicyResponse, error) {
	request, err := c.BaseAPIClient.BuildRequest("POST", "rules-engine/v1.1/policies", policy, true)

	if err != nil {
		return nil, fmt.Errorf("AddPolicy: %v", err)
	}

	// account number hex string -> customer ID
	customerID, err := strconv.ParseInt(accountNumber, 16, 64)

	if err != nil {
		return nil, fmt.Errorf("AddPolicy: ParseInt: %v", err)
	}

	request.Header.Set("Portals_CustomerId", strconv.FormatInt(customerID, 10))
	request.Header.Set("Portals_UserId", customerUserID)
	request.Header.Set("Portals_PortalTypeId", portalTypeID)
	parsedResponse := &AddPolicyResponse{}

	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("AddPolicy: %v", err)
	}

	return parsedResponse, nil
}

//DeployPolicy -
func (c *RulesEngineAPIClient) DeployPolicy(body *AddDeployRequest, accountNumber string, portalTypeID string, customerUserID string) (*AddDeployPolicyResponse, error) {
	request, err := c.BaseAPIClient.BuildRequest("POST", "rules-engine/v1.1/deploy-requests", body, true)

	if err != nil {
		return nil, fmt.Errorf("DeployPolicy: %v", err)
	}

	// account number hex string -> customer ID
	customerID, err := strconv.ParseInt(accountNumber, 16, 64)

	if err != nil {
		return nil, fmt.Errorf("DeployPolicy: ParseInt: %v", err)
	}

	request.Header.Set("Portals_CustomerId", strconv.FormatInt(customerID, 10))
	request.Header.Set("Portals_UserId", customerUserID)
	request.Header.Set("Portals_PortalTypeId", portalTypeID)

	parsedResponse := &AddDeployPolicyResponse{}

	_, err = c.BaseAPIClient.SendRequest(request, &parsedResponse)

	if err != nil {
		return nil, fmt.Errorf("DeployPolicy: %v", err)
	}

	return parsedResponse, nil
}

func removeHexPrefix(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

func formatRulesEngineRelURL(subFormat string, params ...interface{}) string {
	subPath := fmt.Sprintf(subFormat, params...)
	return fmt.Sprintf(rulesEngineRelURLFormat, subPath)
}
