// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

type PopulationResult struct {
	Customer    CustomerResult    `json:"customer,omitempty"`
	Origin      OriginResult      `json:"origin,omitempty"`
	CNAME       CNAMEResult       `json:"cname,omitempty"`
	DNS         DNSResult         `json:"dns,omitempty"`
	WAF         WAFResult         `json:"waf,omitempty"`
	RulesEngine RulseEngineResult `json:"rules_engine,omitempty"`
	OriginV3    OriginV3Result    `json:"origin_v3,omitempty"`
	CPS         CPSResult         `json:"cps,omitempty"`
	BotManager  BotManagerResult  `json:"bot_manager,omitempty"`
}

type CustomerResult struct {
	AccountNumber  string `json:"account_number,omitempty"`
	CustomerUserID int    `json:"customer_user_id,omitempty"`
}

type OriginResult struct {
	OriginID int `json:"origin_id,omitempty"`
}

type CNAMEResult struct {
	EdgeCnameID int `json:"edge_cname_id,omitempty"`
}

type DNSResult struct {
	GroupID                int `json:"group_id,omitempty"`
	MasterServerGroupID    int `json:"master_server_group_id,omitempty"`
	MasterServerA          int `json:"master_server_a,omitempty"`
	MasterServerB          int `json:"master_server_b,omitempty"`
	SecondaryServerGroupID int `json:"secondary_server_group_id,omitempty"`
	TsgID                  int `json:"tsg_id,omitempty"`
	ZoneID                 int `json:"zone_id,omitempty"`
}

type WAFResult struct {
	RateRuleID    string `json:"rate_rule_id,omitempty"`
	AccessRuleID  string `json:"access_rule_id,omitempty"`
	BotRuleID     string `json:"bot_rule_id,omitempty"`
	CustomRuleID  string `json:"custom_rule_id,omitempty"`
	ManagedRuleID string `json:"managed_rule_id,omitempty"`
	ScopesID      string `json:"scopes_id,omitempty"`
}

type RulseEngineResult struct {
	PolicyID string `json:"policy_id,omitempty"`
}

type OriginV3Result struct {
	GroupIdV3  int32 `json:"group_id_v3,omitempty"`
	OriginIdV3 int32 `json:"origin_id_v3,omitempty"`
}

type CPSResult struct {
	CertificateID int64 `json:"certificate_id,omitempty"`
}

type BotManagerResult struct {
	BotManagerID string `json:"bot_manager_id,omitempty"`
}
