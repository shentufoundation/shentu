package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) CreateTheorem(ctx context.Context, proposer sdk.AccAddress, title, desc, code string, submitTime, endTime time.Time) (types.Theorem, error) {
	theoremID, err := k.TheoremID.Next(ctx)
	if err != nil {
		return types.Theorem{}, err
	}

	theorem, err := types.NewTheorem(theoremID, proposer, title, desc, code, submitTime, endTime)
	if err != nil {
		return types.Theorem{}, err
	}

	if err = k.Theorems.Set(ctx, theorem.Id, theorem); err != nil {
		return types.Theorem{}, err
	}

	if err := k.ActiveTheoremsQueue.Set(ctx, collections.Join(endTime, theoremID), theoremID); err != nil {
		return types.Theorem{}, err
	}
	return theorem, nil
}

func (k Keeper) DeleteTheorem(ctx context.Context, theoremID uint64) error {
	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		return err
	}

	err = k.ActiveTheoremsQueue.Remove(ctx, collections.Join(*theorem.EndTime, theorem.Id))
	if err != nil {
		return err
	}
	err = k.Theorems.Remove(ctx, theorem.Id)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) SetTheorem(ctx context.Context, theorem types.Theorem) error {
	return k.Theorems.Set(ctx, theorem.Id, theorem)
}
