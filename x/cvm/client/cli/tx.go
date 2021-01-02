package cli

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/hyperledger/burrow/crypto"
	evm "github.com/hyperledger/burrow/deploy/compile"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cvm/compile"
	"github.com/certikfoundation/shentu/x/cvm/types"
)

const (
	FlagValue    = "value"
	FlagArgs     = "args"
	FlagContract = "contract"
	FlagRaw      = "raw"
	FlagABI      = "abi"
	FlagEWASM    = "ewasm"
	FlagRuntime  = "runtime"
)

var (
	errFileExt = errors.New("contract file extension must be .sol, .ds, .bc .bytecode or .wasm")
)

type abiEntry struct {
	Name string `json:"name"`
	Type string `json:"stateMutability"`
}

// GetTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	ctkTxCmd := &cobra.Command{
		Use:   "cvm",
		Short: "CVM transactions subcommands",
	}

	ctkTxCmd.AddCommand(
		GetCmdCall(),
		GetCmdDeploy(),
	)

	return ctkTxCmd
}

// GetCmdCall returns the CVM contract call transaction command.
func GetCmdCall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call <address> <function> [<params>...]",
		Short: "Call CVM contract",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			from := clientCtx.GetFromAddress()
			accGetter := authtxb.AccountRetriever{}
			if err := accGetter.EnsureExists(clientCtx, from); err != nil {
				return err
			}

			callee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var data []byte
			if viper.GetBool(FlagRaw) {
				if len(args) > 2 {
					return errors.New("cvm call with --raw flag should only have one argument (raw calldata)")
				}

				// support both hex and base64
				var err error
				data, err = hex.DecodeString(args[1])
				if err != nil {
					data, err = base64.StdEncoding.DecodeString(args[1])
					if err != nil {
						fmt.Println("Raw calldata could not be parsed. Use hex or base64. Don't put 0x at the beginning if you're using hex.")
						return err
					}
				}
			} else {
				var err error
				var abiSpec []byte
				abiSpec, data, err = parseCallCmd(clientCtx, args[0], callee, args[1], args[2:])
				if err != nil {
					return err
				}

				// Decode abiSpec to check if the called function's type is view or pure.
				// If it is, reroute to query.
				var abiEntries []abiEntry
				err = json.Unmarshal(abiSpec, &abiEntries)
				if err != nil {
					return err
				}
				for _, entry := range abiEntries {
					if entry.Name != args[1] {
						continue
					}
					if entry.Type != "view" && entry.Type != "pure" {
						break
					}
					fmt.Println(args[1] + " is a " + entry.Type + " function - Attempting to re-route to query")
					queryClient := types.NewQueryClient(clientCtx)
					req := &types.QueryViewRequest{
						Caller:       from.String(),
						Callee:       args[0],
						AbiSpec:      abiSpec,
						FunctionName: args[1],
						Data:         data,
					}
					out, err := queryClient.View(cmd.Context(), req)
					if err != nil {
						return err
					}

					return clientCtx.PrintOutput(out)
				}
			}
			value := viper.GetUint64(FlagValue)
			msg := types.NewMsgCall(from.String(), callee.String(), value, data)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, &msg)
		},
	}
	cmd.Flags().Bool(FlagRaw, false,
		"set this flag to submit raw calldata, otherwise it takes function name and parameters as args")
	cmd.Flags().Uint64(FlagValue, 0, "Value sent with transaction")
	return cmd
}

func parseCallCmd(cliCtx client.Context, calleeString string, calleeAddr sdk.AccAddress, function string, args []string) ([]byte, []byte, error) {
	accGetter := authtxb.AccountRetriever{}
	if err := accGetter.EnsureExists(cliCtx, calleeAddr); err != nil {
		return nil, nil, err
	}

	abiSpec, err := queryAbi(cliCtx, "cvm", calleeString)
	if err != nil {
		return nil, nil, err
	}
	logger := logging.NewNoopLogger()

	data, err := parseData(function, abiSpec, args, logger)
	if err != nil {
		return nil, nil, err
	}
	return abiSpec, data, nil
}

func queryAbi(cliCtx client.Context, queryRoute string, addr string) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/abi/%s", queryRoute, addr), nil)
	if err != nil {
		return nil, err
	}

	var out types.QueryResAbi
	cliCtx.LegacyAmino.MustUnmarshalJSON(res, &out)
	return out.Abi, nil
}

// GetCmdDeploy returns the CVM contract deploy transaction command.
func GetCmdDeploy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <filename> <flags>..",
		Short: "Deploy CVM contract(s)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, file []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			from := clientCtx.GetFromAddress()
			accGetter := authtxb.AccountRetriever{}
			if err := accGetter.EnsureExists(clientCtx, from); err != nil {
				return err
			}

			var msgs []sdk.Msg
			msgs, err = appendDeployMsgs(cmd, clientCtx, msgs, file[0])
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msgs...)
		},
	}
	cmd.Flags().String(FlagABI, "", "name of ABI file (when deploying bytecode)")
	cmd.Flags().Uint64(FlagValue, 0, "value sent with transaction")
	cmd.Flags().String(FlagArgs, "", "constructor arguments")
	cmd.Flags().String(FlagContract, "", "the name of the contract to be deployed")
	cmd.Flags().Bool(FlagEWASM, false, "compile solidity contract to EWASM")
	cmd.Flags().Bool(FlagRuntime, false, "runtime code")

	return cmd
}

