package simulation

import (
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/v2/x/shield/keeper"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

const (
	// C's operations
	OpWeightMsgCreatePool = "op_weight_msg_create_pool"
	OpWeightMsgUpdatePool = "op_weight_msg_update_pool"

	// B and C's operations
	OpWeightMsgDepositCollateral  = "op_weight_msg_deposit_collateral"
	OpWeightMsgWithdrawCollateral = "op_weight_msg_withdraw_collateral"
	OpWeightMsgWithdrawRewards    = "op_weight_msg_withdraw_rewards"

	// P's operations
	OpWeightMsgPurchaseShield     = "op_weight_msg_purchase_shield"
	OpWeightShieldClaimProposal   = "op_weight_msg_submit_claim_proposal"
	OpWeightStakeForShield        = "op_weight_msg_stake_for_shield"
	OpWeightUnstakeFromShield     = "op_weight_msg_unstake_from_shield"
	OpWeightWithdrawReimbursement = "op_weight_msg_withdraw_reimbursement"
)

var (
	DefaultWeightMsgCreatePool            = 10
	DefaultWeightMsgUpdatePool            = 20
	DefaultWeightMsgDepositCollateral     = 20
	DefaultWeightMsgWithdrawCollateral    = 20
	DefaultWeightMsgWithdrawRewards       = 10
	DefaultWeightMsgPurchaseShield        = 20
	DefaultWeightMsgStakeForShield        = 20
	DefaultWeightMsgUnstakeFromShield     = 15
	DefaultWeightShieldClaimProposal      = 5
	DefaultWeightMsgWithdrawReimbursement = 5

	DefaultIntMax = 10000000000
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simulation.WeightedOperations {
	var weightMsgCreatePool int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = DefaultWeightMsgCreatePool
		})
	var weightMsgUpdatePool int
	appParams.GetOrGenerate(cdc, OpWeightMsgUpdatePool, &weightMsgUpdatePool, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePool = DefaultWeightMsgUpdatePool
		})
	var weightMsgDepositCollateral int
	appParams.GetOrGenerate(cdc, OpWeightMsgDepositCollateral, &weightMsgDepositCollateral, nil,
		func(_ *rand.Rand) {
			weightMsgDepositCollateral = DefaultWeightMsgDepositCollateral
		})
	var weightMsgWithdrawCollateral int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawCollateral, &weightMsgWithdrawCollateral, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawCollateral = DefaultWeightMsgWithdrawCollateral
		})
	var weightMsgWithdrawRewards int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawRewards, &weightMsgWithdrawRewards, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawRewards = DefaultWeightMsgWithdrawRewards
		})
	var weightMsgPurchaseShield int
	appParams.GetOrGenerate(cdc, OpWeightMsgPurchaseShield, &weightMsgPurchaseShield, nil,
		func(_ *rand.Rand) {
			weightMsgPurchaseShield = DefaultWeightMsgPurchaseShield
		})
	var weightMsgStakeForShield int
	appParams.GetOrGenerate(cdc, OpWeightStakeForShield, &weightMsgStakeForShield, nil,
		func(_ *rand.Rand) {
			weightMsgStakeForShield = DefaultWeightMsgStakeForShield
		})
	var weightMsgUnstakeFromShield int
	appParams.GetOrGenerate(cdc, OpWeightUnstakeFromShield, &weightMsgUnstakeFromShield, nil,
		func(_ *rand.Rand) {
			weightMsgUnstakeFromShield = DefaultWeightMsgUnstakeFromShield
		})
	var weightMsgWithdrawReimbursement int
	appParams.GetOrGenerate(cdc, OpWeightWithdrawReimbursement, &weightMsgWithdrawReimbursement, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawReimbursement = DefaultWeightMsgWithdrawReimbursement
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreatePool, SimulateMsgCreatePool(k, ak, bk, sk)),
		simulation.NewWeightedOperation(weightMsgCreatePool, SimulateMsgUpdatePool(k, ak, bk, sk)),
		simulation.NewWeightedOperation(weightMsgDepositCollateral, SimulateMsgDepositCollateral(k, ak, bk, sk)),
		simulation.NewWeightedOperation(weightMsgWithdrawCollateral, SimulateMsgWithdrawCollateral(k, ak, bk, sk)),
		simulation.NewWeightedOperation(weightMsgWithdrawRewards, SimulateMsgWithdrawRewards(k, ak)),
		simulation.NewWeightedOperation(weightMsgPurchaseShield, SimulateMsgPurchaseShield(k, ak, bk, sk)),
		simulation.NewWeightedOperation(weightMsgStakeForShield, SimulateMsgStakeForShield(k, ak, bk, sk)),
		simulation.NewWeightedOperation(weightMsgUnstakeFromShield, SimulateMsgUnstakeFromShield(k, ak, bk, sk)),
		simulation.NewWeightedOperation(weightMsgWithdrawReimbursement, SimulateMsgWithdrawReimbursement(k, ak, bk, sk)),
	}
}

