// Package keeper specifies the keeper for the gov module.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

// Keeper implements keeper for the governance module.
type Keeper struct {
	govkeeper.Keeper

	// the reference to the ParamSpace to get and set gov specific params
	paramSpace types.ParamSubspace

	// the SupplyKeeper to reduce the supply of the network
	bankKeeper govtypes.BankKeeper

	// the reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	stakingKeeper types.StakingKeeper

	// the reference to get information about certifiers
	CertKeeper types.CertKeeper

	// the reference to get claim proposal parameters
	ShieldKeeper types.ShieldKeeper

	// the (unexposed) keys used to access the stores from the Context
	storeKey storetypes.StoreKey

	// codec for binary encoding/decoding
	cdc codec.BinaryCodec

	// Legacy Proposal router
	legacyRouter v1beta1.Router

	// Msg server router
	router *baseapp.MsgServiceRouter

	config govtypes.Config
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace types.ParamSubspace, bankKeeper govtypes.BankKeeper,
	stakingKeeper types.StakingKeeper, certKeeper types.CertKeeper, shieldKeeper types.ShieldKeeper,
	authKeeper govtypes.AccountKeeper, legacyRouter v1beta1.Router, router *baseapp.MsgServiceRouter, config govtypes.Config,
) Keeper {
	cosmosKeeper := govkeeper.NewKeeper(cdc, key, paramSpace, authKeeper, bankKeeper, stakingKeeper, legacyRouter, router, config)
	return Keeper{
		Keeper:        cosmosKeeper,
		paramSpace:    paramSpace,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		CertKeeper:    certKeeper,
		ShieldKeeper:  shieldKeeper,
		storeKey:      key,
		cdc:           cdc,
		legacyRouter:  legacyRouter,
		router:        router,
		config:        config,
	}
}

// Tally counts the votes and returns whether the proposal passes and/or if tokens should be burned.
func (k Keeper) Tally(ctx sdk.Context, proposal v1.Proposal) (passes bool, burnDeposits bool, tallyResults v1.TallyResult) {
	return Tally(ctx, k, proposal)
}
