package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

const (
	FlagDescription   = "description"
	FlagTxhash        = "txhash"
	FlagWait          = "wait"
	FlagName          = "name"
	FlagValidDuration = "valid"
)

var FlagForce bool

// NewTxCmd returns the transaction commands for this module.
func NewTxCmd() *cobra.Command {
	oracleTxCmds := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Oracle staking subcommands",
	}

	oracleTxCmds.AddCommand(
		GetCmdCreateOperator(),
		GetCmdRemoveOperator(),
		GetCmdDepositCollateral(),
		GetCmdWithdrawCollateral(),
		GetCmdClaimReward(),
		GetCmdCreateTask(),
		GetCmdRespondToTask(),
		GetCmdDeleteTask(),
	)

	return oracleTxCmds
}

// GetCmdCreateOperator returns command to create on operator.
func GetCmdCreateOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-operator <address> <collateral>",
		Short: "Create an operator and deposit collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			collateral, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			name := viper.GetString(FlagName)
			msg := types.NewMsgCreateOperator(address, collateral, from, name)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagName, "", "name of the operator")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdRemoveOperator returns command to remove an operator.
func GetCmdRemoveOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-operator <address>",
		Short: "Remove an operator and withdraw collateral & rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgRemoveOperator(address, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDepositCollateral returns command to increase an operator's collateral.
func GetCmdDepositCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-collateral <address> <amount>",
		Short: "Increase an operator's collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgAddCollateral(address, coins)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdWithdrawCollateral returns command to reduce an operator's collateral.
func GetCmdWithdrawCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-collateral <address> <amount>",
		Short: "Reduce an operator's collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgReduceCollateral(address, coins)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdClaimReward returns command to claim (withdraw) an operator's accumulated rewards.
func GetCmdClaimReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-reward <address>",
		Short: "Withdraw all of an operator's accumulated rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgWithdrawReward(address)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCreateTask returns command to create a task.
func GetCmdCreateTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-task <contract_address> <function> <bounty>",
		Short: "Create a task",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			bounty, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}
			if !bounty[0].Amount.IsPositive() {
				return fmt.Errorf("bounty amount is required to be positive")
			}

			// Optional flags
			description := viper.GetString(FlagDescription)
			wait := viper.GetInt64(FlagWait)
			hours := viper.GetInt64(FlagValidDuration)
			validDuration := time.Duration(hours) * time.Hour

			msg := types.NewMsgCreateTask(args[0], args[1], bounty, description, from, wait, validDuration)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagDescription, "", "description of the task")
	cmd.Flags().String(FlagWait, "0", "number of blocks between task creation and aggregation")
	cmd.Flags().String(FlagValidDuration, "0", "valid duration of the task result")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRespondToTask returns command to respond to a task.
func GetCmdRespondToTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond-to-task <contract_address> <function> <score>",
		Short: "Respond to a task",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			score, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				panic(err)
			}

			msg := types.NewMsgTaskResponse(args[0], args[1], score, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDeleteTask returns a delete-task command.
func GetCmdDeleteTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-task <contract_address> <function>",
		Short: "delete a finished task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			force := FlagForce

			msg := types.NewMsgDeleteTask(args[0], args[1], force, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().BoolVarP(&FlagForce, "force", "f", false, "force delete")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
