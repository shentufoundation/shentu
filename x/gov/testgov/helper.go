package testgov

import (
	"testing"

	"github.com/shentufoundation/shentu/v2/x/gov/keeper"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/gov"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t *testing.T
	h sdk.Handler
	k keeper.Keeper

	ctx   sdk.Context
	denom string
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper, denom string) *Helper {
	return &Helper{t, gov.NewHandler(k), k, ctx, denom}
}

func (gh *Helper) ShieldClaimProposal(proposer sdk.AccAddress, loss int64, poolID, purchaseID uint64, ok bool) *sdk.Result {
	initDeposit := sdk.NewCoins(sdk.NewInt64Coin(gh.denom, 5000e6))
	lossCoins := sdk.NewCoins(sdk.NewInt64Coin(gh.denom, loss))
	content := shieldtypes.NewShieldClaimProposal(poolID, lossCoins, purchaseID, "test_claim_evidence", "test_claim_description", proposer)
	proposal, err := govtypes.NewMsgSubmitProposal(content, initDeposit, proposer)
	require.NoError(gh.t, err)
	return gh.Handle(proposal, ok)
}

// TurnBlock updates context and calls endblocker.
func (gh *Helper) TurnBlock(ctx sdk.Context) {
	gh.ctx = ctx
	gov.EndBlocker(gh.ctx, gh.k)
}

// Handle calls shield handler on a given message
func (gh *Helper) Handle(msg sdk.Msg, ok bool) *sdk.Result {
	res, err := gh.h(gh.ctx, msg)
	if ok {
		require.NoError(gh.t, err)
		require.NotNil(gh.t, res)
	} else {
		require.Error(gh.t, err)
		require.Nil(gh.t, res)
	}
	return res
}
