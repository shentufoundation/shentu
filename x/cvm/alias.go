package cvm

import (
	"github.com/certikfoundation/shentu/x/cvm/keeper"
	"github.com/certikfoundation/shentu/x/cvm/types"
)

const (
	ModuleName   = types.ModuleName
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
)

var (
	DefaultGenesisState = types.DefaultGenesisState
	ModuleCdc           = types.ModuleCdc
	DefaultParamSpace   = types.DefaultParamspace
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	ValidateGenesis     = types.ValidateGenesis
	NewMsgCall          = types.NewMsgCall
)

type (
	Keeper       = keeper.Keeper
	MsgCall      = types.MsgCall
	MsgDeploy    = types.MsgDeploy
	GenesisState = types.GenesisState
	State        = keeper.State
	QueryResAbi  = types.QueryResAbi
	QueryResView = types.QueryResView
)
