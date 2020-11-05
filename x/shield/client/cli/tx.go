package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/certikfoundation/shentu/x/shield/types"
)

var (
	flagNativeDeposit = "native-deposit"
	flagShield        = "shield"
	flagSponsor       = "sponsor"
	flagDescription   = "description"
	flagShieldLimit   = "shield-limit"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	shieldTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Shield transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	shieldTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreatePool(cdc),
		GetCmdUpdatePool(cdc),
		GetCmdPausePool(cdc),
		GetCmdResumePool(cdc),
		GetCmdDepositCollateral(cdc),
		GetCmdWithdrawCollateral(cdc),
		GetCmdWithdrawRewards(cdc),
		GetCmdWithdrawForeignRewards(cdc),
		GetCmdClearPayouts(cdc),
		GetCmdPurchaseShield(cdc),
		GetCmdWithdrawReimbursement(cdc),
		GetCmdUpdateSponsor(cdc),
		GetCmdStakeForShield(cdc),
		GetCmdUnstakeFromShield(cdc),
	)...)

	return shieldTxCmd
}

// GetCmdSubmitProposal implements the command for submitting a Shield claim proposal.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			proposal, err := ParseShieldClaimProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}
			from := cliCtx.GetFromAddress()
			content := types.NewShieldClaimProposal(proposal.PoolID, proposal.Loss,
				proposal.PurchaseID, proposal.Evidence, proposal.Description, from)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdCreatePool implements the command for creating a Shield pool.
func GetCmdCreatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [shield amount] [sponsor] [sponsor-address]",
		Args:  cobra.ExactArgs(3),
		Short: "create new Shield pool initialized with an validator address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a Shield pool. Can only be executed from the Shield admin address.

Example:
$ %s tx shield create-pool <shield amount> <sponsor> <sponsor-address> --native-deposit <ctk deposit> --shield-limit <shield limit>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			shield, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			sponsor := args[1]

			sponsorAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			nativeDeposit, err := sdk.ParseCoins(viper.GetString(flagNativeDeposit))
			if err != nil {
				return err
			}
			deposit := types.MixedCoins{Native: nativeDeposit}

			description := viper.GetString(flagDescription)

			shieldLimit, ok := sdk.NewIntFromString(viper.GetString(flagShieldLimit))
			if !ok {
				return fmt.Errorf("invalid input for shield limit")
			}

			msg := types.NewMsgCreatePool(fromAddr, shield, deposit, sponsor, sponsorAddr, description, shieldLimit)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagDescription, "", "description for the pool")
	cmd.Flags().String(flagNativeDeposit, "", "CTK deposit amount")
	cmd.Flags().String(flagShieldLimit, "", "the limit of active shield for the pool")
	return cmd
}

// GetCmdUpdatePool implements the command for updating an existing Shield pool.
func GetCmdUpdatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pool [pool id]",
		Args:  cobra.ExactArgs(1),
		Short: "update an existing Shield pool by adding more deposit or updating Shield amount.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update a Shield pool. Can only be executed from the Shield admin address.

Example:
$ %s tx shield update-pool <id> --native-deposit <ctk deposit> --shield <shield amount> --shield-limit <shield limit>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			nativeDeposit, err := sdk.ParseCoins(viper.GetString(flagNativeDeposit))
			if err != nil {
				return err
			}

			shield, err := sdk.ParseCoins(viper.GetString(flagShield))
			if err != nil {
				return err
			}
			deposit := types.MixedCoins{Native: nativeDeposit}

			description := viper.GetString(flagDescription)

			shieldLimit, ok := sdk.NewIntFromString(viper.GetString(flagShieldLimit))
			if !ok {
				return fmt.Errorf("invalid input for shield limit")
			}

			msg := types.NewMsgUpdatePool(fromAddr, shield, deposit, id, description, shieldLimit)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagShield, "", "CTK Shield amount")
	cmd.Flags().String(flagNativeDeposit, "", "CTK deposit amount")
	cmd.Flags().String(flagDescription, "", "description for the pool")
	cmd.Flags().String(flagShieldLimit, "", "the limit of active shield for the pool")
	return cmd
}

// GetCmdPausePool implements the command for pausing a pool.
func GetCmdPausePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause-pool [pool id]",
		Args:  cobra.ExactArgs(1),
		Short: "pause a Shield pool to disallow further Shield purchase.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Pause a Shield pool to prevent new Shield purchases. Can only be executed from the Shield admin address.

