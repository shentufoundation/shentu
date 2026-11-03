package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/log"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	legacywasm "github.com/shentufoundation/shentu/v2/app/legacy/wasm"

	dbm "github.com/cosmos/cosmos-db"
)

func TestShentuAppExportAndBlockedAddrs(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := NewShentuAppWithCustomOptions(t, false, SetupOptions{
		Logger:  logger.With("instance", "first"),
		DB:      db,
		AppOpts: simtestutil.NewAppOptionsWithFlagHome(t.TempDir()),
	})

	// BlockedAddresses returns a map of addresses in app v1 and a map of modules name in app di.
	for acc := range app.ModuleAccountAddrs() {
		var addr sdk.AccAddress
		if modAddr, err := sdk.AccAddressFromBech32(acc); err == nil {
			addr = modAddr
		} else {
			addr = app.AccountKeeper.GetModuleAddress(acc)
		}

		require.True(
			t,
			app.BankKeeper.BlockedAddr(addr),
			fmt.Sprintf("ensure that blocked addresses are properly set in bank keeper: %s should be blocked", acc),
		)
	}

	// A passed proposal remains in the gov store after removal of x/wasm.
	ctx := app.BaseApp.NewContext(false)
	historicalMsg, err := codectypes.NewAnyWithValue(&legacywasm.MsgUpdateParams{
		Authority: app.GovKeeper.GetGovernanceAccount(ctx).GetAddress().String(),
		Params: legacywasm.Params{
			CodeUploadAccess:             legacywasm.AccessConfig{Permission: legacywasm.AccessTypeNobody},
			InstantiateDefaultPermission: legacywasm.AccessTypeEverybody,
		},
	})
	require.NoError(t, err)
	require.NoError(t, app.GovKeeper.Proposals.Set(ctx, 42, govv1.Proposal{
		Id: 42, Status: govv1.StatusPassed, Messages: []*codectypes.Any{historicalMsg},
	}))
	require.NoError(t, app.GovKeeper.ProposalID.Set(ctx, 43))

	// finalize block so we have CheckTx state set
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: 1,
	})
	require.NoError(t, err)

	_, err = app.Commit()
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewShentuApp(logger.With("instance", "second"), db, nil, true, simtestutil.NewAppOptionsWithFlagHome(t.TempDir()))
	exported, err := app2.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
	require.Contains(t, string(exported.AppState), "/cosmwasm.wasm.v1.MsgUpdateParams")

	exported, err = app2.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators for zero height should not have an error")
	require.Contains(t, string(exported.AppState), "/cosmwasm.wasm.v1.MsgUpdateParams")
}
