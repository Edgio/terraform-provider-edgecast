package data

import (
	"fmt"
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"os"
	"time"
)

func account() string {
	return os.Getenv("ACCOUNT_NUMBER")
}

func email() string {
	return fmt.Sprintf("devenablement+testing%d@edgecast.com", time.Now().Unix())
}

func unique(s string) string {
	return fmt.Sprintf("%s%d", s, time.Now().Unix())
}

func Create(cfg edgecast.SDKConfig) {
	accountNumber, customerUser := createCustomerData(cfg)
	fmt.Println("account number:", accountNumber)
	fmt.Println("customer user:", customerUser)

	edgeCnameID := createEdgeCnameData(cfg)
	fmt.Println("edge cname id:", edgeCnameID)

	groupID, masterServerGroupID, masterServerA, masterServerB, secondaryServerGroupID, tsgID, zoneID := createDNSData(cfg)
	fmt.Println("dns group id:", groupID)
	fmt.Println("master server group id:", masterServerGroupID)
	fmt.Println("master server a id:", masterServerA)
	fmt.Println("master server b id:", masterServerB)
	fmt.Println("secondary server group id:", secondaryServerGroupID)
	fmt.Println("tsg id:", tsgID)
	fmt.Println("zone id:", zoneID)

	rateRuleID := createWAFData(cfg)
	fmt.Println("rate rule id:", rateRuleID)

}