func appendDeployMsgs(cmd *cobra.Command, cliCtx client.Context, msgs []sdk.Msg, fileName string) ([]sdk.Msg, error) {
	clientCtx := client.GetClientContextFromCmd(cmd)
	clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
	if err != nil {
		return []sdk.Msg{}, err
	}

	argumentsRaw := viper.GetString(FlagArgs)
	arguments := strings.Split(argumentsRaw, ",")
	deployContract := viper.GetString(FlagContract)

	var target string
	if len(deployContract) > 0 {
		target = strings.ToUpper(deployContract)
	} else {
		target = strings.ToUpper(filepath.Base(fileName))
	}

	resp, err := callEVM(cmd, fileName)
	if err != nil {
		return msgs, err
	}
	value := viper.GetUint64(FlagValue)

	fileNameMatch := false
	for _, object := range resp.Objects {
		code, err := hex.DecodeString(object.Contract.Code())
		if err != nil {
			return msgs, err
		}

		logger := logging.NewNoopLogger()
		metadata, err := object.Contract.GetMetadata(logger)
		if err != nil {
			return msgs, err
		}

		var metas []*payload.ContractMeta
		for codehash, metadata := range metadata {
			metas = append(metas, &payload.ContractMeta{
				CodeHash: codehash.Bytes(),
				Meta:     metadata,
			})
		}

		fileExtensionUpper := filepath.Ext(target)
		fileNameUpper := strings.TrimSuffix(target, fileExtensionUpper)
		objectNameUpper := strings.ToUpper(object.Objectname)
		if fileNameUpper == objectNameUpper || fileExtensionUpper == ".BYTECODE" {
			fileNameMatch = true
			if len(argumentsRaw) > 0 {
				callArgsBytes, err := parseData("", object.Contract.Abi, arguments, logger)
				if err != nil {
					return msgs, err
				}
				code = append(code, callArgsBytes...)
			}
			isEWASM := viper.GetBool(FlagEWASM)
			isRuntime := viper.GetBool(FlagRuntime)
			msg := types.NewMsgDeploy(clientCtx.GetFromAddress().String(), value, code, string(object.Contract.Abi), metas, isEWASM, isRuntime)
			if err := msg.ValidateBasic(); err != nil {
				return msgs, err
			}
			msgs = append(msgs, &msg)
		}
	}
	if !fileNameMatch {
		return msgs, errors.New("contract name does not match the file name")
	}
	return msgs, nil
}

func parseData(function string, abiSpec []byte, args []string, logger *logging.Logger) ([]byte, error) {
	var params []interface{}

	if string(abiSpec) == compile.NoABI {
		panic("No ABI registered for this contract. Use --raw flag to submit raw bytecode.")
	}

	for _, arg := range args {
		var argi interface{}
		argi = arg
		for _, prefix := range []string{common.Bech32MainPrefix, common.Bech32PrefixConsAddr, common.Bech32PrefixAccAddr} {
			if strings.HasPrefix(arg, prefix) && ((len(arg) - len(prefix)) == 39) {
				data, _ := sdk.GetFromBech32(arg, prefix)
				var err error
				argi, err = crypto.AddressFromBytes(data)
				if err != nil {
					return nil, err
				}
				break
			}
		}
		params = append(params, argi)
	}

	data, _, err := abi.EncodeFunctionCall(string(abiSpec), function, logger, params...)
	return data, err
}

func callEVM(cmd *cobra.Command, filename string) (*evm.Response, error) {
	logger := logging.NewNoopLogger()

	basename, workDir, err := compile.ResolveFilename(filename)
	if err != nil {
		return nil, err
	}

	basenameSplit := strings.Split(basename, ".")
	if len(basenameSplit) < 2 {
		return nil, errFileExt
	}

	var resp *evm.Response

	switch fileExt := basenameSplit[len(basenameSplit)-1]; fileExt {
	case "sol":
		if viper.GetBool(FlagEWASM) {
			resp, err = evm.WASM(basename, workDir, logger)
		} else {
			resp, err = evm.EVM(basename, false, workDir, nil, logger)
		}
	case "ds":
		resp, err = compile.DeepseaEVM(basename, workDir, logger)
	case "bc", "bytecode", "wasm":
		abiFile, err := cmd.Flags().GetString(FlagABI)
		if err != nil {
			return nil, err
		}
		resp, err = compile.BytecodeEVM(basename, workDir, abiFile, logger)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errFileExt
	}
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	if len(resp.Objects) < 1 {
		return nil, errors.New("compilation result must contain at least one object")
	}

	return resp, nil

}

// QueryTxCmd implements the default command for a tx query.
func QueryTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Query for a transaction by hash in a committed block",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			output, err := authcli.QueryTx(clientCtx, args[0])
			if err != nil {
				return err
			}

			if output.Empty() {
				return fmt.Errorf("no transaction found with hash %s", args[0])
			}

			return clientCtx.PrintOutput(output)
		},
	}

	cmd.Flags().StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))

	return cmd
}
