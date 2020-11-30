package simapp

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
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	cosmosBank "github.com/cosmos/cosmos-sdk/x/bank"
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
			shield.ProposalHandler,
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

// Verify app interface at compile time
var _ simapp.App = (*SimApp)(nil)

// SimApp is the simulated CertiK Chain application type.
type SimApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	AccountKeeper  auth.AccountKeeper
	BankKeeper     bank.Keeper
	StakingKeeper  staking.Keeper
	SlashingKeeper slashing.Keeper
	MintKeeper     mint.Keeper
	DistrKeeper    distr.Keeper
	CrisisKeeper   crisis.Keeper
	SupplyKeeper   supply.Keeper
	ParamsKeeper   params.Keeper
	UpgradeKeeper  upgrade.Keeper
	GovKeeper      gov.Keeper
	CertKeeper     cert.Keeper
	CvmKeeper      cvm.Keeper
	AuthKeeper     auth.Keeper
	OracleKeeper   oracle.Keeper
	ShieldKeeper   shield.Keeper

	// module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *SimApp {
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
	var app = &SimApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
	}
	// initialize params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	authSubspace := app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.ParamsKeeper.Subspace(crisis.DefaultParamspace)
	oracleSubspace := app.ParamsKeeper.Subspace(oracle.DefaultParamSpace)
	cvmSubspace := app.ParamsKeeper.Subspace(cvm.DefaultParamSpace)
	shieldSubspace := app.ParamsKeeper.Subspace(shield.DefaultParamSpace)

	// initialize keepers
	app.AccountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		authSubspace,
		auth.ProtoBaseAccount,
	)
	app.BankKeeper = bank.NewKeeper(
		app.AccountKeeper,
		&app.CvmKeeper,
		bankSubspace,
		app.BlacklistedAccAddrs(),
	)
	app.SupplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		maccPerms,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		&app.SupplyKeeper,
		stakingSubspace,
	)
	app.DistrKeeper = distr.NewKeeper(
		app.cdc,
		keys[distr.StoreKey],
		distrSubspace,
		&stakingKeeper,
		&app.SupplyKeeper,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	app.CvmKeeper = cvm.NewKeeper(
		app.cdc,
		keys[cvm.StoreKey],
		app.AccountKeeper,
		app.DistrKeeper,
		&app.CertKeeper,
		cvmSubspace,
	)
	app.OracleKeeper = oracle.NewKeeper(
		app.cdc,
		keys[oracle.StoreKey],
		app.AccountKeeper,
		app.DistrKeeper,
		&app.StakingKeeper,
		app.SupplyKeeper,
		oracleSubspace,
	)
	app.MintKeeper = mint.NewKeeper(
		app.cdc,
		keys[mint.StoreKey],
		mintSubspace,
		&stakingKeeper,
		app.SupplyKeeper,
		app.DistrKeeper,
		&app.ShieldKeeper,
		auth.FeeCollectorName,
	)
	app.SlashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		slashingSubspace,
	)
	app.CertKeeper = cert.NewKeeper(
		app.cdc,
		keys[cert.StoreKey],
		app.SlashingKeeper,
		stakingKeeper,
	)
	app.AuthKeeper = auth.NewKeeper(
		app.CertKeeper,
	)
	app.CrisisKeeper = crisis.NewKeeper(
		crisisSubspace,
		invCheckPeriod,
		app.SupplyKeeper,
		auth.FeeCollectorName,
	)
	app.UpgradeKeeper = upgrade.NewKeeper(
		skipUpgradeHeights,
		keys[upgrade.StoreKey],
		app.cdc,
	)
	app.ShieldKeeper = shield.NewKeeper(
		app.cdc,
		keys[shield.StoreKey],
		app.AccountKeeper,
		&stakingKeeper,
		&app.GovKeeper,
		app.SupplyKeeper,
		shieldSubspace,
	)
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference so that it will contain these hooks.
	app.StakingKeeper.Keeper = *stakingKeeper.Keeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
			app.ShieldKeeper.Hooks(),
		),
	)
	app.GovKeeper = gov.NewKeeper(
		app.cdc,
		keys[gov.StoreKey],
		govSubspace,
		app.SupplyKeeper,
		app.StakingKeeper,
		app.CertKeeper,
		app.ShieldKeeper,
		app.UpgradeKeeper,
		cosmosGov.NewRouter().
			AddRoute(gov.RouterKey, gov.ProposalHandler).
			AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
			AddRoute(cert.RouterKey, cert.NewCertifierUpdateProposalHandler(app.CertKeeper)).
			AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper.Keeper)).
			AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
			AddRoute(shield.RouterKey, shield.NewShieldClaimProposalHandler(app.ShieldKeeper)),
	)

	// NOTE: Any module instantiated in the module manager that is
	// later modified must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper, app.CertKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		crisis.NewAppModule(&app.CrisisKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		distr.NewAppModule(app.DistrKeeper, app.AccountKeeper, app.SupplyKeeper, app.StakingKeeper.Keeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.StakingKeeper.Keeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper, app.CertKeeper),
		mint.NewAppModule(app.MintKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper.Keeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.SupplyKeeper),
		cvm.NewAppModule(app.CvmKeeper),
		cert.NewAppModule(app.CertKeeper, app.AccountKeeper),
		oracle.NewAppModule(app.OracleKeeper),
		shield.NewAppModule(app.ShieldKeeper, app.AccountKeeper, app.StakingKeeper, app.SupplyKeeper),
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

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.AccountKeeper, app.CertKeeper),
		cosmosBank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		distr.NewAppModule(app.DistrKeeper, app.AccountKeeper, app.SupplyKeeper, app.StakingKeeper.Keeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.StakingKeeper.Keeper),
		params.NewAppModule(),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper, app.CertKeeper),
		mint.NewAppModule(app.MintKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.SupplyKeeper),
		cvm.NewAppModule(app.CvmKeeper),
		cert.NewAppModule(app.CertKeeper, app.AccountKeeper),
		oracle.NewAppModule(app.OracleKeeper),
		shield.NewAppModule(app.ShieldKeeper, app.AccountKeeper, app.StakingKeeper, app.SupplyKeeper),
	)

	app.sm.RegisterStoreDecoders()

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(app.keys[bam.MainStoreKey]); err != nil {
			tmos.Exit(err.Error())
		}
	}
	return app
}

// BeginBlocker processes application updates at the beginning of each block.
func (app *SimApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker processes application updates at the end of each block.
func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer defines application update at chain initialization
func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
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
func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *SimApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[supply.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}

// Codec returns app.cdc.
func (app *SimApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager returns app.sm.
func (app *SimApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetKey returns the corresponding store key.
func (app *SimApp) GetKey(storeKey string) types.StoreKey {
	return app.keys[storeKey]
}

func (app *SimApp) GetAppModuleGenesis(ctx sdk.Context, moduleName string) module.AppModule {
	return app.mm.Modules[moduleName]
}
