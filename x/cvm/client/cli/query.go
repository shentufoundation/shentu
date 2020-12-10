// Package cli defines the CLI services for the cvm module.
package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/execution/evm/abi"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

const (
	FlagCaller = "caller"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group cvm queries under a subcommand
	cvmQueryCmd := &cobra.Command{
		Use:   "cvm",
		Short: "Querying commands for the CVM module",
	}

	cvmQueryCmd.AddCommand(
		GetCmdCode(),
		GetCmdStorage(),
		GetCmdAbi(),
		GetCmdMeta(),
		GetCmdView(),
		GetCmdAddressTranslate(),
	)

	return cvmQueryCmd
}

// GetCmdView returns the CVM contract view transaction command.
func GetCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <address> <function> [<params>...]",
		Short: "View CVM contract",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// Caller is an optional flag. If not set it becomes the zero address.
			var callerString string
			callerString, err = cmd.Flags().GetString(FlagCaller)
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

			abiSpec, data, err := parseCallCmd(clientCtx, args[0], callee, args[1], args[2:])
			if err != nil {
				return err
			}

			queryPath := fmt.Sprintf("custom/%s/view/%s/%s", queryRoute, callerString, callee)
			return queryContractAndPrint(clientCtx, cdc, queryPath, args[1], abiSpec, data)
		},
	}
	cmd.Flags().String(FlagCaller, "", "optional caller parameter to run the view function with")

	return cmd
}

// Query CVM contract code based on ABI spec and print function output.
func queryContractAndPrint(cliCtx client.Context, cdc *codec.LegacyAmino, queryPath, fname string, abiSpec, data []byte) error {
	res, _, err := cliCtx.QueryWithData(queryPath, data)
	if err != nil {
		return fmt.Errorf("querying CVM contract code: %v", err)
	}
	var out types.QueryResView
	cliCtx.LegacyAmino.MustUnmarshalJSON(res, &out)
	ret, err := abi.DecodeFunctionReturn(string(abiSpec), fname, out.Ret)
	if err != nil {
		return fmt.Errorf("decoding function return: %v", err)
	}
	return cliCtx.PrintOutput(ret)
}

// GetCmdCode returns the CVM code query command.
func GetCmdCode() *cobra.Command {
	return &cobra.Command{
		Use:   "code <address>",
		Short: "Get CVM contract code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			addr := args[0]

			req := &types.QueryCodeRequest{
				Address: addr,
			}

			res, err := queryClient.Code(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}
}

// GetCmdStorage returns the CVM storage query command.
func GetCmdStorage() *cobra.Command {
	return &cobra.Command{
		Use:   "storage <address> <key>",
		Short: "Get CVM storage data",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			addr := args[0]
			key := args[1]

			req := &types.QueryStorageRequest{
				Address: addr,
				Key:     key,
			}

			res, err := queryClient.Storage(cmd.Context(), req)
			if err != nil {
				fmt.Printf("could not get CVM storage\n" + err.Error() + "\n")
				return nil
			}

			return clientCtx.PrintOutput(res)
		},
	}
}

// GetCmdAbi returns the CVM code ABI query command.
func GetCmdAbi() *cobra.Command {
	return &cobra.Command{
		Use:   "abi <address>",
		Short: "Get CVM contract code ABI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			addr := args[0]

			req := &types.QueryAbiRequest{
				Address: addr,
			}

			addrAbi, err := queryClient.Abi(cmd.Context(), req)
			if err != nil {
				fmt.Printf("could not get CVM contract code ABI\n" + err.Error() + "\n")
				return nil
			}
			if addrAbi == nil {
				fmt.Printf("cannot find CVM contract code ABI\n")
				return nil
			}

			return clientCtx.PrintOutput(addrAbi)
		},
	}
}

// GetCmdMeta returns the CVM metadata query command.
func GetCmdMeta() *cobra.Command {
	return &cobra.Command{
		Use:   "meta <address,hash>",
		Short: "Get CVM Metadata hash for an address or Metadata for a hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			addrPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

			queryClient := types.NewQueryClient(clientCtx)
			input := args[0]

			var meta interface{}
			if strings.HasPrefix(input, addrPrefix) {
				request := &types.QueryAddressMetaRequest{Address: input}
				meta, err = queryClient.AddressMeta(cmd.Context(), request)
			} else {
				request := &types.QueryMetaRequest{Hash: input}
				meta, err = queryClient.Meta(cmd.Context(), request)
			}

			if err != nil {
				fmt.Printf("could not get CVM Metadata\n" + err.Error() + "\n")
				return nil
			}

			fmt.Printf("%v\n", meta)
			return nil
		},
	}
}

// GetCmdAddressTranslate is a utility query to translate Bech32 addresses to hex and vice versa.
// It is a pure function and does not interact with the handler or keeper.
func GetCmdAddressTranslate() *cobra.Command {
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
func GetAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [address]",
		Short: "Query contract info",
		Long:  "Query contract info by address, revert to query normal account info if the address is not a contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			queryAccount := &types.QueryAccountRequest{Address: args[0]}
			res, err := queryClient.Account(cmd.Context(), queryAccount)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	return cmd
}
