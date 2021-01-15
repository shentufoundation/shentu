package keeper_test

import (
	gobin "encoding/binary"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/certikfoundation/shentu/x/cvm/compile"
	"github.com/hyperledger/burrow/logging"

	"github.com/stretchr/testify/require"
	"github.com/tmthrgd/go-hex"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/cvm/keeper"
	"github.com/certikfoundation/shentu/x/cvm/types"
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
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	t.Run("should allow call on a contract with no code when calling with empty data (transfer)", func(t *testing.T) {
		result, err := app.CvmKeeper.Call(ctx, addrs[0], addrs[1], 10, []byte{}, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, result)
		require.Nil(t, err)
	})

	t.Run("should not allow call on a contract with no code when calling with data (transfer)", func(t *testing.T) {
		result, err := app.CvmKeeper.Call(ctx, addrs[0], addrs[1], 10, []byte{0x00}, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, result)
		require.NotNil(t, err)
		require.Equal(t, types.ErrCodedError(errors.Codes.CodeOutOfBounds), err)
	})

	t.Run("deploy a contract with regular code and call a function in the contract", func(t *testing.T) {
		code, err := hex.DecodeString(Hello55BytecodeString)

		require.Nil(t, err)

		result, err2 := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, result)
		require.Nil(t, err2)

		sayHiCall, _, err := abi.EncodeFunctionCall(
			Hello55AbiJsonString,
			"sayHi",
			keeper.WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)

		// call its function
		newContractAddress := sdk.AccAddress(result)
		result, err2 = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, sayHiCall, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(55))
		require.Nil(t, err2)
	})

}

func TestProperExecution(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	var newContractAddress sdk.AccAddress
	t.Run("deploy a contract with regular code", func(t *testing.T) {
		code, err := hex.DecodeString(BasicTestsBytecodeString)
		require.Nil(t, err)

		result, err2 := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
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
		result, err := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
	})

	t.Run("call a function that takes parameters and ensure it works properly", func(t *testing.T) {
		addSevenAndEightCall, _, err := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"addTwoNumbers",
			keeper.WrapLogger(ctx.Logger()),
			7, 8,
		)
		require.Nil(t, err)
		result, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, addSevenAndEightCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(15))
	})

	t.Run("call a function that should revert and ensure that it reverts", func(t *testing.T) {
		failureFunctionCall, _, err := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"failureFunction",
			keeper.WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, failureFunctionCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, types.ErrCodedError(errors.Codes.ExecutionReverted), err2)
	})

	t.Run("call a contract with junk callcode and ensure it reverts", func(t *testing.T) {
		_, err := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, []byte("Kanye West"), []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err)
		require.Equal(t, types.ErrCodedError(errors.Codes.ExecutionReverted), err)
	})

	t.Run("write to state and ensure it is reflected in updated state", func(t *testing.T) {
		setMyFavoriteNumberCall, _, err2 := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"setMyFavoriteNumber",
			keeper.WrapLogger(ctx.Logger()),
			777,
		)
		require.Nil(t, err2)
		result, err := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, setMyFavoriteNumberCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		result, err2 = app.CvmKeeper.GetStorage(ctx, crypto.MustAddressFromBytes(newContractAddress), binary.Int64ToWord256(0))
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(777))
	})
}

func TestView(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	var newContractAddress sdk.AccAddress
	t.Run("deploy a contract with regular code", func(t *testing.T) {
		code, err := hex.DecodeString(BasicTestsBytecodeString)
		require.Nil(t, err)

		result, err2 := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	t.Run("write to state while in view mode and ensure it is NOT reflected in updated state", func(t *testing.T) {
		setMyFavoriteNumberCall, _, err2 := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"setMyFavoriteNumber",
			keeper.WrapLogger(ctx.Logger()),
			777,
		)
		require.Nil(t, err2)
		result, err := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, setMyFavoriteNumberCall, []*payload.ContractMeta{}, true, false, false)
		require.NotNil(t, err)
		result, err2 = app.CvmKeeper.GetStorage(ctx, crypto.MustAddressFromBytes(newContractAddress), binary.Int64ToWord256(0))
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(34))
	})
}

