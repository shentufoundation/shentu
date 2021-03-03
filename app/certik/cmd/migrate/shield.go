package migrate

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	shieldtypes "github.com/certikfoundation/shentu/x/shield/types"
)

// Withdraw stores an ongoing withdraw of pool collateral.
type Withdraw struct {
	// Address is the chain address of the provider withdrawing.
	Address sdk.AccAddress `json:"address" yaml:"address"`

	// Amount is the amount of withdraw.
	Amount sdk.Int `json:"amount" yaml:"amount"`

	// CompletionTime is the scheduled withdraw completion time.
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"`
}

type Withdraws []Withdraw

// PoolParams defines the parameters for the shield pool.
type PoolParams struct {
	ProtectionPeriod  time.Duration `json:"protection_period" yaml:"protection_period"`
	ShieldFeesRate    sdk.Dec       `json:"shield_fees_rate" yaml:"shield_fees_rate"`
	WithdrawPeriod    time.Duration `json:"withdraw_period" yaml:"withdraw_period"`
	PoolShieldLimit   sdk.Dec       `json:"pool_shield_limit" yaml:"pool_shield_limit"`
	MinShieldPurchase sdk.Coins     `json:"min_shield_purchase" yaml:"min_shield_purchase"`
}

// ClaimProposalParams defines the parameters for the shield claim proposals.
type ClaimProposalParams struct {
	ClaimPeriod  time.Duration `json:"claim_period" yaml:"claim_period"`
	PayoutPeriod time.Duration `json:"payout_period" yaml:"payout_period"`
	MinDeposit   sdk.Coins     `json:"min_deposit" json:"min_deposit"`
	DepositRate  sdk.Dec       `json:"deposit_rate" yaml:"deposit_rate"`
	FeesRate     sdk.Dec       `json:"fees_rate" yaml:"fees_rate"`
}

// MixedDecCoins defines the struct for mixed coins in decimal with native and foreign decimal coins.
type MixedDecCoins struct {
	Native  sdk.DecCoins `json:"native" yaml:"native"`
	Foreign sdk.DecCoins `json:"foreign" yaml:"foreign"`
}

// Pool contains a shield project pool's data.
type Pool struct {
	// ID is the id of the pool.
	ID uint64 `json:"id" yaml:"id"`

	// Description is the term of the pool.
	Description string `json:"description" yaml:"description"`

	// Sponsor is the project owner of the pool.
	Sponsor string `json:"sponsor" yaml:"sponsor"`

	// SponsorAddress is the CertiK Chain address of the sponsor.
	SponsorAddress sdk.AccAddress `json:"sponsor_address" yaml:"sponsor_address"`

	// ShieldLimit is the maximum shield can be purchased for the pool.
	ShieldLimit sdk.Int `json:"shield_limit" yaml:"shield_limit"`

	// Active means new purchases are allowed.
	Active bool `json:"active" yaml:"active"`

	// Shield is the amount of all active purchased shields.
	Shield sdk.Int `json:"shield" yaml:"shield"`
}

// Provider tracks total delegation, total collateral, and rewards of a provider.
type Provider struct {
	// Address is the address of the provider.
	Address sdk.AccAddress `json:"address" yaml:"address"`

	// DelegationBonded is the amount of bonded delegation.
	DelegationBonded sdk.Int `json:"delegation_bonded" yaml:"delegation_bonded"`

	// Collateral is amount of all collaterals for the provider, including
	// those in withdraw queue but excluding those currently locked, in all
	// pools.
	Collateral sdk.Int `json:"collateral" yaml:"collateral"`

	// TotalLocked is the amount locked for pending claims.
	TotalLocked sdk.Int `json:"total_locked" yaml:"total_locked"`

	// Withdrawing is the amount of collateral in withdraw queues.
	Withdrawing sdk.Int `json:"withdrawing" yaml:"withdrawing"`

	// Rewards is the pooling rewards to be collected.
	Rewards MixedDecCoins `json:"rewards" yaml:"rewards"`
}

// Purchase record an individual purchase.
type Purchase struct {
	// PurchaseID is the purchase_id.
	PurchaseID uint64 `json:"purchase_id" yaml:"purchase_id"`

	// ProtectionEndTime is the time when the protection of the shield ends.
	ProtectionEndTime time.Time `json:"protection_end_time" yaml:"protection_end_time"`

	// DeletionTime is the time when the purchase should be deleted.
	DeletionTime time.Time `json:"deletion_time" yaml:"deletion_time"`

	// Description is the information about the protected asset.
	Description string `json:"description" yaml:"description"`

	// Shield is the unused amount of shield purchased.
	Shield sdk.Int `json:"shield" yaml:"shield"`

	// ServiceFees is the service fees paid by this purchase.
	ServiceFees MixedDecCoins `json:"service_fees" yaml:"service_fees"`
}

