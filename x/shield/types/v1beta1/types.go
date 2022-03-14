package v1beta1

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultShieldRate = sdk.NewDec(5)
)

// NewPool creates a new project pool.
func NewPool(id uint64, description string, sponsorAddress sdk.AccAddress, shield sdk.Int, shieldRate sdk.Dec, shieldLimit sdk.Int) Pool {
	return Pool{
		Id:          id,
		Description: description,
		SponsorAddr: sponsorAddress.String(),
		Active:      true,
		Shield:      shield,
		ShieldRate:  shieldRate,
		ShieldLimit: shieldLimit,
	}
}

// NewReserve initializes a reserve
func NewReserve() Reserve {
	return Reserve{
		Amount: sdk.ZeroDec(),
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

// NewWithdraw creates a new withdraw object.
func NewWithdraw(addr sdk.AccAddress, amount sdk.Int, completionTime time.Time) Withdraw {
	return Withdraw{
		Address:        addr.String(),
		Amount:         amount,
		CompletionTime: completionTime,
	}
}

func NewPurchase(poolID uint64, purchaser sdk.AccAddress, description string, amount, shield sdk.Int) Purchase {
	return Purchase{
		PoolId:      poolID,
		Purchaser:   purchaser.String(),
		Description: description,
		Amount:      amount,
		Shield:      shield,
	}
}
