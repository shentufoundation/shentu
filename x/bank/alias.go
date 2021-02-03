package bank

import (
	"github.com/certikfoundation/shentu/x/bank/keeper"
	"github.com/certikfoundation/shentu/x/bank/types"
)

var (
	// functions aliases
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper = keeper.Keeper
)
