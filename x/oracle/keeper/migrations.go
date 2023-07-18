package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v280 "github.com/shentufoundation/shentu/v2/x/oracle/legacy/v280"

	v2 "github.com/shentufoundation/shentu/v2/x/oracle/legacy/v2"
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
	v2.UpdateParams(ctx, m.keeper.paramSpace)
	return v2.MigrateTaskStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	v280.MigrateAllTaskStore(ctx, m.keeper.storeKey, m.keeper.cdc)
	return v280.MigrateOperatorStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
