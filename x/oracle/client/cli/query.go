package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
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
		GetCmdTxTask(),
		GetCmdTxResponse(),
		GetCmdLeftBounty(),
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
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Operator(
				cmd.Context(),
				&types.QueryOperatorRequest{Address: address.String()},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdOperators returns the operators query command.
func GetCmdOperators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operators",
		Short: "Get operators information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Operators(cmd.Context(), &types.QueryOperatorsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdWithdraws returns the withdrawals query command.
func GetCmdWithdraws() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraws",
		Short: "Get all withdrawals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Withdraws(
				cmd.Context(),
				&types.QueryWithdrawsRequest{},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdTask returns the task query command.
func GetCmdTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task <contract_address> <function>",
		Short: "Get task information",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Task(
				cmd.Context(),
				&types.QueryTaskRequest{Contract: args[0], Function: args[1]},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdTxTask returns the tx task query command.
func GetCmdTxTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx-task <tx_hash>",
		Short: "Get tx task information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.TxTask(
				cmd.Context(),
				&types.QueryTxTaskRequest{TxHash: args[0]},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdResponse returns the response query command.
func GetCmdResponse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "response <operator_address> <contract_address> <function>",
		Short: "Get response information",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Response(
				cmd.Context(),
				&types.QueryResponseRequest{Contract: args[1], Function: args[2], OperatorAddress: args[0]},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagOperator, "", "Provide the operator")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdTxResponse returns the tx response query command.
func GetCmdTxResponse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx-response <operator_address> <tx_hash>",
		Short: "Get tx response information",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.TxResponse(
				cmd.Context(),
				&types.QueryTxResponseRequest{TxHash: args[1], OperatorAddress: args[0]},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdLeftBounty This function fetches the left bounty information for any given address.
func GetCmdLeftBounty() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "left_bounty <address>",
		Short: "Get left bounty information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.LeftBounty(
				cmd.Context(),
				&types.QueryLeftBountyRequest{Address: args[0]},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
