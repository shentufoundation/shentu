package simulation

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
	simutil "github.com/shentufoundation/shentu/v2/x/cvm/simulation"
)

const (
	// OpWeightMsgCreateProgram desc
	OpWeightMsgCreateProgram = "op_weight_msg_create_program"
)

// TestRandReader is for generate pri/pub key
type TestRandReader struct {
	roll rand.Rand
}

func (reader TestRandReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simulation.WeightedOperations {
	var weightMsgCreateProgram int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateProgram, &weightMsgCreateProgram, nil,
		func(_ *rand.Rand) {
			weightMsgCreateProgram = simappparams.DefaultWeightMsgSend
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreateProgram, SimulateMsgCreateProgram(k, ak, bk)),
	}
}

// SimulateMsgCreateProgram generates a MsgCreateProgram object with all of its fields randomized.
// This operation leads a series of future operations.
func SimulateMsgCreateProgram(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		host, _ := simtypes.RandomAcc(r, accs)

		maxDepositAmount := sdk.NewInt(100)

		hostAcc := ak.GetAccount(ctx, host.Address)
		deposit := simtypes.RandSubsetCoins(r, bk.SpendableCoins(ctx, hostAcc.GetAddress()))
		if deposit.Empty() {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: empty deposit, skip this tx", "", false, nil), nil, nil
		}
		if deposit.AmountOf(sdk.DefaultBondDenom).GTE(maxDepositAmount) {
			deposit = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, maxDepositAmount)}
		}

		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, hostAcc.GetAddress()).Sub(deposit))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateProgram, err.Error()), nil, err
		}

		desc := fmt.Sprintf("simulation desc %d", r.Int())

		reader := TestRandReader{roll: *r}
		priKey, _ := ecies.GenerateKey(reader, ecies.DefaultCurve, nil)
		pubKey := crypto.FromECDSAPub(priKey.PublicKey.ExportECDSA())
		commission := sdk.OneDec()

		endTime := ctx.BlockTime().Add(time.Duration((r.Intn(120) + 720) * int(time.Hour)))

		msg, _ := types.NewMsgCreateProgram(host.Address.String(), desc, pubKey, commission, deposit, endTime)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{hostAcc.GetAccountNumber()},
			[]uint64{hostAcc.GetSequence()},
			host.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		programID, err := k.GetNextProgramID(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		futureOperations := []simtypes.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 0, 9),
				Op:          SimulateMsgSubmitFinding(k, ak, bk, priKey, programID, host),
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), futureOperations, nil
	}
}

// SimulateMsgSubmitFinding submits a Finding to the Program
func SimulateMsgSubmitFinding(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, priKey *ecies.PrivateKey, programID uint64, host simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		submitter, _ := simtypes.RandomAcc(r, accs)
		submitterAcc := ak.GetAccount(ctx, submitter.Address)

		title := fmt.Sprintf("finding title %d", r.Int())
		desc := fmt.Sprintf("finding desc %d", r.Int())
		poc := fmt.Sprintf("finding poc %d", r.Int())

		randBytes := make([]byte, 64)

		pubKey := priKey.PublicKey
		r.Read(randBytes)
		reader := bytes.NewReader(randBytes)
		descEnc, err := ecies.Encrypt(reader, &pubKey, []byte(desc), nil, nil)
		if err != nil {
			fmt.Printf("Error on descEnc: %#v\n", err)
		}
		descEnc = append(descEnc, randBytes...)
		descAny, err := codectypes.NewAnyWithValue(&types.EciesEncryptedDesc{
			FindingDesc: descEnc,
		})
		if err != nil {
			fmt.Printf("Error on descAny: %#v\n", err)
		}
		r.Read(randBytes)
		reader = bytes.NewReader(randBytes)
		pocEnc, err := ecies.Encrypt(reader, &pubKey, []byte(poc), nil, nil)
		if err != nil {
			fmt.Printf("Error on pocEnc: %#v\n", err)
		}
		pocEnc = append(pocEnc, randBytes...)
		pocAny, err := codectypes.NewAnyWithValue(&types.EciesEncryptedPoc{
			FindingPoc: pocEnc,
		})
		if err != nil {
			fmt.Printf("Error on pocAny: %#v\n", err)
		}

		serverity := r.Int31n(5)

		msg := types.NewMsgSubmitFinding(submitter.Address.String(), title, descAny, pocAny, programID, serverity)

		fees, _ := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, submitterAcc.GetAddress()))

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{submitterAcc.GetAccountNumber()},
			[]uint64{submitterAcc.GetSequence()},
			submitter.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		futureOperations := []simtypes.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 10, 19),
				Op:          SimulateMsgAcceptFinding(k, ak, bk, priKey, programID, host),
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), futureOperations, nil
	}
}

