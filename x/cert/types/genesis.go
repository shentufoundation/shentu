package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(_ sdk.Coin, startingCertificateID uint64) GenesisState {
	return GenesisState{NextCertificateId: startingCertificateID}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{NextCertificateId: 1}
}

// ValidateGenesis verifies the cert module genesis state is internally consistent:
//   - every certificate has a non-zero ID
//   - certificate IDs are unique
//   - NextCertificateId is strictly greater than every imported certificate ID,
//     so future issuance cannot overwrite an imported record
//   - certifier addresses are non-empty and unique
func ValidateGenesis(data GenesisState) error {
	seenCerts := make(map[uint64]struct{}, len(data.Certificates))
	var maxID uint64
	for _, c := range data.Certificates {
		if c.CertificateId == 0 {
			return fmt.Errorf("certificate has zero id")
		}
		if _, dup := seenCerts[c.CertificateId]; dup {
			return fmt.Errorf("duplicate certificate id in genesis: %d", c.CertificateId)
		}
		seenCerts[c.CertificateId] = struct{}{}
		if c.CertificateId > maxID {
			maxID = c.CertificateId
		}
	}
	if data.NextCertificateId <= maxID {
		return fmt.Errorf("nextCertificateId (%d) must be strictly greater than max imported certificate id (%d)",
			data.NextCertificateId, maxID)
	}

	seenCertifiers := make(map[string]struct{}, len(data.Certifiers))
	for _, c := range data.Certifiers {
		addr, err := sdk.AccAddressFromBech32(c.Address)
		if err != nil {
			return fmt.Errorf("invalid certifier address %q: %w", c.Address, err)
		}
		if _, dup := seenCertifiers[addr.String()]; dup {
			return fmt.Errorf("duplicate certifier in genesis: %s", c.Address)
		}
		seenCertifiers[addr.String()] = struct{}{}
	}
	return nil
}

// GetGenesisStateFromAppState returns cert module GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, certificate := range g.Certificates {
		err := certificate.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}
