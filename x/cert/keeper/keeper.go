// Package keeper specifies the keeper for the cert module.
package keeper

import (
	"context"

	"cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// Keeper manages certifier & security council related logics.
type Keeper struct {
	storeService   store.KVStoreService
	cdc            codec.BinaryCodec
	slashingKeeper types.SlashingKeeper
	stakingKeeper  types.StakingKeeper
}

// NewKeeper creates a new instance of the certifier keeper.
func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, slashingKeeper types.SlashingKeeper, stakingKeeper types.StakingKeeper) Keeper {
	return Keeper{
		cdc:            cdc,
		storeService:   storeService,
		slashingKeeper: slashingKeeper,
		stakingKeeper:  stakingKeeper,
	}
}

// CertifyPlatform certifies a validator host platform by a certifier.
func (k Keeper) CertifyPlatform(ctx context.Context, certifier sdk.AccAddress, validator cryptotypes.PubKey, description string) error {
	if _, err := k.IsCertifier(ctx, certifier); err != nil {
		return types.ErrRejectedValidator
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	pkAny, err := codectypes.NewAnyWithValue(validator)
	if err != nil {
		return err
	}

	bz := k.cdc.MustMarshal(&types.Platform{ValidatorPubkey: pkAny, Description: description})
	return kvStore.Set(types.PlatformStoreKey(validator), bz)
}

// GetPlatform returns the host platform of the validator.
func (k Keeper) GetPlatform(ctx sdk.Context, validator cryptotypes.PubKey) (types.Platform, bool) {
	var platform types.Platform
	var found bool
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PlatformStoreKey(validator))
	if err != nil {
		found = false
	}
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &platform)
		found = true
	}
	return platform, found
}

// GetAllPlatforms gets all platform certificates for genesis export
func (k Keeper) GetAllPlatforms(ctx sdk.Context) (platforms []types.Platform) {
	platforms = make([]types.Platform, 0)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PlatformsStoreKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var platform types.Platform
		k.cdc.MustUnmarshal(iterator.Value(), &platform)
		platforms = append(platforms, platform)
	}
	return platforms
}
