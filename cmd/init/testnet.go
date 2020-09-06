// Package init is used for initializing specialized chain state.
package init

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/app"
	"github.com/certikfoundation/shentu/common"
)

var (
	flagNodeDirPrefix   = "node-dir-prefix"
	flagNumValidators   = "v"
	flagOutputDir       = "output-dir"
	flagNodeDaemonHome  = "node-daemon-home"
	flagNodeCliHome     = "node-cli-home"
	flagServerIPAddress = "server-ip-address"
	flagPortIncrement   = "port-increment"
	flagDefaultPassword = "default-password"
	flagTestnetConfig   = "config"
)

const nodeDirPerm = 0755

var supernodes = []string{
	"CertiK",
	"Flint@Yale",
	"VeriGu@Columbia",
	"DeepSEA",
	"CertiKOS",
	"CVM",
	"DeepWallet",
	"NoOps",
	"Bitmain",
	"Fenbushi Capital",
	"InfStones",
	"SNZ",
	"DHVC",
	"吞吴",
	"Arrington Capital"}

// TestnetFilesCmd returns the CLI command handle that initializes all files for tendermint testnet and application.
func TestnetFilesCmd(ctx *server.Context, cdc *codec.Codec,
	mbm module.BasicManager, genAccIterator genutiltypes.GenesisAccountsIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a certikd testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	certikd testnet --v 4 --output-dir ./output --server-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config := ctx.Config
			return initTestnet(cmd, config, cdc, mbm, genAccIterator)
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with",
	)
	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet",
	)
	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)",
	)
	cmd.Flags().String(flagNodeDaemonHome, "certikd",
		"Home directory of the node's daemon configuration",
	)
	cmd.Flags().String(flagNodeCliHome, "certikcli",
		"Home directory of the node's cli configuration",
	)
	cmd.Flags().String(flagTestnetConfig, "", "Initialization config.")
	ip, err := server.ExternalIP()
	if err != nil {
		panic(err)
	}
	cmd.Flags().String(flagServerIPAddress, ip, "Server IP Address")
	cmd.Flags().Int(flagPortIncrement, 100, "")
	cmd.Flags().String(
		flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created",
	)
	cmd.Flags().String(
		server.FlagMinGasPrices, "",
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum.",
	)
	cmd.Flags().String(flagDefaultPassword, app.DefaultKeyPass, "")
	cmd.Flags().String(flags.FlagKeyringBackend, "test", "Select keyring's backend (os|file|test)")

	return cmd
}

// IPAddress : ec2-5-5-5-5.compute-1.amazonaws.com
type IPAddress = string

// NodeName : e.g. node0
type NodeName = string

// ServerName : e.g. EAST40
type ServerName = string

// TestnetConfig : describe entire testnet
type TestnetConfig struct {
	ChainID string                   `json:"chain_id"`
	Servers map[ServerName]IPAddress `json:"servers"`
	Nodes   map[NodeName]NodeConfig  `json:"nodes"`
}

// NodeConfig : describe a node in testnet
type NodeConfig struct {
	Server             ServerName `json:"server"`
	PortShift          int        `json:"port_shift"`
	InitialCTK         int64      `json:"initial_ctk"`
	Vesting            int64      `json:"vesting"`
	VestingTime        int64      `json:"vesting_time"`
	InitialDelegation  int64      `json:"initial_delegation"`
	StakingDescription string     `json:"staking_description"`
	SentryNodes        []NodeName `json:"sentry_nodes"`
	IsSeed             bool       `json:"is_seed"`
	Continuous         bool       `json:"continuous"`
}

func writeConfig(filename string, defaults map[string]string) error {
	v := viper.New()
	for key, value := range defaults {
		v.Set(key, value)
	}
	return v.WriteConfigAs(filename)
}

// GeneratedNode : Generated per-node data from the configuration
type GeneratedNode struct {
	valPubKey       crypto.PubKey
	nodeID          string
	nodeFullAddress string
	genFile         string
	account         exported.GenesisAccount
}

