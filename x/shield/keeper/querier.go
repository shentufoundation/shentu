package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// NewQuerier creates a querier for shield module.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryProvider:
			return queryProvider(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryProviders:
			return queryProviders(ctx, req, k, legacyQuerierCdc)
		case types.QueryStatus:
			return queryGlobalState(ctx, path[1:], k, legacyQuerierCdc)
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
	var params types.QueryPaginationParams
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

func queryGlobalState(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}

	shieldState := types.NewQueryResStatus(
		k.GetRemainingServiceFees(ctx),
	)

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, shieldState)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
