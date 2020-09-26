package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	sk         types.StakingKeeper
	bk         types.BankKeeper
	paramSpace params.Subspace
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk types.StakingKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		sk:         sk,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (staking.ValidatorI, bool) {
	return k.sk.GetValidator(ctx, addr)
}

func (k Keeper) CreatePool(
	ctx sdk.Context, creator sdk.AccAddress, coverage sdk.Coins, deposit types.MixedCoins, sponsor string) (types.Pool, error) {
	operator := k.GetOperator(ctx)
	if !creator.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool := types.NewPool(coverage, deposit, sponsor)
	k.SetPool(ctx, pool)
	return pool, nil
}

func (k Keeper) UpdatePool(
	ctx sdk.Context, updater sdk.AccAddress, coverage sdk.Coins, deposit types.MixedCoins, sponsor string) (types.Pool, error) {
	operator := k.GetOperator(ctx)
	if !updater.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool := k.GetPool(ctx, sponsor)
	pool.Coverage = pool.Coverage.Add(coverage...)
	pool.Premium = pool.Premium.Add(deposit)
	k.SetPool(ctx, pool)
	return pool, nil
}

func (k Keeper) PausePool(	ctx sdk.Context, updater sdk.AccAddress, sponsor string) (types.Pool, error) {
	operator := k.GetOperator(ctx)
	if !updater.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool := k.GetPool(ctx, sponsor)
	if pool.Active == false {
		return types.Pool{}, types.ErrPoolAlreadyPaused
	}
	pool.Active = false
	k.SetPool(ctx, pool)
	return pool, nil
}

func (k Keeper) ResumePool(	ctx sdk.Context, updater sdk.AccAddress, sponsor string) (types.Pool, error) {
	operator := k.GetOperator(ctx)
	if !updater.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool := k.GetPool(ctx, sponsor)
	if pool.Active == true {
		return types.Pool{}, types.ErrPoolAlreadyActive
	}
	pool.Active = true
	k.SetPool(ctx, pool)
	return pool, nil
}

// set the main record holding validator details
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetPoolKey([]byte(pool.Sponsor)), bz)
}

// set the main record holding validator details
func (k Keeper) GetPool(ctx sdk.Context, sponsor string) (pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey([]byte(sponsor)))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	return
}

// SetPoolParams sets parameters subspace for shield pool parameters.
func (k Keeper) SetPoolParams(ctx sdk.Context, poolParams types.PoolParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyPoolParams, &poolParams)
}

// GetPoolParams returns shield pool parameters.
func (k *Keeper) GetPoolParams(ctx sdk.Context) types.PoolParams {
	var poolParams types.PoolParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyPoolParams, &poolParams)
	return poolParams
}

// SetClaimProposalParams sets parameters subspace for shield claim proposal parameters.
func (k Keeper) SetClaimProposalParams(ctx sdk.Context, claimProposalParams types.ClaimProposalParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyClaimProposalParams, &claimProposalParams)
}

// GetClaimProposalParams returns shield claim proposal parameters.
func (k *Keeper) GetClaimProposalParams(ctx sdk.Context) types.ClaimProposalParams {
	var claimProposalParams types.ClaimProposalParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyClaimProposalParams, &claimProposalParams)
	return claimProposalParams
}
