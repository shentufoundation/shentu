package types

import (
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/shentufoundation/shentu/v2/common"
)

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Minter: types.InitialMinter(math.LegacyNewDecWithPrec(4, 2)),
		Params: types.NewParams(
			common.MicroCTKDenom,
			math.LegacyNewDecWithPrec(10, 2), // max inflation rate change
			math.LegacyNewDecWithPrec(14, 2), // max inflation rate
			math.LegacyNewDecWithPrec(4, 2),  // min inflation rate
			math.LegacyNewDecWithPrec(67, 2), // target staked coin percentage
			common.BlocksPerYear,             // blocks per year, 5 second block time
		),
	}
}
