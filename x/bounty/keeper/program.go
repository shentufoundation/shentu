package keeper

import (
	"encoding/binary"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

func (k Keeper) GetNextProgramID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextProgramIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial program ID hasn't been set")
	}
	return binary.LittleEndian.Uint64(bz), nil
}

func (k Keeper) SetNextProgramID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.NextProgramIDKey, bz)
}

// GetPrograms returns all the programs from store
func (k Keeper) GetPrograms(ctx sdk.Context) (programs types.Programs) {
	k.IteratePrograms(ctx, func(program types.Program) bool {
		programs = append(programs, program)
		return false
	})
	return
}

// IteratePrograms iterates over the all the programs and performs a callback function
func (k Keeper) IteratePrograms(ctx sdk.Context, cb func(program types.Program) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ProgramsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var program types.Program
		k.cdc.MustUnmarshal(iterator.Value(), &program)

		if cb(program) {
			break
		}
	}
}
