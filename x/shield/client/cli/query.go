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

	"github.com/certikfoundation/shentu/x/shield/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	shieldQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the shield module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	shieldQueryCmd.AddCommand(
		GetCmdPool(),
		GetCmdSponsor(),
		GetCmdPools(),
		GetCmdPurchaseList(),
		GetCmdPurchaserPurchases(),
		GetCmdPoolPurchases(),
		GetCmdPurchases(),
		GetCmdProvider(),
		GetCmdProviders(),
		GetCmdPoolParams(),
		GetCmdClaimParams(),
		GetCmdStatus(),
		GetCmdStaking(),
		GetCmdShieldStakingRate(),
		GetCmdReimbursement(),
		GetCmdReimbursements(),
	)

	return shieldQueryCmd
}

// GetCmdPool returns the command for querying the pool.
func GetCmdPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool_ID]",
		Short: "query a pool",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			sponsor := viper.GetString(flagSponsor)
			var id uint64
			if sponsor == "" {
				id, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("no sponsor was provided, and pool id %s is invalid", args[0])
				}
			}

			res, err := queryClient.Pool(
				cmd.Context(),
				&types.QueryPoolRequest{PoolId: id},
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

// GetCmdSponsor returns the command for querying the pools for a sponsor address.
func GetCmdSponsor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sponsor [pool_ID]",
		Short: "query pools for a sponsor",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			sponsor := args[0]
			res, err := queryClient.Sponsor(
				cmd.Context(),
				&types.QuerySponsorRequest{Sponsor: sponsor},
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

// GetCmdPools returns the command for querying a complete list of pools.
func GetCmdPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "query a complete list of pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Pools(cmd.Context(), &types.QueryPoolsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdPurchaseList returns the command for querying purchases
// corresponding to a given pool-purchaser pair.
func GetCmdPurchaseList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-purchaser [pool_ID] [purchaser_address]",
		Short: "get purchases corresponding to a given pool-purchaser pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool id %s is invalid", args[0])
			}
			purchaser, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.PurchaseList(
				cmd.Context(),
				&types.QueryPurchaseListRequest{PoolId: poolID, Purchaser: purchaser.String()},
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

// GetCmdPurchaserPurchases returns the command for querying
// purchases by a given address.
func GetCmdPurchaserPurchases() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchases-by [purchaser_address]",
		Short: "query purchase information of a given account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			purchaser, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.PurchaseLists(
				cmd.Context(),
				&types.QueryPurchaseListsRequest{Purchaser: purchaser.String()},
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

// GetCmdPoolPurchases returns the command for querying
// purchases in a given pool.
func GetCmdPoolPurchases() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-purchases [pool_ID]",
		Short: "query purchases in a given pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool id %s is invalid", args[0])
			}

			res, err := queryClient.PoolPurchaseLists(
				cmd.Context(),
				&types.QueryPoolPurchaseListsRequest{PoolId: poolID},
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

// GetCmdPurchases returns the command for querying all purchases.
func GetCmdPurchases() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchases",
		Short: "query all purchases",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Purchases(cmd.Context(), &types.QueryPurchasesRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdProvider returns the command for querying a provider.
func GetCmdProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [provider_address]",
		Short: "get provider information",
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

			res, err := queryClient.Provider(
				cmd.Context(),
				&types.QueryProviderRequest{Address: address.String()},
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

// GetCmdProviders returns the command for querying all providers.
func GetCmdProviders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Args:  cobra.ExactArgs(0),
		Short: "query all providers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query providers with pagination parameters

Example:
$ %[1]s query shield providers
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

			res, err := queryClient.Providers(cmd.Context(), &types.QueryProvidersRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdPoolParams returns the command for querying pool parameters.
func GetCmdPoolParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-params",
		Short: "get pool parameters",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.PoolParams(cmd.Context(), &types.QueryPoolParamsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdClaimParams returns the command for querying claim parameters.
func GetCmdClaimParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-params",
		Short: "get claim parameters",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.ClaimParams(cmd.Context(), &types.QueryClaimParamsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdStatus returns the command for querying shield status.
func GetCmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "get shield status",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.ShieldStatus(cmd.Context(), &types.QueryShieldStatusRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdStaking returns the command for querying staked-for-shield amounts
// corresponding to a given pool-purchaser pair.
func GetCmdStaking() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staked-for-shield [pool_ID] [purchaser_address]",
		Short: "get staked CTK for shield corresponding to a given pool-purchaser pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool id %s is invalid", args[0])
			}
			purchaser, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.ShieldStaking(
				cmd.Context(),
				&types.QueryShieldStakingRequest{PoolId: poolID, Purchaser: purchaser.String()},
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

// GetCmdShieldStakingRate returns the shield-staking rate for stake-for-shield
func GetCmdShieldStakingRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shield-staking-rate",
		Short: "get shield staking rate for stake-for-shield",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.ShieldStakingRate(cmd.Context(), &types.QueryShieldStakingRateRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdReimbursement returns the command for querying a reimbursement.
func GetCmdReimbursement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reimbursement [proposal ID]",
		Short: "query a reimbursement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool id %s is invalid", args[0])
			}

			res, err := queryClient.Reimbursement(
				cmd.Context(),
				&types.QueryReimbursementRequest{ProposalId: proposalID},
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

// GetCmdReimbursements returns the command for querying reimbursements.
func GetCmdReimbursements() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reimbursements",
		Short: "query all reimbursements",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Reimbursements(cmd.Context(), &types.QueryReimbursementsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
