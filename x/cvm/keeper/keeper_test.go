package keeper_test

import (
	gobin "encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/txs/payload"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	shentuapp "github.com/certikfoundation/shentu/v2/app"
	"github.com/certikfoundation/shentu/v2/common"
	certtypes "github.com/certikfoundation/shentu/v2/x/cert/types"
	. "github.com/certikfoundation/shentu/v2/x/cvm/keeper"
	"github.com/certikfoundation/shentu/v2/x/cvm/types"
)

var (
	uCTKAmount = sdk.NewInt(1005).MulRaw(common.MicroUnit)
)

// getAbi returns the abi at the given address.
// NOTE: Emulates the unexported function in the module.
func getAbi(ctx sdk.Context, key sdk.StoreKey, address crypto.Address) []byte {
	return ctx.KVStore(key).Get(types.AbiStoreKey(address))
}

func getCode(ctx sdk.Context, key sdk.StoreKey, address crypto.Address) []byte {
	return ctx.KVStore(key).Get(types.CodeStoreKey(address))
}

func getAddressMeta(ctx sdk.Context, key sdk.StoreKey, address crypto.Address) []byte {
	return ctx.KVStore(key).Get(types.AddressMetaStoreKey(address))
}

// getAccountSeqNum returns the account sequence number.
// NOTE: Emulates the unexported function in the module.
func getAccountSeqNum(ctx sdk.Context, ak authKeeper.AccountKeeper, address sdk.AccAddress) []byte {
	callerAcc := ak.GetAccount(ctx, address)
	callerSequence := callerAcc.GetSequence()
	accountByte := make([]byte, 8)
	gobin.LittleEndian.PutUint64(accountByte, callerSequence)
	return accountByte
}

// padOrTrim returns (size) bytes from app (bb)
// Short bb gets zeros prefixed, Long bb gets left/MSB bits trimmed
func padOrTrim(bb []byte, size int) []byte {
	l := len(bb)
	if l == size {
		return bb
	}
	if l > size {
		return bb[l-size:]
	}
	tmp := make([]byte, size)
	copy(tmp[size-l:], bb)
	return tmp
}

func TestContractCreation(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	t.Run("should allow call on a contract with no code when calling with empty data (transfer)", func(t *testing.T) {
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], addrs[1], 10, []byte{}, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, result)
		require.Nil(t, err)
	})

	t.Run("should not allow call on a contract with no code when calling with data (transfer)", func(t *testing.T) {
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], addrs[1], 10, []byte{0x00}, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, result)
		require.NotNil(t, err)
		require.Equal(t, types.ErrCodedError(errors.Codes.CodeOutOfBounds), err)
	})

	t.Run("deploy a contract with regular code and call a function in the contract", func(t *testing.T) {
		code, err := hex.DecodeString(Hello55BytecodeString)

		require.Nil(t, err)

		result, err2 := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, result)
		require.Nil(t, err2)

		sayHiCall, _, err := abi.EncodeFunctionCall(
			Hello55AbiJsonString,
			"sayHi",
			WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)

		// call its function
		newContractAddress := sdk.AccAddress(result)
		result, err2 = app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, sayHiCall, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(55))
		require.Nil(t, err2)
	})

}

func TestProperExecution(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	var newContractAddress sdk.AccAddress
	t.Run("deploy a contract with regular code", func(t *testing.T) {
		code, err := hex.DecodeString(BasicTestsBytecodeString)
		require.Nil(t, err)

		result, err2 := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	fmt.Println("")

	t.Run("deploy second contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(Hello55BytecodeString)
		require.Nil(t, err2)
		acc := app.AccountKeeper.GetAccount(ctx, addrs[0])
		_ = acc.SetSequence(acc.GetSequence() + 1)
		app.AccountKeeper.SetAccount(ctx, acc)
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
	})

	t.Run("call a function that takes parameters and ensure it works properly", func(t *testing.T) {
		addSevenAndEightCall, _, err := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"addTwoNumbers",
			WrapLogger(ctx.Logger()),
			7, 8,
		)
		require.Nil(t, err)
		result, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, addSevenAndEightCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(15))
	})

	t.Run("call a function that should revert and ensure that it reverts", func(t *testing.T) {
		failureFunctionCall, _, err := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"failureFunction",
			WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, failureFunctionCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, types.ErrCodedError(errors.Codes.ExecutionReverted), err2)
	})

	t.Run("call a contract with junk callcode and ensure it reverts", func(t *testing.T) {
		_, err := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, []byte("Kanye West"), []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err)
		require.Equal(t, types.ErrCodedError(errors.Codes.ExecutionReverted), err)
	})

	t.Run("write to state and ensure it is reflected in updated state", func(t *testing.T) {
		setMyFavoriteNumberCall, _, err2 := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"setMyFavoriteNumber",
			WrapLogger(ctx.Logger()),
			777,
		)
		require.Nil(t, err2)
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, setMyFavoriteNumberCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		result, err2 = app.CVMKeeper.GetStorage(ctx, crypto.MustAddressFromBytes(newContractAddress), binary.Int64ToWord256(0))
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(777))
	})
}

