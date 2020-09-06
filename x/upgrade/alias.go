package upgrade

import (
	"github.com/certikfoundation/shentu/x/upgrade/client"
	"github.com/certikfoundation/shentu/x/upgrade/internal/keeper"
)

var (
	// function aliases
	NewKeeper = keeper.NewKeeper

	// variable aliases
	ProposalHandler = client.ProposalHandler
)

type (
	Keeper = keeper.Keeper
)
