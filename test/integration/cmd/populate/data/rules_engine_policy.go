// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"fmt"

	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/rulesengine"
)

var fmtPolicyString = `{
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

func createRulesEnginePolicyData(cfg config.Config) RulseEngineResult {
	svc := internal.Check(rulesengine.New(cfg.SDKConfig))
	id := createPolicyV4(svc, cfg.AccountNumber)

	return RulseEngineResult{
		PolicyID: id,
	}
}

func createPolicyV4(
	svc *rulesengine.RulesEngineService,
	accountNumber string,
) (id string) {
	policyString := fmt.Sprintf(
		fmtPolicyString,
		internal.Unique("policy"),
		internal.Unique("rule"))

	params := rulesengine.AddPolicyParams{
		AccountNumber:  accountNumber,
		PolicyAsString: policyString,
	}

	res := internal.Check(svc.AddPolicy(params))
	return res.ID
}
