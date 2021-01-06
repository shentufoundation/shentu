package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm/abi"

	"github.com/certikfoundation/shentu/common"
	cert "github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/cvm/types"
)

func TestContractCreation(t *testing.T) {
	input := CreateTestInput(t)

	t.Run("should allow call on a contract with no code when calling with empty data (transfer)", func(t *testing.T) {
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, result)
		require.Nil(t, err)
	})

	t.Run("should not allow call on a contract with no code when calling with data (transfer)", func(t *testing.T) {
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, result)
		require.NotNil(t, err)
		require.Equal(t, types.ErrCodedError(errors.Codes.CodeOutOfBounds), err)
	})

	t.Run("deploy a contract with regular code and call a function in the contract", func(t *testing.T) {
		code, err := hex.DecodeString(Hello55BytecodeString)

		require.Nil(t, err)

		result, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, result)
		require.Nil(t, err2)

		sayHiCall, _, err := abi.EncodeFunctionCall(
			Hello55AbiJsonString,
			"sayHi",
			WrapLogger(input.Ctx.Logger()),
		)
		require.Nil(t, err)

		// call its function
		newContractAddress := sdk.AccAddress(result)
		result, err2 = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(55))
		require.Nil(t, err2)
	})

}

func TestProperExecution(t *testing.T) {
	input := CreateTestInput(t)

	var newContractAddress sdk.AccAddress
	t.Run("deploy a contract with regular code", func(t *testing.T) {
		code, err := hex.DecodeString(BasicTestsBytecodeString)
		require.Nil(t, err)

		result, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err2)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	fmt.Println("")

	t.Run("deploy second contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(Hello55BytecodeString)
		require.Nil(t, err2)
		acc := input.AccountKeeper.GetAccount(input.Ctx, Addrs[0])
		_ = acc.SetSequence(acc.GetSequence() + 1)
		input.AccountKeeper.SetAccount(input.Ctx, acc)
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err)
		require.NotNil(t, result)
	})

	t.Run("call a function that takes parameters and ensure it works properly", func(t *testing.T) {
		addSevenAndEightCall, _, err := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"addTwoNumbers",
			WrapLogger(input.Ctx.Logger()),
			7, 8,
		)
		require.Nil(t, err)
		result, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err2)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(15))
	})

	t.Run("call a function that should revert and ensure that it reverts", func(t *testing.T) {
		failureFunctionCall, _, err := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"failureFunction",
			WrapLogger(input.Ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, err2)
		require.Equal(t, types.ErrCodedError(errors.Codes.ExecutionReverted), err2)
	})

	t.Run("call a contract with junk callcode and ensure it reverts", func(t *testing.T) {
		_, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, err)
		require.Equal(t, types.ErrCodedError(errors.Codes.ExecutionReverted), err)
	})

	t.Run("write to state and ensure it is reflected in updated state", func(t *testing.T) {
		setMyFavoriteNumberCall, _, err2 := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"setMyFavoriteNumber",
			WrapLogger(input.Ctx.Logger()),
			777,
		)
		require.Nil(t, err2)
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err)
		result, err2 = input.CvmKeeper.GetStorage(input.Ctx, crypto.MustAddressFromBytes(newContractAddress), binary.Int64ToWord256(0))
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(777))
	})
}

func TestView(t *testing.T) {
	input := CreateTestInput(t)

	var newContractAddress sdk.AccAddress
	t.Run("deploy a contract with regular code", func(t *testing.T) {
		code, err := hex.DecodeString(BasicTestsBytecodeString)
		require.Nil(t, err)

		result, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err2)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	t.Run("write to state while in view mode and ensure it is NOT reflected in updated state", func(t *testing.T) {
		setMyFavoriteNumberCall, _, err2 := abi.EncodeFunctionCall(
			BasicTestsAbiJsonString,
			"setMyFavoriteNumber",
			WrapLogger(input.Ctx.Logger()),
			777,
		)
		require.Nil(t, err2)
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, err)
		result, err2 = input.CvmKeeper.GetStorage(input.Ctx, crypto.MustAddressFromBytes(newContractAddress), binary.Int64ToWord256(0))
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(34))
	})
}

