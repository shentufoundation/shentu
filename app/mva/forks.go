package mva

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkauthkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"

	authtypes "github.com/certikfoundation/shentu/v2/x/auth/types"
	bankkeeper "github.com/certikfoundation/shentu/v2/x/bank/keeper"
	stakingkeeper "github.com/certikfoundation/shentu/v2/x/staking/keeper"
)

func RunForkLogic(ctx sdk.Context, ak *sdkauthkeeper.AccountKeeper, bk bankkeeper.Keeper, sk *stakingkeeper.Keeper) {
	ctx.Logger().Info("Applying Shentu MVA upgrade." +
		" Fixing Shentu MVA accounts to correctly update to the right delegation tracking.")
	FixAccounts(ctx, ak, bk, sk)
}

func FixAccounts(ctx sdk.Context, ak *sdkauthkeeper.AccountKeeper, bk bankkeeper.Keeper, sk *stakingkeeper.Keeper) {
	ak.IterateAccounts(ctx, func(account sdktypes.AccountI) (stop bool) {
		mvacc, ok := account.(*authtypes.ManualVestingAccount)
		if !ok {
			return false
		}

		wb, err := MigrateAccount(ctx, mvacc, bk, sk)
		if err != nil {
			panic(err)
		}

		if wb == nil {
			return false
		}

		mvacc, ok = wb.(*authtypes.ManualVestingAccount)

		if !ok {
			panic("couldn't unmarshal resulting account to MVA")
		}

		ak.SetAccount(ctx, mvacc)
		return false
	})

}

func MigrateAccount(ctx sdk.Context, account sdktypes.AccountI, bk bankkeeper.Keeper, sk *stakingkeeper.Keeper) (sdktypes.AccountI, error) {
	bondDenom := sk.BondDenom(ctx)

	asVesting, ok := account.(exported.VestingAccount)
	if !ok {
		return nil, nil
	}

	addr := account.GetAddress()
	balance := bk.GetAllBalances(ctx, addr)

	delegations := sk.GetDelegatorDelegations(ctx, addr, math.MaxUint16)

	delegationsSum := sdk.ZeroInt()
	for _, d := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
		if err != nil {
			panic("cannot convert validator address to sdk.ValAddress")
		}
		val, found := sk.GetValidator(ctx, valAddr)
		if !found {
			panic(fmt.Sprintf("cannot find the validator %s", d.ValidatorAddress))
		}
		tokenAmount := val.TokensFromShares(d.Shares).TruncateInt()
		delegationsSum = delegationsSum.Add(tokenAmount)
	}

	unbondings := sk.GetUnbondingDelegations(ctx, addr, math.MaxUint16)

	unbondingSum := sdk.ZeroInt()
	for _, u := range unbondings {
		for _, e := range u.Entries {
			unbondingSum = unbondingSum.Add(e.Balance)
		}
	}

	delegationsSum = delegationsSum.Add(unbondingSum)
	delegationCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, delegationsSum))

	asVesting, ok = resetVestingDelegatedBalances(asVesting)
	if !ok {
		return nil, nil
	}

	// balance before any delegation includes balance of delegation
	for _, coin := range delegationCoins {
		balance = balance.Add(coin)
	}

	asVesting.TrackDelegation(ctx.BlockTime(), balance, delegationCoins)

	return asVesting.(sdktypes.AccountI), nil
}

func resetVestingDelegatedBalances(evacct exported.VestingAccount) (exported.VestingAccount, bool) {
	switch vacct := evacct.(type) {
	case *authtypes.ManualVestingAccount:
		vacct.DelegatedVesting = sdk.NewCoins()
		vacct.DelegatedFree = sdk.NewCoins()
		return vacct, true
	default:
		return nil, false
	}
}
