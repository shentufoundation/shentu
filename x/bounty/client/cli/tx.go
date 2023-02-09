package cli

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
		NewCancelFindingCmd(),
		NewReleaseFindingCmd(),
		NewTerminateProgramCmd(),
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
			_, ok := types.SeverityLevel_name[severityLevel]
			if !ok {
				return fmt.Errorf("invalid %s value", FlagFindingSeverityLevel)
			}

			poc, _ := cmd.Flags().GetString(FlagFindingPoc)

			descAny, pocAny, err := EncryptMsg(cmd, pid, desc, poc)
			if err != nil {
				return err
			}
			msg := types.NewMsgSubmitFinding(
				submitAddr.String(),
				title,
				descAny,
				pocAny,
				pid,
				severityLevel,
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

func EncryptMsg(cmd *cobra.Command, programID uint64, desc, poc string) (descAny, pocAny *codectypes.Any, err error) {
	eciesEncKey, err := GetEncryptionKey(cmd, programID)
	if err != nil {
		return nil, nil, err
	}

	encryptedDescBytes, err := ecies.Encrypt(rand.Reader, eciesEncKey, []byte(desc), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	encDesc := types.EciesEncryptedDesc{
		FindingDesc: encryptedDescBytes,
	}
	descAny, err = codectypes.NewAnyWithValue(&encDesc)
	if err != nil {
		return nil, nil, err
	}

	encryptedPocBytes, err := ecies.Encrypt(rand.Reader, eciesEncKey, []byte(poc), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	encPoc := types.EciesEncryptedPoc{
		FindingPoc: encryptedPocBytes,
	}
	pocAny, err = codectypes.NewAnyWithValue(&encPoc)
	if err != nil {
		return nil, nil, err
	}
	return descAny, pocAny, nil
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
		RunE: HostAcceptFinding,
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
		RunE: HostRejectFinding,
	}

	cmd.Flags().String(FlagComment, "", "Host's comment on finding")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func HostAcceptFinding(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	findingID, encryptedCommentAny, hostAddr, err := HostProcessFinding(cmd, args)
	if err != nil {
		return err
	}

	msg := types.NewMsgHostAcceptFinding(findingID, encryptedCommentAny, hostAddr)

	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}

func HostRejectFinding(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	findingID, encryptedCommentAny, hostAddr, err := HostProcessFinding(cmd, args)
	if err != nil {
		return err
	}

	msg := types.NewMsgHostRejectFinding(findingID, encryptedCommentAny, hostAddr)

	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}

func HostProcessFinding(cmd *cobra.Command, args []string) (fid uint64,
	commentAny *codectypes.Any, hostAddr sdk.AccAddress, err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return fid, commentAny, hostAddr, err
	}

	// validate that the finding id is uint
	fid, err = strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fid, commentAny, hostAddr, fmt.Errorf("finding-id %s not a valid uint, please input a valid finding-id", args[0])
	}
	// Get host address
	hostAddr = clientCtx.GetFromAddress()
	comment, err := cmd.Flags().GetString(FlagComment)
	if err != nil {
		return fid, commentAny, hostAddr, err
	}

	// comment is empty and does not need to be encrypted
	if len(comment) == 0 {
		return fid, commentAny, hostAddr, nil
	}

	// get eciesEncKey
	finding, err := GetFinding(cmd, fid)
	if err != nil {
		return fid, commentAny, hostAddr, err
	}
	eciesEncKey, err := GetEncryptionKey(cmd, finding.ProgramId)
	if err != nil {
		return fid, commentAny, hostAddr, err
	}

	encryptedComment, err := ecies.Encrypt(rand.Reader, eciesEncKey, []byte(comment), nil, nil)
	if err != nil {
		return fid, commentAny, hostAddr, err
	}
	encComment := types.EciesEncryptedComment{
		FindingComment: encryptedComment,
	}
	commentAny, err = codectypes.NewAnyWithValue(&encComment)
	if err != nil {
		return fid, commentAny, hostAddr, err
	}

	return fid, commentAny, hostAddr, nil
}

func NewCancelFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-finding [finding id]",
		Args:  cobra.ExactArgs(1),
		Short: "cancel the specific finding",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			submitAddr := clientCtx.GetFromAddress()
			fid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			msg := types.NewMsgCancelFinding(submitAddr, fid)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func NewReleaseFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release-finding [finding-id]",
		Args:  cobra.ExactArgs(1),
		Short: "release encrypted part of a finding ",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			hostAddr := clientCtx.GetFromAddress()

			fid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("finding-id %s not a valid uint, please input a valid finding-id", args[0])
			}

			encKeyFile, err := cmd.Flags().GetString(FlagEncKeyFile)
			if err != nil {
				return err
			}

			findingDesc, findingPoc, findingComment, err := GetFindingPlainText(cmd, fid, encKeyFile)
			if err != nil {
				return err
			}

			msg := types.NewReleaseFinding(
				hostAddr.String(),
				fid,
				findingDesc,
				findingPoc,
				findingComment,
			)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagEncKeyFile, "", "The program's encryption key file to decrypt findings")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagEncKeyFile)

	return cmd
}

func GetFindingPlainText(cmd *cobra.Command, fid uint64, encKeyFile string) (
	desc, poc, comment string, err error) {
	// get finding info
	finding, err := GetFinding(cmd, fid)
	if err != nil {
		return "", "", "", err
	}

	prvKey := LoadPrvKey(encKeyFile)

	if finding.FindingDesc == nil {
		desc = ""
	} else {
		var descProto types.EciesEncryptedDesc
		if err = proto.Unmarshal(finding.FindingDesc.GetValue(), &descProto); err != nil {
			return "", "", "", err
		}
		descBytes, err := prvKey.Decrypt(descProto.FindingDesc, nil, nil)
		if err != nil {
			return "", "", "", err
		}
		desc = string(descBytes)
	}

	if finding.FindingPoc == nil {
		poc = ""
	} else {
		var pocProto types.EciesEncryptedPoc
		if err = proto.Unmarshal(finding.FindingPoc.GetValue(), &pocProto); err != nil {
			return "", "", "", err
		}
		pocBytes, err := prvKey.Decrypt(pocProto.FindingPoc, nil, nil)
		if err != nil {
			return "", "", "", err
		}
		poc = string(pocBytes)
	}

	if finding.FindingComment == nil {
		comment = ""
	} else {
		var commentProto types.EciesEncryptedComment
		if err = proto.Unmarshal(finding.FindingComment.GetValue(), &commentProto); err != nil {
			return "", "", "", err
		}
		commentBytes, err := prvKey.Decrypt(commentProto.FindingComment, nil, nil)
		if err != nil {
			return "", "", "", err
		}
		comment = string(commentBytes)
	}
	return desc, poc, comment, nil
}

func NewTerminateProgramCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "terminate-program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "terminate the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress()
			pid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s is not a valid uint", args[0])
			}
			msg := types.NewMsgTerminateProgram(fromAddr.String(), pid)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}
