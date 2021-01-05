package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/auth/types"
)

const (
	flagClientHome    = "home-client"
	flagVestingStart  = "vesting-start-time"
	flagVestingEnd    = "vesting-end-time"
	flagVestingAmt    = "vesting-amount"
	flagPeriod        = "period"
	flagNumberPeriods = "num-periods"
	flagContinuous    = "continuous"
	flagManual        = "manual"
	flagUnlocker      = "unlocker"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(defaultNodeHome string) *cobra.Command {
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
			ctx := client.GetClientContextFromCmd(cmd)
			depCdc := ctx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			config := server.GetServerContextFromCmd(cmd).Config
			config.SetRoot(ctx.HomeDir)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, ctx.HomeDir, inBuf)
				if err != nil {
					return err
				}

				info, err := kb.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			// create concrete account type based on input parameters
			balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			genAccount, err := getVestedAccountFromFlags(baseAccount, coins)
			if err != nil {
				return err
			}
			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			accs = append(accs, genAccount)
			accs = authtypes.SanitizeGenesisAccounts(accs)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := cdc.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(depCdc, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

			bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagPeriod, 0, "set to periodic vesting with period in seconds")
	cmd.Flags().Uint64(flagNumberPeriods, 1, "number of months for monthly vesting")
	cmd.Flags().Bool(flagContinuous, false, "set to continuous vesting.")
	cmd.Flags().Bool(flagManual, false, "set to manual vesting")
	cmd.Flags().String(flagUnlocker, "", "address that can unlock this account's locked coins")
	return cmd
}

func getVestedAccountFromFlags(baseAccount *authtypes.BaseAccount, coins sdk.Coins) (authexported.GenesisAccount, error) {
	vestingStart := viper.GetInt64(flagVestingStart)
	vestingEnd := viper.GetInt64(flagVestingEnd)
	vestingAmt, err := sdk.ParseCoins(viper.GetString(flagVestingAmt))
	if err != nil {
		return nil, fmt.Errorf("failed to parse vesting amount: %w", err)
	}
	period := viper.GetInt64(flagPeriod)
	numberPeriods := viper.GetInt64(flagNumberPeriods)
	continuous := viper.GetBool(flagContinuous)
	manual := viper.GetBool(flagManual)

	if vestingAmt.IsZero() {
		return baseAccount, nil
	}
	vestingAmt = vestingAmt.Sort()
	baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingEnd)

	if coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero() ||
		baseVestingAccount.OriginalVesting.IsAnyGT(coins) {
		return nil, errors.New("vesting amount cannot be greater than total amount")
	}

	if vestingStart == 0 {
		vestingStart = time.Now().Unix() + 10
	}

	switch {
	case manual:
		unlocker, err := sdk.AccAddressFromBech32(viper.GetString(flagUnlocker))
		if err != nil || unlocker.Empty() {
			return nil, errors.New("unlocker address is in incorrect format")
		}
		return types.NewManualVestingAccountRaw(baseVestingAccount, sdk.NewCoins(), unlocker), nil

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
		return authvesting.NewPeriodicVestingAccountRaw(baseVestingAccount, vestingStart, periods), nil

	case continuous && vestingEnd != 0:
		return authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart), nil

	case vestingEnd != 0:
		return authvesting.NewDelayedVestingAccountRaw(baseVestingAccount), nil

	default:
		return nil, errors.New("invalid vesting parameters; must supply start and end time or end time")
	}
}
