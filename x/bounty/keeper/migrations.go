package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v1 "github.com/shentufoundation/shentu/v2/x/bounty/migrations/v1"
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