// PurchaseList is a collection of purchase.
type PurchaseList struct {
	// PoolID is the id of the shield of the purchase.
	PoolID uint64 `json:"pool_id" yaml:"pool_id"`

	// Purchaser is the address making the purchase.
	Purchaser sdk.AccAddress `json:"purchaser" yaml:"purchaser"`

	// Entries stores all purchases by the purchaser in the pool.
	Entries []Purchase `json:"entries" yaml:"entries"`
}

type ShieldStaking struct {
	PoolID            uint64         `json:"pool_id" yaml:"pool_id"`
	Purchaser         sdk.AccAddress `json:"purchaser" yaml:"purchaser"`
	Amount            sdk.Int        `json:"amount" yaml:"amount"`
	WithdrawRequested sdk.Int        `json:"withdraw_requested" yaml:"withdraw_requested"`
}

type OriginalStaking struct {
	PurchaseID uint64
	Amount     sdk.Int
}

// Reimbursement stores information of a reimbursement.
type Reimbursement struct {
	Amount      sdk.Coins      `json:"amount" yaml:"amount"`
	Beneficiary sdk.AccAddress `json:"beneficiary" yaml:"beneficiary"`
	PayoutTime  time.Time      `json:"payout_time" yaml:"payout_time"`
}

// ProposalIDReimbursementPair stores information of a reimbursement and corresponding proposal ID.
type ProposalIDReimbursementPair struct {
	ProposalID    uint64
	Reimbursement Reimbursement
}

// GenesisState defines the shield genesis state.
type ShieldGenesisState struct {
	ShieldAdmin                  sdk.AccAddress                `json:"shield_admin" yaml:"shield_admin"`
	NextPoolID                   uint64                        `json:"next_pool_id" yaml:"next_pool_id"`
	NextPurchaseID               uint64                        `json:"next_purchase_id" yaml:"next_purchase_id"`
	PoolParams                   PoolParams                    `json:"pool_params" yaml:"pool_params"`
	ClaimProposalParams          ClaimProposalParams           `json:"claim_proposal_params" yaml:"claim_proposal_params"`
	TotalCollateral              sdk.Int                       `json:"total_collateral" yaml:"total_collateral"`
	TotalWithdrawing             sdk.Int                       `json:"total_withdrawing" yaml:"total_withdrawing"`
	TotalShield                  sdk.Int                       `json:"total_shield" yaml:"total_shield"`
	TotalClaimed                 sdk.Int                       `json:"total_claimed" yaml:"total_claimed"`
	ServiceFees                  MixedDecCoins                 `json:"service_fees" yaml:"service_fees"`
	RemainingServiceFees         MixedDecCoins                 `json:"remaining_service_fees" yaml:"remaining_service_fees"`
	Pools                        []Pool                        `json:"pools" yaml:"pools"`
	Providers                    []Provider                    `json:"providers" yaml:"providers"`
	PurchaseLists                []PurchaseList                `json:"purchases" yaml:"purchases"`
	Withdraws                    Withdraws                     `json:"withdraws" yaml:"withdraws"`
	LastUpdateTime               time.Time                     `json:"last_update_time" yaml:"last_update_time"`
	ShieldStakingRate            sdk.Dec                       `json:"shield_staking_rate" yaml:"shield_staking_rate"`
	GlobalStakingPool            sdk.Int                       `json:"global_staking_pool" yaml:"global_staking_pool"`
	StakeForShields              []ShieldStaking               `json:"staking_purchases" yaml:"staking_purchases"`
	OriginalStakings             []OriginalStaking             `json:"original_stakings" yaml:"original_stakings"`
	ProposalIDReimbursementPairs []ProposalIDReimbursementPair `json:"proposalID_reimbursement_pairs" yaml:"proposalID_reimbursement_pairs"`
}

func RegisterShieldLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(ShieldClaimProposal{}, "shield/ShieldClaimProposal", nil)
}