func (config TestnetConfig) validateBasic() error {
	for nodename, nodeConfig := range config.Nodes {
		if _, found := config.Servers[nodeConfig.Server]; !found {
			return errors.New("Unknown server " + nodeConfig.Server + " for node " + nodename)
		}
		for _, sentryNode := range nodeConfig.SentryNodes {
			if _, found := config.Nodes[sentryNode]; !found {
				return errors.New("Unknown sentry node " + sentryNode + " for node " + nodename)
			}
		}
	}
	return nil
}

func parseTestnetConfig(cfg *TestnetConfig) error {
	configPath := viper.GetString(flagTestnetConfig)
	if configPath != "" {
		data, err := ioutil.ReadFile(configPath) // #nosec
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return err
		}
	} else {
		// build config from commandline arguments
		cfg.Servers = make(map[ServerName]IPAddress)
		cfg.Nodes = make(map[NodeName]NodeConfig)
		numValidators := viper.GetInt(flagNumValidators)
		portInc := viper.GetInt(flagPortIncrement)
		serverAddressList := strings.Split(viper.GetString(flagServerIPAddress), ",")
		cfg.ChainID = viper.GetString(flags.FlagChainID)
		if cfg.ChainID == "" {
			cfg.ChainID = "chain-" + tmrand.NewRand().Str(6)
		}
		if len(serverAddressList) != 1 && len(serverAddressList) != numValidators {
			return errors.New("address list length does not match --v")
		}
		for i := 0; i < numValidators; i++ {
			nodename := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
			node := NodeConfig{
				Server:             nodename,
				PortShift:          portInc * i,
				InitialCTK:         int64(2000000 * (numValidators - i)),
				InitialDelegation:  int64(1000000 * (numValidators - i)),
				StakingDescription: nodename,
			}
			if len(serverAddressList) > 1 {
				cfg.Servers[nodename] = serverAddressList[i]
			} else {
				cfg.Servers[nodename] = serverAddressList[0]
			}
			if count := len(supernodes); i < count {
				node.StakingDescription = supernodes[i]
			} else {
				node.StakingDescription = fmt.Sprintf("%s-%d", supernodes[i%count], i/count)
			}
			cfg.Nodes[nodename] = node
		}
	}
	return cfg.validateBasic()
}

