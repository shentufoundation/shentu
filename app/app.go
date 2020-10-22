// Package app provides the assets information for server module.
package app

import (
	"fmt"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	cosmosGov "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/certikfoundation/shentu/x/auth"
	"github.com/certikfoundation/shentu/x/auth/vesting"
	"github.com/certikfoundation/shentu/x/bank"
	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/crisis"
	"github.com/certikfoundation/shentu/x/cvm"
	distr "github.com/certikfoundation/shentu/x/distribution"
	"github.com/certikfoundation/shentu/x/gov"
	"github.com/certikfoundation/shentu/x/mint"
	"github.com/certikfoundation/shentu/x/oracle"
	"github.com/certikfoundation/shentu/x/shield"
	"github.com/certikfoundation/shentu/x/slashing"
	"github.com/certikfoundation/shentu/x/staking"
	"github.com/certikfoundation/shentu/x/upgrade"
)

const (
	// AppName specifies the global application name.
	AppName = "CertiK"

	// DefaultKeyPass for certikd node daemon.
	DefaultKeyPass = "12345678"

	keysReserved  = 100
	tkeysReserved = 10
)

var (
	// DefaultCLIHome specifies where the node client data is stored.
	DefaultCLIHome = os.ExpandEnv("$HOME/.certikcli")

	// DefaultNodeHome specifies where the node daemon data is stored.
	DefaultNodeHome = os.ExpandEnv("$HOME/.certikd")

	// ModuleBasics is in charge of setting up basic, non-dependant module
	// elements, such as codec registration and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			distr.ProposalHandler,
			upgrade.ProposalHandler,
			cert.ProposalHandler,
			paramsclient.ProposalHandler,
			// Disabled for phase I.
			// shield.ProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		cvm.NewAppModuleBasic(),
		cert.NewAppModuleBasic(),
		oracle.NewAppModuleBasic(),
		shield.NewAppModuleBasic(),
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		oracle.ModuleName:         {supply.Burner},
		shield.ModuleName:         {supply.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distr.ModuleName:  true,
		oracle.ModuleName: true,
		shield.ModuleName: true,
	}
)

// CertiKApp is the main CertiK Chain application type.
type CertiKApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	crisisKeeper   crisis.Keeper
	supplyKeeper   supply.Keeper
	paramsKeeper   params.Keeper
	upgradeKeeper  upgrade.Keeper
	govKeeper      gov.Keeper
	certKeeper     cert.Keeper
	cvmKeeper      cvm.Keeper
	authKeeper     auth.Keeper
	oracleKeeper   oracle.Keeper
	shieldKeeper   shield.Keeper

	// module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewCertiKApp returns a reference to an initialized CertiKApp.
func NewCertiKApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *CertiKApp {
	// define top-level codec that will be shared between modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol.
	bApp := bam.NewBaseApp(AppName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	ks := []string{
		bam.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		distr.StoreKey,
		mint.StoreKey,
		slashing.StoreKey,
		params.StoreKey,
		upgrade.StoreKey,
		gov.StoreKey,
		cert.StoreKey,
		cvm.StoreKey,
		oracle.StoreKey,
		shield.StoreKey,
	}

	for i := 0; i < keysReserved; i++ {
		ks = append(ks, fmt.Sprintf("reserved%d", i))
	}

	keys := sdk.NewKVStoreKeys(ks...)

	tks := []string{
		staking.TStoreKey,
		params.TStoreKey,
	}

	for i := 0; i < tkeysReserved; i++ {
		tks = append(tks, fmt.Sprintf("reservedT%d", i))
	}

	tkeys := sdk.NewTransientStoreKeys(tks...)

	// initialize application with its store keys
	var app = &CertiKApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
	}
	// initialize params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	oracleSubspace := app.paramsKeeper.Subspace(oracle.DefaultParamSpace)
	cvmSubspace := app.paramsKeeper.Subspace(cvm.DefaultParamSpace)
	shieldSubspace := app.paramsKeeper.Subspace(shield.DefaultParamSpace)

	// initialize keepers
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		authSubspace,
		auth.ProtoBaseAccount,
	)
	app.bankKeeper = bank.NewKeeper(
		app.accountKeeper,
		&app.cvmKeeper,
		bankSubspace,
		app.BlacklistedAccAddrs(),
	)
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		&app.supplyKeeper,
		stakingSubspace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		keys[distr.StoreKey],
		distrSubspace,
		&stakingKeeper,
		&app.supplyKeeper,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	app.cvmKeeper = cvm.NewKeeper(
		app.cdc,
		keys[cvm.StoreKey],
		app.accountKeeper,
		app.distrKeeper,
		&app.certKeeper,
		cvmSubspace,
	)
	app.oracleKeeper = oracle.NewKeeper(
		app.cdc,
		keys[oracle.StoreKey],
		app.accountKeeper,
		app.distrKeeper,
		&app.stakingKeeper,
		app.supplyKeeper,
		oracleSubspace,
	)
	app.mintKeeper = mint.NewKeeper(
		app.cdc,
		keys[mint.StoreKey],
		mintSubspace,
		&stakingKeeper,
		app.supplyKeeper,
		app.distrKeeper,
		auth.FeeCollectorName,
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		slashingSubspace,
	)
	app.certKeeper = cert.NewKeeper(
		app.cdc,
		keys[cert.StoreKey],
		app.slashingKeeper,
		stakingKeeper,
	)
	app.authKeeper = auth.NewKeeper(
		app.certKeeper,
	)
	app.crisisKeeper = crisis.NewKeeper(
		crisisSubspace,
		invCheckPeriod,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)
	app.upgradeKeeper = upgrade.NewKeeper(
		skipUpgradeHeights,
		keys[upgrade.StoreKey],
		app.cdc,
	)
	app.shieldKeeper = shield.NewKeeper(
		app.cdc,
		keys[shield.StoreKey],
		&stakingKeeper,
		&app.govKeeper,
		app.supplyKeeper,
		shieldSubspace,
	)
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference so that it will contain these hooks.
	app.stakingKeeper.Keeper = *stakingKeeper.Keeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks(),
			app.shieldKeeper.Hooks(),
		),
	)
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		keys[gov.StoreKey],
		govSubspace,
		app.supplyKeeper,
		app.stakingKeeper,
		app.certKeeper,
		app.shieldKeeper,
		app.upgradeKeeper,
		cosmosGov.NewRouter().
			AddRoute(gov.RouterKey, gov.ProposalHandler).
			AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
			AddRoute(cert.RouterKey, cert.NewCertifierUpdateProposalHandler(app.certKeeper)).
			AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper.Keeper)).
			AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
			AddRoute(shield.RouterKey, shield.NewShieldClaimProposalHandler(app.shieldKeeper)),
	)

	// NOTE: Any module instantiated in the module manager that is
	// later modified must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper, app.certKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper.Keeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper.Keeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper, app.certKeeper),
		mint.NewAppModule(app.mintKeeper),
		upgrade.NewAppModule(app.upgradeKeeper.Keeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		cvm.NewAppModule(app.cvmKeeper),
		cert.NewAppModule(app.certKeeper, app.accountKeeper),
		oracle.NewAppModule(app.oracleKeeper),
		shield.NewAppModule(app.shieldKeeper, app.accountKeeper, app.stakingKeeper, app.supplyKeeper),
	)

	// NOTE: During BeginBlocker, slashing comes after distr so that
	// there is nothing left over in the validator fee pool, so as to
	// keep the CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(upgrade.ModuleName, mint.ModuleName, distr.ModuleName, slashing.ModuleName,
		supply.ModuleName, oracle.ModuleName, cvm.ModuleName, shield.ModuleName)

	// NOTE: Shield endblocker comes before staking because it queries
	// unbonding delegations that staking endblocker deletes.
	app.mm.SetOrderEndBlockers(crisis.ModuleName, cvm.ModuleName, shield.ModuleName, staking.ModuleName, gov.ModuleName, oracle.ModuleName)

	// NOTE: genutil moodule must occur after staking so that pools
	// are properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		auth.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		supply.ModuleName,
		cvm.ModuleName,
		shield.ModuleName,
		crisis.ModuleName,
		cert.ModuleName,
		genutil.ModuleName,
		oracle.ModuleName,
	)

	app.mm.SetOrderExportGenesis(
		auth.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		supply.ModuleName,
		cvm.ModuleName,
		crisis.ModuleName,
		cert.ModuleName,
		genutil.ModuleName,
		oracle.ModuleName,
		shield.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.accountKeeper, app.certKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper.Keeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper.Keeper),
		params.NewAppModule(),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper, app.certKeeper),
		mint.NewAppModule(app.mintKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		cvm.NewAppModule(app.cvmKeeper),
		cert.NewAppModule(app.certKeeper, app.accountKeeper),
		oracle.NewAppModule(app.oracleKeeper),
		shield.NewAppModule(app.shieldKeeper, app.accountKeeper, app.stakingKeeper, app.supplyKeeper),
	)

	app.sm.RegisterStoreDecoders()

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(app.keys[bam.MainStoreKey]); err != nil {
			tmos.Exit(err.Error())
		}
	}
	return app
}

// BeginBlocker processes application updates at the beginning of each block.
func (app *CertiKApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker processes application updates at the end of each block.
func (app *CertiKApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer defines application update at chain initialization
func (app *CertiKApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	return cdc
}

// LoadHeight loads a particular height
func (app *CertiKApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *CertiKApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *CertiKApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[supply.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}

// Codec returns app.cdc.
func (app *CertiKApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager returns app.sm.
func (app *CertiKApp) SimulationManager() *module.SimulationManager {
	return app.sm
}
