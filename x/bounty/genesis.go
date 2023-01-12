package bounty

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, data *types.GenesisState) {
	k.SetNextProgramID(ctx, data.StartingProgramId)

	// check if the deposits pool account exists
	moduleAcc := k.GetBountyAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	var totalDeposits sdk.Coins
	for _, program := range data.Programs {
		k.SetProgram(ctx, program)
		totalDeposits = totalDeposits.Add(program.Deposit...)
	}

	// if account has zero balance it probably means it's not set, so we set it
	balance := bk.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balance.IsZero() {
		ak.SetModuleAccount(ctx, moduleAcc)
	}

	// check if total deposits equals balance, if it doesn't panic because there were export/import errors
	if !balance.IsEqual(totalDeposits) {
		panic(fmt.Sprintf("expected module account was %s but we got %s", balance.String(), totalDeposits.String()))
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	startingProgramID, _ := k.GetProgramID(ctx)
	programs := k.GetPrograms(ctx)

	return &types.GenesisState{
		StartingProgramId: startingProgramID,
		Programs:          programs,
	}
}
