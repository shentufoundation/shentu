package cert_test

import (
	"testing"

	"github.com/magiconair/properties/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func Test_CertifierStoreKey(t *testing.T) {
	t.Run("test ", func(t *testing.T) {
		acc := sdk.AccAddress([]byte{10})
		tmp := types.CertifierStoreKey(acc)
		assert.Equal(t, tmp, []byte{0, 10})
	})
}
