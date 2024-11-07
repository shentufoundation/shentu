package e2e

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shentugovtypes "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

func modifyGenesis(path, moniker, amountStr string, addrAll []sdk.AccAddress, certAcc sdk.AccAddress, denom string) error {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)
	config.Moniker = moniker

	coins, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}

	var balances []banktypes.Balance
	var genAccounts []*authtypes.BaseAccount
	for _, addr := range addrAll {
		balance := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
		balances = append(balances, balance)
		genAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)
		genAccounts = append(genAccounts, genAccount)
	}

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	certGenState := certtypes.GetGenesisStateFromAppState(cdc, appState)
	certifier := certtypes.Certifier{
		Address: certAcc.String(),
	}
	certGenState.Certifiers = append(certGenState.Certifiers, certifier)
	certGenStateBz, err := cdc.MarshalJSON(&certGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal cert genesis state: %w", err)
	}
	appState[certtypes.ModuleName] = certGenStateBz

	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)
	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return fmt.Errorf("failed to get accounts from any: %w", err)
	}

	for _, addr := range addrAll {
		if accs.Contains(addr) {
			return fmt.Errorf("failed to add account to genesis state; account already exists: %s", addr)
		}
	}

	// Add the new account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	for _, genAcct := range genAccounts {
		accs = append(accs, genAcct)
		accs = authtypes.SanitizeGenesisAccounts(accs)
	}

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

	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances...)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)
	stakingGenState.Params.BondDenom = denom
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal staking genesis state: %s", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	maxDepositPeriod := 10 * time.Minute
	votingPeriod := 15 * time.Second
	expeditedVoting := 13 * time.Second

	govGenState := shentugovtypes.DefaultGenesisState()
	govGenState.Params.MaxDepositPeriod = &maxDepositPeriod
	govGenState.Params.VotingPeriod = &votingPeriod
	govGenState.Params.ExpeditedVotingPeriod = &expeditedVoting

	govGenStateBz, err := cdc.MarshalJSON(govGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}
	genDoc.AppState = appStateJSON

	return genutil.ExportGenesisFile(genDoc, genFile)
}
