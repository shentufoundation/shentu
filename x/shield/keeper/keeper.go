package keeper

import (
	"encoding/binary"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// Keeper implements the shield keeper.
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

// GetValidator returns info of a validator given its operator address.
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
	ret := types.Pool{}
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		if pool.Sponsor == sponsor {
			ret = pool
			return true
		} else {
			return false
		}
	})
	if ret.ID == 0 {
		return ret, false
	}
	return ret, true
}

// DepositNativeServiceFees deposits service fees in native tokens from the shield admin or purchasers.
func (k Keeper) DepositNativeServiceFees(ctx sdk.Context, serviceFees sdk.Coins, from sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, serviceFees)
}

// BondDenom returns staking bond denomination.
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return k.sk.BondDenom(ctx)
}

// GetVotingParams returns gov keeper's voting params.
func (k Keeper) GetVotingParams(ctx sdk.Context) govTypes.VotingParams {
	return k.gk.GetVotingParams(ctx)
}

// SetLastUpdateTime sets the last update time.
// Last update time will be set when the first purchase is made or distributing service fees.
func (k Keeper) SetLastUpdateTime(ctx sdk.Context, prevUpdateTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(prevUpdateTime)
	store.Set(types.GetLastUpdateTimeKey(), bz)
}

// GetLastUpdateTime returns the last update time.
func (k Keeper) GetLastUpdateTime(ctx sdk.Context) (time.Time, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastUpdateTimeKey())
	if bz == nil {
		return time.Time{}, false
	}
	var lastUpdateTime time.Time
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastUpdateTime)
	return lastUpdateTime, true
}
