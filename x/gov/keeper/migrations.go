package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/exported"
	v2 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v2"
	v3 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v3"
	sdkv4 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v4"

	v220 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v220"
	v260 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v260"
	v4 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v4"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.ParamSubspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace exported.ParamSubspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	err := v220.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
	if err != nil {
		return err
	}

	return v2.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {

	//err := v260.MigrateParams(ctx, m.legacySubspace)
	//if err != nil {
	//	return err
	//}
	return v260.MigrateProposalStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate3to4 migrates from version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate4to5 migrates from version 4 to 5.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {

	err := v4.MigrateCustomParams(ctx, m.keeper.storeKey, m.legacySubspace, m.keeper.cdc)
	if err != nil {
		return err
	}
	return sdkv4.MigrateStore(ctx, m.keeper.storeKey, m.legacySubspace, m.keeper.cdc)
}
