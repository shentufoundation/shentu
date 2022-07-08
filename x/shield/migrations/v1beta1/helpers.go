package v1beta1

import (
	"fmt"

	"github.com/gogo/protobuf/grpc"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// We use the baseapp.QueryRouter here to do inter-module state querying.
// PLEASE DO NOT REPLICATE THIS PATTERN IN YOUR OWN APP.
func getBondDenom(ctx sdk.Context, queryServer grpc.Server) (string, error) {
	querier, ok := queryServer.(*baseapp.GRPCQueryRouter)
	if !ok {
		return "", fmt.Errorf("unexpected type: %T wanted *baseapp.GRPCQueryRouter", queryServer)
	}

	queryFn := querier.Route(stakingParamsPath)

	q := &stakingtypes.QueryParamsRequest{}

	b, err := proto.Marshal(q)
	if err != nil {
		return "", fmt.Errorf("cannot marshal staking params query request, %w", err)
	}
	req := abci.RequestQuery{
		Data: b,
		Path: stakingParamsPath,
	}

	resp, err := queryFn(ctx, req)
	if err != nil {
		return "", fmt.Errorf("staking query error, %w", err)
	}

	params := new(stakingtypes.QueryParamsResponse)
	if err := proto.Unmarshal(resp.Value, params); err != nil {
		return "", fmt.Errorf("unable to unmarshal delegator query delegations: %w", err)
	}

	return params.Params.BondDenom, nil
}