func TestGasPrice(t *testing.T) {
	input := CreateTestInput(t)
	var newContractAddress sdk.AccAddress
	t.Run("deploy contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasTestsBytecodeString)
		require.Nil(t, err2)
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)
	})

	addTwoNumbersCall, _, err := abi.EncodeFunctionCall(
		GasTestsAbiJsonString,
		"addTwoNumbers",
		WrapLogger(input.Ctx.Logger()),
		3, 5,
	)
	require.Nil(t, err)

	t.Run("add two numbers with not enough gas and see it fail", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		input.Ctx = input.Ctx.WithGasMeter(NewGasMeter(AddTwoNumbersGasCost - 5000))
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, err2)
		require.Equal(t, err2.Error(), types.ErrCodedError(errors.Codes.InsufficientGas).Error())
	})

	t.Run("add two numbers with the right gas amount", func(t *testing.T) {
		input.Ctx = input.Ctx.WithGasMeter(NewGasMeter(AddTwoNumbersGasCost + 50000))
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err2)
	})

	hashMeCall, _, err := abi.EncodeFunctionCall(
		GasTestsAbiJsonString,
		"hashMe",
		WrapLogger(input.Ctx.Logger()),
		[]byte("abcdefghij"),
	)

	t.Run("hash some bytes with not enough gas and see it fail", func(t *testing.T) {
		require.Nil(t, err)
		input.Ctx = input.Ctx.WithGasMeter(NewGasMeter(HashMeGasCost - 4000))
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, err2)
		require.Equal(t, err2, types.ErrCodedError(errors.Codes.InsufficientGas))
	})

	t.Run("hash some bytes with the right gas amount", func(t *testing.T) {
		input.Ctx = input.Ctx.WithGasMeter(NewGasMeter(HashMeGasCost + 50000))
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err2)
	})

	var deployAnotherContractCall []byte
	t.Run("deploy another contract with not enough gas and see it fail", func(t *testing.T) {
		deployAnotherContractCall, _, err = abi.EncodeFunctionCall(
			GasTestsAbiJsonString,
			"deployAnotherContract",
			WrapLogger(input.Ctx.Logger()),
		)

		require.Nil(t, err)
		input.Ctx = input.Ctx.WithGasMeter(NewGasMeter(DeployAnotherContractGasCost - 150000)) //DeployAnotherContractGasCost - 20))
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, err2)
		require.Equal(t, err2, types.ErrCodedError(errors.Codes.InsufficientGas))
	})

	t.Run("deploy another contract with the right gas amount", func(t *testing.T) {
		input.Ctx = input.Ctx.WithGasMeter(NewGasMeter(DeployAnotherContractGasCost))
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err2)
	})
}

func TestGasRefund(t *testing.T) {
	input := CreateTestInput(t)
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
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
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
			WrapLogger(input.Ctx.Logger()),
			3, 5,
		)
		require.Nil(t, err)
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err2)
		/* TODO, check for gas refunded */
	})

	fmt.Println()
	fmt.Println("Deploy gas refund contract")
	fmt.Println("------------------------")

	t.Run("deploy gas refund contract", func(t *testing.T) {
		code, err2 := hex.DecodeString(GasRefundBytecodeString)
		require.Nil(t, err2)
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[1], false)
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
			WrapLogger(input.Ctx.Logger()),
		)

		require.Nil(t, err)
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[1], false)
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
			WrapLogger(input.Ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[1], false)
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
			WrapLogger(input.Ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[1], false)
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
			WrapLogger(input.Ctx.Logger()),
		)
		require.Nil(t, err)
		_, err2 := input.CvmKeeper.Call(input.Ctx, Addrs[1], false)
		require.Nil(t, err2)
		/* TODO, ensure that refund took place (half the gas should be refunded) */
	})

}

func TestCTKTransfer(t *testing.T) {
	input := CreateTestInput(t)

	t.Run("check to ensure some account1 that has 10000 CTK", func(t *testing.T) {
		balance := input.AccountKeeper.GetAccount(input.Ctx, Addrs[0]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(10000))
		balance = input.AccountKeeper.GetAccount(input.Ctx, Addrs[1]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(10000))
	})

	var newContractAddress sdk.AccAddress

	t.Run("send it to an smart contract that will send half to account2 and keep half (deploy)", func(t *testing.T) {
		code, err2 := hex.DecodeString(CtkTransferTestBytecodeString)
		require.Nil(t, err2)

		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.NotNil(t, result)
		newContractAddress = sdk.AccAddress(result)

		require.Nil(t, err)
	})

	t.Run("send it to an smart contract that will send half to account2 and keep half (execute)", func(t *testing.T) {
		bytesToArr, _ := Addrs[1].Marshal()

		var toAddr crypto.Address
		copy(toAddr[:], padOrTrim(bytesToArr, 32)[12:32])

		sendToAFriendCall, _, err := abi.EncodeFunctionCall(
			CtkTransferTestAbiJsonString,
			"sendToAFriend",
			WrapLogger(input.Ctx.Logger()),
			toAddr,
		)

		require.Nil(t, err)
		_, err = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err)
	})

	t.Run("check to ensure account1 has 0 CTK", func(t *testing.T) {
		balance := input.AccountKeeper.GetAccount(input.Ctx, Addrs[0]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(0))
	})

	t.Run("check to ensure account2 has 15000 CTK", func(t *testing.T) {
		balance := input.AccountKeeper.GetAccount(input.Ctx, Addrs[1]).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(15000))
	})

	t.Run("check to ensure the smart contract has 5000 CTK", func(t *testing.T) {
		balance := input.AccountKeeper.GetAccount(input.Ctx, newContractAddress).GetCoins().AmountOf(common.MicroCTKDenom)
		require.Equal(t, balance, sdk.NewInt(5000))
	})

	t.Run("run code in the smart contract to report its own balance", func(t *testing.T) {
		whatsMyBalanceCall, _, err := abi.EncodeFunctionCall(
			CtkTransferTestAbiJsonString,
			"whatsMyBalance",
			WrapLogger(input.Ctx.Logger()),
		)
		require.Nil(t, err)
		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err)
		require.Equal(t, new(big.Int).SetBytes(result).Int64(), int64(5000))
	})
}

