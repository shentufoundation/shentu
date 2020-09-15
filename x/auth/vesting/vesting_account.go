package vesting

import (
	"encoding/json"
	"errors"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	customauth "github.com/certikfoundation/shentu/x/auth/internal/types"
)

//-----------------------------------------------------------------------------
// Triggered Vesting Account
//

// Compile-time type assertions
var _ vestexported.VestingAccount = (*TriggeredVestingAccount)(nil)
var _ authexported.GenesisAccount = (*TriggeredVestingAccount)(nil)

func init() {
	customauth.RegisterAccountTypeCodec(&vesttypes.BaseVestingAccount{}, "auth/BaseVestingAccount")
	customauth.RegisterAccountTypeCodec(&vesttypes.ContinuousVestingAccount{}, "auth/ContinuousVestingAccount")
	customauth.RegisterAccountTypeCodec(&vesttypes.DelayedVestingAccount{}, "auth/DelayedVestingAccount")
	customauth.RegisterAccountTypeCodec(&vesttypes.PeriodicVestingAccount{}, "auth/PeriodicVestingAccount")
	customauth.RegisterAccountTypeCodec(&TriggeredVestingAccount{}, "auth/TriggeredVestingAccount")
	customauth.RegisterAccountTypeCodec(&supply.ModuleAccount{}, "cosmos-sdk/ModuleAccount")

	authtypes.RegisterAccountTypeCodec(&TriggeredVestingAccount{}, "auth/TriggeredVestingAccount")
}

// TriggeredVestingAccount implements the VestingAccount interface.
// It behaves like PeriodicVestingAccount when activated.
type TriggeredVestingAccount struct {
	*vesttypes.BaseVestingAccount

	StartTime      int64             `json:"start_time" yaml:"start_time"`
	VestingPeriods vesttypes.Periods `json:"vesting_periods" yaml:"vesting_periods"`
	Activated      bool              `json:"activated" yaml:"activated"`
}

func NewTriggeredVestingAccountRaw(bva *vesttypes.BaseVestingAccount, startTime int64,
	periods vesttypes.Periods, activated bool) *TriggeredVestingAccount {
	return &TriggeredVestingAccount{
		BaseVestingAccount: bva,
		StartTime:          startTime,
		VestingPeriods:     periods,
		Activated:          activated,
	}
}

func NewTriggeredVestingAccount(baseAcc *authtypes.BaseAccount, originalVesting sdk.Coins, startTime int64,
	periods vesttypes.Periods, activated bool) *TriggeredVestingAccount {
	endTime := startTime
	for _, p := range periods {
		endTime += p.Length
	}
	baseVestingAcc := &vesttypes.BaseVestingAccount{
		BaseAccount:     baseAcc,
		OriginalVesting: originalVesting,
		EndTime:         endTime,
	}

	return &TriggeredVestingAccount{
		BaseVestingAccount: baseVestingAcc,
		StartTime:          startTime,
		VestingPeriods:     periods,
		Activated:          activated,
	}
}

// Returns the total number of vested coins. If no coins are vested, nil is returned.
func (tva TriggeredVestingAccount) GetVestedCoins(blockTime time.Time) sdk.Coins {
	var vestedCoins sdk.Coins
	startTime := tva.StartTime

	// If not active, assume that nothing has been vested.
	if !tva.Activated {
		return vestedCoins
	}

	// We must handle the case where the start time for a vesting account has
	// been set into the future or when the start of the chain is not exactly
	// known.
	if blockTime.Unix() <= startTime {
		return vestedCoins
	} else if blockTime.Unix() >= tva.EndTime {
		return tva.OriginalVesting
	}

	// Track the start time of the next period.
	currentPeriodStartTime := startTime

	// For each period, if the period is over, add those coins as vested and check the next period.
	for _, period := range tva.VestingPeriods {
		x := blockTime.Unix() - currentPeriodStartTime
		if x < period.Length {
			break
		}

		vestedCoins = vestedCoins.Add(period.Amount...)

		// Update the start time of the next period.
		currentPeriodStartTime += period.Length
	}

	return vestedCoins
}

// Returns the total number of vesting coins. If no coins are vesting, nil is returned.
func (tva TriggeredVestingAccount) GetVestingCoins(blockTime time.Time) sdk.Coins {
	// If not active, assume that none of the original vesting has been vested
	// and all of it is vesting.
	if !tva.Activated {
		return tva.OriginalVesting
	}
	return tva.OriginalVesting.Sub(tva.GetVestedCoins(blockTime))
}

// SpendableCoins returns the total number of spendable coins per denom.
func (tva TriggeredVestingAccount) SpendableCoins(blockTime time.Time) sdk.Coins {
	return tva.BaseVestingAccount.SpendableCoinsVestingAccount(tva.GetVestingCoins(blockTime))
}

