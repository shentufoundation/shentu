package migrate

import (
	"errors"
	"time"

	vesting "github.com/certikfoundation/shentu/x/auth/types"

	"gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	v038auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v038"
	v039auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v039"
	v040auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	v040vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
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

// convertBaseAccount converts a 0.39 BaseAccount to a 0.40 BaseAccount.
func convertBaseAccount(old *v039auth.BaseAccount) *v040auth.BaseAccount {
	var any *codectypes.Any
	// If the old genesis had a pubkey, we pack it inside an Any. Or else, we
	// just leave it nil.
	if old.PubKey != nil {
		var err error
		any, err = codectypes.NewAnyWithValue(old.PubKey)
		if err != nil {
			panic(err)
		}
	}

	return &v040auth.BaseAccount{
		Address:       old.Address.String(),
		PubKey:        any,
		AccountNumber: old.AccountNumber,
		Sequence:      old.Sequence,
	}
}

// convertBaseVestingAccount converts a 0.39 BaseVestingAccount to a 0.40 BaseVestingAccount.
func convertBaseVestingAccount(old *v039auth.BaseVestingAccount) *v040vesting.BaseVestingAccount {
	baseAccount := convertBaseAccount(old.BaseAccount)

	return &v040vesting.BaseVestingAccount{
		BaseAccount:      baseAccount,
		OriginalVesting:  old.OriginalVesting,
		DelegatedFree:    old.DelegatedFree,
		DelegatedVesting: old.DelegatedVesting,
		EndTime:          old.EndTime,
	}
}

// Migrate accepts exported x/auth genesis state from v0.38/v0.39 and migrates
// it to v0.40 x/auth genesis state. The migration includes:
//
// - Removing coins from account encoding.
// - Re-encode in v0.40 GenesisState.
func authMigrate(authGenState v039auth.GenesisState) *v040auth.GenesisState {
	// Convert v0.39 accounts to v0.40 ones.
	var v040Accounts = make([]v040auth.GenesisAccount, len(authGenState.Accounts))
	for i, v039Account := range authGenState.Accounts {
		switch v039Account := v039Account.(type) {
		case *v039auth.BaseAccount:
			{
				v040Accounts[i] = convertBaseAccount(v039Account)
			}
		case *v039auth.ModuleAccount:
			{
				v040Accounts[i] = &v040auth.ModuleAccount{
					BaseAccount: convertBaseAccount(v039Account.BaseAccount),
					Name:        v039Account.Name,
					Permissions: v039Account.Permissions,
				}
			}
		case *v039auth.BaseVestingAccount:
			{
				v040Accounts[i] = convertBaseVestingAccount(v039Account)
			}
		case *v039auth.ContinuousVestingAccount:
			{
				v040Accounts[i] = &v040vesting.ContinuousVestingAccount{
					BaseVestingAccount: convertBaseVestingAccount(v039Account.BaseVestingAccount),
					StartTime:          v039Account.StartTime,
				}
			}
		case *v039auth.DelayedVestingAccount:
			{
				v040Accounts[i] = &v040vesting.DelayedVestingAccount{
					BaseVestingAccount: convertBaseVestingAccount(v039Account.BaseVestingAccount),
				}
			}
		case *v039auth.PeriodicVestingAccount:
			{
				vestingPeriods := make([]v040vesting.Period, len(v039Account.VestingPeriods))
				for j, period := range v039Account.VestingPeriods {
					vestingPeriods[j] = v040vesting.Period{
						Length: period.Length,
						Amount: period.Amount,
					}
				}
				v040Accounts[i] = &v040vesting.PeriodicVestingAccount{
					BaseVestingAccount: convertBaseVestingAccount(v039Account.BaseVestingAccount),
					StartTime:          v039Account.StartTime,
					VestingPeriods:     vestingPeriods,
				}
			}
		case *ManualVestingAccount:
			{
				v040Accounts[i] = &vesting.ManualVestingAccount{
					BaseVestingAccount: convertBaseVestingAccount(v039Account.BaseVestingAccount),
					VestedCoins:        v039Account.VestedCoins,
					Unlocker:           v039Account.Unlocker.String(),
				}
			}
		default:
			panic(sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "got invalid type %T", v039Account))
		}
	}

	// Convert v0.40 accounts into Anys.
	anys := make([]*codectypes.Any, len(v040Accounts))
	for i, v040Account := range v040Accounts {
		any, err := codectypes.NewAnyWithValue(v040Account)
		if err != nil {
			panic(err)
		}

		anys[i] = any
	}

	return &v040auth.GenesisState{
		Params: v040auth.Params{
			MaxMemoCharacters:      authGenState.Params.MaxMemoCharacters,
			TxSigLimit:             authGenState.Params.TxSigLimit,
			TxSizeCostPerByte:      authGenState.Params.TxSizeCostPerByte,
			SigVerifyCostED25519:   authGenState.Params.SigVerifyCostED25519,
			SigVerifyCostSecp256k1: authGenState.Params.SigVerifyCostSecp256k1,
		},
		Accounts: anys,
	}
}
