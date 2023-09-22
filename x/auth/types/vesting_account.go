package types

import (
	"errors"
	"time"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// Compile-time type assertions
var _ vestexported.VestingAccount = (*ManualVestingAccount)(nil)
var _ authexported.GenesisAccount = (*ManualVestingAccount)(nil)

//-----------------------------------------------------------------------------
// Manual Vesting Account
//

// NewManualVestingAccountRaw creates a new ManualVestingAccount object from BaseVestingAccount.
func NewManualVestingAccountRaw(bva *vesttypes.BaseVestingAccount, vestedCoins sdk.Coins, unlocker sdk.AccAddress) *ManualVestingAccount {
	return &ManualVestingAccount{
		BaseVestingAccount: bva,
		VestedCoins:        vestedCoins,
		Unlocker:           unlocker.String(),
	}
}

func NewManualVestingAccount(baseAcc *authtypes.BaseAccount, origVesting, vestedCoins sdk.Coins, unlocker sdk.AccAddress) *ManualVestingAccount {
	baseVestingAcc := &vesttypes.BaseVestingAccount{
		BaseAccount:     baseAcc,
		OriginalVesting: origVesting,
		EndTime:         0,
	}
	return &ManualVestingAccount{
		BaseVestingAccount: baseVestingAcc,
		VestedCoins:        vestedCoins,
		Unlocker:           unlocker.String(),
	}
}

// Returns the total number of vested coins. If no coins are vested, nil is returned.
func (mva ManualVestingAccount) GetVestedCoins(blockTime time.Time) sdk.Coins {
	if !mva.VestedCoins.IsZero() {
		return mva.VestedCoins
	}
	return nil
}

// Returns the total number of vesting coins. If no coins are vesting, nil is returned.
func (mva ManualVestingAccount) GetVestingCoins(blockTime time.Time) sdk.Coins {
	return mva.OriginalVesting.Sub(mva.GetVestedCoins(blockTime))
}

// LockedCoins returns the set of coins that are not spendable.
func (mva ManualVestingAccount) LockedCoins(blockTime time.Time) sdk.Coins {
	return mva.BaseVestingAccount.LockedCoinsFromVesting(mva.GetVestingCoins(blockTime))
}

// TrackDelegation tracks a desired delegation amount by setting the appropriate
// values for the amount of delegated vesting, delegated free, and reducing the
// overall amount of base coins.
func (mva *ManualVestingAccount) TrackDelegation(blockTime time.Time, balance, amount sdk.Coins) {
	mva.BaseVestingAccount.TrackDelegation(balance, mva.GetVestingCoins(blockTime), amount)
}

// GetStartTime returns zero since a manual vesting account has no start time.
func (ManualVestingAccount) GetStartTime() int64 {
	return 0
}

// Validate checks for errors on the account fields.
func (mva ManualVestingAccount) Validate() error {
	if mva.VestedCoins.IsAnyGT(mva.OriginalVesting) {
		return errors.New("vested amount cannot be greater than the original vesting amount")
	}
	return mva.BaseVestingAccount.Validate()
}

type manualVestingAccountYAML struct {
	Address          sdk.AccAddress `json:"address" yaml:"address"`
	PubKey           string         `json:"public_key" yaml:"public_key"`
	AccountNumber    uint64         `json:"account_number" yaml:"account_number"`
	Sequence         uint64         `json:"sequence" yaml:"sequence"`
	OriginalVesting  sdk.Coins      `json:"original_vesting" yaml:"original_vesting"`
	DelegatedFree    sdk.Coins      `json:"delegated_free" yaml:"delegated_free"`
	DelegatedVesting sdk.Coins      `json:"delegated_vesting" yaml:"delegated_vesting"`
	EndTime          int64          `json:"end_time" yaml:"end_time"`
	VestedCoins      sdk.Coins      `json:"vested_coins" yaml:"vested_coins"`
	Unlocker         sdk.AccAddress `json:"unlocker" yaml:"unlocker"`
}

func (mva ManualVestingAccount) String() string {
	out, _ := mva.MarshalYAML()
	return out.(string)
}

func (mva ManualVestingAccount) MarshalYAML() (interface{}, error) {
	accAddr, err := sdk.AccAddressFromBech32(mva.Address)
	if err != nil {
		return nil, err
	}
	unlocker, err := sdk.AccAddressFromBech32(mva.Unlocker)
	if err != nil {
		return nil, err
	}

	alias := manualVestingAccountYAML{
		Address:          accAddr,
		AccountNumber:    mva.AccountNumber,
		Sequence:         mva.Sequence,
		OriginalVesting:  mva.OriginalVesting,
		DelegatedFree:    mva.DelegatedFree,
		DelegatedVesting: mva.DelegatedVesting,
		EndTime:          mva.EndTime,
		VestedCoins:      mva.VestedCoins,
		Unlocker:         unlocker,
	}

	pk := mva.GetPubKey()
	if pk != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pk)
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
