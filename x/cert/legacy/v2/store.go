package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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

		transistCertContent(&cert)

		bz := cdc.MustMarshalLengthPrefixed(&cert)
		oldStore.Set(iterator.Key(), bz)
	}
	return nil
}

func transistCertContent(certificate *types.Certificate) {
	switch certificate.GetContent().(type) {
	case *types.OracleOperator:
		contentShentu, err := common.PrefixToShentu(certificate.GetContentString())
		if err != nil {
			return
		}
		content := types.OracleOperator{Content: contentShentu}
		setContentAny(certificate, &content)
	case *types.ShieldPoolCreator:
		contentShentu, err := common.PrefixToShentu(certificate.GetContentString())
		if err != nil {
			return
		}
		content := types.ShieldPoolCreator{Content: contentShentu}
		setContentAny(certificate, &content)
	case *types.Identity:
		contentShentu, err := common.PrefixToShentu(certificate.GetContentString())
		if err != nil {
			return
		}
		content := types.Identity{Content: contentShentu}
		setContentAny(certificate, &content)
	}
}

func setContentAny(certificate *types.Certificate, content types.Content) {
	any, err := codectypes.NewAnyWithValue(content)
	if err != nil {
		panic(err)
	}
	certificate.Content = any
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

func migrateLibrary(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.LibrariesStoreKey())

	iterator := oldStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var library types.Library
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &library)

		libraryAddr, err := common.PrefixToShentu(library.Address)
		if err != nil {
			return err
		}
		library.Address = libraryAddr

		publisher, err := common.PrefixToShentu(library.Publisher)
		if err != nil {
			return err
		}
		library.Publisher = publisher

		bz := cdc.MustMarshalLengthPrefixed(&library)
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

	if err := migrateLibrary(store, cdc); err != nil {
		return err
	}
	return nil
}
