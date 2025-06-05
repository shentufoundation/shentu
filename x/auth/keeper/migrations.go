package keeper

import (
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	v5 "github.com/cosmos/cosmos-sdk/x/auth/migrations/v5"
	"github.com/cosmos/gogoproto/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	v3 "github.com/cosmos/cosmos-sdk/x/auth/migrations/v3"
	v4 "github.com/cosmos/cosmos-sdk/x/auth/migrations/v4"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	authkeeper.AccountKeeper
	keeper         Keeper
	queryServer    grpc.Server
	legacySubspace exported.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(am authkeeper.AccountKeeper, keeper Keeper, queryServer grpc.Server, ss exported.Subspace) Migrator {
	return Migrator{
		AccountKeeper:  am,
		keeper:         keeper,
		queryServer:    queryServer,
		legacySubspace: ss,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(_ sdk.Context) error {
	//var iterErr error
	//
	//m.keeper.ak.IterateAccounts(ctx, func(account sdktypes.AccountI) (stop bool) {
	//	mvacc, ok := account.(*types.ManualVestingAccount)
	//	if !ok {
	//		return false
	//	}
	//	vestedCoins := mvacc.VestedCoins
	//
	//	dvAcc, err := vestingtypes.NewDelayedVestingAccount(
	//		mvacc.BaseAccount, mvacc.OriginalVesting, math.MaxInt64)
	//
	//	wb, err := v2.MigrateAccount(ctx, dvAcc, m.queryServer)
	//	if err != nil {
	//		iterErr = err
	//		return true
	//	}
	//
	//	if wb == nil {
	//		return false
	//	}
	//
	//	dvAcc, ok = wb.(*vestingtypes.DelayedVestingAccount)
	//	if !ok {
	//		return false
	//	}
	//	unlocker, err := sdk.AccAddressFromBech32(mvacc.Unlocker)
	//	if err != nil {
	//		panic(err)
	//	}
	//	newmvacc := types.NewManualVestingAccount(dvAcc.BaseAccount, dvAcc.OriginalVesting, vestedCoins, unlocker)
	//
	//	m.keeper.ak.SetAccount(ctx, newmvacc)
	//	return false
	//})
	//
	//return iterErr
	return nil
}

// Migrate2to3 migrates from consensus version 2 to version 3. Specifically, for each account
// we index the account's ID to their address.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc)
}

// Migrate3to4 migrates the x/auth module state from the consensus version 3 to
// version 4. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/auth
// module state.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.Migrate(ctx, m.keeper.storeService, m.legacySubspace, m.keeper.cdc)
}

// Migrate4To5 migrates the x/auth module state from the consensus version 4 to 5.
// It migrates the GlobalAccountNumber from being a protobuf defined value to a
// big-endian encoded uint64, it also migrates it to use a more canonical prefix.
func (m Migrator) Migrate4To5(ctx sdk.Context) error {
	return v5.Migrate(ctx, m.keeper.storeService, m.AccountKeeper.AccountNumber)
}
