package keeper

import (
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/shield/types"
)

const (
	QueryPool     = "pool"
	QueryPurchase = "purchase"
)

// NewQuerier creates a querier for shield module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryPool:
			return queryPool(ctx, path[1:], k)
		case QueryPurchase:
			return queryPurchase(ctx, path[1:], k)
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
		pool, err = k.GetPoolBySponsor(ctx, path[1])
		if err != nil {
			return nil, err
		}
	}
	res, err = codec.MarshalJSONIndent(k.cdc, pool)
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

	purchase, err := k.GetPurchase(ctx, path[0])
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, purchase)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
