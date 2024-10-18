package types

import (
	"fmt"

	"cosmossdk.io/math"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// parameter keys
var (
	ParamStoreKeyStakingShieldRate = []byte("stakingshieldrateparams")
)

// ParamKeyTable is the key declaration for parameters.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(ParamStoreKeyStakingShieldRate, math.LegacyDec{}, validateStakingShieldRateParams),
	)
}

func validateStakingShieldRateParams(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("staking shield rate should be greater than 0")
	}
	return nil
}
