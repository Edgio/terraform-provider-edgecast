// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"encoding/base64"

	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/waf"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/access"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/bot"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/custom"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/managed"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules/rate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/scopes"
)

func createWAFData(cfg config.Config) WAFResult {
	svc := internal.Check(waf.New(cfg.SDKConfig))
	wafManagedRuleID := createWAFKManagedRule(svc, cfg.AccountNumber)
	wafAccessRuleID := createWAFAccessRule(svc, cfg.AccountNumber)
	wafBotRuleID := createBotRule(svc, cfg.AccountNumber)
	wafRateRuleID := createWAFRateRule(svc, cfg.AccountNumber)
	wafCustomRuleID := createWAFCustomRule(svc, cfg.AccountNumber)
	wafScopesID := createWAFScopes(
		svc,
		wafRateRuleID,
		wafAccessRuleID,
		wafManagedRuleID,
		wafCustomRuleID,
		wafBotRuleID,
		cfg.AccountNumber)

	return WAFResult{
		RateRuleID:    wafRateRuleID,
		AccessRuleID:  wafAccessRuleID,
		BotRuleID:     wafBotRuleID,
		CustomRuleID:  wafCustomRuleID,
		ManagedRuleID: wafManagedRuleID,
		ScopesID:      wafScopesID,
	}
}

func createWAFRateRule(svc *waf.WafService, accountNumber string) (id string) {
	params := rate.AddRateRuleParams{
		AccountNumber: accountNumber,
		RateRule: rate.RateRule{
			ConditionGroups: []rate.ConditionGroup{
				{
					Conditions: []rate.Condition{
						{
							Target: rate.Target{
								Type:  "REQUEST_HEADERS",
								Value: "User-Agent",
							},
							OP: rate.OP{
								IsCaseInsensitive: internal.Pointer(true),
								IsNegated:         internal.Pointer(false),
								Type:              "EM",
								Values: []string{
									"Mozilla/5.0",
									"Chrome/91.0.4472.114",
								},
							},
						},
					},
					Name: "1",
				},
				{
					Conditions: []rate.Condition{
						{
							Target: rate.Target{
								Type:  "REQUEST_HEADERS",
								Value: "User-Agentz",
							},
							OP: rate.OP{
								IsCaseInsensitive: internal.Pointer(true),
								IsNegated:         internal.Pointer(false),
								Type:              "EM",
								Values: []string{
									"Mozilla/5.0",
									"Chrome/91.0.4472.114",
								},
							},
						},
					},
					Name: "2",
				},
			},
			CustomerID:  accountNumber,
			Disabled:    false,
			DurationSec: 1,
			Keys:        []string{"IP", "USER_AGENT"},
			Name:        internal.Unique("-rate_rule"),
			Num:         10,
		},
	}
	return internal.Check(svc.Rate.AddRateRule(params))
}

func createBotRule(svc *waf.WafService, accountNumber string) (id string) {
	params := bot.AddBotRuleSetParams{
		BotRuleSet: bot.BotRuleSet{
			Name: internal.Unique("-bot_rule_set"),
			Directives: []bot.BotRuleDirective{
				{
					SecRule: &rules.SecRule{
						Name: internal.Unique("-sec_rule"),
						Action: rules.Action{
							ID: "77375686",
							Transformations: []rules.Transformation{
								rules.TransformNone,
							},
						},
						Operator: rules.Operator{
							IsNegated: true,
							Type:      rules.OpNumberEquality,
							Value:     "1",
						},
						Variables: []rules.Variable{
							{
								IsCount: true,
								Type:    rules.VarRequestCookies,
								Matches: []rules.Match{
									{
										IsNegated: false,
										IsRegex:   false,
									},
									{
										IsNegated: true,
										IsRegex:   true,
										Value:     "cookiename",
									},
								},
							},
						},
					},
				},
			},
		},
		AccountNumber: accountNumber,
	}

	id = internal.Check(svc.Bot.AddBotRuleSet(params))

	return
}

