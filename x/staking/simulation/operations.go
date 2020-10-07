package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingSim "github.com/cosmos/cosmos-sdk/x/staking/simulation"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/x/staking/internal/types"
)

const (
	OpWeightMsgCreateValidator = "op_weight_msg_create_validator"
	OpWeightMsgEditValidator   = "op_weight_msg_edit_validator"
	OpWeightMsgDelegate        = "op_weight_msg_delegate"
	OpWeightMsgUndelegate      = "op_weight_msg_undelegate"
	OpWeightMsgBeginRedelegate = "op_weight_msg_begin_redelegate"
)

func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, ak stakingTypes.AccountKeeper, ck types.CertKeeper,
	k staking.Keeper) simulation.WeightedOperations {
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
			SimulateMsgCreateValidator(k, ak, ck),
		),
		simulation.NewWeightedOperation(
			weightMsgEditValidator,
			stakingSim.SimulateMsgEditValidator(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDelegate,
			stakingSim.SimulateMsgDelegate(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUndelegate,
			SimulateMsgUndelegate(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBeginRedelegate,
			SimulateMsgBeginRedelegate(ak, k),
		),
	}
}

func SimulateMsgCreateValidator(k staking.Keeper, ak stakingTypes.AccountKeeper, ck types.CertKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		simAccount, _ := simulation.RandomAcc(r, accs)
		address := sdk.ValAddress(simAccount.Address)

		_, found := k.GetValidator(ctx, address)
		if found {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		if !ck.IsValidatorCertified(ctx, simAccount.PubKey) {
			return simulation.NewOperationMsgBasic(stakingTypes.ModuleName, "NoOp: not a certified validator", "", false, nil), nil, nil
		}

		denom := k.GetParams(ctx).BondDenom
		amount := ak.GetAccount(ctx, simAccount.Address).GetCoins().AmountOf(denom)
		if !amount.IsPositive() {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		amount, err := simulation.RandPositiveInt(r, amount)
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		selfDelegation := sdk.NewCoin(denom, amount)

		account := ak.GetAccount(ctx, simAccount.Address)
		coins := account.SpendableCoins(ctx.BlockTime())

		var fees sdk.Coins
		coins, hasNeg := coins.SafeSub(sdk.Coins{selfDelegation})
		if !hasNeg {
			fees, err = simulation.RandomFees(r, ctx, coins)
			if err != nil {
				return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
			}
		}

		description := stakingTypes.NewDescription(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 10),
		)

		maxCommission := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 0, 100)), 2)
		commission := stakingTypes.NewCommissionRates(
			simulation.RandomDecAmount(r, maxCommission),
			maxCommission,
			simulation.RandomDecAmount(r, maxCommission),
		)

		msg := stakingTypes.NewMsgCreateValidator(address, simAccount.PubKey,
			selfDelegation, description, commission, sdk.OneInt())

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgUndelegate(ak stakingTypes.AccountKeeper, k staking.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		// get random validator
		validator, ok := stakingKeeper.RandomValidator(r, k, ctx)
		if !ok {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}
		valAddr := validator.GetOperator()

		delegations := k.GetValidatorDelegations(ctx, validator.OperatorAddress)

		// get random delegator from validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		totalBond := validator.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		unbondAmt, err := simulation.RandPositiveInt(r, totalBond)
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		if unbondAmt.IsZero() {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		msg := stakingTypes.NewMsgUndelegate(
			delAddr, valAddr, sdk.NewCoin(k.BondDenom(ctx), unbondAmt),
		)

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simulation.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}
		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		account := ak.GetAccount(ctx, delAddr)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas*10,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgBeginRedelegate(ak stakingTypes.AccountKeeper, k staking.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		// get random source validator
		srcVal, ok := stakingKeeper.RandomValidator(r, k, ctx)
		if !ok {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}
		srcAddr := srcVal.GetOperator()

		delegations := k.GetValidatorDelegations(ctx, srcAddr)

		// get random delegator from src validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasReceivingRedelegation(ctx, delAddr, srcAddr) {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil // skip
		}

		// get random destination validator
		destVal, ok := stakingKeeper.RandomValidator(r, k, ctx)
		if !ok {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}
		destAddr := destVal.GetOperator()

		if srcAddr.Equals(destAddr) ||
			destVal.InvalidExRate() ||
			k.HasMaxRedelegationEntries(ctx, delAddr, srcAddr, destAddr) {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		totalBond := srcVal.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		redAmt, err := simulation.RandPositiveInt(r, totalBond)
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		if redAmt.IsZero() {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil
		}

		// check if the shares truncate to zero
		shares, err := srcVal.SharesFromTokens(redAmt)
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		if srcVal.TokensFromShares(shares).TruncateInt().IsZero() {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, nil // skip
		}

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simulation.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}
		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		// get tx fees
		account := ak.GetAccount(ctx, delAddr)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		msg := stakingTypes.NewMsgBeginRedelegate(
			delAddr, srcAddr, destAddr,
			sdk.NewCoin(k.BondDenom(ctx), redAmt),
		)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas*10,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(stakingTypes.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
