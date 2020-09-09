package cert

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/internal/keeper"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

func InitDefaultGenesis(ctx sdk.Context, k keeper.Keeper) {
	InitGenesis(ctx, k, types.DefaultGenesisState())
}

// InitGenesis initialize default parameters and the keeper's address to pubkey map.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	certifiers := data.Certifiers
	validators := data.Validators
	platforms := data.Platforms
	certificates := data.Certificates
	libraries := data.Libraries

	for _, certifier := range certifiers {
		k.SetCertifier(ctx, certifier)
	}
	if len(certifiers) > 0 {
		cert := certifiers[0].Address
		for _, platform := range platforms {
			_ = k.CertifyPlatform(ctx, cert, platform.Address, platform.Description)
		}
	}
	for _, validator := range validators {
		k.SetValidator(ctx, validator.PubKey, validator.Certifier)
	}
	for _, certificate := range certificates {
		k.SetCertificate(ctx, certificate)
	}
	for _, library := range libraries {
		k.SetLibrary(ctx, library.Address, library.Publisher)
	}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	certifiers := k.GetAllCertifiers(ctx)
	validators := k.GetAllValidators(ctx)
	platforms := k.GetAllPlatforms(ctx)
	certificates := k.GetAllCertificates(ctx)
	libraries := k.GetAllLibraries(ctx)

	return GenesisState{
		Certifiers:   certifiers,
		Validators:   validators,
		Platforms:    platforms,
		Certificates: certificates,
		Libraries:    libraries,
	}
}
