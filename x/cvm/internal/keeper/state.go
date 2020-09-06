package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/permission"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

// State is the CVM state object. It implements acmstate.ReaderWriter.
type State struct {
	ctx   sdk.Context
	ak    types.AccountKeeper
	store sdk.KVStore
	cdc   *codec.Codec
}

// NewState returns a new instance of State type data.
func (k Keeper) NewState(ctx sdk.Context) *State {
	return &State{
		ctx:   ctx,
		ak:    k.ak,
		store: ctx.KVStore(k.key),
		cdc:   k.cdc,
	}
}

// GetAccount gets an account by its address and returns nil if it does not
// exist (which should not be an error).
func (s *State) GetAccount(address crypto.Address) (*acm.Account, error) {
	addr := sdk.AccAddress(address.Bytes())
	account := s.ak.GetAccount(s.ctx, addr)
	if account == nil {
		return nil, nil
	}
	balance := account.GetCoins().AmountOf("uctk").Uint64()
	contMeta, err := s.GetAddressMeta(address)
	if err != nil {
		return nil, err
	}

	acc := acm.Account{
		Address: address,
		Balance: balance,
		EVMCode: s.store.Get(types.CodeStoreKey(address)),
		Permissions: permission.AccountPermissions{
			Base: permission.BasePermissions{
				Perms: permission.Call | permission.CreateContract,
			},
		},
		ContractMeta: contMeta,
	}
	return &acc, nil
}

// UpdateAccount updates the fields of updatedAccount by address, creating the
// account if it does not exist.
func (s *State) UpdateAccount(updatedAccount *acm.Account) error {
	if updatedAccount == nil {
		return fmt.Errorf("cannot update nil account")
	}
	address := sdk.AccAddress(updatedAccount.Address.Bytes())
	account := s.ak.GetAccount(s.ctx, address)
	if account == nil {
		account = s.ak.NewAccountWithAddress(s.ctx, address)
	}
	s.store.Set(types.CodeStoreKey(updatedAccount.Address), append([]byte{}, updatedAccount.EVMCode...))
	err := account.SetCoins(sdk.Coins{sdk.NewInt64Coin("uctk", int64(updatedAccount.Balance))})
	if err != nil {
		return err
	}
	s.ak.SetAccount(s.ctx, account)
	return s.SetAddressMeta(updatedAccount.Address, updatedAccount.ContractMeta)
}

// RemoveAccount removes the account at the address.
func (s *State) RemoveAccount(address crypto.Address) error {
	account := s.ak.GetAccount(s.ctx, address.Bytes())
	if account == nil {
		return fmt.Errorf("cannot remove non-existing account %s", address)
	}
	s.store.Delete(types.CodeStoreKey(address))
	s.store.Delete(types.AbiStoreKey(address))
	s.store.Delete(types.AddressMetaStoreKey(address))
	return nil
}

// GetStorage retrieves a 32-byte value stored at the key for the account at the
// address and returns Zero256 if the key does not exist or if the address does
// not exist.
//
// Note: burrow/acm/acmstate.StorageGetter claims that an error should be thrown
// upon non-existing address. However, when being embedded in acmstate.Cache,
// which is the case here, we cannot do that because the contract creation code
// might load from the new contract's storage, while the cache layer caches the
// account creation action hence the embedded storage getter will not be aware
// of it. Returning error in this case would fail the contract deployment.
func (s *State) GetStorage(address crypto.Address, key binary.Word256) (value []byte, err error) {
	bytes := s.store.Get(types.StorageStoreKey(address, key))
	if bytes == nil {
		return binary.Zero256.Bytes(), nil
	}
	return bytes, nil
}

// SetStorage stores a 32-byte value at the key for the account at the address.
// Setting to Zero256 removes the key.
func (s *State) SetStorage(address crypto.Address, key binary.Word256, value []byte) error {
	storeKey := types.StorageStoreKey(address, key)

	zero := true
	for _, b := range value {
		if b != 0 {
			zero = false
			break
		}
	}
	if zero {
		s.store.Delete(storeKey)
		return nil
	}

	s.store.Set(storeKey, value)
	return nil
}

// GetMetadata returns the metadata of the cvm module.
func (s *State) GetMetadata(metahash acmstate.MetadataHash) (string, error) {
	bz := s.store.Get(types.MetaHashStoreKey(metahash))
	res := string(bz)
	return res, nil
}

// SetMetadata sets the metadata of the cvm module.
func (s *State) SetMetadata(metahash acmstate.MetadataHash, metadata string) error {
	bz := []byte(metadata)
	s.store.Set(types.MetaHashStoreKey(metahash), bz)
	return nil
}

// GetAddressMeta returns the metadata hash of an address
func (s *State) GetAddressMeta(address crypto.Address) ([]*acm.ContractMeta, error) {
	bz := s.store.Get(types.AddressMetaStoreKey(address))
	if len(bz) == 0 {
		return []*acm.ContractMeta{}, nil
	}
	var metaList []acm.ContractMeta
	err := s.cdc.UnmarshalBinaryLengthPrefixed(bz, &metaList)
	var res []*acm.ContractMeta
	for i := range metaList {
		res = append(res, &metaList[i])
	}
	return res, err
}

// SetAddressMeta sets the metadata hash for an address
func (s *State) SetAddressMeta(address crypto.Address, contMeta []*acm.ContractMeta) error {
	var metadata []acm.ContractMeta
	for _, meta := range contMeta {
		metadata = append(metadata, *meta)
	}
	bz, err := s.cdc.MarshalBinaryLengthPrefixed(metadata)
	s.store.Set(types.AddressMetaStoreKey(address), bz)
	return err
}
