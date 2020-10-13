package keeper

import (
	"encoding/hex"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// NewQuerier creates a querier for shield module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryPool:
			return queryPool(ctx, path[1:], k)
		case types.QueryPools:
			return queryPools(ctx, path[1:], k)
		case types.QueryPurchase:
			return queryPurchase(ctx, path[1:], k)
		case types.QueryOnesPurchases:
			return queryOnesPurchases(ctx, path[1:], k)
		case types.QueryPoolPurchases:
			return queryPoolPurchases(ctx, path[1:], k)
		case types.QueryOnesCollaterals:
			return queryOnesCollaterals(ctx, path[1:], k)
		case types.QueryPoolCollaterals:
			return queryPoolCollaterals(ctx, path[1:], k)
		case types.QueryProvider:
			return queryProvider(ctx, path[1:], k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func validatePathLength(path []string, length int) error {
	if len(path) != length {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting %d args. Found %d.", length, len(path))
	}
	return nil
}

// queryPool returns information about the queried pool.
func queryPool(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 2); err != nil {
		return nil, err
	}
	var pool types.Pool
	var found bool
	if path[0] == "id" {
		id, err := strconv.ParseUint(path[1], 10, 64)
		if err != nil {
			return nil, err
		}
		pool, err = k.GetPool(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	if path[0] == "sponsor" {
		pool, found = k.GetPoolBySponsor(ctx, path[1])
		if !found {
			return nil, types.ErrNoPoolFoundForSponsor
		}
	}
	res, err = codec.MarshalJSONIndent(k.cdc, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPools returns information about all the pools.
func queryPools(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	res, err = codec.MarshalJSONIndent(k.cdc, k.GetAllPools(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPurchase returns information about a queried purchase.
func queryPurchase(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	txhash, err := hex.DecodeString(path[0])
	if err != nil {
		return nil, err
	}
	purchase, err := k.GetPurchase(ctx, txhash)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, purchase)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryOnesPurchases returns information about a community member's purchases.
func queryOnesPurchases(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, k.GetOnesPurchases(ctx, address))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPoolPurchases queries all purchases in a pool.
func queryPoolPurchases(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, k.GetPoolPurchases(ctx, id))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryOnesCollaterals returns information about one's collaterals.
func queryOnesCollaterals(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, k.GetOnesCollaterals(ctx, address))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPoolCollaterals queries a given pool's collaterals.
func queryPoolCollaterals(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, k.GetAllPoolCollaterals(ctx, pool))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryProvider returns information about a provider.
func queryProvider(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}
	provider, found := k.GetProvider(ctx, addr)
	if !found {
		return nil, types.ErrProviderNotFound
	}

	res, err = codec.MarshalJSONIndent(k.cdc, provider)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