// SimulateMsgAcceptFinding accept finding by host
func SimulateMsgAcceptFinding(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, priKey *ecies.PrivateKey, programID uint64, host simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		hostAcc := ak.GetAccount(ctx, host.Address)
		findingIDs, _ := k.GetPidFindingIDList(ctx, programID)
		findingID, ok := randSliceElem(r, findingIDs)
		if !ok {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: empty finding, skip this tx", "", false, nil), nil, nil
		}
		comment := fmt.Sprintf("finding comment %d", r.Int())
		randBytes := make([]byte, 64)
		r.Read(randBytes)
		reader := bytes.NewReader(randBytes)
		commentEnc, err := ecies.Encrypt(reader, &priKey.PublicKey, []byte(comment), nil, nil)
		if err != nil {
			fmt.Printf("Error on commentEnc: %#v\n", err)
		}
		commentEnc = append(commentEnc, randBytes...)
		commentAny, _ := codectypes.NewAnyWithValue(&types.EciesEncryptedComment{
			FindingComment: commentEnc,
		})
		msg := types.NewMsgHostAcceptFinding(findingID, commentAny, host.Address)
		fees, _ := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, host.Address))
		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{hostAcc.GetAccountNumber()},
			[]uint64{hostAcc.GetSequence()},
			host.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		futureOperations := []simtypes.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 30, 39),
				Op:          SimulateMsgEndProgram(k, ak, bk, programID, host),
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), futureOperations, nil
	}
}

// SimulateMsgRejectFinding reject finding by host
func SimulateMsgRejectFinding(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, priKey *ecies.PrivateKey, programID uint64, host simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		hostAcc := ak.GetAccount(ctx, host.Address)
		findingIDs, _ := k.GetPidFindingIDList(ctx, programID)
		findingID, ok := randSliceElem(r, findingIDs)
		if !ok {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: empty finding, skip this tx", "", false, nil), nil, nil
		}
		finding, _ := k.GetFinding(ctx, findingID)
		if finding.FindingId == 0 {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: zero finding, skip this tx", "", false, nil), nil, nil
		}
		comment := fmt.Sprintf("finding comment %d", r.Int())
		randBytes := make([]byte, 64)
		r.Read(randBytes)
		reader := bytes.NewReader(randBytes)
		commentEnc, err := ecies.Encrypt(reader, &priKey.PublicKey, []byte(comment), nil, nil)
		if err != nil {
			fmt.Printf("Error on commentEnc: %#v\n", err)
		}
		commentEnc = append(commentEnc, randBytes...)
		commentAny, _ := codectypes.NewAnyWithValue(&types.EciesEncryptedComment{
			FindingComment: commentEnc,
		})
		msg := types.NewMsgHostRejectFinding(findingID, commentAny, host.Address)
		fees, _ := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, host.Address))
		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{hostAcc.GetAccountNumber()},
			[]uint64{hostAcc.GetSequence()},
			host.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgEndProgram generates a MsgEndProgram object with all of its fields randomized.
func SimulateMsgEndProgram(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, programID uint64, host simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		hostAcc := ak.GetAccount(ctx, host.Address)
		msg := types.NewMsgEndProgram(host.Address.String(), programID)
		fees, _ := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, host.Address))
		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{hostAcc.GetAccountNumber()},
			[]uint64{hostAcc.GetSequence()},
			host.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func randSliceElem(r *rand.Rand, elems []uint64) (uint64, bool) {
	if len(elems) == 0 {
		var e uint64
		return e, false
	}
	return elems[r.Intn(len(elems))], true
}
