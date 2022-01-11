package keeper_test

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	qtypes "github.com/cosmos/cosmos-sdk/types/query"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/certikfoundation/shentu/v2/x/cert/keeper"
	"github.com/certikfoundation/shentu/v2/x/cert/types"
)

func (suite *KeeperTestSuite) TestQuerier_QueryCertifier() {
	app, ctx := suite.app, suite.ctx
	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCertifier),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())

	tests := []struct {
		path       []string
		shouldPass bool
	}{
		{
			path:       []string{"other"},
			shouldPass: false,
		},
		{
			path:       []string{types.QueryCertifier, suite.address[0].String()},
			shouldPass: true,
		},
	}

	for _, tc := range tests {
		bz, err := querier(ctx, tc.path, query)
		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(bz)
			var jsonResponse map[string]interface{}
			json.Unmarshal(bz, &jsonResponse)
			suite.Require().Equal(tc.path[1], jsonResponse["address"])
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(bz)
		}
	}
}

func (suite *KeeperTestSuite) TestQuerier_QueryCertifiers() {
	app, ctx := suite.app, suite.ctx

	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCertifiers),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())

	tests := []struct {
		path       []string
		certifier  string
		shouldPass bool
	}{
		{
			path:       []string{"other"},
			certifier:  suite.address[0].String(),
			shouldPass: false,
		},
		{
			path:       []string{types.QueryCertifiers},
			certifier:  suite.address[0].String(),
			shouldPass: true,
		},
	}

	for _, tc := range tests {
		bz, err := querier(ctx, tc.path, query)
		if tc.shouldPass {
			var jsonResponse map[string][]map[string]interface{}
			json.Unmarshal(bz, &jsonResponse)
			suite.Require().Equal(tc.certifier, jsonResponse["certifiers"][0]["address"])
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(bz)
		}
	}
}

func (suite *KeeperTestSuite) TestQuerier_QueryCertifierByAlias() {
	app, ctx := suite.app, suite.ctx
	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCertifierByAlias),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())

	tests := []struct {
		path       []string
		shouldPass bool
	}{
		{
			path:       []string{"other"},
			shouldPass: false,
		},
		{
			path:       []string{types.QueryCertifierByAlias, "address1"},
			shouldPass: true,
		},
	}

	for _, tc := range tests {
		bz, err := querier(ctx, tc.path, query)
		if tc.shouldPass {
			var jsonResponse map[string]interface{}
			json.Unmarshal(bz, &jsonResponse)
			suite.Require().Equal(tc.path[1], jsonResponse["alias"])
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(bz)
		}
	}
}

func (suite *KeeperTestSuite) TestQuerier_QueryCertificate() {
	app, ctx := suite.app, suite.ctx
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
	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCertificate),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())

	tests := []struct {
		path       []string
		shouldPass bool
	}{
		{
			path:       []string{"other"},
			shouldPass: false,
		},
		{
			path:       []string{types.QueryCertificate, "0"},
			shouldPass: false,
		},
		{
			path:       []string{types.QueryCertificate, "1"},
			shouldPass: true,
		},
	}

	for _, tc := range tests {
		bz, err := querier(ctx, tc.path, query)
		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(bz)
			var jsonResponse map[string]interface{}
			json.Unmarshal(bz, &jsonResponse)
			suite.Require().Equal(tc.path[1], jsonResponse["certificate_id"])
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(bz)
		}
	}
}

// TODO: Error - invalid character -- Amino:JSON int/int64/uint/uint64 expects quoted values for javascript numeric support, got: 1
func (suite *KeeperTestSuite) TestQuerier_QueryCertificates() {
	app, ctx := suite.app, suite.ctx
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
		path            []string
		certifier       sdk.AccAddress
		certificateType string
		// pagination defines an optional pagination for the request.
		pagination        *query.PageRequest
		totalCertificates int
		shouldPass        bool
	}{
		{
			// invalid path
			path:              []string{"other"},
			certifier:         suite.address[0],
			certificateType:   "auditing",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 2,
			shouldPass:        false,
		},
		{
			path:              []string{types.QueryCertificates},
			certifier:         suite.address[0],
			certificateType:   "auditing",
			pagination:        &query.PageRequest{Offset: 1},
			totalCertificates: 2,
			shouldPass:        true,
		}}
	for _, tc := range tests {
		page, limit, err := qtypes.ParsePagination(tc.pagination)
		suite.Require().NoError(err)
		queryBytes, _ := json.Marshal(&types.QueryCertificatesParams{Certifier: tc.certifier, CertificateType: types.CertificateTypeFromString(tc.certificateType), Page: page, Limit: limit})
		query := abci.RequestQuery{
			Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCertificates),
			Data: queryBytes,
		}
		querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())
		bz, err := querier(ctx, tc.path, query)
		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(bz)
			var jsonResponse map[string]interface{}
			json.Unmarshal(bz, &jsonResponse)
			suite.Require().Equal(tc.totalCertificates, jsonResponse["total"])
		} else {
			suite.Require().Error(err)
		}
	}
}
