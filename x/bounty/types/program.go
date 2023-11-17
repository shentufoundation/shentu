package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Programs is an array of program
type Programs []Program

func NewProgram(pid, name, detail string, admin sdk.AccAddress,
	status ProgramStatus, createTime time.Time) (Program, error) {

	return Program{
		ProgramId:    pid,
		Name:         name,
		Detail:       detail,
		AdminAddress: admin.String(),
		Status:       status,
		CreateTime:   createTime,
	}, nil
}

func (m *Program) ValidateBasic() error {
	if len(m.ProgramId) == 0 {
		return ErrProgramID
	}
	if _, err := sdk.AccAddressFromBech32(m.AdminAddress); err != nil {
		return err
	}

	if !ValidProgramStatus(m.Status) {
		return ErrProgramStatusInvalid
	}

	return nil
}

// ValidProgramStatus returns true if the program status is valid and false
// otherwise.
func ValidProgramStatus(status ProgramStatus) bool {
	if status == ProgramStatusInactive ||
		status == ProgramStatusActive ||
		status == ProgramStatusClosed {
		return true
	}
	return false
}
