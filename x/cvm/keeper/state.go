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

	"github.com/certikfoundation/shentu/v2/x/cvm/types"
)

// State is the CVM state object. It implements acmstate.ReaderWriter.
type State struct {
	ctx   sdk.Context
	ak    types.AccountKeeper
	bk    types.BankKeeper
	sk    types.StakingKeeper
	store sdk.KVStore
	cdc   codec.BinaryCodec
}

// NewState returns a new instance of State type data.
func (k Keeper) NewState(ctx sdk.Context) *State {
	return &State{
		ctx:   ctx,
		ak:    k.ak,
		bk:    k.bk,
		sk:    k.sk,
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
	balance := s.bk.GetBalance(s.ctx, addr, s.sk.BondDenom(s.ctx)).Amount.Uint64()
	contMeta, err := s.GetAddressMeta(address)
	if err != nil {
		return nil, err
	}

	var evmCode, wasmCode acm.Bytecode
	codeData := s.store.Get(types.CodeStoreKey(address))
	if len(codeData) > 0 {
		var cvmCode types.CVMCode
		s.cdc.MustUnmarshal(codeData, &cvmCode)
		if cvmCode.CodeType == types.CVMCodeTypeEVMCode {
			evmCode = cvmCode.Code
		} else {
			wasmCode = cvmCode.Code
		}
	}

	acc := acm.Account{
		Address:  address,
		Balance:  balance,
		EVMCode:  evmCode,
		WASMCode: wasmCode,
		Permissions: permission.AccountPermissions{
			Base: permission.BasePermissions{
				Perms:  permission.Call | permission.CreateContract,
				SetBit: permission.AllPermFlags,
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
	var cvmCode types.CVMCode
	if len(updatedAccount.WASMCode) > 0 {
		cvmCode = types.NewCVMCode(types.CVMCodeTypeEWASMCode, updatedAccount.WASMCode)
	} else {
		cvmCode = types.NewCVMCode(types.CVMCodeTypeEVMCode, updatedAccount.EVMCode)
	}
	s.store.Set(types.CodeStoreKey(updatedAccount.Address), s.cdc.MustMarshal(&cvmCode))
	oldBalance := s.bk.GetBalance(s.ctx, address, s.sk.BondDenom(s.ctx))
	newBalance := sdk.NewInt64Coin(s.sk.BondDenom(s.ctx), int64(updatedAccount.Balance))
	if newBalance.Amount.GT(oldBalance.Amount) {
		coins := sdk.Coins{newBalance.Sub(oldBalance)}
		if err := s.bk.MintCoins(s.ctx, types.ModuleName, coins); err != nil {
			return err
		}
		if err := s.bk.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, address, coins); err != nil {
			return err
		}
	} else if newBalance.Amount.LT(oldBalance.Amount) {
		coins := sdk.Coins{oldBalance.Sub(newBalance)}
		if err := s.bk.SendCoinsFromAccountToModule(s.ctx, address, types.ModuleName, coins); err != nil {
			return err
		}
		if err := s.bk.BurnCoins(s.ctx, types.ModuleName, coins); err != nil {
			return err
		}
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

	var metaList types.ContractMetas
	err := s.cdc.Unmarshal(bz, &metaList)
	if err != nil {
		return nil, err
	}

	var res []*acm.ContractMeta
	copy(res, metaList.Metas)
	return res, err
}

// SetAddressMeta sets the metadata hash for an address
func (s *State) SetAddressMeta(address crypto.Address, contMeta []*acm.ContractMeta) error {
	var metadata types.ContractMetas
	metadata.Metas = append(metadata.Metas, contMeta...)
	bz, err := s.cdc.Marshal(&metadata)
	if err != nil {
		panic(err)
	}
	s.store.Set(types.AddressMetaStoreKey(address), bz)
	return err
}
