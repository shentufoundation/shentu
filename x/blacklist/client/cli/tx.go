package cli

import (
	"fmt"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

// NewTxCmd returns the root command for the blacklist module's transactions.
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Blacklist transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewBanAddressCmd(),
		NewUnbanAddressCmd(),
	)

	return cmd
}

// NewBanAddressCmd returns a CLI command for building a MsgBanAddress
// transaction. The message must be submitted via a governance proposal; use
// --generate-only to produce the message JSON.
func NewBanAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ban-address [address] --category [compromised|suspected]",
		Short: "Build a message to add an address to the blacklist (submit via governance proposal)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			category, err := parseCategory(cmd.Flag(flagCategory).Value.String())
			if err != nil {
				return err
			}

			msg := &types.MsgBanAddress{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Address:   args[0],
				Category:  category,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return printMsgJSON(clientCtx, msg)
		},
	}

	cmd.Flags().String(flagCategory, "", "entry category (compromised or suspected)")
	_ = cmd.MarkFlagRequired(flagCategory)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnbanAddressCmd returns a CLI command for building a MsgUnbanAddress
// transaction.
func NewUnbanAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unban-address [address] --category [compromised|suspected]",
		Short: "Build a message to remove an address from the blacklist (submit via governance proposal)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUnbanAddress{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Address:   args[0],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return printMsgJSON(clientCtx, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

const (
	flagCategory = "category"
)

// printMsgJSON emits the message as typed proto JSON (with @type), suitable
// for pasting into a governance proposal's messages array.
func printMsgJSON(clientCtx client.Context, msg sdk.Msg) error {
	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	bz, err := clientCtx.Codec.MarshalJSON(anyMsg)
	if err != nil {
		return err
	}
	return clientCtx.PrintString(fmt.Sprintf("%s\n", bz))
}

// parseCategory maps a friendly category name to the proto enum value.
func parseCategory(name string) (types.Category, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "compromised":
		return types.CATEGORY_COMPROMISED, nil
	case "suspected", "suspected-leak":
		return types.CATEGORY_SUSPECTED, nil
	default:
		return types.CATEGORY_UNSPECIFIED, fmt.Errorf("invalid category %q (valid: compromised, suspected)", name)
	}
}
