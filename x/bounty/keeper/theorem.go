package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) CreateTheorem(ctx context.Context, proposer sdk.AccAddress, title, desc, code string, submitTime, grantEndTime time.Time, proofTime time.Duration) (types.Theorem, error) {
	theoremID, err := k.TheoremID.Next(ctx)
	if err != nil {
		return types.Theorem{}, err
	}

	theorem, err := types.NewTheorem(theoremID, proposer, title, desc, code, submitTime, grantEndTime, proofTime)
	if err != nil {
		return types.Theorem{}, err
	}

	if err = k.Theorems.Set(ctx, theorem.Id, theorem); err != nil {
		return types.Theorem{}, err
	}

	if err := k.ActiveTheoremsQueue.Set(ctx, collections.Join(submitTime.Add(proofTime), theoremID), theoremID); err != nil {
		return types.Theorem{}, err
	}
	return theorem, nil
}

func (k Keeper) SetTheorem(ctx context.Context, theorem types.Theorem) error {
	return k.Theorems.Set(ctx, theorem.Id, theorem)
}
