package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock executes logics to begin a block
func BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
}

// EndBlock executes logics to begin a block
func EndBlock(ctx sdk.Context, req abci.RequestEndBlock, k Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
