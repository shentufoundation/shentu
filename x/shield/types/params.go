package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/certikfoundation/shentu/common"
)

// default parameter values
var (
	// default values for Shield pool's parameters
	DefaultProtectionPeriod  = time.Hour * 24 * 21                                                   // 21 days
	DefaultShieldFeesRate    = sdk.NewDecWithPrec(769, 5)                                            // 0.769%
	DefaultWithdrawPeriod    = time.Hour * 24 * 21                                                   // 21 days
	DefaultPoolShieldLimit   = sdk.NewDecWithPrec(50, 2)                                             // 50%
	DefaultMinShieldPurchase = sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(50000000))) // 50 CTK

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
	ProtectionPeriod  time.Duration `json:"protection_period" yaml:"protection_period"`
	ShieldFeesRate    sdk.Dec       `json:"shield_fees_rate" yaml:"shield_fees_rate"`
	WithdrawPeriod    time.Duration `json:"withdraw_period" yaml:"withdraw_period"`
	PoolShieldLimit   sdk.Dec       `json:"pool_shield_limit" yaml:"pool_shield_limit"`
	MinShieldPurchase sdk.Coins     `json:"min_shield_purchase" yaml:"min_shield_purchase"`
}

// NewPoolParams creates a new PoolParams object.
func NewPoolParams(protectionPeriod, withdrawPeriod time.Duration, shieldFeesRate sdk.Dec, poolShieldLimit sdk.Dec, minShieldPurchase sdk.Coins) PoolParams {
	return PoolParams{
		ProtectionPeriod:  protectionPeriod,
		ShieldFeesRate:    shieldFeesRate,
		WithdrawPeriod:    withdrawPeriod,
		PoolShieldLimit:   poolShieldLimit,
		MinShieldPurchase: minShieldPurchase,
	}
}

// DefaultPoolParams returns a default PoolParams instance.
func DefaultPoolParams() PoolParams {
	return NewPoolParams(DefaultProtectionPeriod, DefaultWithdrawPeriod, DefaultShieldFeesRate, DefaultPoolShieldLimit, DefaultMinShieldPurchase)
}

func validatePoolParams(i interface{}) error {
	v, ok := i.(PoolParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	protectionPeriod := v.ProtectionPeriod
	shieldFeesRate := v.ShieldFeesRate
	withdrawPeriod := v.WithdrawPeriod
	poolShieldLimit := v.PoolShieldLimit
	minShieldPurchase := v.MinShieldPurchase

	if protectionPeriod <= 0 {
		return fmt.Errorf("protection period must be positive: %s", protectionPeriod)
	}
	if shieldFeesRate.IsNegative() || shieldFeesRate.GT(sdk.OneDec()) {
		return fmt.Errorf("shield fees rate should be positive and less or equal to one but is %s", shieldFeesRate)
	}
	if withdrawPeriod <= 0 {
		return fmt.Errorf("withdraw period must be positive: %s", withdrawPeriod)
	}
	if poolShieldLimit.IsNegative() || poolShieldLimit.GT(sdk.OneDec()) {
		return fmt.Errorf("pool shield limit should be positive and less or equal to one but is %s", poolShieldLimit)
	}
	if !minShieldPurchase.IsValid() {
		return fmt.Errorf("minimum shield purchase must be a valid sdk.Coins, is %s", minShieldPurchase.String())
	}

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
func NewClaimProposalParams(claimPeriod, payoutPeriod time.Duration, minDeposit sdk.Coins, depositRate, feesRate sdk.Dec) ClaimProposalParams {
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
	v, ok := i.(ClaimProposalParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	claimPeriod := v.ClaimPeriod
	payoutPeriod := v.PayoutPeriod
	minDeposit := v.MinDeposit
	depositRate := v.DepositRate
	feesRate := v.FeesRate

	if claimPeriod <= 0 {
		return fmt.Errorf("claim period must be positive: %s", claimPeriod)
	}
	if payoutPeriod <= 0 {
		return fmt.Errorf("payout period must be positive: %s", payoutPeriod)
	}
	if !minDeposit.IsValid() {
		return fmt.Errorf("minimum deposit amount must be a valid sdk.Coins amount, is %s",
			minDeposit.String())
	}
	if depositRate.IsNegative() || depositRate.GT(sdk.OneDec()) {
		return fmt.Errorf("deposit rate should be positive and less or equal to one but is %s",
			depositRate.String())
	}
	if feesRate.IsNegative() || feesRate.GT(sdk.OneDec()) {
		return fmt.Errorf("fees rate should be positive and less or equal to one but is %s",
			feesRate.String())
	}

	return nil
}
