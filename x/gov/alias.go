package gov

import (
	"github.com/certikfoundation/shentu/x/gov/internal/keeper"
	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

const (
	AttributeTxHash = types.AttributeTxHash
)

var (
	// function aliases
	NewKeeper           = keeper.NewKeeper
	ProposalHandler     = types.ProposalHandler
	DefaultGenesisState = types.DefaultGenesisState
	ParamKeyTable       = types.ParamKeyTable
)

type (
	Keeper = keeper.Keeper
)
