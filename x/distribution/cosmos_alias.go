package distribution

import (
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
)

const (
	ModuleName        = distr.ModuleName
	StoreKey          = distr.StoreKey
	RouterKey         = distr.RouterKey
	DefaultParamspace = distr.DefaultParamspace
)

var (
	// function aliases
	NewKeeper                            = distr.NewKeeper
	RegisterCodec                        = distr.RegisterCodec
	NewCommunityPoolSpendProposalHandler = distr.NewCommunityPoolSpendProposalHandler

	// variable aliases
	ProposalHandler = distr.ProposalHandler
)

type (
	Keeper = distr.Keeper
)
