package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// NewQuerier creates a new bounty Querier instance
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryHosts:
			return queryHosts(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case types.QueryHost:
			return queryHost(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case types.QueryPrograms:
			return queryPrograms(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case types.QueryProgram:
			return queryProgram(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case types.QueryFindings:
			return queryFindings(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case types.QueryFinding:
			return queryFinding(ctx, path[1:], req, keeper, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryHosts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	// TODO: implement this
	return nil, nil
}

func queryHost(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	// TODO: implement this
	return nil, nil
}

func queryPrograms(ctx sdk.Context, _ []string, req abci.RequestQuery, keeper Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	// TODO: implement this
	var params types.QueryProgramsParams
	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	programs := keeper.GetProgramsFiltered(ctx, params)
	if programs == nil {
		programs = types.Programs{}
	}

	bz, err := codec.MarshalJSONIndent(cdc, programs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryProgram(ctx sdk.Context, _ []string, req abci.RequestQuery, keeper Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryProgramParams
	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	program, ok := keeper.GetProgram(ctx, params.ProgramID)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrUnknownProgram, "%d", params.ProgramID)
	}

	bz, err := codec.MarshalJSONIndent(cdc, program)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryFindings(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	// TODO: implement this
	return nil, nil
}

func queryFinding(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	// TODO: implement this
	return nil, nil
}
