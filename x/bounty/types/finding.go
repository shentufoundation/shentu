package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewFinding(pid, fid, title, detail, hash string, operator sdk.AccAddress, createTime time.Time, level SeverityLevel) (Finding, error) {
	return Finding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            title,
		FindingHash:      hash,
		SubmitterAddress: operator.String(),
		SeverityLevel:    level,
		Status:           FindingStatusSubmitted,
		Detail:           detail,
		CreateTime:       createTime,
	}, nil
}

func (m Finding) ValidateBasic() error {
	if len(m.ProgramId) == 0 {
		return ErrProgramID
	}
	if len(m.FindingId) == 0 {
		return ErrFindingID
	}

	if _, err := sdk.AccAddressFromBech32(m.SubmitterAddress); err != nil {
		return err
	}

	if !ValidFindingStatus(m.Status) {
		return ErrFindingStatusInvalid
	}
	if !ValidFindingSeverityLevel(m.SeverityLevel) {
		return ErrFindingSeverityLevelInvalid
	}

	return nil
}

// ValidFindingStatus returns true if the finding status is valid and false
// otherwise.
func ValidFindingStatus(status FindingStatus) bool {
	if status == FindingStatusSubmitted ||
		status == FindingStatusActive ||
		status == FindingStatusConfirmed ||
		status == FindingStatusPaid ||
		status == FindingStatusClosed {
		return true
	}
	return false
}

// ValidFindingSeverityLevel returns true if the finding level is valid and false
// otherwise.
func ValidFindingSeverityLevel(level SeverityLevel) bool {
	if level == Unspecified ||
		level == Critical ||
		level == High ||
		level == Medium ||
		level == Low ||
		level == Informational {
		return true
	}
	return false
}

// SeverityLevelFromString returns a SeverityLevel from a string. It returns an error
// if the string is invalid.
func SeverityLevelFromString(str string) (SeverityLevel, error) {
	option, ok := SeverityLevel_value[str]
	if !ok {
		return Unspecified, fmt.Errorf("'%s' is not a valid SeverityLevel option", str)
	}
	return SeverityLevel(option), nil
}

// NormalizeSeverityLevel - normalize user specified severity level
func NormalizeSeverityLevel(level string) string {
	upperLevel := strings.ToUpper(level)
	switch upperLevel {
	case "UNSPECIFIED":
		return Unspecified.String()
	case "CRITICAL":
		return Critical.String()
	case "HIGH":
		return High.String()
	case "MEDIUM":
		return Medium.String()
	case "LOW":
		return Low.String()
	case "INFORMATIONAL":
		return Informational.String()
	default:
		return level
	}
}
