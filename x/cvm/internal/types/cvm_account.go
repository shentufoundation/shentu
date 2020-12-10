package types

import (
	"github.com/hyperledger/burrow/acm"
	yaml "gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Compile-time type assertions
var (
	_ authtypes.AccountI = (*CVMAccount)(nil)
)

//-----------------------------------------------------------------------------
// CVM Account
//

// CVMAccount implements the BaseAccount interface.
// It contains all CVM contract related info.
type CVMAccount struct {
	*authtypes.BaseAccount

	Code acm.Bytecode `json:"code" yaml:"code"`
	Abi  string       `json:"abi" yaml:"abi"`
}

// NewCVMAccount creates a new CVM account
func NewCVMAccount(baseAcc *authtypes.BaseAccount, code acm.Bytecode, abi string) *CVMAccount {
	return &CVMAccount{
		BaseAccount: baseAcc,
		Code:        code,
		Abi:         abi,
	}
}

type cvmAccountYAML struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`

	Code acm.Bytecode `json:"code" yaml:"code"`
	Abi  string       `json:"abi" yaml:"abi"`
}

type cvmAccountJSON struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`

	Code acm.Bytecode `json:"code" yaml:"code"`
	Abi  string       `json:"abi" yaml:"abi"`
}

func (ca CVMAccount) String() string {
	out, _ := ca.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of a CVMAccount.
func (ca CVMAccount) MarshalYAML() (interface{}, error) {
	addr, _ := sdk.AccAddressFromBech32(ca.Address)
	alias := cvmAccountYAML{
		Address:       addr,
		AccountNumber: ca.AccountNumber,
		Sequence:      ca.Sequence,

		Code: ca.Code,
		Abi:  ca.Abi,
	}

	if ca.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, ca.PubKey)
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

// MarshalJSON returns the JSON representation of a CVMAccount.
func (ca CVMAccount) MarshalJSON() ([]byte, error) {
	addr, _ := sdk.AccAddressFromBech32(ca.Address)
	alias := cvmAccountJSON{
		Address:       addr,
		AccountNumber: ca.AccountNumber,
		Sequence:      ca.Sequence,

		Code: ca.Code,
		Abi:  ca.Abi,
	}

	return ModuleCdc.MarshalJSON(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into a CVMAccount.
func (ca *CVMAccount) UnmarshalJSON(bz []byte) error {
	var alias cvmAccountJSON
	if err := codec.Cdc.UnmarshalJSON(bz, &alias); err != nil {
		return err
	}

	ca.BaseAccount = authtypes.NewBaseAccount(
		alias.Address,
		alias.Coins,
		alias.PubKey,
		alias.AccountNumber,
		alias.Sequence,
	)

	ca.Code = alias.Code
	ca.Abi = alias.Abi

	return nil
}
