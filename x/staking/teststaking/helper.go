package teststaking

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/shentufoundation/shentu/v2/x/staking/keeper"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t       *testing.T
	msgSrvr stakingtypes.MsgServer
	k       keeper.Keeper

	ctx        sdk.Context
	Commission stakingtypes.CommissionRates
	// Coin Denomination
	Denom string
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper) *Helper {
	return &Helper{t, stakingkeeper.NewMsgServerImpl(k.Keeper), k, ctx, ZeroCommission(), k.Keeper.BondDenom(ctx)}
}

// CreateValidator calls handler to create a new staking validator
func (sh *Helper) CreateValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, stakeAmount math.Int, ok bool) {
	coin := sdk.NewCoin(sh.Denom, stakeAmount)
	sh.createValidator(addr, pk, coin, ok)
}

// CreateValidatorWithValPower calls handler to create a new staking validator with zero
// commission
func (sh *Helper) CreateValidatorWithValPower(addr sdk.ValAddress, pk cryptotypes.PubKey, valPower int64, ok bool) math.Int {
	amount := sh.k.TokensFromConsensusPower(sh.ctx, valPower)
	coin := sdk.NewCoin(sh.Denom, amount)
	sh.createValidator(addr, pk, coin, ok)
	return amount
}

func (sh *Helper) createValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, coin sdk.Coin, ok bool) {
	msg, err := stakingtypes.NewMsgCreateValidator(addr, pk, coin, stakingtypes.Description{}, sh.Commission, sdk.OneInt())
	require.NoError(sh.t, err)
	res, err := sh.msgSrvr.CreateValidator(sdk.WrapSDKContext(sh.ctx), msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
}

// Delegate calls handler to delegate stake for a validator
func (sh *Helper) Delegate(delegator sdk.AccAddress, val sdk.ValAddress, amount math.Int) {
	coin := sdk.NewCoin(sh.Denom, amount)
	msg := stakingtypes.NewMsgDelegate(delegator, val, coin)
	res, err := sh.msgSrvr.Delegate(sdk.WrapSDKContext(sh.ctx), msg)
	require.NoError(sh.t, err)
	require.NotNil(sh.t, res)
}

// Redelegate calls handler to begin redelegation.
func (sh *Helper) Redelegate(delegator sdk.AccAddress, srcVal, dstVal sdk.ValAddress, amount int64, ok bool) {
	coin := sdk.NewCoin(sh.Denom, sdk.NewInt(amount))
	msg := stakingtypes.NewMsgBeginRedelegate(delegator, srcVal, dstVal, coin)
	res, err := sh.msgSrvr.BeginRedelegate(sdk.WrapSDKContext(sh.ctx), msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
}

// Undelegate calls handler to unbind some stake from a validator.
func (sh *Helper) Undelegate(delegator sdk.AccAddress, val sdk.ValAddress, amount int64, ok bool) {
	unbondAmt := sdk.NewInt64Coin(sh.Denom, amount)
	msg := stakingtypes.NewMsgUndelegate(delegator, val, unbondAmt)
	res, err := sh.msgSrvr.Undelegate(sdk.WrapSDKContext(sh.ctx), msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
}

// CheckValidator asserts that a validor exists and has a given status (if status!="")
// and if it has a right jailed flag.
func (sh *Helper) CheckValidator(addr sdk.ValAddress, status stakingtypes.BondStatus, jailed bool) stakingtypes.Validator {
	v, ok := sh.k.GetValidator(sh.ctx, addr)
	require.True(sh.t, ok)
	require.Equal(sh.t, jailed, v.Jailed, "wrong Jalied status")
	if status >= 0 {
		require.Equal(sh.t, status, v.Status)
	}
	return v
}

// CheckDelegator asserts that a delegator exists
func (sh *Helper) CheckDelegator(delegator sdk.AccAddress, val sdk.ValAddress, found bool) {
	_, ok := sh.k.GetDelegation(sh.ctx, delegator, val)
	require.Equal(sh.t, ok, found)
}

// TurnBlock updates context and calls endblocker.
func (sh *Helper) TurnBlock(ctx sdk.Context) {
	sh.ctx = ctx
	staking.EndBlocker(sh.ctx, sh.k.Keeper)
}

// ZeroCommission constructs a commission rates with all zeros.
func ZeroCommission() stakingtypes.CommissionRates {
	return stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
}