func TestGasPrice(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	var newContractAddress sdk.AccAddress
	t.Run("deploy contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasTestsBytecodeString)
		require.Nil(t, err2)
		result, err := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	addTwoNumbersCall, _, err := abi.EncodeFunctionCall(
		GasTestsAbiJsonString,
		"addTwoNumbers",
		keeper.WrapLogger(ctx.Logger()),
		3, 5,
	)
	require.Nil(t, err)

	t.Run("add two numbers with not enough gas and see it fail", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		ctx = ctx.WithGasMeter(NewGasMeter(AddTwoNumbersGasCost - 5000))
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, addTwoNumbersCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, err2.Error(), types.ErrCodedError(errors.Codes.InsufficientGas).Error())
	})

	t.Run("add two numbers with the right gas amount", func(t *testing.T) {
		ctx = ctx.WithGasMeter(NewGasMeter(AddTwoNumbersGasCost + 50000))
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, addTwoNumbersCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
	})

	hashMeCall, _, err := abi.EncodeFunctionCall(
		GasTestsAbiJsonString,
		"hashMe",
		keeper.WrapLogger(ctx.Logger()),
		[]byte("abcdefghij"),
	)

	t.Run("hash some bytes with not enough gas and see it fail", func(t *testing.T) {
		require.Nil(t, err)
		ctx = ctx.WithGasMeter(NewGasMeter(HashMeGasCost - 4000))
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, hashMeCall, nil, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, err2, types.ErrCodedError(errors.Codes.InsufficientGas))
	})

	t.Run("hash some bytes with the right gas amount", func(t *testing.T) {
		ctx = ctx.WithGasMeter(NewGasMeter(HashMeGasCost + 50000))
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, hashMeCall, nil, false, false, false)
		require.Nil(t, err2)
	})

	var deployAnotherContractCall []byte
	t.Run("deploy another contract with not enough gas and see it fail", func(t *testing.T) {
		deployAnotherContractCall, _, err = abi.EncodeFunctionCall(
			GasTestsAbiJsonString,
			"deployAnotherContract",
			keeper.WrapLogger(ctx.Logger()),
		)

		require.Nil(t, err)
		ctx = ctx.WithGasMeter(NewGasMeter(DeployAnotherContractGasCost - 150000)) //DeployAnotherContractGasCost - 20))
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, deployAnotherContractCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		require.Equal(t, err2, types.ErrCodedError(errors.Codes.InsufficientGas))
	})

	t.Run("deploy another contract with the right gas amount", func(t *testing.T) {
		ctx = ctx.WithGasMeter(NewGasMeter(DeployAnotherContractGasCost))
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, deployAnotherContractCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
	})
}

func TestGasRefund(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	/*
		Currently there is no way to actually test this because gas refunded has no effect,
		ie no CTK will be deducted from or added to a user's wallet upon a CVM transaction
		taking place.

		So gas refunded is just an internal variable within the CVM module and we cannot check
		a function's return value for it.

		It seems writing an actual test for this will have to wait until the economic document is finalized

		For now however we can run the requisite tests without `require` statements and use
		Printlns to read what is going on internally
	*/

	// defer panic("pointless panic") // UNCOMMENT THIS LINE TO SEE PRINTLNS, they will not appear if all tests run successfully

	var newContractAddress sdk.AccAddress
	fmt.Println("Deploy gas test contract")
	fmt.Println("------------------------")

	t.Run("deploy gas test contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasTestsBytecodeString)
		require.Nil(t, err2)
		result, err := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
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
			keeper.WrapLogger(ctx.Logger()),
			3, 5,
		)
		require.Nil(t, err)
		_, err2 := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, addTwoNumbersCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		/* TODO, check for gas refunded */
	})

	fmt.Println()
	fmt.Println("Deploy gas refund contract")
	fmt.Println("------------------------")

	t.Run("deploy gas refund contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasRefundBytecodeString)
		require.Nil(t, err2)
		result, err := app.CvmKeeper.Call(ctx, addrs[1], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
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
			keeper.WrapLogger(ctx.Logger()),
		)

		require.Nil(t, err)
		_, err2 := app.CvmKeeper.Call(ctx, addrs[1], newContractAddress, 0, iWillRevertCall, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, err2)
		/* TODO, check for gas refunded */
	})

	fmt.Println()
	fmt.Println("Call that fails")
	fmt.Println("------------------------")

	t.Run("run call that fails with wayyy too much gas, should not get a refund", func(t *testing.T) {
		iWillFailCall, _, err := abi.EncodeFunctionCall(
			GasRefundAbiJsonString,
			"iWillFail",
			keeper.WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CvmKeeper.Call(ctx, addrs[1], newContractAddress, 0, iWillFailCall, []*payload.ContractMeta{}, false, false, false)
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
			keeper.WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CvmKeeper.Call(ctx, addrs[1], newContractAddress, 0, deleteFromStorageCall, []*payload.ContractMeta{}, false, false, false)
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
			keeper.WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := app.CvmKeeper.Call(ctx, addrs[1], newContractAddress, 0, dieCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err2)
		/* TODO, ensure that refund took place (half the gas should be refunded) */
	})

}

