package test

//ResourceREV4 - RulesEngine V4 Policy
type ResourceREV4 struct {
	policy                 string
	rulesEngineEnvironment string
	credentials            Credentials
	testcustomerinfo       CustomerInfo
}

//ResourceNewCustomer - new MCC customer
type ResourceNewCustomer struct {
	customerInfo NewCustomerInfo
	credential   Credentials
}

// Credentials - credential for API
type Credentials struct {
	apitoken        string
	idsclientsecret string
	idsclientID     string
	idsscope        string
	apiaddress      string
	idsaddress      string
}

// CustomerInfo - target customer for testing
type CustomerInfo struct {
	accountnumber  string
	customeruserID string
	portaltypeID   int
}

// NewCustomerInfo - resource for new customer
type NewCustomerInfo struct {
	companyname      string
	servicelevelcode string
	services         []int
	deliveryregion   int
	accessmodules    []int
}

//ResourceNewCustomer - new MCC customer
type ResourceNewCustomerUser struct {
	customerUserInfo NewCustomerUserInfo
	credential       Credentials
}

// NewCustomerInfo - resource for new customer
type NewCustomerUserInfo struct {
	accountnumber string
	firstname     string
	lastname      string
	email         string
	isadmin       bool
}
