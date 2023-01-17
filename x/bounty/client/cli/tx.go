package cli

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		NewSubmitFindingCmd(),
		NewHostAcceptFindingCmd(),
		NewHostRejectFindingCmd(),
	)

	return bountyTxCmds
}

func NewCreateProgramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-program",
		Short: "create new program initialized with an initial deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			creatorAddr := clientCtx.GetFromAddress()

			desc, _ := cmd.Flags().GetString(FlagDesc)

			encKeyFile, _ := cmd.Flags().GetString(FlagEncKeyFile)
			var encKey []byte
			if encKeyFile == "" {
				decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
				if err != nil {
					return fmt.Errorf("internal error, failed to generate key")
				}
				encKey = crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)

				// TODO: avoid overwriting silently
				SaveKey(decKey, clientCtx.HomeDir)
			} else {
				encKey = LoadPubKey(encKeyFile)
			}

			newRate := sdk.ZeroDec()
			commissionRate, _ := cmd.Flags().GetString(FlagCommissionRate)
			if commissionRate != "" {
				rate, err := sdk.NewDecFromStr(commissionRate)
				if err != nil {
					return fmt.Errorf("invalid new commission rate: %v", err)
				}
				newRate = rate
			}

			depositStr, _ := cmd.Flags().GetString(FlagDeposit)
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			var sET, jET, cET time.Time

			submissionEndTimeStr, _ := cmd.Flags().GetString(FlagSubmissionEndTime)
			if submissionEndTimeStr != "" {
				sET, err = time.Parse(dateLayout, submissionEndTimeStr)
				if err != nil {
					return err
				}
			}

			judgingEndTimeStr, _ := cmd.Flags().GetString(FlagSubmissionEndTime)
			if judgingEndTimeStr != "" {
				sET, err = time.Parse(dateLayout, judgingEndTimeStr)
				if err != nil {
					return err
				}
			}

			claimEndTimeStr, _ := cmd.Flags().GetString(FlagSubmissionEndTime)
			if claimEndTimeStr != "" {
				sET, err = time.Parse(dateLayout, claimEndTimeStr)
				if err != nil {
					return err
				}
			}

			msg, err := types.NewMsgCreateProgram(
				creatorAddr.String(),
				desc,
				encKey,
				newRate,
				deposit,
				sET,
				jET,
				cET,
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagEncKeyFile, "", "The program's encryption key file to encrypt findings")
	cmd.Flags().String(FlagDesc, "", "The program description.")
	cmd.Flags().String(FlagCommissionRate, "", "The commission rate for the program")
	cmd.Flags().String(FlagDeposit, "", "The initial deposit to the program")
	cmd.Flags().String(FlagSubmissionEndTime, "", "The program's submission end time")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagDesc)

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

			desc, _ := cmd.Flags().GetString(FlagFindingDesc)
			title, _ := cmd.Flags().GetString(FlagFindingTitle)

			pid, err := cmd.Flags().GetUint64(FlagProgramID)
			if err != nil {
				return err
			}
			severityLevel, _ := cmd.Flags().GetInt32(FlagFindingSeverityLevel)
			poc, _ := cmd.Flags().GetString(FlagFindingPoc)

			msg := types.NewMsgSubmitFinding(
				submitAddr.String(),
				title,
				desc,
				pid,
				severityLevel,
				poc,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagFindingDesc, "", "The finding's description")
	cmd.Flags().String(FlagFindingTitle, "", "The finding's title")
	cmd.Flags().String(FlagFindingPoc, "", "Ths finding's poc")
	cmd.Flags().Uint64(FlagProgramID, 0, "The program's ID")
	cmd.Flags().Int32(FlagFindingSeverityLevel, 0, "The finding's severity level")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagProgramID)

	return cmd
}

// NewHostAcceptFindingCmd implements accept a finding by host.
func NewHostAcceptFindingCmd() *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// validate that the finding id is uint
			findingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("finding-id %s not a valid uint, please input a valid finding-id", args[0])
			}
			// Get host address
			hostAddr := clientCtx.GetFromAddress()
			comment, err := cmd.Flags().GetString(FlagComment)
			if err != nil {
				return err
			}

			msg := types.NewMsgHostAcceptFinding(findingID, comment, hostAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagComment, "", "Host's comment on finding")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// NewHostRejectFindingCmd implements reject a finding by host.
func NewHostRejectFindingCmd() *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// validate that the finding id is uint
			findingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("finding-id %s not a valid uint, please input a valid finding-id", args[0])
			}
			// Get host address
			hostAddr := clientCtx.GetFromAddress()
			comment, err := cmd.Flags().GetString(FlagComment)
			if err != nil {
				return err
			}

			msg := types.NewMsgHostRejectFinding(findingID, comment, hostAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagComment, "", "Host's comment on finding")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