func initTestnet(cmd *cobra.Command, config *tmconfig.Config, cdc *codec.Codec,
	mbm module.BasicManager, genAccIterator genutiltypes.GenesisAccountsIterator) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())

	outDir := viper.GetString(flagOutputDir)
	defaultPassword := viper.GetString(flagDefaultPassword)
	keyringBackend := viper.GetString(flags.FlagKeyringBackend)

	certikConfig := srvconfig.DefaultConfig()
	certikConfig.MinGasPrices = viper.GetString(server.FlagMinGasPrices)

	// generate private keys, node IDs, and initial transactions
	cfg := TestnetConfig{}
	if err := parseTestnetConfig(&cfg); err != nil {
		return err
	}
	generatedNodes := make(map[string]GeneratedNode)
	for nodename, nodeConfig := range cfg.Nodes {
		nodeDirName := nodename
		fmt.Printf("Generating node %s\n", nodename)
		nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
		nodeCliHomeName := viper.GetString(flagNodeCliHome)
		nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
		clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)
		gentxsDir := filepath.Join(outDir, "gentxs")

		config.SetRoot(nodeDir)

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		if err := os.MkdirAll(clientDir, nodeDirPerm); err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		certikConfigFilePath := filepath.Join(nodeDir, "config/certikd.toml")
		srvconfig.WriteConfigFile(certikConfigFilePath, certikConfig)

		config.Moniker = nodeDirName

		ip := cfg.Servers[nodeConfig.Server]
		PortShift := nodeConfig.PortShift

		nodeID, valPubKey, initErr := genutil.InitializeNodeValidatorFiles(config)
		if initErr != nil {
			_ = os.RemoveAll(outDir)
			return initErr
		}

		nodeFullAddress := fmt.Sprintf("%s@%s:%d", nodeID, ip, 26656+PortShift)
		memo := nodeFullAddress
		genFile := config.GenesisFile()

		keyPass := defaultPassword
		if keyPass == "" {
			buf := bufio.NewReader(os.Stdin)
			prompt := fmt.Sprintf(
				"Password for account '%s' (default %s):", nodeDirName, app.DefaultKeyPass,
			)

			userInput, err := input.GetPassword(prompt, buf)
			keyPass = userInput
			if err != nil && keyPass != "" {
				// An error was returned that either failed to read the password from
				// STDIN or the given password is not empty but failed to meet minimum
				// length requirements.
				return err
			}
		}
		if keyPass == "" {
			keyPass = app.DefaultKeyPass
		}

		kb, err := keys.NewKeyring(strings.ToLower(app.AppName), keyringBackend, clientDir, inBuf)
		if err != nil {
			return err
		}
		addr, secret, err := server.GenerateSaveCoinKey(kb, nodeDirName, keyPass, true)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint); err != nil {
			return err
		}

		// save default client config
		if err := os.MkdirAll(filepath.Join(clientDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}
		if err := writeConfig(filepath.Join(clientDir, "config", "config.toml"), map[string]string{
			flags.FlagChainID:        cfg.ChainID,
			flags.FlagKeyringBackend: keyringBackend,
		}); err != nil {
			return err
		}

		accStakingTokens := sdk.TokensFromConsensusPower(nodeConfig.InitialCTK)
		accCoins := sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, accStakingTokens)}
		account := auth.NewBaseAccount(addr, accCoins, nil, 0, 0)
		if nodeConfig.Vesting > 0 && nodeConfig.Continuous {
			accVestingTokens := sdk.TokensFromConsensusPower(nodeConfig.Vesting)
			accVesting := sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, accVestingTokens)}
			vestingAccount, err := vesting.NewBaseVestingAccount(
				account,
				accVesting,
				time.Now().Unix()+nodeConfig.VestingTime,
			)
			if err != nil {
				return err
			}
			vestingStart := time.Now().Unix() + 10
			delayedVestingAccount := vesting.NewContinuousVestingAccountRaw(vestingAccount, vestingStart)
			generatedNodes[nodename] = GeneratedNode{
				valPubKey:       valPubKey,
				nodeID:          nodeID,
				nodeFullAddress: nodeFullAddress,
				genFile:         genFile,
				account:         delayedVestingAccount,
			}
		} else if nodeConfig.Vesting > 0 {
			accVestingTokens := sdk.TokensFromConsensusPower(nodeConfig.Vesting)
			accVesting := sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, accVestingTokens)}
			vestingAccount, err := vesting.NewBaseVestingAccount(
				account,
				accVesting,
				time.Now().Unix()+nodeConfig.VestingTime,
			)
			if err != nil {
				return err
			}
			delayedVestingAccount := vesting.NewDelayedVestingAccountRaw(vestingAccount)
			generatedNodes[nodename] = GeneratedNode{
				valPubKey:       valPubKey,
				nodeID:          nodeID,
				nodeFullAddress: nodeFullAddress,
				genFile:         genFile,
				account:         delayedVestingAccount,
			}
		} else {
			generatedNodes[nodename] = GeneratedNode{
				valPubKey:       valPubKey,
				nodeID:          nodeID,
				nodeFullAddress: nodeFullAddress,
				genFile:         genFile,
				account:         account,
			}
		}

		if nodeConfig.InitialDelegation == 0 {
			continue
		}
		valTokens := sdk.TokensFromConsensusPower(nodeConfig.InitialDelegation)
		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKey,
			sdk.NewCoin(common.MicroCTKDenom, valTokens),
			staking.NewDescription(nodeConfig.StakingDescription, "", "", "", ""),
			stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			sdk.OneInt(),
		)
		tx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
		txBldr := auth.NewTxBuilderFromCLI(inBuf).WithChainID(cfg.ChainID).WithMemo(memo).WithKeybase(kb)

		signedTx, err := txBldr.SignStdTx(nodeDirName, keyPass, tx, false)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		txBytes, err := cdc.MarshalJSON(signedTx)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		// gather gentxs folder
		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBytes); err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}
	}
	if err := initGenFiles(cdc, mbm, cfg.ChainID, generatedNodes); err != nil {
		return err
	}

	if err := collectGenFiles(
		cdc, config, cfg.ChainID, cfg.Nodes, generatedNodes,
		outDir, viper.GetString(flagNodeDaemonHome), genAccIterator,
	); err != nil {
		return err
	}

	fmt.Printf("Successfully initialized %d node directories\n", len(cfg.Nodes))
	return nil
}

