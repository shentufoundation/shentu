package cert

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/cert/keeper"
	"github.com/certikfoundation/shentu/x/cert/types"
)

func InitDefaultGenesis(ctx sdk.Context, k keeper.Keeper) {
	InitGenesis(ctx, k, *types.DefaultGenesisState())
}

// InitGenesis initialize default parameters and the keeper's address to pubkey map.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	certifiers := data.Certifiers
	validators := data.Validators

	for _, certifier := range certifiers {
		k.SetCertifier(ctx, certifier)
	}
	for _, validator := range validators {
		pk, ok := validator.Pubkey.GetCachedValue().(cryptotypes.PubKey)
		if !ok {
			panic(sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into cryto.PubKey %T", validator.Pubkey))
		}
		certifierAddr, err := sdk.AccAddressFromBech32(validator.Certifier)
		if err != nil {
			panic(err)
		}

		k.SetValidator(ctx, pk, certifierAddr)
	}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	certifiers := k.GetAllCertifiers(ctx)
	validators := k.GetAllValidators(ctx)

	return &types.GenesisState{
		Certifiers: certifiers,
		Validators: validators,
	}
}
