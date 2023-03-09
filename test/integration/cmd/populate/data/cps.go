// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"time"

	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/certificate"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
)

const timestampFormat = "2006-01-02T15_04_05Z07_00"

func createCPSData(cfg config.Config) CPSResult {
	svc := internal.Check(cps.New(cfg.SDKConfig))
	id := createCertificate(svc)

	return CPSResult{
		CertificateID: id,
	}
}

func createCertificate(svc *cps.CpsService) int64 {
	certParams := certificate.NewCertificatePostParams()
	certParams.Certificate = &models.CertificateCreate{
		CertificateLabel:     "C_" + time.Now().Format(timestampFormat),
		Description:          "test cert",
		AutoRenew:            false,
		CertificateAuthority: "DigiCert",
		ValidationType:       "DV",
		DcvMethod:            "DnsTxtToken",
		Domains: []*models.DomainCreateUpdate{
			{
				IsCommonName: true,
				Name:         "testssdomain.com",
			},
		},
	}

	resp := internal.Check(svc.Certificate.CertificatePost(certParams))
	return resp.ID
}