// SimulateMsgCreatePool generates a MsgCreatePool object with all of its fields randomized.
func SimulateMsgCreatePool(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		pools := k.GetAllPools(ctx)
		// restrict number of pools to reduce gas consumptions for unbondings and redelegations
		if len(pools) > 20 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "too many pools"), nil, nil
		}
		// admin
		adminAddr := k.GetAdmin(ctx)
		var simAccount simtypes.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(adminAddr) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)
		bondDenom := sk.BondDenom(ctx)

		// shield
		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		totalClaimed := k.GetTotalClaimed(ctx)
		poolParams := k.GetPoolParams(ctx)
		maxShield := sdk.MinInt(totalCollateral.Sub(totalWithdrawing).Sub(totalClaimed).ToDec().Mul(poolParams.PoolShieldLimit).TruncateInt(), totalCollateral.Sub(totalWithdrawing).Sub(totalClaimed).Sub(totalShield))
		shieldAmount, err := simtypes.RandPositiveInt(r, maxShield)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, err.Error()), nil, nil
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))

		// shield limit
		// No overflow would happen when converting int64 to int in this case.
		shieldLimit := sdk.NewInt(int64(simtypes.RandIntBetween(r, int(maxShield.Int64()), int(maxShield.Int64())*5)))

		// sponsor
		sponsor := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		if _, found := k.GetPoolsBySponsor(ctx, sponsor); found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "pool not found for given sponsor"), nil, nil
		}
		// serviceFees
		nativeAmount := bk.SpendableCoins(ctx, account.GetAddress()).AmountOf(bondDenom)
		if !nativeAmount.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, ""), nil, nil
		}
		nativeAmount, err = simtypes.RandPositiveInt(r, nativeAmount)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, err.Error()), nil, nil
		}
		nativeServiceFees := sdk.NewCoins(sdk.NewCoin(bondDenom, nativeAmount))
		foreignAmount, err := simtypes.RandPositiveInt(r, sdk.NewInt(int64(DefaultIntMax)))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, err.Error()), nil, nil
		}
		foreignServiceFees := sdk.NewCoins(sdk.NewCoin(sponsor, foreignAmount))

		serviceFees := types.MixedCoins{Native: nativeServiceFees, Foreign: foreignServiceFees}
		sponsorAcc, _ := simtypes.RandomAcc(r, accs)
		description := simtypes.RandStringOfLength(r, 42)

		msg := types.NewMsgCreatePool(simAccount.Address, shield, serviceFees, sponsor, sponsorAcc.Address, description, shieldLimit)

		fees := sdk.Coins{}
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgUpdatePool generates a MsgUpdatePool object with all of its fields randomized.
func SimulateMsgUpdatePool(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		adminAddr := k.GetAdmin(ctx)
		var simAccount simtypes.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(adminAddr) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)
		bondDenom := sk.BondDenom(ctx)

		// pool
		poolID, _, found := keeper.RandomPoolInfo(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePool, "random pool info not found"), nil, nil
		}
		pool, found := k.GetPool(ctx, poolID)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePool, "pool not found"), nil, nil
		}

		// shield
		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		totalClaimed := k.GetTotalClaimed(ctx)
		poolParams := k.GetPoolParams(ctx)
		maxShield := computeMaxShield(pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield, poolParams)
		shieldAmount, err := simtypes.RandPositiveInt(r, maxShield)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePool, err.Error()), nil, nil
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))

		// serviceFees
		nativeAmount := bk.SpendableCoins(ctx, account.GetAddress()).AmountOf(bondDenom)
		if !nativeAmount.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePool, ""), nil, nil
		}
		nativeAmount, err = simtypes.RandPositiveInt(r, nativeAmount)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePool, err.Error()), nil, nil
		}
		nativeServiceFees := sdk.NewCoins(sdk.NewCoin(bondDenom, nativeAmount))
		foreignAmount, err := simtypes.RandPositiveInt(r, sdk.NewInt(int64(DefaultIntMax)))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePool, err.Error()), nil, nil
		}
		foreignServiceFees := sdk.NewCoins(sdk.NewCoin(pool.Sponsor, foreignAmount))

		serviceFees := types.MixedCoins{Native: nativeServiceFees, Foreign: foreignServiceFees}
		description := simtypes.RandStringOfLength(r, 42)

		msg := types.NewMsgUpdatePool(simAccount.Address, shield, serviceFees, poolID, description, sdk.ZeroInt())

		fees := sdk.Coins{}
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePool, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgDepositCollateral generates a MsgDepositCollateral object with all of its fields randomized.
func SimulateMsgDepositCollateral(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		delAddr, available, found := keeper.RandomDelegation(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositCollateral, "random delegation not found"), nil, nil
		}
		var simAccount simtypes.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)

		// collateral coins
		provider, found := k.GetProvider(ctx, simAccount.Address)
		if found {
			available = provider.DelegationBonded.Sub(provider.Collateral)
		}
		collateralAmount, err := simtypes.RandPositiveInt(r, available)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositCollateral, err.Error()), nil, nil
		}
		collateral := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), collateralAmount))

		msg := types.NewMsgDepositCollateral(simAccount.Address, collateral)

		fees := sdk.Coins{}
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositCollateral, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgWithdrawCollateral generates a MsgWithdrawCollateral object with all of its fields randomized.
func SimulateMsgWithdrawCollateral(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		provider, found := keeper.RandomProvider(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawCollateral, "random provider not found"), nil, nil
		}

		var simAccount simtypes.Account
		for _, simAcc := range accs {
			providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
			if err != nil {
				panic(err)
			}
			if simAcc.Address.Equals(providerAddr) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)

		// withdraw coins
		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		withdrawable := sdk.MinInt(provider.Collateral.Sub(provider.Withdrawing), totalCollateral.Sub(totalWithdrawing).Sub(totalShield))
		withdrawAmount, err := simtypes.RandPositiveInt(r, withdrawable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawCollateral, err.Error()), nil, nil
		}
		withdraw := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), withdrawAmount))

		msg := types.NewMsgWithdrawCollateral(simAccount.Address, withdraw)

		fees := sdk.Coins{}
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawCollateral, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgWithdrawRewards generates a MsgWithdrawRewards object with all of its fields randomized.
func SimulateMsgWithdrawRewards(k keeper.Keeper, ak types.AccountKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		provider, found := keeper.RandomProvider(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawCollateral, "random provider not found"), nil, nil
		}
		var simAccount simtypes.Account
		for _, simAcc := range accs {
			providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
			if err != nil {
				panic(err)
			}
			if simAcc.Address.Equals(providerAddr) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)

		msg := types.NewMsgWithdrawRewards(simAccount.Address)

		fees := sdk.Coins{}
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawCollateral, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgPurchaseShield generates a MsgPurchaseShield object with all of its fields randomized.
func SimulateMsgPurchaseShield(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		purchaser, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, purchaser.Address)
		bondDenom := sk.BondDenom(ctx)

		poolID, _, found := keeper.RandomPoolInfo(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurchaseShield, "random pool info not found"), nil, nil
		}
		pool, found := k.GetPool(ctx, poolID)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurchaseShield, "pool not found"), nil, nil
		}

		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		totalClaimed := k.GetTotalClaimed(ctx)
		poolParams := k.GetPoolParams(ctx)
		maxShield := computeMaxShield(pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield, poolParams)
		shieldAmount, err := simtypes.RandPositiveInt(r, maxShield)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurchaseShield, err.Error()), nil, nil
		}
		if shieldAmount.ToDec().Mul(poolParams.ShieldFeesRate).GT(bk.SpendableCoins(ctx, account.GetAddress()).AmountOf(bondDenom).ToDec()) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurchaseShield, ""), nil, nil
		}
		if shieldAmount.ToDec().Mul(poolParams.ShieldFeesRate).TruncateInt().IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurchaseShield, ""), nil, nil
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))

		description := simtypes.RandStringOfLength(r, 100)
		msg := types.NewMsgPurchaseShield(poolID, shield, description, purchaser.Address)

		fees := sdk.Coins{}
		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			purchaser.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurchaseShield, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper, sk types.StakingKeeper) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightShieldClaimProposal,
			DefaultWeightShieldClaimProposal,
			SimulateShieldClaimProposalContent(k, sk),
		),
	}
}

