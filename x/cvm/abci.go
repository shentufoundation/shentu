package cvm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cvm/keeper"
)

// BeginBlocker stores previous block's hash into the k-v store.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.StoreLastBlockHash(ctx)
}

// EndBlocker ends the block by sending all coins stored at the zero address to the community pool.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	if err := k.RecycleCoins(ctx); err != nil {
		panic(err)
	}
}