func migrateShield(oldState ShieldGenesisState) *shieldtypes.GenesisState {
	newPoolParams := shieldtypes.PoolParams{
		ProtectionPeriod:  oldState.PoolParams.ProtectionPeriod,
		ShieldFeesRate:    oldState.PoolParams.ShieldFeesRate,
		WithdrawPeriod:    oldState.PoolParams.WithdrawPeriod,
		PoolShieldLimit:   oldState.PoolParams.PoolShieldLimit,
		MinShieldPurchase: oldState.PoolParams.MinShieldPurchase,
	}

	newClaimParams := shieldtypes.ClaimProposalParams{
		ClaimPeriod:  oldState.ClaimProposalParams.ClaimPeriod,
		PayoutPeriod: oldState.ClaimProposalParams.PayoutPeriod,
		MinDeposit:   oldState.ClaimProposalParams.MinDeposit,
		DepositRate:  oldState.ClaimProposalParams.DepositRate,
		FeesRate:     oldState.ClaimProposalParams.FeesRate,
	}

	newServiceFees := shieldtypes.MixedDecCoins{
		Native:  oldState.ServiceFees.Native,
		Foreign: oldState.ServiceFees.Foreign,
	}

	newRemainingServiceFees := shieldtypes.MixedDecCoins{
		Native:  oldState.RemainingServiceFees.Native,
		Foreign: oldState.RemainingServiceFees.Foreign,
	}

	newPools := make([]shieldtypes.Pool, len(oldState.Pools))
	for i, pool := range oldState.Pools {
		newPools[i] = shieldtypes.Pool{
			Id:          pool.ID,
			Description: pool.Description,
			Sponsor:     pool.Sponsor,
			SponsorAddr: pool.SponsorAddress.String(),
			ShieldLimit: pool.ShieldLimit,
			Active:      pool.Active,
			Shield:      pool.Shield,
		}
	}

	newProviders := make([]shieldtypes.Provider, len(oldState.Providers))
	for i, prov := range oldState.Providers {
		newProviders[i] = shieldtypes.Provider{
			Address:          prov.Address.String(),
			DelegationBonded: prov.DelegationBonded,
			Collateral:       prov.Collateral,
			TotalLocked:      prov.TotalLocked,
			Withdrawing:      prov.Withdrawing,
			Rewards: shieldtypes.MixedDecCoins{
				Native:  prov.Rewards.Native,
				Foreign: prov.Rewards.Foreign,
			},
		}
	}

	newPurchaseLists := make([]shieldtypes.PurchaseList, len(oldState.PurchaseLists))
	for i, pl := range oldState.PurchaseLists {
		var newEntries []shieldtypes.Purchase
		for _, entry := range pl.Entries {
			newEntries = append(newEntries, shieldtypes.Purchase{
				PurchaseId:        entry.PurchaseID,
				ProtectionEndTime: entry.ProtectionEndTime,
				DeletionTime:      entry.DeletionTime,
				Description:       entry.Description,
				Shield:            entry.Shield,
				ServiceFees: shieldtypes.MixedDecCoins{
					Native:  entry.ServiceFees.Native,
					Foreign: entry.ServiceFees.Foreign,
				},
			})
		}
		newPurchaseLists[i] = shieldtypes.PurchaseList{
			PoolId:    pl.PoolID,
			Purchaser: pl.Purchaser.String(),
			Entries:   newEntries,
		}
	}

	newWithdraws := make([]shieldtypes.Withdraw, len(oldState.Withdraws))
	for i, wd := range oldState.Withdraws {
		newWithdraws[i] = shieldtypes.Withdraw{
			Address:        wd.Address.String(),
			Amount:         wd.Amount,
			CompletionTime: wd.CompletionTime,
		}
	}

	newStakeForShields := make([]shieldtypes.ShieldStaking, len(oldState.StakeForShields))
	for i, ss := range oldState.StakeForShields {
		newStakeForShields[i] = shieldtypes.ShieldStaking{
			PoolId:            ss.PoolID,
			Purchaser:         ss.Purchaser.String(),
			Amount:            ss.Amount,
			WithdrawRequested: ss.WithdrawRequested,
		}
	}

	newOriginalStakings := make([]shieldtypes.OriginalStaking, len(oldState.OriginalStakings))
	for i, os := range oldState.OriginalStakings {
		newOriginalStakings[i] = shieldtypes.OriginalStaking{
			PurchaseId: os.PurchaseID,
			Amount:     os.Amount,
		}
	}

	newProposalIDReimbursementPairs := make([]shieldtypes.ProposalIDReimbursementPair, len(oldState.ProposalIDReimbursementPairs))
	for i, prp := range oldState.ProposalIDReimbursementPairs {
		newProposalIDReimbursementPairs[i] = shieldtypes.ProposalIDReimbursementPair{
			ProposalId: prp.ProposalID,
			Reimbursement: shieldtypes.Reimbursement{
				Amount:      prp.Reimbursement.Amount,
				Beneficiary: prp.Reimbursement.Beneficiary.String(),
				PayoutTime:  prp.Reimbursement.PayoutTime,
			},
		}
	}

	return &shieldtypes.GenesisState{
		ShieldAdmin:                  oldState.ShieldAdmin.String(),
		NextPoolId:                   oldState.NextPoolID,
		NextPurchaseId:               oldState.NextPurchaseID,
		PoolParams:                   newPoolParams,
		ClaimProposalParams:          newClaimParams,
		TotalCollateral:              oldState.TotalCollateral,
		TotalWithdrawing:             oldState.TotalWithdrawing,
		TotalShield:                  oldState.TotalShield,
		TotalClaimed:                 oldState.TotalClaimed,
		ServiceFees:                  newServiceFees,
		RemainingServiceFees:         newRemainingServiceFees,
		Pools:                        newPools,
		Providers:                    newProviders,
		PurchaseLists:                newPurchaseLists,
		Withdraws:                    newWithdraws,
		LastUpdateTime:               oldState.LastUpdateTime,
		ShieldStakingRate:            oldState.ShieldStakingRate,
		GlobalStakingPool:            oldState.GlobalStakingPool,
		StakeForShields:              newStakeForShields,
		OriginalStakings:             newOriginalStakings,
		ProposalIDReimbursementPairs: newProposalIDReimbursementPairs,
	}
}
