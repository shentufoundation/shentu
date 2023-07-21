package v280

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func migrateCertificate(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.CertificatesStoreKey())

	iterator := oldStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var cert types.Certificate
		cdc.MustUnmarshal(iterator.Value(), &cert)

		shentuAddr, err := common.PrefixToShentu(cert.Certifier)
		if err != nil {
			return err
		}
		cert.Certifier = shentuAddr

		bz := cdc.MustMarshalLengthPrefixed(&cert)
		oldStore.Set(iterator.Key(), bz)
	}
	return nil
}

func migrateCertifier(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.CertifiersStoreKey())

	iterator := oldStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var certifier types.Certifier
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &certifier)

		certifierAddr, err := common.PrefixToShentu(certifier.Address)
		if err != nil {
			return err
		}
		certifier.Address = certifierAddr

		proposalAddr, err := common.PrefixToShentu(certifier.Proposer)
		if err != nil {
			return err
		}
		certifier.Proposer = proposalAddr

		bz := cdc.MustMarshalLengthPrefixed(&certifier)
		oldStore.Set(iterator.Key(), bz)
	}
	return nil
}

func migrateCertifierAlias(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.CertifierAliasesStoreKey())

	iterator := oldStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var certifier types.Certifier
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &certifier)

		certifierAddr, err := common.PrefixToShentu(certifier.Address)
		if err != nil {
			return err
		}
		certifier.Address = certifierAddr

		proposalAddr, err := common.PrefixToShentu(certifier.Proposer)
		if err != nil {
			return err
		}
		certifier.Proposer = proposalAddr

		bz := cdc.MustMarshalLengthPrefixed(&certifier)
		oldStore.Set(iterator.Key(), bz)
	}
	return nil
}

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	if err := migrateCertificate(store, cdc); err != nil {
		return err
	}

	if err := migrateCertifier(store, cdc); err != nil {
		return err
	}

	if err := migrateCertifierAlias(store, cdc); err != nil {
		return err
	}
	return nil
}
