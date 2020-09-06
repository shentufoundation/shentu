// Package utils defines the utils for client services for the gov module.
package utils

import (
	"strings"

	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

// NormalizeVoteOption normalizes user specified vote option
func NormalizeVoteOption(option string) string {
	switch strings.ToLower(option) {
	case "yes":
		return govTypes.OptionYes.String()
	case "abstain":
		return govTypes.OptionAbstain.String()
	case "no":
		return govTypes.OptionNo.String()
	case "nowithveto", "no_with_veto":
		return govTypes.OptionNoWithVeto.String()
	default:
		return option
	}
}

// NormalizeProposalType normalizes user specified proposal type.
func NormalizeProposalType(proposalType string) string {
	switch strings.ToLower(proposalType) {
	case "text":
		return govTypes.ProposalTypeText
	case "softwareupgrade":
		return upgrade.ProposalTypeSoftwareUpgrade
	default:
		return proposalType
	}
}

// NormalizeProposalStatus normalizes user specified proposal status.
func NormalizeProposalStatus(status string) string {
	status = strings.ToUpper(status)
	switch status {
	case "DEPOSITPERIOD", "DEPOSIT_PERIOD":
		return "DepositPeriod"
	case "VOTINGPERIOD", "VOTING_PERIOD":
		return "VotingPeriod"
	case "PASSED":
		return "Passed"
	case "REJECTED":
		return "Rejected"
	case "FAILED":
		return "Failed"
	case "SECONDVOTINGPERIOD", "SECOND_VOTING_PERIOD":
		return "SecondVotingPeriod"
	}
	return ""
}
