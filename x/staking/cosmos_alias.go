package staking

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
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
	NewKeeper            = staking.NewKeeper
	RegisterCodec        = staking.RegisterCodec
	InitGenesis          = staking.InitGenesis
	DefaultParams        = staking.DefaultParams
	NewMultiStakingHooks = staking.NewMultiStakingHooks
)

type (
	Keeper = staking.Keeper
)
