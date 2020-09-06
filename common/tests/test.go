package tests

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MakeTestTimestampCurrent creates a current timestamp.
func MakeTestTimestampCurrent() time.Time {
	return time.Now().UTC()
}

// MakeTestPrivKey creates an instant private key.
func MakeTestPrivKey() crypto.PrivKey {
	return ed25519.GenPrivKey()
}

// MakeTestPubKeyFromPrivKey creates a public key from a private key.
func MakeTestPubKeyFromPrivKey(privKey crypto.PrivKey) crypto.PubKey {
	return privKey.PubKey()
}

// MakeTestPubKey creates an instant public key.
func MakeTestPubKey() crypto.PubKey {
	return ed25519.GenPrivKey().PubKey()
}

// MakeTestAccAddress creates an instant account address.
func MakeTestAccAddress() sdk.AccAddress {
	return sdk.AccAddress(MakeTestPubKey().Address())
}

// MakeTestAccAddressFromPubKey creates a account address from a given public key.
func MakeTestAccAddressFromPubKey(pubKey crypto.PubKey) sdk.AccAddress {
	return sdk.AccAddress(pubKey.Address())
}

// MakeTestValAddressFromPubKey creates a validator address from a given public key.
func MakeTestValAddressFromPubKey(pubKey crypto.PubKey) sdk.ValAddress {
	return sdk.ValAddress(pubKey.Address())
}

// MakeTestConsAddressFromPubKey creates a consensus address from a given public key.
func MakeTestConsAddressFromPubKey(pubKey crypto.PubKey) sdk.ConsAddress {
	return sdk.ConsAddress(pubKey.Address())
}

// MakeTestAccount creates an instant key-pair and associated addresses.
func MakeTestAccount() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress, sdk.ValAddress, sdk.ConsAddress) {
	privKey := MakeTestPrivKey()
	pubKey := MakeTestPubKeyFromPrivKey(privKey)
	return privKey, pubKey, MakeTestAccAddressFromPubKey(pubKey), MakeTestValAddressFromPubKey(pubKey), MakeTestConsAddressFromPubKey(pubKey)
}

// CodecRegister is the alias for module codec register function.
type CodecRegister func(*codec.Codec)

// MakeTestCodec creates an instant codec and registers the modules in it for testing.
func MakeTestCodec(codecRegisters []CodecRegister) *codec.Codec {
	cdc := codec.New()
	for _, RegisterCodec := range codecRegisters {
		RegisterCodec(cdc)
	}
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// MakeTestDB creates an instant mem database.
func MakeTestDB() *dbm.MemDB {
	return dbm.NewMemDB()
}

// MakeTestStore creates an instant private multi-store.
// Note: StoreKey mustn't be newly created / different from keeper's.
func MakeTestStore(
	db *dbm.MemDB,
	storeKeys map[string]*sdk.KVStoreKey,
	transientStoreKeys map[string]*sdk.TransientStoreKey,
) sdk.CommitMultiStore {
	ms := store.NewCommitMultiStore(db)
	for _, storeKey := range storeKeys {
		ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	}
	for _, transientStoreKey := range transientStoreKeys {
		ms.MountStoreWithDB(transientStoreKey, sdk.StoreTypeTransient, db)
	}
	err := ms.LoadLatestVersion()
	if err != nil {
		panic("Oh nooooooooooooooooo!")
	}
	return ms
}

// MakeTestContext creates an context instance for testing.
func MakeTestContext(
	kvStoreKeys map[string]*sdk.KVStoreKey,
	transientStoreKeys map[string]*sdk.TransientStoreKey,
	isCheckTx bool,
	chainID string,
) sdk.Context {
	db := MakeTestDB()
	ms := MakeTestStore(db, kvStoreKeys, transientStoreKeys)
	ctx := sdk.NewContext(
		ms, abci.Header{ChainID: chainID, Time: MakeTestTimestampCurrent()}, isCheckTx, log.NewNopLogger(),
	)
	return ctx
}
