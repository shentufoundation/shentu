package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingSim "github.com/cosmos/cosmos-sdk/x/staking/simulation"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	simutil "github.com/shentufoundation/shentu/v2/x/auth/simulation"
)

const (
	OpWeightMsgCreateValidator = "op_weight_msg_create_validator"
	OpWeightMsgEditValidator   = "op_weight_msg_edit_validator"
	OpWeightMsgDelegate        = "op_weight_msg_delegate"
	OpWeightMsgUndelegate      = "op_weight_msg_undelegate"
	OpWeightMsgBeginRedelegate = "op_weight_msg_begin_redelegate"
)

func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper,
	k *stakingkeeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgCreateValidator int
		weightMsgEditValidator   int
		weightMsgDelegate        int
		weightMsgUndelegate      int
		weightMsgBeginRedelegate int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateValidator, &weightMsgCreateValidator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateValidator = stakingSim.DefaultWeightMsgCreateValidator
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgEditValidator, &weightMsgEditValidator, nil,
		func(_ *rand.Rand) {
			weightMsgEditValidator = stakingSim.DefaultWeightMsgEditValidator
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDelegate, &weightMsgDelegate, nil,
		func(_ *rand.Rand) {
			weightMsgDelegate = stakingSim.DefaultWeightMsgDelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgUndelegate, &weightMsgUndelegate, nil,
		func(_ *rand.Rand) {
			weightMsgUndelegate = stakingSim.DefaultWeightMsgUndelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgBeginRedelegate, &weightMsgBeginRedelegate, nil,
		func(_ *rand.Rand) {
			weightMsgBeginRedelegate = stakingSim.DefaultWeightMsgBeginRedelegate
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateValidator,
			SimulateMsgCreateValidator(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgEditValidator,
			stakingSim.SimulateMsgEditValidator(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDelegate,
			stakingSim.SimulateMsgDelegate(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUndelegate,
			SimulateMsgUndelegate(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBeginRedelegate,
			SimulateMsgBeginRedelegate(ak, bk, k),
		),
	}
}

func SimulateMsgCreateValidator(k *stakingkeeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		address := sdk.ValAddress(simAccount.Address)

		_, found := k.GetValidator(ctx, address)
		if found {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "unable to find a validator"), nil, nil
		}

		denom := k.GetParams(ctx).BondDenom

		balance := bk.GetBalance(ctx, simAccount.Address, denom).Amount
		if !balance.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "balance is negative"), nil, nil
		}

		amount, err := simtypes.RandPositiveInt(r, balance)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "unable to generate positive amount"), nil, err
		}

		selfDelegation := sdk.NewCoin(denom, amount)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		var fees sdk.Coins

		coins, hasNeg := spendable.SafeSub(sdk.Coins{selfDelegation}...)
		if !hasNeg {
			fees, err = simutil.RandomReasonableFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "unable to generate fees"), nil, err
			}
		}

		description := stakingtypes.NewDescription(
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
		)

		maxCommission := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 100)), 2)
		commission := stakingtypes.NewCommissionRates(
			simtypes.RandomDecAmount(r, maxCommission),
			maxCommission,
			simtypes.RandomDecAmount(r, maxCommission),
		)

		msg, err := stakingtypes.NewMsgCreateValidator(address, simAccount.ConsKey.PubKey(), selfDelegation, description, commission)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, msg.Type(), "unable to create CreateValidator message"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    stakingtypes.ModuleName,
		}

		return simulation.GenAndDeliverTx(txCtx, fees)
	}
}

func SimulateMsgUndelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k *stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		val, ok := testutil.RandSliceElem(r, k.GetAllValidators(ctx))
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "validator is not ok"), nil, nil
		}

		valAddr := val.GetOperator()
		delegations := k.GetValidatorDelegations(ctx, val.GetOperator())
		if delegations == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "keeper does have any delegation entries"), nil, nil
		}

		// get random delegator from validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, nil
		}

		totalBond := val.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "total bond is negative"), nil, nil
		}

		unbondAmt, err := simtypes.RandPositiveInt(r, totalBond)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, err
		}

		if unbondAmt.IsZero() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, nil
		}

		msg := stakingtypes.NewMsgUndelegate(
			delAddr, valAddr, sdk.NewCoin(k.BondDenom(ctx), unbondAmt),
		)

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simtypes.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}
		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		account := ak.GetAccount(ctx, delAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      stakingtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgBeginRedelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k *stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// get random source validator
		allVals := k.GetAllValidators(ctx)
		srcVal, ok := testutil.RandSliceElem(r, allVals)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "unable to pick validator"), nil, nil
		}

		srcAddr := srcVal.GetOperator()
		delegations := k.GetValidatorDelegations(ctx, srcAddr)
		if delegations == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "keeper does have any delegation entries"), nil, nil
		}

		// get random delegator from src validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasReceivingRedelegation(ctx, delAddr, srcAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil // skip
		}

		// get random destination validator
		destVal, ok := testutil.RandSliceElem(r, allVals)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "unable to pick validator"), nil, nil
		}

		destAddr := destVal.GetOperator()
		if srcAddr.Equals(destAddr) || destVal.InvalidExRate() || k.HasMaxRedelegationEntries(ctx, delAddr, srcAddr, destAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "checks failed"), nil, nil
		}

		totalBond := srcVal.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "total bond is negative"), nil, nil
		}

		redAmt, err := simtypes.RandPositiveInt(r, totalBond)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, err
		}

		if redAmt.IsZero() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil
		}

		// check if the shares truncate to zero
		shares, err := srcVal.SharesFromTokens(redAmt)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, err
		}

		if srcVal.TokensFromShares(shares).TruncateInt().IsZero() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil // skip
		}

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simtypes.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}
		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		// get tx fees
		account := ak.GetAccount(ctx, delAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := stakingtypes.NewMsgBeginRedelegate(
			delAddr, srcAddr, destAddr,
			sdk.NewCoin(k.BondDenom(ctx), redAmt),
		)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      stakingtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