// TrackDelegation tracks a desired delegation amount by setting the appropriate
// values for the amount of delegated vesting, delegated free, and reducing the
// overall amount of base coins.
func (tva *TriggeredVestingAccount) TrackDelegation(blockTime time.Time, amount sdk.Coins) {
	tva.BaseVestingAccount.TrackDelegation(tva.GetVestingCoins(blockTime), amount)
}

func (tva TriggeredVestingAccount) GetStartTime() int64 {
	return tva.StartTime
}

// Validate checks for errors on the account fields.
func (tva TriggeredVestingAccount) Validate() error {
	if tva.GetStartTime() >= tva.GetEndTime() {
		return errors.New("vesting start-time cannot be before end-time")
	}
	endTime := tva.StartTime
	originalVesting := sdk.NewCoins()
	for _, p := range tva.VestingPeriods {
		endTime += p.Length
		originalVesting = originalVesting.Add(p.Amount...)
	}
	if endTime != tva.EndTime {
		return errors.New("vesting end time does not match length of all vesting periods")
	}
	if !originalVesting.IsEqual(tva.OriginalVesting) {
		return errors.New("original vesting coins does not match the sum of all coins in vesting periods")
	}

	return tva.BaseVestingAccount.Validate()
}

type triggeredVestingAccountPretty struct {
	Address          sdk.AccAddress    `json:"address" yaml:"address"`
	Coins            sdk.Coins         `json:"coins" yaml:"coins"`
	PubKey           string            `json:"public_key" yaml:"public_key"`
	AccountNumber    uint64            `json:"account_number" yaml:"account_number"`
	Sequence         uint64            `json:"sequence" yaml:"sequence"`
	OriginalVesting  sdk.Coins         `json:"original_vesting" yaml:"original_vesting"`
	DelegatedFree    sdk.Coins         `json:"delegated_free" yaml:"delegated_free"`
	DelegatedVesting sdk.Coins         `json:"delegated_vesting" yaml:"delegated_vesting"`
	EndTime          int64             `json:"end_time" yaml:"end_time"`
	StartTime        int64             `json:"start_time" yaml:"start_time"`
	VestingPeriods   vesttypes.Periods `json:"vesting_periods" yaml:"vesting_periods"`
	Activated        bool              `json:"activated" yaml:"activated"`
}

func (tva TriggeredVestingAccount) MarshalJSON() ([]byte, error) {
	alias := triggeredVestingAccountPretty{
		Address:          tva.Address,
		Coins:            tva.Coins,
		AccountNumber:    tva.AccountNumber,
		Sequence:         tva.Sequence,
		OriginalVesting:  tva.OriginalVesting,
		DelegatedFree:    tva.DelegatedFree,
		DelegatedVesting: tva.DelegatedVesting,
		EndTime:          tva.EndTime,
		StartTime:        tva.StartTime,
		VestingPeriods:   tva.VestingPeriods,
		Activated:        tva.Activated,
	}

	if tva.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, tva.PubKey)
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	return json.Marshal(alias)
}

func (tva *TriggeredVestingAccount) UnmarshalJSON(bz []byte) error {
	var alias triggeredVestingAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	var (
		pk  crypto.PubKey
		err error
	)

	if alias.PubKey != "" {
		pk, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, alias.PubKey)
		if err != nil {
			return err
		}
	}

	tva.BaseVestingAccount = &vesttypes.BaseVestingAccount{
		BaseAccount:      authtypes.NewBaseAccount(alias.Address, alias.Coins, pk, alias.AccountNumber, alias.Sequence),
		OriginalVesting:  alias.OriginalVesting,
		DelegatedFree:    alias.DelegatedFree,
		DelegatedVesting: alias.DelegatedVesting,
		EndTime:          alias.EndTime,
	}
	tva.StartTime = alias.StartTime
	tva.VestingPeriods = alias.VestingPeriods
	tva.Activated = alias.Activated

	return nil
}

func (tva TriggeredVestingAccount) String() string {
	out, _ := tva.MarshalYAML()
	return out.(string)
}

func (tva TriggeredVestingAccount) MarshalYAML() (interface{}, error) {
	alias := triggeredVestingAccountPretty{
		Address:          tva.Address,
		Coins:            tva.Coins,
		AccountNumber:    tva.AccountNumber,
		Sequence:         tva.Sequence,
		OriginalVesting:  tva.OriginalVesting,
		DelegatedFree:    tva.DelegatedFree,
		DelegatedVesting: tva.DelegatedVesting,
		EndTime:          tva.EndTime,
		StartTime:        tva.StartTime,
		VestingPeriods:   tva.VestingPeriods,
		Activated:        tva.Activated,
	}

	if tva.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, tva.PubKey)
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}
