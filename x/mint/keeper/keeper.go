package keeper

import (
	"context"

	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/shentufoundation/shentu/v2/x/mint/types"
)

type Keeper struct {
	mintkeeper.Keeper
	dk            types.DistributionKeeper
	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
}

// NewKeeper implements the wrapper newkeeper on top of the original newkeeper with distribution, supply and staking keeper.
func NewKeeper(
	cdc codec.BinaryCodec, storeService storetypes.KVStoreService,
	sk types.StakingKeeper, ak types.AccountKeeper, bk types.BankKeeper, distributionKeeper types.DistributionKeeper,
	feeCollectorName string, authority string) Keeper {
	return Keeper{
		Keeper:        mintkeeper.NewKeeper(cdc, storeService, sk, ak, bk, feeCollectorName, authority),
		dk:            distributionKeeper,
		accountKeeper: ak,
		stakingKeeper: sk,
	}
}

// SendToCommunityPool sends coins to the community pool using FundCommunityPool.
func (k Keeper) SendToCommunityPool(ctx context.Context, amount sdk.Coins) error {
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	if amount.AmountOf(bondDenom).Equal(math.ZeroInt()) {
		return nil
	}
	mintAddress := k.accountKeeper.GetModuleAddress(minttypes.ModuleName)
	return k.dk.FundCommunityPool(ctx, amount, mintAddress)
}

// GetCommunityPoolRatio returns the current ratio of the community pool compared to the total supply.
func (k Keeper) GetCommunityPoolRatio(ctx context.Context) (math.LegacyDec, error) {
	communityPool := k.dk.GetFeePool(ctx).CommunityPool
	for _, coin := range communityPool {
		denom, err := k.stakingKeeper.BondDenom(ctx)
		if err != nil {
			return math.LegacyDec{}, err
		}
		supply, err := k.StakingTokenSupply(ctx)
		if err != nil {
			return math.LegacyZeroDec(), err
		}
		totalBondedTokensDec := math.LegacyNewDecFromInt(supply)
		if coin.Denom == denom {
			ratio := coin.Amount.Quo(totalBondedTokensDec)
			return ratio, nil
		}
	}
	return math.LegacyZeroDec(), nil
}

// GetPoolMint returns Coins that are about to be minted towards the community pool.
func (k Keeper) GetPoolMint(ctx context.Context, ratio math.LegacyDec, mintedCoin sdk.Coin) (sdk.Coin, error) {
	communityPoolMintDec := ratio.MulInt(mintedCoin.Amount)
	amount := communityPoolMintDec.TruncateInt()
	denom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(denom, amount), nil
}