func TestView(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	var newContractAddress sdk.AccAddress
	t.Run("deploy a contract with regular code", func(t *testing.T) {
		code, err := hex.DecodeString(BasicTestsBytecodeString)
		require.Nil(t, err)

		result, err2 := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	t.Run("write to state while in view mode and ensure it is NOT reflected in updated state", func(t *testing.T) {
		setMyFavoriteNumberCall, _, err2 := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"setMyFavoriteNumber",
			WrapLogger(ctx.Logger()),
			777,
		)
		require.Nil(t, err2)
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, setMyFavoriteNumberCall, []*payload.ContractMeta{}, true, false, false)
		require.NotNil(t, err)
		result, err2 = app.CVMKeeper.GetStorage(ctx, crypto.MustAddressFromBytes(newContractAddress), binary.Int64ToWord256(0))
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(34))
	})
}

func TestGasPrice(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	var newContractAddress sdk.AccAddress
	t.Run("deploy contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasTestsBytecodeString)
		require.Nil(t, err2)
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	addTwoNumbersCall, _, err := abi.EncodeFunctionCall(
		GasTestsAbiJsonString,
		"addTwoNumbers",
		WrapLogger(ctx.Logger()),
		3, 5,
	)
	require.Nil(t, err)

	t.Run("add two numbers with not enough gas and see it fail", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(AddTwoNumbersGasCost - 5000))
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, addTwoNumbersCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, err2.Error(), types.ErrCodedError(errors.Codes.InsufficientGas).Error())
	})

	t.Run("add two numbers with the right gas amount", func(t *testing.T) {
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(AddTwoNumbersGasCost + 50000))
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, addTwoNumbersCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
	})

	hashMeCall, _, err := abi.EncodeFunctionCall(
		GasTestsAbiJsonString,
		"hashMe",
		WrapLogger(ctx.Logger()),
		[]byte("abcdefghij"),
	)

	t.Run("hash some bytes with not enough gas and see it fail", func(t *testing.T) {
		require.Nil(t, err)
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(HashMeGasCost - 1500))
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, hashMeCall, nil, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, err2, types.ErrCodedError(errors.Codes.InsufficientGas))
	})

	t.Run("hash some bytes with the right gas amount", func(t *testing.T) {
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(HashMeGasCost + 50000))
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, hashMeCall, nil, false, false, false)
		require.Nil(t, err2)
	})

	var deployAnotherContractCall []byte
	t.Run("deploy another contract with not enough gas and see it fail", func(t *testing.T) {
		deployAnotherContractCall, _, err = abi.EncodeFunctionCall(
			GasTestsAbiJsonString,
			"deployAnotherContract",
			WrapLogger(ctx.Logger()),
		)

		require.Nil(t, err)
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(DeployAnotherContractGasCost - 150000)) //DeployAnotherContractGasCost - 20))
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, deployAnotherContractCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, err2, types.ErrCodedError(errors.Codes.InsufficientGas))
	})

	t.Run("deploy another contract with the right gas amount", func(t *testing.T) {
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(DeployAnotherContractGasCost))
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, deployAnotherContractCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
	})
}

