// Package keeper specifies the keeper for the gov module.
package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/v2/x/gov/types"
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

// Iterators

// IterateActiveProposalsQueue iterates over the proposals in the active proposal queue
// and performs a callback function.
func (k Keeper) IterateActiveProposalsQueue(ctx sdk.Context, endTime time.Time, cb func(proposal types.Proposal) (stop bool)) {
	iterator := k.ActiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := govtypes.SplitActiveProposalQueueKey(iterator.Key())
		proposal, found := k.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateInactiveProposalsQueue iterates over the proposals in the inactive proposal queue
// and performs a callback function.
func (k Keeper) IterateInactiveProposalsQueue(ctx sdk.Context, endTime time.Time, cb func(proposal types.Proposal) (stop bool)) {
	iterator := k.InactiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := govtypes.SplitInactiveProposalQueueKey(iterator.Key())
		proposal, found := k.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateAllDeposits iterates over the all the stored deposits and performs a callback function.
func (k Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit govtypes.Deposit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govtypes.DepositsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit govtypes.Deposit
		k.cdc.MustUnmarshal(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// Tally counts the votes and returns whether the proposal passes and/or if tokens should be burned.
func (k Keeper) Tally(ctx sdk.Context, proposal types.Proposal) (passes bool, burnDeposits bool, tallyResults govtypes.TallyResult) {
	return Tally(ctx, k, proposal)
}

// BondDenom returns the staking denom.
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return k.stakingKeeper.BondDenom(ctx)
}
