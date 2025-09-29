package types

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultStartingTheoremID is 1
	DefaultStartingTheoremID uint64 = 1

	// MaxTermComplexity is the maximum allowed term complexity to prevent overflow
	// and ensure reasonable reward calculations. Set to 1 million as a safe upper bound.
	MaxTermComplexity int64 = 1000000
)

// NewParams creates a new Params instance
func NewParams(minGrant, minDeposit []sdk.Coin, theoremMaxProofPeriod, proofMaxLockPeriod time.Duration, complexityFee sdk.Coin) Params {
	return Params{
		MinGrant:              minGrant,
		MinDeposit:            minDeposit,
		TheoremMaxProofPeriod: &theoremMaxProofPeriod,
		ProofMaxLockPeriod:    &proofMaxLockPeriod,
		ComplexityFee:         complexityFee,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	// Default minimum grant: 1000000uctk
	minGrant := sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000000)))
	// Default minimum deposit: 1000000uctk
	minDeposit := sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000000)))
	// Default theorem max proof period: 120 days
	theoremMaxProofPeriod := 120 * 24 * time.Hour
	// Default proof max lock period: 15 minutes
	proofMaxLockPeriod := 15 * time.Minute
	// Default complexity fee: 10000uctk
	complexityFee := sdk.NewCoin("uctk", sdkmath.NewInt(10000))

	return NewParams(minGrant, minDeposit, theoremMaxProofPeriod, proofMaxLockPeriod, complexityFee)
}

// Validate performs validation on params
func (p Params) Validate() error {
	if err := validateMinGrant(p.MinGrant); err != nil {
		return err
	}

	if err := validateMinDeposit(p.MinDeposit); err != nil {
		return err
	}

	if p.TheoremMaxProofPeriod == nil {
		return fmt.Errorf("theorem max proof period cannot be nil")
	}

	if p.ProofMaxLockPeriod == nil {
		return fmt.Errorf("proof max lock period cannot be nil")
	}

	return nil
}

func validateMinGrant(i interface{}) error {
	v, ok := i.([]sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("min grant cannot be empty")
	}

	for _, coin := range v {
		if !coin.IsValid() {
			return fmt.Errorf("invalid min grant: %s", coin)
		}
	}

	return nil
}

func validateMinDeposit(i interface{}) error {
	v, ok := i.([]sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("min deposit cannot be empty")
	}

	for _, coin := range v {
		if !coin.IsValid() {
			return fmt.Errorf("invalid min deposit: %s", coin)
		}
	}

	return nil
}
