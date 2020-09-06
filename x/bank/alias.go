package bank

import (
	"github.com/certikfoundation/shentu/x/bank/internal/keeper"
	"github.com/certikfoundation/shentu/x/bank/internal/types"
)

var (
	// functions aliases
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec

	// variable aliases
	ModuleCdc = types.ModuleCdc
)
