package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func TestSimAppExport(t *testing.T) {
	encodingConfig := MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := NewShentuApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 1, encodingConfig, EmptyAppOptions{})

	genesisState := ModuleBasics.DefaultGenesis(encodingConfig.Codec)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// make a new app object with the db so that initchain hasn't been called
	app2 := NewShentuApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 1, encodingConfig, EmptyAppOptions{})
	_, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
	_, err = app2.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators for zero height should not have an error")
}

func TestGetMaccPerms(t *testing.T) {
	dup := make(map[string][]string)
	for k, v := range maccPerms {
		dup[k] = v
	}
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}
