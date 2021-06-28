package model

//ResourceREV4 - RulesEngine V4 Policy
type ResourceREV4 struct {
	Policy                 string
	RulesEngineEnvironment string
	Credentials            Credentials
	TestCustomerInfo       CustomerInfo
}

//ResourceNewCustomer - new MCC customer
type ResourceNewCustomer struct {
	CustomerInfo NewCustomerInfo
	Credential   Credentials
}

// Credentials - credential for API
type Credentials struct {
	ApiToken         string
	IdsClientSecret  string
	IdsClientID      string
	IdsScope         string
	ApiAddress       string
	ApiAddressLegacy string
	IdsAddress       string
}

// CustomerInfo - target customer for testing
type CustomerInfo struct {
	AccountNumber  string
	CustomerUserID string
	PortalTypeID   int
}

// NewCustomerInfo - resource for new customer
type NewCustomerInfo struct {
	CompanyName      string
	ServiceLevelCode string
	Services         []int
	DeliveryRegion   int
	AccessModules    []int
}

//ResourceNewCustomer - new MCC customer
type ResourceNewCustomerUser struct {
	CustomerUserInfo NewCustomerUserInfo
	Credential       Credentials
}

// NewCustomerInfo - resource for new customer
type NewCustomerUserInfo struct {
	AccountNumber string
	FirstName     string
	LastName      string
	Email         string
	IsAdmin       bool
}

//MasterServerGroupRequest -
type MasterServerGroupRequest struct {
	Name       string         `json:"Name"`
	Masters    []MasterServer `json:"MasterServers"`
	Credential Credentials
}

// MasterServerGroupResponse -
type MasterServerGroupResponse struct {
	MasterGroupID int            `json:"MasterGroupId"`
	Name          string         `json:"Name"`
	Masters       []MasterServer `json:"Masters"`
}

// MasterServer
type MasterServer struct {
	ID        int    `json:"Id,omitempty"`
	Name      string `json:"Name"`
	IPAddress string `json:"IPAddress"`
}
