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
		paramSpace: paramSpace,
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
