package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryPoolByID            = "pool_id"
	QueryPoolBySponsor       = "pool_sponsor"
	QueryPools               = "pools"
	QueryPurchase            = "purchase"
	QueryPurchaseList        = "purchase_list"
	QueryPurchaserPurchases  = "purchaser_purchases"
	QueryPoolPurchases       = "pool_purchases"
	QueryProviderCollaterals = "provider_collaterals"
	QueryPoolCollaterals     = "pool_collaterals"
	QueryProvider            = "provider"
	QueryPoolParams          = "pool_params"
	QueryClaimParams         = "claim_params"
	QueryStatus              = "status"
)

type QueryResStatus struct {
	TotalCollateral      sdk.Int       `json:"total_collateral" yaml:"total_collateral"`
	TotalShield          sdk.Int       `json:"total_shield" yaml:"total_shield"`
	TotalWithdrawing     sdk.Int       `json:"total_withdrawing" yaml:"total_withdrawing"`
	CurrentServiceFees   MixedDecCoins `json:"current_service_fees" yaml:"current_service_fees"`
	RemainingServiceFees MixedDecCoins `json:"remaining_service_fees" yaml:"remaining_service_fees"`
}

func NewQueryResStatus(totalCollateral, totalShield, totalWithdrawing sdk.Int, currentServiceFees, remainingServiceFees MixedDecCoins) QueryResStatus {
	return QueryResStatus{
		TotalCollateral:      totalCollateral,
		TotalShield:          totalShield,
		TotalWithdrawing:     totalWithdrawing,
		CurrentServiceFees:   currentServiceFees,
		RemainingServiceFees: remainingServiceFees,
	}
}
