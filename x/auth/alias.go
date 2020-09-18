package auth

import (
	"github.com/certikfoundation/shentu/x/auth/client/cli"
	"github.com/certikfoundation/shentu/x/auth/internal/keeper"
	"github.com/certikfoundation/shentu/x/auth/internal/types"
)

var (
	// function aliases
	RegisterCodec            = types.RegisterCodec
	RegisterAccountTypeCodec = types.RegisterAccountTypeCodec
	GetCmdUnlock             = cli.GetCmdUnlock
	NewKeeper                = keeper.NewKeeper

	// variable aliases
	ModuleCdc = types.ModuleCdc
	RouterKey = types.RouterKey
)

type (
	Keeper = keeper.Keeper
)
