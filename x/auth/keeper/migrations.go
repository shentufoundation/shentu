package keeper

import (
	"math"

	"github.com/gogo/protobuf/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v043 "github.com/cosmos/cosmos-sdk/x/auth/legacy/v043"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/certikfoundation/shentu/v2/x/auth/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper      types.AccountKeeper
	queryServer grpc.Server
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, queryServer grpc.Server) Migrator {
	return Migrator{keeper: keeper.ak, queryServer: queryServer}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	var iterErr error

	m.keeper.IterateAccounts(ctx, func(account sdktypes.AccountI) (stop bool) {
		mvacc, ok := account.(*types.ManualVestingAccount)
		if !ok {
			return false
		}

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
		newmvacc := types.NewManualVestingAccount(dvAcc.BaseAccount, dvAcc.OriginalVesting, dvAcc.OriginalVesting, unlocker)

		m.keeper.SetAccount(ctx, newmvacc)
		return false
	})

	return iterErr
}