func TestCTKTransfer(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	t.Run("check to ensure some account1 that has 10000 CTK", func(t *testing.T) {
		balance := app.AccountKeeper.GetAccount(ctx, addrs[0]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(10000))
		balance = app.AccountKeeper.GetAccount(ctx, addrs[1]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(10000))
	})

	var newContractAddress sdk.AccAddress

	t.Run("send it to an smart contract that will send half to account2 and keep half (deploy)", func(t *testing.T) {
		code, err2 := hex.DecodeString(CtkTransferTestBytecodeString)
		require.Nil(t, err2)

		result, err := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)

		require.Nil(t, err)
	})

	t.Run("send it to an smart contract that will send half to account2 and keep half (execute)", func(t *testing.T) {
		bytesToArr, _ := addrs[1].Marshal()

		var toAddr crypto.Address
		copy(toAddr[:], padOrTrim(bytesToArr, 32)[12:32])

		sendToAFriendCall, _, err := abi.EncodeFunctionCall(
			CtkTransferTestAbiJsonString,
			"sendToAFriend",
			keeper.WrapLogger(ctx.Logger()),
			toAddr,
		)

		require.Nil(t, err)
		_, err = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress,
			10000, sendToAFriendCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
	})

	t.Run("check to ensure account1 has 0 CTK", func(t *testing.T) {
		balance := app.AccountKeeper.GetAccount(ctx, addrs[0]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(0))
	})

	t.Run("check to ensure account2 has 15000 CTK", func(t *testing.T) {
		balance := app.AccountKeeper.GetAccount(ctx, addrs[1]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(15000))
	})

	t.Run("check to ensure the smart contract has 5000 CTK", func(t *testing.T) {
		balance := app.AccountKeeper.GetAccount(ctx, newContractAddress).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(5000))
	})

	t.Run("run code in the smart contract to report its own balance", func(t *testing.T) {
		whatsMyBalanceCall, _, err := abi.EncodeFunctionCall(
			CtkTransferTestAbiJsonString,
			"whatsMyBalance",
			keeper.WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		result, err := app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, whatsMyBalanceCall, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(5000))
	})
}

func TestZeroTransfer(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	t.Run("use recycle to send to the community pool", func(t *testing.T) {
		coins := sdk.NewCoins(sdk.NewCoin("uctk", sdk.NewInt(10)))
		err := app.BankKeeper.SendCoins(ctx, addrs[0], crypto.ZeroAddress.Bytes(), coins)
		require.Nil(t, err)
		err = app.CvmKeeper.RecycleCoins(ctx)
		zAcc := app.AccountKeeper.GetAccount(ctx, crypto.ZeroAddress.Bytes())
		err = zAcc.SetCoins(sdk.Coins{})
		require.Nil(t, err)
		app.AccountKeeper.SetAccount(ctx, zAcc)
		require.Nil(t, err)

		communityPool := app.DistrKeeper.GetFeePool(ctx).CommunityPool
		require.Equal(t, sdk.NewDecCoinsFromCoins(coins...), communityPool)
	})
}