func TestZeroTransfer(t *testing.T) {
	input := CreateTestInput(t)

	t.Run("use recycle to send to the community pool", func(t *testing.T) {
		coins := sdk.NewCoins(sdk.NewCoin("uctk", sdk.NewInt(10)))
		err := input.BankKeeper.SendCoins(input.Ctx, Addrs[0], crypto.ZeroAddress.Bytes(), coins)
		require.Nil(t, err)
		err = input.CvmKeeper.RecycleCoins(input.Ctx)
		zAcc := input.AccountKeeper.GetAccount(input.Ctx, crypto.ZeroAddress.Bytes())
		err = zAcc.SetCoins(sdk.Coins{})
		require.Nil(t, err)
		input.AccountKeeper.SetAccount(input.Ctx, zAcc)
		require.Nil(t, err)
		require.Equal(t, *input.DistrKeeper.CommunityPool, coins)
	})
}

func TestStoreLastBlockHash(t *testing.T) {
	input := CreateTestInput(t)

	t.Run("store Ctx block hash", func(t *testing.T) {
		ctx := input.Ctx
		cvmk := input.CvmKeeper

		height1 := ctx.BlockHeight()
		cvmk.StoreLastBlockHash(ctx)
		hash1 := ctx.KVStore(cvmk.key).Get(types.BlockHashStoreKey(height1))

		coins := sdk.NewCoins(sdk.NewCoin("uctk", sdk.NewInt(10)))
		err := input.BankKeeper.SendCoins(input.Ctx, Addrs[0], crypto.ZeroAddress.Bytes(), coins)
		require.Nil(t, err)

		ctx = ctx.WithBlockHeader(abci.Header{LastBlockId: abci.BlockID{Hash: []byte{0x01}}})
		ctx = ctx.WithBlockHeight(2)
		cvmk.StoreLastBlockHash(ctx)
		hash2 := ctx.KVStore(cvmk.key).Get(types.BlockHashStoreKey(2))
		require.NotEqual(t, hash1, hash2)

		ctx = ctx.WithBlockHeader(abci.Header{LastBlockId: abci.BlockID{Hash: nil}})
		ctx = ctx.WithBlockHeight(3)
		cvmk.StoreLastBlockHash(ctx)
		hash3 := ctx.KVStore(cvmk.key).Get(types.BlockHashStoreKey(3))
		require.Equal(t, []uint8(nil), hash3)
	})
}

func TestAbi(t *testing.T) {
	input := CreateTestInput(t)

	t.Run("store Ctx block hash", func(t *testing.T) {
		ctx := input.Ctx
		cvmk := input.CvmKeeper

		oldAbi := []byte(Hello55AbiJsonString)
		addr, err := crypto.AddressFromBytes(Addrs[0].Bytes())
		require.Nil(t, err)
		cvmk.SetAbi(ctx, addr, oldAbi)
		restoreAbi := cvmk.getAbi(ctx, addr)
		require.Equal(t, oldAbi, restoreAbi)
	})
}

