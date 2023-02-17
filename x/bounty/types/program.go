package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Programs is an array of program
type Programs []Program

// Findings is an array of finding
type Findings []Finding

func (m *Program) ValidateBasic() error {
	if m.ProgramId == 0 ||
		m.EncryptionKey == nil {
		return fmt.Errorf("programId or EncryptionKey is error")
	}
	if _, err := sdk.AccAddressFromBech32(m.CreatorAddress); err != nil {
		return err
	}
	for _, deposit := range m.Deposit {
		if !deposit.IsValid() {
			return fmt.Errorf("deposit is invalid [deposit:%s]", deposit.String())
		}
	}
	return nil
}

func (m *Finding) ValidateBasic() error {
	if m.ProgramId == 0 ||
		m.FindingId == 0 ||
		m.SeverityLevel < 0 ||
		int(m.SeverityLevel) >= len(SeverityLevel_name) ||
		int(m.FindingStatus) >= len(FindingStatus_name) {
		return fmt.Errorf("finding data is error:%s", m.String())
	}
	if _, err := sdk.AccAddressFromBech32(m.SubmitterAddress); err != nil {
		return err
	}
	return nil
}
