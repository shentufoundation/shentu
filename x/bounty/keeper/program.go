package keeper

import (
	"encoding/binary"

	errorsmod "cosmossdk.io/errors"

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

func (k Keeper) GetAllPrograms(ctx sdk.Context) []types.Program {
	store := ctx.KVStore(k.storeKey)

	var programs []types.Program
	var program types.Program

	iterator := sdk.KVStorePrefixIterator(store, types.ProgramsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshal(iterator.Value(), &program)
		programs = append(programs, program)
	}
	return programs
}

func (k Keeper) SetProgram(ctx sdk.Context, program types.Program) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&program)
	store.Set(types.GetProgramKey(program.ProgramId), bz)
}

func (k Keeper) GetNextProgramID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetNextProgramIDKey())
	if bz == nil {
		return 1, errorsmod.Wrap(types.ErrInvalidGenesis, "initial next finding ID hasn't been set")
	}
	return binary.LittleEndian.Uint64(bz), nil
}

func (k Keeper) SetNextProgramID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.GetNextProgramIDKey(), bz)
}

func (k Keeper) EndProgram(ctx sdk.Context, caller sdk.AccAddress, id uint64) error {
	program, found := k.GetProgram(ctx, id)
	if !found {
		return types.ErrProgramNotExists
	}
	host, err := sdk.AccAddressFromBech32(program.CreatorAddress)
	if err != nil {
		return types.ErrProgramCreatorInvalid
	}
	if !caller.Equals(host) && !k.certKeeper.IsCertifier(ctx, caller) {
		return types.ErrProgramNotAllowed
	}
	if !program.Active {
		return types.ErrProgramInactive
	}
	if ctx.BlockTime().After(program.SubmissionEndTime) {
		return types.ErrProgramExpired
	}
	program.Active = false
	k.SetProgram(ctx, program)
	return nil
}