func TestCode(t *testing.T) {
	input := CreateTestInput(t)

	t.Run("deploy contract testCheck", func(t *testing.T) {
		ctx := input.Ctx
		cvmk := input.CvmKeeper
		bytecode, err := hex.DecodeString(Hello55BytecodeString)
		require.Nil(t, err)
		addr, err := crypto.AddressFromBytes(Addrs[0].Bytes())
		require.Nil(t, err)
		_, err = cvmk.Call(ctx, Addrs[0], false)
		require.Nil(t, err)

		seqNum := cvmk.getAccountSeqNum(ctx, Addrs[0])
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
	input := CreateTestInput(t)
	t.Run("store Ctx block hash", func(t *testing.T) {
		ctx := input.Ctx
		cvmk := input.CvmKeeper

		acc := input.AccountKeeper.GetAccount(ctx, Addrs[0])
		coins := acc.GetCoins()
		require.Greater(t, coins.AmountOf("uctk").Int64(), int64(9999))
		err := cvmk.Send(ctx, Addrs[0], Addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", 100000)})
		require.NotNil(t, err)
		err = cvmk.Send(ctx, Addrs[0], Addrs[1], sdk.Coins{sdk.NewInt64Coin("uatk", 100000)})
		require.NotNil(t, err)

		err = cvmk.Send(ctx, Addrs[0], Addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", 10)})
		require.Nil(t, err)
	})
}

// padOrTrim returns (size) bytes from input (bb)
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

func TestPrecompiles(t *testing.T) {
	input := CreateTestInput(t)
	t.Run("deploy and call native contracts", func(t *testing.T) {
		code, err := hex.DecodeString(testCheckBytecodeString)
		require.Nil(t, err)

		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress := sdk.AccAddress(result)

		certAddr, err := sdk.AccAddressFromBech32("certik1q7h5e2gwpykq27etnd5k5fcvc2azpueujqcvap")
		proofAddr, err := sdk.AccAddressFromBech32("certik1m95kfvajw5dmnu9e6h365tqcnazm9auddngklv")
		everythingAddr, err := sdk.AccAddressFromBech32("certik1nzgvd4k34zzf6vk3qevuh5xx8xshg6uy0l8rd5")
		if err != nil {
			panic(err)
		}
		input.CertKeeper.SetCertifier(input.Ctx, cert.Certifier{Address: certAddr})
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

		_, err = input.CertKeeper.IssueCertificate(input.Ctx, auditingCert1)
		_, err = input.CertKeeper.IssueCertificate(input.Ctx, auditingCert2)
		_, err = input.CertKeeper.IssueCertificate(input.Ctx, proofCert)
		_, err = input.CertKeeper.IssueCertificate(input.Ctx, proofCert2)
		_, err = input.CertKeeper.IssueCertificate(input.Ctx, compCert)
		if err != nil {
			panic(err)
		}

		callCheckCall, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"callCheck",
			WrapLogger(input.Ctx.Logger()),
		)
		callCheckNotCertified, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"callCheckNotCertified",
			WrapLogger(input.Ctx.Logger()),
		)
		proofCheck, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"proofCheck",
			WrapLogger(input.Ctx.Logger()),
		)
		compCheck, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"compilationCheck",
			WrapLogger(input.Ctx.Logger()),
		)
		bothCheck, _, err := abi.EncodeFunctionCall(
			testCheckAbiJsonString,
			"proofAndAuditingCheck",
			WrapLogger(input.Ctx.Logger()),
		)
		require.Nil(t, err)
		result, err = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Equal(t, []byte{0x01}, result)
		result, err = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Equal(t, []byte{0x00}, result)
		result, err = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Equal(t, []byte{0x01}, result)
		result, err = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Equal(t, []byte{0x01}, result)
		result, err = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Equal(t, []byte{0x01}, result)
		require.Nil(t, err)
	})

	t.Run("deploy and call certify validator native contract", func(t *testing.T) {
		valStr := "certikvalconspub1zcjduepq32v65eegk2yvgzdya5dqnlnc063u7mt3dh66z2xyv9rddgm6t94s4pjeat"
		code, err := hex.DecodeString(testCertifyValidatorString)
		require.Nil(t, err)

		result, err := input.CvmKeeper.Call(input.Ctx, Addrs[1], false)
		require.Nil(t, err)
		require.NotNil(t, result)
		newContractAddress := sdk.AccAddress(result)

		certAddr, err := sdk.AccAddressFromBech32("certik1nzgvd4k34zzf6vk3qevuh5xx8xshg6uy0l8rd5")
		if err != nil {
			panic(err)
		}
		input.CertKeeper.SetCertifier(input.Ctx, cert.Certifier{Address: certAddr})
		require.True(t, input.CertKeeper.IsCertifier(input.Ctx, certAddr))

		certifyValidator, _, err := abi.EncodeFunctionCall(
			testCertifyValidatorAbiJsonString,
			"certifyValidator",
			WrapLogger(input.Ctx.Logger()),
		)
		_, err = input.BankKeeper.AddCoins(input.Ctx, certAddr, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 12345)})
		require.NoError(t, err)
		certAcc := input.AccountKeeper.GetAccount(input.Ctx, certAddr)
		_ = certAcc.SetSequence(1)
		input.AccountKeeper.SetAccount(input.Ctx, certAcc)

		result, err = input.CvmKeeper.Call(input.Ctx, Addrs[0], false)
		require.Equal(t, []byte{0x00}, result)
		result, err = input.CvmKeeper.Call(input.Ctx, certAddr, false)
		fmt.Println(result)
		fmt.Println(err)
		require.Equal(t, []byte{0x01}, result)

		validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, valStr)
		require.Nil(t, err)
		require.True(t, input.CertKeeper.IsValidatorCertified(input.Ctx, validator))
	})
}
