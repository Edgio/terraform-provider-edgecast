package data

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf"
	"terraform-provider-edgecast/test/integration/cmd/create-import-data/internal"
)

func createWAFData(cfg edgecast.SDKConfig) (rateRuleID string) {
	svc := internal.Check(waf.New(cfg)).(*waf.WAFService)
	rateRuleID = createWAFRateRule(svc)
	return
}

func createWAFRateRule(svc *waf.WAFService) (id string) {
	var tru bool = true
	var fal bool = false

	params := waf.AddRateRuleParams{
		AccountNumber: account(),
		RateRule: waf.RateRule{
			ConditionGroups: []waf.ConditionGroup{
				{
					Conditions: []waf.Condition{
						{
							Target: waf.Target{
								Type:  "REQUEST_HEADERS",
								Value: "User-Agent",
							},
							OP: waf.OP{
								IsCaseInsensitive: &tru,
								IsNegated:         &fal,
								Type:              "EM",
								Values:            []string{"Mozilla/5.0", "Chrome/91.0.4472.114"},
							},
						},
					},
					Name: "1",
				},
				{
					Conditions: []waf.Condition{
						{
							Target: waf.Target{
								Type:  "REQUEST_HEADERS",
								Value: "User-Agentz",
							},
							OP: waf.OP{
								IsCaseInsensitive: &tru,
								IsNegated:         &fal,
								Type:              "EM",
								Values:            []string{"Mozilla/5.0", "Chrome/91.0.4472.114"},
							},
						},
					},
					Name: "2",
				},
			},
			CustomerID:  account(),
			Disabled:    false,
			DurationSec: 1,
			Keys:        []string{"IP", "USER_AGENT"},
			Name:        "rate rule 1",
			Num:         10,
		},
	}
	return internal.Check(svc.AddRateRule(params)).(string)
}
