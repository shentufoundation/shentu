package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Admin(ctx context.Context, request *types.QueryAdminRequest) (*types.QueryAdminResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}
	admin, err := k.GetAdmin(sdkContext, addr)
	return &types.QueryAdminResponse{Admin: admin}, err
}

func (k Keeper) Admins(ctx context.Context, _ *types.QueryAdminsRequest) (*types.QueryAdminsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	admins := k.GetAdmins(sdkContext)
	return &types.QueryAdminsResponse{Admins: admins}, nil
}
