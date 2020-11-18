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
