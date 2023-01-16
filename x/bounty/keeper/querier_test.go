package keeper_test

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestQueries() {
	var (
		programPath  = []string{types.QueryProgram}
		programsPath = []string{types.QueryPrograms}
	)

	app, ctx := suite.app, suite.ctx
	legacyQuerierCdc := app.LegacyAmino()
	querier := keeper.NewQuerier(app.BountyKeeper, legacyQuerierCdc)

	// default
	req := abci.RequestQuery{
		Data: []byte{},
		Path: "",
	}
	bz, err := querier(ctx, []string{"other"}, req)
	suite.Require().Error(err)
	suite.Require().Nil(bz)

	// program
	queryProgramParams := &types.QueryProgramParams{ProgramID: 1}
	req = abci.RequestQuery{
		Data: legacyQuerierCdc.MustMarshalJSON(queryProgramParams),
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProgram),
	}
	bz, err = querier(ctx, programPath, req)
	suite.Require().Error(err)
	suite.Require().Nil(bz)

	// create programs
	//suite.CreatePrograms()

	// programs
	queryProgramsParams := &types.QueryProgramsParams{
		Page:  1,
		Limit: 100,
	}
	req = abci.RequestQuery{
		Data: legacyQuerierCdc.MustMarshalJSON(queryProgramsParams),
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPrograms),
	}

	bz, err = querier(ctx, programsPath, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(bz)
}
