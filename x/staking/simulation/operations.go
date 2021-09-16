package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingSim "github.com/cosmos/cosmos-sdk/x/staking/simulation"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	OpWeightMsgCreateValidator = "op_weight_msg_create_validator"
	OpWeightMsgEditValidator   = "op_weight_msg_edit_validator"
	OpWeightMsgDelegate        = "op_weight_msg_delegate"
	OpWeightMsgUndelegate      = "op_weight_msg_undelegate"
	OpWeightMsgBeginRedelegate = "op_weight_msg_begin_redelegate"
)

func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper,
	k stakingkeeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgCreateValidator int
		weightMsgEditValidator   int
		weightMsgDelegate        int
		weightMsgUndelegate      int
		weightMsgBeginRedelegate int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateValidator, &weightMsgCreateValidator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateValidator = simappparams.DefaultWeightMsgCreateValidator
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgEditValidator, &weightMsgEditValidator, nil,
		func(_ *rand.Rand) {
			weightMsgEditValidator = simappparams.DefaultWeightMsgEditValidator
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDelegate, &weightMsgDelegate, nil,
		func(_ *rand.Rand) {
			weightMsgDelegate = simappparams.DefaultWeightMsgDelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgUndelegate, &weightMsgUndelegate, nil,
		func(_ *rand.Rand) {
			weightMsgUndelegate = simappparams.DefaultWeightMsgUndelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgBeginRedelegate, &weightMsgBeginRedelegate, nil,
		func(_ *rand.Rand) {
			weightMsgBeginRedelegate = simappparams.DefaultWeightMsgBeginRedelegate
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

func SimulateMsgCreateValidator(k stakingkeeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		address := sdk.ValAddress(simAccount.Address)

		_, found := k.GetValidator(ctx, address)
		if found {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "unable to find a validator"), nil, nil
		}

		denom := k.GetParams(ctx).BondDenom
		amount := bk.GetBalance(ctx, simAccount.Address, denom).Amount
		if !amount.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, ""), nil, nil
		}

		amount, err := simtypes.RandPositiveInt(r, amount)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, ""), nil, err
		}

		selfDelegation := sdk.NewCoin(denom, amount)

		account := ak.GetAccount(ctx, simAccount.Address)
		coins := bk.SpendableCoins(ctx, account.GetAddress())

		var fees sdk.Coins
		coins, hasNeg := coins.SafeSub(sdk.Coins{selfDelegation})
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, ""), nil, err
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

		msg, err := stakingtypes.NewMsgCreateValidator(address, simAccount.ConsKey.PubKey(), selfDelegation, description, commission, sdk.OneInt())
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, msg.Type(), "unable to create CreateValidator message"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func SimulateMsgUndelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// get random validator
		validator, ok := stakingkeeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, nil
		}
		valAddr := validator.GetOperator()

		valOperAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, err.Error()), nil, nil
		}
		delegations := k.GetValidatorDelegations(ctx, valOperAddr)

		// get random delegator from validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, nil
		}

		totalBond := validator.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, nil
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
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas*10,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func SimulateMsgBeginRedelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// get random source validator
		srcVal, ok := stakingkeeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil
		}
		srcAddr := srcVal.GetOperator()

		delegations := k.GetValidatorDelegations(ctx, srcAddr)

		// get random delegator from src validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasReceivingRedelegation(ctx, delAddr, srcAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil // skip
		}

		// get random destination validator
		destVal, ok := stakingkeeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil
		}
		destAddr := destVal.GetOperator()

		if srcAddr.Equals(destAddr) ||
			destVal.InvalidExRate() ||
			k.HasMaxRedelegationEntries(ctx, delAddr, srcAddr, destAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil
		}

		totalBond := srcVal.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, nil
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
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, err
		}

		msg := stakingtypes.NewMsgBeginRedelegate(
			delAddr, srcAddr, destAddr,
			sdk.NewCoin(k.BondDenom(ctx), redAmt),
		)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas*10,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}
