package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
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

	certiksimapp "github.com/certikfoundation/shentu/simapp"
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

	app := certiksimapp.NewSimApp(logger, db, nil, true, map[int64]bool{},
		DefaultNodeHome, simapp.FlagPeriodValue, certiksimapp.MakeTestEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.AppCodec(), app.SimulationManager()),
		RandomAccounts, simapp.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(), config, app.AppCodec(),
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

	//invCheckPeriod := simapp.FlagPeriodValue
	var invCheckPeriod uint = 1
	app := certiksimapp.NewSimApp(logger, db, nil, true, map[int64]bool{},
		DefaultNodeHome, invCheckPeriod, certiksimapp.MakeTestEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.AppCodec(), app.SimulationManager()),
		RandomAccounts, simapp.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(), config, app.AppCodec(),
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

	newApp := certiksimapp.NewSimApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{},
		DefaultNodeHome, invCheckPeriod, certiksimapp.MakeTestEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)

	require.Equal(t, AppName, newApp.Name())

	var genesisState simapp.GenesisState
	err = json.Unmarshal(appState.AppState, &genesisState)
	require.NoError(t, err)

	ctxA := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	ctxB := newApp.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	newApp.ModuleManager().InitGenesis(ctxB, app.AppCodec(), genesisState)
	newApp.StoreConsensusParams(ctxB, appState.ConsensusParams)

	fmt.Printf("comparing stores...\n")

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.GetKey(auth.StoreKey), newApp.GetKey(auth.StoreKey), [][]byte{}},
		{app.GetKey(staking.StoreKey), newApp.GetKey(staking.StoreKey), [][]byte{
			staking.UnbondingQueueKey, staking.RedelegationQueueKey, staking.ValidatorQueueKey,
			staking.HistoricalInfoKey,
		}},
		{app.GetKey(distr.StoreKey), newApp.GetKey(distr.StoreKey), [][]byte{}},
		{app.GetKey(mint.StoreKey), newApp.GetKey(mint.StoreKey), [][]byte{}},
		{app.GetKey(slashing.StoreKey), newApp.GetKey(slashing.StoreKey), [][]byte{}},
		{app.GetKey(bank.StoreKey), newApp.GetKey(bank.StoreKey), [][]byte{bank.BalancesPrefix}},
		{app.GetKey(params.StoreKey), newApp.GetKey(params.StoreKey), [][]byte{}},
		{app.GetKey(upgrade.StoreKey), newApp.GetKey(upgrade.StoreKey), [][]byte{}},
		{app.GetKey(gov.StoreKey), newApp.GetKey(gov.StoreKey), [][]byte{}},
		{app.GetKey(cert.StoreKey), newApp.GetKey(cert.StoreKey), [][]byte{}},
		{app.GetKey(cvm.StoreKey), newApp.GetKey(cvm.StoreKey), [][]byte{}},
		{app.GetKey(oracle.StoreKey), newApp.GetKey(oracle.StoreKey), [][]byte{oracle.TaskStoreKeyPrefix, oracle.ClosingTaskStoreKeyPrefix}},
		{app.GetKey(shield.StoreKey), newApp.GetKey(shield.StoreKey), [][]byte{shield.WithdrawQueueKey, shield.PurchaseQueueKey, shield.BlockServiceFeesKey}},
		{app.GetKey(evidence.StoreKey), newApp.GetKey(evidence.StoreKey), [][]byte{}},
		{app.GetKey(capability.StoreKey), newApp.GetKey(capability.StoreKey), [][]byte{}},
		{app.GetKey(ibchost.StoreKey), newApp.GetKey(ibchost.StoreKey), [][]byte{}},
		{app.GetKey(ibctransfer.StoreKey), newApp.GetKey(ibctransfer.StoreKey), [][]byte{}},
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

	app := certiksimapp.NewSimApp(logger, db, nil, true, map[int64]bool{},
		DefaultNodeHome, simapp.FlagPeriodValue, certiksimapp.MakeTestEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.AppCodec(), app.SimulationManager()),
		RandomAccounts, simapp.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(), config, app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}

	fmt.Printf("exporting genesis...\n")

	appState, err := app.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	_, newDB, newDir, _, _, err := simapp.SetupSimulation("leveldb-app-sim-2", "Simulation-2") // nolint
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := certiksimapp.NewSimApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{},
		DefaultNodeHome, simapp.FlagPeriodValue, certiksimapp.MakeTestEncodingConfig(), EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, AppName, newApp.Name())

	newApp.InitChain(abci.RequestInitChain{
		AppStateBytes: appState.AppState,
	})

	_, _, err = simulation.SimulateFromSeed(
		t, os.Stdout, newApp.BaseApp, simapp.AppStateFn(app.AppCodec(), app.SimulationManager()),
		RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simapp.SimulationOperations(newApp, newApp.AppCodec(), config),
		app.ModuleAccountAddrs(), config, app.AppCodec(),
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
		app := certiksimapp.NewSimApp(logger, db, nil, true, map[int64]bool{},
			DefaultNodeHome, simapp.FlagPeriodValue, certiksimapp.MakeTestEncodingConfig(), EmptyAppOptions{}, interBlockCacheOpt())

		fmt.Printf(
			"running non-determinism simulation; seed %d: attempt: %d/%d\n",
			config.Seed, j+1, numTimesToRunPerSeed,
		)

		_, _, err := simulation.SimulateFromSeed(
			t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.AppCodec(), app.SimulationManager()),
			RandomAccounts, // Replace with own random account function if using keys other than secp256k1
			simapp.SimulationOperations(app, app.AppCodec(), config),
			app.ModuleAccountAddrs(), config, app.AppCodec(),
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

// RandomAccounts generates n random accounts
func RandomAccounts(r *rand.Rand, n int) []simtypes.Account {
	accs := make([]simtypes.Account, n)

	for i := 0; i < n; i++ {
		// don't need that much entropy for simulation
		privkeySeed := make([]byte, 15)
		r.Read(privkeySeed)

		accs[i].PrivKey = secp256k1.GenPrivKeyFromSecret(privkeySeed)
		accs[i].PubKey = accs[i].PrivKey.PubKey()
		accs[i].Address = sdk.AccAddress(accs[i].PubKey.Address())

		accs[i].ConsKey = ed25519.GenPrivKeyFromSecret(privkeySeed)
	}

	return accs
}
