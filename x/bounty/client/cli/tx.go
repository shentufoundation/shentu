package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// NewTxCmd returns the transaction commands for the certification module.
func NewTxCmd() *cobra.Command {
	bountyTxCmds := &cobra.Command{
		Use:   "bounty",
		Short: "Bounty transactions subcommands",
	}

	bountyTxCmds.AddCommand(
		NewCreateProgramCmd(),
		NewEditProgramCmd(),
		NewActivateProgramCmd(),
		NewCloseProgramCmd(),
		NewSubmitFindingCmd(),
		NewEditFindingCmd(),
		NewActivateFindingCmd(),
		NewConfirmFindingCmd(),
		NewConfirmFindingPaidCmd(),
		NewCloseFindingCmd(),
		NewPublishFindingCmd(),
	)

	return bountyTxCmds
}

func NewCreateProgramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-program",
		Short: "create new program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creatorAddr := clientCtx.GetFromAddress()
			pid, err := cmd.Flags().GetString(FlagProgramID)
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			detail, err := cmd.Flags().GetString(FlagDetail)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateProgram(pid, name, detail, creatorAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagProgramID, "", "The program's id")
	cmd.Flags().String(FlagName, "", "The program's name")
	cmd.Flags().String(FlagDetail, "", "The program's detail")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagProgramID)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagDetail)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewEditProgramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-program",
		Short: "edit a program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creatorAddr := clientCtx.GetFromAddress()
			pid, err := cmd.Flags().GetString(FlagProgramID)
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			detail, err := cmd.Flags().GetString(FlagDetail)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditProgram(pid, name, detail, creatorAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagProgramID, "", "The program's id")
	cmd.Flags().String(FlagName, "", "The program's name")
	cmd.Flags().String(FlagDetail, "", "The program's detail")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagProgramID)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewActivateProgramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "activate the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgActivateProgram(args[0], fromAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCloseProgramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "close the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgCloseProgram(args[0], fromAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewSubmitFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-finding",
		Short: "submit finding for a program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()

			pid, err := cmd.Flags().GetString(FlagProgramID)
			if err != nil {
				return err
			}
			fid, err := cmd.Flags().GetString(FlagFindingID)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(FlagFindingTitle)
			if err != nil {
				return err
			}
			detail, err := cmd.Flags().GetString(FlagDetail)
			if err != nil {
				return err
			}
			severityLevel, err := cmd.Flags().GetString(FlagFindingSeverityLevel)
			if err != nil {
				return err
			}
			byteSeverityLevel, err := types.SeverityLevelFromString(types.NormalizeSeverityLevel(severityLevel))
			if err != nil {
				return err
			}

			desc, err := cmd.Flags().GetString(FlagFindingDescription)
			if err != nil {
				return err
			}
			poc, err := cmd.Flags().GetString(FlagFindingProofOfContent)
			if err != nil {
				return err
			}
			hash := sha256.Sum256([]byte(desc + poc + submitAddr.String()))

			msg := types.NewMsgSubmitFinding(pid, fid, title, detail, hex.EncodeToString(hash[:]), submitAddr, byteSeverityLevel)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagProgramID, "", "The program's ID")
	cmd.Flags().String(FlagFindingID, "", "The finding's ID")
	cmd.Flags().String(FlagFindingTitle, "", "The finding's title")
	cmd.Flags().String(FlagFindingDescription, "", "The finding's description")
	cmd.Flags().String(FlagFindingProofOfContent, "", "The finding's proof of content")
	cmd.Flags().String(FlagDetail, "", "The finding's detail")
	cmd.Flags().String(FlagFindingSeverityLevel, "unspecified", "The finding's severity level")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagProgramID)
	_ = cmd.MarkFlagRequired(FlagFindingID)
	_ = cmd.MarkFlagRequired(FlagFindingTitle)
	_ = cmd.MarkFlagRequired(FlagFindingDescription)
	_ = cmd.MarkFlagRequired(FlagFindingProofOfContent)

	return cmd
}

func NewEditFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-finding",
		Short: "edit finding for a program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()

			fid, err := cmd.Flags().GetString(FlagFindingID)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(FlagFindingTitle)
			if err != nil {
				return err
			}
			detail, err := cmd.Flags().GetString(FlagDetail)
			if err != nil {
				return err
			}
			severityLevel, err := cmd.Flags().GetString(FlagFindingSeverityLevel)
			if err != nil {
				return err
			}
			byteSeverityLevel, err := types.SeverityLevelFromString(types.NormalizeSeverityLevel(severityLevel))
			if err != nil {
				return err
			}

			desc, err := cmd.Flags().GetString(FlagFindingDescription)
			if err != nil {
				return err
			}
			poc, err := cmd.Flags().GetString(FlagFindingProofOfContent)
			if err != nil {
				return err
			}
			if len(desc) > 0 {
				if len(poc) == 0 {
					return errors.New("poc is empty")
				}
			}
			if len(poc) > 0 {
				if len(desc) == 0 {
					return errors.New("desc is empty")
				}
			}

			hashString := ""
			if len(desc) > 0 && len(poc) > 0 {
				hash := sha256.Sum256([]byte(desc + poc + submitAddr.String()))
				hashString = hex.EncodeToString(hash[:])
			}

			paymentHash, err := cmd.Flags().GetString(FlagFindingPaymentHash)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditFinding(fid, title, detail, hashString, paymentHash, submitAddr, byteSeverityLevel)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagFindingID, "", "The finding's ID")
	cmd.Flags().String(FlagFindingTitle, "", "The finding's title")
	cmd.Flags().String(FlagFindingDescription, "", "The finding's description")
	cmd.Flags().String(FlagFindingProofOfContent, "", "The finding's proof of content")
	cmd.Flags().String(FlagDetail, "", "The finding's detail")
	cmd.Flags().String(FlagFindingSeverityLevel, "unspecified", "The finding's severity level")
	cmd.Flags().String(FlagFindingPaymentHash, "", "The finding's payment hash")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagFindingID)

	return cmd
}

func NewActivateFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-finding [finding id]",
		Args:  cobra.ExactArgs(1),
		Short: "activate the specific finding",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()
			msg := types.NewMsgActivateFinding(args[0], submitAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func NewConfirmFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-finding [finding id]",
		Args:  cobra.ExactArgs(1),
		Short: "confirm the specific finding",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()
			fingerprint, err := cmd.Flags().GetString(FlagFindingFingerprint)
			if err != nil {
				return err
			}
			msg := types.NewMsgConfirmFinding(args[0], fingerprint, submitAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagFindingFingerprint, "", "The finding's fingerprint")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagFindingFingerprint)

	return cmd
}

func NewConfirmFindingPaidCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-finding-paid [finding id]",
		Args:  cobra.ExactArgs(1),
		Short: "confirm the specific finding paid",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()
			msg := types.NewMsgConfirmFindingPaid(args[0], submitAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCloseFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-finding [finding id]",
		Args:  cobra.ExactArgs(1),
		Short: "close the specific finding",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()
			msg := types.NewMsgCloseFinding(args[0], submitAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func NewPublishFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish-finding [finding id]",
		Args:  cobra.ExactArgs(1),
		Short: "publish the specific finding",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()
			desc, err := cmd.Flags().GetString(FlagFindingDescription)
			if err != nil {
				return err
			}
			poc, err := cmd.Flags().GetString(FlagFindingProofOfContent)
			if err != nil {
				return err
			}
			msg := types.NewMsgPublishFinding(args[0], desc, poc, submitAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagFindingDescription, "", "The finding's description")
	cmd.Flags().String(FlagFindingProofOfContent, "", "The finding's poc")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagFindingDescription)
	_ = cmd.MarkFlagRequired(FlagFindingProofOfContent)

	return cmd
}
