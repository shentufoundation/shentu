package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Findings is an array of finding
type Findings []Finding

func NewFinding(pid, fid, title, desc string, operator sdk.ValAddress, submitTime time.Time, level SeverityLevel) (Finding, error) {

	hash := sha256.Sum256([]byte(title + desc))
	bzHash := hash[:]
	hashString := hex.EncodeToString(bzHash)

	return Finding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            title,
		Description:      desc,
		SubmitterAddress: operator.String(),
		CreateTime:       submitTime,
		Status:           FindingStatusReported,
		FindingHash:      hashString,
		SeverityLevel:    level,
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
