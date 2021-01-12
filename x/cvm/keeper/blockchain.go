package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/execution/errors"

	"github.com/certikfoundation/shentu/x/cvm/types"
)

// Blockchain implements the blockchain interface from burrow to make state queries.
type Blockchain struct {
	ctx sdk.Context
	k   Keeper
}

// NewBlockChain returns the pointer to a new BlockChain type data.
func NewBlockChain(ctx sdk.Context, k Keeper) *Blockchain {
	return &Blockchain{
		ctx: ctx,
		k:   k,
	}
}

// LastBlockHeight returns the last block height of the chain.
func (bc *Blockchain) LastBlockHeight() uint64 {
	return uint64(bc.ctx.BlockHeight())
}

// LastBlockTime return the unix Time type for the last block.
func (bc *Blockchain) LastBlockTime() time.Time {
	return bc.ctx.BlockHeader().Time
}

// BlockHash returns the block's hash at the provided height.
func (bc *Blockchain) BlockHash(height uint64) ([]byte, error) {
	if height > uint64(bc.ctx.BlockHeight()) {
		return nil, errors.Codes.InvalidBlockNumber
	}
	return bc.ctx.KVStore(bc.k.key).Get(types.BlockHashStoreKey(int64(height))), nil
}
