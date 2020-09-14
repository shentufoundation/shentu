package oracle

import (
	"github.com/certikfoundation/shentu/x/oracle/internal/keeper"
	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

const (
	ModuleName        = types.ModuleName
	QuerierRoute      = types.QuerierRoute
	StoreKey          = types.StoreKey
	DefaultParamSpace = types.ModuleName
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewMsgTaskResponse  = types.NewMsgTaskResponse
	DefaultGenesisState = types.DefaultGenesisState
)

type (
	Keeper          = keeper.Keeper
	MsgTaskResponse = types.MsgTaskResponse
	MsgCreateTask   = types.MsgCreateTask
)
