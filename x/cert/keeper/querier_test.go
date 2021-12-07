package keeper_test

import (
	"encoding/json"
	"fmt"

	"github.com/certikfoundation/shentu/v2/x/cert/keeper"
	"github.com/certikfoundation/shentu/v2/x/cert/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
