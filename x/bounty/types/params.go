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
func NewParams(minGrant, minDeposit []sdk.Coin, theoremMaxProofPeriod, proofMaxLockPeriod time.Duration, complexityFee sdk.Coin, maxComplexity int64, complexityFeeRocq, complexityFeeLean sdk.Coin) Params {
	return Params{
		MinGrant:              minGrant,
		MinDeposit:            minDeposit,
		TheoremMaxProofPeriod: &theoremMaxProofPeriod,
		ProofMaxLockPeriod:    &proofMaxLockPeriod,
		ComplexityFee:         complexityFee,
		MaxComplexity:         maxComplexity,
		ComplexityFeeRocq:     complexityFeeRocq,
		ComplexityFeeLean:     complexityFeeLean,
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
	// Default complexity fee: 10000uctk (deprecated, kept for compatibility)
	complexityFee := sdk.NewCoin("uctk", sdkmath.NewInt(10000))
	// Default max complexity: 1000000
	maxComplexity := int64(1000000)
	// Default per-type complexity fees
	complexityFeeRocq := sdk.NewCoin("uctk", sdkmath.NewInt(10000))
	complexityFeeLean := sdk.NewCoin("uctk", sdkmath.NewInt(10000))

	return NewParams(minGrant, minDeposit, theoremMaxProofPeriod, proofMaxLockPeriod, complexityFee, maxComplexity, complexityFeeRocq, complexityFeeLean)
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

	if !p.ComplexityFee.IsValid() {
		return fmt.Errorf("complexity fee is invalid: %s", p.ComplexityFee)
	}

	if p.MaxComplexity <= 0 {
		return fmt.Errorf("max complexity must be positive, got %d", p.MaxComplexity)
	}

	if !p.ComplexityFeeRocq.IsValid() {
		return fmt.Errorf("complexity fee rocq is invalid: %s", p.ComplexityFeeRocq)
	}

	if !p.ComplexityFeeLean.IsValid() {
		return fmt.Errorf("complexity fee lean is invalid: %s", p.ComplexityFeeLean)
	}

	return nil
}

// GetComplexityFeeByType returns the complexity fee for the given theorem type.
func (p Params) GetComplexityFeeByType(t TheoremType) (sdk.Coin, error) {
	switch t {
	case TheoremType_THEOREM_TYPE_ROCQ:
		return p.ComplexityFeeRocq, nil
	case TheoremType_THEOREM_TYPE_LEAN:
		return p.ComplexityFeeLean, nil
	default:
		return sdk.Coin{}, fmt.Errorf("unsupported theorem type: %s", t)
	}
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
