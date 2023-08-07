package v2_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/common"
	v2 "github.com/shentufoundation/shentu/v2/x/cert/legacy/v2"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func makeCertificate(certType string) types.Certificate {
	_, _, addr := testdata.KeyTestPubAddr()
	certifier, _ := common.PrefixToCertik(addr.String())

	content := types.AssembleContent(certType, certifier)
	msg, ok := content.(proto.Message)
	if !ok {
		panic(fmt.Errorf("%T does not implement proto.Message", content))
	}
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}

	return types.Certificate{
		CertificateId:      rand.Uint64(),
		Content:            any,
		CompilationContent: &types.CompilationContent{Compiler: "", BytecodeHash: ""},
		Description:        "for test",
		Certifier:          certifier,
	}
}

func makeCertifier(certifierAddr, alias string) types.Certifier {
	_, _, addr := testdata.KeyTestPubAddr()
	proposer, _ := common.PrefixToCertik(addr.String())

	return types.Certifier{
		Address:     certifierAddr,
		Alias:       alias,
		Proposer:    proposer,
		Description: "unit test",
	}
}

func saveLibrary(cdc codec.BinaryCodec, storeKey sdk.KVStore) types.Library {
	_, _, libraryAddr := testdata.KeyTestPubAddr()
	libraryAddrStr, _ := common.PrefixToCertik(libraryAddr.String())

	_, _, publisherAddr := testdata.KeyTestPubAddr()
	publisherStr, _ := common.PrefixToCertik(publisherAddr.String())

	library := types.Library{
		Address:   libraryAddrStr,
		Publisher: publisherStr,
	}

	bz := cdc.MustMarshalLengthPrefixed(&library)
	storeKey.Set(types.LibraryStoreKey(libraryAddr), bz)
	return library
}

func TestMigrateStore(t *testing.T) {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)

	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cdc := shentuapp.MakeEncodingConfig().Marshaler

	oldCertificate := makeCertificate("IDENTITY")

	oldCertifier := makeCertifier(oldCertificate.Certifier, "test_certifier")

	app.CertKeeper.SetCertificate(ctx, oldCertificate)
	app.CertKeeper.SetCertifier(ctx, oldCertifier)

	store := ctx.KVStore(app.GetKey(types.StoreKey))

	oldLibrary := saveLibrary(cdc, store)

	err := v2.MigrateStore(ctx, app.GetKey(types.StoreKey), cdc)
	require.NoError(t, err)

	//check for Certificate
	var cert types.Certificate
	bz := store.Get(types.CertificateStoreKey(oldCertificate.CertificateId))
	cdc.MustUnmarshal(bz, &cert)

	newCertifierAddr, _ := common.PrefixToShentu(oldCertificate.Certifier)
	require.Equal(t, newCertifierAddr, cert.Certifier)

	var certifier types.Certifier
	certifierAcc, err := sdk.AccAddressFromBech32(cert.Certifier)
	require.NoError(t, err)
	bz = store.Get(types.CertifierStoreKey(certifierAcc))
	cdc.MustUnmarshalLengthPrefixed(bz, &certifier)
	require.Equal(t, newCertifierAddr, certifier.Address)
	require.Equal(t, newCertifierAddr, cert.GetContentString())

	newCertifierProposer, _ := common.PrefixToShentu(oldCertifier.Proposer)
	require.Equal(t, newCertifierProposer, certifier.Proposer)

	bz = store.Get(types.CertifierAliasStoreKey(oldCertifier.Alias))
	var certifierAlias types.Certifier
	cdc.MustUnmarshalLengthPrefixed(bz, &certifierAlias)
	require.Equal(t, newCertifierAddr, certifierAlias.Address)
	require.Equal(t, newCertifierProposer, certifierAlias.Proposer)

	libraryAddr, _ := sdk.AccAddressFromBech32(oldLibrary.Address)
	bz = store.Get(types.LibraryStoreKey(libraryAddr))
	var library types.Library
	cdc.MustUnmarshalLengthPrefixed(bz, &library)
	newLibraryAddr, _ := common.PrefixToShentu(oldLibrary.Address)
	newPublisherAddr, _ := common.PrefixToShentu(oldLibrary.Publisher)
	require.Equal(t, library.Address, newLibraryAddr)
	require.Equal(t, library.Publisher, newPublisherAddr)
}
