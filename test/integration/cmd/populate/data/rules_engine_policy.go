package data

import (
	"fmt"
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/rulesengine"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

var policyString = `{
		"@type": "policy-create",
		"name": "a%s",
		"description": "This is a test of the policy-create process.",
		"platform": "http_large",
		"state": "locked",
		"rules": [
			{
				"name": "Deny POST - %s",
				"description": "Allow all POST requests",
				"matches": [{
					"type": "match.request.request-method.literal",
					"value" : "GET",
					"features": [{
						"type": "feature.access.deny-access",
						"enabled": false
					}]
				}]
			}
		]
	}`

func createRulesEnginePolicyData(cfg edgecast.SDKConfig) (id string) {
	svc := internal.Check(rulesengine.New(cfg))
	id = createPolicyV4(svc)
	return
}

func createPolicyV4(svc *rulesengine.RulesEngineService) (id string) {
	params := rulesengine.AddPolicyParams{
		AccountNumber:  account(),
		PolicyAsString: fmt.Sprintf(policyString, unique("policy"), unique("rule")),
	}

	res := internal.Check(svc.AddPolicy(params))
	return res.ID
}
