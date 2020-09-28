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
	flagNativeDeposit    = "native-deposit"
	flagForeignDeposit   = "foreign-deposit"
	flagShield           = "shield"
	flagSponsor          = "sponsor"
	flagTimeOfCoverage   = "time-of-coverage"
	flagBlocksOfCoverage = "blocks-of-coverage"
)

// GetTxCmd returns the transaction commands for this module
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
		GetCmdPurchaseShield(cdc),
	)...)

	return shieldTxCmd
}

// GetCmdSubmitProposal implements the command to submit a shield-claim proposal.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shield-claim [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a shield claim proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a shield claim proposal along with an initial deposit.
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
			content := types.NewShieldClaimProposal(proposal.PoolID, proposal.Loss, proposal.Evidence,
				proposal.PurchaseTxHash, proposal.Description, from, proposal.Deposit)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdCreatePool implements the create pool command handler.
func GetCmdCreatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [shield amount] [sponsor]",
		Args:  cobra.ExactArgs(2),
		Short: "create new Shield pool initialized with an validator address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a shield pool. Can only be executed from the shield operator address.

Example:
$ %s tx shield create-pool <shield amount> <sponsor> --native-deposit <ctk deposit> --foreign-deposit <external deposit>
--time-of-coverage <period in seconds>
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

			nativeDeposit, err := sdk.ParseCoins(viper.GetString(flagNativeDeposit))
			if err != nil {
				return err
			}

			foreignDeposit, err := sdk.ParseCoins(viper.GetString(flagForeignDeposit))
			if err != nil {
				return err
			}

			deposit := types.MixedCoins{
				Native:  nativeDeposit,
				Foreign: foreignDeposit,
			}

			timeOfCoverage := viper.GetInt64(flagTimeOfCoverage)
			blocksOfCoverage := viper.GetInt64(flagBlocksOfCoverage)

			msg := types.NewMsgCreatePool(fromAddr, shield, deposit, sponsor, timeOfCoverage, blocksOfCoverage)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagNativeDeposit, "", "CTK deposit amount")
	cmd.Flags().String(flagForeignDeposit, "", "foreign coins deposit amount")
	cmd.Flags().Int64(flagTimeOfCoverage, 0, "time of coverage")
	cmd.Flags().Int64(flagBlocksOfCoverage, 0, "blocks of coverage")

	return cmd
}

// GetCmdUpdatePool implements the create pool command handler.
func GetCmdUpdatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pool [pool id]",
		Args:  cobra.ExactArgs(1),
		Short: "update new Shield pool through adding more deposit or updating shield amount.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update a shield pool. Can only be executed from the shield operator address.

Example:
$ %s tx shield update-pool <id> --native-deposit <ctk deposit> --foreign-deposit <external deposit> --shield <shield amount> 
--time-of-coverage <additional period>
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

			foreignDeposit, err := sdk.ParseCoins(viper.GetString(flagForeignDeposit))
			if err != nil {
				return err
			}

			shield, err := sdk.ParseCoins(viper.GetString(flagShield))
			if err != nil {
				return err
			}

			deposit := types.MixedCoins{
				Native:  nativeDeposit,
				Foreign: foreignDeposit,
			}

			if deposit.Native == nil && deposit.Foreign == nil && shield == nil {
				return types.ErrNoUpdate
			}

			timeOfCoverage := viper.GetInt64(flagTimeOfCoverage)
			blocksOfCoverage := viper.GetInt64(flagBlocksOfCoverage)

			msg := types.NewMsgUpdatePool(fromAddr, shield, deposit, id, timeOfCoverage, blocksOfCoverage)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagShield, "", "CTK shield amount")
	cmd.Flags().String(flagNativeDeposit, "", "CTK deposit amount")
	cmd.Flags().String(flagForeignDeposit, "", "foreign coins deposit amount")
	cmd.Flags().Int64(flagTimeOfCoverage, 0, "additional time of coverage")
	cmd.Flags().Int64(flagBlocksOfCoverage, 0, "additional blocks of coverage")
	return cmd
}

// GetCmdPausePool implements the pause pool command handler.
func GetCmdPausePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause-pool [pool id]",
		Args:  cobra.ExactArgs(1),
		Short: "pause a Shield pool to disallow further shield purchase.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Pause a shield pool to prevent and new shield purchases for the pool. Can only be executed from the shield operator address.

Example:
$ %s tx shield pause-pool <pool id>
`,
				version.ClientName,
			),
		),
		RunE: PauseOrResume(cdc, false),
	}
	return cmd
}

// GetCmdResumePool implements the resume pool command handler.
func GetCmdResumePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume-pool [pool id]",
		Args:  cobra.ExactArgs(1),
		Short: "resume a Shield pool to allow shield purchase for an existing pool.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Resume a shield pool to reactivate shield purchase. Can only be executed from the shield operator address.

Example:
$ %s tx shield resume-pool <pool id>
`,
				version.ClientName,
			),
		),
		RunE: PauseOrResume(cdc, true),
	}
	return cmd
}

func PauseOrResume(cdc *codec.Codec, active bool) func(cmd *cobra.Command, args []string) error {
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
		Use:   "deposit-collateral [pool id] [collateral]",
		Short: "join a Shield pool as a community member by depositing collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositCollateral(fromAddr, id, collateral)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdPurchaseShield implements the purchase shield command handler.
func GetCmdPurchaseShield(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase [pool id] [shield amount] [description]",
		Args:  cobra.ExactArgs(3),
		Short: "purchase shield",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Purchase shield. Requires purchaser to provide descriptions of accounts to be protected.

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