// SimulateShieldClaimProposalContent generates random shield claim proposal content
func SimulateShieldClaimProposalContent(k keeper.Keeper, sk types.StakingKeeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		bondDenom := sk.BondDenom(ctx)
		purchaseList, found := keeper.RandomPurchaseList(r, k, ctx)
		if len(purchaseList.Entries) == 0 {
			return nil
		}
		i := r.Intn(len(purchaseList.Entries))
		poolID := purchaseList.PoolId
		purchaser := purchaseList.Purchaser
		purchaserAddr, err := sdk.AccAddressFromBech32(purchaser)
		if err != nil {
			panic(err)
		}
		purchase := purchaseList.Entries[i]
		if !found || purchase.ProtectionEndTime.Before(ctx.BlockTime()) {
			return nil
		}
		lossAmount, err := simtypes.RandPositiveInt(r, purchase.Shield)
		if err != nil {
			return nil
		}
		return types.NewShieldClaimProposal(
			poolID,
			sdk.NewCoins(sdk.NewCoin(bondDenom, lossAmount)),
			purchase.PurchaseId,
			simtypes.RandStringOfLength(r, 500),
			simtypes.RandStringOfLength(r, 500),
			purchaserAddr,
		)
	}
}

// SimulateMsgStakeForShield generates a MsgPurchaseShield object with all of its fields randomized.
func SimulateMsgStakeForShield(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		purchaser, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, purchaser.Address)
		bondDenom := sk.BondDenom(ctx)

		poolID, _, found := keeper.RandomPoolInfo(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStakeForShield, "random pool info not found"), nil, nil
		}
		pool, found := k.GetPool(ctx, poolID)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStakeForShield, "pool not found"), nil, nil
		}

		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		totalClaimed := k.GetTotalClaimed(ctx)
		poolParams := k.GetPoolParams(ctx)
		maxShield := computeMaxShield(pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield, poolParams)
		accountMax := sdk.OneDec().Quo(k.GetShieldStakingRate(ctx)).MulInt(bk.GetAllBalances(ctx, account.GetAddress()).AmountOf(k.BondDenom(ctx))).TruncateInt()
		max := sdk.MinInt(accountMax, maxShield)
		shieldAmount, err := simtypes.RandPositiveInt(r, max)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStakeForShield, err.Error()), nil, nil
		}
		rate := k.GetShieldStakingRate(ctx)
		maxShieldAmt := bk.SpendableCoins(ctx, account.GetAddress()).AmountOf(bondDenom).ToDec().Quo(rate).TruncateInt()
		if shieldAmount.GT(maxShieldAmt) {
			shieldAmount = maxShieldAmt
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))
		if shield.IsZero() || k.GetShieldStakingRate(ctx).MulInt(shield.AmountOf(bondDenom)).TruncateInt().IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStakeForShield, ""), nil, nil
		}

		description := simtypes.RandStringOfLength(r, 100)
		msg := types.NewMsgStakeForShield(poolID, shield, description, purchaser.Address)

		fees := sdk.Coins{}
		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			purchaser.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStakeForShield, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgUnstakeFromShield generates a MsgUnstakeFromShield object with all of its fields randomized.
