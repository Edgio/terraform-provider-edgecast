package data

import (
	"fmt"
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/rulesengine"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
	"time"
)

const emptyPolicy string = "{\"@type\":\"policy-create\",\"name\":\"Terraform Placeholder - %s\",\"platform\":\"%s\",\"rules\":[{\"@type\":\"rule-create\",\"description\":\"Placeholder rule created by the Edgecast Terraform Provider\",\"matches\":[{\"features\":[{\"type\":\"feature.comment\",\"value\":\"Empty policy created on %s\"}],\"ordinal\":1,\"type\":\"match.always\"}],\"name\":\"Placeholder Rule\"}],\"state\":\"locked\"}"
const samplePolicy = `{
    "name": "test policy %s",
    "description": "This is a test policy!",
    "platform": "http_large",
    "rules": [
        {
            "name": "rule1",
            "description": "This is a test rule.",
            "matches": [
                {
                    "type": "match.always",
                    "features": [
                        {
                            "type": "feature.comment",
                            "value": "Update this comment!"
                        }
                    ]
                }
            ]
        }
    ]
}
`

func createRulesEnginePolicyData(cfg edgecast.SDKConfig) (id string) {
	svc := internal.Check(rulesengine.New(cfg))
	id = createPolicyV4(svc)
	return
}

func createPolicyV4(svc *rulesengine.RulesEngineService) (id string) {
	timestamp := time.Now().Format(time.RFC3339)

	params := rulesengine.AddPolicyParams{
		AccountNumber:  account(),
		CustomerUserID: account(),
		PortalTypeID:   "1",
		PolicyAsString: fmt.Sprintf(
			samplePolicy,
			timestamp),
	}

	res := internal.Check(svc.AddPolicy(params))
	return res.ID
}
