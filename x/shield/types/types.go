package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Pool struct {
	PoolID           uint64
	Active           bool
	Description      string
	Sponsor          string
	Premium          MixedDecCoins
	StartBlockHeight int64
	TotalCollateral  sdk.Coins
	Available        sdk.Int
	Shield           sdk.Coins
	EndTime          int64
	EndBlockHeight   int64
}

func NewPool(
	shield sdk.Coins, deposit MixedDecCoins, sponsor string,
	endTime, startBlockHeight, endBlockHeight int64, id uint64) Pool {
	return Pool{
		Shield:           shield,
		Premium:          deposit,
		Sponsor:          sponsor,
		Active:           true,
		TotalCollateral:  shield,
		EndTime:          endTime,
		StartBlockHeight: startBlockHeight,
		EndBlockHeight:   endBlockHeight,
		PoolID:           id,
	}
}

type Collateral struct {
	PoolID            uint64
	Provider          sdk.AccAddress
	Amount            sdk.Coins
	LockedCollaterals []LockedCollateral
}

func NewCollateral(pool Pool, provider sdk.AccAddress, amount sdk.Coins) Collateral {
	return Collateral{
		PoolID:   pool.PoolID,
		Provider: provider,
		Amount:   amount,
	}
}

type PendingPayout struct {
	Amount sdk.Dec
	ToAddr string
}

type PendingPayouts []PendingPayout

func NewPendingPayouts(amount sdk.Dec, to string) PendingPayout {
	return PendingPayout{
		Amount: amount,
		ToAddr: to,
	}
}

// Provider tracks A or C's total delegation, total collateral,
// and rewards.
type Provider struct {
	Address          sdk.AccAddress
	DelegationBonded sdk.Coins
	Collateral       sdk.Coins
	TotalLocked      sdk.Coins
	Rewards          MixedDecCoins
}

func NewProvider(addr sdk.AccAddress) Provider {
	return Provider{
		Address:          addr,
		DelegationBonded: sdk.Coins{},
		Collateral:       sdk.Coins{},
		TotalLocked:      sdk.Coins{},
	}
}

type Purchase struct {
	TxHash             []byte
	PoolID             uint64
	Shield             sdk.Coins
	StartBlockHeight   int64
	ProtectionEndTime  time.Time
	ClaimPeriodEndTime time.Time
	Description        string
	Purchaser          sdk.AccAddress
}

func NewPurchase(
	txhash []byte, poolID uint64, shield sdk.Coins, startBlockHeight int64, protectionEndTime, claimPeriodEndTime time.Time,
	description string, purchaser sdk.AccAddress) Purchase {
	return Purchase{
		TxHash:             txhash,
		PoolID:             poolID,
		Shield:             shield,
		StartBlockHeight:   startBlockHeight,
		ProtectionEndTime:  protectionEndTime,
		ClaimPeriodEndTime: claimPeriodEndTime,
		Description:        description,
		Purchaser:          purchaser,
	}
}

// Withdrawal stores an ongoing withdrawal of pool collateral.
type Withdrawal struct {
	PoolID  uint64         `json:"pool_id" yaml:"pool_id"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Amount  sdk.Coins      `json:"amount" yaml:"amount"`
}

func NewWithdrawal(poolID uint64, addr sdk.AccAddress, amount sdk.Coins) Withdrawal {
	return Withdrawal{
		PoolID:  poolID,
		Address: addr,
		Amount:  amount,
	}
}
