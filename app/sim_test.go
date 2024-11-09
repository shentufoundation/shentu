package app

//
//import (
//	"encoding/json"
//	"fmt"
//	"math/rand"
//	"os"
//	"runtime/debug"
//	"strings"
//	"testing"
//
//	"cosmossdk.io/log"
//	dbm "github.com/cometbft/cometbft-db"
//	abci "github.com/cometbft/cometbft/abci/types"
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//	"github.com/stretchr/testify/require"
//
//	"cosmossdk.io/store"
//	storetypes "cosmossdk.io/store/types"
//
//	evidence "cosmossdk.io/x/evidence/types"
//
//	"github.com/cosmos/cosmos-sdk/baseapp"
//	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
//	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
//	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
//	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
//	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
//	mint "github.com/cosmos/cosmos-sdk/x/mint/types"
//	params "github.com/cosmos/cosmos-sdk/x/params/types"
//	"github.com/cosmos/cosmos-sdk/x/simulation"
//	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
//	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
//	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
//	capability "github.com/cosmos/ibc-go/modules/capability/types"
//	"github.com/cosmos/ibc-go/v8/testing/simapp"
//
//	cert "github.com/shentufoundation/shentu/v2/x/cert/types"
//	oracle "github.com/shentufoundation/shentu/v2/x/oracle/types"
//)
//
//// SimAppChainID hardcoded chainID for simulation
//const SimAppChainID = "simulation-app"
//
//type StoreKeysPrefixes struct {
//	A        storetypes.StoreKey
//	B        storetypes.StoreKey
//	Prefixes [][]byte
//}
//
//func init() {
//	simcli.GetSimulatorFlags()
//}
//
//// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
//// an IAVLStore for faster simulation speed.
//func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
//	bapp.SetFauxMerkleMode()
//}
//
//// interBlockCacheOpt returns a BaseApp option function that sets the persistent
//// inter-block write-through cache.
//func interBlockCacheOpt() func(*baseapp.BaseApp) {
//	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
//}
//
//func TestFullAppSimulation(t *testing.T) {
//	config := simcli.NewConfigFromFlags()
//	config.ChainID = SimAppChainID
//
//	db, dir, logger, skip, err := simtestutil.SetupSimulation(config, "leveldb-app-sim", "Simulation", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
//	if skip {
//		t.Skip("skipping application simulation")
//	}
//	require.NoError(t, err, "simulation setup failed")
//
//	defer func() {
//		db.Close()
//		require.NoError(t, os.RemoveAll(dir))
//	}()
//	app := NewShentuApp(logger, db, nil, true, map[int64]bool{},
//		dir, simcli.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
//	require.Equal(t, AppName, app.Name())
//
//	// run randomized simulation
//	_, simParams, simErr := simulation.SimulateFromSeed(
//		t,
//		os.Stdout,
//		app.BaseApp,
//		simapp.AppStateFn(app.Codec(), app.SimulationManager()),
//		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
//		simtestutil.SimulationOperations(app, app.Codec(), config),
//		app.ModuleAccountAddrs(),
//		config,
//		app.Codec(),
//	)
//
//	// export state and simParams before the simulation error is checked
//	err = simtestutil.CheckExportSimulation(app, config, simParams)
//	require.NoError(t, err)
//	require.NoError(t, simErr)
//
//	if config.Commit {
//		simtestutil.PrintStats(db)
//	}
//}
//
//func TestAppImportExport(t *testing.T) {
//	config := simcli.NewConfigFromFlags()
//	config.ChainID = SimAppChainID
//
//	db, dir, logger, skip, err := simtestutil.SetupSimulation(config, "leveldb-app-sim", "Simulation", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
//	if skip {
//		t.Skip("skipping application import/export simulation")
//	}
//	require.NoError(t, err, "simulation setup failed")
//
//	defer func() {
//		db.Close()
//		require.NoError(t, os.RemoveAll(dir))
//	}()
//
//	app := NewShentuApp(logger, db, nil, true, map[int64]bool{},
//		dir, simcli.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
//	require.Equal(t, "Shentu", app.Name())
//
//	// Run randomized simulation
//	_, simParams, simErr := simulation.SimulateFromSeed(
//		t,
//		os.Stdout,
//		app.BaseApp,
//		simapp.AppStateFn(app.Codec(), app.SimulationManager()),
//		simtypes.RandomAccounts,
//		simtestutil.SimulationOperations(app, app.Codec(), config),
//		app.ModuleAccountAddrs(),
//		config,
//		app.Codec(),
//	)
//
//	// export state and simParams before the simulation error is checked
//	err = simtestutil.CheckExportSimulation(app, config, simParams)
//	require.NoError(t, err)
//	require.NoError(t, simErr)
//
//	if config.Commit {
//		simtestutil.PrintStats(db)
//	}
//
//	fmt.Printf("exporting genesis...\n")
//	exported, err := app.ExportAppStateAndValidators(false, []string{}, []string{})
//	require.NoError(t, err)
//
//	fmt.Printf("importing genesis...\n")
//
//	newDB, newDir, _, _, err := simtestutil.SetupSimulation(config, "leveldb-app-sim-2", "Simulation-2", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
//	require.NoError(t, err, "simulation setup failed")
//
//	defer func() {
//		newDB.Close()
//		require.NoError(t, os.RemoveAll(newDir))
//	}()
//
//	newApp := NewShentuApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{},
//		newDir, simcli.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
//
//	var genesisState GenesisState
//	err = json.Unmarshal(exported.AppState, &genesisState)
//	require.NoError(t, err)
//
//	defer func() {
//		if r := recover(); r != nil {
//			err := fmt.Sprintf("%v", r)
//			if !strings.Contains(err, "validator set is empty after InitGenesis") {
//				panic(r)
//			}
//			logger.Info("Skipping simulation as all validators have been unbonded")
//			logger.Info("err", err, "stacktrace", string(debug.Stack()))
//		}
//	}()
//
//	header := tmproto.Header{Height: app.LastBlockHeight(), ChainID: SimAppChainID}
//	ctxA := app.NewContext(true, header)
//	ctxB := newApp.NewContext(true, header)
//	newApp.mm.InitGenesis(ctxB, app.Codec(), genesisState)
//	newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)
//
//	fmt.Printf("comparing stores...\n")
//
//	storeKeysPrefixes := []StoreKeysPrefixes{
//		{app.GetKey(auth.StoreKey), newApp.GetKey(auth.StoreKey), [][]byte{}},
//		{app.GetKey(staking.StoreKey), newApp.GetKey(staking.StoreKey), [][]byte{
//			staking.UnbondingQueueKey, staking.RedelegationQueueKey, staking.ValidatorQueueKey,
//			staking.HistoricalInfoKey,
//		}},
//		{app.GetKey(distr.StoreKey), newApp.GetKey(distr.StoreKey), [][]byte{}},
//		{app.GetKey(mint.StoreKey), newApp.GetKey(mint.StoreKey), [][]byte{}},
//		{app.GetKey(slashing.StoreKey), newApp.GetKey(slashing.StoreKey), [][]byte{}},
//		{app.GetKey(bank.StoreKey), newApp.GetKey(bank.StoreKey), [][]byte{bank.BalancesPrefix}},
//		{app.GetKey(gov.StoreKey), newApp.GetKey(gov.StoreKey), [][]byte{}},
//		{app.GetKey(cert.StoreKey), newApp.GetKey(cert.StoreKey), [][]byte{}},
//		{app.GetKey(oracle.StoreKey), newApp.GetKey(oracle.StoreKey), [][]byte{oracle.TaskStoreKeyPrefix, oracle.ClosingTaskStoreKeyPrefix, oracle.ClosingTaskStoreKeyTimedPrefix, oracle.ExpireTaskStoreKeyPrefix}},
//		{app.GetKey(evidence.StoreKey), newApp.GetKey(evidence.StoreKey), [][]byte{}},
//		{app.GetKey(capability.StoreKey), newApp.GetKey(capability.StoreKey), [][]byte{}},
//		{app.GetKey(params.StoreKey), newApp.GetKey(params.StoreKey), [][]byte{}},
//	}
//
//	for _, skp := range storeKeysPrefixes {
//		storeA := ctxA.KVStore(skp.A)
//		storeB := ctxB.KVStore(skp.B)
//
//		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
//		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")
//
//		fmt.Printf("compared %d different key/value pairs between %s and %s\n", len(failedKVAs), skp.A, skp.B)
//		require.Equal(t, len(failedKVAs), 0, simtestutil.GetSimulationLog(skp.A.Name(), app.SimulationManager().StoreDecoders, failedKVAs, failedKVBs))
//	}
//}
//
//func TestAppSimulationAfterImport(t *testing.T) {
//	config := simcli.NewConfigFromFlags()
//	config.ChainID = SimAppChainID
//
//	db, dir, logger, skip, err := simtestutil.SetupSimulation(config, "leveldb-app-sim", "Simulation", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
//	if skip {
//		t.Skip("skipping application import/export simulation")
//	}
//	require.NoError(t, err, "simulation setup failed")
//
//	app := NewShentuApp(logger, db, nil, true, map[int64]bool{},
//		dir, simcli.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
//	require.Equal(t, AppName, app.Name())
//
//	// run randomized simulation
//	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
//		t,
//		os.Stdout,
//		app.BaseApp,
//		simapp.AppStateFn(app.Codec(), app.SimulationManager()),
//		simtypes.RandomAccounts,
//		simtestutil.SimulationOperations(app, app.Codec(), config),
//		app.ModuleAccountAddrs(),
//		config,
//		app.Codec(),
//	)
//
//	// export state and simParams before the simulation error is checked
//	err = simtestutil.CheckExportSimulation(app, config, simParams)
//	require.NoError(t, err)
//	require.NoError(t, simErr)
//
//	if config.Commit {
//		simtestutil.PrintStats(db)
//	}
//
//	if stopEarly {
//		fmt.Println("can't export or import a zero-validator genesis, exiting test...")
//		return
//	}
//
//	fmt.Printf("exporting genesis...\n")
//
//	appState, err := app.ExportAppStateAndValidators(false, []string{}, []string{})
//	require.NoError(t, err)
//
//	fmt.Printf("importing genesis...\n")
//
//	newDB, newDir, _, _, err := simtestutil.SetupSimulation(config, "leveldb-app-sim-2", "Simulation-2", simcli.FlagVerboseValue, simcli.FlagEnabledValue) // nolint
//	require.NoError(t, err, "simulation setup failed")
//
//	defer func() {
//		newDB.Close()
//		require.NoError(t, os.RemoveAll(newDir))
//	}()
//
//	newApp := NewShentuApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{},
//		DefaultNodeHome, simcli.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
//	require.Equal(t, AppName, newApp.Name())
//
//	newApp.InitChain(abci.RequestInitChain{
//		ChainId:       SimAppChainID,
//		AppStateBytes: appState.AppState,
//	})
//
//	_, _, err = simulation.SimulateFromSeed(
//		t,
//		os.Stdout,
//		newApp.BaseApp,
//		simapp.AppStateFn(app.Codec(), app.SimulationManager()),
//		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
//		simtestutil.SimulationOperations(newApp, newApp.Codec(), config),
//		app.ModuleAccountAddrs(),
//		config,
//		app.Codec(),
//	)
//	require.NoError(t, err)
//}
//
//func TestAppStateDeterminism(t *testing.T) {
//	if !simcli.FlagEnabledValue {
//		t.Skip("skipping application simulation")
//	}
//
//	config := simcli.NewConfigFromFlags()
//	config.InitialBlockHeight = 1
//	config.ExportParamsPath = ""
//	config.OnOperation = false
//	config.AllInvariants = false
//	config.ChainID = SimAppChainID
//
//	numSeeds := 3
//	numTimesToRunPerSeed := 5
//	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)
//
//	for i := 0; i < numSeeds; i++ {
//		//nolint: gosec
//		config.Seed = rand.Int63()
//
//		for j := 0; j < numTimesToRunPerSeed; j++ {
//			logger := log.NewNopLogger()
//			if simcli.FlagVerboseValue {
//				logger = log.TestingLogger()
//			} else {
//				logger = log.NewNopLogger()
//			}
//
//			db := dbm.NewMemDB()
//			app := NewShentuApp(logger, db, nil, true, map[int64]bool{},
//				DefaultNodeHome, simcli.FlagPeriodValue, MakeEncodingConfig(), EmptyAppOptions{}, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
//
//			fmt.Printf(
//				"running non-determinism simulation; seed %d: attempt: %d/%d\n",
//				config.Seed, j+1, numTimesToRunPerSeed,
//			)
//
//			_, _, err := simulation.SimulateFromSeed(
//				t,
//				os.Stdout,
//				app.BaseApp,
//				simapp.AppStateFn(app.Codec(), app.SimulationManager()),
//				simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
//				simtestutil.SimulationOperations(app, app.Codec(), config),
//				app.ModuleAccountAddrs(),
//				config,
//				app.Codec(),
//			)
//			require.NoError(t, err)
//
//			if config.Commit {
//				simtestutil.PrintStats(db)
//			}
//
//			appHash := app.LastCommitID().Hash
//			appHashList[j] = appHash
//
//			if j != 0 {
//				require.Equal(
//					t, appHashList[0], appHashList[j],
//					"non-determinism in seed %d: attempt: %d/%d\n", config.Seed, j+1, numTimesToRunPerSeed,
//				)
//			}
//		}
//	}
//}
