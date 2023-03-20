package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

// Proposal flags
const (
	flagVoter     = "voter"
	flagDepositor = "depositor"
	flagStatus    = "status"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := cli.GetQueryCmd()

	govQueryCmd.RemoveCommand(
		cli.GetCmdQueryParam(),
		cli.GetCmdQueryParams(),
	)

	govQueryCmd.AddCommand(
		GetCmdQueryParam(),
		GetCmdQueryParams(),
		GetCmdCertVoted(),
	)

	return govQueryCmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.

Example:
$ %s query gov params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			// Query store for all 3 params
			votingRes, err := queryClient.Params(
				cmd.Context(),
				&govtypes.QueryParamsRequest{ParamsType: "voting"},
			)
			if err != nil {
				return err
			}

			tallyRes, err := queryClient.Params(
				cmd.Context(),
				&govtypes.QueryParamsRequest{ParamsType: "tallying"},
			)
			if err != nil {
				return err
			}

			depositRes, err := queryClient.Params(
				cmd.Context(),
				&govtypes.QueryParamsRequest{ParamsType: "deposit"},
			)
			if err != nil {
				return err
			}

			customRes, err := queryClient.Params(
				cmd.Context(),
				&govtypes.QueryParamsRequest{ParamsType: "custom"},
			)
			if err != nil {
				return err
			}

			res := &types.QueryParamsResponse{
				VotingParams:  votingRes.GetVotingParams(),
				DepositParams: depositRes.GetDepositParams(),
				TallyParams:   tallyRes.GetTallyParams(),
				CustomParams:  customRes.GetCustomParams(),
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParam implements the query param command.
func GetCmdQueryParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param [param-type]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the parameters (voting|tallying|deposit) of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.

Example:
$ %[1]s query gov param voting
$ %[1]s query gov param tallying
$ %[1]s query gov param deposit
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			// Query store
			res, err := queryClient.Params(
				cmd.Context(),
				&govtypes.QueryParamsRequest{ParamsType: args[0]},
			)
			if err != nil {
				return err
			}

			var out fmt.Stringer
			switch args[0] {
			case govtypes.ParamVoting:
				out = res.GetVotingParams()
			case govtypes.ParamTallying:
				out = res.GetTallyParams()
			case govtypes.ParamDeposit:
				out = res.GetDepositParams()
			case types.ParamCustom:
				out = res.GetCustomParams()
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit|custom), was %s", args[0])
			}

			return cliCtx.PrintObjectLegacy(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdCertVoted implements the query param command.
func GetCmdCertVoted() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cert-voted [proposa-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query if the certifiers voted on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query if the certifiers voted on a proposal.

Example:
$ %[1]s query gov cert-voted 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			proposalId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			// Query store
			res, err := queryClient.CertVoted(
				cmd.Context(),
				&types.QueryCertVotedRequest{ProposalId: proposalId},
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
