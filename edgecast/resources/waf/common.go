// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

import (
	"errors"
	"fmt"
	"terraform-provider-edgecast/edgecast/api"
	"terraform-provider-edgecast/edgecast/helper"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	sdkwaf "github.com/EdgeCast/ec-sdk-go/edgecast/waf"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf/rules"
)

// buildWAFService builds the SDK WAF service to managed WAF resources
func buildWAFService(config api.ClientConfig) (*sdkwaf.WafService, error) {

	idsCredentials := edgecast.IDSCredentials{
		ClientID:     config.IdsClientID,
		ClientSecret: config.IdsClientSecret,
		Scope:        config.IdsScope,
	}

	sdkConfig := edgecast.NewSDKConfig()
	sdkConfig.APIToken = config.APIToken
	sdkConfig.IDSCredentials = idsCredentials
	sdkConfig.BaseAPIURL = *config.APIURL
	sdkConfig.BaseAPIURLLegacy = *config.APIURLLegacy
	sdkConfig.BaseIDSURL = *config.IdsURL

	return sdkwaf.New(sdkConfig)
}

func expandSecRule(attr interface{}) (*rules.SecRule, error) {

	if attr == nil {
		return nil, nil
	}

	curr, err := helper.ConvertSingletonSetToMap(attr)
	if err != nil {
		return nil, fmt.Errorf("error expanding sec_rule: %w", err)
	}

	// Empty map
	if len(curr) == 0 {
		return nil, nil
	}

	secRule := rules.SecRule{
		Name: curr["name"].(string),
	}

	actionMap, err := helper.ConvertSingletonSetToMap(curr["action"])
	if err != nil {
		return nil, fmt.Errorf("error expanding action: %w", err)
	}

	operatorMap, err := helper.ConvertSingletonSetToMap(curr["operator"])
	if err != nil {
		return nil, fmt.Errorf("error expanding operator: %w", err)
	}

	chainedRule, err := expandChainedRules(curr["chained_rule"])
	if err != nil {
		return nil, fmt.Errorf("error expanding chained_rule: %w", err)
	}
	secRule.ChainedRules = *chainedRule

	variables, err := expandVariables(curr["variable"])
	if err != nil {
		return nil, err
	}
	secRule.Variables = *variables

	if actionId, ok := actionMap["id"]; ok {
		secRule.Action.ID = actionId.(string)
	}

	if actionMsg, ok := actionMap["msg"]; ok {
		secRule.Action.Message = actionMsg.(string)
	}

	if actionT, ok := actionMap["transformations"]; ok {
		arr, err := expandTransformations(actionT)
		if err == nil {
			secRule.Action.Transformations = arr
		} else {
			return nil,
				fmt.Errorf("error reading 'transformations': %w", err)
		}
	}

	if v, ok := operatorMap["is_negated"]; ok {
		secRule.Operator.IsNegated = v.(bool)
	}

	if operatorType, ok := operatorMap["type"]; ok {
		if s, ok := operatorType.(string); ok {
			secRule.Operator.Type = rules.ConvertToOperatorType(s)
		} else {
			return nil, fmt.Errorf("operator type is not a string")
		}
	}

	if operatorValue, ok := operatorMap["value"]; ok {
		secRule.Operator.Value = operatorValue.(string)
	}

	return &secRule, nil
}

// convertToTransformations converts Terraform's
// TypeList and TypeSet collections into a []rules.Transformation.
func expandTransformations(v interface{}) ([]rules.Transformation, error) {
	ts, err := helper.ConvertTFCollectionToSlice(v)

	if err != nil {
		return nil, fmt.Errorf("error converting Transformations: %w", err)
	}

	result := make([]rules.Transformation, len(ts))
	for i, v := range ts {
		if s, ok := v.(string); ok {
			result[i] = rules.ConvertToTransformation(s)
		} else {
			return nil,
				fmt.Errorf("transformation was not a string: %+v", v)
		}
	}

	return result, nil
}

