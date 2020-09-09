// Package keeper specifies the keeper for the cert module.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// Keeper manages certifier & security council related logics.
type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	slashingKeeper types.SlashingKeeper
	stakingKeeper  types.StakingKeeper
}

// NewKeeper creates a new instance of the certifier keeper.
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, slashingKeeper types.SlashingKeeper, stakingKeeper types.StakingKeeper) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		slashingKeeper: slashingKeeper,
		stakingKeeper:  stakingKeeper,
	}
}

// CertifyPlatform certifies a validator host platform by a certifier.
func (k Keeper) CertifyPlatform(ctx sdk.Context, certifier sdk.AccAddress, validator []byte, description string) error {
	if !k.IsCertifier(ctx, certifier) {
		return types.ErrRejectedValidator
	}
	ctx.KVStore(k.storeKey).Set(types.PlatformStoreKey(validator), []byte(description))
	return nil
}

// GetPlatform returns the host platform of the validator.
func (k Keeper) GetPlatform(ctx sdk.Context, validator []byte) (string, bool) {
	if platform := ctx.KVStore(k.storeKey).Get(types.PlatformStoreKey(validator)); platform != nil {
		return string(platform), true
	}
	return "", false
}

// GetAllPlatforms gets all platform certificates for genesis export
func (k Keeper) GetAllPlatforms(ctx sdk.Context) (platforms []types.Platform) {
	platforms = make([]types.Platform, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PlatformsStoreKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Key()[len(types.PlatformsStoreKey()):]
		description := iterator.Value()
		platforms = append(platforms, types.Platform{Address: address, Description: string(description)})
	}
	return platforms
}
