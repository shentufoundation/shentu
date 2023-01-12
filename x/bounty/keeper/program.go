package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) GetProgram(ctx sdk.Context, id uint64) (types.Program, bool) {
	store := ctx.KVStore(k.storeKey)

	pBz := store.Get(types.GetProgramKey(id))
	if pBz == nil {
		return types.Program{}, false
	}

	var program types.Program
	k.cdc.MustUnmarshal(pBz, &program)
	return program, true
}

func (k Keeper) SetProgram(ctx sdk.Context, program types.Program) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&program)
	store.Set(types.GetProgramKey(program.ProgramId), bz)
}

func (k Keeper) GetNextProgramID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextProgramIDKey)
	return binary.LittleEndian.Uint64(bz)
}

func (k Keeper) SetNextProgramID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.NextProgramIDKey, bz)
}

// GetPrograms returns all the program from store
func (k Keeper) GetPrograms(ctx sdk.Context) (proposals types.Programs) {
	//TODO implement this
	return
}
