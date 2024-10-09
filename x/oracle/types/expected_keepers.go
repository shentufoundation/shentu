package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ParamSubspace interface {
	Get(ctx context.Context, key []byte, ptr interface{})
	Set(ctx context.Context, key []byte, param interface{})
}

type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

type CertKeeper interface {
	IsCertifier(ctx context.Context, addr sdk.AccAddress) (bool, error)
}

type DistrKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
}
