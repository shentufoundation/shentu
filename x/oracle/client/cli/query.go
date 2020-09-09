package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

const (
	FlagOperator = "operator"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	oracleQueryCmds := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Oracle staking subcommands",
	}

	oracleQueryCmds.AddCommand(flags.GetCommands(
		GetCmdOperator(queryRoute, cdc),
		GetCmdOperators(queryRoute, cdc),
		GetCmdWithdraws(queryRoute, cdc),
		GetCmdTask(queryRoute, cdc),
		GetCmdResponse(queryRoute, cdc),
	)...)

	return oracleQueryCmds
}

// GetCmdOperator returns the operator query command.
func GetCmdOperator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator <address>",
		Short: "Get operator information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/operator/%s", queryRoute, args[0]), nil)
			if err != nil {
				return err
			}
			var out types.Operator
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdOperators returns the operators query command.
func GetCmdOperators(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operators",
		Short: "Get operators information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/operators", queryRoute), nil)
			if err != nil {
				return err
			}
			var out types.Operators
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdWithdraws returns the withdrawals query command.
func GetCmdWithdraws(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraws",
		Short: "Get all withdrawals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/withdraws", queryRoute), nil)
			if err != nil {
				return err
			}
			var out types.Withdraws
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdTask returns the task query command.
func GetCmdTask(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task <flags>",
		Short: "Get task information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			contract := viper.GetString(FlagContract)
			if contract == "" {
				return fmt.Errorf("contract address is required")
			}
			function := viper.GetString(FlagFunction)
			if function == "" {
				return fmt.Errorf("function is required")
			}

			params := types.NewQueryTaskParams(contract, function)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/task", queryRoute), bz)
			if err != nil {
				return err
			}
			var out types.Task
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().String(FlagContract, "", "Provide the contract address")
	cmd.Flags().String(FlagFunction, "", "Provide the function")
	return cmd
}

// GetCmdResponse returns the response query command.
func GetCmdResponse(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "response <flags>",
		Short: "Get response information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			contract := viper.GetString(FlagContract)
			if contract == "" {
				return fmt.Errorf("contract address is required")
			}
			function := viper.GetString(FlagFunction)
			if function == "" {
				return fmt.Errorf("function is required")
			}
			operatorStr := viper.GetString(FlagOperator)
			if operatorStr == "" {
				return fmt.Errorf("opeartor Address is required")
			}
			operatorAddress, err := sdk.AccAddressFromBech32(operatorStr)
			if err != nil {
				return err
			}
			params := types.NewQueryResponseParams(contract, function, operatorAddress)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/response", queryRoute), bz)
			if err != nil {
				return err
			}
			var out types.Response
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().String(FlagContract, "", "Provide the contract address")
	cmd.Flags().String(FlagFunction, "", "Provide the function")
	cmd.Flags().String(FlagOperator, "", "Provide the operator")
	return cmd
}
