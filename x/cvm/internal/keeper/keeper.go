// Package keeper specifies the keeper for the cvm module.
package keeper

import (
	"bytes"
	gobin "encoding/binary"
	"encoding/hex"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/native"
	"github.com/hyperledger/burrow/execution/wasm"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/vm"
	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

// TransactionGasLimit is the gas limit of the block.
const TransactionGasLimit = uint64(5000000)

// Keeper implements SDK Keeper.
type Keeper struct {
	cdc        *codec.Codec
	key        sdk.StoreKey
	ak         types.AccountKeeper
	dk         types.DistributionKeeper
	ck         types.CertKeeper
	paramSpace params.Subspace
}

// NewKeeper creates a new instance of the CVM keeper.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, ak types.AccountKeeper, dk types.DistributionKeeper,
	ck types.CertKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        key,
		ak:         ak,
		dk:         dk,
		ck:         ck,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// Call executes the CVM call from caller to callee with the given data and gas limit.
func (k *Keeper) Call(ctx sdk.Context, caller, callee sdk.AccAddress, value uint64, data []byte, payloadMeta []*payload.ContractMeta,
	view, isEWASM, isRuntime bool) ([]byte, error) {
	state := k.NewState(ctx)

	callframe := engine.NewCallFrame(state, acmstate.Named("TxCache"))
	callerAddr := crypto.MustAddressFromBytes(caller)
	cache := callframe.Cache

	var sequenceBytes []byte
	if view {
		callframe.ReadOnly()
		sequenceBytes = make([]byte, 1)
	} else {
		sequenceBytes = k.getAccountSeqNum(ctx, caller)
	}

	var calleeAddr crypto.Address
	var code acm.Bytecode
	var err error
	if callee == nil {
		calleeAddr = crypto.NewContractAddress(callerAddr, sequenceBytes)
		if err = native.CreateAccount(cache, calleeAddr); err != nil {
			return nil, types.ErrCodedError(errors.GetCode(err))
		}
		code = data
		err = native.UpdateContractMeta(cache, state, calleeAddr, payloadMeta)
	} else {
		calleeAddr = crypto.MustAddressFromBytes(callee)
		calleeAddr, code, isEWASM, err = getCallee(callee, cache)
		if len(code) == 0 && !bytes.Equal(data, []byte{}) {
			return nil, types.ErrCodedError(errors.Codes.CodeOutOfBounds)
		}
	}
	if err != nil {
		return nil, types.ErrCodedError(errors.GetCode(err))
	}

	gasRate := k.GetGasRate(ctx)
	originalGas, err := k.getOriginalGas(ctx, gasRate)
	if err != nil {
		return nil, types.ErrCodedError(errors.GetCode(err))
	}
	gasTracker := originalGas

	callParams := engine.CallParams{
		Origin: callerAddr,
		Caller: callerAddr,
		Callee: calleeAddr,
		Input:  data,
		Value:  value,
		Gas:    &gasTracker,
	}
	options := vm.CVMOptions{
		Nonce: sequenceBytes,
	}
	cc := CertificateCallable{
		ctx:        ctx,
		certKeeper: k.ck,
	}
	registerCVMNative(&options, cc)

	newCVM := vm.NewCVM(options)
	logger := ctx.Logger()
	txHash := tmhash.Sum(ctx.TxBytes())
	bc := NewBlockChain(ctx, *k)

	logger.Info("CVM Start", "txHash", hex.EncodeToString(txHash))

	var ret []byte
	if isEWASM {
		if isRuntime {
			ret = code
		} else {
			ret, err = wasm.RunWASM(cache, callParams, code)
		}
	} else {
		ret, err = newCVM.Execute(cache, bc, NewEventSink(ctx), callParams, code)
	}
	defer func() {
		logger.Info("CVM Stop", "result", hex.EncodeToString(ret))
	}()

	// Refund cannot exceed half of the total gas cost.
	// Only refund when there is no error.
	if err != nil {
		gasTracker = gasTracker + vm.Min((originalGas-gasTracker)/2, newCVM.GetRefund())
	}

	// GasTracker is guaranteed to not underflow during CVM execution.
	fee := originalGas - gasTracker
	ctx.GasMeter().ConsumeGas((fee+gasRate-1)/gasRate, "CVM execution fee")
	if err != nil {
		return nil, types.ErrCodedError(errors.GetCode(err))
	}

	if callee == nil {
		if isEWASM {
			err = native.InitWASMCode(cache, calleeAddr, ret)
		} else {
			err = native.InitEVMCode(cache, calleeAddr, ret)
		}
		if err != nil {
			return nil, types.ErrCodedError(errors.GetCode(err))
		}
		ret = calleeAddr.Bytes()
	}
	if err = cache.Sync(state); err != nil {
		return nil, types.ErrCodedError(errors.GetCode(err))
	}

	return ret, nil
}

// Send executes the send transaction from caller to callee with the given amount of tokens.
func (k Keeper) Send(ctx sdk.Context, caller, callee sdk.AccAddress, coins sdk.Coins) error {
	value := coins.AmountOf(common.MicroCTKDenom).Uint64()
	if value <= 0 {
		return sdkerrors.ErrInvalidCoins
	}
	_, err := k.Call(ctx, caller, callee, value, nil, nil, false, false, false)
	return err
}

// GetCode returns the code at the given account address.
func (k Keeper) GetCode(ctx sdk.Context, addr crypto.Address) ([]byte, error) {
	state := k.NewState(ctx)
	callframe := engine.NewCallFrame(state, acmstate.Named("TxCache"))
	cache := callframe.Cache
	acc, err := cache.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, nil
	}
	return acc.EVMCode, nil
}

