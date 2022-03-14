package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

// Keeper implements the shield keeper.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	ak         types.AccountKeeper
	bk         types.BankKeeper
	sk         types.StakingKeeper
	gk         types.GovKeeper
	paramSpace types.ParamSubspace
}

// NewKeeper creates a shield keeper.
func NewKeeper(cdc codec.BinaryCodec, shieldStoreKey sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper,
	sk types.StakingKeeper, gk types.GovKeeper, paramSpace types.ParamSubspace) Keeper {
	return Keeper{
		storeKey:   shieldStoreKey,
		cdc:        cdc,
		ak:         ak,
		bk:         bk,
		sk:         sk,
		gk:         gk,
		paramSpace: paramSpace,
	}
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

// GetPoolsBySponsor search store for a pool object with given pool ID.
func (k Keeper) GetPoolsBySponsor(ctx sdk.Context, sponsorAddr string) ([]v1beta1.Pool, bool) {
	var ret []v1beta1.Pool
	found := false
	k.IterateAllPools(ctx, func(pool v1beta1.Pool) bool {
		if pool.SponsorAddr == sponsorAddr {
			ret = append(ret, pool)
			found = true
		}
		return false
	})
	return ret, found
}

// BondDenom returns staking bond denomination.
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return k.sk.BondDenom(ctx)
}
