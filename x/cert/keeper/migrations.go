package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v3 "github.com/shentufoundation/shentu/v2/x/cert/legacy/v3"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2 (no-op; address prefix migration was already applied).
func (m Migrator) Migrate1to2(_ sdk.Context) error {
	return nil
}

// Migrate2to3 migrates from version 2 to 3. The migration body lives in
// x/cert/legacy/v3; this method just wires keeper collections writes into
// that package via closures.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	return v3.MigrateStore(v3.StoreWriter{
		Cdc:   m.keeper.cdc,
		Store: store,
		WriteIndexes: func(cert types.Certificate) error {
			return m.keeper.writeCertificateIndexes(ctx, cert)
		},
		SetCertifier: func(c types.Certifier) error {
			return m.keeper.SetCertifier(ctx, c)
		},
		SetCertificate: func(id uint64, cert types.Certificate) error {
			return m.keeper.Certificates.Set(ctx, id, cert)
		},
		SetNextCertID: func(id uint64) error {
			return m.keeper.NextCertificateID.Set(ctx, id)
		},
	})
}
