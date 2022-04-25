package keeper

import (
	"github.com/gogo/protobuf/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"

	v231 "github.com/certikfoundation/shentu/v2/x/shield/migrations/v231"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper      Keeper
	queryServer grpc.Server
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, queryServer grpc.Server) Migrator {
	return Migrator{keeper: keeper, queryServer: queryServer}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v231.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.keeper.paramSpace, m.queryServer)
}
