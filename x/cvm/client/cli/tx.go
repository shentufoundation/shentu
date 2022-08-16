package cli

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/txs/payload"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/cvm/compile"
	"github.com/shentufoundation/shentu/v2/x/cvm/types"
)

const (
	FlagValue    = "value"
	FlagArgs     = "args"
	FlagContract = "contract"
	FlagRaw      = "raw"
	FlagABI      = "abi"
	FlagEWASM    = "ewasm"
	FlagRuntime  = "runtime"
	FlagMetadata = "metadata"
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
			clientCtx, err := client.GetClientTxContext(cmd)
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

					return clientCtx.PrintProto(out)
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
	flags.AddTxFlagsToCmd(cmd)
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
		Use:   "deploy <filename>",
		Short: "Deploy CVM contract(s)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, file []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			from := clientCtx.GetFromAddress()
			accGetter := authtxb.AccountRetriever{}
			if err := accGetter.EnsureExists(clientCtx, from); err != nil {
				return err
			}

			msgs, err := appendDeployMsgs(cmd, file[0])
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msgs...)
		},
	}
	cmd.Flags().String(FlagABI, "", "name of ABI file (when deploying bytecode)")
	cmd.Flags().Uint64(FlagValue, 0, "value sent with transaction")
	cmd.Flags().String(FlagArgs, "", "constructor arguments")
	cmd.Flags().Bool(FlagEWASM, false, "compile solidity contract to EWASM")
	cmd.Flags().Bool(FlagRuntime, false, "runtime code")
	cmd.Flags().String(FlagMetadata, "", "the metadata files to be deployed along with the contract")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func appendDeployMsgs(cmd *cobra.Command, fileName string) ([]sdk.Msg, error) {
	var msgs []sdk.Msg
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return []sdk.Msg{}, err
	}

	argsRaw := viper.GetString(FlagArgs)
	arguments := strings.Split(argsRaw, ",")
	code, err := ioutil.ReadFile(fileName)
	codeStr := strings.Trim(string(code), "\n")
	if err != nil {
		return msgs, err
	}
	code, err = hex.DecodeString(codeStr)
	if err != nil {
		return msgs, err
	}
	codehash := crypto.Keccak256(code)

	value := viper.GetUint64(FlagValue)
	metadataFile := viper.GetString(FlagMetadata)
	var metas []*payload.ContractMeta
	if metadataFile != "" {
		metadataBytes, err := ioutil.ReadFile(metadataFile)
		if err != nil {
			return msgs, err
		}
		metadataString := string(metadataBytes)
		metas = append(metas, &payload.ContractMeta{
			CodeHash: codehash,
			Meta:     metadataString,
		})
	}

	abiFile, err := cmd.Flags().GetString(FlagABI)
	if err != nil {
		return msgs, err
	}
	var abiBytes []byte
	if abiFile != "" {
		abiBytes, err = ioutil.ReadFile(abiFile)
		if err != nil {
			return msgs, err
		}
	}

	logger := logging.NewNoopLogger()

	if len(argsRaw) > 0 {
		callArgsBytes, err := parseData("", abiBytes, arguments, logger)
		if err != nil {
			return msgs, err
		}
		code = append(code, callArgsBytes...)
	}
	isEWASM := viper.GetBool(FlagEWASM)
	isRuntime := viper.GetBool(FlagRuntime)
	msg := types.NewMsgDeploy(clientCtx.GetFromAddress().String(), value, code, string(abiBytes), metas, isEWASM, isRuntime)
	if err := msg.ValidateBasic(); err != nil {
		return msgs, err
	}
	msgs = append(msgs, &msg)
	return msgs, nil
}

func parseData(function string, abiSpec []byte, args []string, logger *logging.Logger) ([]byte, error) {
	params := make([]interface{}, len(args))

	if string(abiSpec) == compile.NoABI {
		panic("No ABI registered for this contract. Use --raw flag to submit raw bytecode.")
	}

	for i, arg := range args {
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
		params[i] = argi
	}

	data, _, err := abi.EncodeFunctionCall(string(abiSpec), function, logger, params...)
	return data, err
}
