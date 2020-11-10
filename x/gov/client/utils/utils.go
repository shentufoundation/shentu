// Package utils defines the utils for client services for the gov module.
package utils

import (
	"strings"

	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
