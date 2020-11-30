package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govUtils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/client/utils"
	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:                        govTypes.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryProposal(queryRoute, cdc),
		GetCmdQueryProposals(queryRoute, cdc),
		cli.GetCmdQueryVote(queryRoute, cdc),
		GetCmdQueryVotes(queryRoute, cdc),
		GetCmdQueryParam(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryProposer(queryRoute, cdc),
		cli.GetCmdQueryDeposit(queryRoute, cdc),
		GetCmdQueryDeposits(queryRoute, cdc),
		cli.GetCmdQueryTally(queryRoute, cdc))...)

	return govQueryCmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint, please input a valid proposal-id", args[0])
			}

			// query the proposal
			res, err := govUtils.QueryProposalByID(proposalID, cliCtx, queryRoute)
			if err != nil {
				return err
			}

			var proposal types.Proposal
			cdc.MustUnmarshalJSON(res, &proposal)
			return cliCtx.PrintOutput(proposal)
		},
	}
}

// GetCmdQueryProposals implements a query proposals command.
func GetCmdQueryProposals(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechDepositorAddr := viper.GetString(flagDepositor)
			bechVoterAddr := viper.GetString(flagVoter)
			strProposalStatus := viper.GetString(flagStatus)
			page := viper.GetInt(flags.FlagPage)
			limit := viper.GetInt(flags.FlagLimit)

			var depositorAddr sdk.AccAddress
			var voterAddr sdk.AccAddress
			var proposalStatus types.ProposalStatus

			params := types.NewQueryProposalsParams(page, limit, proposalStatus, voterAddr, depositorAddr)

			if bechDepositorAddr != "" {
				depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
				if err != nil {
					return err
				}
				params.Depositor = depositorAddr
			}

			if bechVoterAddr != "" {
				voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
				params.Voter = voterAddr
			}

			if strProposalStatus != "" {
				proposalStatus, err := types.ProposalStatusFromString(govUtils.NormalizeProposalStatus(strProposalStatus))
				if err != nil {
					return err
				}
				params.ProposalStatus = proposalStatus
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposals", queryRoute), bz)
			if err != nil {
				return err
			}

			var matchingProposals types.Proposals
			err = cdc.UnmarshalJSON(res, &matchingProposals)
			if err != nil {
				return err
			}

			if len(matchingProposals) == 0 {
				return fmt.Errorf("no matching proposals found")
			}

			return cliCtx.PrintOutput(matchingProposals)
		},
	}

	cmd.Flags().Int(flags.FlagPage, 1, "pagination page of proposals to to query for")
	cmd.Flags().Int(flags.FlagLimit, 100, "pagination limit of proposals to query for")
	cmd.Flags().String(flagDepositor, "", "(optional) filter by proposals deposited on by depositor")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected")

	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			page := viper.GetInt(flags.FlagPage)
			limit := viper.GetInt(flags.FlagLimit)

			params := govTypes.NewQueryProposalVotesParams(proposalID, page, limit)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			// check to see if the proposal is in the store
			res, err := govUtils.QueryProposalByID(proposalID, cliCtx, queryRoute)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			var proposal types.Proposal
			cdc.MustUnmarshalJSON(res, &proposal)

			propStatus := proposal.Status
			if !(propStatus == types.StatusCertifierVotingPeriod ||
				propStatus == types.StatusValidatorVotingPeriod ||
				propStatus == types.StatusDepositPeriod) {
				res, err = govUtils.QueryVotesByTxQuery(cliCtx, params)
			} else {
				res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/votes", queryRoute), bz)
			}

			if err != nil {
				return err
			}

			var votes types.Votes
			cdc.MustUnmarshalJSON(res, &votes)
			return cliCtx.PrintOutput(votes)
		},
	}
	cmd.Flags().Int(flags.FlagPage, 1, "pagination page of votes to to query for")
	cmd.Flags().Int(flags.FlagLimit, 100, "pagination limit of votes to query for")
	return cmd
}

// GetCmdQueryDeposits implements the command to query for proposal deposits.
func GetCmdQueryDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposits [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query deposits on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for all deposits on a proposal.
You can find the proposal-id by running "%[1]s query gov proposals".

Example:
$ %[1]s query gov deposits 1
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint, please input a valid proposal-id", args[0])
			}

			params := govTypes.NewQueryProposalParams(proposalID)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			// check to see if the proposal is in the store
			res, err := govUtils.QueryProposalByID(proposalID, cliCtx, queryRoute)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal with id %d: %s", proposalID, err)
			}

			var proposal types.Proposal
			cdc.MustUnmarshalJSON(res, &proposal)

			propStatus := proposal.Status
			if !(propStatus == types.StatusCertifierVotingPeriod ||
				propStatus == types.StatusValidatorVotingPeriod ||
				propStatus == types.StatusDepositPeriod) {
				res, err = govUtils.QueryDepositsByTxQuery(cliCtx, params)
			} else {
				res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposits", queryRoute), bz)
			}

			if err != nil {
				return err
			}

			var dep types.Deposits
			cdc.MustUnmarshalJSON(res, &dep)
			return cliCtx.PrintOutput(dep)
		},
	}
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.

Example:
$ %s query gov params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/tallying", queryRoute), nil)
			if err != nil {
				return err
			}
			dp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/deposit", queryRoute), nil)
			if err != nil {
				return err
			}
			vp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/voting", queryRoute), nil)
			if err != nil {
				return err
			}

			var tallyParams types.TallyParams
			cdc.MustUnmarshalJSON(tp, &tallyParams)
			var depositParams types.DepositParams
			cdc.MustUnmarshalJSON(dp, &depositParams)
			var votingParams govTypes.VotingParams
			cdc.MustUnmarshalJSON(vp, &votingParams)

			return cliCtx.PrintOutput(types.NewParams(votingParams, tallyParams, depositParams))
		},
	}
}

// GetCmdQueryParam implements the query param command.
func GetCmdQueryParam(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/%s", queryRoute, args[0]), nil)
			if err != nil {
				return err
			}
			var out fmt.Stringer
			switch args[0] {
			case "voting":
				var param govTypes.VotingParams
				cdc.MustUnmarshalJSON(res, &param)
				out = param
			case "tallying":
				var param govTypes.TallyParams
				cdc.MustUnmarshalJSON(res, &param)
				out = param
			case "deposit":
				var param types.DepositParams
				cdc.MustUnmarshalJSON(res, &param)
				out = param
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit), was %s", args[0])
			}

			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryProposer implements the query proposer command.
func GetCmdQueryProposer(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proposer [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the proposer of a governance proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query which address proposed a proposal with a given ID.

Example:
$ %s query gov proposer 1
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the proposalID is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint", args[0])
			}

			prop, err := utils.QueryProposerByTxQuery(cliCtx, proposalID, queryRoute)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(prop)
		},
	}
}