func initGenFiles(
	cdc *codec.Codec, mbm module.BasicManager, chainID string, generatedNodes map[string]GeneratedNode,
) error {
	appGenState := mbm.DefaultGenesis()
	accounts := make([]exported.GenesisAccount, 0)
	for nodename := range generatedNodes {
		accounts = append(accounts, generatedNodes[nodename].account)
	}

	// set the accounts in the genesis state
	var authGenState auth.GenesisState
	cdc.MustUnmarshalJSON(appGenState[auth.ModuleName], &authGenState)
	authGenState.Accounts = accounts
	appGenState[auth.ModuleName] = cdc.MustMarshalJSON(authGenState)

	appGenStateJSON, err := codec.MarshalJSONIndent(cdc, appGenState)
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for nodename := range generatedNodes {
		if err := genDoc.SaveAs(generatedNodes[nodename].genFile); err != nil {
			return err
		}
	}

	return nil
}

func collectGenFiles(
	cdc *codec.Codec, config *tmconfig.Config, chainID string,
	nodeConfigs map[NodeName]NodeConfig,
	generatedNodes map[string]GeneratedNode,
	outputDir, nodeDaemonHome string,
	genAccIterator genutiltypes.GenesisAccountsIterator,
) error {
	var appState json.RawMessage
	genTime := tmtime.Now()

	publicPeers, seeds, sentryToValidatorMap := make([]string, 0), make([]string, 0), make(map[string]string)
	for nodename, nodeConfig := range nodeConfigs {
		nodeFullAddress := generatedNodes[nodename].nodeFullAddress
		if len(nodeConfig.SentryNodes) == 0 && !nodeConfig.IsSeed {
			publicPeers = append(publicPeers, nodeFullAddress)
		}
		if nodeConfig.IsSeed {
			seeds = append(seeds, nodeFullAddress)
		}
		for _, sentryNode := range nodeConfig.SentryNodes {
			sentryToValidatorMap[sentryNode] = nodeFullAddress
		}
	}
	config.P2P.Seeds = strings.Join(seeds, ",")

	for nodename, nodeConfig := range nodeConfigs {
		gentxsDir := filepath.Join(outputDir, "gentxs")
		config.Moniker = nodename
		config.ProfListenAddress = fmt.Sprintf("localhost:%d", 6060+nodeConfig.PortShift)
		config.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", 26656+nodeConfig.PortShift)
		config.P2P.AddrBookStrict = false
		config.P2P.AllowDuplicateIP = true
		config.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", 26657+nodeConfig.PortShift)
		config.ProxyApp = fmt.Sprintf("tcp://127.0.0.1:%d", 26658+nodeConfig.PortShift)
		config.Instrumentation.PrometheusListenAddr = fmt.Sprintf(":%d", 26660+nodeConfig.PortShift)
		config.DBBackend = "goleveldb"
		config.P2P.SeedMode = nodeConfig.IsSeed
		config.SetRoot(filepath.Join(outputDir, nodename, nodeDaemonHome))

		generatedNode := generatedNodes[nodename]

		initCfg := genutil.NewInitConfig(chainID, gentxsDir, nodename, generatedNode.nodeID, generatedNode.valPubKey)

		genDoc, err := types.GenesisDocFromFile(config.GenesisFile())
		if err != nil {
			return err
		}
		if validatorID, ok := sentryToValidatorMap[nodename]; ok {
			config.P2P.PrivatePeerIDs = validatorID
		} else {
			config.P2P.PrivatePeerIDs = ""
		}
		persistentPeers := make([]string, 0)
		if len(nodeConfig.SentryNodes) > 0 {
			for _, sentryNode := range nodeConfig.SentryNodes {
				persistentPeers = append(persistentPeers, generatedNodes[sentryNode].nodeFullAddress)
			}
			config.P2P.PexReactor = false
		} else {
			persistentPeers = publicPeers
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(cdc, config, initCfg, *genDoc, genAccIterator)
		config.P2P.PersistentPeers = strings.Join(persistentPeers, ",")
		tmconfig.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := config.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(name, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	if err := tmos.EnsureDir(writePath, 0700); err != nil {
		return err
	}

	if err := tmos.WriteFile(file, contents, 0600); err != nil {
		return err
	}

	return nil
}
