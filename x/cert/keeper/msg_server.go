package keeper

import (
	"context"

	"github.com/certikfoundation/shentu/x/cert/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the cert MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ProposeCertifier(goCtx context.Context, msg *types.MsgProposeCertifier) (*types.MsgProposeCertifierResponse, error) {
	return &types.MsgProposeCertifierResponse{}, nil
}
