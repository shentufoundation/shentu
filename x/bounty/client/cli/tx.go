package cli

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		NewWithdrawalFindingCmd(),
		NewReactivateFindingCmd(),
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

			msg, err := types.NewMsgSubmitFinding(
				submitAddr.String(),
				title,
				desc,
				pid,
				severityLevel,
				poc,
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagFindingDesc, "", "The finding description")
	cmd.Flags().String(FlagFindingTitle, "", "The finding's title")
	cmd.Flags().String(FlagFindingPoc, "", "The finding's poc")
	cmd.Flags().Uint64(FlagProgramID, 0, "The program's ID")
	cmd.Flags().Int32(FlagFindingSeverityLevel, 0, "The finding's severity level")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagProgramID)

	return cmd
}

func NewWithdrawalFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawal-finding",
		Short: "withdrawal the specific finding",
		RunE:  setFindingActiveStatus(false),
	}

	cmd.Flags().Uint64(FlagFindingID, 0, "The finding's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagFindingID)
	return cmd
}

func NewReactivateFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reactivate-finding",
		Short: "reactivate the specific finding",
		RunE:  setFindingActiveStatus(true),
	}

	cmd.Flags().Uint64(FlagFindingID, 0, "The finding's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagFindingID)
	return cmd
}

func setFindingActiveStatus(active bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}
		fromAddr := clientCtx.GetFromAddress()
		fid, err := cmd.Flags().GetUint64(FlagFindingID)
		if err != nil {
			return err
		}
		var msg sdk.Msg
		if active {
			// Reactivate
			msg = types.NewMsgReactivateFinding(fromAddr, fid)
		} else {
			// Withdrawal
			msg = types.NewMsgWithdrawalFinding(fromAddr, fid)
		}
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	}
}
