package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	capability "github.com/cosmos/cosmos-sdk/x/capability/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidence "github.com/cosmos/cosmos-sdk/x/evidence/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfer "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	mint "github.com/cosmos/cosmos-sdk/x/mint/types"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	cert "github.com/certikfoundation/shentu/x/cert/types"
	cvm "github.com/certikfoundation/shentu/x/cvm/types"
	oracle "github.com/certikfoundation/shentu/x/oracle/types"
	shield "github.com/certikfoundation/shentu/x/shield/types"
	//"github.com/certikfoundation/shentu/x/staking"
)

type StoreKeysPrefixes struct {
	A        sdk.StoreKey
	B        sdk.StoreKey
	Prefixes [][]byte
}

func init() {
	simapp.GetSimulatorFlags()
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func TestFullAppSimulation(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewCertiKApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.SimulationManager()),
		simtypes.RandomAccounts, simapp.SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config, app.Codec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}
}

func TestAppImportExport(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application import/export simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewCertiKApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.SimulationManager()),
		simtypes.RandomAccounts, simapp.SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config, app.Codec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}

	fmt.Printf("exporting genesis...\n")

	appState, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	_, newDB, newDir, _, _, err := simapp.SetupSimulation("leveldb-app-sim-2", "Simulation-2") // nolint
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := NewCertiKApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)

	require.Equal(t, AppName, newApp.Name())

	var genesisState simapp.GenesisState
	err = app.Codec().UnmarshalJSON(appState.AppState, &genesisState)
	require.NoError(t, err)

	ctxA := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	ctxB := newApp.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, genesisState)

	fmt.Printf("comparing stores...\n")

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.keys[auth.StoreKey], newApp.keys[auth.StoreKey], [][]byte{}},
		{app.keys[staking.StoreKey], newApp.keys[staking.StoreKey], [][]byte{
			staking.UnbondingQueueKey, staking.RedelegationQueueKey, staking.ValidatorQueueKey,
			staking.HistoricalInfoKey,
		}},
		{app.keys[distr.StoreKey], newApp.keys[distr.StoreKey], [][]byte{}},
		{app.keys[mint.StoreKey], newApp.keys[mint.StoreKey], [][]byte{}},
		{app.keys[slashing.StoreKey], newApp.keys[slashing.StoreKey], [][]byte{}},
		{app.keys[params.StoreKey], newApp.keys[params.StoreKey], [][]byte{}},
		{app.keys[upgrade.StoreKey], newApp.keys[upgrade.StoreKey], [][]byte{}},
		{app.keys[gov.StoreKey], newApp.keys[gov.StoreKey], [][]byte{}},
		{app.keys[cert.StoreKey], newApp.keys[cert.StoreKey], [][]byte{}},
		{app.keys[cvm.StoreKey], newApp.keys[cvm.StoreKey], [][]byte{}},
		{app.keys[oracle.StoreKey], newApp.keys[oracle.StoreKey], [][]byte{oracle.TaskStoreKeyPrefix, oracle.ClosingTaskStoreKeyPrefix}},
		{app.keys[shield.StoreKey], newApp.keys[shield.StoreKey], [][]byte{shield.WithdrawQueueKey, shield.PurchaseQueueKey, shield.BlockServiceFeesKey}},
		{app.keys[evidence.StoreKey], newApp.keys[evidence.StoreKey], [][]byte{}},
		{app.keys[capability.StoreKey], newApp.keys[capability.StoreKey], [][]byte{}},
		{app.keys[ibchost.StoreKey], newApp.keys[ibchost.StoreKey], [][]byte{}},
		{app.keys[ibctransfer.StoreKey], newApp.keys[ibctransfer.StoreKey], [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")
		if len(failedKVAs) != 0 {
			fmt.Printf("found %d non-equal key/value pairs between %s and %s\n", len(failedKVAs), skp.A.Name(), skp.B.Name())
		}
		require.Equal(t, len(failedKVAs), 0, simapp.GetSimulationLog(skp.A.Name(),
			app.SimulationManager().StoreDecoders, failedKVAs, failedKVBs))
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation after import")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewCertiKApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.SimulationManager()),
		simtypes.RandomAccounts, simapp.SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config, app.Codec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}

	fmt.Printf("exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	_, newDB, newDir, _, _, err := simapp.SetupSimulation("leveldb-app-sim-2", "Simulation-2") // nolint
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := NewCertiKApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, newApp.Name())

	newApp.InitChain(abci.RequestInitChain{
		AppStateBytes: appState,
	})

	_, _, err = simulation.SimulateFromSeed(
		t, os.Stdout, newApp.BaseApp, simapp.AppStateFn(app.Codec(), app.SimulationManager()),
		simapp.SimulationOperations(newApp, newApp.Codec(), config),
		newApp.ModuleAccountAddrs(), config,
	)
	require.NoError(t, err)
}

func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = helpers.SimAppChainID

	numTimesToRunPerSeed := 2
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for j := 0; j < numTimesToRunPerSeed; j++ {
		logger := log.NewNopLogger()
		db := dbm.NewMemDB()
		app := NewCertiKApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, interBlockCacheOpt())

		fmt.Printf(
			"running non-determinism simulation; seed %d: attempt: %d/%d\n",
			config.Seed, j+1, numTimesToRunPerSeed,
		)

		_, _, err := simulation.SimulateFromSeed(
			t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.SimulationManager()),
			simapp.SimulationOperations(app, app.Codec(), config),
			app.ModuleAccountAddrs(), config,
		)
		require.NoError(t, err)

		appHash := app.LastCommitID().Hash
		appHashList[j] = appHash

		if j != 0 {
			require.Equal(
				t, appHashList[0], appHashList[j],
				"non-determinism in seed %d: attempt: %d/%d\n", config.Seed, j+1, numTimesToRunPerSeed,
			)
		}
	}
}
