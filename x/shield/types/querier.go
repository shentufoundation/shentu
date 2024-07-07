package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryProvider  = "provider"
	QueryProviders = "providers"
	QueryStatus    = "status"
)

type QueryResStatus struct {
	RemainingServiceFees sdk.DecCoins `json:"remaining_service_fees" yaml:"remaining_service_fees"`
}

func NewQueryResStatus(remainingServiceFees sdk.DecCoins) QueryResStatus {
	return QueryResStatus{
		RemainingServiceFees: remainingServiceFees,
	}
}

// QueryPaginationParams provides basic pagination parameters
// for queries in shield module.
type QueryPaginationParams struct {
	Page  int
	Limit int
}
