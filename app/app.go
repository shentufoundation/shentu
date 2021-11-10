// Package app provides the assets information for server module.
package app

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	sdkauthkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	sdkauthz "github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"
	sdkbanktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	sdkdistr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	sdkfeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrant "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	sdkgovtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkminttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/ibc-go/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/modules/core"
	ibcclient "github.com/cosmos/ibc-go/modules/core/02-client"
	ibcconnectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"
	ibcmock "github.com/cosmos/ibc-go/testing/mock"

	appparams "github.com/certikfoundation/shentu/v2/app/params"
	"github.com/certikfoundation/shentu/v2/x/auth"
	authkeeper "github.com/certikfoundation/shentu/v2/x/auth/keeper"
	"github.com/certikfoundation/shentu/v2/x/bank"
	bankkeeper "github.com/certikfoundation/shentu/v2/x/bank/keeper"
	"github.com/certikfoundation/shentu/v2/x/cert"
	certclient "github.com/certikfoundation/shentu/v2/x/cert/client"
	certkeeper "github.com/certikfoundation/shentu/v2/x/cert/keeper"
	certtypes "github.com/certikfoundation/shentu/v2/x/cert/types"
	"github.com/certikfoundation/shentu/v2/x/cvm"
	cvmkeeper "github.com/certikfoundation/shentu/v2/x/cvm/keeper"
	cvmtypes "github.com/certikfoundation/shentu/v2/x/cvm/types"
	distr "github.com/certikfoundation/shentu/v2/x/distribution"
	"github.com/certikfoundation/shentu/v2/x/gov"
	govkeeper "github.com/certikfoundation/shentu/v2/x/gov/keeper"
	govtypes "github.com/certikfoundation/shentu/v2/x/gov/types"
	"github.com/certikfoundation/shentu/v2/x/mint"
	mintkeeper "github.com/certikfoundation/shentu/v2/x/mint/keeper"
	"github.com/certikfoundation/shentu/v2/x/oracle"
	oraclekeeper "github.com/certikfoundation/shentu/v2/x/oracle/keeper"
	oracletypes "github.com/certikfoundation/shentu/v2/x/oracle/types"
	"github.com/certikfoundation/shentu/v2/x/shield"
	shieldclient "github.com/certikfoundation/shentu/v2/x/shield/client"
	shieldkeeper "github.com/certikfoundation/shentu/v2/x/shield/keeper"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/slashing"
	"github.com/certikfoundation/shentu/v2/x/staking"
	stakingkeeper "github.com/certikfoundation/shentu/v2/x/staking/keeper"

	// unnamed import of statik for swagger UI support
	_ "github.com/certikfoundation/shentu/v2/docs/statik"
)

const (
	// AppName specifies the global application name.
	AppName = "Shentu"
	upgradeName = "Shentu-upgrade"

	// DefaultKeyPass for certik node daemon.
	DefaultKeyPass = "12345678"
)

var (
	// DefaultNodeHome specifies where the node daemon data is stored.
	DefaultNodeHome = os.ExpandEnv("$HOME/.certik")

	// ModuleBasics is in charge of setting up basic, non-dependant module
	// elements, such as codec registration and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		authz.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		feegrant.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler,
			distrclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			upgradeclient.CancelProposalHandler,
			certclient.ProposalHandler,
			shieldclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		cvm.NewAppModuleBasic(),
		cert.NewAppModuleBasic(),
		oracle.NewAppModuleBasic(),
		shield.NewAppModuleBasic(),
		evidence.AppModuleBasic{},
		ibc.AppModuleBasic{},
		transfer.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		sdkminttypes.ModuleName:        {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		sdkgovtypes.ModuleName:         {authtypes.Burner},
		oracletypes.ModuleName:         {authtypes.Burner},
		shieldtypes.ModuleName:         {authtypes.Burner},
		cvmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}
)

