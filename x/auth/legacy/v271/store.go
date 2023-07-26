package v271

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/shentufoundation/shentu/v2/common"
	authtypes "github.com/shentufoundation/shentu/v2/x/auth/types"
)

func MigrateAccount(ctx sdk.Context, ak authtypes.AccountKeeper) error {
	authKey := sdk.NewKVStoreKey(sdktypes.StoreKey)
	store := ctx.KVStore(authKey)

	iterator := sdk.KVStorePrefixIterator(store, sdktypes.AddressStoreKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		acc, err := ak.UnmarshalAccount(iterator.Value())
		if err != nil {
			return err
		}
		oldAddr := acc.GetAddress()
		newAddr, err := common.PrefixToShentu(oldAddr.String())
		if err != nil {
			return err
		}

		switch account := acc.(type) {
		case *sdktypes.BaseAccount:
			account.Address = newAddr
		case *sdktypes.ModuleAccount:
			account.Address = newAddr
		case *authtypes.ManualVestingAccount:
			newUnlocker, err := common.PrefixToShentu(account.Unlocker)
			if err != nil {
				return err
			}
			account.Address = newAddr
			account.Unlocker = newUnlocker
		default:
			return errors.New("unknown account type")
		}

		bz, err := ak.MarshalAccount(acc)
		if err != nil {
			return err
		}

		accAddress, err := sdk.AccAddressFromBech32(newAddr)
		if err != nil {
			return err
		}
		store.Set(sdktypes.AddressStoreKey(accAddress), bz)
		store.Delete(sdktypes.AddressStoreKey(oldAddr))
	}

	return nil
}
