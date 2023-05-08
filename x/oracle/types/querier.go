package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Querier routes for the oracle module
const (
	QueryOperator    = "operator"
	QueryOperators   = "operators"
	QueryWithdrawals = "withdrawals"
	QueryTask        = "task"
	QueryResponse    = "response"
)

type QueryTaskParams struct {
	Contract string
	Function string
}

// NewQueryTaskParams returns a QueryTaskParams object.
func NewQueryTaskParams(contract string, function string) QueryTaskParams {
	return QueryTaskParams{
		Contract: contract,
		Function: function,
	}
}

type QueryResponseParams struct {
	Contract string
	Function string
	Operator sdk.AccAddress
}

// NewQueryResponseParams returns a QueryResponseParams.
func NewQueryResponseParams(contract string, function string, operator sdk.AccAddress) QueryResponseParams {
	return QueryResponseParams{
		Contract: contract,
		Function: function,
		Operator: operator,
	}
}
