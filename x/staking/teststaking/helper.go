package teststaking

import (
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/x/staking/keeper"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t *testing.T
	h sdk.Handler
	k keeper.Keeper

	ctx        sdk.Context
	Commission stakingtypes.CommissionRates
	// Coin Denomination
	Denom string
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper) *Helper {
	return &Helper{t, staking.NewHandler(k.Keeper), k, ctx, ZeroCommission(), k.Keeper.BondDenom(ctx)}
}

// CreateValidator calls handler to create a new staking validator
func (sh *Helper) CreateValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, stakeAmount int64, ok bool) {
	coin := sdk.NewCoin(sh.Denom, sdk.NewInt(stakeAmount))
	sh.createValidator(addr, pk, coin, ok)
}

// CreateValidatorWithValPower calls handler to create a new staking validator with zero
// commission
func (sh *Helper) CreateValidatorWithValPower(addr sdk.ValAddress, pk cryptotypes.PubKey, valPower int64, ok bool) sdk.Int {
	amount := sdk.TokensFromConsensusPower(valPower)
	coin := sdk.NewCoin(sh.Denom, amount)
	sh.createValidator(addr, pk, coin, ok)
	return amount
}

func (sh *Helper) createValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, coin sdk.Coin, ok bool) {
	msg, err := stakingtypes.NewMsgCreateValidator(addr, pk, coin, stakingtypes.Description{}, sh.Commission, sdk.OneInt())
	if err != nil {
		panic(err)
	}
	sh.Handle(msg, ok)
}

// Delegate calls handler to delegate stake for a validator
func (sh *Helper) Delegate(delegator sdk.AccAddress, val sdk.ValAddress, amount int64) {
	coin := sdk.NewCoin(sh.Denom, sdk.NewInt(amount))
	msg := stakingtypes.NewMsgDelegate(delegator, val, coin)
	sh.Handle(msg, true)
}

// Redelegate calls handler to begin redelegation.
func (sh *Helper) Redelegate(delegator sdk.AccAddress, srcVal, dstVal sdk.ValAddress, amount int64, ok bool) {
	coin := sdk.NewCoin(sh.Denom, sdk.NewInt(amount))
	msg := stakingtypes.NewMsgBeginRedelegate(delegator, srcVal, dstVal, coin)
	sh.Handle(msg, ok)
}

// Undelegate calls handler to unbound some stake from a validator.
func (sh *Helper) Undelegate(delegator sdk.AccAddress, val sdk.ValAddress, amount int64, ok bool) *sdk.Result {
	unbondAmt := sdk.NewInt64Coin(sh.Denom, amount)
	msg := stakingtypes.NewMsgUndelegate(delegator, val, unbondAmt)
	return sh.Handle(msg, ok)
}

// Handle calls staking handler on a given message
func (sh *Helper) Handle(msg sdk.Msg, ok bool) *sdk.Result {
	res, err := sh.h(sh.ctx, msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
	return res
}

// CheckValidator asserts that a validor exists and has a given status (if status!="")
// and if has a right jailed flag.
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
