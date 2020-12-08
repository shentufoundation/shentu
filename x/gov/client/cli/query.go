package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govUtils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:                        govTypes.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(
		GetCmdQueryProposal(),
		GetCmdQueryProposals(),
		cli.GetCmdQueryVote(),
		GetCmdQueryVotes(),
		GetCmdQueryParam(),
		GetCmdQueryParams(),
		cli.GetCmdQueryProposer(),
		cli.GetCmdQueryDeposit(),
		GetCmdQueryDeposits(),
		cli.GetCmdQueryTally(),
	)

	return govQueryCmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal() *cobra.Command {
	return &cobra.Command{
		Use:   "proposal [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a proposal. You can find the
proposal-id by running "%[1]s query gov proposals".

Example:
$ %[1]s query gov proposal 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint, please input a valid proposal-id", args[0])
			}

			// query the proposal
			res, err := queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(&res.Proposal)

		},
	}
}

// GetCmdQueryProposals implements a query proposals command.
func GetCmdQueryProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated proposals that match optional filters:

Example:
$ %[1]s query gov proposals --depositor certik1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %[1]s query gov proposals --voter certik1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %[1]s query gov proposals --status (DepositPeriod|CertifierVotingPeriod|ValidatorVotingPeriod|Passed|Rejected)
$ %[1]s query gov proposals --page=2 --limit=100
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechDepositorAddr := viper.GetString(flagDepositor)
			bechVoterAddr := viper.GetString(flagVoter)
			strProposalStatus := viper.GetString(flagStatus)

			var proposalStatus types.ProposalStatus
			var err error
			if bechDepositorAddr != "" {
				_, err = sdk.AccAddressFromBech32(bechDepositorAddr)
				if err != nil {
					return err
				}
			}

			if bechVoterAddr != "" {
				_, err = sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			if strProposalStatus != "" {
				proposalStatus, err = types.ProposalStatusFromString(govUtils.NormalizeProposalStatus(strProposalStatus))
				if err != nil {
					return err
				}
			}

			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Proposals(
				context.Background(),
				&types.QueryProposalsRequest{
					ProposalStatus: proposalStatus,
					Voter:          bechVoterAddr,
					Depositor:      bechDepositorAddr,
					Pagination:     pageReq,
				},
			)
			if err != nil {
				return err
			}

			if len(res.GetProposals()) == 0 {
				return fmt.Errorf("no proposals found")
			}

			return cliCtx.PrintOutput(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "proposals")
	cmd.Flags().String(flagDepositor, "", "(optional) filter by proposals deposited on by depositor")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected")

	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query vote details for a single proposal by its identifier.

Example:
$ %[1]s query gov votes 1
$ %[1]s query gov votes 1 --page=2 --limit=100
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// TODO Query tx depending on proposal status?

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Votes(
				context.Background(),
				&types.QueryVotesRequest{ProposalId: proposalID, Pagination: pageReq},
			)

			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)				
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "votes")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDeposits implements the command to query for proposal deposits.
func GetCmdQueryDeposits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query deposits on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for all deposits on a proposal.
You can find the proposal-id by running "%[1]s query gov proposals".

Example:
$ %[1]s query gov deposits 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint, please input a valid proposal-id", args[0])
			}

			// params := govTypes.NewQueryProposalParams(proposalID)
			// bz, err := cdc.MarshalJSON(params)
			// if err != nil {
			// 	return err
			// }

			// check to see if the proposal is in the store
			proposalRes, err := queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// TODO
			//propStatus := proposalRes.GetProposal().Status
			// if !(propStatus == types.StatusCertifierVotingPeriod ||
			// 	propStatus == types.StatusValidatorVotingPeriod ||
			// 	propStatus == types.StatusDepositPeriod) {
			// 	res, err = govUtils.QueryDepositsByTxQuery(cliCtx, params)
			// } else {
			//	res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposits", queryRoute), bz)
			// }

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Deposits(
				context.Background(),
				&types.QueryDepositsRequest{ProposalId: proposalID, Pagination: pageReq},
			)

			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "deposits")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
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
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			// Query store for all 3 params
			votingRes, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: "voting"},
			)
			if err != nil {
				return err
			}

			tallyRes, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: "tallying"},
			)
			if err != nil {
				return err
			}

			depositRes, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: "deposit"},
			)
			if err != nil {
				return err
			}

			params := types.NewParams(
				votingRes.GetVotingParams(),
				tallyRes.GetTallyParams(),
				depositRes.GetDepositParams(),
			)

			return cliCtx.PrintOutputLegacy(params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParam implements the query param command.
func GetCmdQueryParam() *cobra.Command {
	return &cobra.Command{
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
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			// Query store
			res, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: args[0]},
			)
			if err != nil {
				return err
			}

			var out fmt.Stringer
			switch args[0] {
			case "voting":
				out = res.GetVotingParams()
			case "tallying":
				out = res.GetTallyParams()
			case "deposit":
				out = res.GetDepositParams()
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit), was %s", args[0])
			}

			return cliCtx.PrintOutputLegacy(out)
		},
	}
}
