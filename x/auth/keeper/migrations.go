package keeper

import (
	"math"

	"github.com/gogo/protobuf/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v043 "github.com/cosmos/cosmos-sdk/x/auth/migrations/v043"
	v046 "github.com/cosmos/cosmos-sdk/x/auth/migrations/v046"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/shentufoundation/shentu/v2/x/auth/types"
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
	var iterErr error

	m.keeper.ak.IterateAccounts(ctx, func(account sdktypes.AccountI) (stop bool) {
		mvacc, ok := account.(*types.ManualVestingAccount)
		if !ok {
			return false
		}
		vestedCoins := mvacc.VestedCoins

		dvAcc := vestingtypes.NewDelayedVestingAccount(
			mvacc.BaseAccount, mvacc.OriginalVesting, math.MaxInt64)

		wb, err := v043.MigrateAccount(ctx, dvAcc, m.queryServer)
		if err != nil {
			iterErr = err
			return true
		}

		if wb == nil {
			return false
		}

		dvAcc, ok = wb.(*vestingtypes.DelayedVestingAccount)
		if !ok {
			return false
		}
		unlocker, err := sdk.AccAddressFromBech32(mvacc.Unlocker)
		if err != nil {
			panic(err)
		}
		newmvacc := types.NewManualVestingAccount(dvAcc.BaseAccount, dvAcc.OriginalVesting, vestedCoins, unlocker)

		m.keeper.ak.SetAccount(ctx, newmvacc)
		return false
	})

	return iterErr
}

// Migrate2to3 migrates from consensus version 2 to version 3. Specifically, for each account
// we index the account's ID to their address.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v046.MigrateStore(ctx, m.keeper.key, m.keeper.cdc)
}
