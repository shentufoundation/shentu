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
)

// NewParams creates a new Params instance
func NewParams(minGrant, minDeposit []sdk.Coin, theoremMaxProofPeriod, proofMaxLockPeriod time.Duration, checkerRate sdkmath.LegacyDec) Params {
	return Params{
		MinGrant:              minGrant,
		MinDeposit:            minDeposit,
		TheoremMaxProofPeriod: &theoremMaxProofPeriod,
		ProofMaxLockPeriod:    &proofMaxLockPeriod,
		CheckerRate:           checkerRate,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	// Default minimum grant: 1000000 uctk
	minGrant := sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000000)))
	// Default minimum deposit: 1000000 uctk
	minDeposit := sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000000)))
	// Default theorem max proof period: 120 days
	theoremMaxProofPeriod := 120 * 24 * time.Hour
	// Default proof max lock period: 15 minutes
	proofMaxLockPeriod := 15 * time.Minute
	// Default checker rate: 0%
	checkerRate := sdkmath.LegacyMustNewDecFromStr("0.0")

	return NewParams(minGrant, minDeposit, theoremMaxProofPeriod, proofMaxLockPeriod, checkerRate)
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

	return validateCheckerRate(p.CheckerRate)
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

func validateCheckerRate(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("checker rate cannot be negative: %s", v)
	}

	if v.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("checker rate cannot be greater than 1: %s", v)
	}

	return nil
}
