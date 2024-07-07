package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// parameter keys
var (
	ParamStoreKeyStakingShieldRate = []byte("stakingshieldrateparams")
)

// ParamKeyTable is the key declaration for parameters.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(ParamStoreKeyStakingShieldRate, sdk.Dec{}, validateStakingShieldRateParams),
	)
}

func validateStakingShieldRateParams(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LTE(sdk.ZeroDec()) {
		return fmt.Errorf("staking shield rate should be greater than 0")
	}
	return nil
}
