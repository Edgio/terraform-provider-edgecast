package data

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/origin"
	"github.com/EdgeCast/ec-sdk-go/edgecast/shared/enums"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

func createOriginData(cfg edgecast.SDKConfig) (id int) {
	/*TODO: Repair
	--
	svc := internal.Check(origin.New(cfg))

	id = createOrigin(svc)
	*/
	return
}

func createOrigin(svc *origin.OriginService) int {
	params := origin.AddOriginParams{
		AccountNumber: account(),
		MediaTypeID:   enums.HttpLarge,
		Origin: origin.Origin{
			DirectoryName:   "www",
			FollowRedirects: false,
			HostHeader:      "home.edgecast.com:80",
			HTTPHostnames: []origin.Hostname{
				{
					Name:      "http://origin1.customer.com",
					IsPrimary: 1,
					Ordinal:   0,
				},
			},
			HTTPLoadBalancing:    "PF",
			HTTPSHostnames:       []origin.Hostname{},
			HTTPSLoadBalancing:   "",
			NetworkConfiguration: 1,
			ValidationURL:        "",
			ShieldPOPs: []origin.ShieldPOP{
				{
					Name:    "",
					POPCode: "MCO",
				},
				{
					Name:    "",
					POPCode: "NYC",
				},
			},
		},
	}

	return *internal.Check(svc.AddOrigin(params))
}