// getCallee returns the callee address and bytecode of a given account address.
func getCallee(callee sdk.AccAddress, cache *acmstate.Cache) (crypto.Address, acm.Bytecode, bool, error) {
	calleeAddr := crypto.MustAddressFromBytes(callee)
	acc, err := cache.GetAccount(calleeAddr)
	if err != nil {
		return crypto.Address{}, nil, false, err
	}
	if len(acc.WASMCode) != 0 {
		return calleeAddr, acc.WASMCode, true, err
	}
	return calleeAddr, acc.EVMCode, false, err
}

// getOriginalGas returns the original gas cost.
func (k Keeper) getOriginalGas(ctx sdk.Context, gasRate uint64) (uint64, error) {
	gasCurrent := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
	originalGas := gasCurrent * gasRate
	if originalGas < gasCurrent {
		return 0, types.ErrCodedError(errors.Codes.IntegerOverflow)
	}
	originalGas = vm.Min(originalGas, TransactionGasLimit)
	return originalGas, nil
}

// getAccount returns the account at a given address.
func (k Keeper) getAccount(ctx sdk.Context, address crypto.Address) *acm.Account {
	account, _ := k.NewState(ctx).GetAccount(address)
	return account
}

// GetStorage returns the value stored given the address and key.
func (k Keeper) GetStorage(ctx sdk.Context, address crypto.Address, key binary.Word256) ([]byte, error) {
	return k.NewState(ctx).GetStorage(address, key)
}

// SetAbi stores the abi for the address.
func (k Keeper) SetAbi(ctx sdk.Context, address crypto.Address, abi []byte) {
	ctx.KVStore(k.key).Set(types.AbiStoreKey(address), abi)
}

// getAbi returns the abi at the given address.
func (k Keeper) getAbi(ctx sdk.Context, address crypto.Address) []byte {
	return ctx.KVStore(k.key).Get(types.AbiStoreKey(address))
}

// getAddrMeta returns the meta hash of an address.
func (k Keeper) getAddrMeta(ctx sdk.Context, address crypto.Address) ([]*acm.ContractMeta, error) {
	state := k.NewState(ctx)
	return state.GetAddressMeta(address)
}

// getMeta returns the meta data of a meta hash.
func (k Keeper) getMeta(ctx sdk.Context, metaHash acmstate.MetadataHash) (string, error) {
	state := k.NewState(ctx)
	return state.GetMetadata(metaHash)
}

// StoreLastBlockHash stores the last block hash.
func (k Keeper) StoreLastBlockHash(ctx sdk.Context) {
	if ctx.BlockHeight() == 0 {
		return
	}
	block := ctx.BlockHeader().LastBlockId
	hash := block.GetHash()
	if hash == nil {
		return
	}
	ctx.KVStore(k.key).Set(types.BlockHashStoreKey(ctx.BlockHeight()), hash)
}

type logger struct {
	log.Logger
}

// Log implements github.com/go-kit/kit/log.Logger.
func (l *logger) Log(keyvals ...interface{}) error {
	l.Info("CVM", keyvals...)
	return nil
}

// WrapLogger converts a Tendermint logger into Burrow logger.
func WrapLogger(l log.Logger) *logging.Logger {
	return logging.NewLogger(&logger{l})
}

// getAccountSeqNum returns the account sequence number.
func (k Keeper) getAccountSeqNum(ctx sdk.Context, address sdk.AccAddress) []byte {
	callerAcc := k.ak.GetAccount(ctx, address)
	callerSequence := callerAcc.GetSequence()
	accountByte := make([]byte, 8)
	gobin.LittleEndian.PutUint64(accountByte, callerSequence)
	return accountByte
}

