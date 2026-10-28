package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

// InitGenesis initializes the blacklist module state from genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	for _, entry := range data.Entries {
		if err := k.SetEntry(ctx, entry); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis exports the blacklist module state to genesis.
func (k Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	entries, err := k.GetAllEntries(ctx)
	if err != nil {
		panic(err)
	}
	return &types.GenesisState{Entries: entries}
}
