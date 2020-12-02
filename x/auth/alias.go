package auth

import (
	"github.com/certikfoundation/shentu/x/auth/client/cli"
	"github.com/certikfoundation/shentu/x/auth/client/rest"
	"github.com/certikfoundation/shentu/x/auth/keeper"
	"github.com/certikfoundation/shentu/x/auth/types"
)

var (
	// function aliases
	GetCmdUnlock   = cli.GetCmdUnlock
	NewKeeper      = keeper.NewKeeper
	RegisterRoutes = rest.RegisterRoutes

	// variable aliases
	ModuleCdc = types.ModuleCdc
	RouterKey = types.RouterKey
)

type (
	Keeper = keeper.Keeper
)
