package types

import "github.com/irisnet/irismod/modules/nft/types"

func NewGenesisState(collections []types.Collection, admins []Admin) *GenesisState {
	return &GenesisState{
		Collections: collections,
		Admin: admins,
	}
}