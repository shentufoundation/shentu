package keeper_test

import (
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

type TestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	accounts    []sdk.AccAddress
	vals        []stakingtypes.Validator
	queryClient types.QueryClient
}

func setup() TestSuite {
	var t TestSuite
	t.app = simapp.Setup(false)
	t.ctx = t.app.BaseApp.NewContext(false, tmproto.Header{})
	t.keeper = t.app.ShieldKeeper
	t.accounts = []sdk.AccAddress{acc1, acc2, acc3, acc4}

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := t.app.BankKeeper.AddCoins(
			t.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin(t.app.StakingKeeper.BondDenom(t.ctx), sdk.NewInt(10000000000)), // 1,000 stake
			),
		)
		if err != nil {
			panic(err)
		}
	}
	return t
}

func OneMixedCoins(nativeDenom string) types.MixedCoins {
	native := sdk.NewCoins(sdk.NewCoin(nativeDenom, sdk.NewInt(1)))
	foreign := sdk.NewCoins(sdk.NewCoin("dummy", sdk.NewInt(1)))
	return types.MixedCoins{
		Native:  native,
		Foreign: foreign,
	}
}

func OneMixedDecCoins(nativeDenom string) types.MixedDecCoins {
	native := sdk.NewDecCoins(sdk.NewDecCoin(nativeDenom, sdk.NewInt(1)))
	foreign := sdk.NewDecCoins(sdk.NewDecCoin("dummy", sdk.NewInt(1)))
	return types.MixedDecCoins{
		Native:  native,
		Foreign: foreign,
	}
}

func DummyPool() types.Pool {
	return types.Pool{
		Id:          1,
		Description: "w",
		Sponsor:     acc1.String(),
		SponsorAddr: acc2.String(),
		ShieldLimit: sdk.NewInt(1),
		Active:      false,
		Shield:      sdk.NewInt(1),
	}
}
