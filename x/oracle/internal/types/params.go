package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	ParamsStoreKeyTaskParams = []byte("taskparams")
	ParamsStoreKeyPoolParams = []byte("poolparams")
)

// Default parameters
var (
	MinScore                 = sdk.NewInt(0)
	MaxScore                 = sdk.NewInt(100)
	DefaultThresholdScore    = sdk.NewInt(50)
	DefaultAggregationResult = sdk.NewInt(50)

	DefaultExpirationDuration = time.Duration(24) * time.Hour
	DefaultAggregationWindow  = int64(20)
	DefaultEpsilon1           = sdk.NewInt(1)
	DefaultEpsilon2           = sdk.NewInt(100)

	DefaultLockedInBlocks    = int64(30)
	DefaultMinimumCollateral = int64(50000)
)

// ParamKeyTable is the key declaration for parameters.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamsStoreKeyTaskParams, TaskParams{}, validateTaskParams),
		params.NewParamSetPair(ParamsStoreKeyPoolParams, LockedPoolParams{}, validatePoolParams),
	)
}

// NewTaskParams returns a TaskParams object.
func NewTaskParams(expirationDuration time.Duration, aggregationWindow int64, aggregationResult,
	thresholdScore, epsilon1, epsilon2 sdk.Int) TaskParams {
	return TaskParams{
		ExpirationDuration: expirationDuration,
		AggregationWindow:  aggregationWindow,
		AggregationResult:  aggregationResult,
		ThresholdScore:     thresholdScore,
		Epsilon1:           epsilon1,
		Epsilon2:           epsilon2,
	}
}

// DefaultTaskParams generates default set for TaskParams.
func DefaultTaskParams() TaskParams {
	return NewTaskParams(DefaultExpirationDuration, DefaultAggregationWindow,
		DefaultAggregationResult, DefaultThresholdScore, DefaultEpsilon1, DefaultEpsilon2)
}

func validateTaskParams(i interface{}) error {
	taskParams, ok := i.(TaskParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if taskParams.ExpirationDuration < 0 ||
		taskParams.AggregationWindow < 0 ||
		taskParams.ThresholdScore.GT(MaxScore) ||
		taskParams.Epsilon1.LT(sdk.NewInt(0)) ||
		taskParams.Epsilon2.LT(sdk.NewInt(0)) {
		return ErrInvalidTaskParams
	}
	return nil
}

// NewLockedPoolParams returns a LockedPoolParams object.
func NewLockedPoolParams(lockedInBlocks, minimumCollateral int64) LockedPoolParams {
	return LockedPoolParams{
		LockedInBlocks:    lockedInBlocks,
		MinimumCollateral: minimumCollateral,
	}
}

// DefaultLockedPoolParams generates default set for LockedPoolParams
func DefaultLockedPoolParams() LockedPoolParams {
	return NewLockedPoolParams(DefaultLockedInBlocks, DefaultMinimumCollateral)
}

func validatePoolParams(i interface{}) error {
	poolParams, ok := i.(LockedPoolParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if poolParams.LockedInBlocks < 0 {
		return ErrInvalidPoolParams
	}
	if poolParams.MinimumCollateral < 0 {
		return ErrInvalidPoolParams
	}
	return nil
}
