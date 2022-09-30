package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/shentufoundation/shentu/v2/x/mint/types"
)

type Keeper struct {
	mintkeeper.Keeper
	dk            types.DistributionKeeper
	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
	shieldKeeper  types.ShieldKeeper
}

// NewKeeper implements the wrapper newkeeper on top of the original newkeeper with distribution, supply and staking keeper.
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	sk types.StakingKeeper, ak types.AccountKeeper, bk types.BankKeeper, distributionKeeper types.DistributionKeeper, shieldKeeper types.ShieldKeeper,
	feeCollectorName string) Keeper {
	return Keeper{
		Keeper:        mintkeeper.NewKeeper(cdc, key, paramSpace, sk, ak, bk, feeCollectorName),
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
	mintAddress := k.accountKeeper.GetModuleAddress(minttypes.ModuleName)
	return k.dk.FundCommunityPool(ctx, amount, mintAddress)
}

// SendToShieldRewards sends coins to the shield rewards using FundShieldBlockRewards.
func (k Keeper) SendToShieldRewards(ctx sdk.Context, amount sdk.Coins) error {
	if amount.AmountOf(k.stakingKeeper.BondDenom(ctx)).Equal(sdk.ZeroInt()) {
		return nil
	}
	mintAddress := k.accountKeeper.GetModuleAddress(minttypes.ModuleName)
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

// GetShieldRatio returns the current ratio of
// shield staking pool compared to the total supply.
func (k Keeper) GetShieldRatio(ctx sdk.Context) sdk.Dec {
	return k.shieldKeeper.GetShieldBlockRewardRatio(ctx)
}

// GetPoolMint returns Coins that are about to be minted towards the community pool.
func (k Keeper) GetPoolMint(ctx sdk.Context, ratio sdk.Dec, mintedCoin sdk.Coin) sdk.Coin {
	communityPoolMintDec := ratio.MulInt(mintedCoin.Amount)
	amount := communityPoolMintDec.TruncateInt()
	return sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), amount)
}
