package keeper

import (
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// NewQuerier creates a querier for shield module.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryPoolByID:
			return queryPoolByID(ctx, path[1:], k)
		case types.QueryPoolBySponsor:
			return queryPoolBySponsor(ctx, path[1:], k)
		case types.QueryPools:
			return queryPools(ctx, path[1:], k)
		case types.QueryPurchaseList:
			return queryPurchaseList(ctx, path[1:], k)
		case types.QueryPurchaserPurchases:
			return queryPurchaserPurchases(ctx, path[1:], k)
		case types.QueryPoolPurchases:
			return queryPoolPurchases(ctx, path[1:], k)
		case types.QueryPurchases:
			return queryPurchases(ctx, path[1:], k)
		case types.QueryProvider:
			return queryProvider(ctx, path[1:], k)
		case types.QueryProviders:
			return queryProviders(ctx, req, k)
		case types.QueryPoolParams:
			return queryPoolParams(ctx, path[1:], k)
		case types.QueryClaimParams:
			return queryClaimParams(ctx, path[1:], k)
		case types.QueryStatus:
			return queryGlobalState(ctx, path[1:], k)
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

func queryPoolByID(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}
	pool, found := k.GetPool(ctx, id)
	if !found {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryPoolBySponsor(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	pool, found := k.GetPoolBySponsor(ctx, path[0])
	if !found {
		return nil, types.ErrNoPoolFoundForSponsor
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

func queryPurchaseList(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 2); err != nil {
		return nil, err
	}

	poolID, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}
	purchaser, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, err
	}
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return []byte{}, types.ErrPurchaseNotFound
	}

	res, err = codec.MarshalJSONIndent(k.cdc, purchaseList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPurchaserPurchases returns information about a community member's purchases.
func queryPurchaserPurchases(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, k.GetPurchaserPurchases(ctx, address))
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

	res, err = codec.MarshalJSONIndent(k.cdc, k.GetPoolPurchaseLists(ctx, id))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPurchases queries all purchases.
func queryPurchases(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, k.GetAllPurchases(ctx))
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

func queryProviders(ctx sdk.Context, req abci.RequestQuery, k Keeper) (res []byte, err error) {
	var params types.QueryPaginationParams
	err = k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	providers := k.GetProvidersPaginated(ctx, uint(params.Page), uint(params.Limit))

	res, err = codec.MarshalJSONIndent(k.cdc, providers)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryPoolParams(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	params := k.GetPoolParams(ctx)

	res, err = codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryClaimParams(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	params := k.GetClaimProposalParams(ctx)

	res, err = codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryGlobalState(ctx sdk.Context, path []string, k Keeper) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	shieldState := types.NewQueryResStatus(
		k.GetTotalCollateral(ctx),
		k.GetTotalShield(ctx),
		k.GetTotalWithdrawing(ctx),
		k.GetServiceFees(ctx),
		k.GetRemainingServiceFees(ctx),
	)

	res, err = codec.MarshalJSONIndent(k.cdc, shieldState)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
