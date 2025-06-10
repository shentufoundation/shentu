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
		GetCmdQueryProgram(),
		GetCmdQueryPrograms(),
		GetCmdQueryFinding(),
		GetCmdQueryFindings(),
		GetCmdQueryFindingFingerprint(),
		GetCmdQueryProgramFingerprint(),
		GetCmdQueryTheorem(),
		GetCmdQueryProof(),
		GetCmdQueryTheorems(),
		GetCmdQueryRewards(),
		GetCmdQueryParams(),
		GetCmdQueryProofs(),
		GetCmdQueryGrants(),
	)

	return bountyQueryCmd
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

			// Query the program
			res, err := queryClient.Program(
				cmd.Context(),
				&types.QueryProgramRequest{ProgramId: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryPrograms implements the query programs command.
func GetCmdQueryPrograms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "programs",
		Short: "Query all programs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated programs that match optional filters.

Example:
$ %s query bounty programs --page=1 --limit=100
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

	flags.AddPaginationFlagsToCmd(cmd, "programs")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFinding implements the query finding command.
func GetCmdQueryFinding() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finding [finding-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single finding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a finding. You can find the finding-id by running "%s query bounty findings".
Example:
$ %s query bounty finding 1
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

			// Query the finding
			res, err := queryClient.Finding(
				cmd.Context(),
				&types.QueryFindingRequest{
					FindingId: args[0],
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFindings implements the query findings command.
func GetCmdQueryFindings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "findings",
		Short: "Query findings with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated findings that match optional filters.

Example:
$ %s query bounty findings --program-id 1
$ %s query bounty findings --submitter-address cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query bounty findings --page=1 --limit=100
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			// validate that the program-id
			pid, err := cmd.Flags().GetString(FlagProgramID)
			if err != nil {
				return err
			}

			submitterAddr, err := cmd.Flags().GetString(FlagSubmitterAddress)
			if err != nil {
				return err
			}
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

			req := &types.QueryFindingsRequest{
				SubmitterAddress: submitterAddr,
				Pagination:       pageReq,
			}
			if len(pid) != 0 {
				req.ProgramId = pid
			}

			if len(req.ProgramId) == 0 && len(req.SubmitterAddress) == 0 {
				return fmt.Errorf("invalid request")
			}
			res, err := queryClient.Findings(cmd.Context(), req)
			if err != nil {
				return err
			}

			if len(res.GetFindings()) == 0 {
				return fmt.Errorf("no finding found")
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagProgramID, "", "(optional) filter by programs find by program id")
	cmd.Flags().String(FlagSubmitterAddress, "", "(optional) filter by programs find by submitter address")
	flags.AddPaginationFlagsToCmd(cmd, "findings")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFindingFingerprint implements the query finding fingerPrint command.
func GetCmdQueryFindingFingerprint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finding-fingerprint [finding-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query fingerPrint of a single finding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query fingerPrint for a finding. You can find the finding-id by running "%s query bounty findings".
Example:
$ %s query bounty finding 1
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

			// Query the finding
			res, err := queryClient.FindingFingerprint(
				cmd.Context(),
				&types.QueryFindingFingerprintRequest{
					FindingId: args[0],
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryProgramFingerprint implements the query program fingerPrint command.
func GetCmdQueryProgramFingerprint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program-fingerprint [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query fingerPrint of a single program",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query fingerPrint for a program. You can find the program-id by running "%s query bounty findings".
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

			// Query the program
			res, err := queryClient.ProgramFingerprint(
				cmd.Context(),
				&types.QueryProgramFingerprintRequest{
					ProgramId: args[0],
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTheorem implements the query theorem command.
func GetCmdQueryTheorem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theorem [theorem-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a theorem",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a theorem. You can find the theorem-id by running "%s query theorem".
Example:
$ %s query bounty theorem 1
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

			theoremID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("theorem-id %s is not a valid uint, please input a valid theorem-id", args[0])
			}

			// Query the program
			res, err := queryClient.Theorem(
				cmd.Context(),
				&types.QueryTheoremRequest{TheoremId: theoremID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryProof implements the query proof command.
func GetCmdQueryProof() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proof [proof-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a proof",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a proof. You can find the proof-id by running "%s query proof".
Example:
$ %s query bounty proof "hash"
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

			// Query the program
			res, err := queryClient.Proof(
				cmd.Context(),
				&types.QueryProofRequest{ProofId: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTheorems implements the query all theorems command.
func GetCmdQueryTheorems() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theorems",
		Short: "Query all theorems",
		Long:  "Query all theorems with optional pagination",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Theorems(cmd.Context(), &types.QueryTheoremsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "theorems")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRewards implements the query rewards command.
func GetCmdQueryRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query rewards for an address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query rewards for a given address.

Example:
$ %s query bounty rewards [address]
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query the rewards
			res, err := queryClient.Reward(
				cmd.Context(),
				&types.QueryRewardsRequest{Address: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current bounty module parameters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the current bounty module parameters.

Example:
$ %s query bounty params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query the params
			res, err := queryClient.Params(
				cmd.Context(),
				&types.QueryParamsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryProofs implements the query proofs command.
func GetCmdQueryProofs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proofs [theorem-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all proofs for a theorem",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all proofs for a theorem by theorem ID.

Example:
$ %s query bounty proofs 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			theoremID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("theorem-id %s not a valid uint", args[0])
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Proofs(
				cmd.Context(),
				&types.QueryProofsRequest{
					TheoremId:  theoremID,
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "proofs")
	return cmd
}

// GetCmdQueryGrants implements the query grants command.
func GetCmdQueryGrants() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grants [theorem-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query grants for a theorem",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query grants for a theorem by theorem ID.

Example:
$ %s query bounty grants 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			theoremID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("theorem-id %s is not a valid uint, please input a valid theorem-id", args[0])
			}

			// Query the grants
			res, err := queryClient.Grants(
				cmd.Context(),
				&types.QueryGrantsRequest{
					TheoremId: theoremID,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
