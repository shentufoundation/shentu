package simulation

import (
	"github.com/cosmos/cosmos-sdk/client"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
)

// WeightedOperations defers to the upstream gov simulation for the
// deposit/vote/cancel ops. Shentu's keeper embeds *govkeeper.Keeper
// directly, so the upstream sim only needs a pointer to the embedded
// field — it has no awareness of (and no reason to touch) the shentu
// certifier-round logic. Cert-update proposals never show up here
// anyway because upstream's proposal msgs are plain text/bank ops.
func WeightedOperations(
	appParams simtypes.AppParams,
	txGen client.TxConfig,
	ak govtypes.AccountKeeper,
	bk govtypes.BankKeeper,
	k keeper.Keeper,
	wMsgs []simtypes.WeightedProposalMsg,
	wContents []simtypes.WeightedProposalContent, //nolint:staticcheck // used for legacy testing
) []simtypes.WeightedOperation {
	return govsim.WeightedOperations(appParams, txGen, ak, bk, &k.Keeper, wMsgs, wContents)
}
