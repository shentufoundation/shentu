// Package keeper specifies the keeper for the gov module.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

// Keeper implements keeper for the governance module.
type Keeper struct {
	govkeeper.Keeper

	// the SupplyKeeper to reduce the supply of the network
	bankKeeper govtypes.BankKeeper

	// the reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	stakingKeeper types.StakingKeeper

	// the reference to get information about certifiers
	CertKeeper types.CertKeeper

	// the (unexposed) keys used to access the stores from the Context
	storeKey storetypes.StoreKey

	// codec for binary encoding/decoding
	cdc codec.BinaryCodec

	// Legacy Proposal router
	legacyRouter v1beta1.Router

	// Msg server router
	router *baseapp.MsgServiceRouter

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
	cdc codec.BinaryCodec, key storetypes.StoreKey, bankKeeper govtypes.BankKeeper,
	stakingKeeper types.StakingKeeper, certKeeper types.CertKeeper,
	authKeeper govtypes.AccountKeeper, legacyRouter v1beta1.Router, router *baseapp.MsgServiceRouter,
	config govtypes.Config, authority string,
) Keeper {
	cosmosKeeper := govkeeper.NewKeeper(cdc, key, authKeeper, bankKeeper, stakingKeeper, router, config, authority)
	return Keeper{
		Keeper:        *cosmosKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		CertKeeper:    certKeeper,
		storeKey:      key,
		cdc:           cdc,
		legacyRouter:  legacyRouter,
		router:        router,
		config:        config,
		authority:     authority,
	}
}
