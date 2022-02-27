package keeper

import (
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// NewQuerier creates a querier for shield module.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case v1beta1.QueryPoolByID:
			return queryPoolByID(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryPoolBySponsor:
			return queryPoolBySponsor(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryPools:
			return queryPools(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryPurchaserPurchases:
			return queryPurchaserPurchases(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryPoolPurchases:
			return queryPoolPurchases(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryPurchases:
			return queryPurchases(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryProvider:
			return queryProvider(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryProviders:
			return queryProviders(ctx, req, k, legacyQuerierCdc)
		case v1beta1.QueryPoolParams:
			return queryPoolParams(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryClaimParams:
			return queryClaimParams(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryStatus:
			return queryGlobalState(ctx, path[1:], k, legacyQuerierCdc)
		case v1beta1.QueryStakedForShield:
			return queryPurchase(ctx, path[1:], k, legacyQuerierCdc)
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

func queryPoolByID(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
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

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryPoolBySponsor(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	pool, found := k.GetPoolsBySponsor(ctx, path[0])
	if !found {
		return nil, types.ErrNoPoolFoundForSponsor
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPools returns information about all the pools.
func queryPools(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, k.GetAllPools(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPurchaserPurchases returns information about a community member's purchases.
func queryPurchaserPurchases(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, k.GetPurchaserPurchases(ctx, address))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPoolPurchases queries all purchases in a pool.
func queryPoolPurchases(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, k.GetPoolPurchases(ctx, id))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPurchases queries all purchases.
func queryPurchases(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, k.GetAllPurchase(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryProvider returns information about a provider.
func queryProvider(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
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

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, provider)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryProviders(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	var params v1beta1.QueryPaginationParams
	err = legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	providers := k.GetProvidersPaginated(ctx, uint(params.Page), uint(params.Limit))

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, providers)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryPoolParams(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	params := k.GetPoolParams(ctx)

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryClaimParams(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	params := k.GetClaimProposalParams(ctx)

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryGlobalState(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	shieldState := v1beta1.NewQueryResStatus(
		k.GetTotalCollateral(ctx),
		k.GetTotalShield(ctx),
		k.GetTotalWithdrawing(ctx),
		k.GetServiceFees(ctx),
		k.GetGlobalStakingPool(ctx),
	)

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, shieldState)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryPurchase queries staked-for-shield for pool-purchaser pair.
func queryPurchase(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
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
	purchaseList, found := k.GetPurchase(ctx, poolID, purchaser)
	if !found {
		return []byte{}, types.ErrPurchaseNotFound
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, purchaseList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
