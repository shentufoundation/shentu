package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/certikfoundation/shentu/v2/x/shield/keeper"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

// shared setup
type PurchaseTestSuite struct {
	suite.Suite

	app                   *simapp.SimApp
	ctx                   sdk.Context
	keeper                keeper.Keeper
	address               []sdk.AccAddress
	queryClient           types.QueryClient
	shieldAdminAccAddress sdk.AccAddress
}

func (suite *PurchaseTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
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
	_, err := suite.keeper.CreatePool(suite.ctx, types.MsgCreatePool{From: suite.shieldAdminAccAddress.String(), SponsorAddr: suite.address[1].String(), Description: "", ShieldRate: sdk.NewDec(1)})
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
