package mint

import (
	"github.com/certikfoundation/shentu/x/mint/internal/keeper"
)

var (
	// function aliases
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
