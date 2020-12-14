package keeper

import (
	"github.com/certikfoundation/shentu/x/mint"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/certikfoundation/shentu/x/mint/internal/types"
)

type Keeper struct {
	mintKeeper.Keeper
	dk            types.DistributionKeeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	shieldKeeper  types.ShieldKeeper
}

// NewKeeper implements the wrapper newkeeper on top of the original newkeeper with distribution, supply and staking keeper.
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	sk types.StakingKeeper, ak types.AccountKeeper, bk types.BankKeeper, distributionKeeper types.DistributionKeeper, shieldKeeper types.ShieldKeeper,
	feeCollectorName string) Keeper {
	return Keeper{
		Keeper:        mintKeeper.NewKeeper(cdc, key, paramSpace, sk, ak, bk, feeCollectorName),
		dk:            distributionKeeper,
		accountKeeper: ak,
		stakingKeeper: sk,
		shieldKeeper:  shieldKeeper,
	}
}

// SendToCommunityPool sends coins to the community pool using FundCommunityPool.
func (k Keeper) SendToCommunityPool(ctx sdk.Context, amount sdk.Coins) error {
	if amount.AmountOf(k.stakingKeeper.BondDenom(ctx)).Equal(sdk.ZeroInt()) {
		return nil
	}
	mintAddress := k.accountKeeper.GetModuleAddress(mint.ModuleName)
	return k.dk.FundCommunityPool(ctx, amount, mintAddress)
}

// SendToShieldRewards sends coins to the shield rewards using FundShieldBlockRewards.
func (k Keeper) SendToShieldRewards(ctx sdk.Context, amount sdk.Coins) error {
	if amount.AmountOf(k.stakingKeeper.BondDenom(ctx)).Equal(sdk.ZeroInt()) {
		return nil
	}
	mintAddress := k.accountKeeper.GetModuleAddress(mint.ModuleName)
	return k.shieldKeeper.FundShieldBlockRewards(ctx, amount, mintAddress)
}

// GetCommunityPoolRatio returns the current ratio of the community pool compared to the total supply.
func (k Keeper) GetCommunityPoolRatio(ctx sdk.Context) sdk.Dec {
	communityPool := k.dk.GetFeePool(ctx).CommunityPool
	for _, coin := range communityPool {
		totalBondedTokensDec := k.StakingTokenSupply(ctx).ToDec()
		if coin.Denom == k.stakingKeeper.BondDenom(ctx) {
			ratio := coin.Amount.Quo(totalBondedTokensDec)
			return ratio
		}
	}
	return sdk.NewDec(0)
}

// GetShieldStakeForShieldPoolRatio returns the current ratio of the community pool compared to the total supply.
func (k Keeper) GetShieldStakeForShieldPoolRatio(ctx sdk.Context) sdk.Dec {
	pool := k.shieldKeeper.GetGlobalShieldStakingPool(ctx)
	totalBondedTokensDec := k.StakingTokenSupply(ctx).ToDec()
	return pool.ToDec().Quo(totalBondedTokensDec)
}

// GetPoolMint returns Coins that are about to be minted towards the community pool.
func (k Keeper) GetPoolMint(ctx sdk.Context, ratio sdk.Dec, mintedCoin sdk.Coin) sdk.Coins {
	communityPoolMintDec := ratio.MulInt(mintedCoin.Amount)
	amount := communityPoolMintDec.TruncateInt()
	return sdk.Coins{sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), amount)}
}
