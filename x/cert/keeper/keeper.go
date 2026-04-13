// Package keeper specifies the keeper for the cert module.
package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// Keeper manages certifier & certificate related logics.
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	authority    string

	Schema            collections.Schema
	Certifiers        collections.Map[sdk.AccAddress, types.Certifier]
	Certificates      collections.Map[uint64, types.Certificate]
	NextCertificateID collections.Sequence
}

// NewKeeper creates a new instance of the certifier keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		Certifiers: collections.NewMap(
			sb,
			collections.NewPrefix([]byte{types.CertifierStoreKeyPrefix}),
			"certifiers",
			sdk.AccAddressKey,
			codec.CollValue[types.Certifier](cdc),
		),
		Certificates: collections.NewMap(
			sb,
			collections.NewPrefix([]byte{types.CertificateStoreKeyPrefix}),
			"certificates",
			collections.Uint64Key,
			codec.CollValue[types.Certificate](cdc),
		),
		NextCertificateID: collections.NewSequence(
			sb,
			collections.NewPrefix([]byte{types.NextCertificateIDKeyPrefix}),
			"next_certificate_id",
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetCodec returns the keeper's binary codec.
func (k Keeper) GetCodec() codec.BinaryCodec { return k.cdc }
