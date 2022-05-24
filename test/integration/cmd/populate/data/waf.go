package data

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

func createWAFData(cfg edgecast.SDKConfig) (rateRuleID, wafAccessRuleID, wafCustomRuleID, wafManagedRuleID, wafScopesID string) {
	svc := internal.Check(waf.New(cfg))
	wafManagedRuleID = createWAFKManagedRule(svc)
	wafAccessRuleID = createWAFAccessRule(svc)
	rateRuleID = createWAFRateRule(svc)
	wafScopesID = createWAFScopes(svc)
	wafCustomRuleID = createWAFCustomRule(svc)
	return
}

func createWAFRateRule(svc *waf.WAFService) (id string) {

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
								IsCaseInsensitive: internal.Pointer(true),
								IsNegated:         internal.Pointer(false),
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
								IsCaseInsensitive: internal.Pointer(true),
								IsNegated:         internal.Pointer(false),
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
	return internal.Check(svc.AddRateRule(params))
}

func createWAFAccessRule(svc *waf.WAFService) (id string) {
	params := waf.AddAccessRuleParams{
		AccountNumber: account(),
		AccessRule: waf.AccessRule{
			AllowedHTTPMethods:         []string{"GET", "POST"},
			AllowedRequestContentTypes: []string{"application/json"},
			ASNAccessControls:          nil,
			CookieAccessControls:       nil,
			CountryAccessControls:      nil,
			CustomerID:                 account(),
			DisallowedExtensions:       nil,
			DisallowedHeaders:          nil,
			IPAccessControls:           nil,
			MaxFileSize:                0,
			Name:                       "access rule 1",
			RefererAccessControls:      nil,
			ResponseHeaderName:         "",
			URLAccessControls:          nil,
			UserAgentAccessControls:    nil,
		},
	}
	return internal.Check(svc.AddAccessRule(params))
}

func createWAFCustomRule(svc *waf.WAFService) (id string) {
	params := waf.AddCustomRuleSetParams{
		AccountNumber: account(),
		CustomRuleSet: waf.CustomRuleSet{
			Directives: []waf.Directive{
				{
					SecRule: waf.SecRule{
						Action: waf.Action{
							Message:         "Invalid user agent",
							Transformations: []string{"NONE"},
						},
						Name: "a name",
						Operator: waf.Operator{
							IsNegated: false,
							Type:      "EQ",
							Value:     "bot",
						},
						Variables: []waf.Variable{
							{
								Type: "REQUEST_HEADERS",
								Matches: []waf.Match{
									{
										IsNegated: false,
										IsRegex:   false,
										Value:     "User-Agent",
									},
								},
								IsCount: false,
							},
						},
					},
				},
			},
			Name: unique("My-Rule"),
		},
	}

	return internal.Check(svc.AddCustomRuleSet(params))
}

func createWAFKManagedRule(svc *waf.WAFService) (id string) {
	params := waf.AddManagedRuleParams{
		AccountNumber: account(),
		ManagedRule: waf.ManagedRule{
			Name:           "Terraform Managed Rule",
			RulesetID:      "ECRS",
			RulesetVersion: "2020-05-01",
			DisabledRules:  nil,
			GeneralSettings: waf.GeneralSettings{
				AnomalyThreshold:     10,
				ArgLength:            8000,
				ArgNameLength:        1024,
				CombinedFileSizes:    6291456,
				IgnoreCookie:         []string{"ignoredCookie"},
				IgnoreHeader:         []string{"ignoredHeaders"},
				IgnoreQueryArgs:      []string{"ignoredQuery"},
				JsonParser:           true,
				MaxFileSize:          0,
				MaxNumArgs:           512,
				ParanoiaLevel:        1,
				ProcessRequestBody:   true,
				ResponseHeaderName:   "X-EC-Security-Audit",
				TotalArgLength:       640000,
				ValidateUtf8Encoding: true,
				XmlParser:            true,
			},
			Policies: []string{
				"r4020_tw_cpanel.conf.json",
				"r4040_tw_drupal.conf.json",
				"r4030_tw_iis.conf.json",
				"r4070_tw_joomla.conf.json",
				"r4050_tw_microsoft_sharepoint.conf.json",
				"r4060_tw_wordpress.conf.json",
				"r5040_cross_site_scripting.conf.json",
				"r2000_ec_custom_rule.conf.json",
				"r5021_http_attack.conf.json",
				"r5020_http_protocol_violation.conf.json",
				"r5043_java_attack.conf.json",
				"r5030_local_file_inclusion.conf.json",
				"r5033_php_injection.conf.json",
				"r5032_remote_code_execution.conf.json",
				"r5031_remote_file_inclusion.conf.json",
				"r5010_scanner_detection.conf.json",
				"r5042_session_fixation.conf.json",
				"r5041_sql_injection.conf.json",
			},
			RuleTargetUpdates: []waf.RuleTargetUpdate{
				{
					IsNegated:     false,
					IsRegex:       true,
					ReplaceTarget: "",
					Target:        "ARGS",
					TargetMatch:   "ignoredArgumentException",
				},
			},
		},
	}

	return internal.Check(svc.AddManagedRule(params))
}

func createWAFScopes(svc *waf.WAFService) (id string) {
	params := waf.Scopes{
		CustomerID: account(),
		Scopes: []waf.Scope{
			{
				Name: "scopes-web-security",
				Host: waf.MatchCondition{
					IsCaseInsensitive: internal.Pointer(false),
					Type:              "EM",
					Values:            &[]string{"site1.com/path1", "site2.com"},
				},
				Limits: &[]waf.Limit{
					{
						ID: "one",
						Action: waf.LimitAction{
							DurationSec:        60,
							ENFType:            "CUSTOM_RESPONSE",
							Name:               "limit action custom",
							ResponseBodyBase64: internal.Pointer("SGVsbG8sIHdvcmxkIQo="),
							ResponseHeaders:    nil,
							Status:             internal.Pointer(404),
						},
					},
				},
				Path: waf.MatchCondition{
					IsCaseInsensitive: internal.Pointer(false),
					IsNegated:         internal.Pointer(false),
					Type:              "GLOB",
					Value:             internal.Pointer("*"),
				},
				//	ACLAuditAction: &waf.AuditAction{
				//			ID:   "",
				//				Name: "",
				//					Type: "ALERT",
				//				},
				//ACLAuditID: internal.Pointer("rule1"),
				ACLProdAction: &waf.ProdAction{
					Name:    "acl action",
					ENFType: "ALERT",
				},
				//	ACLProdID: internal.Pointer("access rule id"),
				//	ProfileAuditAction: &waf.AuditAction{
				//		Type: "ALERT",
				//		},
				//ProfileAuditID: internal.Pointer("audit id"),
				ProfileProdAction: &waf.ProdAction{
					Name:    "custom rule action",
					ENFType: "BLOCK_REQUEST",
				},
				ProfileProdID: nil,
				//RuleAuditAction: nil,
				//RuleAuditID:     internal.Pointer("<Custom Rule ID>"),
				RuleProdAction: &waf.ProdAction{
					Name:    "custom rule action",
					ENFType: "BLOCK_REQUEST",
				},
				RuleProdID: internal.Pointer("<Custom Rule ID>"),
			},
		},
	}

	return internal.Check(svc.ModifyAllScopes(params)).ID

}