func SimulateMsgUnstakeFromShield(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		bondDenom := sk.BondDenom(ctx)
		stakeForShields := k.GetAllStakeForShields(ctx)
		if len(stakeForShields) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstakeFromShield, "no stake for shields found"), nil, nil
		}
		index := simtypes.RandIntBetween(r, 0, len(stakeForShields))
		sfs := stakeForShields[index]

		withdrawable := sfs.Amount.Sub(sfs.WithdrawRequested)
		withdrawableCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, withdrawable))
		shield := simtypes.RandSubsetCoins(r, withdrawableCoins)
		purchaserAddr, err := sdk.AccAddressFromBech32(sfs.Purchaser)
		if err != nil {
			panic(err)
		}
		msg := types.NewMsgUnstakeFromShield(sfs.PoolId, shield, purchaserAddr)

		var account authtypes.AccountI
		var simAcc simtypes.Account
		for _, acc := range accs {
			if acc.Address.Equals(purchaserAddr) {
				account = ak.GetAccount(ctx, acc.Address)
				simAcc = acc
			}
		}
		if account == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstakeFromShield, "account is nil"), nil, nil
		}

		fees := sdk.Coins{}
		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstakeFromShield, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgWithdrawReimbursement generates a MsgWithdrawReimbursement object with randomized fields.
func SimulateMsgWithdrawReimbursement(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		prPair, found := keeper.RandomMaturedProposalIDReimbursementPair(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawReimbursement, "no mature proposal id - reimbursement pair found"), nil, nil
		}

		var simAccount simtypes.Account
		for _, simAcc := range accs {
			beneficiaryAddr, err := sdk.AccAddressFromBech32(prPair.Reimbursement.Beneficiary)
			if err != nil {
				panic(err)
			}
			if simAcc.Address.Equals(beneficiaryAddr) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)

		msg := types.NewMsgWithdrawReimbursement(prPair.ProposalId, simAccount.Address)

		fees := sdk.Coins{}
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawReimbursement, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func computeMaxShield(pool types.Pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield sdk.Int, poolParams types.PoolParams) sdk.Int {
	poolLimit := pool.ShieldLimit.Sub(pool.Shield)
	globalLimit := sdk.MinInt(totalCollateral.Sub(totalWithdrawing).Sub(totalClaimed).ToDec().Mul(poolParams.PoolShieldLimit).TruncateInt().Sub(pool.Shield),
		totalCollateral.Sub(totalWithdrawing).Sub(totalClaimed).Sub(totalShield))
	return sdk.MinInt(poolLimit, globalLimit)
}
