package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	Get(ctx context.Context, key []byte, ptr interface{})
	Set(ctx context.Context, key []byte, param interface{})
	HasKeyTable() bool
}

type CertKeeper interface {
	IsBountyAdmin(ctx context.Context, addr sdk.AccAddress) bool
}
