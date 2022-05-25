package data

import (
	"fmt"
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"os"
	"time"
)

var accountNumber = os.Getenv("ACCOUNT_NUMBER")

func account() string {
	return accountNumber
}

func email() string {
	return fmt.Sprintf("devenablement+testing%d@edgecast.com", time.Now().Unix())
}

func unique(s string) string {
	n := fmt.Sprintf("%d", time.Now().Unix())
	n = n[len(n)-4:]
	return fmt.Sprintf("%s%s", n, s)
}

func Create(cfg edgecast.SDKConfig) {

	rulesEnginePolicyID := createRulesEnginePolicyData(cfg)
	fmt.Println("rules engine policy id:", rulesEnginePolicyID)

	accountNumber, customerUser := createCustomerData(cfg)
	fmt.Println("account number:", accountNumber)
	fmt.Println("customer user:", customerUser)

	originID := createOriginData(cfg)
	fmt.Println("origin id:", originID)

	edgeCnameID := createEdgeCnameData(cfg)
	fmt.Println("edge cname id:", edgeCnameID)

	groupID, masterServerGroupID, masterServerA, masterServerB, secondaryServerGroupID, tsgID, zoneID := createDNSData(cfg)
	fmt.Println("dns group id:", groupID)
	fmt.Println("dns master server group id:", masterServerGroupID)
	fmt.Println("dns master server a id:", masterServerA)
	fmt.Println("dns master server b id:", masterServerB)
	fmt.Println("dns secondary server group id:", secondaryServerGroupID)
	fmt.Println("dns tsg id:", tsgID)
	fmt.Println("dns zone id:", zoneID)

	wafRateRuleID, wafAccessRuleID, wafCustomRuleID, wafManagedRuleID, wafScopesID := createWAFData(cfg)
	fmt.Println("waf access rule id:", wafAccessRuleID)
	fmt.Println("waf custom rule id:", wafCustomRuleID)
	fmt.Println("waf managed rule b id:", wafManagedRuleID)
	fmt.Println("waf rate rule id:", wafRateRuleID)
	fmt.Println("waf scopes id:", wafScopesID)

}
