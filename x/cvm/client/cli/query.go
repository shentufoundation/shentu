// Package cli defines the CLI services for the cvm module.
package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"

	"github.com/hyperledger/burrow/execution/evm/abi"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cvm/client/utils"
	"github.com/certikfoundation/shentu/x/cvm/types"
)

const (
	FlagCaller = "caller"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group cvm queries under a subcommand
	cvmQueryCmd := &cobra.Command{
		Use:   "cvm",
		Short: "Querying commands for the CVM module",
	}

	cvmQueryCmd.AddCommand(flags.GetCommands(
		GetCmdCode(queryRoute, cdc),
		GetCmdStorage(queryRoute, cdc),
		GetCmdAbi(queryRoute, cdc),
		GetCmdMeta(queryRoute, cdc),
		GetCmdView(queryRoute, cdc),
		GetCmdAddressTranslate(queryRoute, cdc),
	)...)

	return cvmQueryCmd
}

// GetCmdView returns the CVM contract view transaction command.
func GetCmdView(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <address> <function> [<params>...]",
		Short: "View CVM contract",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Caller is an optional flag. If not set it becomes the zero address.
			var callerString string
			callerString, err := cmd.Flags().GetString(FlagCaller)
			if err != nil {
				return err
			}

			if callerString != "" {
				_, err := sdk.AccAddressFromBech32(callerString)
				if err != nil {
					return err
				}
			}

			callee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			abiSpec, data, err := parseCallCmd(cliCtx, args[0], callee, args[1], args[2:])
			if err != nil {
				return err
			}

			queryPath := fmt.Sprintf("custom/%s/view/%s/%s", queryRoute, callerString, callee)
			return queryContractAndPrint(cliCtx, cdc, queryPath, args[1], abiSpec, data)
		},
	}
	cmd.Flags().String(FlagCaller, "", "optional caller parameter to run the view function with")

	return cmd
}

// Query CVM contract code based on ABI spec and print function output.
func queryContractAndPrint(cliCtx context.CLIContext, cdc *codec.Codec, queryPath, fname string, abiSpec, data []byte) error {
	res, _, err := cliCtx.QueryWithData(queryPath, data)
	if err != nil {
		return fmt.Errorf("querying CVM contract code: %v", err)
	}
	var out types.QueryResView
	cdc.MustUnmarshalJSON(res, &out)
	ret, err := abi.DecodeFunctionReturn(string(abiSpec), fname, out.Ret)
	if err != nil {
		return fmt.Errorf("decoding function return: %v", err)
	}
	err = cliCtx.PrintOutput(ret)
	if err != nil {
		return fmt.Errorf("printing output: %v", err)
	}
	return nil
}

// GetCmdCode returns the CVM code query command.
func GetCmdCode(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "code <address>",
		Short: "Get CVM contract code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addr := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/code/%s", queryRoute, addr), nil)
			if err != nil {
				fmt.Printf("could not get CVM contract code\n" + err.Error() + "\n")
				return nil
			}

			var out types.QueryResCode
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdStorage returns the CVM storage query command.
func GetCmdStorage(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "storage <address> <key>",
		Short: "Get CVM storage data",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addr := args[0]
			key := args[1]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/storage/%s/%s", queryRoute, addr, key), nil)
			if err != nil {
				fmt.Printf("could not get CVM storage\n" + err.Error() + "\n")
				return nil
			}

			var out types.QueryResStorage
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdAbi returns the CVM code ABI query command.
func GetCmdAbi(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "abi <address>",
		Short: "Get CVM contract code ABI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addr := args[0]

			bytes, err := queryAbi(cliCtx, queryRoute, addr)
			if err != nil {
				fmt.Printf("could not get CVM contract code ABI\n" + err.Error() + "\n")
				return nil
			}
			if bytes == nil {
				fmt.Printf("cannot find CVM contract code ABI\n")
				return nil
			}
			fmt.Printf("%s\n", string(bytes))
			return nil
		},
	}
}

func queryAbi(cliCtx context.CLIContext, queryRoute string, addr string) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/abi/%s", queryRoute, addr), nil)
	if err != nil {
		return nil, err
	}

	var out types.QueryResAbi
	cliCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Abi, nil
}

