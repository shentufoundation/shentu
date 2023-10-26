package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

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
		NewOpenProgramCmd(),
		NewCloseProgramCmd(),
		NewSubmitFindingCmd(),
		NewAcceptFindingCmd(),
		NewRejectFindingCmd(),
		NewCloseFindingCmd(),
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
			name, _ := cmd.Flags().GetString(FlagName)
			pid, _ := cmd.Flags().GetString(FlagProgramID)
			desc, _ := cmd.Flags().GetString(FlagDesc)
			members, _ := cmd.Flags().GetStringArray(FlagMembers)

			msg, err := types.NewMsgCreateProgram(name, desc, pid, creatorAddr, members, nil)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetProgramDetail())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagDesc)

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
			name, _ := cmd.Flags().GetString(FlagName)
			pid, _ := cmd.Flags().GetString(FlagProgramID)
			desc, _ := cmd.Flags().GetString(FlagDesc)
			members, _ := cmd.Flags().GetStringArray(FlagMembers)

			msg, err := types.NewMsgEditProgram(name, pid, desc, creatorAddr, members, nil)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetProgramDetail())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagDesc)

	return cmd
}

func NewOpenProgramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "open the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgOpenProgram(fromAddr, args[0])
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
		Short: "end the program",
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
			title, _ := cmd.Flags().GetString(FlagFindingTitle)

			severityLevel, _ := cmd.Flags().GetInt32(FlagFindingSeverityLevel)
			_, ok := types.SeverityLevel_name[severityLevel]
			if !ok {
				return fmt.Errorf("invalid %s value", FlagFindingSeverityLevel)
			}

			desc, _ := cmd.Flags().GetString(FlagDesc)

			msg := types.NewMsgSubmitFinding(pid, fid, title, desc, submitAddr, types.SeverityLevel(severityLevel))
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagDesc, "", "The finding's description")
	cmd.Flags().String(FlagFindingTitle, "", "The finding's title")
	cmd.Flags().String(FlagFindingPoc, "", "Ths finding's poc")
	cmd.Flags().Uint64(FlagProgramID, 0, "The program's ID")
	cmd.Flags().Int32(FlagFindingSeverityLevel, 0, "The finding's severity level")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagProgramID)

	return cmd
}

// NewAcceptFindingCmd implements accept a finding by host.
func NewAcceptFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-finding [finding-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Host accept a finding for a program",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Host accept a finding for a program.Meantime, you can also add some comments, which will be encrypted.
Example:
$ %s tx bounty accept-finding 1 --comment "Looks good to me"
`,
				version.AppName,
			),
		),
		RunE: AcceptFinding,
	}

	cmd.Flags().String(FlagComment, "", "Host's comment on finding")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// NewRejectFindingCmd implements reject a finding by host.
func NewRejectFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-finding [finding-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Host reject a finding for a program",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Host reject a finding for a program.Meantime, you can also add some comments, which will be encrypted.
Example:
$ %s tx bounty reject-finding 1 --comment "Verified to be an invalid finding"
`,
				version.AppName,
			),
		),
		RunE: RejectFinding,
	}

	cmd.Flags().String(FlagComment, "", "Host's comment on finding")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func AcceptFinding(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	// Get host address
	hostAddr := clientCtx.GetFromAddress()
	msg := types.NewMsgAcceptFinding(args[0], hostAddr)

	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}

func RejectFinding(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	// Get host address
	hostAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgRejectFinding(args[0], hostAddr)

	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
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
