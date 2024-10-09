// Package keeper specifies the keeper for the gov module.
package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

// Keeper implements keeper for the governance module.
type Keeper struct {
	govkeeper.Keeper

	authKeeper types.AccountKeeper
	bankKeeper govtypes.BankKeeper

	// the reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	stakingKeeper types.StakingKeeper

	// the reference to get information about certifiers
	CertKeeper types.CertKeeper

	// the (unexposed) keys used to access the stores from the Context
	storeService store.KVStoreService

	// The codec for binary encoding/decoding.
	cdc codec.Codec

	// Legacy Proposal router
	legacyRouter v1beta1.Router

	// Msg server router
	router baseapp.MessageRouter

	config govtypes.Config

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(
	cdc codec.Codec, storeService store.KVStoreService, bankKeeper govtypes.BankKeeper,
	stakingKeeper types.StakingKeeper, certKeeper types.CertKeeper,
	authKeeper govtypes.AccountKeeper, distrKeeper types.DistributionKeeper, legacyRouter v1beta1.Router,
	router baseapp.MessageRouter, config govtypes.Config, authority string,
) Keeper {
	cosmosKeeper := govkeeper.NewKeeper(cdc, storeService, authKeeper, bankKeeper, stakingKeeper, distrKeeper, router, config, authority)
	return Keeper{
		Keeper:        *cosmosKeeper,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		CertKeeper:    certKeeper,
		storeService:  storeService,
		cdc:           cdc,
		legacyRouter:  legacyRouter,
		router:        router,
		config:        config,
		authority:     authority,
	}
}