// GetCmdMeta returns the CVM metadata query command.
func GetCmdMeta(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "meta <address,hash>",
		Short: "Get CVM Metadata hash for an address or Metadata for a hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			input := args[0]

			var meta string
			var err error
			if strings.HasPrefix(input, common.Bech32PrefixAccAddr) {
				meta, err = queryAddrMeta(cliCtx, queryRoute, input)
			} else {
				meta, err = queryMeta(cliCtx, queryRoute, input)
			}

			if err != nil {
				fmt.Printf("could not get CVM Metadata\n" + err.Error() + "\n")
				return nil
			}
			if len(meta) == 0 {
				fmt.Printf("cannot find CVM Metadata\n")
				return nil
			}

			fmt.Printf("%v\n", meta)
			return nil
		},
	}
}

func queryAddrMeta(cliCtx context.CLIContext, queryRoute string, addr string) (string, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/address-meta/%s", queryRoute, addr), nil)
	if err != nil {
		return "", err
	}

	var out types.QueryResAddrMeta
	err = cliCtx.Codec.UnmarshalJSON(res, &out)
	return out.Metahash, err
}

func queryMeta(cliCtx context.CLIContext, queryRoute string, addr string) (string, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/meta/%s", queryRoute, addr), nil)
	if err != nil {
		return "", err
	}

	var out types.QueryResMeta
	err = cliCtx.Codec.UnmarshalJSON(res, &out)
	return out.Meta, err
}

// GetCmdAddressTranslate is a utility query to translate Bech32 addresses to hex and vice versa.
// It is a pure function and does not interact with the handler or keeper.
func GetCmdAddressTranslate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address-translate <address>",
		Short: "Translate a Bech32 address to hex and vice versa",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := args[0]
			errorMsg := `Address is not in a readable format.
			Please supply either a Bech32 address ("certik1...", "certikvaloper1...")
			or a 20-byte hex address ("0x" prefix not necessary).`

			config := sdk.GetConfig()

			if strings.HasPrefix(addr, config.GetBech32ConsensusAddrPrefix()) { // Bech32 to hex
				consAddr, err := sdk.ConsAddressFromBech32(addr)
				if err != nil {
					fmt.Println(errorMsg)
					return err
				}
				fmt.Println(hex.EncodeToString(consAddr))
				return nil
			} else if strings.HasPrefix(addr, config.GetBech32ValidatorAddrPrefix()) { // Bech32 to hex
				valAddr, err := sdk.ValAddressFromBech32(addr)
				if err != nil {
					fmt.Println(errorMsg)
					return err
				}
				fmt.Println(hex.EncodeToString(valAddr))
				return nil
			} else if strings.HasPrefix(addr, config.GetBech32AccountAddrPrefix()) { // Bech32 to hex
				accAddr, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					fmt.Println(errorMsg)
					return err
				}
				fmt.Println(hex.EncodeToString(accAddr))
				return nil
			} else { // hex to Bech32
				if len(strings.TrimSpace(addr)) != sdk.AddrLen*2 {
					fmt.Println(errorMsg)
					return errors.New("address needs to be 20 bytes")
				}
				accAddr, err := sdk.AccAddressFromHex(addr)
				if err != nil {
					fmt.Println(errorMsg)
					return err
				}
				fmt.Println(accAddr.String())
				valAddr, err := sdk.ValAddressFromHex(addr)
				if err != nil {
					fmt.Println(errorMsg)
					return err
				}
				fmt.Println(valAddr.String())
				consAddr, err := sdk.ConsAddressFromHex(addr)
				if err != nil {
					fmt.Println(errorMsg)
					return err
				}
				fmt.Println(consAddr.String())

				return nil
			}
		},
	}
}

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
func GetAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [address]",
		Short: "Query contract info",
		Long:  "Query contract info by address, revert to query normal account info if the address is not a contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			var account exported.Account
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cvm/account/%s", address), nil)
			if err != nil {
				return err
			}
			cdc.MustUnmarshalJSON(res, &account)

			cvmAcc, err := utils.QueryCVMAccount(cliCtx, address, account)
			if err == nil {
				return cliCtx.PrintOutput(cvmAcc)
			}

			return cliCtx.PrintOutput(account)
		},
	}

	return flags.GetCommands(cmd)[0]
}
