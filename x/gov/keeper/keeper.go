// Package keeper specifies the keeper for the gov module.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

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
	storeKey sdk.StoreKey

	// codec for binary encoding/decoding
	cdc codec.BinaryCodec

	// Proposal router
	router govtypes.Router
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace types.ParamSubspace, bankKeeper govtypes.BankKeeper,
	stakingKeeper types.StakingKeeper, certKeeper types.CertKeeper, shieldKeeper types.ShieldKeeper,
	authKeeper govtypes.AccountKeeper, router govtypes.Router,
) Keeper {
	cosmosKeeper := govkeeper.NewKeeper(cdc, key, paramSpace, authKeeper, bankKeeper, stakingKeeper, router)
	return Keeper{
		Keeper:        cosmosKeeper,
		storeKey:      key,
		paramSpace:    paramSpace,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		CertKeeper:    certKeeper,
		ShieldKeeper:  shieldKeeper,
		cdc:           cdc,
		router:        router,
	}
}

// BondDenom returns the staking denom.
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return k.stakingKeeper.BondDenom(ctx)
}

// IsCertifier checks if the input address is a certifier.
func (k Keeper) IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertifier(ctx, addr)
}

// IsCertifiedIdentity checks if the input address is a certified identity.
func (k Keeper) IsCertifiedIdentity(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertified(ctx, addr.String(), "identity")
}
