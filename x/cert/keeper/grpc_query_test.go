package keeper_test

import (
	"github.com/certikfoundation/shentu/v2/x/cert/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (suite *KeeperTestSuite) TestQueryCertifier() {
	ctx, queryClient := suite.ctx, suite.queryClient

	tests := []struct {
		address    string
		alias      string
		shouldPass bool
	}{
		{
			address:    "",
			alias:      "",
			shouldPass: false,
		},
		{
			address:    suite.address[0].String(),
			alias:      "",
			shouldPass: true,
		},
		{
			address:    "",
			alias:      "address1",
			shouldPass: true,
		},
		{
			address:    "",
			alias:      "invalid",
			shouldPass: false,
		},
	}

	for _, tc := range tests {
		queryResponse, err := queryClient.Certifier(ctx.Context(), &types.QueryCertifierRequest{Address: tc.address, Alias: tc.alias})
		if tc.shouldPass {
			suite.Require().NoError(err)
			if tc.address != "" {
				suite.Equal(tc.address, queryResponse.Certifier.Address)
			} else {
				suite.Equal(tc.alias, queryResponse.Certifier.Alias)
			}
		} else {
			suite.Require().Error(err)
		}
	}
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

	tests := []struct {
		certificateId uint
		shouldPass    bool
	}{
		{
			certificateId: 0,
			shouldPass:    false,
		},
		{
			certificateId: 1,
			shouldPass:    true,
		},
		{
			certificateId: 10,
			shouldPass:    false,
		},
	}

	for _, tc := range tests {
		_, err = queryClient.Certificate(ctx.Context(), &types.QueryCertificateRequest{CertificateId: uint64(tc.certificateId)})
		if tc.shouldPass {
			suite.Require().NoError(err)
		} else {
			suite.Require().Error(err)
		}
	}
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
			// type auditing, but unqualified certifier suite.address[1]
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

	// intitalize certificate ID
	suite.keeper.SetNextCertificateID(suite.ctx, 1)
	for _, cert := range allCertificates {
		want, err := types.NewCertificate(cert.certTypeStr, cert.contStr, cert.compiler, cert.bytecodeHash, cert.description, cert.certifier)
		suite.Require().NoError(err)
		suite.keeper.IssueCertificate(suite.ctx, want)
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
			certifier:         suite.address[0].String(),
			certificateType:   "auditing",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 0,
			shouldPass:        false,
		},
		{
			certifier:         suite.address[1].String(),
			certificateType:   "auditing",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 0,
			shouldPass:        true,
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
