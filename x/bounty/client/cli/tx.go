package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"
	"time"

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
			if err != nil {
				return err
			}

			desc, _ := cmd.Flags().GetString(FlagDesc)

			encKeyFile, _ := cmd.Flags().GetString(FlagEncKeyFile)
			var encKey cryptotypes.PubKey
			if encKeyFile == "" {
				decKey := secp256k1.GenPrivKey()
				encKey = decKey.PubKey()

				// TODO: avoid overwriting silently
				SaveKeys(encKey, decKey, clientCtx.HomeDir, clientCtx.Codec)
			} else {
				encKey = LoadPubKey(encKeyFile, clientCtx.Codec)
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