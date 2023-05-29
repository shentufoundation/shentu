package simulation

import (
	"bytes"
	cryptorand "crypto/rand"
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
	OpWeightMsgCreateProgram = "op_weight_msg_create_program"
	OpWeightMsgSubmitFinding = "op_weight_msg_submit_finding"
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simulation.WeightedOperations {
	var weightMsgCreateProgram int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateProgram, &weightMsgCreateProgram, nil,
		func(_ *rand.Rand) {
			weightMsgCreateProgram = simappparams.DefaultWeightMsgSend
		},
	)

	var weightMsgSubmitFinding int
	appParams.GetOrGenerate(cdc, OpWeightMsgSubmitFinding, &weightMsgSubmitFinding, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitFinding = simappparams.DefaultWeightMsgSend
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

		hostAcc := ak.GetAccount(ctx, host.Address)
		deposit := simtypes.RandSubsetCoins(r, bk.SpendableCoins(ctx, hostAcc.GetAddress()))
		if deposit.Empty() {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: empty deposit, skip this tx", "", false, nil), nil, nil
		}

		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, hostAcc.GetAddress()).Sub(deposit))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateProgram, err.Error()), nil, err
		}

		desc := fmt.Sprintf("simulation desc %d", r.Int())

		priKey, _ := ecies.GenerateKey(cryptorand.Reader, ecies.DefaultCurve, nil)
		pubKey := crypto.FromECDSAPub(&priKey.ExportECDSA().PublicKey)

		commission := sdk.NewDec(r.Int63n(10))

		endTime := time.Now().Add(time.Duration((r.Intn(120) + 10) * int(time.Hour)))

		msg, _ := types.NewMsgCreateProgram(host.Address.String(), desc, pubKey, commission, deposit, endTime, endTime, endTime)

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
			// {
			// 	BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 0, 20),
			// 	Op:          SimulateMsgSubmitFinding(k, ak, bk, priKey.PublicKey),
			// },
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), futureOperations, nil
	}
}

func SimulateMsgSubmitFinding(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, pubKey ecies.PublicKey) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		submitter, _ := simtypes.RandomAcc(r, accs)
		submitterAcc := ak.GetAccount(ctx, submitter.Address)

		title := fmt.Sprintf("finding title %d", r.Int())
		desc := fmt.Sprintf("finding desc %d", r.Int())
		poc := fmt.Sprintf("finding poc %d", r.Int())

		randBytes := make([]byte, 64)

		cryptorand.Read(randBytes)
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
		cryptorand.Read(randBytes)
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

		rsp, _ := k.Programs(sdk.WrapSDKContext(ctx), &types.QueryProgramsRequest{})
		programID := rsp.Programs[r.Intn(len(rsp.Programs))].ProgramId

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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}
