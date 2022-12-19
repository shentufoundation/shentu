package keeper

import (
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// NewQuerier creates a querier for bounty module.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryProgramByID:
			return queryProgramByID(ctx, path[1:], k, legacyQuerierCdc)
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

func queryProgramByID(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}
	program, found := k.GetProgram(ctx, id)
	if !found {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, program)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
