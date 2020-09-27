package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	sk           types.StakingKeeper
	supplyKeeper types.SupplyKeeper
	paramSpace   params.Subspace
}

// NewKeeper creates a slashing keeper
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, sk types.StakingKeeper, supplyKeeper types.SupplyKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		sk:           sk,
		supplyKeeper: supplyKeeper,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (staking.ValidatorI, bool) {
	return k.sk.GetValidator(ctx, addr)
}

func (k Keeper) CreatePool(
	ctx sdk.Context, creator sdk.AccAddress, shield sdk.Coins, deposit types.MixedCoins, sponsor string,
	timeOfCoverage, blocksOfCoverage int64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !creator.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldOperator
	}

	// Store endTime. If not available, store endBlockHeight.
	var endTime, endBlockHeight int64
	startBlockHeight := ctx.BlockHeight()
	if timeOfCoverage != 0 {
		endTime = ctx.BlockTime().Unix() + timeOfCoverage
	} else if blocksOfCoverage != 0 {
		endBlockHeight = startBlockHeight + blocksOfCoverage
	}

	id := k.GetNextPoolID(ctx)
	depositDec := types.MixedDecCoinsFromMixedCoins(deposit)

	pool := types.NewPool(admin, shield, depositDec, sponsor, endTime, startBlockHeight, endBlockHeight, id)

	// TODO: ensure shield is backed by staking from admin?

	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, id+1)
	return pool, nil
}

func (k Keeper) UpdatePool(
	ctx sdk.Context, updater sdk.AccAddress, shield sdk.Coins, deposit types.MixedCoins, id uint64,
	additionalTime, additionalBlocks int64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldOperator
	}

	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return types.Pool{}, err
	}

	// Extend EndTime. If not available, extend EndBlockHeight.
	if additionalTime != 0 {
		if pool.EndTime == 0 {
			return types.Pool{}, types.ErrCannotExtend
		}
		pool.EndTime += additionalTime
	} else if additionalBlocks != 0 {
		if pool.EndBlockHeight == 0 {
			return types.Pool{}, types.ErrCannotExtend
		}
		pool.EndBlockHeight += additionalBlocks
	}

	pool.Shield = pool.Shield.Add(shield...)
	pool.Premium = pool.Premium.Add(types.MixedDecCoinsFromMixedCoins(deposit))
	k.SetPool(ctx, pool)
	return pool, nil
}

func (k Keeper) PausePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return types.Pool{}, err
	}
	if !pool.Active {
		return types.Pool{}, types.ErrPoolAlreadyPaused
	}
	pool.Active = false
	k.SetPool(ctx, pool)
	return pool, nil
}

func (k Keeper) ResumePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return types.Pool{}, err
	}
	if pool.Active {
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
	store.Set(types.GetPoolKey(pool.PoolID), bz)
}

// set the main record holding validator details
func (k Keeper) GetPool(ctx sdk.Context, id uint64) (types.Pool, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(id))
	if bz != nil {
		var pool types.Pool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
		return pool, nil
	}
	return types.Pool{}, types.ErrNoPoolFound
}

// GetAllPools retrieves all pools in the store.
func (k Keeper) GetAllPools(ctx sdk.Context) (pools []types.Pool) {
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		pools = append(pools, pool)
		return false
	})
	return pools
}

// IterateAllPools iterates over the all the stored pools and performs a callback function.
func (k Keeper) IterateAllPools(ctx sdk.Context, callback func(certificate types.Pool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &pool)

		if callback(pool) {
			break
		}
	}
}

// DepositCollateral deposits a community member's collateral for a pool.
func (k Keeper) DepositCollateral(ctx sdk.Context, from sdk.AccAddress, id uint64, amount sdk.Coins) error {
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return err
	}

	// update the pool - update or create collateral entry
	found := false
	for i, collateral := range pool.Community {
		if collateral.Provider.Equals(from) {
			pool.Community[i].Amount = pool.Community[i].Amount.Add(amount...)
			found = true
		}
	}
	if !found {
		pool.Community = append(pool.Community, types.NewCollateral(from, amount))
	}

	pool.TotalCollateral = pool.TotalCollateral.Add(amount...)
	k.SetPool(ctx, pool)

	return nil
}

// SetLatestPoolID sets the latest pool ID to store.
func (k Keeper) SetNextPoolID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.GetNextPoolIDKey(), bz)
}

// GetLatestPoolID gets the latest pool ID from store.
func (k Keeper) GetNextPoolID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.GetNextPoolIDKey())
	return binary.LittleEndian.Uint64(opBz)
}

// GetPoolByID search store for a pool object with given pool ID.
func (k Keeper) GetPoolBySponsor(ctx sdk.Context, sponsor string) (types.Pool, error) {
	ret := types.Pool{
		PoolID: 0,
	}
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		if pool.Sponsor == sponsor {
			ret = pool
			return true
		} else {
			return false
		}
	})
	if ret.PoolID == 0 {
		return ret, types.ErrNoPoolFound
	}
	return ret, nil
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