// ShentuApp is the main Shentu Chain application type.
type ShentuApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	accountKeeper    sdkauthkeeper.AccountKeeper
	authzKeeper      authzkeeper.Keeper
	bankKeeper       bankkeeper.Keeper
	stakingKeeper    stakingkeeper.Keeper
	slashingKeeper   slashingkeeper.Keeper
	mintKeeper       mintkeeper.Keeper
	distrKeeper      distrkeeper.Keeper
	feegrantKeeper   feegrantkeeper.Keeper
	paramsKeeper     paramskeeper.Keeper
	upgradeKeeper    upgradekeeper.Keeper
	govKeeper        govkeeper.Keeper
	certKeeper       certkeeper.Keeper
	authKeeper       authkeeper.Keeper
	evidenceKeeper   evidencekeeper.Keeper
	ibcKeeper        *ibckeeper.Keeper
	transferKeeper   ibctransferkeeper.Keeper
	capabilityKeeper *capabilitykeeper.Keeper
	cvmKeeper        cvmkeeper.Keeper
	oracleKeeper     oraclekeeper.Keeper
	shieldKeeper     shieldkeeper.Keeper

	// make scoped keepers public for test purposes
	scopedIBCKeeper      capabilitykeeper.ScopedKeeper
	scopedTransferKeeper capabilitykeeper.ScopedKeeper
	scopedIBCMockKeeper  capabilitykeeper.ScopedKeeper

	// module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
	configurator module.Configurator
}

