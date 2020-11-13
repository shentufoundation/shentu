package simulation

import (
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
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

	DefaultIntMax = 100000000000
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.WeightedOperations {
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
		simulation.NewWeightedOperation(weightMsgCreatePool, SimulateMsgCreatePool(k, ak, sk)),
		simulation.NewWeightedOperation(weightMsgCreatePool, SimulateMsgUpdatePool(k, ak, sk)),
		simulation.NewWeightedOperation(weightMsgDepositCollateral, SimulateMsgDepositCollateral(k, ak, sk)),
		simulation.NewWeightedOperation(weightMsgWithdrawCollateral, SimulateMsgWithdrawCollateral(k, ak, sk)),
		simulation.NewWeightedOperation(weightMsgWithdrawRewards, SimulateMsgWithdrawRewards(k, ak)),
		simulation.NewWeightedOperation(weightMsgPurchaseShield, SimulateMsgPurchaseShield(k, ak, sk)),
		simulation.NewWeightedOperation(weightMsgStakeForShield, SimulateMsgStakeForShield(k, ak, sk)),
		simulation.NewWeightedOperation(weightMsgUnstakeFromShield, SimulateMsgUnstakeFromShield(k, ak, sk)),
		simulation.NewWeightedOperation(weightMsgWithdrawReimbursement, SimulateMsgWithdrawReimbursement(k, ak, sk)),
	}
}

