package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNormalizeVoteOption(t *testing.T) {
	t.Run("option is valid for yes", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("yes")
		assert.Equal(t, gotOptionType, "Yes")
	})
	t.Run("option is valid for abstain", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("Abstain")
		assert.Equal(t, gotOptionType, "Abstain")
	})
	t.Run("option is valid for no", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("No")
		assert.Equal(t, gotOptionType, "No")
	})

	t.Run("option is valid for noWithVeto", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("no_with_veto")
		assert.Equal(t, gotOptionType, "NoWithVeto")
	})
	t.Run("option is valid for default", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("")
		assert.Equal(t, gotOptionType, "")
	})
}

func TestNormalizeProposalType(t *testing.T) {
	t.Run("type is valid for text", func(t *testing.T) {
		gotOptionType := NormalizeProposalType("text")
		assert.Equal(t, gotOptionType, "Text")
	})
	t.Run("type is valid for software upgrade", func(t *testing.T) {
		gotOptionType := NormalizeProposalType("softwareupgrade")
		assert.Equal(t, gotOptionType, "SoftwareUpgrade")
	})
	t.Run("type is valid for default", func(t *testing.T) {
		gotOptionType := NormalizeProposalType("")
		assert.Equal(t, gotOptionType, "")
	})
}

func TestNormalizeProposalStatus(t *testing.T) {
	t.Run("status is valid for DepositPeriod", func(t *testing.T) {
		gotOptionType := NormalizeProposalStatus("DepositPeriod")
		assert.Equal(t, gotOptionType, "DepositPeriod")
	})
	t.Run("status is valid for VotingPeriod", func(t *testing.T) {
		gotOptionType := NormalizeProposalStatus("VotingPeriod")
		assert.Equal(t, gotOptionType, "VotingPeriod")
	})
	t.Run("status is valid for Passed", func(t *testing.T) {
		gotOptionType := NormalizeProposalStatus("Passed")
		assert.Equal(t, gotOptionType, "Passed")
	})
	t.Run("status is valid for Rejected", func(t *testing.T) {
		gotOptionType := NormalizeProposalStatus("Rejected")
		assert.Equal(t, gotOptionType, "Rejected")
	})
	t.Run("status is valid for Failed", func(t *testing.T) {
		gotOptionType := NormalizeProposalStatus("Failed")
		assert.Equal(t, gotOptionType, "Failed")
	})
	t.Run("status is valid for SecondVotingPeriod", func(t *testing.T) {
		gotOptionType := NormalizeProposalStatus("SecondVotingPeriod")
		assert.Equal(t, gotOptionType, "SecondVotingPeriod")
	})
	t.Run("status is valid for default", func(t *testing.T) {
		gotOptionType := NormalizeProposalStatus("")
		assert.Equal(t, gotOptionType, "")
	})
}
