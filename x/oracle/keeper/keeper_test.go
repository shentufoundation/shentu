package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite

	// cdc    *codec.LegacyAmino
	app           *shentuapp.ShentuApp
	ctx           sdk.Context
	params        types.LockedPoolParams
	keeper        keeper.Keeper
	address       []sdk.AccAddress
	minCollateral int64
	queryClient   types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.OracleKeeper
	suite.params = suite.keeper.GetLockedPoolParams(suite.ctx)
	suite.minCollateral = suite.keeper.GetLockedPoolParams(suite.ctx).MinimumCollateral // 0.01 CTK

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.OracleKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := testutil.FundAccount(
			suite.app.BankKeeper,
			suite.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin("uctk", sdk.NewInt(10000000000)), // 1,000 CTK
			),
		)
		if err != nil {
			panic(err)
		}
	}

	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
