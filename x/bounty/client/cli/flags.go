package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName        = "name"
	FlagMembers     = "members"
	FlagDesc        = "desc"
	FlagScopeRules  = "scope-rules"
	FlagKnownIssues = "known-issues"

	FlagFindingTitle = "title"

	FlagFindingSeverityLevel = "severity-level"
	FlagFindingPoc           = "poc"
	FlagComment              = "comment"

	FlagFindingAddress   = "finding-address"
	FlagSubmitterAddress = "submitter-address"
	FlagProgramID        = "program-id"
	FlagFindingID        = "finding-id"
)

func FlagSetProgramDetail() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagDesc, "", "The program's description")
	fs.String(FlagScopeRules, "", "Indicating the scope of bug hunting")
	fs.String(FlagKnownIssues, "", "The existing known concerns")

	return fs
}
