package staking

import (
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/staking/internal/keeper"
)

const (
	ModuleName        = staking.ModuleName
	StoreKey          = staking.StoreKey
	BondedPoolName    = staking.BondedPoolName
	NotBondedPoolName = staking.NotBondedPoolName
	TStoreKey         = staking.TStoreKey
	DefaultParamspace = staking.DefaultParamspace
)

var (
	// function aliases
	NewKeeper            = keeper.NewKeeper
	RegisterCodec        = staking.RegisterCodec
	InitGenesis          = staking.InitGenesis
	DefaultParams        = staking.DefaultParams
	NewMultiStakingHooks = staking.NewMultiStakingHooks
)

type (
	Keeper = keeper.Keeper
)
