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

	// TODO: make an array of testcases
	queryParameters := struct {
		Certifier       string
		CertificateType string
		// pagination defines an optional pagination for the request.
		Pagination *query.PageRequest
	}{
		Certifier:       suite.address[0].String(),
		CertificateType: "auditing",
		Pagination:      &query.PageRequest{Offset: 1},
	}

	queryResponse, err := queryClient.Certificates(ctx.Context(), &types.QueryCertificatesRequest{Certifier: queryParameters.Certifier, CertificateType: queryParameters.CertificateType, Pagination: queryParameters.Pagination})
	suite.Require().NoError(err)
	// TODO: add the hardcoded value in testcases
	suite.Require().Equal(2, int(queryResponse.Total))

}
