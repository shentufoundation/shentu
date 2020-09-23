package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Pool struct {
	// TODO: better field names
	Community         []Collateral
	BlockChainCompany PoolCreator
	CertiK            Collateral
}

func NewPool(accAddr sdk.AccAddress, coverage, deposit sdk.Coins) Pool {
	creator := PoolCreator{
		Creator:  accAddr,
		Coverage: coverage,
		Premium:  deposit,
	}
	return Pool{
		BlockChainCompany: creator,
	}
}

type PoolCreator struct {
	Creator  sdk.AccAddress
	Coverage sdk.Coins
	Premium  sdk.Coins
}

type Collateral struct {
	Provider sdk.AccAddress
	Amount   sdk.Coins
}
