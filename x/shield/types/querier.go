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
	QueryGlobalState         = "global_state"
)

type QueryResShieldState struct {
	TotalCollateral    sdk.Int       `json:"total_collateral" yaml:"total_collateral"`
	TotalShield        sdk.Int       `json:"total_shield" yaml:"total_shield"`
	TotalWithdrawing   sdk.Int       `json:"total_withdrawing" yaml:"total_withdrawing"`
	CurrentServiceFees MixedDecCoins `json:"current_service_fees" yaml:"current_service_fees"`
	ServiceFeesLeft    MixedDecCoins `json:"service_fees_left" yaml:"service_fees_left"`
}

func NewQueryResShieldState(totalCollateral, totalShield, totalWithdrawing sdk.Int, currentServiceFees, serviceFeesLeft MixedDecCoins) QueryResShieldState {
	return QueryResShieldState{
		TotalCollateral:    totalCollateral,
		TotalShield:        totalShield,
		TotalWithdrawing:   totalWithdrawing,
		CurrentServiceFees: currentServiceFees,
		ServiceFeesLeft:    serviceFeesLeft,
	}
}
