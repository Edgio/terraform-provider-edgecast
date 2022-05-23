package data

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

func createDNSData(cfg edgecast.SDKConfig) (groupID, masterServerGroupID, masterServerA, masterServerB, secondaryServerGroupID, tsgID, zoneID int) {
	svc := internal.Check(routedns.New(cfg))

	tsgID = createTSIG(svc)
	masterServerGroupID, masterServerA, masterServerB = createMasterServerGroup(svc)
	secondaryServerGroupID = createSecondaryServerGroup(svc, tsgID, masterServerGroupID, masterServerA, masterServerB)
	/* TODO Repair
	----
	groupID = createGroup(svc)
	zoneID = createZone(svc)
	----
	*/
	return
}

func createZone(svc *routedns.RouteDNSService) int {
	params := routedns.AddZoneParams{
		AccountNumber: account(),
		Zone: routedns.Zone{
			DomainName: "test.edgecast.com",
			Status:     1,
			ZoneType:   1,
			Records:    routedns.DNSRecords{},
			Groups: []routedns.DnsRouteGroup{
				{
					Name:             "",
					GroupTypeID:      1,
					GroupProductType: 1,
					GroupComposition: routedns.DNSGroupRecords{},
				},
			},
		},
	}

	return *internal.Check(svc.AddZone(params))
}

func createTSIG(svc *routedns.RouteDNSService) (id int) {
	params := routedns.AddTSIGParams{
		AccountNumber: account(),
		TSIG: routedns.TSIG{
			Alias:       "test",
			KeyName:     "key3",
			KeyValue:    "0011223344552",
			AlgorithmID: routedns.HMAC_SHA512,
		},
	}

	return *internal.Check(svc.AddTSIG(params))
}

func createMasterServerGroup(svc *routedns.RouteDNSService) (id, msA, msB int) {
	params := &routedns.AddMasterServerGroupParams{
		AccountNumber: account(),
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

func createSecondaryServerGroup(svc *routedns.RouteDNSService, tsgID, msgID, serverA, serverB int) (id int) {
	params := &routedns.AddSecondaryZoneGroupParams{
		AccountNumber: account(),

		SecondaryZoneGroup: routedns.SecondaryZoneGroup{
			Name: unique("SZG"),
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
						DomainName: unique("second49.com"),
						Status:     1,
					},
				},
			},
		},
	}

	szg := internal.Check(svc.AddSecondaryZoneGroup(*params))
	return szg.ID
}

func createGroup(routeDNSService *routedns.RouteDNSService) (groupID int) {
	group := routedns.DnsRouteGroup{
		Name:             "DNS GROUP",
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
	params.AccountNumber = account()
	params.Group = group

	return *internal.Check(routeDNSService.AddGroup(*params))
}
