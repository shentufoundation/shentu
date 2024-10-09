package app

//
//import (
//	"os"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	dbm "github.com/cometbft/cometbft-db"
//	"cosmossdk.io/log"
//)
//
//func TestSimAppExport(t *testing.T) {
//	encCfg := MakeEncodingConfig()
//	db := dbm.NewMemDB()
//	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
//	app := NewSimappWithCustomOptions(t, false, SetupOptions{
//		Logger:             logger,
//		DB:                 db,
//		InvCheckPeriod:     0,
//		EncConfig:          encCfg,
//		HomePath:           DefaultNodeHome,
//		SkipUpgradeHeights: map[int64]bool{},
//		AppOpts:            EmptyAppOptions{},
//	})
//
//	for acc := range maccPerms {
//		require.True(
//			t,
//			app.BankKeeper.BlockedAddr(app.AccountKeeper.GetModuleAddress(acc)),
//			"ensure that blocked addresses are properly set in bank keeper",
//		)
//	}
//
//	app.Commit()
//
//	logger2 := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
//	// make a new app object with the db so that initchain hasn't been called
//	app2 := NewShentuApp(logger2, db, nil, true, map[int64]bool{}, DefaultNodeHome, 1, encCfg, EmptyAppOptions{})
//	_, err := app2.ExportAppStateAndValidators(false, []string{})
//	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
//	//_, err = app2.ExportAppStateAndValidators(true, []string{})
//	//require.NoError(t, err, "ExportAppStateAndValidators for zero height should not have an error")
//}
//
//func TestGetMaccPerms(t *testing.T) {
//	dup := make(map[string][]string)
//	for k, v := range maccPerms {
//		dup[k] = v
//	}
//	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
//}
