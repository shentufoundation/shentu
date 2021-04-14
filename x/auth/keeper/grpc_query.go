package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/auth/types"
)

var _ types.QueryServer = Keeper{}

// Validator returns Bech32 certik address corresponding to Bech32 certikvaloper address.
func (q Keeper) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorResponse{AccountAddress: sdk.AccAddress(valAddr).String()}, nil
}