func TestGasRefund(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	var newContractAddress sdk.AccAddress
	fmt.Println("Deploy gas test contract")
	fmt.Println("------------------------")

	t.Run("deploy gas test contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasTestsBytecodeString)
		require.Nil(t, err2)
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	fmt.Println()
	fmt.Println("Add two numbers w too much gas")
	fmt.Println("------------------------")

	t.Run("add two numbers with wayyy too much gas and hope to get a refund", func(t *testing.T) {
		addTwoNumbersCall, _, err := abi.EncodeFunctionCall(
			GasTestsAbiJsonString,
			"addTwoNumbers",
			WrapLogger(ctx.Logger()),
			3, 5,
		)
		require.Nil(t, err)
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, addTwoNumbersCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		/* TODO, check for gas refunded */
	})

	fmt.Println()
	fmt.Println("Deploy gas refund contract")
	fmt.Println("------------------------")

	t.Run("deploy gas refund contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasRefundBytecodeString)
		require.Nil(t, err2)
		result, err := app.CVMKeeper.Tx(ctx, addrs[1], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	fmt.Println()
	fmt.Println("Add two numbers and revert")
	fmt.Println("------------------------")

	t.Run("add two numbers with wayyy too much gas, followed by a revert, and hope to get a refund", func(t *testing.T) {
		iWillRevertCall, _, err := abi.EncodeFunctionCall(
			GasRefundAbiJsonString,
			"iWillRevert",
			WrapLogger(ctx.Logger()),
		)

		require.Nil(t, err)
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[1], newContractAddress, 0, iWillRevertCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		/* TODO, check for gas refunded */
	})

	fmt.Println()
	fmt.Println("Tx that fails")
	fmt.Println("------------------------")

	t.Run("run call that fails with wayyy too much gas, should not get a refund", func(t *testing.T) {
		iWillFailCall, _, err := abi.EncodeFunctionCall(
			GasRefundAbiJsonString,
			"iWillFail",
			WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[1], newContractAddress, 0, iWillFailCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		/* TODO, ensure that no refund took place */
	})

	fmt.Println()
	fmt.Println("Delete from storage")
	fmt.Println("------------------------")

	t.Run("delete from storage, should get a refund", func(t *testing.T) {
		deleteFromStorageCall, _, err := abi.EncodeFunctionCall(
			GasRefundAbiJsonString,
			"deleteFromStorage",
			WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[1], newContractAddress, 0, deleteFromStorageCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		/* TODO, ensure that refund took place (half the gas should be refunded) */
	})

	fmt.Println()
	fmt.Println("Self destruct")
	fmt.Println("------------------------")

	t.Run("self destruct, should get a refund", func(t *testing.T) {
		dieCall, _, err := abi.EncodeFunctionCall(
			GasRefundAbiJsonString,
			"die",
			WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CVMKeeper.Tx(ctx, addrs[1], newContractAddress, 0, dieCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		/* TODO, ensure that refund took place (half the gas should be refunded) */
	})

}

func TestCTKTransfer(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	t.Run("check to ensure some account1 that has 10000 CTK", func(t *testing.T) {
		balance := app.BankKeeper.GetAllBalances(ctx, addrs[0]).AmountOf(app.StakingKeeper.BondDenom(ctx))
		require.Equal(t, balance, sdk.NewInt(80000000000))
		balance = app.BankKeeper.GetAllBalances(ctx, addrs[1]).AmountOf(app.StakingKeeper.BondDenom(ctx))
		require.Equal(t, balance, sdk.NewInt(80000000000))
	})

	var newContractAddress sdk.AccAddress

	t.Run("send it to an smart contract that will send half to account2 and keep half (deploy)", func(t *testing.T) {
		code, err2 := hex.DecodeString(CtkTransferTestBytecodeString)
		require.Nil(t, err2)

		result, err := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, result)
		newContractAddress = result

		require.Nil(t, err)
	})

	t.Run("send it to an smart contract that will send half to account2 and keep half (execute)", func(t *testing.T) {
		bytesToArr, _ := addrs[1].Marshal()

		var toAddr crypto.Address
		copy(toAddr[:], padOrTrim(bytesToArr, 32)[12:32])

		sendToAFriendCall, _, err := abi.EncodeFunctionCall(
			CtkTransferTestAbiJsonString,
			"sendToAFriend",
			WrapLogger(ctx.Logger()),
			toAddr,
		)

		require.Nil(t, err)
		_, err = app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress,
			10000, sendToAFriendCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
	})

	t.Run("check to ensure account1 has 0", func(t *testing.T) {
		balance := app.BankKeeper.GetAllBalances(ctx, addrs[0]).AmountOf(app.StakingKeeper.BondDenom(ctx))
		require.Equal(t, balance, sdk.NewInt(79999990000))
	})

	t.Run("check to ensure account2 has 15000", func(t *testing.T) {
		balance := app.BankKeeper.GetAllBalances(ctx, addrs[1]).AmountOf(app.StakingKeeper.BondDenom(ctx))
		require.Equal(t, balance, sdk.NewInt(80000005000))
	})

	t.Run("check to ensure the smart contract has 5000", func(t *testing.T) {

		balance := app.BankKeeper.GetAllBalances(ctx, newContractAddress).AmountOf(app.StakingKeeper.BondDenom(ctx))
		require.Equal(t, balance, sdk.NewInt(5000))
	})

	t.Run("run code in the smart contract to report its own balance", func(t *testing.T) {
		whatsMyBalanceCall, _, err := abi.EncodeFunctionCall(
			CtkTransferTestAbiJsonString,
			"whatsMyBalance",
			WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		result, err := app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, whatsMyBalanceCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(5000))
	})
}

func TestStoreLastBlockHash(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	cvmk := app.CVMKeeper

	blockChain := NewBlockChain(ctx, cvmk)

	t.Run("store Ctx block hash", func(t *testing.T) {
		ctx = ctx.WithBlockHeader(tmproto.Header{LastBlockId: tmproto.BlockID{Hash: []byte{0x00}}, Height: 1})
		cvmk.StoreLastBlockHash(ctx)
		hash1, _ := blockChain.BlockHash(1)

		coins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), sdk.NewInt(10)))
		err := app.BankKeeper.SendCoins(ctx, addrs[0], addrs[1], coins)
		require.Nil(t, err)

		ctx = ctx.WithBlockHeader(tmproto.Header{LastBlockId: tmproto.BlockID{Hash: []byte{0x01}}, Height: 2})
		blockChain = NewBlockChain(ctx, cvmk)
		cvmk.StoreLastBlockHash(ctx)
		hash2, err := blockChain.BlockHash(2)
		require.Nil(t, err)
		require.NotEqual(t, hash1, hash2)

		ctx = ctx.WithBlockHeader(tmproto.Header{LastBlockId: tmproto.BlockID{Hash: nil}, Height: 3})
		cvmk.StoreLastBlockHash(ctx)
		hash3, _ := blockChain.BlockHash(3)
		require.Equal(t, []uint8(nil), hash3)
	})
}

func TestAbi(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	t.Run("store Ctx block hash", func(t *testing.T) {
		ctx := ctx
		cvmk := app.CVMKeeper

		oldAbi := []byte(Hello55AbiJsonString)
		addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
		require.Nil(t, err)
		cvmk.SetAbi(ctx, addr, oldAbi)
		restoreAbi := cvmk.GetAbi(ctx, addr)
		require.Equal(t, oldAbi, restoreAbi)
	})
}

func TestCode(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	t.Run("deploy contract testCheck", func(t *testing.T) {
		ctx := ctx
		cvmk := app.CVMKeeper
		bytecode, err := hex.DecodeString(Hello55BytecodeString)
		require.Nil(t, err)
		addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
		require.Nil(t, err)
		_, err = cvmk.Tx(ctx, addrs[0], nil, 0, bytecode, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)

		seqNum := cvmk.GetAccountSeqNum(ctx, addrs[0])
		calleeAddr := crypto.NewContractAddress(addr, seqNum)

		code, err := cvmk.GetCode(ctx, calleeAddr)
		require.Nil(t, err)
		require.Greater(t, len(bytecode), len(code))

		// doesn't contain constructor after deployment
		require.Equal(t, bytecode[:18], code[:18])
		require.Equal(t, bytecode[48:], code[18:])
	})
}

func TestSend(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	t.Run("test sending through CVM", func(t *testing.T) {
		ctx := ctx
		cvmk := app.CVMKeeper

		coins := app.BankKeeper.GetAllBalances(ctx, addrs[0])
		fmt.Println(coins.String())
		require.Greater(t, coins.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64(), int64(9999))
		err := cvmk.Send(ctx, addrs[0], addrs[1], sdk.Coins{sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 80000000001)})
		require.NotNil(t, err)
		err = cvmk.Send(ctx, addrs[0], addrs[1], sdk.Coins{sdk.NewInt64Coin("uatk", 100000)})
		require.NotNil(t, err)

		err = cvmk.Send(ctx, addrs[0], addrs[1], sdk.Coins{sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 10)})
		require.Nil(t, err)
	})
}

func TestPrecompiles(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(80000*1e6))

	t.Run("deploy and call native contracts", func(t *testing.T) {
		code, err := hex.DecodeString(TestCheckBytecodeString)
		require.Nil(t, err)

		result, err := app.CVMKeeper.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress := sdk.AccAddress(result)

		certAddr, err := sdk.AccAddressFromBech32("cosmos17w5kw28te7r5vn4qu08hu6a4crcvwrrgzmsrrn")
		proofAddr, err := sdk.AccAddressFromBech32("cosmos1r60hj2xaxn79qth4pkjm9t27l985xfsmnz9paw")
		everythingAddr, err := sdk.AccAddressFromBech32("cosmos1xxkueklal9vejv9unqu80w9vptyepfa95pd53u")
		if err != nil {
			panic(err)
		}
		app.CertKeeper.SetCertifier(ctx, certtypes.Certifier{Address: certAddr.String()})
		auditingCert1, err := certtypes.NewCertificate("auditing", certAddr.String(), "", "", "WOW", certAddr)
		if err != nil {
			panic(err)
		}
		auditingCert2, err := certtypes.NewCertificate("auditing", everythingAddr.String(), "", "", "WOW", certAddr)
		if err != nil {
			panic(err)
		}
		proofCert, err := certtypes.NewCertificate("proof", proofAddr.String(), "", "", "testproof", certAddr)
		if err != nil {
			panic(err)
		}
		proofCert2, err := certtypes.NewCertificate("proof", everythingAddr.String(), "", "", "testproof", certAddr)
		if err != nil {
			panic(err)
		}
		compCert, err := certtypes.NewCertificate("compilation", "dummysourcecodehash", "testproof",
			"bch", "dummydesc", certAddr)
		if err != nil {
			panic(err)
		}

		_, err = app.CertKeeper.IssueCertificate(ctx, auditingCert1)
		_, err = app.CertKeeper.IssueCertificate(ctx, auditingCert2)
		_, err = app.CertKeeper.IssueCertificate(ctx, proofCert)
		_, err = app.CertKeeper.IssueCertificate(ctx, proofCert2)
		_, err = app.CertKeeper.IssueCertificate(ctx, compCert)
		if err != nil {
			panic(err)
		}

		callCheckCall, _, err := abi.EncodeFunctionCall(
			TestCheckAbiJsonString,
			"callCheck",
			WrapLogger(ctx.Logger()),
		)
		callCheckNotCertified, _, err := abi.EncodeFunctionCall(
			TestCheckAbiJsonString,
			"callCheckNotCertified",
			WrapLogger(ctx.Logger()),
		)
		proofCheck, _, err := abi.EncodeFunctionCall(
			TestCheckAbiJsonString,
			"proofCheck",
			WrapLogger(ctx.Logger()),
		)
		compCheck, _, err := abi.EncodeFunctionCall(
			TestCheckAbiJsonString,
			"compilationCheck",
			WrapLogger(ctx.Logger()),
		)
		bothCheck, _, err := abi.EncodeFunctionCall(
			TestCheckAbiJsonString,
			"proofAndAuditingCheck",
			WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		result, err = app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, callCheckCall, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		result, err = app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, callCheckNotCertified, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x00}, result)
		result, err = app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, proofCheck, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		result, err = app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, compCheck, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		result, err = app.CVMKeeper.Tx(ctx, addrs[0], newContractAddress, 0, bothCheck, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		require.Nil(t, err)
	})
}
