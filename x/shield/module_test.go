package shield_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := shentuapp.Setup(t, false)

	ctx := app.BaseApp.NewContext(false)
	acc := app.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	require.NotNil(t, acc)
	acc1 := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc1)
}