// SimulateMsgCreatePool generates a MsgCreatePool object with all of its fields randomized.
func SimulateMsgCreatePool(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		pools := k.GetAllPools(ctx)
		// restrict number of pools to reduce gas consumptions for unbondings and redelegations
		if len(pools) > 20 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		// admin
		adminAddr := k.GetAdmin(ctx)
		var simAccount simulation.Account
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
		shieldAmount, err := simulation.RandPositiveInt(r, maxShield)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))

		// shield limit
		// No overflow would happen when converting int64 to int in this case.
		shieldLimit := sdk.NewInt(int64(simulation.RandIntBetween(r, int(maxShield.Int64()), int(maxShield.Int64())*5)))

		// sponsor
		sponsor := strings.ToLower(simulation.RandStringOfLength(r, 10))
		if _, found := k.GetPoolBySponsor(ctx, sponsor); found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		// serviceFees
		nativeAmount := account.SpendableCoins(ctx.BlockTime()).AmountOf(bondDenom)
		if !nativeAmount.IsPositive() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		nativeAmount, err = simulation.RandPositiveInt(r, nativeAmount)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		nativeServiceFees := sdk.NewCoins(sdk.NewCoin(bondDenom, nativeAmount))
		foreignAmount, err := simulation.RandPositiveInt(r, sdk.NewInt(int64(DefaultIntMax)))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		foreignServiceFees := sdk.NewCoins(sdk.NewCoin(sponsor, foreignAmount))

		serviceFees := types.MixedCoins{Native: nativeServiceFees, Foreign: foreignServiceFees}
		sponsorAcc, _ := simulation.RandomAcc(r, accs)
		description := simulation.RandStringOfLength(r, 42)

		msg := types.NewMsgCreatePool(simAccount.Address, shield, serviceFees, sponsor, sponsorAcc.Address, description, shieldLimit)

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgUpdatePool generates a MsgUpdatePool object with all of its fields randomized.
func SimulateMsgUpdatePool(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		adminAddr := k.GetAdmin(ctx)
		var simAccount simulation.Account
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
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		pool, found := k.GetPool(ctx, poolID)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// shield
		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		totalClaimed := k.GetTotalClaimed(ctx)
		poolParams := k.GetPoolParams(ctx)
		maxShield := computeMaxShield(pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield, poolParams)
		shieldAmount, err := simulation.RandPositiveInt(r, maxShield)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))

		// serviceFees
		nativeAmount := account.SpendableCoins(ctx.BlockTime()).AmountOf(bondDenom)
		if !nativeAmount.IsPositive() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		nativeAmount, err = simulation.RandPositiveInt(r, nativeAmount)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		nativeServiceFees := sdk.NewCoins(sdk.NewCoin(bondDenom, nativeAmount))
		foreignAmount, err := simulation.RandPositiveInt(r, sdk.NewInt(int64(DefaultIntMax)))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		foreignServiceFees := sdk.NewCoins(sdk.NewCoin(pool.Sponsor, foreignAmount))

		serviceFees := types.MixedCoins{Native: nativeServiceFees, Foreign: foreignServiceFees}
		description := simulation.RandStringOfLength(r, 42)

		msg := types.NewMsgUpdatePool(simAccount.Address, shield, serviceFees, poolID, description, sdk.ZeroInt())

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDepositCollateral generates a MsgDepositCollateral object with all of its fields randomized.
func SimulateMsgDepositCollateral(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		delAddr, available, found := keeper.RandomDelegation(r, k, ctx)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		var simAccount simulation.Account
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
		collateralAmount, err := simulation.RandPositiveInt(r, available)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		collateral := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), collateralAmount))

		msg := types.NewMsgDepositCollateral(simAccount.Address, collateral)

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawCollateral generates a MsgWithdrawCollateral object with all of its fields randomized.
func SimulateMsgWithdrawCollateral(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		provider, found := keeper.RandomProvider(r, k, ctx)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		var simAccount simulation.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(provider.Address) {
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
		withdrawAmount, err := simulation.RandPositiveInt(r, withdrawable)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		withdraw := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), withdrawAmount))

		msg := types.NewMsgWithdrawCollateral(simAccount.Address, withdraw)

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawRewards generates a MsgWithdrawRewards object with all of its fields randomized.
func SimulateMsgWithdrawRewards(k keeper.Keeper, ak types.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		provider, found := keeper.RandomProvider(r, k, ctx)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		var simAccount simulation.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(provider.Address) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)

		msg := types.NewMsgWithdrawRewards(simAccount.Address)

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgPurchaseShield generates a MsgPurchaseShield object with all of its fields randomized.
func SimulateMsgPurchaseShield(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		purchaser, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, purchaser.Address)
		bondDenom := sk.BondDenom(ctx)

		poolID, _, found := keeper.RandomPoolInfo(r, k, ctx)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		pool, found := k.GetPool(ctx, poolID)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		totalClaimed := k.GetTotalClaimed(ctx)
		poolParams := k.GetPoolParams(ctx)
		maxShield := computeMaxShield(pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield, poolParams)
		shieldAmount, err := simulation.RandPositiveInt(r, maxShield)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		if shieldAmount.ToDec().Mul(poolParams.ShieldFeesRate).GT(account.SpendableCoins(ctx.BlockTime()).AmountOf(bondDenom).ToDec()) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		if shieldAmount.ToDec().Mul(poolParams.ShieldFeesRate).TruncateInt().IsZero() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))

		description := simulation.RandStringOfLength(r, 100)
		msg := types.NewMsgPurchaseShield(poolID, shield, description, purchaser.Address)

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			purchaser.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper, sk types.StakingKeeper) []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightShieldClaimProposal,
			DefaultWeight:      DefaultWeightShieldClaimProposal,
			ContentSimulatorFn: SimulateShieldClaimProposalContent(k, sk),
		},
	}
}

// SimulateShieldClaimProposalContent generates random shield claim proposal content
func SimulateShieldClaimProposalContent(k keeper.Keeper, sk types.StakingKeeper) simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simulation.Account) govtypes.Content {
		if ctx.BlockHeight() < common.Update1Height {
			return nil
		}
		bondDenom := sk.BondDenom(ctx)
		purchaseList, found := keeper.RandomPurchaseList(r, k, ctx)
		if len(purchaseList.Entries) == 0 {
			return nil
		}
		i := r.Intn(len(purchaseList.Entries))
		poolID := purchaseList.PoolID
		purchaser := purchaseList.Purchaser
		purchase := purchaseList.Entries[i]
		if !found || purchase.ProtectionEndTime.Before(ctx.BlockTime()) {
			return nil
		}
		lossAmount, err := simulation.RandPositiveInt(r, purchase.Shield)
		if err != nil {
			return nil
		}
		return types.NewShieldClaimProposal(
			poolID,
			sdk.NewCoins(sdk.NewCoin(bondDenom, lossAmount)),
			purchase.PurchaseID,
			simulation.RandStringOfLength(r, 500),
			simulation.RandStringOfLength(r, 500),
			purchaser,
		)
	}
}

