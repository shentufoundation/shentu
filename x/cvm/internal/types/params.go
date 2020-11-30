package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultGasRate uint64 = 1
)

// Parameter keys
var (
	ParamStoreKeyGasRate = []byte("GasRate")
)

var _ subspace.ParamSet = &Params{}

// Params defines the parameters for the cvm module.
type Params struct {
	GasRate uint64 `json:"gas_rate"`
}

// NewParams creates a new Params object.
func NewParams(gasRate uint64) Params {
	return Params{
		GasRate: gasRate,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs pairs of cvm module's parameters.
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyGasRate, &p.GasRate, validateGasRate),
	}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateGasRate(p.GasRate); err != nil {
		return err
	}
	return nil
}

func validateGasRate(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("invalid gas rate: %d", v)
	}
	return nil
}

// ParamKeyTable for auth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}
