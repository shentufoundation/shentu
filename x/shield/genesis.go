package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// InitGenesis initialize store values with genesis states.
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	shieldOperator := data.ShieldOperator
	poolParams := data.PoolParams
	claimProposalParams := data.ClaimProposalParams

	k.SetOperator(ctx, shieldOperator)
	k.SetPoolParams(ctx, poolParams)
	k.SetClaimProposalParams(ctx, claimProposalParams)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	shieldOperator := k.GetOperator(ctx)
	poolParams := k.GetPoolParams(ctx)
	claimProposalParams := k.GetClaimProposalParams(ctx)

	return types.NewGenesisState(shieldOperator, poolParams, claimProposalParams)
}
