package keeper

import (
	"encoding/binary"

	govTypes "github.com/cosmos/cosmos-sdk/x/gov"

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
	gk           types.GovKeeper
	supplyKeeper types.SupplyKeeper
	paramSpace   params.Subspace
}

// NewKeeper creates a shield keeper.
func NewKeeper(cdc *codec.Codec, shieldStoreKey sdk.StoreKey, sk types.StakingKeeper, gk types.GovKeeper, supplyKeeper types.SupplyKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:     shieldStoreKey,
		cdc:          cdc,
		sk:           sk,
		gk:           gk,
		supplyKeeper: supplyKeeper,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (staking.ValidatorI, bool) {
	return k.sk.GetValidator(ctx, addr)
}

// SetLatestPoolID sets the latest pool ID to store.
func (k Keeper) SetNextPoolID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.GetNextPoolIDKey(), bz)
}

// GetNextPoolID gets the latest pool ID from store.
func (k Keeper) GetNextPoolID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.GetNextPoolIDKey())
	return binary.LittleEndian.Uint64(opBz)
}

// GetPoolBySponsor search store for a pool object with given pool ID.
func (k Keeper) GetPoolBySponsor(ctx sdk.Context, sponsor string) (types.Pool, bool) {
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
		return ret, false
	}
	return ret, true
}

// DepositNativePremium deposits premium in native tokens from the shield admin or purchasers.
func (k Keeper) DepositNativePremium(ctx sdk.Context, premium sdk.Coins, from sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, premium)
}

// BondDenom returns staking bond denomination
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return k.sk.BondDenom(ctx)
}

// GetVotingParams returns gov keeper's voting params.
func (k Keeper) GetVotingParams(ctx sdk.Context) govTypes.VotingParams {
	return k.gk.GetVotingParams(ctx)
}
