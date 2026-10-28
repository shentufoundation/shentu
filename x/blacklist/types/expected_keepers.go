package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the account keeper interface needed by the blacklist
// module. It is used to prevent blacklisting module accounts.
type AccountKeeper interface {
	// GetAccount retrieves an account for a given address, returning nil if
	// the account does not exist.
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}
