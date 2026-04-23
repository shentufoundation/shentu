package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// SetCertifier sets a certifier.
func (k Keeper) SetCertifier(ctx context.Context, certifier types.Certifier) error {
	certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
	if err != nil {
		return err
	}
	return k.Certifiers.Set(ctx, certifierAddr, certifier)
}

// deleteCertifier deletes a certifier.
func (k Keeper) deleteCertifier(ctx context.Context, certifierAddress sdk.AccAddress) error {
	return k.Certifiers.Remove(ctx, certifierAddress)
}

// IsCertifier checks if an address is a certifier.
func (k Keeper) IsCertifier(ctx context.Context, address sdk.AccAddress) (bool, error) {
	return k.Certifiers.Has(ctx, address)
}

// GetCertifier returns the certification information for a certifier.
func (k Keeper) GetCertifier(ctx context.Context, certifierAddress sdk.AccAddress) (types.Certifier, error) {
	certifier, err := k.Certifiers.Get(ctx, certifierAddress)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Certifier{}, types.ErrCertifierNotExists
		}
		return types.Certifier{}, err
	}
	return certifier, nil
}

// UpdateCertifier applies an add/remove certifier operation through one keeper path.
func (k Keeper) UpdateCertifier(ctx context.Context, operation types.AddOrRemove, certifier types.Certifier) error {
	certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
	if err != nil {
		return err
	}

	switch operation {
	case types.Add:
		isCertifier, err := k.IsCertifier(ctx, certifierAddr)
		if err != nil {
			return err
		}
		if isCertifier {
			return types.ErrCertifierAlreadyExists
		}
		return k.SetCertifier(ctx, certifier)
	case types.Remove:
		isCertifier, err := k.IsCertifier(ctx, certifierAddr)
		if err != nil {
			return err
		}
		if !isCertifier {
			return types.ErrCertifierNotExists
		}
		certifiers := k.GetAllCertifiers(ctx)
		if len(certifiers) == 1 {
			return types.ErrOnlyOneCertifier
		}
		return k.deleteCertifier(ctx, certifierAddr)
	default:
		return types.ErrAddOrRemove
	}
}

// IterateAllCertifiers iterates over the all the stored certifiers and performs a callback function.
func (k Keeper) IterateAllCertifiers(ctx context.Context, callback func(certifier types.Certifier) (stop bool)) {
	err := k.Certifiers.Walk(ctx, nil, func(_ sdk.AccAddress, certifier types.Certifier) (bool, error) {
		return callback(certifier), nil
	})
	if err != nil {
		panic(err)
	}
}

// GetAllCertifiers gets all certifiers.
func (k Keeper) GetAllCertifiers(ctx context.Context) types.Certifiers {
	certifiers := types.Certifiers{}
	k.IterateAllCertifiers(ctx, func(certifier types.Certifier) bool {
		certifiers = append(certifiers, certifier)
		return false
	})
	return certifiers
}
