package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v1 "github.com/shentufoundation/shentu/v2/x/bounty/migrations/v1"
	v2 "github.com/shentufoundation/shentu/v2/x/bounty/migrations/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v1.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc)
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc)
}
