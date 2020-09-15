package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/certikfoundation/shentu/app"
	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/auth/vesting"
)

const (
	flagClientHome    = "home-client"
	flagVestingStart  = "vesting-start-time"
	flagVestingEnd    = "vesting-end-time"
	flagVestingAmt    = "vesting-amount"
	flagPeriod        = "period"
	flagNumberPeriods = "num-periods"
	flagContinuous    = "continuous"
	flagTriggered     = "triggered"
	flagManual        = "manual"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
the precedence rule is period > continuous > endtime.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				// attempt to lookup address from Keybase if no address was provided
				kb, err1 := keys.NewKeyring(
					strings.ToLower(app.AppName),
					viper.GetString(flags.FlagKeyringBackend),
					viper.GetString(flagClientHome),
					inBuf,
				)
				if err1 != nil {
					return err1
				}
				info, err2 := kb.Get(args[0])
				if err2 != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err2)
				}
				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			// create concrete account type based on input parameters
			baseAccount := auth.NewBaseAccount(addr, coins.Sort(), nil, 0, 0)
			genAccount, err := getVestedAccountFromFlags(baseAccount, coins)
			if err != nil {
				return err
			}
			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			authGenState := auth.GetGenesisStateFromAppState(cdc, appState)
			if authGenState.Accounts.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}
			// add the new account to the set of genesis accounts
			// and sanitize the accounts afterwards
			authGenState.Accounts = append(authGenState.Accounts, genAccount)
			authGenState.Accounts = auth.SanitizeGenesisAccounts(authGenState.Accounts)
			authGenStateBz, err := cdc.MarshalJSON(authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[auth.ModuleName] = authGenStateBz
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagPeriod, 0, "set to periodic vesting with period in seconds")
	cmd.Flags().Uint64(flagNumberPeriods, 1, "number of months for monthly vesting")
	cmd.Flags().Bool(flagContinuous, false, "set to continuous vesting.")
	cmd.Flags().Bool(flagTriggered, true, "set to false to deactivate periodic vesting until manually triggered")
	cmd.Flags().Bool(flagManual, false, "set to manual vesting")
	return cmd
}

func getVestedAccountFromFlags(baseAccount *authtypes.BaseAccount, coins sdk.Coins) (authexported.GenesisAccount, error) {
	vestingStart := viper.GetInt64(flagVestingStart)
	vestingEnd := viper.GetInt64(flagVestingEnd)
	vestingAmt, err := sdk.ParseCoins(viper.GetString(flagVestingAmt))
	period := viper.GetInt64(flagPeriod)
	numberPeriods := viper.GetInt64(flagNumberPeriods)
	continuous := viper.GetBool(flagContinuous)
	vestingTriggered := viper.GetBool(flagTriggered)
	manual := viper.GetBool(flagManual)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vesting amount: %w", err)
	}
	if vestingAmt.IsZero() {
		return baseAccount, nil
	}
	vestingAmt = vestingAmt.Sort()
	baseVestingAccount, err := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt, vestingEnd)
	if err != nil {
		return nil, err
	}

	if coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero() ||
		baseVestingAccount.OriginalVesting.IsAnyGT(coins) {
		return nil, errors.New("vesting amount cannot be greater than total amount")
	}

	if vestingStart == 0 {
		vestingStart = time.Now().Unix() + 10
	}

	switch {
	case manual:
		return vesting.NewManualVestingAccountRaw(baseVestingAccount, sdk.NewCoins()), nil

	case period != 0:
		periods := authvesting.Periods{}
		remaining := vestingAmt
		monthlyAmount := common.DivideCoins(vestingAmt, numberPeriods)

		for i := int64(0); i < numberPeriods-1; i++ {
			periods = append(periods, authvesting.Period{Length: period, Amount: monthlyAmount})
			remaining = remaining.Sub(monthlyAmount)
		}
		periods = append(periods, authvesting.Period{Length: period, Amount: remaining})
		endTime := vestingStart
		for _, p := range periods {
			endTime += p.Length
		}
		baseVestingAccount.EndTime = endTime
		if !vestingTriggered {
			return vesting.NewTriggeredVestingAccountRaw(baseVestingAccount, vestingStart, periods, false), nil
		}
		return authvesting.NewPeriodicVestingAccountRaw(baseVestingAccount, vestingStart, periods), nil

	case continuous && vestingEnd != 0:
		return authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart), nil

	case vestingEnd != 0:
		return authvesting.NewDelayedVestingAccountRaw(baseVestingAccount), nil

	default:
		return nil, errors.New("invalid vesting parameters; must supply start and end time or end time")
	}
}
