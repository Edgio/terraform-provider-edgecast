// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
)

func createDNSData(cfg config.Config) DNSResult {
	svc := internal.Check(routedns.New(cfg.SDKConfig))

	tsgID := createTSIG(svc, cfg.AccountNumber)

	masterServerGroupID, masterServerA, masterServerB :=
		createMasterServerGroup(svc, cfg.AccountNumber)

	secondaryServerGroupID := createSecondaryServerGroup(
		svc,
		tsgID,
		masterServerGroupID,
		masterServerA,
		masterServerB,
		cfg.AccountNumber)

	groupID := createGroup(svc, cfg.AccountNumber)
	zoneID := createZone(svc, cfg.AccountNumber)

	return DNSResult{
		GroupID:                groupID,
		MasterServerGroupID:    masterServerGroupID,
		MasterServerA:          masterServerA,
		MasterServerB:          masterServerB,
		SecondaryServerGroupID: secondaryServerGroupID,
		TsgID:                  tsgID,
		ZoneID:                 zoneID,
	}
}

func createZone(svc *routedns.RouteDNSService, accountNumber string) int {
	params := routedns.AddZoneParams{
		AccountNumber: accountNumber,
		Zone: routedns.Zone{
			DomainName: internal.Unique("devenbltt.com"),
			Status:     1,
			ZoneType:   1,
			Comment:    "SDK test zone 1",
			Records: routedns.DNSRecords{
				A: []routedns.DNSRecord{
					{
						Name:         "testarecord2",
						TTL:          300,
						Rdata:        "54.11.33.29",
						RecordTypeID: routedns.A,
					},
				},
			},
		},
	}

	return *internal.Check(svc.AddZone(params))
}

func createTSIG(svc *routedns.RouteDNSService, accountNumber string) (id int) {
	params := routedns.AddTSIGParams{
		AccountNumber: accountNumber,
		TSIG: routedns.TSIG{
			Alias:       "test",
			KeyName:     "key3",
			KeyValue:    "0011223344552",
			AlgorithmID: routedns.HMAC_SHA512,
		},
	}

	return *internal.Check(svc.AddTSIG(params))
}

func createMasterServerGroup(
	svc *routedns.RouteDNSService,
	accountNumber string,
) (id, msA, msB int) {
	params := &routedns.AddMasterServerGroupParams{
		AccountNumber: accountNumber,
		MasterServerGroup: routedns.MasterServerGroupAddRequest{
			Name: "msg3000",
			Masters: []routedns.MasterServer{
				{
					Name:      "data.test.com",
					IPAddress: "10.10.10.2",
				},
				{
					Name:      "data2.test.com",
					IPAddress: "10.10.10.3",
				},
			},
		},
	}

	msg := internal.Check(svc.AddMasterServerGroup(*params))
	return msg.MasterGroupID, msg.Masters[0].ID, msg.Masters[1].ID
}

func createSecondaryServerGroup(
	svc *routedns.RouteDNSService,
	tsgID,
	msgID,
	serverA,
	serverB int,
	accountNumber string,
) (id int) {
	params := &routedns.AddSecondaryZoneGroupParams{
		AccountNumber: accountNumber,

		SecondaryZoneGroup: routedns.SecondaryZoneGroup{
			Name: internal.Unique("SZG"),
			ZoneComposition: routedns.ZoneComposition{
				MasterGroupID: msgID,
				MasterServerTSIGs: []routedns.MasterServerTSIGIDs{
					{
						MasterServer: routedns.MasterServerID{
							ID: serverA,
						},
						TSIG: routedns.TSIGID{
							ID: tsgID,
						},
					},
					{
						MasterServer: routedns.MasterServerID{
							ID: serverB,
						},
						TSIG: routedns.TSIGID{
							ID: tsgID,
						},
					},
				},
				Zones: []routedns.SecondaryZone{
					{
						Comment:    "comment",
						DomainName: internal.Unique("second49.com"),
						Status:     1,
					},
				},
			},
		},
	}

	szg := internal.Check(svc.AddSecondaryZoneGroup(*params))
	return szg.ID
}

func buildLoadbalancedGroup(
	groupTypeID routedns.GroupType,
) routedns.DnsRouteGroup {
	// Load Balanced Group with Records
	cnameRecord1 := routedns.DNSRecord{
		Name:         "testcnamerecord1",
		TTL:          300,
		Rdata:        "lb1.sdkzone.com",
		RecordTypeID: routedns.CNAME,
		Weight:       50,
	}
	cnameRecord2 := routedns.DNSRecord{
		Name:         "testcnamerecord2",
		TTL:          300,
		Rdata:        "lb2.sdkzone.com",
		RecordTypeID: routedns.CNAME,
		Weight:       50,
	}

	lbGroupRecord1 := routedns.DNSGroupRecord{
		Record: cnameRecord1,
	}

	lbGroupRecord2 := routedns.DNSGroupRecord{
		Record: cnameRecord2,
	}

	lbGroupRecords := routedns.DNSGroupRecords{}
	lbGroupRecords.CNAME = append(
		lbGroupRecords.CNAME,
		lbGroupRecord1,
		lbGroupRecord2,
	)

	lbGroup := routedns.DnsRouteGroup{
		Name:             internal.Unique("sdklbgroup01"),
		GroupTypeID:      groupTypeID,
		GroupProductType: routedns.LoadBalancing,
		GroupComposition: lbGroupRecords,
	}

	return lbGroup
}

func createGroup(
	routeDNSService *routedns.RouteDNSService,
	accountNumber string,
) (groupID int) {
	_ = routedns.DnsRouteGroup{
		Name:             internal.Unique("DNS GROUP"),
		GroupTypeID:      routedns.CName,
		GroupProductType: routedns.LoadBalancing,
		GroupComposition: routedns.DNSGroupRecords{
			A:    nil,
			AAAA: nil,
			CNAME: []routedns.DNSGroupRecord{
				{
					Record: routedns.DNSRecord{
						Name:           "test.",
						TTL:            2400,
						Rdata:          "",
						Weight:         100,
						RecordTypeID:   routedns.CNAME,
						RecordTypeName: "CNAME",
					},
					HealthCheck: nil,
					IsPrimary:   true,
				},
			},
		},
	}

	params := routedns.NewAddGroupParams()
	params.AccountNumber = accountNumber
	params.Group = buildLoadbalancedGroup(routedns.CName)

	return *internal.Check(routeDNSService.AddGroup(*params))
}
