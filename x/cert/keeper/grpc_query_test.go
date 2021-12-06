package keeper_test

import (
	"github.com/certikfoundation/shentu/v2/x/cert/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (suite *KeeperTestSuite) TestQueryCertifier() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// empty address string
	_, err := queryClient.Certifier(ctx.Context(), &types.QueryCertifierRequest{})
	suite.Require().Error(err)

	// valid address
	queryResponse, err := queryClient.Certifier(ctx.Context(), &types.QueryCertifierRequest{Address: suite.address[0].String(), Alias: ""})
	suite.Require().NoError(err)
	suite.Equal(suite.address[0].String(), queryResponse.Certifier.Address)
}

func (suite *KeeperTestSuite) TestQueryCertifiers() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// TODO: Need clarification
	// _, err := queryClient.Certifiers(ctx.Context(), nil)
	// suite.Require().Error(err)

	// valid request
	_, err := queryClient.Certifiers(ctx.Context(), &types.QueryCertifiersRequest{})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestQueryCertificate() {
	ctx, queryClient := suite.ctx, suite.queryClient
	cert := cert{
		certTypeStr:  "compilation",
		contStr:      "sourcodehash0",
		compiler:     "compiler1",
		bytecodeHash: "bytecodehash1",
		description:  "",
		certifier:    suite.address[0],
	}
	certificate, err := types.NewCertificate(cert.certTypeStr, cert.contStr, cert.compiler, cert.bytecodeHash, cert.description, cert.certifier)
	suite.Require().NoError(err)
	suite.keeper.SetNextCertificateID(suite.ctx, 1)
	suite.keeper.IssueCertificate(ctx, certificate)

	// invalid request
	// _, err = queryClient.Certificate(ctx.Context(), nil)
	// suite.Require().Error(err)

	// id not found
	_, err = queryClient.Certificate(ctx.Context(), &types.QueryCertificateRequest{CertificateId: 10})
	suite.Require().Error(err)

	// valid request
	_, err = queryClient.Certificate(ctx.Context(), &types.QueryCertificateRequest{CertificateId: 1})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestQueryCertificates() {
	ctx, queryClient := suite.ctx, suite.queryClient

	allCertificates := []cert{
		{
			// type auditing, certifier suite.address[0]
			certTypeStr:  "auditing",
			contStr:      "sourcodehash0",
			compiler:     "compiler1",
			bytecodeHash: "bytecodehash1",
			description:  "",
			certifier:    suite.address[0],
			inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx),
		},
		{
			// type auditing, but different certifier
			certTypeStr:  "auditing",
			contStr:      "sourcodehash0",
			compiler:     "compiler1",
			bytecodeHash: "bytecodehash1",
			description:  "",
			certifier:    suite.address[1],
			inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx) + 1,
		},
		{
			// different type
			certTypeStr:  "compilation",
			contStr:      "sourcodehash0",
			compiler:     "compiler1",
			bytecodeHash: "bytecodehash1",
			description:  "",
			certifier:    suite.address[0],
			inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx) + 2,
		},
		{
			certTypeStr:  "auditing",
			contStr:      "sourcodehash0",
			compiler:     "compiler1",
			bytecodeHash: "bytecodehash1",
			description:  "",
			certifier:    suite.address[0],
			inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx) + 3,
		},
	}

	for _, cert := range allCertificates {
		want, err := types.NewCertificate(cert.certTypeStr, cert.contStr, cert.compiler, cert.bytecodeHash, cert.description, cert.certifier)
		suite.Require().NoError(err)
		// TODO: maybe remove this and just use issueCertificate
		want.CertificateId = cert.inputCertId
		// set the cert
		suite.keeper.SetCertificate(suite.ctx, want)
	}

	tests := []struct {
		certifier       string
		certificateType string
		// pagination defines an optional pagination for the request.
		pagination        *query.PageRequest
		totalCertificates int
		shouldPass        bool
	}{
		{
			certifier:         suite.address[0].String(),
			certificateType:   "auditing",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 2,
			shouldPass:        true,
		},
		{
			certifier:         suite.address[0].String(),
			certificateType:   "compilation",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 1,
			shouldPass:        true,
		},
		{
			certifier:         suite.address[1].String(),
			certificateType:   "auditing",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 1,
			shouldPass:        true,
		},
		{
			certifier:         suite.address[0].String(),
			certificateType:   "auditing",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 0,
			shouldPass:        false,
		},
	}

	for _, tc := range tests {
		queryResponse, err := queryClient.Certificates(ctx.Context(), &types.QueryCertificatesRequest{Certifier: tc.certifier, CertificateType: tc.certificateType, Pagination: tc.pagination})
		suite.Require().NoError(err)
		if tc.shouldPass {
			suite.Require().Equal(tc.totalCertificates, int(queryResponse.Total))
		} else {
			suite.Require().NotEqual(tc.totalCertificates, int(queryResponse.Total))
		}

	}

}
