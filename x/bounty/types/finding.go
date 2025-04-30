package types

import (
	"fmt"
	"strings"
)

// ValidFindingSeverityLevel returns true if the finding level is valid and false otherwise.
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

// SeverityLevelFromString returns a SeverityLevel from a string. It returns an error if the string is invalid.
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
