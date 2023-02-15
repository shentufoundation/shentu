package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Programs is an array of program
type Programs []Program

// Findings is an array of finding
type Findings []Finding

func (m *Program) Valid() bool {
	if m.ProgramId == 0 ||
		m.EncryptionKey == nil {
		return false
	}
	if _, err := sdk.AccAddressFromBech32(m.CreatorAddress); err != nil {
		return false
	}
	for _, deposit := range m.Deposit {
		if !deposit.IsValid() {
			return false
		}
	}
	return true
}

func (m *Finding) Valid() bool {
	if m.ProgramId == 0 ||
		m.FindingId == 0 ||
		m.SeverityLevel < 0 ||
		int(m.SeverityLevel) >= len(SeverityLevel_name) ||
		int(m.FindingStatus) >= len(FindingStatus_name) {
		return false
	}
	if _, err := sdk.AccAddressFromBech32(m.SubmitterAddress); err != nil {
		return false
	}
	return true
}