// RecycleCoins transfers tokens from the zero address to the community pool.
func (k Keeper) RecycleCoins(ctx sdk.Context) error {
	zeroAddrBytes := crypto.ZeroAddress.Bytes()
	acc := k.ak.GetAccount(ctx, zeroAddrBytes)
	if acc == nil {
		return nil
	}
	coins := acc.GetCoins()
	if coins.IsZero() {
		return nil
	}
	return k.dk.FundCommunityPool(ctx, coins, zeroAddrBytes)
}

// GetAllContracts gets all contracts for genesis export.
func (k Keeper) GetAllContracts(ctx sdk.Context) []types.Contract {
	contracts := make([]types.Contract, 0)
	store := ctx.KVStore(k.key)
	contractIterator := sdk.KVStorePrefixIterator(store, types.CodeStoreKeyPrefix)
	defer contractIterator.Close()

	for ; contractIterator.Valid(); contractIterator.Next() {
		addressBytes := contractIterator.Key()[len(types.CodeStoreKeyPrefix):]
		address, err := crypto.AddressFromBytes(addressBytes)
		if err != nil {
			panic(err)
		}
		var code types.CVMCode
		k.cdc.MustUnmarshalBinaryLengthPrefixed(contractIterator.Value(), &code)
		abi := k.getAbi(ctx, address)
		addrMeta, err := k.getAddrMeta(ctx, address)

		var meta []types.ContractMeta
		for _, adm := range addrMeta {
			meta = append(meta, types.ContractMeta{CodeHash: adm.CodeHash, MetadataHash: adm.MetadataHash})
		}

		if err != nil {
			panic(err)
		}
		storeIterator := sdk.KVStorePrefixIterator(store, append(types.StorageStoreKeyPrefix, addressBytes...))
		prefixAddrLen := len(append(types.StorageStoreKeyPrefix, addressBytes...))
		var storage []types.Storage
		for ; storeIterator.Valid(); storeIterator.Next() {
			keyBytes := storeIterator.Key()[prefixAddrLen:]
			var key binary.Word256
			copy(key[:], keyBytes)
			storage = append(storage, types.Storage{Key: key, Value: storeIterator.Value()})
		}
		storeIterator.Close()
		contracts = append(contracts, types.Contract{
			Address: address,
			Code:    code,
			Storage: storage,
			Abi:     abi,
			Meta:    meta,
		})
	}
	return contracts
}

// GetAllMetas gets all metadata for genesis export.
func (k Keeper) GetAllMetas(ctx sdk.Context) []types.Metadata {
	contracts := make([]types.Metadata, 0)
	store := ctx.KVStore(k.key)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaHashStoreKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		metahashBytes := iterator.Key()[len(types.MetaHashStoreKeyPrefix):]
		var metahash acmstate.MetadataHash
		copy(metahash[:], metahashBytes)

		meta, err := k.getMeta(ctx, metahash)
		if err != nil {
			panic(err)
		}
		contracts = append(contracts, types.Metadata{Hash: metahash, Metadata: meta})
	}
	return contracts
}

// AuthKeeper returns keeper's AccountKeeper.
func (k Keeper) AuthKeeper() types.AccountKeeper {
	return k.ak
}

// RegisterGlobalPermissionAcc registers the zero address as the global permission account.
func RegisterGlobalPermissionAcc(ctx sdk.Context, k Keeper) {
	state := k.NewState(ctx)
	st := engine.State{
		CallFrame:  engine.NewCallFrame(state).WithMaxCallStackDepth(0),
		Blockchain: NewBlockChain(ctx, k),
		EventSink:  NewEventSink(ctx),
	}
	gpacc, err := st.CallFrame.GetAccount(acm.GlobalPermissionsAddress)
	if err != nil {
		panic(err)
	}
	if gpacc == nil {
		if err = native.CreateAccount(st.CallFrame, acm.GlobalPermissionsAddress); err != nil {
			panic(err)
		}
	}
	if err := st.Sync(); err != nil {
		panic(err)
	}
}

// SetGasRate sets parameters subspace for gas rate.
func (k Keeper) SetGasRate(ctx sdk.Context, gasRate uint64) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyGasRate, &gasRate)
}

// GetGasRate returns the gas rate in parameters subspace.
func (k *Keeper) GetGasRate(ctx sdk.Context) uint64 {
	var gasRate uint64
	k.paramSpace.Get(ctx, types.ParamStoreKeyGasRate, &gasRate)
	return gasRate
}