// NewShentuApp returns a reference to an initialized ShentuApp.
func NewShentuApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool, homePath string,
	invCheckPeriod uint, encodingConfig appparams.EncodingConfig, appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) *ShentuApp {
	// define top-level codec that will be shared between modules
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	// BaseApp handles interactions with Tendermint through the ABCI protocol.
	bApp := baseapp.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	ks := []string{
		authtypes.StoreKey,
		authzkeeper.StoreKey,
		sdkbanktypes.StoreKey,
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		sdkfeegrant.StoreKey,
		sdkminttypes.StoreKey,
		slashingtypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		sdkgovtypes.StoreKey,
		certtypes.StoreKey,
		cvmtypes.StoreKey,
		oracletypes.StoreKey,
		shieldtypes.StoreKey,
		evidencetypes.StoreKey,
		ibchost.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
	}

	keys := sdk.NewKVStoreKeys(ks...)

	tks := []string{
		paramstypes.TStoreKey,
	}
	tkeys := sdk.NewTransientStoreKeys(tks...)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// initialize application with its store keys
	var app = &ShentuApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}
	// initialize params keeper and subspaces
	app.paramsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.capabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	scopedIBCMockKeeper := app.capabilityKeeper.ScopeToModule(ibcmock.ModuleName)

	// initialize keepers
	app.accountKeeper = sdkauthkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)
	app.bankKeeper = bankkeeper.NewKeeper(
		appCodec,
		keys[sdkbanktypes.StoreKey],
		app.accountKeeper,
		&app.cvmKeeper,
		app.GetSubspace(sdkbanktypes.ModuleName),
		app.ModuleAccountAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)
	app.distrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.accountKeeper, app.bankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.feegrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[sdkfeegrant.StoreKey],
		app.accountKeeper,
	)
	app.cvmKeeper = cvmkeeper.NewKeeper(
		appCodec,
		keys[cvmtypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		app.distrKeeper,
		&app.certKeeper,
		&app.stakingKeeper,
		app.GetSubspace(cvmtypes.ModuleName),
	)
	app.oracleKeeper = oraclekeeper.NewKeeper(
		appCodec,
		keys[oracletypes.StoreKey],
		app.accountKeeper,
		app.distrKeeper,
		&app.stakingKeeper,
		app.bankKeeper,
		app.GetSubspace(oracletypes.ModuleName),
	)
	app.slashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.certKeeper = certkeeper.NewKeeper(
		appCodec,
		keys[certtypes.StoreKey],
		app.slashingKeeper,
		stakingKeeper,
	)
	app.authKeeper = authkeeper.NewKeeper(
		app.accountKeeper,
		app.certKeeper,
	)
	app.authzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
	)
	app.upgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
	)
	app.shieldKeeper = shieldkeeper.NewKeeper(
		appCodec,
		keys[shieldtypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		&stakingKeeper,
		&app.govKeeper,
		app.GetSubspace(shieldtypes.ModuleName),
	)
	app.mintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[sdkminttypes.StoreKey], app.GetSubspace(sdkminttypes.ModuleName), &stakingKeeper,
		app.accountKeeper, app.bankKeeper, app.distrKeeper, app.shieldKeeper, authtypes.FeeCollectorName,
	)
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference so that it will contain these hooks.
	app.stakingKeeper.Keeper = *stakingKeeper.Keeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks(),
			app.shieldKeeper.Hooks(),
		),
	)

	// Create IBC Keeper
	app.ibcKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.stakingKeeper, app.upgradeKeeper, scopedIBCKeeper,
	)

	govRouter := sdkgovtypes.NewRouter()
	govRouter.AddRoute(sdkgovtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distrtypes.RouterKey, sdkdistr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.ibcKeeper.ClientKeeper)).
		AddRoute(shieldtypes.RouterKey, shield.NewShieldClaimProposalHandler(app.shieldKeeper)).
		AddRoute(certtypes.RouterKey, cert.NewCertifierUpdateProposalHandler(app.certKeeper))

	app.govKeeper = govkeeper.NewKeeper(
		appCodec,
		keys[sdkgovtypes.StoreKey],
		app.GetSubspace(sdkgovtypes.ModuleName),
		app.bankKeeper,
		app.stakingKeeper,
		app.certKeeper,
		app.shieldKeeper,
		app.accountKeeper,
		govRouter,
	)

	// Create Transfer Keepers
	app.transferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.ibcKeeper.ChannelKeeper, &app.ibcKeeper.PortKeeper,
		app.accountKeeper, app.bankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.transferKeeper)

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	mockModule := ibcmock.NewAppModule(scopedIBCMockKeeper, &app.ibcKeeper.PortKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	ibcRouter.AddRoute(ibcmock.ModuleName, mockModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.stakingKeeper.Keeper, app.slashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.evidenceKeeper = *evidenceKeeper

	// Add empty upgrade handler to bump to 0.44.0
	cfg := module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.setUpgradeHandler(cfg)

	/****  Module Options ****/

	// NOTE: Any module instantiated in the module manager that is
	// later modified must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(appCodec, app.authKeeper, app.accountKeeper, app.bankKeeper, app.certKeeper, authsims.RandomGenesisAccounts),
		authz.NewAppModule(appCodec, app.authzKeeper, app.accountKeeper, app.bankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(appCodec, *app.capabilityKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper.Keeper),
		feegrant.NewAppModule(appCodec, app.accountKeeper, app.bankKeeper, app.feegrantKeeper, app.interfaceRegistry),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper.Keeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		mint.NewAppModule(appCodec, app.mintKeeper, app.accountKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		gov.NewAppModule(appCodec, app.govKeeper, app.accountKeeper, app.bankKeeper),
		cvm.NewAppModule(app.cvmKeeper, app.bankKeeper),
		cert.NewAppModule(app.certKeeper, app.accountKeeper, app.bankKeeper),
		oracle.NewAppModule(app.oracleKeeper, app.bankKeeper),
		shield.NewAppModule(app.shieldKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		transferModule,
	)

	// NOTE: During BeginBlocker, slashing comes after distr so that
	// there is nothing left over in the validator fee pool, so as to
	// keep the CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(upgradetypes.ModuleName, capabilitytypes.ModuleName, sdkminttypes.ModuleName, distrtypes.ModuleName,
		slashingtypes.ModuleName, evidencetypes.ModuleName, oracletypes.ModuleName, cvmtypes.ModuleName, stakingtypes.ModuleName,
		shieldtypes.ModuleName, ibchost.ModuleName)

	// NOTE: Shield endblocker comes before staking because it queries
	// unbonding delegations that staking endblocker deletes.
	app.mm.SetOrderEndBlockers(cvmtypes.ModuleName, shieldtypes.ModuleName, stakingtypes.ModuleName, sdkgovtypes.ModuleName,
		oracletypes.ModuleName)

	// NOTE: genutil moodule must occur after staking so that pools
	// are properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		sdkbanktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		sdkgovtypes.ModuleName,
		sdkminttypes.ModuleName,
		cvmtypes.ModuleName,
		shieldtypes.ModuleName,
		certtypes.ModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		oracletypes.ModuleName,
		sdkauthz.ModuleName,
		ibctransfertypes.ModuleName,
		sdkfeegrant.ModuleName,
	)

	app.mm.SetOrderExportGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		sdkbanktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		sdkgovtypes.ModuleName,
		sdkminttypes.ModuleName,
		cvmtypes.ModuleName,
		certtypes.ModuleName,
		genutiltypes.ModuleName,
		oracletypes.ModuleName,
		shieldtypes.ModuleName,
		ibchost.ModuleName,
		sdkauthz.ModuleName,
		ibctransfertypes.ModuleName,
		sdkfeegrant.ModuleName,
		evidencetypes.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.authKeeper, app.accountKeeper, app.bankKeeper, app.certKeeper, authsims.RandomGenesisAccounts),
		authz.NewAppModule(appCodec, app.authzKeeper, app.accountKeeper, app.bankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(appCodec, *app.capabilityKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper.Keeper),
		feegrant.NewAppModule(appCodec, app.accountKeeper, app.bankKeeper, app.feegrantKeeper, app.interfaceRegistry),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper.Keeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		mint.NewAppModule(appCodec, app.mintKeeper, app.accountKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		gov.NewAppModule(appCodec, app.govKeeper, app.accountKeeper, app.bankKeeper),
		cvm.NewAppModule(app.cvmKeeper, app.bankKeeper),
		cert.NewAppModule(app.certKeeper, app.accountKeeper, app.bankKeeper),
		oracle.NewAppModule(app.oracleKeeper, app.bankKeeper),
		shield.NewAppModule(app.shieldKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		transferModule,
	)

	app.sm.RegisterStoreDecoders()

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	// The AnteHandler handles signature verification and transaction pre-processing
	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.accountKeeper,
			BankKeeper:      app.bankKeeper,
			FeegrantKeeper:  app.feegrantKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		tmos.Exit(err.Error())
	}
	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	app.upgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, _ upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
			app.ibcKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

			fromVM := make(map[string]uint64)
			for moduleName := range app.mm.Modules {
				fromVM[moduleName] = 1
			}
			// override versions for _new_ modules as to not skip InitGenesis
			fromVM[sdkauthz.ModuleName] = 0
			fromVM[sdkfeegrant.ModuleName] = 0

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}
	app.scopedIBCKeeper = scopedIBCKeeper
	app.scopedTransferKeeper = scopedTransferKeeper

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	app.scopedIBCMockKeeper = scopedIBCMockKeeper

	return app
}

// BeginBlocker processes application updates at the beginning of each block.
func (app *ShentuApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker processes application updates at the end of each block.
func (app *ShentuApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer defines application update at chain initialization
func (app *ShentuApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// MakeCodecs constructs the *std.Marshaler and *codec.LegacyAmino instances used by
// app. It is useful for tests and clients who do not want to construct the
// full app
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

// LoadHeight loads a particular height
func (app *ShentuApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *ShentuApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ShentuApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ShentuApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// Codec returns app.legacyAmino.
func (app *ShentuApp) Codec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns the app's InterfaceRegistry
func (app *ShentuApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// SimulationManager returns app.sm.
func (app *ShentuApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *ShentuApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *ShentuApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *ShentuApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(sdkbanktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(sdkminttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(sdkgovtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(oracletypes.ModuleName).WithKeyTable(oracletypes.ParamKeyTable())
	paramsKeeper.Subspace(cvmtypes.ModuleName).WithKeyTable(cvmtypes.ParamKeyTable())
	paramsKeeper.Subspace(shieldtypes.ModuleName).WithKeyTable(shieldtypes.ParamKeyTable())

	return paramsKeeper
}
