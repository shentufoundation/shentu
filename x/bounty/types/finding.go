package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Findings is an array of finding
type Findings []Finding

func NewFinding(programId, findingId string, title, submitterAddr string, detail FindingDetail, submitTime time.Time) (Finding, error) {

	return Finding{
		FindingId:        findingId,
		ProgramId:        programId,
		Title:            title,
		Detail:           detail,
		CreateTime:       submitTime,
		SubmitterAddress: submitterAddr,
		Status:           FindingStatusReported,
	}, nil
}

func (m *Finding) ValidateBasic() error {
	if len(m.ProgramId) == 0 {
		return ErrProgramID
	} else if len(m.FindingId) == 0 {
		return ErrFindingID
	} else if m.Detail.SeverityLevel < 0 || int(m.Detail.SeverityLevel) >= len(SeverityLevel_name) {
		return fmt.Errorf("invalid SeverityLevel:%d", m.Detail.SeverityLevel)
	} else if int(m.Status) >= len(FindingStatus_name) {
		return fmt.Errorf("invalid finding status:%d", m.Status)
	}

	if _, err := sdk.AccAddressFromBech32(m.SubmitterAddress); err != nil {
		return err
	}
	return nil
}

func NewFindingDetail(desc, poc string, targets []string, level SeverityLevel) FindingDetail {

	return FindingDetail{
		Description:    desc,
		ProofOfConcept: poc,
		ProgramTargets: targets,
		SeverityLevel:  level,
	}
}
