package keeper_test

import (
	"github.com/certikfoundation/shentu/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/certikfoundation/shentu/x/shield/keeper"
)

type TestSuite struct {
	app    *simapp.SimApp
	ctx    sdk.Context
	keeper keeper.Keeper
}

func setup() TestSuite {
	var t TestSuite
	t.app = simapp.Setup(false)
	t.ctx = t.app.BaseApp.NewContext(false, tmproto.Header{})
	t.keeper = t.app.ShieldKeeper
	return t
}