// SimulateMsgStakeForShield generates a MsgPurchaseShield object with all of its fields randomized.
func SimulateMsgStakeForShield(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		if ctx.BlockHeight() < common.Update1Height {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		purchaser, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, purchaser.Address)
		bondDenom := sk.BondDenom(ctx)

		poolID, _, found := keeper.RandomPoolInfo(r, k, ctx)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		pool, found := k.GetPool(ctx, poolID)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		totalCollateral := k.GetTotalCollateral(ctx)
		totalWithdrawing := k.GetTotalWithdrawing(ctx)
		totalShield := k.GetTotalShield(ctx)
		totalClaimed := k.GetTotalClaimed(ctx)
		poolParams := k.GetPoolParams(ctx)
		maxShield := computeMaxShield(pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield, poolParams)
		accountMax := sdk.OneDec().Quo(k.GetShieldStakingRate(ctx)).MulInt(account.GetCoins().AmountOf(k.BondDenom(ctx))).TruncateInt()
		max := sdk.MinInt(accountMax, maxShield)
		shieldAmount, err := simulation.RandPositiveInt(r, max)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		rate := k.GetShieldStakingRate(ctx)
		maxShieldAmt := account.SpendableCoins(ctx.BlockTime()).AmountOf(bondDenom).ToDec().Quo(rate).TruncateInt()
		if shieldAmount.GT(maxShieldAmt) {
			shieldAmount = maxShieldAmt
		}
		shield := sdk.NewCoins(sdk.NewCoin(bondDenom, shieldAmount))
		if shield.IsZero() || k.GetShieldStakingRate(ctx).MulInt(shield.AmountOf(bondDenom)).TruncateInt().IsZero() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		description := simulation.RandStringOfLength(r, 100)
		msg := types.NewMsgStakeForShield(poolID, shield, description, purchaser.Address)

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			purchaser.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgUnstakeFromShield generates a MsgUnstakeFromShield object with all of its fields randomized.
func SimulateMsgUnstakeFromShield(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		if ctx.BlockHeight() < common.Update1Height {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		bondDenom := sk.BondDenom(ctx)
		stakeForShields := k.GetAllStakeForShields(ctx)
		if len(stakeForShields) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		index := simulation.RandIntBetween(r, 0, len(stakeForShields))
		sfs := stakeForShields[index]

		withdrawable := sfs.Amount.Sub(sfs.WithdrawRequested)
		withdrawableCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, withdrawable))
		shield := simulation.RandSubsetCoins(r, withdrawableCoins)
		msg := types.NewMsgUnstakeFromShield(sfs.PoolID, shield, sfs.Purchaser)

		var account authexported.Account
		var simAcc simulation.Account
		for _, acc := range accs {
			if acc.Address.Equals(sfs.Purchaser) {
				account = ak.GetAccount(ctx, acc.Address)
				simAcc = acc
			}
		}
		if account == nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAcc.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawReimbursement generates a MsgWithdrawReimbursement object with randomized fields.
func SimulateMsgWithdrawReimbursement(k keeper.Keeper, ak types.AccountKeeper, sk types.StakingKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		if ctx.BlockHeight() < common.Update1Height {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		prPair, found := keeper.RandomMaturedProposalIDReimbursementPair(r, k, ctx)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		var simAccount simulation.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(prPair.Reimbursement.Beneficiary) {
				simAccount = simAcc
				break
			}
		}
		account := ak.GetAccount(ctx, simAccount.Address)

		msg := types.NewMsgWithdrawReimbursement(prPair.ProposalID, simAccount.Address)

		fees := sdk.Coins{}
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func computeMaxShield(pool types.Pool, totalCollateral, totalWithdrawing, totalClaimed, totalShield sdk.Int, poolParams types.PoolParams) sdk.Int {
	poolLimit := pool.ShieldLimit.Sub(pool.Shield)
	globalLimit := sdk.MinInt(totalCollateral.Sub(totalWithdrawing).Sub(totalClaimed).ToDec().Mul(poolParams.PoolShieldLimit).TruncateInt().Sub(pool.Shield),
		totalCollateral.Sub(totalWithdrawing).Sub(totalClaimed).Sub(totalShield))
	return sdk.MinInt(poolLimit, globalLimit)
}
