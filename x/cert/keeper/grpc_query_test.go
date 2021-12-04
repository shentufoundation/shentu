package keeper_test

import (
	"github.com/certikfoundation/shentu/v2/x/cert/types"
)

func (suite *KeeperTestSuite) TestQueryCertifier() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// empty address string
	_, err := queryClient.Certifier(ctx.Context(), &types.QueryCertifierRequest{})
	suite.Require().Error(err)

	// valid address
	queryResponse, err := queryClient.Certifier(ctx.Context(), &types.QueryCertifierRequest{Address: string(suite.address[0]), Alias: ""})
	suite.Require().NoError(err)
	suite.Equal(acc1, queryResponse.Certifier.Address)
}

func (suite *KeeperTestSuite) TestQueryCertifiers() {
	ctx, queryClient := suite.ctx, suite.queryClient

	_, err := queryClient.Certifiers(ctx.Context(), nil)
	suite.Require().Error(err)

	// valid request
	_, err = queryClient.Certifiers(ctx.Context(), &types.QueryCertifiersRequest{})
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
	_, err = queryClient.Certificate(ctx.Context(), nil)
	suite.Require().Error(err)

	// id not found
	_, err = queryClient.Certificate(ctx.Context(), &types.QueryCertificateRequest{CertificateId: 10})
	suite.Require().Error(err)

	// valid request
	_, err = queryClient.Certificate(ctx.Context(), &types.QueryCertificateRequest{CertificateId: 1})
	suite.Require().NoError(err)
}
