package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

const (
	FlagContract      = "contract"
	FlagFunction      = "function"
	FlagBounty        = "bounty"
	FlagDescription   = "description"
	FlagScore         = "score"
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
		GetCmdInquiry(),
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
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
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
			collateral, err := sdk.ParseCoins(args[1])
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

	return cmd
}

// GetCmdRemoveOperator returns command to remove an operator.
func GetCmdRemoveOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-operator <address>",
		Short: "Remove an operator and withdraw collateral & rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
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
	return cmd
}

// GetCmdDepositCollateral returns command to increase an operator's collateral.
func GetCmdDepositCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-collateral <address> <amount>",
		Short: "Increase an operator's collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			coins, err := sdk.ParseCoins(args[1])
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
	return cmd
}

// GetCmdWithdrawCollateral returns command to reduce an operator's collateral.
func GetCmdWithdrawCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-collateral <address> <amount>",
		Short: "Reduce an operator's collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			coins, err := sdk.ParseCoins(args[1])
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
	return cmd
}

// GetCmdClaimReward returns command to claim (withdraw) an operator's accumulated rewards.
func GetCmdClaimReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-reward <address>",
		Short: "Withdraw all of an operator's accumulated rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
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
	return cmd
}

// GetCmdCreateTask returns command to create a task.
func GetCmdCreateTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-task <flags>",
		Short: "Create a task",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			// Required flags
			contract := viper.GetString(FlagContract)
			if contract == "" {
				return fmt.Errorf("contract address is required to submit a task")
			}
			function := viper.GetString(FlagFunction)
			if function == "" {
				return fmt.Errorf("function is required to submit a task")
			}
			bountyStr := viper.GetString(FlagBounty)
			if bountyStr == "" {
				return fmt.Errorf("bounty is required to submit a task")
			}
			bounty, err := sdk.ParseCoins(bountyStr)
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

			msg := types.NewMsgCreateTask(contract, function, bounty, description, from, wait, validDuration)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagContract, "", "target contract address")
	cmd.Flags().String(FlagFunction, "", "target function")
	cmd.Flags().String(FlagBounty, "", "bounty for operators working on the task")
	cmd.Flags().String(FlagDescription, "", "description of the task")
	cmd.Flags().String(FlagWait, "0", "number of blocks between task creation and aggregation")
	cmd.Flags().String(FlagValidDuration, "0", "valid duration of the task result")

	return cmd
}

// GetCmdRespondToTask returns command to respond to a task.
func GetCmdRespondToTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond-to-task <flags>",
		Short: "Respond to a task",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			contract := viper.GetString(FlagContract)
			if contract == "" {
				return fmt.Errorf("contract address is required to respond to a task")
			}
			function := viper.GetString(FlagFunction)
			if function == "" {
				return fmt.Errorf("function is required to respond to a task")
			}
			scoreStr := viper.GetString(FlagScore)
			if scoreStr == "" {
				return fmt.Errorf("score is required to respond to a task")
			}
			score := viper.GetInt64(FlagScore)

			msg := types.NewMsgTaskResponse(contract, function, score, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagContract, "", "contract address")
	cmd.Flags().String(FlagFunction, "", "function")
	cmd.Flags().String(FlagScore, "", "score")

	return cmd
}

// GetCmdInquiry returns a inquiry-task command.
func GetCmdInquiry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inquiry-task <flags>",
		Short: "Inquiry a task",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			contract := viper.GetString(FlagContract)
			if contract == "" {
				return fmt.Errorf("contract address is required to inquiry a task")
			}

			function := viper.GetString(FlagFunction)
			if function == "" {
				return fmt.Errorf("function is required to inquiry a task")
			}

			txhash := viper.GetString(FlagTxhash)
			if txhash == "" {
				return fmt.Errorf("txhash is required to inquiry a task")
			}

			msg := types.NewMsgInquiryTask(contract, function, txhash, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}
	cmd.Flags().String(FlagContract, "", "contract address")
	cmd.Flags().String(FlagFunction, "", "function")
	cmd.Flags().String(FlagTxhash, "", "txhash")
	return cmd
}

// GetCmdDeleteTask returns a delete-task command.
func GetCmdDeleteTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-task <flags>",
		Short: "delete a finished task",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			contract := viper.GetString(FlagContract)
			if contract == "" {
				return fmt.Errorf("contract address is required to delete a task")
			}

			function := viper.GetString(FlagFunction)
			if function == "" {
				return fmt.Errorf("function is required to delete a task")
			}
			force := FlagForce
			msg := types.NewMsgDeleteTask(contract, function, force, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}
	cmd.Flags().String(FlagContract, "", "contract address")
	cmd.Flags().String(FlagFunction, "", "function")
	cmd.Flags().BoolVarP(&FlagForce, "force", "f", false, "compulsory delete")
	return cmd
}
