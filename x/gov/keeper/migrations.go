package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v043 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v043"
	v046 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v046"

	v220 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v220"
	v260 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v260"
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
	err := v220.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
	if err != nil {
		return err
	}

	return v043.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {

	err := v260.MigrateParams(ctx, m.keeper.paramSpace)
	if err != nil {
		return err
	}
	return v260.MigrateProposalStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v046.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
