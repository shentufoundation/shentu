package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Programs is an array of program
type Programs []Program

func NewProgram(pid, name, desc string,
	admin sdk.ValAddress, members []string, status ProgramStatus, levels []BountyLevel) (Program, error) {

	return Program{
		ProgramId:      pid,
		Name:           name,
		Description:    desc,
		AdminAddress:   admin.String(),
		MemberAccounts: members,
		Status:         status,
		BountyLevels:   levels,
	}, nil
}

func (m *Program) ValidateBasic() error {
	if len(m.ProgramId) == 0 {
		return ErrProgramID
	}
	if _, err := sdk.AccAddressFromBech32(m.AdminAddress); err != nil {
		return err
	}

	return nil
}
