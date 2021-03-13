package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNormalizeVoteOption(t *testing.T) {
	t.Run("option is valid for yes", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("yes")
		assert.Equal(t, gotOptionType, "VOTE_OPTION_YES")
	})
	t.Run("option is valid for abstain", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("Abstain")
		assert.Equal(t, gotOptionType, "VOTE_OPTION_ABSTAIN")
	})
	t.Run("option is valid for no", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("No")
		assert.Equal(t, gotOptionType, "VOTE_OPTION_NO")
	})

	t.Run("option is valid for noWithVeto", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("no_with_veto")
		assert.Equal(t, gotOptionType, "VOTE_OPTION_NO_WITH_VETO")
	})
	t.Run("option is valid for default", func(t *testing.T) {
		gotOptionType := NormalizeVoteOption("")
		assert.Equal(t, gotOptionType, "")
	})
}
