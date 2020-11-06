package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	shieldQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the shield module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	shieldQueryCmd.AddCommand(flags.GetCommands(
		GetCmdPool(queryRoute, cdc),
		GetCmdPools(queryRoute, cdc),
		GetCmdPurchaseList(queryRoute, cdc),
		GetCmdPurchaserPurchases(queryRoute, cdc),
		GetCmdPoolPurchases(queryRoute, cdc),
		GetCmdPurchases(queryRoute, cdc),
		GetCmdProvider(queryRoute, cdc),
		GetCmdProviders(queryRoute, cdc),
		GetCmdPoolParams(queryRoute, cdc),
		GetCmdClaimParams(queryRoute, cdc),
		GetCmdStatus(queryRoute, cdc),
		GetCmdStaking(queryRoute, cdc),
		GetCmdShieldStakingRate(queryRoute, cdc),
		GetCmdReimbursement(queryRoute, cdc),
		GetCmdReimbursements(queryRoute, cdc),
	)...)

	return shieldQueryCmd
}

// GetCmdPool returns the command for querying the pool.
func GetCmdPool(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool_ID]",
		Short: "query a pool",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var res []byte
			var err error
			if len(args) == 1 {
				route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryPoolByID, args[0])
				res, _, err = cliCtx.QueryWithData(route, nil)
				if err != nil {
					return err
				}
			} else {
				sponsor := viper.GetString(flagSponsor)
				if sponsor == "" {
					return fmt.Errorf("either poolID or sponsor is required to query pool")
				}

				route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryPoolBySponsor, sponsor)
				res, _, err = cliCtx.QueryWithData(route, nil)
				if err != nil {
					return err
				}
			}

			var out types.Pool
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().String(flagSponsor, "", "use sponsor to query the pool info")

	return cmd
}

// GetCmdPools returns the command for querying a complete list of pools.
func GetCmdPools(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "query a complete list of pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPools), nil)
			if err != nil {
				return err
			}

			var out []types.Pool
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}

	return cmd
}

// GetCmdPurchaseList returns the command for querying purchases
// corresponding to a given pool-purchaser pair.
func GetCmdPurchaseList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-purchaser [pool_ID] [purchaser_address]",
		Short: "get purchases corresponding to a given pool-purchaser pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryPurchaseList, args[0], args[1])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.PurchaseList
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdPurchaserPurchases returns the command for querying
// purchases by a given address.
func GetCmdPurchaserPurchases(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchases-by [purchaser_address]",
		Short: "query purchase information of a given account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryPurchaserPurchases, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out []types.PurchaseList
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}

	return cmd
}

// GetCmdPoolPurchases returns the command for querying
// purchases in a given pool.
func GetCmdPoolPurchases(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-purchases [pool_ID]",
		Short: "query purchases in a given pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryPoolPurchases, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out []types.PurchaseList
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}

	return cmd
}

// GetCmdPurchases returns the command for querying all purchases.
func GetCmdPurchases(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchases",
		Short: "query all purchases",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPurchases)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out []types.Purchase
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}

	return cmd
}

// GetCmdProvider returns the command for querying a provider.
func GetCmdProvider(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [provider_address]",
		Short: "get provider information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryProvider, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.Provider
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}

	return cmd
}

// GetCmdProviders returns the command for querying all providers.
func GetCmdProviders(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Args:  cobra.ExactArgs(0),
		Short: "query all providers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query providers with pagination parameters

Example:
$ %[1]s query shield providers
$ %[1]s query shield providers --page=2 --limit=100
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(flags.FlagPage)
			limit := viper.GetInt(flags.FlagLimit)

			params := types.NewQueryPaginationParams(page, limit)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryProviders)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var out []types.Provider
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().Int(flags.FlagPage, 1, "pagination page of providers to to query for")
	cmd.Flags().Int(flags.FlagLimit, 100, "pagination limit of providers to query for")
	return cmd
}

// GetCmdPoolParams returns the command for querying pool parameters.
func GetCmdPoolParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-params",
		Short: "get pool parameters",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPoolParams)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.PoolParams
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdClaimParams returns the command for querying claim parameters.
func GetCmdClaimParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-params",
		Short: "get claim parameters",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryClaimParams)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.ClaimProposalParams
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdStatus returns the command for querying shield status.
func GetCmdStatus(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "get shield status",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryStatus)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.QueryResStatus
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdStaking returns the command for querying staked-for-shield amounts
// corresponding to a given pool-purchaser pair.
func GetCmdStaking(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staked-for-shield [pool_ID] [purchaser_address]",
		Short: "get staked CTK for shield corresponding to a given pool-purchaser pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryStakedForShield, args[0], args[1])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.ShieldStaking
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdShieldStakingRate returns the shield-staking rate for stake-for-shield
func GetCmdShieldStakingRate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shield-staking-rate",
		Short: "get shield staking rate for stake-for-shield",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryShieldStakingRate)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out sdk.Dec
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	return cmd
}

// GetCmdReimbursement returns the command for querying a reimbursement.
func GetCmdReimbursement(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reimbursement [proposal ID]",
		Short: "query a reimbursement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryReimbursement, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.Reimbursement
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}

	return cmd
}

// GetCmdReimbursements returns the command for querying reimbursements.
func GetCmdReimbursements(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reimbursements",
		Short: "query all reimbursements",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryReimbursements), nil)
			if err != nil {
				return err
			}

			var out []types.ProposalIDReimbursementPair
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}

	return cmd
}
