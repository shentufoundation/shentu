package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

const (
	QueryOperator  = "operator"
	QueryOperators = "operators"
	QueryWithdraws = "withdraws"
	QueryTask      = "task"
	QueryResponse  = "response"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryOperator:
			return queryOperator(ctx, path[1:], keeper, legacyQuerierCdc)
		case QueryOperators:
			return queryOperators(ctx, path[1:], keeper, legacyQuerierCdc)
		case QueryWithdraws:
			return queryWithdraws(ctx, path[1:], keeper, legacyQuerierCdc)
		case QueryTask:
			return queryTask(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case QueryResponse:
			return queryResponse(ctx, path[1:], req, keeper, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

// ValidatePathLength validates the length of a given path.
func validatePathLength(path []string, length int) error {
	if len(path) != length {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting %d args. Found %d.", length, len(path))
	}
	return nil
}

// queryOperator returns information of an operator.
func queryOperator(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return nil, err
	}
	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, operator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryOperators returns information of all operators.
func queryOperators(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, k.GetAllOperators(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryWithdraws returns information of all withdrawals.
func queryWithdraws(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, k.GetAllWithdraws(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryResponse returns information of a response.
func queryResponse(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	var params types.QueryResponseParams
	err = legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	task, err := k.GetTask(ctx, types.NewTaskID(params.Contract, params.Function))
	if err != nil {
		return nil, err
	}
	for _, response := range task.GetResponses() {
		operatorAddr, err := sdk.AccAddressFromBech32(response.Operator)
		if err != nil {
			panic(err)
		}
		if operatorAddr.Equals(params.Operator) {
			res, err = codec.MarshalJSONIndent(legacyQuerierCdc, response)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if res == nil {
		return nil, fmt.Errorf("there is no response from this operator")
	}
	return res, err
}

// queryTask returns information of a task.
func queryTask(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	var params types.QueryTaskParams
	err = legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	task, err := k.GetTask(ctx, types.NewTaskID(params.Contract, params.Function))
	if err != nil {
		return nil, err
	}
	res, err = nil, fmt.Errorf("failed to cast to concrete task")
	if smartContractTask, ok := task.(*types.Task); ok {
		res, err = codec.MarshalJSONIndent(legacyQuerierCdc, *smartContractTask)
	}
	return res, err
}
