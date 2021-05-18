package keeper

import (
	"context"


	"github.com/irisnet/irismod/modules/nft/keeper"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

type msgServer struct {
	nfttypes.MsgServer
	Keeper
}

func (m msgServer) CreateAdmin(ctx context.Context, admin *types.MsgCreateAdmin) (*types.MsgIssueAdminResponse, error) {
	return &types.MsgIssueAdminResponse{}, nil
}

func (m msgServer) RevokeAdmin(ctx context.Context, admin *types.MsgRevokeAdmin) (*types.MsgRevokeAdminResponse, error) {
	return &types.MsgRevokeAdminResponse{}, nil
}

var _ types.MsgServer = msgServer{}

func (k Keeper) NewMsgServerImpl() types.MsgServer {
	basicServer := keeper.NewMsgServerImpl(k.Keeper)
	return msgServer {
		MsgServer: basicServer,
		Keeper: k,
	}
}