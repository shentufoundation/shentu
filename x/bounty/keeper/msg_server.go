package keeper

import (
	"context"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CreateProgram(ctx context.Context, program *types.MsgCreateProgram) (*types.MsgCreateProgramResponse, error) {
	//TODO implement create program
	panic("implement me")
}