func TestStoreLastBlockHash(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	t.Run("store Ctx block hash", func(t *testing.T) {
		cvmk := app.CvmKeeper
		key := app.GetKey(types.StoreKey)

		height1 := ctx.BlockHeight()
		cvmk.StoreLastBlockHash(ctx)
		hash1 := ctx.KVStore(key).Get(types.BlockHashStoreKey(height1))

		coins := sdk.NewCoins(sdk.NewCoin("uctk", sdk.NewInt(10)))
		err := app.BankKeeper.SendCoins(ctx, addrs[0], crypto.ZeroAddress.Bytes(), coins)
		require.Nil(t, err)

		ctx = ctx.WithBlockHeader(abci.Header{LastBlockId: abci.BlockID{Hash: []byte{0x01}}})
		ctx = ctx.WithBlockHeight(2)
		cvmk.StoreLastBlockHash(ctx)
		hash2 := ctx.KVStore(key).Get(types.BlockHashStoreKey(2))
		require.NotEqual(t, hash1, hash2)

		ctx = ctx.WithBlockHeader(abci.Header{LastBlockId: abci.BlockID{Hash: nil}})
		ctx = ctx.WithBlockHeight(3)
		cvmk.StoreLastBlockHash(ctx)
		hash3 := ctx.KVStore(key).Get(types.BlockHashStoreKey(3))
		require.Equal(t, []uint8(nil), hash3)
	})
}

func TestAbi(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	t.Run("store Ctx block hash", func(t *testing.T) {
		cvmk := app.CvmKeeper

		oldAbi := []byte(Hello55AbiJsonString)
		addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
		require.Nil(t, err)
		cvmk.SetAbi(ctx, addr, oldAbi)
		restoreAbi := getAbi(ctx, app.GetKey(types.StoreKey), addr)
		require.Equal(t, oldAbi, restoreAbi)
	})
}

