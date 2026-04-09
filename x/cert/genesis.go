package cert

import (
	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitDefaultGenesis(ctx sdk.Context, k keeper.Keeper) {
	InitGenesis(ctx, k, *types.DefaultGenesisState())
}

// InitGenesis initialize default parameters and the keeper's address to pubkey map.
// Platform and library fields in GenesisState are accepted but ignored; they have been
// removed from the module's runtime state. The migration step handles deletion of any
// existing on-chain platform/library entries.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	for _, certifier := range data.Certifiers {
		if err := k.SetCertifier(ctx, certifier); err != nil {
			panic(err)
		}
	}
	for _, certificate := range data.Certificates {
		if err := k.SetCertificate(ctx, certificate); err != nil {
			panic(err)
		}
	}
	if err := k.SetNextCertificateID(ctx, data.NextCertificateId); err != nil {
		panic(err)
	}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	certifiers := k.GetAllCertifiers(ctx)
	certificates := k.GetAllCertificates(ctx)
	nextCertificateID, _ := k.GetNextCertificateID(ctx)

	return &types.GenesisState{
		Certifiers:        certifiers,
		Certificates:      certificates,
		NextCertificateId: nextCertificateID,
	}
}
