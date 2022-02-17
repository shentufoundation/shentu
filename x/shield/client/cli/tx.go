package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

var (
	flagDescription = "description"
	flagShieldRate  = "shield-rate"
	flagActive      = "active"
)

// NewTxCmd returns the transaction commands for this module.
func NewTxCmd() *cobra.Command {
	shieldTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Shield transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	shieldTxCmd.AddCommand(
		GetCmdCreatePool(),
		GetCmdUpdatePool(),
		GetCmdDepositCollateral(),
		GetCmdWithdrawCollateral(),
		GetCmdWithdrawRewards(),
		GetCmdWithdrawForeignRewards(),
		GetCmdPurchaseShield(),
		GetCmdUpdateSponsor(),
		GetCmdUnstake(),
		GetCmdDoante(),
	)

	return shieldTxCmd
}

// GetCmdSubmitProposal implements the command for submitting a Shield claim proposal.
func GetCmdSubmitProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shield-claim [proposal file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a Shield claim proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a Shield claim proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal shield-claim <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "pool_id": 1,
  "loss": [
    {
      "denom": "ctk",
      "amount": "1000"
    }
  ],
  "evidence": "Attack happened on <time> caused loss of <amount> to <account> by <txhashes>",
  "purchase_txhash": "7D5C90FBD3082D2CD763FA1580BBA29568D0749D76C7CD627B841F2FAB22BBEA",
  "description": "Details of the attack",
  "deposit": [
    {
      "denom": "ctk",
      "amount": "100"
    }
  ]
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			proposal, err := ParseShieldClaimProposalJSON(args[0])
			if err != nil {
				return err
			}
			from := cliCtx.GetFromAddress()
			content := types.NewShieldClaimProposal(proposal.PoolID, proposal.Loss,
				proposal.Evidence, proposal.Description, from)

			msg, err := govtypes.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	return cmd
}

// GetCmdCreatePool implements the command for creating a Shield pool.
func GetCmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [shield amount] [sponsor] [sponsor-address]",
		Args:  cobra.ExactArgs(3),
		Short: "create new Shield pool initialized with an validator address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a Shield pool. Can only be executed from the Shield admin address.

Example:
$ %s tx shield create-pool <shield amount> <sponsor> <sponsor-address> --shield-rate <shield rate>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			sponsorAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			description, _ := cmd.Flags().GetString(flagDescription)
			flagShieldRateExtract, _ := cmd.Flags().GetString(flagShieldRate)
			shieldRate, err := sdk.NewDecFromStr(flagShieldRateExtract)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePool(fromAddr, sponsorAddr, description, shieldRate)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(flagDescription, "", "description for the pool")
	cmd.Flags().String(flagShieldRate, "", "Shield Rate")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdUpdatePool implements the command for updating an existing Shield pool.
func GetCmdUpdatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pool [pool id]",
		Args:  cobra.ExactArgs(1),
		Short: "update an existing Shield pool by adding more deposit or updating Shield amount.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update a Shield pool. Can only be executed from the Shield admin address.

Example:
$ %s tx shield update-pool <id> --native-deposit <ctk deposit> --shield <shield amount> --shield-rate <shield rate>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			description, _ := cmd.Flags().GetString(flagDescription)
			flagShieldRateExtract, _ := cmd.Flags().GetString(flagShieldRate)
			var shieldRate sdk.Dec
			if shieldRateInput := flagShieldRateExtract; shieldRateInput != "" {
				shieldRate, err = sdk.NewDecFromStr(shieldRateInput)
				if err != nil {
					return err
				}
			}
			active, err := cmd.Flags().GetBool(flagActive)
			if err != nil {
				panic(err)
			}

			msg := types.NewMsgUpdatePool(fromAddr, id, description, active, shieldRate)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(flagDescription, "", "description for the pool")
	cmd.Flags().String(flagShieldRate, "", "Shield Rate")
	cmd.Flags().Bool(flagActive, true, "new pool status. default true.")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDepositCollateral implements command for community member to
// join a pool by depositing collateral.
func GetCmdDepositCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-collateral [collateral]",
		Short: "join a Shield pool as a community member by depositing collateral",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			collateral, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositCollateral(fromAddr, collateral)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdWithdrawCollateral implements command for community member to
// withdraw deposited collateral from Shield pool.
func GetCmdWithdrawCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-collateral [collateral]",
		Short: "withdraw deposited collateral from Shield pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			collateral, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawCollateral(fromAddr, collateral)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdWithdrawRewards implements command for requesting to withdraw native tokens rewards.
func GetCmdWithdrawRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-rewards",
		Short: "withdraw CTK rewards",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgWithdrawRewards(fromAddr)

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdWithdrawForeignRewards implements command for requesting to withdraw foreign tokens rewards.
func GetCmdWithdrawForeignRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-foreign-rewards [denom] [address]",
		Short: "withdraw foreign rewards coins to their original chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()
			denom := args[0]
			addr := args[1]

			msg := types.NewMsgWithdrawForeignRewards(fromAddr, denom, addr)

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdPurchaseShield implements the command for purchasing Shield.
func GetCmdPurchaseShield() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase [pool id] [shield amount] [description]",
		Args:  cobra.ExactArgs(3),
		Short: "purchase Shield",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Purchase Shield. Requires purchaser to provide descriptions of accounts to be protected.

Example:
$ %s tx shield purchase <pool id> <shield amount> <description>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			shield, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			description := args[2]
			if description == "" {
				return types.ErrPurchaseMissingDescription
			}

			msg := types.NewMsgPurchase(poolID, shield, description, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdUnstake implements the command for purchasing Shield.
func GetCmdUnstake() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unstake-from-shield [pool id] [amount] ",
		Args:  cobra.ExactArgs(2),
		Short: "unstake staked-for-shield coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw staking from shield. Requires existing shield purchase through staking.

Example:
$ %s tx shield withdraw-staking <pool id> <shield amount> 
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			shield, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnstake(poolID, shield, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdUpdateSponsor implements the command for updating a pool's sponsor.
func GetCmdUpdateSponsor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-sponsor [pool id] [new_sponsor] [new_sponsor_address]",
		Args:  cobra.ExactArgs(3),
		Short: "update the sponsor of an existing pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update a pool's sponsor. Can only be executed from the Shield admin address.
Example:
$ %s tx shield update-sponsor <id> <new_sponsor_name> <new_sponsor_address> --from=<key_or_address>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			sponsorAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateSponsor(poolID, args[1], sponsorAddr, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDoante implements donating to Shield Donation Pool.
func GetCmdDoante() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate [amount]",
		Short: "donate to Shield Donation Pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			donation, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDonate(fromAddr, donation)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
