package keeper_test

import (
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

	bz, err := querier(ctx, []string{"other"}, query)
	suite.Require().Error(err)
	suite.Require().Nil(bz)

	path := []string{types.QueryCertifier, suite.address[0].String()}

	bz, err = querier(ctx, path, query)
	suite.Require().NoError(err)
	suite.Require().NotNil(bz)

}

func (suite *KeeperTestSuite) TestQuerier_QueryCertifiers() {
	app, ctx := suite.app, suite.ctx

	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCertifiers),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())

	bz, err := querier(ctx, []string{"other"}, query)
	suite.Require().Error(err)
	suite.Require().Nil(bz)

	path := []string{types.QueryCertifiers}

	bz, err = querier(ctx, path, query)
	suite.Require().NoError(err)
	suite.Require().NotNil(bz)

}

func (suite *KeeperTestSuite) TestQuerier_QueryCertifierByAlias() {
	app, ctx := suite.app, suite.ctx
	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCertifierByAlias),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())

	bz, err := querier(ctx, []string{"other"}, query)
	suite.Require().Error(err)
	suite.Require().Nil(bz)

	path := []string{types.QueryCertifierByAlias, "address1"}

	bz, err = querier(ctx, path, query)
	suite.Require().NoError(err)
	suite.Require().NotNil(bz)

}

// func (suite *KeeperTestSuite) TestQuerier_QueryPlatform() {
// 	app, ctx := suite.app, suite.ctx
// 	query := abci.RequestQuery{
// 		Path: "",
// 		Data: []byte{},
// 	}

// 	querier := keeper.NewQuerier(app.CertKeeper, app.LegacyAmino())

// 	bz, err := querier(ctx, []string{"other"}, query)
// 	suite.Require().Error(err)
// 	suite.Require().Nil(bz)

// 	// TODO: Clarification about input, Should we test queryPlatform?
// 	_, pubkey, _ := testdata.KeyTestPubAddr()
// 	path := []string{types.QueryPlatform, pubkey.String()}

// 	bz, err = querier(ctx, path, query)
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(bz)

// }
