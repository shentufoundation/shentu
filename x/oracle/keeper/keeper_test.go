package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	"github.com/stretchr/testify/suite"

	"github.com/cometbft/cometbft/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"

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
	suite.ctx = suite.app.BaseApp.NewContext(false)
	suite.keeper = suite.app.OracleKeeper
	suite.params = suite.keeper.GetLockedPoolParams(suite.ctx)
	suite.minCollateral = suite.keeper.GetLockedPoolParams(suite.ctx).MinimumCollateral // 0.01 CTK

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.OracleKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := testutil.FundAccount(
			suite.ctx,
			suite.app.BankKeeper,
			acc,
			sdk.NewCoins(
				sdk.NewCoin("stake", math.NewInt(10000000000)),
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
