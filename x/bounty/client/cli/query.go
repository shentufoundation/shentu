package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	bountyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bounty module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bountyQueryCmd.AddCommand(
		GetCmdQueryHost(),
		GetCmdQueryHosts(),
		GetCmdQueryProgram(),
		GetCmdQueryPrograms(),
		GetCmdQueryFinding(),
		GetCmdQueryFindings(),
	)

	return bountyQueryCmd
}

// GetCmdQueryHost implements the query host command.
func GetCmdQueryHost() *cobra.Command {
	//TODO implement me
	cmd := &cobra.Command{}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryHosts implements the query hosts command. Command to Get a
// Host Information list.
func GetCmdQueryHosts() *cobra.Command {
	//TODO implement me
	cmd := &cobra.Command{}

	flags.AddPaginationFlagsToCmd(cmd, "hosts")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryProgram implements the query program command.
func GetCmdQueryProgram() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single program",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a program. You can find the program-id by running "%s query bounty program".
Example:
$ %s query bounty program 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal_id is an uint
			programID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid uint, please input a valid program-id", args[0])
			}

			// Query the program
			res, err := queryClient.Program(
				cmd.Context(),
				&types.QueryProgramRequest{ProgramId: programID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Program)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryPrograms implements the query programs command.
func GetCmdQueryPrograms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "programs",
		Short: "Query programs with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated programs that match optional filters:

Example:
$ %s query bounty programs --submitter cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query bounty programs --active (true|false)
$ %s query bounty programs --page=2 --limit=100
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			submitterAddr, _ := cmd.Flags().GetString(FlagSubmitter)

			if len(submitterAddr) != 0 {
				_ = sdk.MustAccAddressFromBech32(submitterAddr)
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Programs(
				cmd.Context(),
				&types.QueryProgramsRequest{
					Submitter:  submitterAddr,
					Pagination: pageReq,
				})
			if err != nil {
				return err
			}

			if len(res.GetPrograms()) == 0 {
				return fmt.Errorf("no programs found")
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagSubmitter, "", "(optional) filter by programs submitted on by submitter")
	flags.AddPaginationFlagsToCmd(cmd, "programs")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFinding implements the query finding command.
func GetCmdQueryFinding() *cobra.Command {
	//TODO implement me
	cmd := &cobra.Command{}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFindings implements the query findings command.
func GetCmdQueryFindings() *cobra.Command {
	//TODO implement me
	cmd := &cobra.Command{}

	flags.AddPaginationFlagsToCmd(cmd, "findings")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
