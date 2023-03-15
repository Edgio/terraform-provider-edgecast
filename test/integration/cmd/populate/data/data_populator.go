// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"terraform-provider-edgecast/test/integration/cmd/populate/config"

	"github.com/kr/pretty"
)

type DataPopulator struct {
	Config config.Config
}

// PopulateResource determines whether a resource should be populated. If the
// PopulateOnly config is not set, then this function always returns true. If
// the PopulateOnly config is set, this function returns true if the
// resourceName is present in the populate flags
func (d DataPopulator) PopulateResource(resourceName string) bool {
	return len(d.Config.PopulateFlags) == 0 ||
		d.Config.PopulateFlags[resourceName]
}

func (d DataPopulator) Populate() {
	if len(d.Config.PopulateFlags) > 0 {
		fmt.Printf(
			"populating data for only: %v\n",
			pretty.Formatter(d.Config.PopulateFlags.ToList()))
	}

	result := PopulationResult{}

	if d.PopulateResource("customer") {
		fmt.Println("populating data: customer")
		result.Customer = createCustomerData(d.Config)
		fmt.Printf("Customer result: %v\n", result.Customer)
	} else {
		fmt.Println("skipping populate: customer")
	}

	if d.PopulateResource("origin") {
		fmt.Println("populating data: origin")
		result.Origin = createOriginData(d.Config)
		fmt.Printf("Origin result: %v\n", result.Origin)
	} else {
		fmt.Println("skipping populate: origin")
	}

	if d.PopulateResource("cname") {
		fmt.Println("populating data: cname")
		result.CNAME = createEdgeCnameData(d.Config)
		fmt.Printf("CNAME result: %v\n", result.CNAME)
	} else {
		fmt.Println("skipping populate: cname")
	}

	if d.PopulateResource("dns") {
		fmt.Println("populating data: dns")
		result.DNS = createDNSData(d.Config)
		fmt.Printf("DNS result: %v\n", result.DNS)
	} else {
		fmt.Println("skipping populate: dns")
	}

	if d.PopulateResource("waf") {
		fmt.Println("populating data: waf")
		result.WAF = createWAFData(d.Config)
		fmt.Printf("WAF result: %v\n", result.WAF)
	} else {
		fmt.Println("skipping populate: waf")
	}

	if d.PopulateResource("rules_engine") {
		fmt.Println("populating data: rules_engine")
		result.RulesEngine = createRulesEnginePolicyData(d.Config)
		fmt.Printf("Rules Engine result: %v\n", result.RulesEngine)
	} else {
		fmt.Println("skipping populate: rules_engine")
	}

	if d.PopulateResource("originv3") {
		fmt.Println("populating data: originv3")
		result.OriginV3 = createOriginV3Data(d.Config)
		fmt.Printf("OriginV3 result: %v\n", result.OriginV3)
	} else {
		fmt.Println("skipping populate: originv3")
	}

	if d.PopulateResource("cps") {
		fmt.Println("populating data: cps")
		result.CPS = createCPSData(d.Config)
		fmt.Printf("CPS result: %v\n", result.CPS)
	} else {
		fmt.Println("skipping populate: cps")
	}

	if d.PopulateResource("waf_botmanager") {
		fmt.Println("populating data: waf_botmanager")
		result.BotManager = createWAFBotManagerData(d.Config)
		fmt.Printf("WAF Bot Manager result: %v\n", result.BotManager)
	} else {
		fmt.Println("skipping populate: waf_botmanager")
	}

	// Write results to a file
	jsonBytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Printf("error marshaling results to JSON: %v\n", err)
		log.Printf("result: %# v\n", pretty.Formatter(result))
		os.Exit(1)
	}

	file := fmt.Sprintf("result-%s.json", time.Now().Format(timestampFormat))
	err = os.WriteFile(file, jsonBytes, 0o644)
	if err != nil {
		log.Printf("error writing results file: %v\n", err)
		log.Printf("result json: %s\n", string(jsonBytes))
		os.Exit(1)
	}

	workingDirectory, _ := os.Getwd()
	fmt.Println("result file: ", path.Join(workingDirectory, file))
}
