package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govUtils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// Proposal flags
const (
	flagVoter     = "voter"
	flagDepositor = "depositor"
	flagStatus    = "status"
)

// TODO remove

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:                        govtypes.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(
		GetCmdQueryProposal(),
		GetCmdQueryProposals(),
		GetCmdQueryVote(),
		GetCmdQueryVotes(),
		GetCmdQueryParam(),
		GetCmdQueryParams(),
		GetCmdQueryTally(),
		GetCmdCertVoted(),
		GetCmdQueryCustomParam(),
	)

	return govQueryCmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal() *cobra.Command {
	cmd := &cobra.Command{
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
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := govtypesv1.NewQueryClient(cliCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint, please input a valid proposal-id", args[0])
			}

			// query the proposal
			res, err := queryClient.Proposal(
				cmd.Context(),
				&govtypesv1.QueryProposalRequest{ProposalId: proposalID},
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

// GetCmdQueryProposals implements a query proposals command.
func GetCmdQueryProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated proposals that match optional filters:

Example:
$ %[1]s query gov proposals --depositor shentu1skjwj5whet0lpe65qaq4rpq03hjxlwd9ma4udt
$ %[1]s query gov proposals --voter shentu1skjwj5whet0lpe65qaq4rpq03hjxlwd9ma4udt
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

			var proposalStatus govtypesv1.ProposalStatus
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
				proposalStatus, err = govtypesv1.ProposalStatusFromString(govUtils.NormalizeProposalStatus(strProposalStatus))
				if err != nil {
					return err
				}
			}

			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := govtypesv1.NewQueryClient(cliCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Proposals(
				cmd.Context(),
				&govtypesv1.QueryProposalsRequest{
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

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "proposals")
	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(flagDepositor, "", "(optional) filter by proposals deposited on by depositor")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected")

	return cmd
}

// GetCmdQueryVote implements the query proposal vote command. Command to Get a
// Proposal Information.
func GetCmdQueryVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [voter-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single vote on a proposal given its identifier.

Example:
$ %s query gov vote 1 shentu16gzt5vd0dd5c98ajl3ld2ltvcahxgyygd58n3m
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := govtypesv1.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			ctx := cmd.Context()
			_, err = queryClient.Proposal(
				ctx,
				&govtypesv1.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			voterAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.Vote(
				ctx,
				&govtypesv1.QueryVoteRequest{ProposalId: proposalID, Voter: args[1]},
			)
			if err != nil {
				return err
			}

			vote := res.GetVote()
			if vote.Empty() {
				params := govtypesv1.NewQueryVoteParams(proposalID, voterAddr)
				resByTxQuery, err := govUtils.QueryVoteByTxQuery(clientCtx, params)

				if err != nil {
					return err
				}

				if err := clientCtx.Codec.UnmarshalJSON(resByTxQuery, vote); err != nil {
					return err
				}
			}

			return clientCtx.PrintProto(res.Vote)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

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
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := govtypesv1.NewQueryClient(cliCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			proposalRes, err := queryClient.Proposal(
				cmd.Context(),
				&govtypesv1.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// TODO Query tx depending on proposal status?
			propStatus := proposalRes.GetProposal().Status
			if !(propStatus == govtypesv1.StatusVotingPeriod || propStatus == govtypesv1.StatusDepositPeriod) {
				page, _ := cmd.Flags().GetInt(flags.FlagPage)
				limit, _ := cmd.Flags().GetInt(flags.FlagLimit)

				params := govtypesv1.NewQueryProposalVotesParams(proposalID, page, limit)
				resByTxQuery, err := govUtils.QueryVotesByTxQuery(cliCtx, params)
				if err != nil {
					return err
				}

				var votes govtypesv1.Votes
				// TODO migrate to use JSONCodec (implement MarshalJSONArray
				// or wrap lists of proto.Message in some other message)
				cliCtx.LegacyAmino.MustUnmarshalJSON(resByTxQuery, &votes)
				return cliCtx.PrintObjectLegacy(votes)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Votes(
				cmd.Context(),
				&govtypesv1.QueryVotesRequest{ProposalId: proposalID, Pagination: pageReq},
			)

			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "votes")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTally implements the command to query for proposal tally result.
func GetCmdQueryTally() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tally [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get the tally of a proposal vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query tally of votes on a proposal. You can find
the proposal-id by running "%s query gov proposals".

Example:
$ %s query gov tally 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := govtypesv1.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			ctx := cmd.Context()
			_, err = queryClient.Proposal(
				ctx,
				&govtypesv1.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// Query store
			res, err := queryClient.TallyResult(
				ctx,
				&govtypesv1.QueryTallyResultRequest{ProposalId: proposalID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Tally)
		},
	}

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
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := govtypesv1.NewQueryClient(cliCtx)

			// Query store for all 3 params
			votingRes, err := queryClient.Params(
				cmd.Context(),
				&govtypesv1.QueryParamsRequest{ParamsType: "voting"},
			)
			if err != nil {
				return err
			}

			tallyRes, err := queryClient.Params(
				cmd.Context(),
				&govtypesv1.QueryParamsRequest{ParamsType: "tallying"},
			)
			if err != nil {
				return err
			}

			depositRes, err := queryClient.Params(
				cmd.Context(),
				&govtypesv1.QueryParamsRequest{ParamsType: "deposit"},
			)
			if err != nil {
				return err
			}

			res := &govtypesv1.QueryParamsResponse{
				VotingParams:  votingRes.GetVotingParams(),
				DepositParams: depositRes.GetDepositParams(),
				TallyParams:   tallyRes.GetTallyParams(),
				Params:        tallyRes.Params,
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
			queryClient := govtypesv1.NewQueryClient(cliCtx)

			// Query store
			res, err := queryClient.Params(
				cmd.Context(),
				&govtypesv1.QueryParamsRequest{ParamsType: args[0]},
			)
			if err != nil {
				return err
			}

			var out fmt.Stringer
			switch args[0] {
			case govtypesv1.ParamVoting:
				out = res.GetVotingParams()
			case govtypesv1.ParamTallying:
				out = res.GetTallyParams()
			case govtypesv1.ParamDeposit:
				out = res.GetDepositParams()
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit), was %s", args[0])
			}

			return cliCtx.PrintObjectLegacy(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryCustomParam implements the query param command.
func GetCmdQueryCustomParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-param [param-type]",
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
			queryClient := typesv1.NewCustomQueryClient(cliCtx)

			// Query store
			res, err := queryClient.CustomParams(
				cmd.Context(),
				&govtypesv1.QueryParamsRequest{ParamsType: args[0]},
			)
			if err != nil {
				return err
			}

			var out fmt.Stringer
			switch args[0] {
			case typesv1.ParamCustom:
				out = res.GetCustomParams()
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit), was %s", args[0])
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
			queryClient := typesv1.NewCustomQueryClient(cliCtx)

			proposalID, err := strconv.ParseUint(args[0], 10, 64)

			if err != nil {
				return err
			}

			// Query store
			res, err := queryClient.CertVoted(
				cmd.Context(),
				&typesv1.QueryCertVotedRequest{ProposalId: proposalID},
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
