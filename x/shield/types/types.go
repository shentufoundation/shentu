package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPool creates a new project pool.
func NewPool(id uint64, description, sponsor string, sponsorAddress sdk.AccAddress, shieldLimit sdk.Int, shield sdk.Int) Pool {
	return Pool{
		Id:          id,
		Description: description,
		Sponsor:     sponsor,
		SponsorAddr: sponsorAddress.String(),
		ShieldLimit: shieldLimit,
		Active:      true,
		Shield:      shield,
	}
}

// NewProvider creates a new provider object.
func NewProvider(addr sdk.AccAddress) Provider {
	return Provider{
		Address:          addr.String(),
		DelegationBonded: sdk.ZeroInt(),
		Collateral:       sdk.ZeroInt(),
		TotalLocked:      sdk.ZeroInt(),
		Withdrawing:      sdk.ZeroInt(),
	}
}

// NewPurchase creates a new purchase object.
func NewPurchase(purchaseID uint64, protectionEndTime, deletionTime time.Time, description string, shield sdk.Int, nativeServiceFee sdk.DecCoins) Purchase {
	return Purchase{
		PurchaseId:        purchaseID,
		ProtectionEndTime: protectionEndTime,
		DeletionTime:      deletionTime,
		Description:       description,
		Shield:            shield,
		NativeServiceFee:  nativeServiceFee,
	}
}

// NewPurchaseList creates a new purchase list.
func NewPurchaseList(poolID uint64, purchaser sdk.AccAddress, purchases []Purchase) PurchaseList {
	return PurchaseList{
		PoolId:    poolID,
		Purchaser: purchaser.String(),
		Entries:   purchases,
	}
}

// NewWithdraw creates a new withdraw object.
func NewWithdraw(addr sdk.AccAddress, amount sdk.Int, completionTime time.Time) Withdraw {
	return Withdraw{
		Address:        addr.String(),
		Amount:         amount,
		CompletionTime: completionTime,
	}
}

func NewShieldStaking(poolID uint64, purchaser sdk.AccAddress, amount sdk.Int) ShieldStaking {
	return ShieldStaking{
		PoolId:            poolID,
		Purchaser:         purchaser.String(),
		Amount:            amount,
		WithdrawRequested: sdk.NewInt(0),
	}
}
