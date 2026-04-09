package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	AccountKeeper interface {
		GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	}

	BankKeeper interface {
		SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	}
)
