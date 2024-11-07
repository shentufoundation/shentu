package e2e

import (
	"fmt"
	"os"

	"cosmossdk.io/log"
	evidencetypes "cosmossdk.io/x/evidence/types"
	feegrant "cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsproptypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	shentu "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/app/params"
	bountytypes "github.com/shentufoundation/shentu/v2/x/bounty/types"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	oracletypes "github.com/shentufoundation/shentu/v2/x/oracle/types"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var (
	encodingConfig params.EncodingConfig
	cdc            codec.Codec
)

func init() {
	encodingConfig = params.MakeEncodingConfig()

	stakingtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	banktypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	authvesting.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	stakingtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	evidencetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	govv1types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	govv1beta1types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	paramsproptypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	upgradetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	distribtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	txtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	feegrant.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	bountytypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	certtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	oracletypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	cdc = encodingConfig.Codec
}

type chain struct {
	dataDir         string
	id              string
	validators      []*validator
	accounts        []*account
	genesisAccounts []*account
	certifier       *account
}

func newChain() (*chain, error) {
	tmpDir, err := os.MkdirTemp("", "shentu-e2e-testnet-")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + tmrand.NewRand().Str(6),
		dataDir: tmpDir,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int) error {
	app := shentu.NewShentuApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		simtestutil.NewAppOptionsWithFlagHome(shentu.DefaultNodeHome),
	)

	defer func() {
		if err := app.Close(); err != nil {
			panic(err)
		}
	}()

	genesisState := app.BasicModuleManager.DefaultGenesis(encodingConfig.Codec)

	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(genesisState); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAndInitValidatorsWithMnemonics(count int, mnemonics []string) error {
	app := shentu.NewShentuApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		simtestutil.NewAppOptionsWithFlagHome(shentu.DefaultNodeHome),
	)

	defer func() {
		if err := app.Close(); err != nil {
			panic(err)
		}
	}()

	genesisState := app.BasicModuleManager.DefaultGenesis(encodingConfig.Codec)

	for i := 0; i < count; i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(genesisState); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKeyFromMnemonic("val", mnemonics[i]); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createValidator(index int) *validator {
	return &validator{
		chain:   c,
		index:   index,
		moniker: fmt.Sprintf("%s-shentu-%d", c.id, index),
	}
}
