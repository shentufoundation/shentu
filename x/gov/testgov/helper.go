package testgov

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/shentufoundation/shentu/v2/x/gov"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t       *testing.T
	msgSrvr govtypesv1beta1.MsgServer
	k       keeper.Keeper

	ctx   sdk.Context
	denom string
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper, denom string) *Helper {
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(govtypes.ModuleName)))
	msgSrvr := keeper.NewMsgServerImpl(k)
	return &Helper{t, keeper.NewLegacyMsgServerImpl(moduleAcct.String(), msgSrvr, k), k, ctx, denom}
}

func (gh *Helper) ShieldClaimProposal(proposer sdk.AccAddress, loss int64, poolID, purchaseID uint64, ok bool) {
	initDeposit := sdk.NewCoins(sdk.NewInt64Coin(gh.denom, 5000e6))
	lossCoins := sdk.NewCoins(sdk.NewInt64Coin(gh.denom, loss))
	content := shieldtypes.NewShieldClaimProposal(poolID, lossCoins, purchaseID, "test_claim_evidence", "test_claim_description", proposer)
	msg, err := govtypesv1beta1.NewMsgSubmitProposal(content, initDeposit, proposer)
	require.NoError(gh.t, err)
	res, err := gh.msgSrvr.SubmitProposal(sdk.WrapSDKContext(gh.ctx), msg)
	require.NoError(gh.t, err)
	require.NotNil(gh.t, res)
}

// TurnBlock updates context and calls endblocker.
func (gh *Helper) TurnBlock(ctx sdk.Context) {
	gh.ctx = ctx
	gov.EndBlocker(gh.ctx, gh.k)
}
