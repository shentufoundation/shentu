package testshield

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield"
	"github.com/certikfoundation/shentu/x/shield/types"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t *testing.T
	h sdk.Handler
	k shield.Keeper

	Ctx        sdk.Context
	Denom      string
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k shield.Keeper, denom string) *Helper {
	return &Helper{t, shield.NewHandler(k), k, ctx, denom}
}

func (sh *Helper) DepositCollateral(addr sdk.AccAddress, amount int64, ok bool) {
	coins := sdk.NewCoins(sdk.NewInt64Coin(sh.Denom, amount))
	msg := types.NewMsgDepositCollateral(addr, coins)
	sh.Handle(msg, ok)
}

func (sh *Helper) WithdrawCollateral(addr sdk.AccAddress, amount int64, ok bool) {
	coins := sdk.NewCoins(sdk.NewInt64Coin(sh.Denom, amount))
	msg := types.NewMsgWithdrawCollateral(addr, coins)
	sh.Handle(msg, ok)
}

func (sh *Helper) CreatePool(addr, sponsorAddr sdk.AccAddress, nativeDeposit, shield, shieldLimit int64, sponsor, description, bondDenom string) {
	shieldCoins := sdk.NewCoins(sdk.NewInt64Coin(bondDenom, shield))
	depositCoins := types.MixedCoins{Native: sdk.NewCoins(sdk.NewInt64Coin(bondDenom, nativeDeposit))}
	limit := sdk.NewInt(shieldLimit)
	msg := types.NewMsgCreatePool(addr, shieldCoins, depositCoins, sponsor, sponsorAddr, description, limit)
	sh.Handle(msg, true)
}

func (sh *Helper) PurchaseShield(purchaser sdk.AccAddress, shield, poolID int64, ok bool) {
	shieldCoins := sdk.NewCoins(sdk.NewInt64Coin(bondDenom, shield))
	msg := types.NewMsgPurchaseShield(poolID, shieldCoins, "test_purchase", purchaser)
	sh.Handle(msg, ok)
}

// Handle calls shield handler on a given message
func (sh *Helper) Handle(msg sdk.Msg, ok bool) *sdk.Result {
	res, err := sh.h(sh.Ctx, msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
	return res
}
