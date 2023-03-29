package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool
	IsCertified(ctx sdk.Context, content string, certType string) bool
	GetAllCertifiers(ctx sdk.Context) certtypes.Certifiers
}

type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

type StakingKeeper interface {
	BondDenom(ctx sdk.Context) (res string)
}
