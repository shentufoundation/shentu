package testshield

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/shield"
	"github.com/shentufoundation/shentu/v2/x/shield/keeper"
	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t  *testing.T
	h  sdk.Handler
	ph govtypes.Handler
	k  keeper.Keeper

	ctx   sdk.Context
	denom string
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper, denom string) *Helper {
	return &Helper{t, shield.NewHandler(k), shield.NewShieldClaimProposalHandler(k), k, ctx, denom}
}

func (sh *Helper) DepositCollateral(addr sdk.AccAddress, amount int64, ok bool) {
	coins := sdk.NewCoins(sdk.NewInt64Coin(sh.denom, amount))
	msg := types.NewMsgDepositCollateral(addr, coins)
	sh.Handle(msg, ok)
}

func (sh *Helper) WithdrawCollateral(addr sdk.AccAddress, amount int64, ok bool) {
	coins := sdk.NewCoins(sdk.NewInt64Coin(sh.denom, amount))
	msg := types.NewMsgWithdrawCollateral(addr, coins)
	sh.Handle(msg, ok)
}

func (sh *Helper) CreatePool(addr, sponsorAddr sdk.AccAddress, nativeDeposit, shield, shieldLimit int64, sponsor, description string) {
	shieldCoins := sdk.NewCoins(sdk.NewInt64Coin(sh.denom, shield))
	depositCoins := types.MixedCoins{Native: sdk.NewCoins(sdk.NewInt64Coin(sh.denom, nativeDeposit))}
	limit := sdk.NewInt(shieldLimit)
	msg := types.NewMsgCreatePool(addr, shieldCoins, depositCoins, sponsor, sponsorAddr, description, limit)
	sh.Handle(msg, true)
}

func (sh *Helper) PurchaseShield(purchaser sdk.AccAddress, shield int64, poolID uint64, ok bool) {
	shieldCoins := sdk.NewCoins(sdk.NewInt64Coin(sh.denom, shield))
	msg := types.NewMsgPurchaseShield(poolID, shieldCoins, "test_purchase", purchaser)
	sh.Handle(msg, ok)
}

func (sh *Helper) ShieldClaimProposal(proposer sdk.AccAddress, loss int64, poolID, purchaseID uint64, ok bool) {
	lossCoins := sdk.NewCoins(sdk.NewInt64Coin(sh.denom, loss))
	proposal := types.NewShieldClaimProposal(poolID, lossCoins, purchaseID, "test_claim_evidence", "test_claim_description", proposer)
	sh.HandleProposal(proposal, ok)
}

func (sh *Helper) WithdrawReimbursement(purchaser sdk.AccAddress, proposalID uint64, ok bool) {
	msg := types.NewMsgWithdrawReimbursement(proposalID, purchaser)
	sh.Handle(msg, ok)
}

// TurnBlock updates context and calls endblocker.
func (sh *Helper) TurnBlock(ctx sdk.Context) {
	sh.ctx = ctx
	shield.EndBlocker(sh.ctx, sh.k)
}

// Handle calls shield handler on a given message
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

// HandleProposal calls shield proposal handler on a given proposal.
func (sh *Helper) HandleProposal(content govtypes.Content, ok bool) {
	err := sh.ph(sh.ctx, content)
	if ok {
		require.NoError(sh.t, err)
	} else {
		require.Error(sh.t, err)
	}
}