// expandChainedRules converts the Terraform representation of Chained Rules
// into the ChainedRule API Model
func expandChainedRules(attr interface{}) (*[]rules.ChainedRule, error) {

	if items, ok := attr.([]interface{}); ok {
		chainedRules := make([]rules.ChainedRule, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			chainedRule := rules.ChainedRule{}

			actionMap, err := helper.ConvertSingletonSetToMap(curr["action"])
			if err != nil {
				return nil, err
			}

			operatorMap, err := helper.ConvertSingletonSetToMap(curr["operator"])
			if err != nil {
				return nil, err
			}

			variables, err := expandVariables(curr["variable"])
			if err != nil {
				return nil, err
			}
			chainedRule.Variables = *variables

			if actionId, ok := actionMap["id"]; ok {
				chainedRule.Action.ID = actionId.(string)
			}

			if actionMsg, ok := actionMap["msg"]; ok {
				chainedRule.Action.Message = actionMsg.(string)
			}

			if actionT, ok := actionMap["transformations"]; ok {
				arr, err := expandTransformations(actionT)
				if err == nil {
					chainedRule.Action.Transformations = arr
				} else {
					return nil,
						fmt.Errorf("error reading transformations: %w", err)
				}
			}

			if v, ok := operatorMap["is_negated"]; ok {
				chainedRule.Operator.IsNegated = v.(bool)
			}

			if operatorType, ok := operatorMap["type"]; ok {
				if s, ok := operatorType.(string); ok {
					chainedRule.Operator.Type = rules.ConvertToOperatorType(s)
				} else {
					return nil, fmt.Errorf("operator type is not a string")
				}
			}

			if operatorValue, ok := operatorMap["value"]; ok {
				chainedRule.Operator.Value = operatorValue.(string)
			}

			chainedRules = append(chainedRules, chainedRule)
		}

		return &chainedRules, nil

	} else {
		return nil,
			errors.New("ExpandChainedRules: attr input was not a []interface{}")
	}
}

// expandVariables converts the Terraform representation of Variables into
// the Variable API Model
func expandVariables(attr interface{}) (*[]rules.Variable, error) {

	if items, ok := attr.([]interface{}); ok {

		variables := make([]rules.Variable, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			variable := rules.Variable{
				Type:    rules.ConvertToVariableType(curr["type"].(string)),
				IsCount: curr["is_count"].(bool),
			}

			matches, err := expandMatches(curr["match"])
			if err != nil {
				return nil, err
			}
			variable.Matches = *matches

			variables = append(variables, variable)
		}

		return &variables, nil

	} else {
		return nil,
			errors.New("ExpandVariables: attr input was not a []interface{}")
	}
}

// expandMatches converts the Terraform representation of Matches into
// the Match API Model
func expandMatches(attr interface{}) (*[]rules.Match, error) {
	if items, ok := attr.([]interface{}); ok {
		matches := make([]rules.Match, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			match := rules.Match{
				IsNegated: curr["is_negated"].(bool),
				IsRegex:   curr["is_regex"].(bool),
				Value:     curr["value"].(string),
			}
			matches = append(matches, match)
		}

		return &matches, nil

	} else {
		return nil,
			errors.New("ExpandMatches: attr input was not a []interface{}")
	}
}

// flattenSecRule converts the SecRule API Model
// into a format that Terraform can work with
func flattenSecRule(secrule rules.SecRule) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	m["name"] = secrule.Name
	m["action"] = flattenAction(secrule.Action)
	m["chained_rule"] = flattenChainedRules(secrule.ChainedRules)
	m["operator"] = flattenOperator(secrule.Operator)
	m["variable"] = flattenVariable(secrule.Variables)

	flattened = append(flattened, m)

	return flattened
}

// FlattenAction converts the Action API Model
// into a format that Terraform can work with
func flattenAction(action rules.Action) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	if action.ID != "" {
		m["id"] = action.ID
	}

	if action.Message != "" {
		m["msg"] = action.Message
	}

	ft := make([]string, 0)

	for _, t := range action.Transformations {
		ft = append(ft, t.String())
	}

	m["transformations"] = ft

	flattened = append(flattened, m)

	return flattened
}

// FlattenChainrule converts the ChainedRule API Model
// into a format that Terraform can work with
func flattenChainedRules(
	chainedRules []rules.ChainedRule,
) []map[string]interface{} {

	flattened := make([]map[string]interface{}, 0)

	for _, cg := range chainedRules {
		m := make(map[string]interface{})
		m["action"] = flattenAction(cg.Action)
		m["operator"] = flattenOperator(cg.Operator)
		m["variable"] = flattenVariable(cg.Variables)
		flattened = append(flattened, m)
	}

	return flattened
}

// FlattenAction converts the Operator API Model
// into a format that Terraform can work with
func flattenOperator(operator rules.Operator) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	m["is_negated"] = operator.IsNegated
	m["type"] = operator.Type.String()
	m["value"] = operator.Value

	flattened = append(flattened, m)

	return flattened
}

// FlattenVariable converts the Variable API Model
// into a format that Terraform can work with
func flattenVariable(variables []rules.Variable) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range variables {
		m := make(map[string]interface{})

		m["type"] = v.Type.String()
		m["match"] = flattenMatch(v.Matches)
		m["is_count"] = v.IsCount

		flattened = append(flattened, m)
	}

	return flattened
}

// FlattenMatch converts the Match API Model
// into a format that Terraform can work with
func flattenMatch(matches []rules.Match) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)

	for _, v := range matches {
		m := make(map[string]interface{})

		m["is_negated"] = v.IsNegated
		m["is_regex"] = v.IsRegex
		m["value"] = v.Value

		flattened = append(flattened, m)
	}

	return flattened
}
