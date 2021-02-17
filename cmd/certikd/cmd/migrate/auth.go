package migrate

import (
	"errors"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v038auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v038"
	v039auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v039"
	"gopkg.in/yaml.v2"
)

func RegisterAuthLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cryptocodec.RegisterCrypto(cdc)
	cdc.RegisterInterface((*v038auth.GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*v038auth.Account)(nil), nil)
	cdc.RegisterConcrete(&v039auth.BaseAccount{}, "cosmos-sdk/Account", nil)
	cdc.RegisterConcrete(&v039auth.BaseVestingAccount{}, "cosmos-sdk/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&v039auth.ContinuousVestingAccount{}, "cosmos-sdk/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&v039auth.DelayedVestingAccount{}, "cosmos-sdk/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(&v039auth.PeriodicVestingAccount{}, "cosmos-sdk/PeriodicVestingAccount", nil)
	cdc.RegisterConcrete(&v039auth.ModuleAccount{}, "cosmos-sdk/ModuleAccount", nil)
	cdc.RegisterConcrete(&ManualVestingAccount{}, "auth/ManualVestingAccount", nil)
}

//-----------------------------------------------------------------------------
// Manual Vesting Account
//

// ManualVestingAccount implements the VestingAccount interface.
type ManualVestingAccount struct {
	*v039auth.BaseVestingAccount

	VestedCoins sdk.Coins      `json:"vested_coins" yaml:"vested_coins"`
	Unlocker    sdk.AccAddress `json:"unlocker" yaml:"unlocker"`
}

// NewManualVestingAccountRaw creates a new ManualVestingAccount object from BaseVestingAccount.
func NewManualVestingAccountRaw(bva *v039auth.BaseVestingAccount, vestedCoins sdk.Coins, unlocker sdk.AccAddress) *ManualVestingAccount {
	return &ManualVestingAccount{
		BaseVestingAccount: bva,
		VestedCoins:        vestedCoins,
		Unlocker:           unlocker,
	}
}

func NewManualVestingAccount(baseAcc *v039auth.BaseAccount, vestedCoins sdk.Coins, unlocker sdk.AccAddress) *ManualVestingAccount {
	baseVestingAcc := &v039auth.BaseVestingAccount{
		BaseAccount:     baseAcc,
		OriginalVesting: baseAcc.Coins,
		EndTime:         0,
	}
	return &ManualVestingAccount{
		BaseVestingAccount: baseVestingAcc,
		VestedCoins:        vestedCoins,
		Unlocker:           unlocker,
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

// GetStartTime returns zero since a manual vesting account has no start time.
func (mva ManualVestingAccount) GetStartTime() int64 {
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
	Coins            sdk.Coins      `json:"coins" yaml:"coins"`
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

type manualVestingAccountJSON struct {
	Address          sdk.AccAddress     `json:"address" yaml:"address"`
	Coins            sdk.Coins          `json:"coins" yaml:"coins"`
	PubKey           cryptotypes.PubKey `json:"public_key" yaml:"public_key"`
	AccountNumber    uint64             `json:"account_number" yaml:"account_number"`
	Sequence         uint64             `json:"sequence" yaml:"sequence"`
	OriginalVesting  sdk.Coins          `json:"original_vesting" yaml:"original_vesting"`
	DelegatedFree    sdk.Coins          `json:"delegated_free" yaml:"delegated_free"`
	DelegatedVesting sdk.Coins          `json:"delegated_vesting" yaml:"delegated_vesting"`
	EndTime          int64              `json:"end_time" yaml:"end_time"`
	VestedCoins      sdk.Coins          `json:"vested_coins" yaml:"vested_coins"`
	Unlocker         sdk.AccAddress     `json:"unlocker" yaml:"unlocker"`
}

func (mva ManualVestingAccount) MarshalJSON() ([]byte, error) {
	alias := manualVestingAccountJSON{
		Address:          mva.Address,
		Coins:            mva.Coins,
		PubKey:           mva.PubKey,
		AccountNumber:    mva.AccountNumber,
		Sequence:         mva.Sequence,
		OriginalVesting:  mva.OriginalVesting,
		DelegatedFree:    mva.DelegatedFree,
		DelegatedVesting: mva.DelegatedVesting,
		EndTime:          mva.EndTime,
		VestedCoins:      mva.VestedCoins,
		Unlocker:         mva.Unlocker,
	}

	return legacy.Cdc.MarshalJSON(alias)
}

func (mva *ManualVestingAccount) UnmarshalJSON(bz []byte) error {
	var alias manualVestingAccountJSON
	if err := legacy.Cdc.UnmarshalJSON(bz, &alias); err != nil {
		return err
	}

	mva.BaseVestingAccount = &v039auth.BaseVestingAccount{
		BaseAccount:      v039auth.NewBaseAccount(alias.Address, alias.Coins, alias.PubKey, alias.AccountNumber, alias.Sequence),
		OriginalVesting:  alias.OriginalVesting,
		DelegatedFree:    alias.DelegatedFree,
		DelegatedVesting: alias.DelegatedVesting,
		EndTime:          alias.EndTime,
	}
	mva.VestedCoins = alias.VestedCoins
	mva.Unlocker = alias.Unlocker

	return nil
}

func (mva ManualVestingAccount) String() string {
	out, _ := mva.MarshalYAML()
	return out.(string)
}

func (mva ManualVestingAccount) MarshalYAML() (interface{}, error) {
	alias := manualVestingAccountYAML{
		Address:          mva.Address,
		Coins:            mva.Coins,
		AccountNumber:    mva.AccountNumber,
		Sequence:         mva.Sequence,
		OriginalVesting:  mva.OriginalVesting,
		DelegatedFree:    mva.DelegatedFree,
		DelegatedVesting: mva.DelegatedVesting,
		EndTime:          mva.EndTime,
		VestedCoins:      mva.VestedCoins,
		Unlocker:         mva.Unlocker,
	}

	if mva.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, mva.PubKey)
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