func createWAFAccessRule(
	svc *waf.WafService,
	accountNumber string,
) (id string) {
	params := access.AddAccessRuleParams{
		AccountNumber: accountNumber,
		AccessRule: access.AccessRule{
			AllowedHTTPMethods:         []string{"GET", "POST"},
			AllowedRequestContentTypes: []string{"application/json"},
			ASNAccessControls:          nil,
			CookieAccessControls:       nil,
			CountryAccessControls:      nil,
			CustomerID:                 accountNumber,
			DisallowedExtensions:       nil,
			DisallowedHeaders:          nil,
			IPAccessControls:           nil,
			MaxFileSize:                0,
			Name:                       internal.Unique("-access_rule"),
			RefererAccessControls:      nil,
			ResponseHeaderName:         "",
			URLAccessControls:          nil,
			UserAgentAccessControls:    nil,
		},
	}
	return internal.Check(svc.Access.AddAccessRule(params))
}

func createWAFCustomRule(
	svc *waf.WafService,
	accountNumber string,
) (id string) {
	params := custom.AddCustomRuleSetParams{
		AccountNumber: accountNumber,
		CustomRuleSet: custom.CustomRuleSet{
			Directives: []custom.CustomRuleDirective{
				{
					SecRule: rules.SecRule{
						Action: rules.Action{
							Message: "Invalid user agent",
							Transformations: []rules.Transformation{
								rules.TransformNone,
							},
						},
						Name: internal.Unique("-sec_rule"),
						Operator: rules.Operator{
							IsNegated: false,
							Type:      rules.OpNumberEquality,
							Value:     "bot",
						},
						Variables: []rules.Variable{
							{
								Type: rules.VarRequestHeaders,
								Matches: []rules.Match{
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
			Name: internal.Unique("-custom_rule"),
		},
	}

	return internal.Check(svc.Custom.AddCustomRuleSet(params))
}

func createWAFKManagedRule(
	svc *waf.WafService,
	accountNumber string,
) (id string) {
	params := managed.AddManagedRuleParams{
		AccountNumber: accountNumber,
		ManagedRule: managed.ManagedRule{
			Name:           internal.Unique("-managed_rule"),
			RulesetID:      "ECRS",
			RulesetVersion: "2020-05-01",
			DisabledRules:  nil,
			GeneralSettings: managed.GeneralSettings{
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
			RuleTargetUpdates: []managed.RuleTargetUpdate{
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

	return internal.Check(svc.Managed.AddManagedRule(params))
}

func createWAFScopes(
	svc *waf.WafService,
	rateRuleID,
	accessRuleID,
	managedRuleID,
	customRuleID,
	botRuleID string,
	accountNumber string,
) (id string) {
	trueVar := true
	encodedMessage := base64.StdEncoding.EncodeToString([]byte("hello!"))
	status404 := 404
	redirectURL := "https://www.devenblment.com/redirected"

	params := scopes.Scopes{
		CustomerID: accountNumber,
		Scopes: []scopes.Scope{
			{
				Name: internal.Unique("-scope"),
				Host: scopes.MatchCondition{
					Type:              "EM",
					IsCaseInsensitive: &trueVar,
					Values: &[]string{
						"devenblment.com",
						"devenblment2.com",
					},
				},
				Path: scopes.MatchCondition{
					Type:   "EM",
					Values: &[]string{"/account", "/admin"},
				},
				Limits: &[]scopes.Limit{
					{
						ID: rateRuleID,
						Action: scopes.LimitAction{
							Name:               "Custom action",
							DurationSec:        10,
							ENFType:            "CUSTOM_RESPONSE",
							ResponseBodyBase64: &encodedMessage,
							ResponseHeaders: &map[string]string{
								"key1": "value1",
							},
							Status: &status404,
						},
					},
				},
				ACLAuditID: &accessRuleID,
				ACLProdID:  &accessRuleID,
				ACLProdAction: &scopes.ProdAction{
					Name:    "Access Rule Action",
					ENFType: "REDIRECT_302",
					URL:     &redirectURL,
				},
				ProfileAuditID: &managedRuleID,
				ProfileProdID:  &managedRuleID,
				ProfileProdAction: &scopes.ProdAction{
					Name:    "Managed Rule Action",
					ENFType: "BLOCK_REQUEST",
				},
				RuleAuditID: &customRuleID,
				RuleProdID:  &customRuleID,
				RuleProdAction: &scopes.ProdAction{
					Name:    "Custom Rule Action",
					ENFType: "ALERT",
				},
			},
		},
	}

	return internal.Check(svc.Scopes.ModifyAllScopes(params)).ID
}
