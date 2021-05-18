package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

func NewGenesisState(collections []nfttypes.Collection, admins []Admin) *GenesisState {
	return &GenesisState{
		Collections: collections,
		Admin:       admins,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState([]nfttypes.Collection{}, []Admin{})
}

func ValidateGenesis(data GenesisState) error {
	for _, c := range data.Collections {
		if err := nfttypes.ValidateDenomID(c.Denom.Name); err != nil {
			return err
		}

		for _, nft := range c.NFTs {
			if nft.GetOwner().Empty() {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner")
			}

			if err := nfttypes.ValidateTokenID(nft.GetID()); err != nil {
				return err
			}

			if err := nfttypes.ValidateTokenURI(nft.GetURI()); err != nil {
				return err
			}
		}
	}
	for _, a := range data.Admin {
		_, err := sdk.AccAddressFromBech32(a.Address)
		if err != nil {
			return err
		}
	}
	return nil
}