Example:
$ %s tx shield pause-pool <pool id>
`,
				version.ClientName,
			),
		),
		RunE: pauseOrResume(cdc, false),
	}
	return cmd
}

// GetCmdResumePool implements the command for resuming a pool.
func GetCmdResumePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume-pool [pool id]",
		Args:  cobra.ExactArgs(1),
		Short: "resume a Shield pool to allow Shield purchase.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Resume a Shield pool to reactivate Shield purchase. Can only be executed from the Shield admin address.

Example:
$ %s tx shield resume-pool <pool id>
`,
				version.ClientName,
			),
		),
		RunE: pauseOrResume(cdc, true),
	}
	return cmd
}

func pauseOrResume(cdc *codec.Codec, active bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		inBuf := bufio.NewReader(cmd.InOrStdin())
		txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
		cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

		fromAddr := cliCtx.GetFromAddress()

		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return err
		}

		var msg sdk.Msg
		if active {
			msg = types.NewMsgResumePool(fromAddr, id)
		} else {
			msg = types.NewMsgPausePool(fromAddr, id)
		}
		if err := msg.ValidateBasic(); err != nil {
			return err
		}

		return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
	}
}

// GetCmdDepositCollateral implements command for community member to
// join a pool by depositing collateral.
func GetCmdDepositCollateral(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-collateral [collateral]",
		Short: "join a Shield pool as a community member by depositing collateral",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			collateral, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositCollateral(fromAddr, collateral)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdWithdrawCollateral implements command for community member to
// withdraw deposited collateral from Shield pool.
func GetCmdWithdrawCollateral(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-collateral [collateral]",
		Short: "withdraw deposited collateral from Shield pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			collateral, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawCollateral(fromAddr, collateral)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdWithdrawRewards implements command for requesting to withdraw native tokens rewards.
func GetCmdWithdrawRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-rewards",
		Short: "withdraw CTK rewards",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgWithdrawRewards(fromAddr)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdWithdrawForeignRewards implements command for requesting to withdraw foreign tokens rewards.
func GetCmdWithdrawForeignRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-foreign-rewards [denom] [address]",
		Short: "withdraw foreign rewards coins to their original chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()
			denom := args[0]
			addr := args[1]

			msg := types.NewMsgWithdrawForeignRewards(fromAddr, denom, addr)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdClearPayouts implements command for requesting to clear out pending payouts.
func GetCmdClearPayouts(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear-payouts [denom]",
		Short: "clear pending payouts after they have been distributed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()
			denom := args[0]

			msg := types.NewMsgClearPayouts(fromAddr, denom)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdPurchaseShield implements the command for purchasing Shield.
func GetCmdPurchaseShield(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase [pool id] [shield amount] [description]",
		Args:  cobra.ExactArgs(3),
		Short: "purchase Shield",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Purchase Shield. Requires purchaser to provide descriptions of accounts to be protected.

Example:
$ %s tx shield purchase <pool id> <shield amount> <description>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			shield, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}
			description := args[2]
			if description == "" {
				return types.ErrPurchaseMissingDescription
			}

			msg := types.NewMsgPurchaseShield(poolID, shield, description, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdWithdrawReimbursement the command for withdrawing reimbursement.
func GetCmdWithdrawReimbursement(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-reimbursement [proposal id]",
		Args:  cobra.ExactArgs(1),
		Short: "withdraw reimbursement",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw reimbursement by proposal id.

Example:
$ %s tx shield withdraw-reimbursement <proposal id>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawReimbursement(proposalID, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdStakeForShield implements the command for purchasing Shield.
func GetCmdStakeForShield(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-for-shield [pool id] [shield amount] [description]",
		Args:  cobra.ExactArgs(3),
		Short: "obtain shield through staking CTK",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Obtain shield through staking. Requires purchaser to provide descriptions of accounts to be protected.

Example:
$ %s tx shield stake-for-shield <pool id> <shield amount> <description>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			shield, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}
			description := args[2]
			if description == "" {
				return types.ErrPurchaseMissingDescription
			}

			msg := types.NewMsgStakeForShield(poolID, shield, description, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdUnstakeFromShield implements the command for purchasing Shield.
func GetCmdUnstakeFromShield(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unstake-from-shield [pool id] [amount] ",
		Args:  cobra.ExactArgs(2),
		Short: "unstake staked-for-shield coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw staking from shield. Requires existing shield purchase through staking.

Example:
$ %s tx shield withdraw-staking <pool id> <shield amount> 
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			shield, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnstakeFromShield(poolID, shield, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdUpdateSponsor implements the command for updating a pool's sponsor.
func GetCmdUpdateSponsor(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-sponsor [pool id] [new_sponsor] [new_sponsor_address]",
		Args:  cobra.ExactArgs(3),
		Short: "update the sponsor of an existing pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update a pool's sponsor. Can only be executed from the Shield admin address.
Example:
$ %s tx shield update-sponsor <id> <new_sponsor_name> <new_sponsor_address> --from=<key_or_address>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

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
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
