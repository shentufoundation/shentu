package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/certikfoundation/shentu/common"
)

// default parameter values
var (
	// default values for Shield pool's parameters
	DefaultProtectionPeriod = time.Hour * 24 * 14      // 14 days
	DefaultMinPoolLife      = time.Hour * 24 * 56      // 56 days
	DefaultShieldFeesRate   = sdk.NewDecWithPrec(1, 2) // 1%

	// default values for Shield claim proposal's parameters
	DefaultClaimPeriod              = time.Hour * 24 * 21                                                    // 21 days
	DefaultPayoutPeriod             = time.Hour * 24 * 56                                                    // 56 days
	DefaultMinClaimProposalDeposit  = sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(100000000))) // 100 CTK
	DefaultClaimProposalDepositRate = sdk.NewDecWithPrec(10, 2)                                              // 10%
	DefaultClaimProposalFeesRate    = sdk.NewDecWithPrec(1, 2)                                               // 1%
)

// parameter keys
var (
	ParamStoreKeyPoolParams          = []byte("shieldpoolparams")
	ParamStoreKeyClaimProposalParams = []byte("claimproposalparams")
)

// ParamKeyTable is the key declaration for parameters.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamStoreKeyPoolParams, PoolParams{}, validatePoolParams),
		params.NewParamSetPair(ParamStoreKeyClaimProposalParams, ClaimProposalParams{}, validateClaimProposalParams),
	)
}

// PoolParams defines the parameters for the shield pool.
type PoolParams struct {
	ProtectionPeriod time.Duration `json:"protection_period" yaml:"protection_period"`
	MinPoolLife      time.Duration `json:"min_pool_life" yaml:"min_pool_life"`
	ShieldFeesRate   sdk.Dec       `json:"shield_fees_rate" yaml:"shield_fees_rate"`
}

// NewPoolParams creates a new PoolParams object.
func NewPoolParams(protectionPeriod, minPoolLife time.Duration, shieldFeesRate sdk.Dec) PoolParams {
	return PoolParams{
		ProtectionPeriod: protectionPeriod,
		MinPoolLife:      minPoolLife,
		ShieldFeesRate:   shieldFeesRate,
	}
}

// DefaultClaimProposalParams returns a default PoolParams instance.
func DefaultPoolParams() PoolParams {
	return NewPoolParams(DefaultProtectionPeriod, DefaultMinPoolLife, DefaultShieldFeesRate)
}

func validatePoolParams(i interface{}) error {
	// TODO
	return nil
}

// ClaimProposalParams defines the parameters for the shield claim proposals.
type ClaimProposalParams struct {
	ClaimPeriod  time.Duration `json:"claim_period" yaml:"claim_period"`
	PayoutPeriod time.Duration `json:"payout_period" yaml:"payout_period"`
	MinDeposit   sdk.Coins     `json:"min_deposit" json:"min_deposit"`
	DepositRate  sdk.Dec       `json:"deposit_rate" yaml:"deposit_rate"`
	FeesRate     sdk.Dec       `json:"fees_rate" yaml:"fees_rate"`
}

// NewClaimProposalParams creates a new ClaimProposalParams instance.
func NewClaimProposalParams(claimPeriod, payoutPeriod time.Duration,
	minDeposit sdk.Coins, depositRate, feesRate sdk.Dec) ClaimProposalParams {
	return ClaimProposalParams{
		ClaimPeriod:  claimPeriod,
		PayoutPeriod: payoutPeriod,
		MinDeposit:   minDeposit,
		DepositRate:  depositRate,
		FeesRate:     feesRate,
	}
}

// DefaultClaimProposalParams returns a default ClaimProposalParams instance.
func DefaultClaimProposalParams() ClaimProposalParams {
	return NewClaimProposalParams(DefaultClaimPeriod, DefaultPayoutPeriod,
		DefaultMinClaimProposalDeposit, DefaultClaimProposalDepositRate, DefaultClaimProposalFeesRate)
}

func validateClaimProposalParams(i interface{}) error {
	// TODO
	return nil
}
