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

func (m msgServer) CreateNFTAdmin(ctx context.Context, admin *types.MsgCreateNFTAdmin) (*types.MsgIssueNFTAdminResponse, error) {

	return &types.MsgIssueNFTAdminResponse{}, nil
}

func (m msgServer) RevokeNFTAdmin(ctx context.Context, admin *types.MsgRevokeNFTAdmin) (*types.MsgRevokeAdminResponse, error) {
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