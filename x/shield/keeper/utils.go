package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPremiumRate(days uint64) sdk.Dec {
	return sdk.NewDecFromBigIntWithPrec(big.NewInt(4), 2) //placeholder 4% for now
}
