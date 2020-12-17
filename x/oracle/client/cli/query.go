package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

const (
	FlagOperator = "operator"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	oracleQueryCmds := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Oracle staking subcommands",
	}

	oracleQueryCmds.AddCommand(
		GetCmdOperator(),
		GetCmdOperators(),
		GetCmdWithdraws(),
		GetCmdTask(),
		GetCmdResponse(),
	)

	return oracleQueryCmds
}

// GetCmdOperator returns the operator query command.
func GetCmdOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator <address>",
		Short: "Get operator information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			
			res, err := queryClient.Operator(
				context.Background(),
				&types.QueryOperatorRequest{Address: address.String()},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
	return cmd
}

// GetCmdOperators returns the operators query command.
func GetCmdOperators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operators",
		Short: "Get operators information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Operators(context.Background(), &types.QueryOperatorsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
	return cmd
}

// GetCmdWithdraws returns the withdrawals query command.
func GetCmdWithdraws() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraws",
		Short: "Get all withdrawals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Withdraws(
				context.Background(),
				&types.QueryWithdrawsRequest{},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
	return cmd
}

// GetCmdTask returns the task query command.
func GetCmdTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task <flags>",
		Short: "Get task information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			contract := viper.GetString(FlagContract)
			if contract == "" {
				return fmt.Errorf("contract address is required")
			}
			function := viper.GetString(FlagFunction)
			if function == "" {
				return fmt.Errorf("function is required")
			}
			
			res, err := queryClient.Task(
				context.Background(),
				&types.QueryTaskRequest{Contract: contract, Function: function},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
	cmd.Flags().String(FlagContract, "", "Provide the contract address")
	cmd.Flags().String(FlagFunction, "", "Provide the function")
	return cmd
}

// GetCmdResponse returns the response query command.
func GetCmdResponse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "response <flags>",
		Short: "Get response information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

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

			res, err := queryClient.Response(
				context.Background(),
				&types.QueryResponseRequest{Contract: contract, Function: function, OperatorAddress: operatorAddress.String()},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
	cmd.Flags().String(FlagContract, "", "Provide the contract address")
	cmd.Flags().String(FlagFunction, "", "Provide the function")
	cmd.Flags().String(FlagOperator, "", "Provide the operator")
	return cmd
}
