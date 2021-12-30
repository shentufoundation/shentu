package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"


	"github.com/certikfoundation/shentu/v2/simapp"
)

func TestFundCommunityPool(t *testing.T) {
	t.Log("Test keeper FundCommunityPool")
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	coins := sdk.Coins{sdk.NewInt64Coin("uctk", 80000*1e6)}
	require.NoError(t, sdksimapp.FundModuleAccount(app.BankKeeper, ctx, "mint", coins))

	t.Run("Funding community pool", func(t *testing.T) {
		moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte("mint")))
		coins100 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 100*1e6))
		bal0 := app.BankKeeper.GetBalance(ctx, moduleAcct, "uctk")
		err := app.MintKeeper.SendToCommunityPool(ctx, coins100)
		bal1 := app.BankKeeper.GetBalance(ctx, moduleAcct, "uctk")
		require.NoError(t, err)
		require.Equal(t, bal0.Sub(bal1), sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 100*1e6))
	})
}