func TestCode(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	t.Run("deploy contract testCheck", func(t *testing.T) {
		cvmk := app.CvmKeeper
		bytecode, err := hex.DecodeString(Hello55BytecodeString)
		require.Nil(t, err)
		addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
		require.Nil(t, err)
		_, err = cvmk.Call(ctx, addrs[0], nil, 0, bytecode, nil, false, false, false)
		require.Nil(t, err)

		seqNum := getAccountSeqNum(ctx, app.AccountKeeper, addrs[0])
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
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	t.Run("store Ctx block hash", func(t *testing.T) {
		cvmk := app.CvmKeeper

		acc := app.AccountKeeper.GetAccount(ctx, addrs[0])
		coins := acc.GetCoins()
		require.Greater(t, coins.AmountOf("uctk").Int64(), int64(9999))
		err := cvmk.Send(ctx, addrs[0], addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", 100000)})
		require.NotNil(t, err)
		err = cvmk.Send(ctx, addrs[0], addrs[1], sdk.Coins{sdk.NewInt64Coin("uatk", 100000)})
		require.NotNil(t, err)

		err = cvmk.Send(ctx, addrs[0], addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", 10)})
		require.Nil(t, err)
	})
}

func TestPrecompiles(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	t.Run("deploy and call native contracts", func(t *testing.T) {
		code, err := hex.DecodeString(testCheckBytecodeString)
		require.Nil(t, err)

		result, err := app.CvmKeeper.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress := sdk.AccAddress(result)

		certAddr, err := sdk.AccAddressFromBech32("certik1q7h5e2gwpykq27etnd5k5fcvc2azpueujqcvap")
		proofAddr, err := sdk.AccAddressFromBech32("certik1m95kfvajw5dmnu9e6h365tqcnazm9auddngklv")
		everythingAddr, err := sdk.AccAddressFromBech32("certik1nzgvd4k34zzf6vk3qevuh5xx8xshg6uy0l8rd5")
		if err != nil {
			panic(err)
		}
		app.CertKeeper.SetCertifier(ctx, cert.Certifier{Address: certAddr})
		auditingCert1, err := cert.NewGeneralCertificate("auditing", "address", certAddr.String(), "WOW", certAddr)
		if err != nil {
			panic(err)
		}
		auditingCert2, err := cert.NewGeneralCertificate("auditing", "address", everythingAddr.String(), "WOW", certAddr)
		if err != nil {
			panic(err)
		}
		proofCert, err := cert.NewGeneralCertificate("proof", "address", proofAddr.String(), "testproof", certAddr)
		if err != nil {
			panic(err)
		}
		proofCert2, err := cert.NewGeneralCertificate("proof", "address", everythingAddr.String(), "testproof", certAddr)
		if err != nil {
			panic(err)
		}
		compCert := cert.NewCompilationCertificate(
			cert.CertificateTypeCompilation, "dummysourcecodehash", "testproof",
			"bch", "dummydesc", certAddr)

		_, err = app.CertKeeper.IssueCertificate(ctx, auditingCert1)
		_, err = app.CertKeeper.IssueCertificate(ctx, auditingCert2)
		_, err = app.CertKeeper.IssueCertificate(ctx, proofCert)
		_, err = app.CertKeeper.IssueCertificate(ctx, proofCert2)
		_, err = app.CertKeeper.IssueCertificate(ctx, compCert)
		if err != nil {
			panic(err)
		}

		callCheckCall, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"callCheck",
			keeper.WrapLogger(ctx.Logger()),
		)
		callCheckNotCertified, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"callCheckNotCertified",
			keeper.WrapLogger(ctx.Logger()),
		)
		proofCheck, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"proofCheck",
			keeper.WrapLogger(ctx.Logger()),
		)
		compCheck, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"compilationCheck",
			keeper.WrapLogger(ctx.Logger()),
		)
		bothCheck, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"proofAndAuditingCheck",
			keeper.WrapLogger(ctx.Logger()),
		)
		require.Nil(t, err)
		result, err = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, callCheckCall, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		result, err = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, callCheckNotCertified, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x00}, result)
		result, err = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, proofCheck, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		result, err = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, compCheck, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		result, err = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, bothCheck, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x01}, result)
		require.Nil(t, err)
	})

	t.Run("deploy and call certify validator native contract", func(t *testing.T) {
		valStr := "certikvalconspub1zcjduepq32v65eegk2yvgzdya5dqnlnc063u7mt3dh66z2xyv9rddgm6t94s4pjeat"
		code, err := hex.DecodeString(testCertifyValidatorString)
		require.Nil(t, err)

		result, err := app.CvmKeeper.Call(ctx, addrs[1], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress := sdk.AccAddress(result)

		certAddr, err := sdk.AccAddressFromBech32("certik1nzgvd4k34zzf6vk3qevuh5xx8xshg6uy0l8rd5")
		if err != nil {
			panic(err)
		}
		app.CertKeeper.SetCertifier(ctx, cert.Certifier{Address: certAddr})
		require.True(t, app.CertKeeper.IsCertifier(ctx, certAddr))

		certifyValidator, _, err := abi.EncodeFunctionCall(
			testCertifyValidatorAbiJsonString,
			"certifyValidator",
			keeper.WrapLogger(ctx.Logger()),
		)
		_, err = app.BankKeeper.AddCoins(ctx, certAddr, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 12345)})
		require.NoError(t, err)
		certAcc := app.AccountKeeper.GetAccount(ctx, certAddr)
		_ = certAcc.SetSequence(1)
		app.AccountKeeper.SetAccount(ctx, certAcc)

		result, err = app.CvmKeeper.Call(ctx, addrs[0], newContractAddress, 0, certifyValidator, []*payload.ContractMeta{}, false, false, false)
		require.Equal(t, []byte{0x00}, result)
		result, err = app.CvmKeeper.Call(ctx, certAddr, newContractAddress, 0, certifyValidator, []*payload.ContractMeta{}, false, false, false)
		fmt.Println(result)
		fmt.Println(err)
		require.Equal(t, []byte{0x01}, result)

		validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, valStr)
		require.Nil(t, err)
		require.True(t, app.CertKeeper.IsValidatorCertified(ctx, validator))
	})
}

func TestEWASM(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	k := app.CvmKeeper

	t.Run("test eWASM contracts", func(t *testing.T) {
		basename, workDir, _ := compile.ResolveFilename("tests/for-r.wasm")
		testeWASMForStringResp, err := compile.BytecodeEVM(basename, workDir, "tests/for.abi", logging.NewNoopLogger())
		require.Nil(t, err)

		code, err := hex.DecodeString(testeWASMForStringResp.Objects[0].Contract.Code())
		require.Nil(t, err)

		result, err := k.Call(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, true, true)
		require.Nil(t, err)
		require.NotNil(t, result)

		contractAddr := sdk.AccAddress(result)
		fmt.Println(contractAddr.String())

		logger := logging.NewNoopLogger()
		callcode, _, err := abi.EncodeFunctionCall(testeWASMForAbiJsonString, "multiply", logger, "3", "2")
		result, err = k.Call(ctx, addrs[0], contractAddr, 0, callcode, nil, false, false, false)
		require.Nil(t, err)
		fmt.Println(result)
	})
}
