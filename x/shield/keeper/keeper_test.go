package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	//tmproto "github.com/tendermint/tendermint/proto/types"

	//"github.com/cosmos/cosmos-sdk/simapp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	//"github.com/cosmos/cosmos-sdk/x/staking/teststaking"

	"github.com/certikfoundation/shentu/common/tests"
	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/staking/teststaking"
)

// TestWithdraw tests withdraws triggered by staking undelegation.
func TestWithdraw(t *testing.T) {
	//_ = cosmosSimApp.Setup(false)
	app := simapp.Setup(false)
	//app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})
	//ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	p := app.StakingKeeper.GetParams(ctx)
	p.MaxValidators = 5
	app.StakingKeeper.SetParams(ctx, p)

	//addrDels := simapp.AddTestAddrsIncremental(app, ctx, 6, sdk.TokensFromConsensusPower(200))
	//valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)

	delAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]

	accAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]
	valAddr := sdk.ValAddress(accAddr)

	pubKey := tests.MakeTestPubKey()
	//accAddr := tests.MakeTestAccAddressFromPubKey(pubKey)
	//valAddr := tests.MakeTestValAddressFromPubKey(pubKey)
	
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// Set up a validator
	tstaking.CreateValidatorWithValPower(valAddr, pubKey, 100, true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)

	// Delegate some to the validator
	tstaking.CheckDelegator(delAddr, valAddr, false)
	tstaking.Delegate(delAddr, valAddr, 100)
	tstaking.CheckDelegator(delAddr, valAddr, true)


	// Deposit collateral
	// TODO: Create shield test helper
	err := app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(100))
	require.Nil(t, err)
	
	// Undelegate some
	//tstaking.Undelegate(sdk.AccAddress(valAddr), valAddr, sdk.TokensFromConsensusPower(1), true)
	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(50), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)
	tstaking.CheckDelegator(delAddr, valAddr, true)
	// TODO check amount

	fmt.Printf("\n asdfasdfasdfasdf \n")


	// Check shield withdrawal
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	withdraws := app.ShieldKeeper.GetAllWithdraws(ctx) // GetAllWithdraws NOT WORKING?
	fmt.Printf("\n WITHDRAWS: %v\n", withdraws)

	withdraws = app.ShieldKeeper.GetAllWithdrawsRevised(ctx)
	fmt.Printf("\n AGAIN WITHDRAWS: %v\n", withdraws)
}
