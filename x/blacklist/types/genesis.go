package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState returns an empty blacklist genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{Entries: []Entry{}}
}

// ValidateGenesis performs basic validation of the blacklist genesis state.
func ValidateGenesis(data GenesisState) error {
	seen := make(map[string]struct{}, len(data.Entries))
	for _, entry := range data.Entries {
		if _, err := sdk.AccAddressFromBech32(entry.Address); err != nil {
			return err
		}
		if entry.Category == CATEGORY_UNSPECIFIED {
			return ErrInvalidCategory
		}
		if _, dup := seen[entry.Address]; dup {
			return ErrEntryNotFound.Wrapf("duplicate entry for %s", entry.Address)
		}
		seen[entry.Address] = struct{}{}
	}
	return nil
}
