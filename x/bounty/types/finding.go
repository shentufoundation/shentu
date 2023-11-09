package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Findings is an array of finding
type Findings []Finding

func NewFinding(pid, fid, title, detail, hash string, operator sdk.AccAddress, submitTime time.Time, level SeverityLevel) (Finding, error) {

	return Finding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            title,
		FindingHash:      hash,
		SubmitterAddress: operator.String(),
		SeverityLevel:    level,
		Status:           FindingStatusSubmitted,
		Detail:           detail,
		CreateTime:       submitTime,
	}, nil
}

func (m Finding) ValidateBasic() error {
	if len(m.ProgramId) == 0 {
		return ErrProgramID
	} else if len(m.FindingId) == 0 {
		return ErrFindingID
	} else if m.SeverityLevel < 0 || int(m.SeverityLevel) >= len(SeverityLevel_name) {
		return fmt.Errorf("invalid SeverityLevel:%d", m.SeverityLevel)
	} else if int(m.Status) >= len(FindingStatus_name) {
		return fmt.Errorf("invalid finding status:%d", m.Status)
	}

	if _, err := sdk.AccAddressFromBech32(m.SubmitterAddress); err != nil {
		return err
	}
	return nil
}
