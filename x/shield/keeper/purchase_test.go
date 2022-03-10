package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/certikfoundation/shentu/v2/app"
	"github.com/certikfoundation/shentu/v2/x/shield/keeper"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// shared setup
type PurchaseTestSuite struct {
	suite.Suite

	app                   *shentuapp.ShentuApp
	ctx                   sdk.Context
	keeper                keeper.Keeper
	address               []sdk.AccAddress
	queryClient           types.QueryClient
	shieldAdminAccAddress sdk.AccAddress
}

func (suite *PurchaseTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.ShieldKeeper
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.ShieldKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}

	// Set address[3] as admin
	suite.app.ShieldKeeper.SetAdmin(suite.ctx, suite.address[3])
	suite.shieldAdminAccAddress = suite.address[3]

	// create pool
	pool := types.NewMsgCreatePool(
		suite.shieldAdminAccAddress,
		suite.address[1],
		"",
		sdk.NewDec(1),
		sdk.NewInt(1e10),
	)
	_, err := suite.keeper.CreatePool(suite.ctx, *pool)
	suite.Require().NoError(err)
}

func (suite *PurchaseTestSuite) TestPurchaseShield() {
	tests := []struct {
		name        string
		description string
		purchaser   sdk.AccAddress
		amount      sdk.Coins
		poolID      uint64
		shouldPass  bool
	}{
		{
			name:        "Sufficient purchase amount, valid PoolID",
			description: "",
			purchaser:   suite.address[0],
			amount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (50)*1e6)),
			poolID:      1,
			shouldPass:  true,
		},
		{
			name:        "Insufficient purchase amount",
			description: "",
			purchaser:   suite.address[0],
			amount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (10)*1e6)),
			poolID:      1,
			shouldPass:  false,
		},
		{
			name:        "Invalid poolID",
			description: "",
			purchaser:   suite.address[0],
			amount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (50)*1e6)),
			poolID:      10,
			shouldPass:  false,
		},
	}

	for _, tc := range tests {
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, tc.purchaser, tc.amount))
		_, err := suite.keeper.PurchaseShield(suite.ctx, tc.poolID, tc.amount, tc.description, tc.purchaser)

		if tc.shouldPass {
			suite.Require().NoError(err)
		} else {
			suite.Require().Error(err)
		}
	}
}

func TestPurchaseTestSuite(t *testing.T) {
	suite.Run(t, new(PurchaseTestSuite))
}
